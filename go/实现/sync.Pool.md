---
tags:
  - go
  - cache
date created: 2023-04-13 00:57
date modified: 2023-04-13 02:05
version: 1.2
---

```go
type Pool struct {
	noCopy noCopy

	local     unsafe.Pointer // 每个P的本地队列，实际类型为 [P]poolLocal
	localSize uintptr        // [P]poolLocal 大小

	victim     unsafe.Pointer 
	victimSize uintptr        

	// 自定义对象创建回调函数，当 pool 中无可用对象时会调用此函数
	New func() any
}
```

在新一轮 GC 到来时，victim 和 victimSize 会分别接管 local 和 localsize，victim 机制减少 GC 后冷启动导致的性能抖动，让对象分配更平滑
>[!tip] victim cache
>victim Cache 是计算机架构的一个概念。
>
>Victim Cache 是一种高速缓存技术，用于减少由于缓存缺失导致的处理器停顿时间。当缓存缺失时，处理器需要从主存储器中检索缺失的数据，这将导致较长的停顿时间。使用 Victim Cache 技术可以减少这种停顿时间。Victim Cache 是一个小而快速的缓存，通常比 L1 高速缓存更小，但比 L2 缓存更快。当 L1 高速缓存缺失时，数据可以从 Victim Cache 中检索，从而避免了从主存储器中检索数据所需的较长停顿时间。 Victim Cache 中的缓存行通常包含最近被淘汰的缓存行的数据，因此它们通常被称为“牺牲品缓存”。

```go
type poolLocalInternal struct {
	private any       // P 的私有缓存区，访问时不需要加锁
	shared  poolChain // 本地P 可以 pushHead/popHead; 其他 P 只能 popTail.
}

type poolLocal struct {
	poolLocalInternal

	// 在大多数平台上，128 mod (cache line size) == 0 可防止伪共享
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```

>[!tip] 伪共享
>伪共享（False sharing）指的是在多处理器系统中，由于不同的处理器缓存会缓存相同的缓存行（Cache line），而这个缓存行中只有部分数据被不同的处理器使用，因此可能会导致性能的降低。
>
当一个处理器对缓存行进行更新时，由于缓存一致性协议，其他处理器的缓存也需要更新这个缓存行中的数据，即使这些处理器并不需要更新这个数据。这个过程涉及到缓存一致性协议的通信，会消耗额外的带宽和处理器时间。
>
伪共享问题通常是由于两个或更多的变量被放置在同一个缓存行中而引起的，当其中一个变量被一个处理器更新时，就会导致其他处理器需要更新缓存行，即使其他变量在其他处理器中是不相关的。
>
为了避免伪共享问题，可以使用填充（padding）的方式来确保变量被放置在不同的缓存行中，以避免处理器之间不必要的缓存一致性通信，提高程序的性能。

^161773

```go
func (p *Pool) pinSlow() (*poolLocal, int) {
	.
	.
	.
	if p.local == nil {
		allPools = append(allPools, p)
	}
	.
	.
	.
}
```

```go
func (p *Pool) Get() any {
	.
	.
	.
	l, pid := p.pin()
	x := l.private
	l.private = nil
	if x == nil {
		// Try to pop the head of the local shard. We prefer
		// the head over the tail for temporal locality of
		// reuse.
		x, _ = l.shared.popHead()
		if x == nil {
			x = p.getSlow(pid)
		}
	}
	.
	.
	.
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}
```

```go
func (p *Pool) getSlow(pid int) any {
	.
	.
	.

	// Try the victim cache. We do this after attempting to steal
	// from all primary caches because we want objects in the
	// victim cache to age out if at all possible.
	size = atomic.LoadUintptr(&p.victimSize)
	if uintptr(pid) >= size {
		return nil
	}
	locals = p.victim
	l := indexLocal(locals, pid)
	if x := l.private; x != nil {
		l.private = nil
		return x
	}
	
	.
	.
	.
	return nil
}
```

```go
func poolCleanup() {
	// This function is called with the world stopped, at the beginning of a garbage collection.
	// It must not allocate and probably should not call any runtime functions.

	// Because the world is stopped, no pool user can be in a
	// pinned section (in effect, this has all Ps pinned).

	// Drop victim caches from all pools.
	for _, p := range oldPools {
		p.victim = nil
		p.victimSize = 0
	}

	// Move primary cache to victim cache.
	for _, p := range allPools {
		p.victim = p.local
		p.victimSize = p.localSize
		p.local = nil
		p.localSize = 0
	}

	// The pools with non-empty primary caches now have non-empty
	// victim caches and no pools have primary caches.
	oldPools, allPools = allPools, nil
}
```

- 当发生 GC 时，会调用 poolCleanup 主要是将 local 和 victim 做交换，不让 GC 一次就将 pool 清空
- 当调用 Get 时，因为 GC 后 local 变为 nil 所以从 victim 中获取对象
- 当调用 Put 时，local 变为 nil 后，经过 `pinSlow` 会将 p 追加进 local（这里就避免了一次 GC，暂时没有被调用的对象会被 GC 掉）

通过使用 victim cache 技术，使得 1 次 GC 的开销被拉长到 2 次，减小了 GC 压力，提高了对象复用
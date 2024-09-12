---
tags: []
date created: 2023-04-12 23:24
date modified: 2023-04-13 00:38
go version: 1.20
---

```go
type WaitGroup struct {
	noCopy noCopy

	state atomic.Uint64 // 高32位是运行计数，低32位是等待计数
	sema  uint32 // 信号计数
}
```

## wait

```go
func (wg *WaitGroup) Wait() {
	if race.Enabled {
		race.Disable()
	}
	for {
		state := wg.state.Load()
		v := int32(state >> 32) // 运行计数
		w := uint32(state) // 等待计数
		if v == 0 {
			// Counter is 0, no need to wait.
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
		// Increment waiters count.
		if wg.state.CompareAndSwap(state, state+1) {
			runtime_Semacquire(&wg.sema)
			if wg.state.Load() != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
	}
}
```

比较运行计数，每次完成检查后增加等待计数，为了其他 goroutine 能得到充分调度，调用 `runtime_Semacquire`,等待计数记录成功后直接休眠

## Add

```go
func (wg *WaitGroup) Add(delta int) {
	state := wg.state.Add(uint64(delta) << 32)
	v := int32(state >> 32)
	w := uint32(state)
	
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if wg.state.Load() != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	wg.state.Store(0)
	for ; w != 0; w-- {
		runtime_Semrelease(&wg.sema, false, 0)
	}
}
```

- 运行计数不能为负
- `运行计数>0` 或者 `等待计数==0` 则返回
- 通过信号计数通知所有等待的 goroutine

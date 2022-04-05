## Mutex
锁易错场景：
- Lock/Unlock 不是成对出现
- copy 已使用的mutex
> mutex是一个有状态的对象，state字段记录了锁的状态，如果是复制的，则会把state的状态也一并复制过去
- 重入
> mutex 没有记录哪个goroutine拥有这把锁，任何goroutine都可以随意的Unlock这把锁，没有办法计算重入条件
- 死锁


Mutex 结构
```go
type Mutex struct {
	state int32
	sema  uint32
}
```
state 字段分为四个含义
![[Pasted image 20220331180811.png]]

```go
// sync/mutex.go

func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}


func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false // goroutine的饥饿标记
	awoke := false // 唤醒标记
	iter := 0 // 自旋次数
	old := m.state // 当前锁的状态
	for {
		// 锁是非饥饿状态，锁还没有被释放，尝试自旋
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state // 再次获取锁的状态
			continue
		}
		new := old
		// 非饥饿状态加锁
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		// waiter 数量加1
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		
		// 如果是饥饿状态并且还是被锁住，则设置饥饿标记
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		
		if awoke {
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken // 清除唤醒标记
		}
		
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// 锁被释放且不是饥饿状态，正常获取锁
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			
			// 如果以前就在队列，加入到队列头
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			// 阻塞等待
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			// 唤醒之后检查是否应该处于饥饿状态
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			// 如果是饥饿状态直接获取锁
			if old&mutexStarving != 0 {
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				// waiter数量减1
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				// 不饥饿或者是最后一个waiter
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
}
```

## RWMutex
rwmutex 结构
```go
type RWMutex struct {
	w           Mutex  // 互斥锁
	writerSem   uint32 // writer信号量
	readerSem   uint32 // reader信号量
	readerCount int32  // reader数量
	readerWait  int32  // writer等待reader完成的数量
}
```
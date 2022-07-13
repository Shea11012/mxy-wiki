---
date created: 2022-02-28 19:47
date modified: 2022-03-28 13:24
title: GPM调度器
---
## Goroutine & Scheduler
### 操作系统三种线程模型：
#### 内核级线程模型
此模型下线程的切换调度由系统内核完成，系统内核负责将多个线程执行的任务映射到各个 cpu 中执行，直接使用操作系统内核来创建、销毁、切换和调度，对性能影响很大。
#### 用户级线程模型
多个用户线程一般从属于单个进程且多线程的调度是由用户自己完成的。这种实现方式相比内核线程很轻量，对系统资源的消耗很小，上下文切换代价也很小。但是并不能真正意义上做到并发，如果进程内的某个线程发生阻塞调用从而被 CPU 调度，会导致整个进程被挂起。
#### 混合线程模型
一个进程可以与多个内核线程 KSE（kernel schedule entity）关联。进程里的线程并不与 KSE 唯一绑定，而是多个用户线程映射到同一个 KSE，当某个 KSE 因为其绑定的线程阻塞操作被内核调度出 CPU 时，其关联的进程中其余用户线程可以重新与其它 KSE 绑定运行。
goroutine 和 go scheduler 在底层实现上是属于混合线程模型

## GPM 模型
每个 OS 线程的大小根据操作系统的不同大约在 1~8MB，这对于 goroutine 来说很大。
在 Go 中，每一个 goroutine 的初始大小为 2kb，采取动态扩容的方式，随着任务执行按需增长，且完全由 Go Scheduler 调度。
### G
表示 goroutine，每个 goroutine 对应一个 G 结构体，G 存储的 Goroutine 的运行堆栈、状态以及任务函数，可重用。G 并非执行体，每个 G 需要绑定到 P 才能调度执行。
### P
表示逻辑处理器（Processor），对 G 来说，P 相当于 CPU，G 只有绑定到 P（P 的 local runq 中）才能被调度。对 M 来说，P 提供了相关的执行环境。P 的数量决定了系统内最大可并行的 G 数量，P 的数量由 GOMAXPROCS 决定，但 P 的最大数量为 256.

### M
OS 线程抽象，表示真正的执行计算资源，在绑定有效的 P 后，进入 schedule 循环。
schedule 循环机制大致是从 Global 队列、P 的 local 队列以及 wait 队列中获取 G，切换到 G 的执行栈上并执行 G 的函数，调用 goexit 做清理工作并回到 M，如此反复。M 并不保留 G 状态，M 的数量是不确定的，由 Go runtime 调整，M 的最大数量为 10000 个。

![[Pasted image 20220326150634.png]]

## 调度时机
### 使用 go 关键字
使用关键字 go，一旦一个新的协程创建，会给调度器一个机会去执行调度

### GC
因为 GC 时会运行一组自己的 goroutines，这些 goroutines 需要时间去 M 上运行。这时调度器会做出调度

### system calls
如果一个 goroutine 发生系统调用将导致 M 阻塞，此时调度器会将 goroutine 解绑换一个新的 goroutine 到 M 上执行

### 同步和编排
如果一个 atomic、mutex、channel 导致 goroutine 阻塞，调度器会切换一个新的 goroutine 去运行。一旦可运行了就会重新入队等待被执行。

### 系统监控
go 程序启动时会有一个 sysmon 线程进行监控（定期垃圾回收和调度）


## go 进程启动
go 中的入口函数在 `asm_amd64.s` 中定义，main 包中的 main 函数是由 `runtime.main` 启动的
runtime/asm_amd64.s
```s
// runtime·rt0_go
.
.
.
get_tls(BX)  
LEAQ runtime·g0(SB), CX  
MOVQ CX, g(BX)  
LEAQ runtime·m0(SB), AX  

// 将m0和g0绑定到全局变量上
// save m->g0 = g0  
MOVQ CX, m_g0(AX)  
// save m0 to g0->m  
MOVQ AX, g_m(CX)

.
.
.

CALL runtime·check(SB)  
  
MOVL 24(SP), AX // copy argc  
MOVL AX, 0(SP)  
MOVQ 32(SP), AX // copy argv  
MOVQ AX, 8(SP)  
CALL runtime·args(SB)
// os初始化
CALL runtime·osinit(SB)  
// 调度初始化
CALL runtime·schedinit(SB)  
  
// 创建一个新goroutine执行程序
// create a new goroutine to start program  
MOVQ $runtime·mainPC(SB), AX // entry  
PUSHQ AX  
CALL runtime·newproc(SB)  
POPQ AX  
  
// 启动线程
// start this M  
CALL runtime·mstart(SB)  
  
CALL runtime·abort(SB)  // mstart should never return  
RET

```
主要分为四步：
1. 调用 `runtime.osinit` 获取系统的 CPU 个数
2. 调用 `runtime.schedinit` 初始化调度系统，进行 p 的初始化
3. `runtime.newproc` 新建一个 goroutine，参数 fn 为 `runtime.main`,建好后插入 m0 绑定的 p 的本地队列中
4. `runtime.mstart` 启动 m，进入调度系统

### mstart
`runtime.mstart` 由两个子函数组成：
- `runtime.mstart0`: 初始化栈边界
- `runtime.mstart1`: 初始化信号，并执行 `runtime.main` 

### getg 函数
getg 返回只想当前的 g 的指针，由汇编指令实现（可能来自 TLS 或专用寄存器）。
获取当前用户堆栈的 g，使用 `getg().m.curg`。当在系统或信号堆栈上执行时，则分别返回当前 m 的 g0 或 gsignal。要确定 g 是在用户堆栈或是系统堆栈上，可以使用 `getg() == getg().m.curg`,相等表示在用户态堆栈，不想等表示在系统堆栈。

### g0 和 m0
#### m0
- m0 表示进程启动的第一个线程，它是通过汇编直接复制给 m0
- 其他的 m 都是 runtime 内创建的
- 一个 go 进程只有一个 m0

#### g0
每个 m 都有一个 g0，因为每个线程都有一个系统堆栈，g0 与其它 g 的区别在栈的差别。g0 是系统分配的栈，在 Linux 上默认大小为 8MB，不能扩展也不能缩小。而其余的 g 一开始大小只有 2KB。g0 上没有任务函数也没有状态，并且不能被调度程序抢占。因为调度就是运行在 g0 上的

```go
// runtime/proc.go
var (
	m0           m  
	g0           g
)
```
全局变量上的 m0 代表主线程，g0 代表主线程的堆栈
### schedule 函数
```go
// runtime/proc.go

if gp == nil {  
   // Check the global runnable queue once in a while to ensure fairness.  
 // Otherwise two goroutines can completely occupy the local runqueue // by constantly respawning each other. 
	if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {  
	  lock(&sched.lock)  
	  gp = globrunqget(_g_.m.p.ptr(), 1)  
	  unlock(&sched.lock)  
   	}  
}  
if gp == nil {  
   gp, inheritTime = runqget(_g_.m.p.ptr())  
   // We can see gp != nil here even if the M is spinning,  
 // if checkTimers added a local goroutine via goready.}  
if gp == nil {  
   gp, inheritTime = findrunnable() // blocks until work is available  
}

execute(gp, inheritTime)
```
调度器查找 g 的流程如下：
1. 每隔 61 次调度从全局队列中找，避免全局队列 G 饿死
2. 从 p.runnext 获取 g，从 p 的本地队获取
3. 调用 `findrunnable` 找 g，找不到的话就将 m 休眠，等待唤醒

找到一个 g 后，就会调用 `execute` 去执行 g

### runtime.main 函数执行
```go
// runtime/proc.go
// 在系统栈上运行 sysmon
systemstack(func() {
		// 分配一个新的m，运行sysmom系统后台监控
		// 定期垃圾回收和调度抢占
		newm(sysmon, nil, -1)
})

// 确保运行在主线程上
if g.m != &m0 {
		throw("runtime.main not on m0")
}

// 执行runtime中的init函数
doInit(&runtime_inittask)

// 初始化完runtime的init后，启动gc的goroutine
gcenable()


// 执行编译器动态生成、用户自定义的所有init函数
doInit(&main_inittask)

// 开始执行 package main 下的main函数
fn := main_main // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime  
fn()
```

### G 的创建
```go
go func(){}()
```

g 一旦被创建就被保存进 `allgs` 全局变量中，且不会被销毁
g 的创建调用了 `runtime.newproc`,具体流程如下：
1. 用 `systemstack` 切换到系统堆栈，调用 `newproc1` 获取 g
2. 尝试从 p 的本地空闲链表和全局空闲链表中找到一个 g 的实例
3. 如果未找到，则调用 `malg` 生成新的 g 实例，且分配 g 的栈和设置栈边界，在添加到 `allgs` 中
4. 保存 g 切换的上下文
5. 生成唯一的 gid
6. 调用 `runqput` 将 g 插入队列中，如果本地队列还有剩余位置，将 G 插入本地队列，否则插入全局队列
7. 如果有空闲的 p 且 m 没有处于自旋且 `main goroutine` 已经启动，则唤醒或新建一个 m 来执行任务

```go
// runtime/proc.go
func newproc(fn *funcval) {  
   gp := getg()  
   pc := getcallerpc() 
	// 使用g0或 gsignal栈创建g对象
   systemstack(func() {  
      newg := newproc1(fn, gp, pc)  
  
      _p_ := getg().m.p.ptr()
	   // 将G入队
      runqput(_p_, newg, true)  
  
	  // 有空闲的g则唤醒或新建m执行
      if mainStarted {  
         wakep()  
      }  
   })  
}

func newproc1(fn *funcval, callergp *g, callerpc uintptr) *g {
	_g_ := getg()

	if fn == nil {
		_g_.m.throwing = -1 // do not dump full stacks
		throw("go of nil func value")
	}
	acquirem() // disable preemption because it can be holding p in a local var

	_p_ := _g_.m.p.ptr()
	// 从本地空闲链表和全局链表中找到一个g
	newg := gfget(_p_)
	if newg == nil {
		// 创建一个新的G，并加入全局变量 allgs 中
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead)
		allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
	}
	if newg.stack.hi == 0 {
		throw("newproc1: newg missing stack")
	}

	if readgstatus(newg) != _Gdead {
		throw("newproc1: new g is not Gdead")
	}

	// 参数大小+最小的栈尺寸
	totalSize := uintptr(4*goarch.PtrSize + sys.MinFrameSize) // extra space in case of reads slightly beyond frame
	totalSize = alignUp(totalSize, sys.StackAlign)
	// 新协程栈顶计算，将栈顶减去参数占用空间
	sp := newg.stack.hi - totalSize
	spArg := sp
	if usesLR {
		// caller's LR
		*(*uintptr)(unsafe.Pointer(sp)) = 0
		prepGoExitFrame(sp)
		spArg += sys.MinFrameSize
	}

	// 初始化g的gobuf，保存sp、pc、任务函数等
	memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
	newg.sched.sp = sp
	newg.stktopsp = sp
	// 保存goexit地址，后面会将goexit作为任务函数返回后执行的地址
	newg.sched.pc = abi.FuncPCABI0(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	// 保存当前新的G
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	// 将当前的pc压入栈，保存g的任务函数为pc
	gostartcallfn(&newg.sched, fn)
	// gopc保存newproc的pc
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	// 任务函数地址
	newg.startpc = fn.fn
	
	// 判断g的任务函数是否是runtime系统的任务函数
	if isSystemGoroutine(newg, false) {
		atomic.Xadd(&sched.ngsys, +1)
	} else {
		// Only user goroutines inherit pprof labels.
		if _g_.m.curg != nil {
			newg.labels = _g_.m.curg.labels
		}
	}
	// Track initial transition?
	newg.trackingSeq = uint8(fastrand())
	if newg.trackingSeq%gTrackingPeriod == 0 {
		newg.tracking = true
	}
	casgstatus(newg, _Gdead, _Grunnable)
	gcController.addScannableStack(_p_, int64(newg.stack.hi-newg.stack.lo))

	if _p_.goidcache == _p_.goidcacheend {
		// Sched.goidgen is the last allocated id,
		// this batch must be [sched.goidgen+1, sched.goidgen+GoidCacheBatch].
		// At startup sched.goidgen=0, so main goroutine receives goid=1.
		_p_.goidcache = atomic.Xadd64(&sched.goidgen, _GoidCacheBatch)
		_p_.goidcache -= _GoidCacheBatch - 1
		_p_.goidcacheend = _p_.goidcache + _GoidCacheBatch
	}
	// gid
	newg.goid = int64(_p_.goidcache)
	_p_.goidcache++
	if raceenabled {
		newg.racectx = racegostart(callerpc)
	}
	if trace.enabled {
		traceGoCreate(newg, newg.startpc)
	}
	releasem(_g_.m)

	return newg
}
```

## P 的创建
p 的初始化是在 `schedinit` 函数中调用的，最终调用 `procresize` 实现。P 的数量默认等于系统的 CPU 数。所有的 P 在程序启动时就已经创建完毕，并用一个全局变量 `allp` 维护
## M 的创建
通过 `newm` 来新建 M，最终通过 `newosproc` 函数来实现新建线程
```go
// runtime/proc.go

func newm(fn func(), _p_ *p, id int64) {
	acquirem()

	// 根据fn和p绑定一个m对象
	mp := allocm(_p_, fn, id)
	// 设置当前m的下一个p为_p_
	mp.nextp.set(_p_)
	mp.sigmask = initSigmask
	...
	// 分配os thread
	newm1(mp)
	releasem(getg().m)
}

func newm1(mp *m) {
	...
	execLock.rlock() // Prevent process clone.
	// 创建一个系统线程
	newosproc(mp)
	execLock.runlock()
}
```

M 是真正的执行者，它负责整个调度工作，调度的本质就是查找可以运行的 G，然后去执行 G 上的任务函数。可以理解为 M 循环执行着 `schedule` 函数，schedule 中 `findrunable` 函数执行流程：
1. 调用 runqget，尝试从 p 本地队列获取 G，获取到就返回
2. 调用 globalrunqget，尝试从全局队列获取 G，获取到就返回
3. 从网络轮询器中找到就绪的 G，把这个 G 变为可运行的 G 并返回
4. 如果不是所有的 P 都是空闲，随机选一个 P（最多四次），尝试从这个 P 中偷去一些 G，获取到就返回
5. 上面都找不到 G，判断此时 P 是否处于 GC mark 阶段，是则可以进行扫描和标记对象
6. 再次从全局队列中获取 G
7. 再次检查所有 P，查看有没有可运行的 G
8. 再次检查网络轮询器
9. 实在找不到可运行的 G，进入休眠 `stopm`

```go
// Finds a runnable goroutine to execute.
// Tries to steal from other P's, get g from local or global queue, poll network.
func findrunnable() (gp *g, inheritTime bool) {
	_g_ := getg()

	// The conditions here and in handoffp must agree: if
	// findrunnable would return a G to run, handoffp must start
	// an M.

top:
	_p_ := _g_.m.p.ptr()
	if sched.gcwaiting != 0 {
		gcstopm()
		goto top
	}
	if _p_.runSafePointFn != 0 {
		runSafePointFn()
	}

	now, pollUntil, _ := checkTimers(_p_, 0)

	if fingwait && fingwake {
		if gp := wakefing(); gp != nil {
			ready(gp, 0, true)
		}
	}
	if *cgo_yield != nil {
		asmcgocall(*cgo_yield, nil)
	}

	// local runq
	// 尝试从本地队列中获取G
	if gp, inheritTime := runqget(_p_); gp != nil {
		return gp, inheritTime
	}

	// global runq
	// 尝试从全局队列中获取G
	if sched.runqsize != 0 {
		lock(&sched.lock)
		gp := globrunqget(_p_, 0)
		unlock(&sched.lock)
		if gp != nil {
			return gp, false
		}
	}

	
	// 从网络轮询器中获取G，并将G变为可运行
	if netpollinited() && atomic.Load(&netpollWaiters) > 0 && atomic.Load64(&sched.lastpoll) != 0 {
		if list := netpoll(0); !list.empty() { // non-blocking
			gp := list.pop()
			injectglist(&list)
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false
		}
	}

	// Spinning Ms: steal work from other Ps.
	//
	// Limit the number of spinning Ms to half the number of busy Ps.
	// This is necessary to prevent excessive CPU consumption when
	// GOMAXPROCS>>1 but the program parallelism is low.
	procs := uint32(gomaxprocs)
	// 当前的M没有在自旋或者自旋的M数量小于繁忙P数量
	if _g_.m.spinning || 2*atomic.Load(&sched.nmspinning) < procs-atomic.Load(&sched.npidle) {
		if !_g_.m.spinning {
			_g_.m.spinning = true
			atomic.Xadd(&sched.nmspinning, 1)
		}

		// 从P中偷去一半的G，最多4次
		gp, inheritTime, tnow, w, newWork := stealWork(now)
		now = tnow
		if gp != nil {
			// Successfully stole.
			return gp, inheritTime
		}
		if newWork {
			// There may be new timer or GC work; restart to
			// discover.
			goto top
		}
		if w != 0 && (pollUntil == 0 || w < pollUntil) {
			// Earlier timer to wait for.
			pollUntil = w
		}
	}

	// We have nothing to do.
	//
	// If we're in the GC mark phase, can safely scan and blacken objects,
	// and have work to do, run idle-time marking rather than give up the
	// P.
	// 判断是否可以进行GC扫描和标记工作，并返回GC worker运行
	if gcBlackenEnabled != 0 && gcMarkWorkAvailable(_p_) {
		node := (*gcBgMarkWorkerNode)(gcBgMarkWorkerPool.pop())
		if node != nil {
			_p_.gcMarkWorkerMode = gcMarkWorkerIdleMode
			gp := node.gp.ptr()
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false
		}
	}

	...
	
	// 尝试从全局队列获取G
	if sched.runqsize != 0 {
		gp := globrunqget(_p_, 0)
		unlock(&sched.lock)
		return gp, false
	}
	// 将当前p和m解绑
	if releasep() != _p_ {
		throw("findrunnable: wrong p")
	}
	// 将p放入空闲链表
	pidleput(_p_)
	unlock(&sched.lock)

	
	// M 取消自旋
	wasSpinning := _g_.m.spinning
	if _g_.m.spinning {
		_g_.m.spinning = false
		if int32(atomic.Xadd(&sched.nmspinning, -1)) < 0 {
			throw("findrunnable: negative nmspinning")
		}

		// 检查所有的P，查看有没有可运行G
		_p_ = checkRunqsNoP(allpSnapshot, idlepMaskSnapshot)
		if _p_ != nil {
			acquirep(_p_)
			_g_.m.spinning = true
			atomic.Xadd(&sched.nmspinning, 1)
			goto top
		}

		// Check for idle-priority GC work again.
		// 检查GC work
		_p_, gp = checkIdleGCNoP()
		if _p_ != nil {
			acquirep(_p_)
			_g_.m.spinning = true
			atomic.Xadd(&sched.nmspinning, 1)

			// Run the idle worker.
			_p_.gcMarkWorkerMode = gcMarkWorkerIdleMode
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false
		}

		// Finally, check for timer creation or expiry concurrently with
		// transitioning from spinning to non-spinning.
		//
		// Note that we cannot use checkTimers here because it calls
		// adjusttimers which may need to allocate memory, and that isn't
		// allowed when we don't have an active P.
		pollUntil = checkTimersNoP(allpSnapshot, timerpMaskSnapshot, pollUntil)
	}

	// Poll network until next timer.
	// 再次检查网络轮询器
	if netpollinited() && (atomic.Load(&netpollWaiters) > 0 || pollUntil != 0) && atomic.Xchg64(&sched.lastpoll, 0) != 0 {
		atomic.Store64(&sched.pollUntil, uint64(pollUntil))
		if _g_.m.p != 0 {
			throw("findrunnable: netpoll with p")
		}
		if _g_.m.spinning {
			throw("findrunnable: netpoll with spinning")
		}
		delay := int64(-1)
		if pollUntil != 0 {
			if now == 0 {
				now = nanotime()
			}
			delay = pollUntil - now
			if delay < 0 {
				delay = 0
			}
		}
		if faketime != 0 {
			// When using fake time, just poll.
			delay = 0
		}
		list := netpoll(delay) // block until new work is available
		atomic.Store64(&sched.pollUntil, 0)
		atomic.Store64(&sched.lastpoll, uint64(nanotime()))
		if faketime != 0 && list.empty() {
			// Using fake time and nothing is ready; stop M.
			// When all M's stop, checkdead will call timejump.
			stopm()
			goto top
		}
		lock(&sched.lock)
		_p_ = pidleget()
		unlock(&sched.lock)
		if _p_ == nil {
			injectglist(&list)
		} else {
			acquirep(_p_)
			if !list.empty() {
				gp := list.pop()
				injectglist(&list)
				casgstatus(gp, _Gwaiting, _Grunnable)
				if trace.enabled {
					traceGoUnpark(gp, 0)
				}
				return gp, false
			}
			if wasSpinning {
				_g_.m.spinning = true
				atomic.Xadd(&sched.nmspinning, 1)
			}
			goto top
		}
	} else if pollUntil != 0 && netpollinited() {
		pollerPollUntil := int64(atomic.Load64(&sched.pollUntil))
		if pollerPollUntil == 0 || pollerPollUntil > pollUntil {
			netpollBreak()
		}
	}
	
	// 实在找不到进入休眠
	stopm()
	goto top
}
```

## 抢占的实现
sysmon 会先检测程序是否死锁

```go
// runtime/proc.go

func sysmon() {
	lock(&sched.lock)
	sched.nmsys++
	// 检测程序是否死锁
	checkdead()
	unlock(&sched.lock)

	lasttrace := int64(0)
	idle := 0 // how many cycles in succession we had not wokeup somebody
	delay := uint32(0)

	for {
		if idle == 0 { // start with 20us sleep...
			delay = 20
		} else if idle > 50 { // start doubling the sleep after 1ms...
			delay *= 2
		}
		if delay > 10*1000 { // up to 10ms
			delay = 10 * 1000
		}
		// 休眠 delay us
		usleep(delay)

		
		...
		// Update now in case we blocked on sysmonnote or spent a long time
		// blocked on schedlock or sysmonlock above.
		now = nanotime()

		// trigger libc interceptors if needed
		if *cgo_yield != nil {
			asmcgocall(*cgo_yield, nil)
		}
		// poll network if not polled for more than 10ms
		lastpoll := int64(atomic.Load64(&sched.lastpoll))
		// 如果超过了10ms都没进行netpool，则强制执行一次
		// 
		if netpollinited() && lastpoll != 0 && lastpoll+10*1000*1000 < now {
			atomic.Cas64(&sched.lastpoll, uint64(lastpoll), uint64(now))
			list := netpoll(0) // non-blocking - returns list of goroutines
			if !list.empty() {
				// Need to decrement number of idle locked M's
				// (pretending that one more is running) before injectglist.
				// Otherwise it can lead to the following situation:
				// injectglist grabs all P's but before it starts M's to run the P's,
				// another M returns from syscall, finishes running its G,
				// observes that there is no work to do and no other running M's
				// and reports deadlock.
				incidlelocked(-1)
				injectglist(&list)
				incidlelocked(1)
			}
		}
		...
		// retake P's blocked in syscalls
		// and preempt long running G's
		// 抢夺阻塞和长时间运行的G
		if retake(now) != 0 {
			idle = 0
		} else {
			idle++
		}
		...
	}
}
```

```go
// runtime/proc.go

func retake(now int64) uint32 {
	n := 0
	// Prevent allp slice changes. This lock will be completely
	// uncontended unless we're already stopping the world.
	lock(&allpLock)
	// We can't use a range loop over allp because we may
	// temporarily drop the allpLock. Hence, we need to re-fetch
	// allp each time around the loop.
	for i := 0; i < len(allp); i++ {
		_p_ := allp[i]
		if _p_ == nil {
			// This can happen if procresize has grown
			// allp but not yet created new Ps.
			continue
		}
		pd := &_p_.sysmontick
		s := _p_.status
		sysretake := false
		if s == _Prunning || s == _Psyscall {
			// Preempt G if it's running for too long.
			t := int64(_p_.schedtick)
			if int64(pd.schedtick) != t {
				pd.schedtick = uint32(t)
				pd.schedwhen = now
			} else if pd.schedwhen+forcePreemptNS <= now {
				preemptone(_p_)
				// In case of syscall, preemptone() doesn't
				// work, because there is no M wired to P.
				sysretake = true
			}
		}
		if s == _Psyscall {
			// Retake P from syscall if it's there for more than 1 sysmon tick (at least 20us).
			t := int64(_p_.syscalltick)
			if !sysretake && int64(pd.syscalltick) != t {
				pd.syscalltick = uint32(t)
				pd.syscallwhen = now
				continue
			}
			// On the one hand we don't want to retake Ps if there is no other work to do,
			// but on the other hand we want to retake them eventually
			// because they can prevent the sysmon thread from deep sleep.
			if runqempty(_p_) && atomic.Load(&sched.nmspinning)+atomic.Load(&sched.npidle) > 0 && pd.syscallwhen+10*1000*1000 > now {
				continue
			}
			// Drop allpLock so we can take sched.lock.
			unlock(&allpLock)
			// Need to decrement number of idle locked M's
			// (pretending that one more is running) before the CAS.
			// Otherwise the M from which we retake can exit the syscall,
			// increment nmidle and report deadlock.
			incidlelocked(-1)
			if atomic.Cas(&_p_.status, s, _Pidle) {
				if trace.enabled {
					traceGoSysBlock(_p_)
					traceProcStop(_p_)
				}
				n++
				_p_.syscalltick++
				handoffp(_p_)
			}
			incidlelocked(1)
			lock(&allpLock)
		}
	}
	unlock(&allpLock)
	return uint32(n)
}
```

`retake` 函数会遍历所有的 P，如果一个 P 处于 running 或 syscall 状态且运行时间超过了 10ms，调用 `preemptone` 设置 p 的 `gp.stackguard0=stackPreempt` 和 `gp.preempt = true` （导致该 P 正在执行的 G 进行下一次函数调用时栈空间检查失败，触发 morestack 然后进行一系列的系统调用，最终调用 `goschedImpl` 解除 P 与当前 M 的关联，让该 G 进入 \_Grunnable 状态，插入全局队列），处于阻塞状态调用 `handoff` 检查该 P 下是否有其他可运行的 G，如果有的话则调用 `startm` 来获取或新建一个 M 来服务。
```mermaid
graph LR
	morestack --> newstack -->gopreempt_m --> goschedImpl --> schedule
```

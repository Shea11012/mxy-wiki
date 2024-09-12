---
tags: []
date created: 2022-03-26 14:50
date modified: 2023-04-12 13:59
go version: 1.20
---

## 应用场景

**停止信号**
关闭 channel 或者向 channel 发送一个元素，使得接受方得知信息，进而做一些其他操作。

**定时任务**
与定时器结合，实现超时控制、实现定期执行某个任务。

**解耦生产和消费方**
启动 N 个 worker 协程，这些协程工作在一个 for 无限循环里，从 channel 消费工作任务

**控制并发数**
```go
token := make(chan struct{},3)

for _,w := range work {
	go func(){
		token <- struct{}
		w()
		<-token
	}()
}
```

## 结构

```go
type hchan struct {
	qcount   uint           // 元素个数
	dataqsiz uint           // 循环数组长度
	buf      unsafe.Pointer // 指向循环数组的指针
	elemsize uint16	// 元素大小
	closed   uint32 // chan 是否被关闭
	elemtype *_type // 元素类型
	sendx    uint   // 已发送索引
	recvx    uint   // 接受索引
	recvq    waitq  // 等待接收的goroutine队列
	sendq    waitq  // 等待发送的goroutine队列
	// 保护hchan中的所有字段
	lock mutex
}

type waitq struct {  
	first *sudog  
	last *sudog  
}
```

- buf：指向底层循环数组，只有缓冲型的 channel 才有效
- sendx，recvx：指向底层循环数组，表示当前可以发送和接收的元素位置索引
- senq,recvq：分别表示被阻塞的 goroutine，这些 goroutine 尝试向 channel 读取或发送数据而被阻塞。
- waitq：是 sudog 的一个双向链表，sudog 是对 goroutine 的封装
- lock：保证每个读或写都是原子操作

## 发送数据

- 当存在等待的接收者时，会直接将数据发送给阻塞的接收者
- 当缓冲区存在剩余空间时，将发送的数据写入 channel 的缓冲区
- 当不存在缓冲区或者缓冲区已满时，等待其他 goroutine 从 channel 接收数据
```go
// runtime/chan.go

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	...

	// 等待队列不为空时，直接发送
	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	// 缓冲区未满时，写入缓冲区
	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	// 阻塞发送，当channel中没有接收者能处理数据时，向channel发送数据会被阻塞
	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	closed := !mysg.success
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	if closed {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	return true
}
```

## 接收数据

- 如果 `ep==nil`，表示调用者忽略接收值
- 如果 `block==false`，表示非阻塞调用，在没有数据可接收的情况下返回 `false,false`
- c 是关闭状态，将 ep 指向的地址清零，返回 `true,false`，否则用返回值填充 ep 地址，返回 `true,true`


```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	if debugChan {
		print("chanrecv: chan=", c, "\n")
	}

	if c == nil {
		if !block {
			return
		}
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	if !block && empty(c) {
		if atomic.Load(&c.closed) == 0 {
			return
		}
		if empty(c) {
			// The channel is irreversibly closed and empty.
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			return true, false
		}
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	// channel已关闭且循环数组里没有元素，给ep赋值该类型的零值
	if c.closed != 0 && c.qcount == 0 {
		if raceenabled {
			raceacquire(c.raceaddr())
		}
		unlock(&c.lock)
		if ep != nil {
			typedmemclr(c.elemtype, ep)
		}
		return true, false
	}

	// 直接接收，当channel中的sendq不为空时，取出队头的goroutine
	if sg := c.sendq.dequeue(); sg != nil {
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}

	// 缓冲区，当channel的缓冲区已经包含数据时，接收数据会直接从缓冲区中 recx 的索引位置取出数据进行处理
	if c.qcount > 0 {
		// Receive directly from queue
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
		}
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemclr(c.elemtype, qp)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	if !block {
		unlock(&c.lock)
		return false, false
	}

	// 阻塞接受，当channel发送队列中不存在等待的goroutine并且缓冲区也不存在任何数据，从管道接受会阻塞
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	c.recvq.enqueue(mysg)
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)

	// someone woke us up
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	success := mysg.success
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, success
}
```

## 关闭 channel

```go
func closechan(c *hchan) {
	if c == nil {
		panic(plainError("close of nil channel"))
	}

	lock(&c.lock)
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}

	if raceenabled {
		callerpc := getcallerpc()
		racewritepc(c.raceaddr(), callerpc, abi.FuncPCABIInternal(closechan))
		racerelease(c.raceaddr())
	}

	c.closed = 1

	var glist gList

	// 将 channel 所有等待接收队列里的 sudog 释放
	for {
		sg := c.recvq.dequeue()
		if sg == nil {
			break
		}
		// elem 不为 nil，说明该 receiver 未忽略接收数据
		// 给它赋值一个类型零值
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		// 追加进链表
		glist.push(gp)
	}

	// 释放写，并将 goroutine 追加进 glist 中
	for {
		sg := c.sendq.dequeue()
		if sg == nil {
			break
		}
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}
	unlock(&c.lock)

	// 这里会唤醒所有的sudog，读会返回零值，但写关闭的channel会panic
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
```
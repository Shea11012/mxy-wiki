---
date created: 2021-11-30 21:22
date modified: 2022-03-26 14:49
title: channel
---
# channel

```go
var c chan string // 未初始化，为 nil 的 channel
<-c // 阻塞
c <- "a" // 阻塞

c = make(chan string)
close(c)
<-c // 关闭后的channel，可读
c <- "a" // panic
```

读取管道时，阻塞条件有：

- 管道无缓冲区
- 管道的缓冲区中无数据
- 管道的值为 nil

写入管道时，阻塞条件：

- 管道无缓冲区
- 管道的缓冲区已满
- 管道的值为 nil

### done channel

done channel 负责关闭创建的协程

```go
newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited. noleak")
			defer close(randStream)
      /* 	// 此处如果不适用 done 控制这个协程，则在守护进程会导致内存泄漏
      	for {
      		randStream <- rand.Int()	
      	}
      */
			for {
				select {
				case <-done:
					return
				case randStream <- rand.Int():

				}
			}
		}()

		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n",i,<-randStream)
	}
	close(done)
	time.Sleep(time.Second)
```

### or channel

or channel 在其中一个 channel 任务完成后其余 channel 都会结束

```go
package main

import (
	"fmt"
	"time"
)

var or func(channels ...<-chan interface{}) <-chan interface{}

func main() {
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {	// 此处是递归的终止条件
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {	// 这里每次监听三个任务，一旦任务完成这个协程就执行完毕，会将 orDone 关闭，并触发递归监听 ordone 的 select
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): // 此处发生递归

				}
			}
		}()

		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Second),
		sig(3*time.Hour),
		sig(4*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

```

### err handling

```go
package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error error
	Response *http.Response
}
func main() {
	checkStatus := func(done <-chan interface{},urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _,url := range urls {
				var result Result
				resp,err := http.Get(url)
				result = Result{Error: err,Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.baidu.com","https://badhost"}
	for result := range checkStatus(done,urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v",result.Error)
			continue
		}
		fmt.Printf("Response: %v\n",result.Response.Status)
	}
}

```

### pipeline 

pipeline 流式处理，会让数据变得有顺序，但可能会因为当前的流处理花费的时间过长导致上流被堵塞。

```go
func generator(done <-chan interface{},integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _,i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()

	return intStream
}

func multiplyChan(done <-chan interface{},intStream <-chan int,multipiler int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i*multipiler:
			}
		}
	}()
	return multipliedStream
}

func addChan(done <-chan interface{},intStream <-chan int,additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i+ additive:
			}
		}
	}()

	return addedStream
}

func ChanPipeline() {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done,1,2,3,4)
	pipeline := multiplyChan(done,addChan(done,multiplyChan(done,intStream,2),1),2)
	for range pipeline {

	}
}
```

### fan-out & fan-in

管道的属性可以让它的各个阶段都独立处理，可以多次重复使用管道的各个阶段。在多个 Goroutine 上重用管道的单个阶段实现并行化，可以提升管道的性能。

fan-out（扇出）：用于描述启动多个 goroutines 以处理来自管道的输入过程。

fan-in（扇入）：将多个结果组合到一个通道的过程。

在不依赖模块之前的计算结果，运行需要很长时间时，可以使用该模式。

```go
repeatFn := func(
		done <-chan interface{},
		fn func() (interface{}, ),
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()

		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:

				}
			}
		}()
		return takeStream
	}

	toInt := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()

		return intStream
	}

	primeFinder := func(
		done <-chan interface{},
		intStream <-chan int,
	) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer -= 1
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer % divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()

		return primeStream
	}

	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
		) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _,c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	rands := func() interface{} {
		return rand.Intn(50000000)
	}
	randIntStream := toInt(done,repeatFn(done,rands))

	// 正常的管道模式
	idx1 := 0
	for prime := range take(done,primeFinder(done,randIntStream),12) {
		idx1++
		fmt.Printf("\t%d: %d\n",idx1,prime)
	}
	fmt.Printf("管道模式耗费时间：%v",time.Since(start))

	 // fan-out fan-in 模式
	numFinders := runtime.NumCPU()
	fmt.Printf("找 %d prime\n",numFinders)
	finders := make([]<-chan interface{},numFinders) // 1 使用多个查找质数的模块做处理，体现为 fan-out
	fmt.Println("primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done,randIntStream)
	}

	idx := 0
	for prime := range take(done,fanIn(done,finders...),numFinders) { // 2 通过 fan-in 将多个管道整合成一个输出
		idx++
		fmt.Printf("\t%d: %d\n",idx,prime)
	}

	fmt.Printf("消耗时间：%v\n",time.Since(start))
```

### orDone

通过 done 控制可以快速的退出嵌套循环

```go
orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()

		return valStream
	}

	data := make(chan interface{})
	go func() {
		a := []int{1,2,3,4,5,6,7}
		for _,v:= range a {
			data <- v
		}
	}()

	done := make(chan interface{})
	for val := range orDone(done,data) {
		if val.(int) == 5 {
			close(done)
		}
		fmt.Println(val.(int))
	}
```

### tee channel

将从通道接收到的值，分别发送至两个独立的区域，如：从上一个通道接收到值，经过处理后，将值分别发送给执行者和记录操作日志。tee channel 的功能与 unix 的 tee 命令功能类似。

### bridge channel

将一系列通道拆解为一个简单的通道，称为桥接通道

### queue channel

队列不是减少了某个阶段的运行时间，而是减少了它处于阻塞状态的时间。队列将操作流程分离，以便一个阶段的运行时间不会影响另一个阶段的运行时间。

过早的添加队列会隐藏同步问题，如死锁和活锁。

Go 的 bufio 就是利了用这种特性
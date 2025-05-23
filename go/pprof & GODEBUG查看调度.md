---
tags: 
date created: 2021-11-30 21:22
date modified: 2023-10-31 16:40
title: pprof & GODEBUG查看调度
---

# PProf

采样方式：

- runtime/pprof：采集程序指定区块的数据进行分析
- net/http/pprof：基于 HTTP Server，且可以采集运行时的数据进行分析
- go test：基于测试用例，指定所需标识进行采集

使用模式：

- Report Generation： 报告生成
- Interactive Terminal User： 交互式终端
- Web Interface：web 界面

可分析内容：

- CPU Profiling：CPU 分析。按照一定的频率采集所监听的应用程序 CPU，确定应用程序在主动消耗 CPU 周期时花费时间的位置
- Memory Profiling：内存分析。在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏。
- Block Profiling：阻塞分析。记录 Goroutine 阻塞等待同步的位置，默认不开启，需要调用 **runtime.SetBlockProfileRate** 进行设置。
- Mutex Profiling：互斥锁分析。报告互斥锁的竞争情况，默认不开启，需要调用 **runtime.SetMutexProfileFraction** 设置。

可查看信息：

- allocs：查看过去所有内存分配样本
- block：查看导致阻塞同步的堆栈跟踪

- cmdline：当前程序命令行完整调用路径
- goroutine：查看当前所有运行的 goroutine 堆栈跟踪
- heap：查看活动对象的内存分配情况
- mutex：查看导致互斥锁的竞争持有者的堆栈跟踪

- profile：默认进行 30s 的 CPU profiling
- threadcreate：查看创建新 OS 线程的堆栈跟踪

## 使用方式

> 1. 通过浏览器访问
>
>    http://127.0.0.1:{port}/debug/pprof/ # 端口自定义
>
>    在对应的可查看信息后面加上 ?debug=1 可以直接在浏览器查看，如 http://127.0.0.1:6060/debug/pprof/allocs?debug=1
>
> 2. 交互式终端
>
>    go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=60
>
> 3. 下载 pprof 文件分析
>
>    使用 wget、curl 等命令下载文件进行分析
>
>    go tool pprof 下载的文件

## CPU profiling

> go tool pprof http://127.0.0.1:6060/debug/pprof/profile?second=60

参数含义：

- flat：函数运行耗时
- flat%：函数占 CPU 运行总耗时的比例
- sum%：函数累积占 CPU 运行总耗时比例
- cum：函数及调用函数占用 CPU 总耗时
- cum%：函数及调用函数占 CPU 运行总耗时的比例
- Name：函数名

> 在 pprof 终端内，执行 list 命令查看堆栈分析时，其默认的依赖时当前运行主机上得源码路径，也就是二进制文件的路径。从远程拉取的 profile，本机和原本运行主机的路径一般是不同的，此时会导致 list 命令不可用需。 go tool pprof 指定 source_path 参数

## Heap profiling

> go tool pprof http://127.0.0.1:6060/debug/pprof/heap
>
> // inuse_space 分析内存常驻内存占用情况
>
> go tool pprof -inuse_space http://127.0.0.1:6060/debug/pprof/heap
>
> // alloc_objects 分析应用程序的内存临时分配情况
>
> go tool pprof -alloc_objects http://127.0.0.1:6060/debug/pprof/heap
>
> // inuse_objects 查看每个函数对象数量
>
> // alloc_space 查看每个函数分配内存空间大小

## Goroutine Profiling

> go tool pprof http://127.0.0.1:6060/debug/pprof/goroutine
>
> 查看 goroutine 时可以使用 traces 命令，这个命令会打印处所有的调用栈，以及指标信息。
>
> 调用栈的展示顺序是自下而上的

## Mutex Profiling

> go tool pprof http://127.0.0.1:6060/debug/pprof/mutex
>
> **在分析 mutex 时，需要设置 runtime.SetMutexProfileFraction(1)**
>
> 使用 top 查看互斥量排名
>
> list func|address 查看具体的函数，以及锁开销的位置

## Block Profiling

> go tool pprof http://127.0.0.1:6060/debug/pprof/block
>
> 查看 block 分析时，需设置 **runtime.SetBlockProfileRate(1)**
>
> 查看方式等同 mutex

## 使用可视化界面查看

>先下载一个文件，使用 wget 或 curl
>
>go tool pprof -http=:6001 file|url
>
>// 如果遇到下面的提示，可以去官网下载 [graphviz](http://www.graphviz.org/download/)
>
>Could not execute dot; may need to install graphviz.

## 测试用例分析

通过测试用例可以跟精准的分析流程和函数

```go
func TestAdd(t *testing.T) {
    _ = Add("fhasddfa")
}

func BenchmarkAdd(b *testing.B) {
    for i:=0;i<b.N;i++ {
        Add("fasfasfas")
    }
}

// 执行命令
// go test -bech=. -cpuprofile=cpu.profile
// 执行完毕后，可以看到cpu.profile 文件
```

## 通过 Lookup 写入文件进行分析

pprof 提供了 lookup 方法对相关内容进行采集和调用，一共支持 6 中类型：

- goroutine
- threadcreate
- heap
- block
- mutex

```go
pprof.Lookup(name) // 此处的 name 就是上述类型
```

# Trace

> wget 或 curl 下载 http://127.0.0.1:6060/debug/pprof/trace?seconds=60 trace.out
>
> go tool trace trace.out

主要查看 goroutine 的创建、阻塞、解除阻塞，syscall 的进入、退出、阻止，GC 事件、调整 Heap 的大小，以及 Processor 启动、停止等。

可视界面列表：

- view trace：查看跟踪
- goroutine analysis：goroutine 分析
- network blocking profile：网络阻塞概况 (优先查看)
- synchronization blocking profile：同步阻塞概况
- syscall blocking profile：系统调用阻塞概况
- scheduler latency profile：调度延迟概况，(优先查看)
- user defined tasks：用户自定义任务
- user defined regions：用户自定义区域
- minimum mutator utilization：最低 mutator 利用率

## GODEBUG

### schedtrace

参数：

- schedtrace：设置 schedtrace=X，可以在运行时每 X 毫秒发出一行调度器的摘要信息到标准 err 输出中。
- scheddetail：设置 schedtrace=X 和 scheddetail=1，可以在运行时每 X 毫秒发出一次详细的多行信息，信息内容包括调度程序、处理器、OS 线程和 goroutine 的状态。

查看调度器或垃圾回收等详细信息

`GODEBUG=schedtrace=1000 go run main.go` 输出调度器的摘要信息

- sched：每一行都代表调度器的调试信息，后面提示的毫秒数表示从启动到现在的运行时间，输出时间受到 schedtrace 值影响

- gomaxprocs：当前的 CPU 核心数，受 GOMAXPROCS 控制
- idleprocs：空闲的处理器数量，表示当前空闲的处理器数量
- threads：OS 线程数量，表示当前正在运行的线程数量
- spinningthreads：自旋状态的 OS 线程数量
- idlethreads：空闲的线程数量
- runqueue：全局队列的 goroutine 数量，[0 0 1 1] 分别代表 4 个 P 的本地队列正在运行的 goroutine 数量
- runqsize：运行队列中的 G 数量
- gfreecnt：可用的 G，状态为 Gdead
- schedtick：P 的调度次数
- syscalltick：系统调用次数
- preemptoff：如果不等于空字符串，则保持 curg 在这个 M 上运行

增加 scheddetail 参数，可以提供处理器 P，线程 M 和 goroutine G 的细节

```go
GODEBUG=schedtrace=1000,scheddetail=1 go run main.go
```

摘要信息中，G 的 status 状态列表：

| 状态             | 值   | 含义                                                         |
| ---------------- | ---- | ------------------------------------------------------------ |
| Gidle            | 0    | 刚刚被分配，还没有进行初始化                                 |
| Grunnable        | 1    | 已经在运行队列中，还没有执行用户代码                         |
| Grunning         | 2    | 不在运行队列中，已经可以执行用户代码，此时已经分配了 M 和 P     |
| Gsyscall         | 3    | 正在执行系统调用，此时分配了 M                                |
| Gwaiting         | 4    | 在运行时被阻止，没有执行用户代码，也不再运行队列中，此时正在某处阻塞等待中 |
| Gmoribund_unused | 5    | 尚未使用，但是在 gdb 中进行了硬编码                            |
| Gdead            | 6    | 尚未使用，这个状态可能时刚退出或是刚被初始化，此时它没有执行用户代码，有可能是没有分配堆栈 |
| Genqueue_unused  | 7    | 尚未使用                                                     |
| Gcopystack       | 8    | 正在复制堆栈，并没有执行用户代码，也不在运行队列中           |

p 的状态

| 状态     | 值   | 含义                                                       |
| -------- | ---- | ---------------------------------------------------------- |
| Pidle    | 0    | 刚刚被分配，还没有进行初始化                               |
| Prunning | 1    | 当 M 与 P 绑定调用 acquirep 时，P 的状态会变为 Prunning          |
| Psyscall | 2    | 正在执行系统调用                                           |
| Pgcstop  | 3    | 暂停运行，此时系统正在进行 GC，直至 GC 结束后才会进入下一阶段 |
| Pdead    | 4    | 废弃，不再使用                                             |

### gctrace

`GODEBUG='gctrace=1' exec`
设置 `gctrace=1` 使垃圾回收器每次回收时汇总所回收内存大小及耗时，并输出

输出格式
```
gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
```

|             | 含义                               |
| ----------- | ---------------------------------- |
| gc #        | GC 次数编号，每次 GC 递增             |
| @\#s        | 距离程序开始执行时的时间           |
| \#%         | GC 占用执行时间百分比               |
| \#\+...\+\# | GC 使用时间                         |
| \#->#-># MB | GC 开始，结束，及当前活跃堆内存大小 |
| # MB goal   | 全局堆内存大小                     |
| # P         | processor 数量                      |
>[!tip]
>如果以 `(forced)` 结尾，则是由 `runtime.GC` 触发


>[!example]
> `gc 17 @0.149s 1%: 0.004+36+0.003 ms clock, 0.009+0/0.051/36+0.006 ms cpu, 181->181->101 MB, 182 MB goal, 2 P`
> 
> `gc 17`：GC 编号 17
> `@0.149s`：程序已经执行 0.149s
> `1%`：gc 占用 1% 时间
> `0.004+36+0.003 ms clock`：垃圾回收时间，分别为 STW 清扫时间 + 并发标记和扫描时间 +STW 标记时间
> `0.009+0/0.051/36+0.006 ms cpu`：垃圾回收占用 CPU 时间
> `181->181->101 MB`：GC 开始前堆内存 181M，GC 结束后堆内存 181M，当前活跃堆内存 101M
> `182 MB goal`：全局堆内存大小

## gops 进程诊断工具

gops (go process status)。gops 由 google 官方推出的一个命令行工具

```go
// go get -u github.com/google/gops

import "github.com/google/gops/agent"

func main() {
    if err := agent.Listen(agent.Option{});err != nil {
        log.Fatal(err)
    }
    
    http.HandleFunc("/hello",func(w http.ResponseWriter,r *http.Request) {
        _,_ = w.Write([]byte("gops test"))
    })
    
    _ = http.ListenAndServe(":6060",http.DefaultServeMux)
}

// 启动后可以通过 gops 查看详细信息
```


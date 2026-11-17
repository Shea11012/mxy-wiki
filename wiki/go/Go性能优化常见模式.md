---
title: "Go性能优化常见模式"
created: 2026-05-10
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - go
source: "https://goperf.dev/01-common-patterns/"
confidence: INFERRED
---
# Go 性能优化常见模式 

来源：[Go Optimization Guide](https://goperf.dev/01-common-patterns/)  
共 15 种模式，分 4 大类，每篇含：问题场景 → 核心原理 → 关键代码 → 性能数据 → 注意事项。

---

# 一、内存管理与效率（7 篇）

## 1.1 对象池化 (Object Pooling) — sync.Pool

### 问题场景
高吞吐服务（HTTP 处理、JSON 编解码、日志格式化）中每个请求都需要临时缓冲区（`bytes.Buffer`）。每次 `new(bytes.Buffer{})` = 一次堆分配 + GC 回收。QPS 数千时，分配抖动显著推高 GC 频率和 CPU。

### 核心原理
`sync.Pool` 是线程安全对象池：
- `Get()` 取出对象（池空时调 `New` 工厂函数创建）
- 使用完毕后 `Put()` 归还
- 关键：配合 `Reset()` 重置对象状态而不释放底层内存 —— `bytes.Buffer.Reset()` 只归零偏移，底层 `[]byte` 保留

### 关键代码
```go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handleRequest(payload []byte) {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()         // 只复位偏移，不释放底层数组
    buf.Write(payload)  // 复用已有容量，零分配
    fmt.Println(buf.Len())
    bufPool.Put(buf)
}
```

### 性能数据
写入 4KB 负载：无池化 864ns + 1 alloc + 4096B；有池化 42ns + 0 alloc + 0B。**吞吐提升约 20 倍**，完全消除逐次分配。

### 注意事项
- 仅适合短生命周期、可复用的对象。长期持有/跨 goroutine 共享不适合
- `Get` 返回的对象状态不确定 —— 必须 Reset
- 池中对象可能随时被 GC 回收，不要依赖持久存在
- 先用 pprof 确认分配热点，再引入 Pool

---

## 1.2 内存预分配 (Memory Preallocation)

### 问题场景
循环向 slice/map 追加元素时，Go 运行时自动扩容。slice 小容量翻倍增长（2→4→8...），超 1024 后增速降为约 25%；map 涉及 rehashing。已知容量时这些反复分配和拷贝可避免。

### 核心原理
`make([]T, 0, capacity)` 或 `make(map[K]V, capacity)` 提前告知容量。运行时一次性分配足够内存，后续 append 无需扩容。

### 关键代码
```go
// 差：反复扩容
var result []int
for i := 0; i < 10000; i++ {
    result = append(result, i)
}

// 好：一次分配
result := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    result = append(result, i)
}

// 更好：索引直接赋值（无边界检查）
result := make([]int, 10000)
for i := range result {
    result[i] = i
}

// map 预分配 (Go 1.11+)
m := make(map[int]string, 10000)
```

### 性能数据
追加 10000 个 int：无预分配 28539ns + 19 allocs + 357KB；有预分配 7093ns + 1 alloc + 82KB。**速度提升约 4 倍**，分配从 19 降到 1。

### 注意事项
- 容量预估不准时：预分配太少仍扩容，太多浪费内存
- 数据量波动大（100~10万），固定预分配 10 万可能严重浪费
- 先 profiling，再优化

---

## 1.3 结构体字段对齐 (Struct Field Alignment)

### 问题场景
Go 编译器按平台对齐规则在字段间插入 padding。排列不当会显著膨胀结构体。例如 64 位：`bool→int64→byte` = 24 字节；`int64→bool→byte` = 16 字节（减少 33%）。并发下还有「伪共享」：两字段落同一 64B cache line，不同 goroutine 分别写入导致 cache line 反复失效。

### 核心原理
按字段大小**从大到小排列**。大对齐字段（`int64`、`*T`、`float64`）放前，小对齐字段（`bool`、`byte`）放后自然填补。并发场景在共享字段间插入 `[N]byte` padding 分隔到不同 cache line。

### 关键代码
```go
// 差：24 字节
type PoorlyAligned struct {
    flag  bool
    count int64  // 前有7字节 padding
    id    byte   // 后有7字节 padding
}

// 好：16 字节
type WellAligned struct {
    count int64
    flag  bool
    id    byte   // flag+id 紧挨
}

// 避免伪共享
type SharedCounterGood struct {
    a int64
    _ [56]byte  // padding 到不同 cache line
    b int64
}
```

### 性能数据
1000 万个实例：对齐差 240MB；对齐好 160MB —— 节省 **80MB**。消除伪共享速度提升约 3.8%。

### 注意事项
- **零成本优化** —— 只改字段顺序，不改逻辑
- 用 `fieldalignment` linter 自动检测
- 伪共享 padding 仅在两个字段被不同 goroutine 频繁并发写入时有效
- cache line 通常 64B（Apple M 系列 128B）

---

## 1.4 避免接口装箱 (Avoiding Interface Boxing)

### 问题场景
将具体类型值赋给接口变量时（`var i interface{} = someStruct`），Go 执行「装箱」—— 创建含类型描述符和数据指针的两字结构。大结构体非指针赋值会在堆上分配副本并拷贝整个结构体。循环中或热路径上隐式堆分配显著增加 GC 压力。

### 核心原理
接口在 Go 中表示为 `(type, data)` 二元组。值类型赋给接口时调用 `runtime.convT2E` 等函数，堆分配 + 拷贝。改用指针（`&s`）只存 8 字节指针，避免结构体拷贝。

### 关键代码
```go
// 差：每次 append 拷贝整个 4KB 结构体
jobs := make([]Worker, 0, 1000)
for j := 0; j < 1000; j++ {
    var job LargeJob             // 4KB 结构体
    jobs = append(jobs, job)     // 装箱：堆分配 + 拷贝 4KB
}

// 好：只存 8 字节指针
for j := 0; j < 1000; j++ {
    job := &LargeJob{}
    jobs = append(jobs, job)     // 接口只存指针
}
```

### 性能数据
1000 个 4KB 结构体装箱：值版本 404µs，指针版本 340µs —— 快约 **19%**。

### 注意事项
- Go 1.4 后所有值类型装箱必然触发堆分配
- 小型值（`int`、`float64`）装箱开销极小，别为消除它牺牲可读性
- 热路径中已知具体类型直接用，避免接口间接调用
- 泛型替代 `[]interface{}` 可获得类型安全 + 零装箱
- pprof 查 `runtime.convT2E` 定位装箱热点

---

## 1.5 零拷贝技术 (Zero-Copy)

### 问题场景
I/O 密集型应用频繁在内核空间和用户空间之间拷贝数据。每次拷贝消耗 CPU、内存带宽、触发缓存失效。高吞吐场景下成为瓶颈。

### 核心原理
直接在已有内存上操作，避免中间拷贝：
1. **Slice 重切片** —— 新 slice 共享底层数组，不复制
2. **`io.CopyBuffer`** —— 复用用户提供的缓冲区
3. **`mmap`** —— 文件直接映射到进程地址空间（真正零拷贝）
4. **流式 `io.Reader/Writer`** —— 数据边读边传，不积攒

### 关键代码
```go
// slice 零拷贝
func process(buffer []byte) []byte {
    return buffer[128:256]  // 不拷贝，只返回新视图
}

// 复用缓冲区
func StreamData(src io.Reader, dst io.Writer) error {
    buf := make([]byte, 4096)
    _, err := io.CopyBuffer(dst, src, buf)
    return err
}

// mmap 真正零拷贝
data, _ := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ, unix.MAP_SHARED)
defer unix.Munmap(data)
h := xxhash.New()
h.Write(data[:4*1024*1024])  // 直接读映射页
```

### 性能数据
| 场景 | 零拷贝？ | 效果 |
|------|----------|------|
| slice 重切片 vs 拷贝 | 是 | **7000 倍**（0.59ns vs 4246ns） |
| mmap + XXHash vs ReadAt | 是 | **~2 倍**（281µs vs 539µs） |
| mmap + SHA256 vs ReadAt | 是 | 差距缩小（计算密集型稀释收益） |

### 注意事项
- 切片共享底层数组 → 修改互相影响，需要明确所有权
- `mmap` 不是免费的 —— 占用虚拟地址空间
- `golang.org/x/exp/mmap` 的 `ReadAt` 不暴露映射内存，仍有拷贝
- 只在 benchmark 证明收益显著的路径上引入

---

## 1.6 GC 调优 (Memory Efficiency: Mastering Go's GC)

### 问题场景
GC 在以下场景成为瓶颈：(1) 高并发低分配 —— 几乎不产生垃圾但 GC 周期性运行占用 CPU；(2) 容器化部署 —— OOM killer 可能在 GC 反应前杀进程；(3) 延迟敏感系统 —— STW 即使微秒级也影响尾延迟；(4) 大堆场景 —— 标记扫描遍历大量存活对象。

### 核心原理
Go 使用**非分代、并发、三色标记-清扫** GC：
- 不分代（团队认为分代未在 Go 中展示连续优势）
- 并发执行（标记和清扫与应用 goroutine 交替，仅少量 STW）
- 默认 GOGC=100（堆增长 100% 触发下一次 GC）
- Go 1.19+ 引入 `GOMEMLIMIT` 软内存限制

### 关键配置
```bash
# 不那么频繁 GC（CPU 密集型）—— Cloudflare 设 GOGC=11300 获 22 倍吞吐
GOGC=11300

# 容器场景：防 OOM
GOMEMLIMIT=400MiB

# 组合：只在接近 2GiB 时触发
GOMEMLIMIT=2GiB GOGC=off

# 程序内设置
import "runtime/debug"
debug.SetMemoryLimit(2 << 30)
```

### 性能数据
- **Uber**：动态调整 GOGC，30 个关键服务节省约 7 万个 CPU 核心
- **Cloudflare**：高并发加密场景，GOGC=11300，单核性能提升 **22 倍**

### 注意事项
- **先 profiling，再调优**。默认 GC 对大多数场景足够好
- `GOGC=off` 仅极特殊场景（短生命周期 CLI）使用
- `GOMEMLIMIT` 设太低 → GC 过于激进 → CPU 飙升
- 减少 GC 压力的最佳方式是减少堆分配（栈分配、sync.Pool、预分配等）
- Go 1.24 新增 `weak` 包，适合缓存/去重映射场景

---

## 1.7 栈分配与逃逸分析 (Stack Allocations & Escape Analysis)

### 问题场景
栈分配随函数返回自动释放（零 GC），堆分配需 GC 追踪回收。编译器通过逃逸分析决定分配位置。本意栈分配的变量因返回指针、接口装箱、闭包捕获等原因「逃逸」到堆。

### 核心原理
逃逸分析判断变量在其声明函数返回后是否仍被引用。生命周期不超过函数作用域 → 栈；否则 → 堆。常见逃逸原因：返回局部变量指针、闭包捕获、接口转换、赋值给全局变量/堆字段。

### 关键代码
```go
// 栈分配：返回值
func StackAlloc() Data { return Data{1, 2, 3} }

// 真正逃逸：赋值给全局变量
var sink *Data
func HeapAllocEscape() {
    d := &Data{1, 2, 3}
    sink = d  // d 逃逸到堆
}
```

### 查看逃逸分析
```bash
go build -gcflags="-m" ./...
# 关注 "moved to heap" 信息
```

### 性能数据
栈分配 0.26ns + 0 alloc；真正堆逃逸 10.55ns + 1 alloc + 24B —— 约 **40 倍**差距。

### 注意事项
- 并非所有逃逸都需消除 —— `NewXXX()` 返回指针是 Go 惯用法
- 小型低频分配即使堆分配也无妨
- 强迫避免逃逸可能导致大结构体拷贝反而更差
- 编译器逃逸分析持续改进，不要基于旧版本过早优化
- 用 `-gcflags="-m"` 诊断，benchmark 验证，pprof 定位

---

# 二、并发与同步（5 篇）

## 2.1 Goroutine 工作池 (Worker Pool)

### 问题场景
为每个任务启动 goroutine 处理大量并发任务时，内存膨胀、GC 压力增大、上下文切换频繁、CPU 缓存局部性下降。任务数无界可能直接 OOM。

### 核心原理
固定数量 worker goroutine 消费共享 job channel，形成「背压」（backpressure）。并发度限制在 worker 数以内，无论收到多少任务，同时运行的 goroutine 不变。

Worker 数与 CPU 核心数关系：
- CPU 密集型：≈ `runtime.NumCPU()` 或略少
- I/O 密集型：可超过核心数（worker 阻塞等待 I/O）

### 关键代码
```go
for w := 1; w <= 5; w++ {
    go worker(w, jobs, results)
}
for j := 1; j <= 10; j++ {
    jobs <- j
}
close(jobs)
```

### 性能数据
| 方案 | 耗时 (ns/op) | 内存 (B/op) | 分配次数 |
|------|-------------|-------------|---------|
| 无界 goroutine | 2,499,213 | 639,350 | 39,754 |
| 10-worker 池 | 1,791,772 | 320,707 | 19,762 |

工作池快约 **40%**，内存占用减半。

### 注意事项
- Worker 超合理数量后吞吐反而下降（调度争抢、CPU 上下文切换、cache thrashing）
- 不适合延迟敏感场景（排队等待引入额外延迟）
- 任务量小且确定时直接 goroutine 更简单

---

## 2.2 原子操作 (Atomic Operations)

### 问题场景
高并发下对共享计数器、状态标志等简单变量做安全操作。`sync.Mutex` 对「加 1」「读 flag」这种极轻量操作杀鸡用牛刀 —— 锁竞争导致上下文切换。

### 核心原理
`sync/atomic` 利用 CPU 级原子指令（CAS — Compare-And-Swap）直接操作内存，完全绕过锁。Go 默认实现**顺序一致性**（等价 C++ `memory_order_seq_cst`），提供全内存屏障。

### 关键代码
```go
// 高吞吐计数器
var requests atomic.Int64
requests.Add(1)

// 无锁 shutdown 标志
var shutdown atomic.Int32
if shutdown.Load() == 1 { break }
shutdown.Store(1)

// 无锁栈（CAS 循环）
head.CompareAndSwap(old, n)

// 原子 + 锁混合：快路径过滤
if !atomic.CompareAndSwapInt32(&flag, 0, 1) { return }
mu.Lock()
// 执行一次性初始化
mu.Unlock()
```

### 性能数据
14 核并行自增：`atomic.AddInt64` 80.40ns；`sync.Mutex` + `counter++` 110.7ns。**原子操作快约 27%**，高竞争下差距更显著。

### 注意事项
- 只能操作简单类型（单次读写/CAS），无法保护多步临界区
- 没有 rollback 机制
- 不要通过 `go:linkname` 访问内部 relaxed atomics
- 手动实现的初始化协议比 `sync.Once` 灵活但更易出错

---

## 2.3 延迟初始化 (Lazy Init — sync.Once)

### 问题场景
数据库连接池、大型缓存、配置加载等资源创建成本高。如果某些代码路径永不触发，启动时初始化就是浪费。

### 核心原理
Go 提供三层 API：
- `sync.Once` —— 保证函数只执行一次（无返回值）
- `sync.OnceValue` (Go 1.21+) —— 返回单值
- `sync.OnceValues` (Go 1.21+) —— 返回多值（如 `(*Config, error)`）

### 关键代码
```go
// sync.Once 经典用法
var (
    resource *MyResource
    once     sync.Once
)
func getResource() *MyResource {
    once.Do(func() { resource = expensiveInit() })
    return resource
}

// OnceValue (Go 1.21+) 更简洁
var getResource = sync.OnceValue(func() *MyResource {
    return expensiveInit()
})

// OnceValues 处理带错误的初始化
var getConfig = sync.OnceValues(func() (*Config, error) {
    return loadConfig("config.yml")
})
```

### 性能数据
收益体现在避免了不必要的资源创建、加速启动、降低内存，而非直接的运行时吞吐差异。应在生产环境 profiling 验证。

### 注意事项
- 手动 `atomic.CompareAndSwap` 替代 `sync.Once` 可行但丢失安全保证 —— 除非 profiling 明确显示瓶颈，否则不要这样做

---

## 2.4 不可变数据共享 (Immutable Data Sharing)

### 问题场景
多 goroutine 频繁读取共享数据（全局配置、路由表、feature flags）。`sync.RWMutex` 保护下「纯读」也要获取读锁，高并发下锁竞争成为瓶颈。

### 核心原理
数据创建后永不修改（immutable）→ 读取方无需任何同步。更新时构建全新副本，用 `atomic.Pointer` / `atomic.Value` 原子替换指针 —— **Copy-on-Write** 模式。读路径零锁开销。

### 关键代码
```go
// 深拷贝 map 确保不可变
func NewConfig(features map[string]bool) *Config {
    copied := make(map[string]bool, len(features))
    for k, v := range features { copied[k] = v }
    return &Config{Features: copied}
}

// 原子替换
var currentConfig atomic.Pointer[Config]
currentConfig.Store(cfg)
cfg := currentConfig.Load()  // 零锁读

// 路由表同理
var currentRoutes atomic.Pointer[RoutingTable]
```

### 注意事项
- map/slice 必须深拷贝，否则引用类型被外部修改破坏不可变性
- 更新频繁时全量拷贝成本高 → 分段（tenant 隔离）、索引化（map 代替 slice）等策略降级
- 读路径天然支持最终一致性

---

## 2.5 高效 Context 管理 (Context Management)

### 问题场景
HTTP 请求、数据库查询、worker 协程、流式处理中需要统一管理生命周期：超时取消、客户端断开后停止无用计算、优雅关闭通知所有子 goroutine、跨层传递请求级元数据（trace ID 等）。

### 核心原理
`context` 包提供树形传播机制：
- `WithCancel`：手动取消，所有派生 context 收到信号
- `WithTimeout` / `WithDeadline`：到时自动取消
- `WithValue`：附加请求级元数据（非业务参数）
- `signal.NotifyContext`：OS 信号转 context 取消

### 关键代码
```go
// 超时数据库查询
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
rows, err := db.QueryContext(ctx, "SELECT * FROM users")

// 优雅关闭 CLI
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
<-ctx.Done()

// worker 管理：一次 cancel，全部退出
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx, 1)
cancel()

// 流式管道检查取消
select {
case <-ctx.Done(): return
case data := <-ch: process(data)
}
```

### 注意事项
- context 必须作为函数**第一个参数**显式传递，**禁止**存 struct 字段或全局变量
- `WithValue` 只用于请求级元数据，不塞业务逻辑或大对象
- 区分 `context.Canceled`（主动取消）和 `context.DeadlineExceeded`（超时）
- 每次 `WithTimeout`/`WithCancel` 都要 `defer cancel()` 防 goroutine 泄漏

---

# 三、I/O 优化与吞吐（2 篇）

## 3.1 高效缓冲 (Buffered I/O)

### 问题场景
高频小量 I/O 操作（循环写文件）每次 `Write` 触发一次系统调用，涉及用户态/内核态切换。每秒数千次积累严重拖垮吞吐。

### 核心原理
`bufio.Writer` 维护内部缓冲区（默认 4KB）。写入数据先填充内存缓冲，仅缓冲满或 `Flush()` 时触发一次系统调用。4KB 对齐 OS 内存页大小和文件系统块大小。高吞吐可调大（16KB~32KB），低延迟可调小。

### 关键代码
```go
// 无缓冲：每次 Write 触发 syscall
f, _ := os.Create("output.txt")
for i := 0; i < 10000; i++ {
    f.Write([]byte("line\n"))
}

// 有缓冲：大幅减少 syscall
f, _ := os.Create("output.txt")
buf := bufio.NewWriterSize(f, 16*1024)
for i := 0; i < 10000; i++ {
    buf.WriteString("line\n")
}
buf.Flush() // 必须！
```

### 性能数据
| 方案 | 耗时 | 分配次数 |
|------|------|---------|
| 无缓冲 | 23,672,792 ns | 10,007 |
| 有缓冲 (16KB) | 379,703 ns | 10,008 |

缓冲版本快约 **62 倍**。

### 注意事项
- **必须调用 `Flush()`** —— 不会在 Close 时自动 flush
- 缓冲大小需通过 profiling 和 benchmark 测量确定，不硬编码
- 低延迟场景（交互式 CLI）慎用大缓冲

---

## 3.2 批量操作 (Batching Operations)

### 问题场景
高频调用外部系统（文件、数据库、RPC）时每次调用携带固定开销：系统调用、网络往返、锁竞争。日志服务每秒数千次 `WriteString`，文件系统被海量细小写请求淹没。

### 核心原理
将多个小操作聚合为一次大操作。累积 100 条日志后一次写出。本质是**用内存换 syscall 次数**，用微小的批延迟换取巨大吞吐提升。

### 关键代码
```go
// 无批量
func logLine(line string) {
    f.WriteString(line + "\n")
}

// 有批量：累积 100 条后一次写入
var batch []string
func logBatch(line string) {
    batch = append(batch, line)
    if len(batch) >= 100 {
        f.WriteString(strings.Join(batch, "\n") + "\n")
        batch = batch[:0]
    }
}
```

### 性能数据
| 场景 | 未批量 | 批量 | 加速 | 分配减少 |
|------|--------|------|------|---------|
| 内存处理 | 2,028,492 ns | 2,094,168 ns | ~持平 | 50x |
| **文件 I/O** | 12,766,433 ns | 993,912 ns | **12.8x** | 48x |
| SHA256 | 1,232,242 ns | 675,303 ns | **1.8x** | 75x |

文件 I/O 收益最大（12 倍+），即使纯内存处理 GC 压力也大幅降低。

### 注意事项
- **数据丢失风险**：进程崩溃时内存中未落盘数据丢失
- Mutex 不可重入：不要在 `flush()` 回调中再调用 `Batcher.Add()`
- 不适合要求每条数据立即生效的场景（实时交易）

---

# 四、编译器级优化（2 篇）

## 4.1 编译器优化标志 (Compiler Flags)

### 问题场景
大多数 Go 性能调优聚焦运行时，编译器层面优化常被忽略。合理使用编译/链接标志可在二进制体积、调试体验、跨平台构建、静态链接等方面获得显著收益。

### 核心原理
三组关键标志：
- `-ldflags`：链接器行为（剥离调试信息、注入变量、静态链接）
- `-gcflags`：编译器优化级别（禁用内联/优化以便调试）
- `GOOS`/`GOARCH`：交叉编译目标

### 关键命令
```bash
# 缩小二进制体积（生产发布）- 减 30-40%
go build -ldflags="-s -w" -o app main.go

# 调试友好构建（配合 Delve）
go build -gcflags="all=-N -l" -o app main.go

# 交叉编译
GOOS=linux GOARCH=arm64 go build -o app main.go

# 注入构建时变量
go build -ldflags="-X main.version=1.0.0" -o app main.go

# 完全静态链接（scratch 容器）
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  go build -tags netgo \
  -ldflags="-linkmode=external -extldflags '-static'" \
  -o app main.go
```

### 注意事项
- `-s -w` 去除调试信息后无法用 Delve 源码级调试
- 静态链接需系统上有对应 C 库的 `.a` 文件
- `netgo` build tag 强制 Go 原生 DNS 解析器，scratch 容器必备
- `-gcflags="all=-N -l"` 显著降低运行时性能，仅调试用

---

## 4.2 栈分配与逃逸分析 (Stack Allocations & Escape Analysis)

> 与 §1.7 相同，索引页将其同时归入内存管理和编译器优化两个类别。

见 **§1.7**。核心要点重申：

```bash
# 查看逃逸分析结果
go build -gcflags="-m" ./...

# 更详细
go build -gcflags="-m -m" ./...
```

| 分配方式 | 耗时 | 分配 | 差距 |
|---------|------|------|------|
| 栈分配 | 0.26 ns | 0 | — |
| 堆逃逸 | 10.55 ns | 1 alloc + 24B | **40 倍** |

---

# 总结速查

## 四大类别核心思想

| 类别 | 核心策略 | 最显著收益 |
|------|---------|-----------|
| 内存管理 | 减少堆分配、复用内存 | sync.Pool **20x**、零拷贝 slice **7000x**、GOGC 调优 **22x** |
| 并发同步 | 无锁化、不可变化、固定并发度 | Worker Pool **40%** 提速、atomic 比 mutex **27%** 快 |
| I/O 优化 | 减少系统调用、批量处理 | bufio **62x**、batching 文件 I/O **12.8x** |
| 编译器 | 理解编译器行为、构建期优化 | 堆逃逸比栈慢 **40x**、`-s -w` 减 **30-40%** 体积 |

## 优先级建议

```
第一优先（低投入高回报）:
  → make 预分配、字段对齐、bufio 缓冲、sync.Once

第二优先（需要 profiling 确认热点）:
  → sync.Pool、零拷贝、Worker Pool、不可变数据

第三优先（系统级调优）:
  → GOGC/GOMEMLIMIT、编译器标志、逃逸分析消除

永远不要做的:
  → 不 profiling 就优化
  → 为消除逃逸写出晦涩代码
  → 在低 QPS 路径上引入 Pool 增加复杂度
  → 把 context 存到 struct 字段
```

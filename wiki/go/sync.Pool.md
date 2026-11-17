---
title: "sync.Pool"
created: 2023-04-13
updated: 2026-10-07
status: snapshot
type: concept
tags:
  - go/runtime
  - go/memory
  - go/performance
source: "raw/2026/10/07/sync.Pool 视频英文字幕-fwHok9ZhQaY.md"
sources:
  - "raw/2026/10/07/sync.Pool 视频英文字幕-fwHok9ZhQaY.md"
  - "https://github.com/golang/go/blob/go1.27.1/src/sync/pool.go"
  - "https://github.com/golang/go/blob/go1.27.1/src/bytes/buffer.go"
  - "https://github.com/golang/go/blob/go1.27.1/src/net/http/server.go"
  - "https://github.com/golang/go/blob/go1.27.1/src/runtime/runtime2.go"
  - "https://github.com/golang/go/blob/go1.27.1/src/runtime/iface.go"
source_kind: synthesis
as_of: 2026-10-07
confidence: SYNTHESIZED
contested: false
contradictions: []
related:
  - "垃圾回收"
  - "Go性能优化常见模式"
---

# sync.Pool

`sync.Pool` 缓存已经分配、暂时不用的对象，供后续操作复用，摊薄分配成本并减轻 GC 压力。它是可丢失的临时对象缓存，不是保证容量、对象身份或存活时间的资源池。

本文整合 Phuong Le 的视频 [sync.Pool in Go Explained & Avoid The Heap Allocation Mistake](https://www.youtube.com/watch?v=fwHok9ZhQaY)（2025-11-06 发布），以及旧笔记的 per-P / victim 源码分析。视频摘要依据英文字幕；源码核对使用本机 Go 1.27.1 的 `src/` 文件，以上链接对应同版本标签。实现细节不构成跨版本 API 保证，视频也明确声明部分解释经过简化。

## 1. 为什么需要对象复用（00:00）

> 摘录：The trade-off is extra overhead.

频繁创建堆对象有两类成本：分配内存，以及对象不再可达后由 GC 回收。视频还介绍了标记、清扫、写屏障和阶段切换中的 STW；对象复用的目标是减少重复分配及相应回收工作，而不是取消 GC。

视频用“值需要比创建它的函数活得更久”解释堆分配。这是入门直觉，不是完整的逃逸分析判定规则。不能据此断言每个临时对象必然分配在堆上，也不能断言引入 Pool 必然更快。

适合候选：每次请求的格式化缓冲区、编码临时对象，以及初始化或扩容成本值得复用的对象。判断收益应基于实际分配热点与基准测试，而非仅看代码里是否出现 `new`。

相关背景见 [[垃圾回收]]；工程模式见 [[Go性能优化常见模式]]。这两页只作延伸阅读，不替代本页的核心解释。

## 2. New、Get、Put 与使用流程（01:27）

> 摘录：There are just three things.

| API | 行为 | 注意事项 |
|---|---|---|
| `New func() any` | `Get` 没有取到对象时调用的可选工厂 | 无参数；不能与 `Get` 并发修改 |
| `Get() any` | 从池取出并移除一个任意对象；未命中时调用 `New` | 没有 `New` 且未命中时返回 `nil`；不保证返回先前 `Put` 的对象 |
| `Put(x any)` | 将对象放回池中供后续复用 | `Put(nil)` 无效果；归还后不能继续当作自己的对象使用 |

Pool 的零值可用，不需要构造函数。首次使用后不能复制 Pool。Pool 支持多个 goroutine 并发调用，但这不等于池中对象允许并发读写。

源码文档规定：`Put(x)` synchronizes before 返回同一个 `x` 的 `Get`；`New` 返回 `x` 也 synchronizes before 返回该 `x` 的 `Get`。这是对象交接的内存可见性关系，不是对归还后继续访问的许可。

### 示例：复用 bytes.Buffer

下面是根据视频思路整理的示例，不是视频画面代码的逐字转录。

```go
var bufferPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func writePayload(dst io.Writer, payload []byte) (int, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    _, _ = buf.Write(payload)
    return dst.Write(buf.Bytes())
}
```

`Get → Reset → 使用 → Put` 是基本流程。示例在 `Write` 返回后才归还；如果需要让数据在函数返回后继续存在，就不能把仍引用该 Buffer 底层数组的切片随池对象一起交还。`defer` 保证正常返回和错误返回都执行归还，但不保证复用命中。

### 为什么有些池不设置 New

`New` 没有参数，但对象创建有时依赖当前参数。视频举了标准库 HTTP 包的例子：在包装函数里先 `Get`，未命中再用当前参数创建，命中则重置已有对象。

Go 1.27.1 `net/http/server.go:newBufioWriterSize(w, size)` 根据大小选择池；命中后执行 `bw.Reset(w)`，否则 `bufio.NewWriterSize(w, size)`。归还函数 `putBufioWriter` 先 `Reset(nil)`，避免一直保留旧 Writer 引用。

## 3. 什么对象值得放入池（04:14）

> 摘录：that are created and discarded frequently.

关键是**单次使用短、反复使用频繁、复用能摊薄成本**。短生命周期指每次借用时间短，不意味着底层对象只能存在很短时间。

独立请求可以分别借用同类临时对象；使用者负责清理旧状态，避免下一次操作读到上一次的数据。长时间闲置的对象可能在下一次使用前已被池丢弃，池化收益因此下降。

Pool 不是连接池、限流器或持久缓存：对象允许无通知地丢失，不提供容量上限，也没有资源关闭回调。需要确定数量、等待可用连接或按时关闭资源时，不能用 Pool 的缓存老化替代相应管理机制。

## 4. per-P 本地缓存与 Get/Put 路径

以下保留旧笔记的实现主题，并按 Go 1.27.1 源码核对。`P` 是 Go 调度器的处理器，不是固定的 goroutine 私有空间。

```go
type Pool struct {
    noCopy noCopy
    local     unsafe.Pointer // 实际指向每个 P 的 poolLocal 数组
    localSize uintptr
    victim     unsafe.Pointer // 上一轮留下的 local
    victimSize uintptr
    New func() any
}

type poolLocalInternal struct {
    private any
    shared  poolChain
}

type poolLocal struct {
    poolLocalInternal
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```

`pin()` 将当前 goroutine 暂时固定到当前 P 并禁用抢占，访问对应的 `poolLocal`，结束后 `runtime_procUnpin()`。

- `private` 只供对应 P 使用，避免跨 P 争抢。
- `shared` 是共享链：当前 P 可 `pushHead/popHead`，其他 P 可 `popTail`。
- `Put` 优先填空的 `private`，否则放入本地 `shared` 的头部。
- `Get` 先取本地 `private`，再取本地 `shared` 的头部；头部优先有利于时间局部性。
- `getSlow` 先从其他 P 的主缓存 `shared` 尾部偷取，再尝试 victim 中当前 P 的 `private` 和 victim 的共享链。
- 仍未取到且设置了 `New`，才调用 `New`，直接返回结果；**不会自动把结果放入 local**。只有调用方后续 `Put` 才归还。

`pinSlow` 在 `p.local == nil` 时把 **Pool 指针加入 `allPools`**，并分配新的 per-P 数组。旧笔记“把 p 追加进 local”的说法不准确：注册池与分配本地数组是两个动作。GC 后 `Get` 和 `Put` 都可能触发这个慢路径。

> [!tip] 伪共享
> 多个处理器频繁修改同一缓存行内互不相关的数据时，会因缓存一致性流量产生额外成本。`poolLocal` 的 padding 用于在源码注释所描述的常见缓存行大小条件下避免伪共享；不是说所有平台的缓存行都等于 128 字节。

## 5. GC 与 victim cache（04:54）

> 摘录：to the garbage collection cycle.

API 文档的保证是：池中任意对象可能随时、无通知地被移除。以下 local/victim 两阶段机制只能解释所核对版本的实现，不能改写为“保证至少保留两轮 GC”。

Go 1.27.1 的 `poolCleanup()` 在 GC 开始的 STW 中执行：

```go
func poolCleanup() {
    for _, p := range oldPools {
        p.victim = nil
        p.victimSize = 0
    }
    for _, p := range allPools {
        p.victim = p.local
        p.victimSize = p.localSize
        p.local = nil
        p.localSize = 0
    }
    oldPools, allPools = allPools, nil
}
```

| 时点 | 未再取出的对象 x | 后续含义 |
|---|---|---|
| `Put(x)` 后 | 在主缓存 local 中 | 可以被 `Get` 取出 |
| 下一次 GC 开始 | 原 local 转为 victim，local 清空 | 主缓存未命中后仍可尝试 victim |
| 再下一次 GC 开始 | 旧 victim 被清空 | Pool 不再保留 x；是否可回收取决于其他引用 |

这里不是交换两组数组：先丢弃旧 victim，再把主缓存转交给新的 victim。重新使用时按需建立 local。如果期间 x 从 victim 被取出、使用并 `Put`，会重新进入主缓存，不能再按原来那次归还的两轮来判断它何时丢失。

已经借出的对象不在池中，GC 清空池缓存不会使调用者手上的可达对象失效。清空缓存是移除引用，不是承诺对象在该瞬间立即完成内存回收。

> [!warning] 视频简化与 API 边界
> 视频 05:42 将 `New` 创建和 `Put` 都描述成进入 local；实际 `New` 的结果直接交给调用方。06:24 的“至少两轮 GC”应理解为缓存老化示意，不是对象存活承诺。源码中 race 模式的 `Put` 还可能随机丢弃对象；`pinSlow` 注释也说明 GOMAXPROCS 变化时可能丢失旧数组。

victim 可理解为上一代缓存，为仍有复用机会的对象提供缓冲。旧笔记中关于 CPU 硬件 victim cache 的长篇描述只是类比，不能用来推导 Go Pool 的具体行为。收益应表述为缓解 GC 后的缓存冷启动、保留复用机会，而不是“把一次 GC 的开销拉长到两次”。

## 6. 为什么 bytes.Buffer 适合复用（06:49）

> 摘录：to zero while keeping its capacity.

Buffer 内部包含字节切片及读写状态。数据增长超过可用容量时，可能需要分配更大的底层数组并拷贝数据；每次从空 Buffer 开始，可能重复承担扩容和回收成本。

Go 1.27.1 `bytes/buffer.go:Buffer.Reset` 的实现：

```go
func (b *Buffer) Reset() {
    b.buf = b.buf[:0]
    b.off = 0
    b.lastRead = opInvalid
}
```

重置的是逻辑长度和状态，已有容量保留，因此下一次相似大小的操作能复用扩容成果。它不保证下一次绝不分配：请求超过容量、池未命中等情况下仍会分配。

> [!note] SYNTHESIZED：保留容量也是内存权衡
> 保留大容量能减少扩容，也可能使偶发的大请求留下过大的闲置数组。应用可按自身负载决定是否拒绝归还超大 Buffer；这属于应用策略，不是 Pool 自带的容量限制。`Reset` 也不是安全擦除数据，不会逐字节清零底层数组。

## 7. 接口装箱与额外堆分配陷阱（08:15）

> 摘录：data field is itself a pointer.

`Put` 接收 `any`。Go 1.27.1 `runtime/runtime2.go:eface` 的表示为：

```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

视频用整数说明值语义：把整数赋给接口，接口保存的是该值，不会随原整数变量的后续修改而变化。对需要间接存储的值，运行时可能为接口的数据准备独立存储。

### 值放回池可能重新装箱

```go
var valuePool = sync.Pool{
    New: func() any { return 500 },
}

v := valuePool.Get().(int)
v = 500
valuePool.Put(v)
```

尽管整数看上去被复用了，`Put(v)` 仍需把当前值转换成接口。视频以 500 展示这一步产生堆上的值副本。性能陷阱在于：复用逻辑没有消除每次归还的装箱分配。

`runtime/iface.go:convT64` 对小整数有 `staticuint64s` 路径，对超出缓存范围的值才走分配路径。所以不能推广成“所有值装箱都必然分配”。转换是否被优化、存储是否逃逸，也必须结合具体代码判断。

### 指针避免再次分配值的装箱副本

```go
var pointerPool = sync.Pool{
    New: func() any {
        v := 500
        return &v
    },
}

p := pointerPool.Get().(*int)
*p = 500
pointerPool.Put(p)
```

对象首次创建仍可能需要堆分配；区别是归还已有指针时，接口直接表示该指针，不必每次再创建整数副本。它不是“使用指针就没有任何堆分配”，而是避免已有对象在归还时的额外装箱分配。

### 切片也有这个问题

切片值是描述底层数组的 header，不是数组本身。直接 `Put([]byte)` 时，即使底层数组复用了，header 装箱仍可能分配。`runtime/iface.go:convTslice` 对非 nil 底层数组的切片具有分配 header 存储的路径。

可以池化 `*[]byte`、包含切片的结构体指针，或者 `*bytes.Buffer`。注意：切片 header 的拷贝不等于底层数组被深拷贝；“池化指针”也不免除重置状态和处理别名的责任。

### 如何验证，而不是只看直觉

```bash
go build -gcflags='-m=2' ./...
go test -run '^$' -bench . -benchmem -count=5 ./...
```

观察 `B/op` 和 `allocs/op`，比较值版本与指针版本，并区分冷启动、缓存命中和并发场景。视频解释：`a escapes to heap` 不一定意味着原变量本身移到堆上，也可能是在接口转换点为值安排了堆存储；不能忽略具体诊断位置。

视频用“goroutine 栈私有，所以共享必须放堆”作直觉说明，但更准确的判定依据是编译器能否证明相应存储的生命周期满足使用要求。不能仅凭发生跨 goroutine 使用，就替代逃逸分析或实测。

## 8. 复习要点

- Pool 缓存的是临时对象，不承诺命中、顺序或永久保留。
- 使用者负责 `Reset`、对象独占以及在正确时点归还。
- per-P 分片减少争抢；Get 先主缓存、后 victim、最后 New。
- `New` 返回对象不会自动入池；`pinSlow` 注册的是 `allPools`。
- victim 的两代老化是实现机制，不是“两轮 GC 存活保证”。
- Buffer 复用的收益主要来自保留底层容量，但要权衡内存驻留。
- 优先复用指针形态的对象；不要让值或切片 header 的装箱抵消池化收益。
- 最终判断依赖具体程序的逃逸分析与基准测试，不将教学示例的结果泛化到所有负载。

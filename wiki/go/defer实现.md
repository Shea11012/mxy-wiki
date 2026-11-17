---
title: "defer实现"
created: 2022-03-31
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - go
  - go/runtime
confidence: INFERRED
related:
  - "panic和recover"
---
defer 会在当前函数返回前执行传入的函数，经常被用于关闭文件描述符、关闭数据库以及解锁资源、捕获 panic

defer 会对函数参数进行预计算，即如果函数参数是一个表达式会立即求值并进行参数拷贝

## 结构
```go
type _defer struct {
	heap      bool
	rangefunc bool               // 标记是否为 range-over-func 的 defer
	sp        uintptr            // sp at time of defer
	pc        uintptr            // pc at time of defer
	fn        func()             // 执行的函数，可以为 nil（开放编码 defer）
	link      *_defer            // 链表下一项，可指向堆或栈

	// rangefunc 为 true 时，*head 是 range-over-func 执行期间的原子链表头
	head *atomic.Pointer[_defer]
}
```

与旧版相比，结构体从 10+ 个字段精简为 7 个：
- **移除** `siz`（参数大小）、`started`（是否已开始执行）、`openDefer`（是否开放编码）、`_panic`（关联的 panic）、`fd`（funcdata）、`varp`（栈帧 varp）、`framepc`（栈帧 PC）
- **新增** `rangefunc` 和 `head`，用于支持 Go 1.23 引入的 range-over-func
- `fn` 类型从 `*funcval` 简化为 `func()`

## 执行方式
defer 有三种执行方式：
- 堆分配
- 栈分配
- 开放编码

### 堆分配
默认方式，通过 `deferproc(fn func())` 调用。运行时使用 per-P 池（`p.deferpool`）复用 `_defer` 结构体，减少堆分配开销。池满时将一半转移至全局调度器池（`sched.deferpool`）。

### 栈分配
当 defer 关键字在函数体中最多执行一次时，编译期间会将结构体分配到栈上，通过 `deferprocStack(d *_defer)` 调用。此时 `d.heap = false`，执行完毕后无需归还池。

### 开放编码
使用代码内联优化 defer 关键字，避免运行时链表操作的开销。

编译器在函数入口创建 `deferBits`（`uint8` 类型）变量作为位掩码，每个 defer 对应一个 bit。defer 语句执行时设置对应 bit，函数返回时按 bit 位逆序内联调用 defer 函数。8 个 bit 最多支持 8 个 defer。

开放编码的 defer 信息不再存储在 `_defer` 结构体中，而是在 panic/Goexit 处理时由 `_panic` 结构体中的 `deferBitsPtr` 和 `slotsPtr` 字段，结合 `FUNCDATA_OpenCodedDeferInfo` 在栈上定位。

#### 满足开放编码的条件
全部条件必须同时满足：
- 未禁用优化（`-N` 标志未设置）
- 函数的 defer 数量 <= 8
- 函数的 defer 关键字不能在循环中执行（由逃逸分析判定）
- 函数的 return 语句数量 × defer 语句数量 <= 15
- 非 386 架构的 shared/dynlink 模式
- 未启用 race 检测（race 检测需要插桩函数出入口）
- 函数的所有结果参数都在栈上（不能有堆分配的结果参数需要回拷）
- 不使用 `deferprocat`（range-over-func 的 defer 机制）

### range-over-func 的 defer
Go 1.23 引入的 range-over-func 会为循环体合成闭包，循环体内的 defer 需要挂载到外层函数而非闭包上。为此新增：

- `deferrangefunc()` — 在外层函数入口调用，创建一个 `rangefunc=true` 的 `_defer` 节点，并初始化 `head` 原子链表
- `deferprocat(fn, frame)` — 循环体内的 defer 通过原子操作追加到 `head` 链表
- `deferconvert(d0)` — 当 panic/defer 处理到达 `rangefunc` 节点时，将原子链表中的 defer 转换为普通链表插入主 defer 链

## defer 的执行与 panic
`deferreturn`、`gopanic`、`Goexit` 共享统一的 defer 遍历逻辑，由 `_panic.nextDefer()` 方法驱动：
1. 优先检查当前栈帧的开放编码 defer（通过 `deferBitsPtr` 位掩码）
2. 再检查链表 defer（通过 `gp._defer` 链表，匹配 `sp`）
3. 遇到 `rangefunc` 节点时调用 `deferconvert` 转换后重新检查
4. 当前帧处理完毕后通过 `_panic.nextFrame()` 向上遍历调用栈

refs:
- [理解 Go 语言 defer 关键字的原理 | Go 语言设计与实现 (draveness.me)](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)
- [深入 Go 语言 defer 实现原理 - luozhiyun`s Blog](https://www.luozhiyun.com/archives/523)
- [Go runtime: runtime/panic.go](https://github.com/golang/go/blob/master/src/runtime/panic.go)
- [Go runtime: runtime/runtime2.go](https://github.com/golang/go/blob/master/src/runtime/runtime2.go)
- [Go compiler: cmd/compile/internal/ssagen/ssa.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssagen/ssa.go)

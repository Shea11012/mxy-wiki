---
date created: 2022-04-05 04:23
date modified: 2022-04-05 05:04
title: panic和recover
---
panic 和 recover 的使用范围：
- panic 只会触发当前 goroutine 的 defer
- recover 只有在 defer 中调用才会生效
- panic 允许在 defer 中嵌套多次调用

## 结构
```go
type _panic struct {
	argp      unsafe.Pointer // 指向defer调用时参数的指针
	arg       interface{}    // 调用panic时传入的参数
	link      *_panic        // 链表，指向了前一个panic
	pc        uintptr        // where to return to in runtime if this panic is bypassed
	sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
	recovered bool           // 当前panic是否被recover恢复
	aborted   bool           // 当前panic是否被强行终止
	goexit    bool
}
```


## panic 流程
1. panic 关键字转为 `runtime.gopanic` 调用，在循环中不断从当前 G 的 defer 链表中获取 defer 执行
2. 在 defer 函数中含有 recover，那么就会调用 `runtime.gorecover` 会修改 `_panic` 的 recovered 字段为 true
3. 调用完 defer 后会回到 `runtime.gopanic` 主逻辑中，检查 recovered 字段为 true，调用 `runtime.recovery` 函数恢复程序并将函数返回值设置为 1
4. 当 `runtime.deferproc` 函数返回 1 时，编译器生成的代码会直接跳转到调用方函数返回之前并执行 `runtime.deferreturn`，然后从 panic 中恢复并执行正常的逻辑
5. `runtime.gopanic` 执行完所有的 defer 且没有遇到 recover，就会执行 `runtime.fatalpanic` 终止程序
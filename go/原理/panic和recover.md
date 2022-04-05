panic和recover的使用范围：
- panic只会触发当前goroutine的defer
- recover只有在defer中调用才会生效
- panic允许在defer中嵌套多次调用

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
1. panic关键字转为`runtime.gopanic`调用，在循环中不断从当前G的defer链表中获取defer执行
2. 在defer函数中含有recover，那么就会调用 `runtime.gorecover` 会修改`_panic`的 recovered 字段为 true
3. 调用完defer后会回到 `runtime.gopanic` 主逻辑中，检查recovered字段为true，调用 `runtime.recovery` 函数恢复程序并将函数返回值设置为1
4. 当 `runtime.deferproc` 函数返回1时，编译器生成的代码会直接跳转到调用方函数返回之前并执行 `runtime.deferreturn`，然后从panic中恢复并执行正常的逻辑
5. `runtime.gopanic` 执行完所有的defer且没有遇到recover，就会执行 `runtime.fatalpanic` 终止程序
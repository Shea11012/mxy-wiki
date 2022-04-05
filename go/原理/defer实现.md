defer 会在当前函数返回前执行传入的函数，经常被用于关闭文件描述符、关闭数据库以及解锁资源、捕获panic

defer 会对函数参数进行预计算，即如果函数参数是一个表达式会立即求值并进行参数拷贝

## 结构
```go
type _defer struct {
	siz     int32 // includes both arguments and results
	started bool
	heap    bool
	// openDefer indicates that this _defer is for a frame with open-coded
	// defers. We have only one defer record for the entire frame (which may
	// currently have 0, 1, or more defers active).
	openDefer bool
	sp        uintptr  // sp at time of defer
	pc        uintptr  // pc at time of defer
	fn        *funcval // 执行的函数地址
	_panic    *_panic  // panic that is running defer
	link      *_defer

	// If openDefer is true, the fields below record values about the stack
	// frame and associated function that has the open-coded defer(s). sp
	// above will be the sp for the frame, and pc will be address of the
	// deferreturn call in the function.
	fd   unsafe.Pointer // funcdata for the function associated with the frame
	varp uintptr        // value of varp for the stack frame
	// framepc is the current pc associated with the stack frame. Together,
	// with sp above (which is the sp associated with the stack frame),
	// framepc/sp can be used as pc/sp pair to continue a stack trace via
	// gentraceback().
	framepc uintptr
}
```

## 执行方式
defer有三种执行方式：
- 堆分配
- 栈分配
- 开放编码

### 堆分配


### 栈分配
当defer关键字在函数体中最多执行一次时，编译期间会将结构体分配到栈上调用

### 开放编码
使用代码内联优化defer关键字，并引入函数数据funcdata管理panic调用

如果defer关键字的执行可以在编译期间确定，会在函数返回前直接插入相应的代码

通过使用 `deferbits` 延迟比特判断条件分支是否执行。延迟比特为8个bit，所以最多只能优化8个defer

满足开放编码的条件：
- 函数的defer数量<=8
- 函数的defer关键字不能在循环中执行
- 函数的return语句与defer语句的乘积<=5

refs:
- [理解 Go 语言 defer 关键字的原理 | Go 语言设计与实现 (draveness.me)](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)
- [深入 Go 语言 defer 实现原理 - luozhiyun`s Blog](https://www.luozhiyun.com/archives/523)
### nil 的零值

在 go 中，任何类型在未初始化时都对应一个零值：

- 布尔类型是 `false`
- 整型是 `0`
- 字符串是 `""`
- 指针、函数、interface、slice、channel、map 是 `nil`

```go
// nil 没有默认的类型，尽管是多个类型的零值，必须显式或隐式指定每个 nil 的用法
// 显式
_ = (*struct{})(nil)
_ = []int(nil)
_ = map[int]bool(int)
_ = chan string(nil)
_ = (func())(nil)
_ = interface{}(nil)

// 隐式
var _ *struct{} = nil
var _ []int = nil
var _ map[int]bool = nil
var _ chan string = nil
var _ func() = nil
var _ interface{} = nil
```

### nil 值比较

1. 不同类型的 nil 值不能比较，**两个不同可比较类型的值只能在一个值可以隐式转换成另一种类型的情况下进行比较**
   1. 两个值之一的类型是另一个的基础类型
   2. 两个值之一的类型实现了另一个值的类型（必须是接口类型）

```go
type IntPtr *int
IntPtr(nil) == (*int)(nil) // true ,IntPtr 的基础类型是 *int,所以可比较
(interface{})(nil) == (*int)(nil) // false
```


2. 同一类型的 nil 有可能无法比较，go 中 map、slice、func 类型是不可比较类型，不可比较类型可以与纯 nil 进行比较

   ```go
   (map[string]int)(nil) == (map[string]int)(nil) // 比较是非法的
   (func())(nil) == nil // true， 纯 nil 可以比较
   ```

   
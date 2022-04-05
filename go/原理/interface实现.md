go中有两种略微不同的接口：
- 带有方法的接口
- 不带任何方法的接口

判断一个接口类型的值是否等于nil，需要动态类型和动态值都为nil时才可以

```go
// 带有方法的接口
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

// 不带方法的接口
type eface struct {
	_type *_type
	data  unsafe.Pointer
}

type _type struct {
	size       uintptr	// 类型大小
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32	// 快速确定类型是否相等
	tflag      tflag	// 反射相关
	align      uint8	// 内存对齐
	fieldAlign uint8	
	kind       uint8	// 基础类型
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte	// GC类型的数据
	str       nameOff	// 名称字符串在二进制文件段中的偏移量
	ptrToThis typeOff 	// 类型元信息指针在二进制文件段中的偏移量
}

type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```
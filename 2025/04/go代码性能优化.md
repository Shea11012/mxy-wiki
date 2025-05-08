---
created: 2025-04-08
tags: [go]
---

## 缓存池
在面对短周期可重复对象、局部的生命周期对象、减少堆分配时使用sync.Pool

## 字段内存对齐
控制对齐可以减少内存占用，提高缓存局部性
```go
type PoorlyAligned struct {
	flag bool
	count int64
	id byte
}

type WellAligned struct {
	count int64
	flag bool
	id byte
}
```

在64位的系统中，PoorlyAligned会占用24bytes，WellAligned只占用16bytes
当并发访问时也会出现伪共享的问题
![[../../go/实现/sync.Pool#^161773]]

在现代CPU中，缓存行宽度为64字节，当CPU加载时会加在包含该结构的整行。如果一个字段在该缓存行被修改，会导致整个 缓存行失效而导致性能下降。
```go
type SharedBad struct {
	a int64
	b int64
}

type SharedGood struct {
	a int64
	_ [56]byte
	b int64
}
```


## 零拷贝
io.CopyBuffer 重用一个buffer，避免重复分配和中转拷贝 
```go
func StreamData(src io.Reader,  dst io.Writer) error {
	buf := make([]byte, 4096) // 可重用buffer
	_,err := io.CopyBuffer(dst, src, buf)
	return err
}
```

mmap
```go
import "golang.org/x/exp/mmap"

func ReadFileZeroCopy(path string) ([]byte,error) {
	r,err := mmap.Open(path)
	if err != nil {
		return nil,err
	}
	defer r.Close()

	data := make([]byte,r.Len())
	_,err := r.ReadAt(data,0)
	return data,err
}
```


 references:
 - [Common Go Patterns for Performance - Go Optimization Guide](https://goperf.dev/01-common-patterns/)
 
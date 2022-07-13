---
date created: 2021-09-01 22:08
date modified: 2021-11-28 21:49
title: go  内置工具
---
### 输出一个指定的函数的 SSA
```go
func makeSlice() {  
 s := make([]interface{},0,3)  
 s = append(s,1,"3",'a')  
}
```

```bash
GOSSAFUNC=makeSlice go run main.go
```

### 逃逸分析命令
```bash
go tool compile -m -l xxx.go

go build -gcflags "-m -l" xxx.go
```

### 获取 go 的汇编代码
```bash
go tool compile -NlS  xx.go
```
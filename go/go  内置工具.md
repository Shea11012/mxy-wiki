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

### 获取go的汇编代码
```bash
go tool compile -NlS  xx.go
```
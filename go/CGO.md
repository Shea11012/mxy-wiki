---
date created: 2021-11-30 21:22
date modified: 2022-03-26 14:50
title: CGO
---
## import "C"

### 内联 C 代码

`import "C"` 表示使用 CGO 特性，紧跟在这行语句前面的注释是一种特殊语法，里面包含正常的 C 语言代码。

```go
/*
#include <stdio.h>

void printint(int v) {
	printf("printint:%d\n",v);
}
*/
import "C"

func main() {
  v := 42
  C.printint(C.int(v))
}
```

### 独立 C 源代码

```c
/** 目录结构

example/
 hello.c
 hello.h
 main.go
 
 **/

// foo.h
void Hello(const char* s);

// foo.c
void Hello(const char* s) {
    puts(s);
}

// main.go
// #include "foo.h"
import "C"
func main() {
    C.Hello(C.CString("hello world\n"))
}
 
```

> 引入多个文件后不能直接使用 go run 命令运行，需要使用 go build，go build 会检测 cgo 引用文件，.c .h 等文件也会被一起编译

### 调用外部库

```
目录结构：
main.go
	hello/
		hello.h
		hello.c
```

> 因为 c 代码不在 main.go 下，所以 hello 就需要变成外部库来引用
>
> gcc -c foo.c -o foo.o
>
> ar -crs libfoo.a foo.o

```go
// #cgo CFLAGS: -I./hello
// #cgo LDFLAGS: -L./hello -lhello
// #include "hello.h"
import "C"

func main() {
    C.Hello(C.CString("hello world\n"))
}
```



`import "C"` 语句前的注释中可以通过 `#cgo` 语句设置编译阶段和链接阶段的相关参数。

- 编译阶段的参数主要用于定义相关宏和指定头文件检索路径，CFLAGS
- 链接阶段的参数主要指定库文件检索路径和要链接的库文件，LDFLAGS

```go
//#cgo CFLAGS: -DPNG_DEBUG=1 -I./include
//#cgo LDFLAGS: -L/usr/local/lib -lpng
//#include <png.h>
import "C"
```

`-D`：定义宏 PNG_DEBUG 值为 1

`-I`：定义头文件包含的检索目录，c 头文件检索可以是相对路径

`-L`：指定了链接时库文件检索目录，库文件检索必须是绝对路径

`-l`：指定链接时需要链接的库名称

如果引用的头文件或者库在系统默认的目录下，（/usr/include，/usr/local/include，/usr/lib，/usr/local/lib）则可以不用指定目录

**cgo 可以添加限制平台的参数**，具体限制参考 [hdr-build-constraints](https://golang.org/cmd/go/#hdr-Build_constraints)，同时使用 `${SRCDIR}` 替代源码目录

```
// #cgo linux CFLAGS: -I$./hello
// #cgo darwin CFLAGS: -I$./hello
// #cgo LDFLAGS: -L${SRCDIR}/hello -lhello
```



### references

[go-cgo-c](https://bastengao.com/blog/2017/12/go-cgo-c.html)


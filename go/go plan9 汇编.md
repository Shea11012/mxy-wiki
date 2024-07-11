---
date created: 2021-11-30 21:22
date modified: 2024-05-29 05:38
title: go plan9 汇编
---

## Plan9

Go 语言使用 plan9 汇编

### 获取 go 的汇编方式

- `go build -gcflags='-N -l -S' main.go`

- go tool compile -N -l xx.go # 编译程序
  go tool objdump xx.o	# 使用 objdump 反汇编出代码
  go tool objdump -s {symbol} xx.0 # 可以输出指定的 symbol
  
- `go build -gcflags -S xx.go`

### 常量

常量总是 64 位无符号整数

### 寄存器

| AMD64  | GO 汇编 |
| ------ | ------ |
| rax    | AX     |
| rbx    | BX     |
| rcx    | CX     |
| rdx    | DX     |
| rsi    | SI     |
| rdi    | DI     |
| rbp    | BP     |
| rsp    | SP     |
| r8~r15 | R8~R15 |
| rip    | PC     | 

### 符号

go 预定义了一些伪寄存器，这些伪寄存器适用所有的架构

- FP：Frame pointer（栈帧指针）：arguments and locals
- PC：Program counter（程序计数器）：jumps and branches
- SB：Static base pointer（静态指针）：global symbols
- SP：Stack pointer（栈顶指针）：top of stack

#### SB

> 用户定义的符号都会按照偏移量写在 FP 和 SB 中
>
> SB 可以被看做是起始内存，foo(SB) foo 是一个内存地址中的别名。经常当做全局函数和变量。
>
> foo<>(SB) 使这个 foo 仅在当前的 source file 可见
>
> foo+4(SB) 指从 foo 这个地址偏移 4 个字节

#### FP

> FP 是一个虚拟的栈帧指针，用来引用函数参数
>
> 0(FP) 是函数的第一个参数，8(FP) 是函数的第二个参数 ( 64 bit 机器)
>
> 如果这样调用 0(FP) 是会报错，它们需要一个具体的名字来使用 first_arg+0(FP) 表示第一个函数的地址。
>
> FP 总是一个伪地址即使是在一个具有 FP 的架构上。

#### SP

> SP 一个虚拟的栈顶指针，指向一个局部栈帧和参数。因为一直指向栈顶所以在使用时应该是，x-8(SP),y-4(SP)
>
> 在一个具有 SP 的架构上，SP 的名字前缀将伪寄存器与实际寄存器分开
>
> - x-8(SP) 表示一个伪寄存器
> - -8(SP) 实际寄存器
>
> 一般 SP 和 PC 都是一个物理寄存器的别名。在 go 汇编中区别对待了它们。go 中使用需要一个 symbol，first_arg+0(FP)，而为了访问实际寄存器上的值，需要使用 R 。如 R13，R15

#### TLS

> 由 runtime 维护的伪寄存器，保存了指向当前 g 的指针，这个 g 的数据结构会跟踪 goroutine 运行时的所有状态值。

### 指令

#### TEXT 

text 指令的最后一条指令必须是 jump，ret 是一个伪指令。

如果不添加，linker 会自动添加一个 jump-to-itself 的指令

```
TEXT runtime·profileloop(SB),NOSPLIT,$8
	MOVQ	$runtime·profileloop1(SB), CX
	MOVQ	CX, 0(SP)
	CALL	runtime·externalthreadhandler(SB)
	RET
```

栈帧大小即是参数大小，通过 - 分隔。如：$24-8 表示函数有 24 byte 和参数有 8 byte，都存在于调用者的栈帧中。

如果没有 NOSPLIT 在 text 中，则参数大小必须指定

符号名字使用 *·* 分隔组件和名字。如 `$runtime·profileloop1(SB)` 表示这个函数会调用来自 runtime 包的一个变量名为 profileloop 的变量

#### DATA

全局变量被 DATA 和 GLOBL 指令定义

每个 DATA 指令都会初始化一块对应的内存，没有明确初始化的内存会被清零

```assembly
DATA	symbol+offset(SB)/width, value
# symbol 在 SB 的 offset 上，value 的大小是 width

DATA divtab<>+0(SB)/4, $1
```

#### GLOBL

```assembly
GLOBL divtab<>(SB), RODATA, $64

GLOBL runtime·tlsoffset(SB), NOPTR, $4
runtime·tlsoffset，4 byte，不包含指针，隐式清零
```

可能会有一个或者两个指令，如果有两个，第一个表示位掩码标记。

它们的被定义在 `#include textflag.h` 中：

- DUPOK = 2

  > 在一个二进制文件中具有多个重复符号实例是合法的。由链接器选择使用其中的一个

- NOSPLIT = 4

  > 针对 TEXT items。不会插入前导检查堆栈是否拆分。这个栈调用可以加上任何参数，必须对应栈顶段剩余空间。经常用来保护栈代码划分。

- RODATA = 8

  > 针对 DATA 和 GLOBL。因为这个数据不包含指针所以不需要 gc 扫描

- WRAPPER = 32

  > 针对 TEXT。这是一个包装函数，不能被作为一个 disabling recover

- NEEDCTXT = 64

  > 针对 TEXT。这是一个闭包用来作为上下文寄存器

- LOCAL = 128

  > 该符号位于本地的动态共享库

- TLSBSS = 256

  > 针对 DATA 和 GLOBL。将数据放入线程中存储。

- NOFRAME = 512

  > 针对 TEXT。不要插入指令去分配栈帧并保存/恢复返回地址。仅在函数声明了一个大小为 0 的栈帧。

- TOPFRAME = 2048

  > 针对 TEXT。函数在调用堆栈顶部。回溯应该在此函数停止。

#### PCDATA

> PCDATA tableid, tableoffset 第一个参数为表格的类型，第二个为表格的地址
>
> 类型：
>
> - PCDATA_StackMapIndex
> - PCDATA_InlTreeIndex 主要用于内联函数的表格
>
> PCDATA 表格主要包含文件路径、行号和函数信息

#### FUNCDATA

> FUNCDATA tableid, tableoffset 表格类型，表格地址
>
> 类型：
>
> - FUNCDATA_ArgsPointerMaps 表示函数参数的指针信息
> - FUNCDATA_LocalsPointermaps 表示局部指针信息表
> - FUNCDATA_InlTree 表示被内联展开的指针信息表
>
> FUNC 表格，可以让 Go 的垃圾回收期跟踪全部指针的生命周期，同时根据指针指向的地址是否在被移动的栈范围来确定是否要进行指针移动

### interacting with Go Types and constants

如果一个包有 `.s` 文件，则在编译时会直接调用 `go_asm.h` ，然后这个 `.s` 文件会被 `#include`

### 示例

```go
//go:noinline
func add(a, b int32) (int32,bool) {
    return a + b, true
}

func main() {
    add(10, 32)
}
```

编译后汇编如下：

```assembly
# 编译
# GOOS=linux GOARCH=amd64 go tool compile -S add.go

"".add STEXT nosplit size=20 args=0x10 locals=0x0 funcid=0x0
        0x0000 00000 (.\main.go:4)      TEXT    "".add(SB), NOSPLIT|ABIInternal, $0-16
        0x0000 00000 (.\main.go:4)      FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (.\main.go:4)      FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (.\main.go:5)      MOVL    "".b+12(SP), AX
        0x0004 00004 (.\main.go:5)      MOVL    "".a+8(SP), CX
        0x0008 00008 (.\main.go:5)      ADDL    CX, AX
        0x000a 00010 (.\main.go:5)      MOVL    AX, "".~r2+16(SP)
        0x000e 00014 (.\main.go:5)      MOVB    $1, "".~r3+20(SP)
        0x0013 00019 (.\main.go:5)      RET

"".main STEXT size=66 args=0x0 locals=0x18 funcid=0x0
        0x0000 00000 (.\main.go:8)      TEXT    "".main(SB), ABIInternal, $24-0
        0x000f 00015 (.\main.go:8)      SUBQ    $24, SP
        0x0013 00019 (.\main.go:8)      MOVQ    BP, 16(SP)
        0x0018 00024 (.\main.go:8)      LEAQ    16(SP), BP
        0x001d 00029 (.\main.go:8)      FUNCDATA        $0, 
        0x001d 00029 (.\main.go:8)      MOVQ    $137438953482, AX
        0x0027 00039 (.\main.go:8)      MOVQ    AX, (SP)
        0x002b 00043 (.\main.go:8)      PCDATA  $1, $0
        0x002b 00043 (.\main.go:8)      CALL    "".add(SB)
        0x0030 00048 (.\main.go:8)      MOVQ    16(SP), BP
        0x0035 00053 (.\main.go:8)      ADDQ    $24, SP
        0x0039 00057 (.\main.go:8)      RET
```

#### add 函数

```assembly
0x0000 00000 (.\main.go:4)      TEXT    "".add(SB), NOSPLIT|ABIInternal, $0-16
```

- `0x0000 00000`：当前指令的偏移，相对于当前函数的开始

- `TEXT    "".add(SB)`：text 指令声明 `"".add` 符号，`""` 当链接时会被当前包名替换 `main.add`，最终链接进二进制包中

- `SB`：伪寄存器，`"".add(SB)` 通过链接器计算，从程序的起始空间开始定位常量和函数的位置，它是一个绝对地址。

  ```
  # 通过 objdump 可以看到程序的函数的绝对地址
  objdump -j .text -t add | grep "main.add"
  ```

  

### Links

[go asm](https://golang.org/doc/asm)

[go-internals](https://github.com/teh-cmc/go-internals)

[plan9-assembly](https://mioto.me/2021/01/plan9-assembly/)


## gdb 调试工具
调试可执行文件
```shell
gdb program
```

### 运行
- run（r）：运行程序，程序会运行至断点处
- continue（c）：继续执行至下一个断点处
- next（n）：单步跟踪程序，当遇到函数时，不会进入函数体，而是直接调用函数
- step（s）：单步调试，遇到函数会进入函数体内
- until：运行程序直至退出循环体
- until+行号：运行至某行
- finish：运行程序，直到当前函数完成返回，并打印函数返回时的堆栈地址和返回值及参数值等信息
- call：调用程序中可见的函数，并传递参数; `call test(55)`
- quit（q）：退出gdb

### 设置断点
- break n（b n）：在第n行处设置断点，可以带上代码路径和代码名称 `b test.go:14`
- b fn1 if a > b：条件断点设置
- break func：在函数func() 入口处设置断点，`break test`
- delete n：删除第n个断点
- disable n：暂停第n个断点
- enable n：开启第n个断点
- clear n：清除第n行断点
- info b：显示当前程序的断点设置情况
- delete breakpoints：清除所有断点

### 查看源代码
- list（l）：列出程序的源代码，默认每次显示10行
- list 行号：显示当前文件以行号为中心的前后10行代码
- list 函数名：显示函数名的源代码
- list：接着上次list命令，继续输出下面的内容

### 打印表达式
- print（p）：表达式可以是任何当前正在被测试程序的有效表达式
- print a：显示a的值
- print ++a：将a的值加1显示
- print name：显示name的值
- print test(22)：将22作为参数调用 test 函数
- print test(a)：将变量a作为参数调用 test 函数
- display 表达式：使用display设置一个表达式后，将在每次单步进行指令后，紧接着输出被设置的表达式及值。如：`display a`
- watch 表达式：设置一个监视点，一旦被监视的表达式值改变，gdb将终止正在被调试的程序
- whatis：查询变量或函数
- info func：查询函数
- info locals：显示当前堆栈页的所有变量

### 查询运行信息
- where/bt：当前运行堆栈列表
- bt backtrace：显示当前调用堆栈
- up/down：改变堆栈显示深度
- set args：指定运行时参数
- show args：查看设置好的参数
- info program：查看程序是否在运行，进程号，被暂停的原因

### 分割窗口
- layout：分割窗口
- layout src：显示源代码窗口
- layout asm：显示反汇编窗口
- layout regs：显示源代码/反汇编和CPU寄存器窗口
- layout split：显示源代码和反汇编窗口
- ctrl+L：刷新窗口
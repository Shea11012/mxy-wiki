---
date created: 2021-10-27 00:48
date modified: 2022-03-29 15:11
title: process
---
process id，进程数量的最大上限被定义在 `/proc/sys/kernel/pid_max`

### 进程的内存布局
每个进程是由多个部分组成，通常被称为 **segments**。有如下 segments：
- **text segment**: 包含了程序的机器级指令。只读，不会被程序的指令修改值。因为一些进程可能运行同样的程序，text segment 是可共享的，可以被映射进这些程序的虚拟内存地址空间。
- **initialized data segment**: 包含初始化的全局和静态变量。当程序载入内存时会从可执行文件中读取这些值。
- **uninitialized data segment (bss segment)**: 包含未被初始化的全局和静态变量。
- **stack**：动态增长和缩容的 segment。会为每一个被调用的函数分配栈帧。一个帧存储了函数本地变量，参数和返回值。
- **heap**：在运行时动态分配。堆顶被称为程序中断 (program break)。
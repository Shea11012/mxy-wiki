---
tags: ["linux"]
date created: 2021-10-27 00:48
date modified: 2021-10-27 00:48
title: system call
---
# system call
- 系统调用会改变进程状态，从用户态进入内核态，这样 CPU 才可以保护内核空间
- 一组系统调用是固定的
- 每个系统调用会将数据从用户态传入内核态
- 所有的系统调用会以同样的方法进入内核，系统调用数字会写入 %eax 寄存器中

`ldd <execute file>` 会列出可执行文件的动态库依赖

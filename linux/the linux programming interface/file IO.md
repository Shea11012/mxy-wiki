---
tags: ["linux","IO"]
date created: 2021-10-27 00:48
date modified: 2022-03-29 15:11
title: file I/O
---
# file I/O
所有打开的文件都使用文件描述符作为引用。
文件描述符包含所有文件类型：pipe、FIFO、sockets、terminals、devices、regular files。每个进程有自己的一组文件描述符。
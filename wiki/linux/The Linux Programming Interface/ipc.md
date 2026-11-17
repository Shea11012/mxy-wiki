---
title: "IPC（interprocess communication)"
created: 2026-01-13
updated: 2026-07-22
status: needs-review
type: source-summary
tags:
  - linux
  - linux/system-programming
  - linux/process
collection: "The Linux Programming Interface"
collection_kind: book
collection_order: 6
collection_home: "The Linux Programming Interface（总览）"
confidence: INFERRED
---

# IPC（interprocess communication)
linux 中每个进程都是相互独立的，为了实现进程间通信和同步操作，类 UNIX 提供了一组 IPC 工具。
- signals： 用来表示某一个事件被触发
- pipes：在两个进程间传递消息
- sockets：通过网络实现在进程间通信
- file locking：允许一个进程锁住原文件，阻止另一个进程读或更新操作
- message queues：用来在进程间交换消息
- semaphores：用来同步进程间操作
- shared memory：允许两个或更多的进程共享一块内存，当一个进程改变了共享内存的内容，其他进程能立即看见改变
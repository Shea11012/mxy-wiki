---
title: "file I/O"
created: 2026-01-13
updated: 2026-07-22
status: needs-review
type: source-summary
tags:
  - linux
  - linux/system-programming
collection: "The Linux Programming Interface"
collection_kind: book
collection_order: 2
collection_home: "The Linux Programming Interface（总览）"
confidence: INFERRED
---

# file I/O
所有打开的文件都使用文件描述符作为引用。
文件描述符包含所有文件类型：pipe、FIFO、sockets、terminals、devices、regular files。每个进程有自己的一组文件描述符。
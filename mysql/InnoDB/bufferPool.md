---
tags: []
date created: 2022-03-29 15:09
date modified: 2023-04-05 04:32
title: bufferPool
---
InnoDB 存储引擎通过 buffer pool，提高数据库的读写性能，buffer pool 的作用：
- 当读取数据时，数据存在 buffer pool 中，直接读取 buffer pool 数据，否则去磁盘读取
- 当修改数据时，先修改 buffer pool 中数据所在的页，然后将其页设置为脏页，最后通过刷盘机制将脏页写入到磁盘

## Buffer Pool

buffer pool 申请一片连续的内存，然后按照默认的 16KB 大小划分出页，buffer pool 中的页就叫缓存页。
![[attachments/Pasted image 20230405043056.png]]

### free 链表

刚完成初始化的 buffer pool，所有的缓冲页都是空闲的，所以每一个缓冲页对应的控制块都会加入到 free 链表中
当需要一个缓冲页时，会有一个以表空间号 + 页号作为 key，缓冲页地址作为 val 的 hash 表，通过这个 hash 表就可以拿到空闲缓冲页

### flush 链表

当修改了某个缓冲页的数据，它与磁盘上的页不一致了，就被称为脏页（dirty page）
每次修改缓冲页并不是立即刷新到磁盘，而是会使用一个 flush 链表，将被修改过的缓冲页对应控制块作为一个节点加入到链表中
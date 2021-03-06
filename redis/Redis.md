---
date created: 2021-12-03 20:20
date modified: 2022-04-10 00:56
title: Redis
---

- Redis 是一个开源的支持多种数据类型的 key=>value 的存储数据库。支持字符串、列表、集合、有序集合、哈希五种类型
- Redis 是一个开源，基于内存的结构化数据存储媒介，可以作为数据库、缓存服务或消息服务使用。
- Redis 支持多种数据结构，包括字符串、哈希表、链表、集合、有序集合、位图、Hyperloglogs
- Redis 具备 LRU 淘汰、事务实现以及不同级别的硬盘持久化等能力，并且支持副本集和通过 Redis Sentinel 实现的高可用方案，同时还支持 Redis Cluster 实现的数据自动分片能力
- Redis 主要功能都基于单线程模型实现，表示 Redis 使用一个线程来服务所有的客户端请求，同时 Redis 采用了非阻塞式 IO，并优化各种命令的算法时间复杂度
- Redis 是线程安全（单线程），所有的操作都是原子的，不会因为并发产生数据异常
- Redis 的速度非常快（非阻塞式 IO，且大部分命令的时间复杂度都是 $O(1)$
- 使用高耗时 Redis 命令很危险，会占用唯一的一个线程大量处理时间，导致所有的请求都会被拖慢，例如：时间复杂度为 $O(N)$ 的 keys 命令，严格禁止在生产环境中使用
- 纯内存访问，Redis 将数据放在内存中，内存响应时间短
- 非阻塞 I/O，Redis 使用 epoll 模型，加上 Redis 自身的事件处理模型将 epoll 中的连接、读写、关闭都转换为事件，不会在网络 I/O 上浪费过多的时间
- 单线程避免了线程切换和竞态产生的消耗（现在多了一个线程处理网络 io，会有些许的切换消耗）


## redis 和 memcache 区别

1. Redis 不仅仅支持简单的 k/v 类型的数据，同时还提供 list，set，hash 等数据结构的存储。 

2. Redis 支持数据的备份，即 master-slave 模式的数据备份。

3. Redis 支持数据的持久化，可以将内存中的数据保持在磁盘中，重启的时候可以再次加载进行使用。


![[数据类型]]

![[内部结构]]

![[持久化]]

![[内存溢出和删除策略]]

![[缓存设计]]


![[注意事项]]

---
date created: 2021-12-03 20:20
date modified: 2022-03-09 19:04
title: zookeeper
---
#### zookeeper 模型

zookeeper 使用树形结构，类似 Linux 的文件系统。

树的每个节点叫 znode。每个 znode 都可以保存数据。每个节点都有一个版本号，从 0 开始计数。

znode 只支持全量的读取和写入

#### znode 分类

持久性 znode，在创建之后即使发生 zookeeper 集群宕机或者 client 宕机也不会丢失。

临时性 znode，client 宕机或者 client 在指定的 timeout 时间内没有给 zookeeper 集群发消息，znode 会消失。

> znode 可以是顺序性的。每一个顺序性的 znode 关联一个唯一的单调递增整数（单调递增整数是 znode 名字的后缀）。

持久顺序性 znode：具备持久性，名字具备顺序性。

临时顺序性 znode：具备临时性，名字具备顺序性。



使用 zookeeper 锁

```
# a
1 create -e /lock # 创建一个临时锁
4 可以主动释放锁或者客户端退出

# b
2 create -e /lock # 此时锁被持有
3 stat -w /lock  # 观察锁状态
```


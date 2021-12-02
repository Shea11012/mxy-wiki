#### zookeeper 模型

zookeeper 使用树形结构，类似Linux的文件系统。

树的每个节点叫znode。每个znode都可以保存数据。每个节点都有一个版本号，从0开始计数。

znode只支持全量的读取和写入

#### znode 分类

持久性znode，在创建之后即使发生zookeeper集群宕机或者client宕机也不会丢失。

临时性znode，client宕机或者client在指定的timeout时间内没有给zookeeper集群发消息，znode会消失。

> znode可以是顺序性的。每一个顺序性的znode关联一个唯一的单调递增整数（单调递增整数是znode名字的后缀）。

持久顺序性znode：具备持久性，名字具备顺序性。

临时顺序性znode：具备临时性，名字具备顺序性。



使用zookeeper锁

```
# a
1 create -e /lock # 创建一个临时锁
4 可以主动释放锁或者客户端退出

# b
2 create -e /lock # 此时锁被持有
3 stat -w /lock  # 观察锁状态
```


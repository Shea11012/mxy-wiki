---
date created: 2021-12-03 20:20
date modified: 2022-05-04 04:12
title: kafka
---
## Kafka

支持多分区、多副本，分布式流平台（a distributed streaming plattorm)，也是 **基于发布订阅模式的消息引擎系统**

### 术语
**event 消息：** 表示一个记录或消息，包含：key、value、timestamp 和可选的 metadata

**topic 主题：** 承载消息的容器，多用来区分具体的业务

**partition 分区：** 一个有序不变的消息序列，每个主题是可以有多个分区

**offset 消息位移：** 分区中每条消息的位置信息，是一个单调递增的值

**replica 副本：** 同一条消息能够被拷贝到多个地方以提供数据冗余。副本分为领导者副本和追随者副本。副本是在分区层级下的，每个分区可配置多个副本。

**producer 生产者：** 向主题发布消息

**consumer 消费者：** 订阅主题消费消息

**consumer offset 消费者位移：** 每个消费者都有自己的消费者位移

**consumer group 消费者组：** 多个消费者组成一个组，同时消费多个分区提高吞吐

**rebalance 重平衡：** 消费者组内某个消费者实例挂掉后，其他消费者实例自动重新分配订阅主题分区的过程


### 消费者位移
`__consumer_offsets` kafka 内部的位移主题，用于保存 kafka 消费者的位移信息。

位移主题的 key 内容为<group_id、主题名、分区号>，value 则除了保存位移值，还保存了一些其他的元数据，如时间戳和用户自定义数据。

位移的消息格式：
- 保存位移值的消息格式
- 保存 consumer_group 信息的消息格式
- 删除 group 过期位移甚至是删除 group 的消息格式


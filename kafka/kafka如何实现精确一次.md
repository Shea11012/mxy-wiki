---
tags: []
date created: 2023-05-20 03:44
date modified: 2023-05-20 03:56
---
- 最多一次（at most once）：消息可能丢失，但绝不会被重复发送
- 至少一次（at least once）：消息不会丢失，但可能被重复发送
- 精确一次（exactly once）：消息不会丢失，也不会被重复发送

kafka 默认提供的是 at least once

>kafka 通过幂等性（idempotence）和事务（transactions）实现了精确一次

## 幂等性（idempotence）

kafka 在 `0.11.0.0` 版本引入了幂等性，通过设置参数 `enable.idempotence = true`，生产者自动变为幂等性生产者。

kafka 为了实现幂等性，引入了 producerID 和 SequenceNumber。
producer 需要做两件事：
1. 初始化时向 broker 申请一个 producerID
2. 为每条消息绑定一个 SequenceNumber
在 kafka broker 收到消息后会以 producerID 为单位存储 SequenceNumber，即使 Producer 重复发送，broker 也会将其过滤掉。
这种实现有很大的限制：
- **只能保证单分区的幂等性**，因为 SequenceNumber 是以 Topic+Partition 为单位单调递增的，如果一条消息被发送到多个分区必然会分配到不同的 SequenceNumber，导致重复
- **只能实现单会话幂等性**，当重启 Producer 进程后，这种幂等性就丧失了，因为重启后会重新分配一个 ProducerID

## 事务（transaction）

事务型 producer 能够保证消息原子性的写入到多个分区，也不会被进程重启所影响。
设置事务型Producer：
- 开启 `enable.idempotence = true`
- 设置Producer参数 `transacal.id`，最好赋值一个有意义的名字

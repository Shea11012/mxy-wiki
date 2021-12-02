## Kafka

支持多分区、多副本，分布式流平台（a distributed streaming plattorm)，也是 **基于发布订阅模式的消息引擎系统**



#### event 消息

> 表示一个记录或者消息，包含：key、value、timestamp 和 可选的metadata。

#### producers 生产者 

> 往kafka 某个主题内publish（写）events

#### consumers  消费者

> 订阅主题消息 subscribe（读和处理）events，处理生产者产生的消息

#### topics 主题

> 类似一个文件系统的文件夹，events都被存放在topics中。
>
> topic 可以拥有多个producers和consumers
>
> topics 是会被分区的（partitioned），表示一个topic会分散在不同的brokers中。（简单理解就是kafka根据event的key做了一致性哈希）。同时kafka保证了在读取一个topic分区的events时是会按照它们的写入顺序。

#### partition 分区

> 主题可以被分为若干个分区，同一个主题中的分区可以不在一个机器上，有可能会部署在多个机器上，由此实现kafka的伸缩性。
>
> 单一主题中单一分区有序，但是无法保证主题中所有分区有序

#### group 消费者群组

> 消费者群组由一个或多个消费者组成的群体

#### offset 偏移量

> 元数据，是一个不断递增的整数值，用来记录消费者发生重平衡时的位置，以便用来恢复数据。

#### broker

> 一个独立的 kafka 服务器就被称为 broker，broker 接收来自生产者的消息，为消息设置偏移，并提交消息到磁盘保存。

#### broker 集群

> broker 集群由一个或多个 broker 组成，每个集群都有一个 broker 同时充当了 集群控制器的角色。（根据分布式算法选举出来）

#### replica 副本

> kafka 中消息的备份叫做 replica。副本的数量可配置。
>
> kafka 定义了两类副本：
>
> - 领导者副本（leader），对外提供服务
> - 追随者副本（follower），被动跟随

#### rebalance 重平衡

> 消费者组内某个消费者实例挂掉后，其他消费者实例自动重新分配订阅主题分区的过程。rebalance 是 kafka 消费者端实现高可用的重要手段。

### kafka 消息队列

- 点对点模式（消息队列）：一个消费者群组消费一个主题中的消息
- 发布订阅模式：一个主题中的消息被多个消费者群组共同消费

消费者组

<img src="kafka.assets/1515111-20191128124553311-1167412207.png" style="zoom:50%;" />

> 如果应用需要读取全量消息，那么为该应用设置一个消费组；如果该应用消费能力不足，那么可以考虑在这个组增加消费者

消费者重平衡

![](kafka.assets/1515111-20191128124603277-431025951.png)



### 常用命令

| 说明          | 命令                                                         |
| ------------- | ------------------------------------------------------------ |
| 创建 topic    | kafka-topics --create --bootstrap-server :9092 --replication-facotr 1 --partitions 1 --topic test-topic |
| 查看 topic    | kafka-topics --list --bootstrap-server :9092                 |
| 查看指定topic | kafka-topics  --describe --topic test-topic --bootstrap-server :9092 |
| 产生消息      | kafka-console-producer --topice test-topic --bootstrap-server :9092 |
| 消费消息      | kafka-console-consumer --topic test-topic --from-beginning --bootstrap-server :9092 --partition 0 |
|               |                                                              |
|               |                                                              |
|               |                                                              |
|               |                                                              |


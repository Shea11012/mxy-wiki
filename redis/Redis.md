## redis

- Redis是一个开源的支持多种数据类型的key=>value的存储数据库。支持字符串、列表、集合、有序集合、哈希五种类型
- Redis 是一个开源，基于内存的结构化数据存储媒介，可以作为数据库、缓存服务或消息服务使用。
- Redis 支持多种数据结构，包括字符串、哈希表、链表、集合、有序集合、位图、Hyperloglogs
- Redis 具备 LRU 淘汰、事务实现以及不同级别的硬盘持久化等能力，并且支持副本集和通过 Redis Sentinel 实现的高可用方案，同时还支持 Redis Cluster 实现的数据自动分片能力
- Redis 主要功能都基于单线程模型实现，表示 Redis 使用一个线程来服务所有的客户端请求，同时 Redis 采用了非阻塞式 IO，并优化各种命令的算法时间复杂度
- Redis 是线程安全（单线程），所有的操作都是原子的，不会因为并发产生数据异常
- Redis  的速度非常快（非阻塞式 IO，且大部分命令的时间复杂度都是 $O(1)$
- 使用高耗时 Redis 命令很危险，会占用唯一的一个线程大量处理时间，导致所有的请求都会被拖慢，例如：时间复杂度为 $O(N)$ 的 keys 命令，严格禁止在生产环境中使用

![[redis.assets/redis.png]]

### redis 和memcache区别

1. Redis不仅仅支持简单的k/v类型的数据，同时还提供list，set，hash等数据结构的存储。 

2. Redis支持数据的备份，即master-slave模式的数据备份。

3. Redis支持数据的持久化，可以将内存中的数据保持在磁盘中，重启的时候可以再次加载进行使用。



### redis五种类型

#### 字符串 String

> - set：为一个 key 设置 value，可以配置 EX/PX 参数指定 key 的有效期，$O(1)$
> - get：获取某个 key 对应的 value，$O(1)$
> - getset：为一个 key 设置 value ，并返回该 key 的原 value，$O(1)$
> - mset：为多个 key 设置 value，$O(N)$
> - msetnx：同 mset，如果指定的 key 中有任意一个已存在，则不进行任何操作，$O(N)$
> - mget：获取多个 key 对应的 value，$O(N)$
> - incr：将 key 对应的 value 自增 1，并返回自增后的值（只对可以转换为整型 string 数据起作用）$O(1)$
> - incrby：将 key 对应的 value 值自增指定的整型值，并返回自增的值（只对可以转换为整型 string 数据起作用）$O(1)$
> - decr/decrby：同 incr/incrby，自增改为自减

### 列表 List

链表型的数据结构，list 也支持在特定的 index 上插入和读取元素的功能，但时间复杂度较高 $O(N)$

> - lpush：向指定的 list 的左侧插入 1 个或多个元素，返回插入后的长度 $O(N)$
> - rpush：同 lpush
> - lpop：从指定的 list 左侧移除一个元素并返回，$O(1)$
> - rpop：同 lpop
> - lpushx、rpushx：如果操作的 key 不存在，则不会进行任何操作
> - llen：返回指定 list 的长度 $O(1)$
> - lrange：返回指定 list 中指定范围的元素，$O(N)$ ，应尽可能的控制一次获取的元素数量，一次获取大范围的 list 元素会导致延迟
> - lindex：返回指定 list 指定 index 上的元素，如果 index 越界，返回 nil。index 数值是回环的，即 -1 代表最后一个位置，-2 代表倒数第二，$O(N)$
> - lset：将指定 list 指定 index 上的元素设置为 value，如果 index越界则返回错误，$O(N)$
> - linsert：向指定 list 中指定元素之前、之后插入一个新元素，并返回操作后的 list 长度，如果指定的元素不存在，返回 -1。如果指定 key 不存在，则不会进行任何操作，$O(N)$
> - 

#### 集合 Set

> - sadd：向指定 set 中添加 1 个或多个 member，如果指定 set 不存在，则会自动创建一个，$O(N)$
> - srem：从指定 set 中移除 1 个或多个 member，$O(N)$
> - srandmember：从指定 set 中随机返回一个或多个 member $O(N)$
> - spop：从指定 set 中随机移除并返回 count 个 member $O(N)$
> - scard：返回指定 set 中的 member个数，$O(1)$
> - sismember：判断指定的 value 是否存在指定 set 中，$O(1)$
> - smove：将指定 member 从一个 set 移至另一个 set
> - smembers：返回指定 hash 中所有的 member，$O(N)$
> - sunion、sunionstore：计算多个 set 的并集并返回、存储至另一个 set 中，$O(N)$
> - sinter、sinterstore：计算多个 set 的交集并返回、存储至另一个 set 中，$O(N)$
> - sdiff、sdiffstore：计算一个 set 或多个 set 的差集并返回、存储至另一个 set $O(N)$

#### 有序集合 Sorted Set

sorted set 中的每个元素都需要指派一个分数（score），sorted set 会根据 score 对元素进行升序排序。如果多个 member 拥有相同的 score，则以字典序进行升序排序，适用场景：实现排名

> - zadd：向指定 sorted set 中添加一个或多个 member，$O(M_log(N))$ M 为添加的 member 数量，N 为 sorted set 中 member 数量
>
> - zrem：从指定 sorted set 中删除一个或多个 member，$O(M_log(N))$ M 为删除的 member 数量，N 为 sorted set 中的 member 数量
> - zcount：返回指定 sorted set 中指定 score 范围内的 member 数量，$O(log(N))$
> - zcard：返回指定 sorted set 中的 member 数量，$O(1)$
> - zscore：返回指定 sorted set 中指定 member 的 score，$O(1)$
> - zrank、zrevrank：返回指定 member 在 sorted set 中的排名，zrank 返回升序排序的排名，zrevrank 返回降序，$O(log(N))$
> - zincrby：同 incrby
> - zrange、zrevrange：返回指定 sorted set 中指定排名范围内的所有 member，zrange 为按 score 升序排序，zrevrange 按 score 降序排序，$O(log(N)+M)$  M 为本次返回的 member 数
> - zrangebyscore、zrevrangebyscore：返回指定 sorted set 中指定 score 范围内所有 member，返回结构以升序、降序排序，min 和 max 可以指定为 -int 和 +int ，代表返回所有的 member，$O(log(N)+M)$
> - zremrangebyrank、zremrangebyscore：移除 sorted set 中指定排名范围、指定 score 范围内的所有 member，$O(log(N)+M)$

#### 哈希 Hash

Hash 是一种 field-value 型的数据结构，使用用于表现对象类型的数据

> - hset：将 key 对应的 hash 中的 field 设置为 value，如果该 hash 不存在，会自动创建一个，$O(1)$
> - hget：返回指定 hash 中 field 字段的值，$O(1)$
> - hmset、hmget：同 hset 和 hget，可以批量操作同一个 key 下的多个 field，$O(N)$
> - hsetnx：如果 field 已经存在，则不会进行任何操作，$O(1)$
> - hexists：判断指定 hash 中 field 是否存在，存在返回 1 ，不存咋返回 0，$O(1)$
> - hdel：删除指定的 hash 中的 field ，$O(N)$
> - hincrby：同 incrby 命令
> - hgetall：返回指定 hash 中所有的 field-value 对，返回结果为数组，数组中 field 和 value 交替出现 $O(N)$
> - hkeys、hvals：返回指定 hash 中所有 field、value，$O(N)$
>
> tips：hgetall、hkeys、hvals 都会对 hash 进行完整遍历，对于尺寸不可预知的 hash，需要避免使用，改为使用 hscan 遍历

#### 其余常用命令：

> - exists：判断指定 key 是否存在，返回 1 代表存在，0 代表不存在，$O(1)$
> - del：删除指定 key 及其对应的 value，$O(N)$
> - expire、pexpire：为一个 key 设置有效期，单位秒或毫秒，$O(1)$
> - ttl、pttl：返回一个 key 剩余的有效时间，单位秒或毫秒，$O(1)$
> - rename、renamenx：将 key 重命名为 newkey。使用 rename 时，如果 newkey 存在其值会被覆盖；使用 renamenx，则不会进行任何操作，$O(1)$
> - type：返回指定 key 的类型，$O(1)$
> - config get：获得 redis 某配置项的当前值，可以使用通配符 * $O(1)$
> - config set：为 redis 某个配置设置新值 $O(1)$
> - config rewrite：让 redis 重新加载 redis.con 中配置

### 数据持久化

#### RDB

优点：

- 对性能影响最小，redis 在保存 RDB 快照时会 fork 子进程进行，几乎不会影响 redis 处理客户端请求的效率
- 每次快照会生成一个完整的数据快照文件，所以可以使用 bash 脚本按照时间点保存快照
- RDB 文件数据恢复比 AOF 快

缺点：

- 快照是定期生成，Redis crash 时会或多或少丢失一部分的数据
- 如果数据集非常大且 CPU 不够强，redis 在 fork 子进程时会消耗相对较长的时间，会影响此期间的客户端请求

#### AOF

AOF 提供三种 fsync 配置：

> - appendfsync no：不进行 fsync，将 flush 文件的时机交个 OS 决定，速度最快
> - appendfsync always：每写入一条日志就进行一次 fsync 操作，数据安全性最高，速度最慢
> - appendfsync everysync：折中做法，交由后台线程每秒 fsync 一次
>
> 随着不停的记录写操作日志，会出现一些无用的日志，Redis 提供 AOF rewrite 功能，可以重写 AOF 文件，只保留能够把数据恢复到最新状态的最小写操作集
>
> AOF rewrite 可以通过 BGREWRITEAOF 命令触发，可以配置 Redis 定期进行

优点：

- 安全，在启用 appendfsync always 时，任何已写入的数据都不会丢失，使用 appendfsync everysec 也只会丢失 1 秒的数据
- AOF 文件在发生断电等问题时也不会损坏，即使出现某条日志只写入了一半的情况，也可以使用 redis-check-aof 工具修复
- AOF 文件易读，可修改，在进行某些错误的数据清除操作，只要AOF 文件没有 rewrite 就可以把 AOF 文件备份出来，把错误的命令删除，然后恢复数据

缺点：

- AOF 文件比 RDB 文件更大
- 性能消耗比 RDB 高
- 数据恢复速度比 RDB 慢

### Redis 过期删除策略

redis 采用的是定期删除 + 惰性删除策略

- **定时删除**：设置键的过期时间同时创建一个定时器。此策略十分消耗 CPU 时间，影响服务器的吞吐量，一般不采取此策略
- **惰性删除**：获取键时，每次都检查取得的键是否过期，过期则删除。此策略浪费内存，有内存泄漏的风险。
- **定期删除**：redis 默认每个 100ms 检查，是否有过期的 key，有则删除。redis 不是每隔 100ms 将所有的键检查一次，而是随机抽取检查，随机抽查会导致很多键到时间没有删除

### 内存淘汰机制
在 redis.conf 配置 `maxmemory-policy volatile-lru`，配置参数如下：
  - noeviction：当内存不足以容纳新写入数据时，新写入操作会报错
  - allkeys-lru：当内存不足以容纳新写入数据时，在键空间中，移除最少使用的键
  - allkeys-random：当内存不足以容纳新写入数据时，在键空间中，随机移除某个键
  - volatile-lru：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，移除最少使用的键。此种情况一般是把 redis 当缓存，也当作持久化存储时才使用。
  - volatile-random：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，随机移除某个键
  - volatile-ttl：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，有更早过期时间的键优先移除




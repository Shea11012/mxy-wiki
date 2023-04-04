---
tags: []
date created: 2022-03-30 16:59
date modified: 2023-04-05 04:26
title: binlog日志
---
MySQL 在完成一条操作后，Server 层会生成一条 binlog，等事务提交时，会将该事务执行过程中产生的所有 binlog 统一写入到 binlog 文件。

binlog 主要用途：
- 复制：如一主多从
- 掉电恢复

binlog 的查询语句
```mysql
show binlog events [IN 'log-name'] [From pos] [Limit [offset], row_count]
```

## 记录格式

binlog 不同格式，`--binlog-format=`：
- `statement`：生成的 binlog**基于语句的日志**。此时会将一条 SQL 语句完整的记录到 binlog 上，而不管该语句影响了多少记录
- `row`：生成的 binlog**基于行的日志**，会将该语句所改动的全部信息都记录上
- `mixed`：生成的 binlog**基于行的日志**，一般会 **基于语句的日志**，特殊情况下会自动转为 **基于行的日志**

## 刷盘时机

事务执行过程，先把日志写到 `binlog cache` 中，事务提交时，再把 `binlog cache` 写入到 binlog 文件中。

因为一个事务的 binlog 不能被拆开，无论事务有多大，也要确保一次性写入，所以系统给每个线程分配一个块内存作为 `binlog cache`

通过 `binlog_cache_size` 控制单个线程的 binlog cache 大小，如果存储内容超过了这个参数就要暂存到磁盘（swap）

![[Pasted image 20220410014021.png]]

write 将日志写入到文件系统的 page cache，并没有把数据持久化到磁盘
fsync 将数据持久化到磁盘
write 和 fsync 的时机，由参数 `sync_binlog` 控制，默认为 0
- 0，表示每次提交事务都只 write，由系统判断什么时候执行 fsync
- 1，每次提交事务都会执行 fsync，最安全但性能损耗最大
- N(N>1)，表示每次提交事务都 write，但累积 N 个事务后才 fsync，为了提高写入性能，一般会设置 N 为 100~1000 中的某个数值

---
tags:
  - 数据库
  - mysql
date created: 2022-03-29 15:49
date modified: 2023-11-21 15:21
title: redo日志
---

# redo 日志

redo log 保证了事务的持久性。

为了避免修改一个字节就将脏也刷到磁盘，将一个修改操作按照一定的格式记录下来，当服务器发生崩溃后能使用 redo 日志复原结果，这就是 redo 日志。

redo log 写入方式使用追加操作，所以磁盘操作是顺序写。MySQL 写入数据的操作是随机写，因此 redo log 写入磁盘的开销更小。

redo log 也有缓存，**redo log buffer**，每当产生一条 redo log 时，会先写入到 redo log buffer 中，然后再持久化到磁盘。

redo log buffer 大小默认为 16MB，通过 `innodb_log_buffer_size` 调整大小

## 刷盘时机

- log buffer 空间不足时，当 redo 日志占了 log buffer 总容量的 50% 时，执行一次刷盘
- 事务提交时
- 后台线程以每秒一次的频率将 redo 日志进行刷盘
- checkpoint 时

### checkpoint

redo 日志是为了系统崩溃后恢复脏页使用，如果对应的脏页已经刷到磁盘中，即使系统崩溃，在重启后也不会使用该 redo 日志恢复该页面了，所以该 redo 日志就没有存在的必要的，它所占用的磁盘空间就可以被后续的 redo 日志覆盖
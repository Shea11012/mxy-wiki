---
date created: 2022-03-29 15:49
date modified: 2022-03-29 16:46
title: redo日志
---
## redo 日志
为了避免修改一个字节就将脏也刷到磁盘，将一个修改操作按照一定的格式记录下来，当服务器发生崩溃后能使用 redo 日志复原结果，这就是 redo 日志。
## 日志格式
redo 日志通用结构,根据类型的不同，日志的格式也有所不同
![[Pasted image 20220329155422.png]]
- type：日志的类型
- space ID：表空间 ID
- page number：页号
- data：具体内容

type 为 `MLOG_8BYTE` 类型时，表示在页面的某个偏移量处写入 8 字节的数据
![[Pasted image 20220329155954.png]]

type 为 `MLOG_WRITE_STRING` 时，因为是写入字节序列不能确定大小所以需要个 len 字段
![[Pasted image 20220329160416.png]]

## 日志的写入过程
将对底层页面进行一次原子访问的过程称为，Mini-Transaction（MTR）。
一个事务可以包含若干个语句，每一条语句又包含若干个 MTR，每一个 MTR 又包含若干条 redo 日志。
![[Pasted image 20220329162318.png]]

redo 日志并不是直接写入磁盘的，会使用一个 redo log buffer 的连续内存空间，这片内存空间被划分为若干个连续的 redo log block。默认值为 16MB
日志的写入并不是以每生成一个 redo 日志就写入到 redo log 日志中的，而是一个 MTR 为一组，当 MTR 结束时才将这一组 redo 日志写入

### 刷盘时机
- log buffer 空间不足时，当 redo 日志占了 log buffer 总容量的 50% 时，执行一次刷盘
- 事务提交时
- 后台线程以每秒一次的频率将 redo 日志进行刷盘
- checkpoint 时

## checkpoint
redo 日志是为了系统崩溃后恢复脏页使用，如果对应的脏页已经刷到磁盘中，即使系统崩溃，在重启后也不会使用该 redo 日志恢复该页面了，所以该 redo 日志就没有存在的必要的，它所占用的磁盘空间就可以被后续的 redo 日志覆盖
---
tags: []
date created: 2022-03-29 16:49
date modified: 2023-04-04 14:16
title: undo日志
---
undo log，是一种用于撤销回退的日志，在事务没提交前，MySQL 会先记录更新前的数据到 undo log 文件里，当事务回滚时，可以利用 undo log 进行回滚。undo log 提供了原子性保证。

## undo 页面链表

一个事务执行中可能混着执行 INSERT、DELETE、UPDATE 语句，这会产生不同类型的 undo 日志，但是同一个 undo 页面只能存储同一种类型的 undo 日志，所以一个事务的执行过程中可能需要 2 个 undo 页面的链表：insert undo 链表、update undo 链表，且对普通表和临时表的记录改动时所产生的 undo 日志要分别记录，所以一个事务中最多有 4 个以 undo 页面为节点组成的链表。分配策略为：
- 开启事务时，不分配 undo 页面
- 事务执行时向普通表插入记录或执行更新记录主键的操作之后，就会分配一个普通表的 insert undo 链表
- 事务执行时删除或更新普通表中的记录，会分配一个 update undo 链表
- 事务执行时向临时表插入记录或执行更新记录主键的操作，会分配一个临时表的 insert undo 链表
- 事务执行时删除或更新临时表中的记录，会分配一个临时表的 update undo 链表

>[!tip]
>undo log 格式中有一个 roll_pointer 和 trx_id 字段：
>- trx_id 指向事务 ID
>- roll_pointer 指向前一个 undo log，形成链表
>
>同时可以通过 readview+undo log 实现 MVCC

## undo 日志在崩溃恢复时的作用

如在服务器崩溃的而恢复的过程中，首先按照 redo 日志将各个页面的数据恢复到崩溃之前的状态，这样保证了已经提交的事务的持久性。但是那些没有提交的事务写的 redo 日志可能也已经刷盘，此时为了保证事务的原子性，可以通过 undo 日志进行回滚。  
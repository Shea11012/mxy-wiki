---
tags: []
date created: 2023-06-14 13:34
date modified: 2023-06-16 14:39
---

## MySQL 的执行流程

[[../mysql/MySQL执行流程|MySQL执行流程]]

## CPU 升高后，如何排查

- 通过 `show processlist` 或 `show full processlist` 检查正在运行的查询
- 检查慢查询日志：配置慢查询日志，记录查询执行时间超过某个阈值的查询
- 查看 MySQL 的性能模式：使用 `show status` 查看 MySQL 运行状态，其中包括很多性能相关的指标，如 InnoDB 的读写情况、查询缓存的命中情况、索引的使用情况等
- 使用性能诊断工具：如 `performance_schema` 或者借助第三方工具
- 优化查询：根据上面的分析，优化消耗大量 CPU 资源的查询。包括修改查询语句，创建或修改索引，修改表结构，调整 MySQL 的配置。

## 索引类型

[[../mysql/索引#索引类型|索引类型]]

## 回表

[[../mysql/索引#回表|回表]]

## 索引下推

[[../mysql/索引#索引下推|索引下推]]

## 索引失效

[[../mysql/索引失效|索引失效]]

## 锁

[[../mysql/InnoDB/锁|锁]]

## 事务隔离级别

[[../mysql/InnoDB/事务#事务隔离级别|事务隔离级别]]

## MVCC

[[../mysql/InnoDB/事务#MVCC|MVCC]]

## redolog

[[../mysql/InnoDB/redo日志|redo日志]]

## undolog

[[../mysql/InnoDB/undo日志|undo日志]]

## binlog

[[../mysql/binlog日志|binlog日志]]

## 主库和从库数据存在延迟问题如何解决

## 主从复制原理是什么

## 常见的分片算法有哪些

## 分库分表会有什么问题
读扩散

## 分库分表后，数据怎么迁移


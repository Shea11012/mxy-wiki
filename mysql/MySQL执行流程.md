---
tags: []
date created: 2023-03-29 16:30
date modified: 2023-03-29 17:29
---
MySQL 分为两层：**Server 层和存储引擎层**
- Server 层负责建立连接、分析和执行 SQL。包括连接器、查询缓存、解析器、预处理器、优化器、执行器等。
- 存储引擎层负责数据的存储和提取。支持 InnoDB、MyISAM、Memory 等多个存储引擎。

## 连接器

连接到 MySQL 服务器

## 查询缓存

MySQL 解析 SQL 语句，如果是 select 语句，就会先去缓存查找数据，查询缓存以 key-value 形式保存在内存中，key 为 SQL 语句，value 为 SQL 查询结果
但因为对更新比较频繁的表，查询缓存命中率很低，因为只要一个表有更新，那么这个表的查询缓存就会被清空，所以 **MySQL8 删除了查询缓存**
>[!tip]
>这里的查询缓存指 server 层，并不是 Innodb 存储引擎的 buffer pool

## 解析 SQL

1. 词法分析：根据 SQL，构建 SQL 语法树
2. 语法分析：根据语法树，判断 SQL 语法是否正确

## 执行 SQL

1. prepare 阶段
2. optimize 阶段
3. execute 阶段

### prepare 阶段

- 检查 SQL 语句中的表或字段是否存在
- 将 `select *` 符号，扩展为表上的所有列

### optimize 阶段

主要负责将 SQL 查询语句的执行方案确定下来。如在表中有多个索引选择时，会基于查询成本考虑，决定使用哪个索引。

### execute阶段

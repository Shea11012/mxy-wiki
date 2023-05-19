---
tags: 
date created: 2022-06-10 17:35
date modified: 2023-05-18 22:13
title: mysql优化
---

## Mysql 优化

### 避免使用 select *

`select *` 不走覆盖索引，会出现回表操作，而导致查询 SQL 的性能很低

### 小表驱动大表

如：orders 表有 100000 条数据，users 表有 100 条数据

```mysql
select * from order where user_id in (select id from user where status=1)
```

```mysql
select * from order where exists (select 1 from user where order.user_id = user.id and status=1)
```

sql 语句中包含 in 关键字，则会优先执行 in 里面的子查询语句，然后再执行 in 外面的语句。如果 in 里面的数据量很少，作为条件查询速度就会变快。

sql 中包含 exists 关键字，优先执行 exists 左边的语句（即主查询语句）然后把它作为条件去和右边的语句匹配。

- in 适用于左边大表，右边小表
- exists 适用于左边小表，右边大表

### 批量操作

将通过 for 循环一次插入一条，变为一次插入多条，每批数据尽量控制在 500 以内。如果数据大于 500 则分多批次处理。

### 多用 limit

在取单条数据时，使用 `limit 1`，提高 sql 返回效率

### 避免 in 中值太多

如果 in 中的查询值过多，返回的数量太大，网络传输也是非常消耗性能的。可以考虑分批次查询

### 增量查询

```mysql
# bad
select * from user;

# good
select * from user where id >{lastid} and create_at > {lastCreateTime} limit 100;
```

通过增量查询方式，能够提升单词查询效率

### 高效分页

数据量不大时可以使用 limit 进行分页，当数据量变大时,如：

```mysql
select id,name,age from user limit 1000000,20;
```

这种查询 MySQL 会查到 1000020 条数据，然后丢弃前面的 1000000 条，只查后面的 20 条数据，这是非常浪费资源的。

优化:

```mysql
select id,name,age from user where id > 1000000 limit 20;
```

### 用连接查询代替子查询

子查询:

```mysql
select * from order where user_id in (select id from user where status=1);
```

子查询语句可以通过 in 关键字实现，一个查询语句的条件落在另一个 select 语句的查询结果中。程序先运行嵌套在最内层的语句，再运行外层的语句。
子查询语句的优点是简单，结构化，如果涉及的表数量不多。
子查询语句缺点是 mysql 执行子查询时，需要创建临时表，查询完毕后，需要再删除这些临时表，有一些额外的性能消耗。

```mysql
select o.* from order o inner join user u on o.user_id = u.id where u.status=1;
```

### join 表数量不宜过多

如果 join 太多，MySQL 在选择索引时会非常复杂，很容易选错索引。如果没有命中索引，nested loop join 就是分别从两个表读一行数据进行两两对比，复杂度是 $n^2$。
阿里巴巴开发者手册规定，join 表的数量不应该超过 3 个。

#### join 时注意事项

- `left join`：求两个表的交集外加左表剩下的数据
- `inner join`：求两个表交集的数据

inner join:

```mysql
select o.id,o.code,u.name from order o inner join user u on o.user_id = u.id where u.status=1;
```

两张表使用 inner join 关联，mysql 会自动选择两张表中的小表，去驱动大表。

left join:

```mysql
select o.id,o.code,u.name from order o left join user u on o.user_id = u.id where u.status=1;
```

使用 left join 关键字会默认使用左表去驱动它右边的表，如果左表数据很多时，就会出现性能问题。

### 控制索引数量

阿里巴巴开发者手册规定，单表的索引数量应该尽量控制在 5 个以内，单个索引中的字段数不能超过 5 个。

mysql 使用 B+ 树结构来保存索引，在 insert、update、delete 时，需要更新 B+ 树索引。如果索引过多，会消耗很多额外的性能。

### 索引优化

通过 [explain](explain使用.md) 查看 mysql 的执行计划，判断是否用到索引。
还有可能是索引失效。

## 总结

- 建立合适的索引：根据查询的具体情况和数据的分布情况选择合适的索引可以提高查询效率
- 避免全表扫描
- 使用合适的数据类型
- 优化查询语句，避免使用子查询，尽量将查询拆解成简单的
- 适当使用 join，合理使用 join 可以减少查询的数据冗余，但是如果 join 引入大量的数据组合，可能会到性能下降
- 分析查询计划，使用 explain 查看查询计划


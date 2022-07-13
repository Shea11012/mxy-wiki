---
date created: 2022-03-08 16:04
date modified: 2022-06-11 12:56
title: explain使用
---

## explain 使用

### id（select 唯一标识）

id 相同时，执行顺序从上到下

```mysql
explain select * from test1 t1 inner join test1 t2 on t1.id=t2.id
```

![](attachments/Pasted%20image%2020220610193635.png)

id 不同时，从下到上执行

```mysql
explain select * from test1 t1 where t1.id = (select id from test1 t2 where t2.id=2);
```

![](attachments/Pasted%20image%2020220610193436.png)

id 即有相同的也有不同的时候，先按照序号从大到下执行，遇到相同的再从上到下执行。

```mysql
explain  
select t1.* from test1 t1  
inner join (select max(id) mid from test1 group by id) t2  
on t1.id=t2.mid
```

![](attachments/Pasted%20image%2020220610193823.png)

### select_type(select 类型)
- simple：简单 select
- primary：最外层 select
- union：union 后的第二个或者更后面的 select
- dependent union：union 后面的第二个或更后面的 select，取决于外部查询
- union result：union 的结果
- subquery：第一个 select 子查询
- dependent subquery：第一个 select 子查询，取决于外部查询
- derived：派生表
- materialized：物化子查询
- unchangeable subquery：子查询，其结果无法缓存
- unchangeable union：union 后的第二个或更后面的 select，其结果无法缓存

### table（名称）
表示输出行所引用的表名称

### partitions（匹配的分区）
表示查询将从中匹配记录的分区

### type（连接类型）
- system：只有一条记录
- const：通过索引一次就能找到
- eq_ref：用于主键或唯一索引扫描
- ref：用于非主键和唯一索引扫描
- fulltext：全文扫描
- ref_or_null：类似 ref，还会扫描 null 数据
- index_merge：多种索引合并的扫描方式
- unique_subquery：类似 eq_ref，条件用于 in 子查询
- index_subquery：类似 unique_subquery，条件用于 in 子查询，但条件非主键或唯一索引
- range：范围扫描
- index：全索引扫描
- all：全表扫描

### possible_keys（可能的索引选择）
表示可能的索引选择

### key（实际用到的索引）
可能会出现 possible_key 列为 null，但是 key 不为 null 的情况

### key_len（实际索引长度）
决定 key_len 值得三个因素：
- 字符集
- 类型长度
- 是否为空

如：

```mysql
# 编码为utf-8
code char(30)
name varchar(30),not null
idx_code_name(code,name)
```

utf-8 编码每个字符默认为 3 字节

> 如果字段类型允许为空则加 1 字节

varchar 额外使用了 2 个字节

$key_len=30*3+1+30*3+2$

对于联合索引，可能存在 key_len 的长度不完全等于联合索引的长度，这表示只用到了联合索引的一部分

### ref（与索引比较的列）
表示索引命中的列或者常量

### rows（预计要检查的行数）
表示 mysql 认为执行查询必须检查的行数
对于 innodb 表，此数字是估值，并不总是准确

### filtered（按表条件过滤的行百分比）
表示按条件过滤的表行的估计百分比，最大值为 100，表示未过滤行。

### extra（附加信息）
包含 mysql 如何解析查询的其他信息
- impossible where：表示 where 后面的条件一直都是 false
- using filesort：表示按文件排序，一般是在指定的排序和索引排序不一致的情况下才会出现
- using index：表示是否用了覆盖索引
- using temporary：表示是否启用了临时表，一般多见于 order by 和 group by 语句
- using where：表示使用了 where 条件过滤
- using join buffer：表示是否使用连接缓冲。来自较早连接的表被部分读取到连接缓冲区，然后从缓冲区中使用它们的行来与当前表执行连接。
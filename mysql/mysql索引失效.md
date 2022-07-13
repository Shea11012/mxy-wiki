---
date created: 2022-06-11 14:29
date modified: 2022-06-11 15:13
title: mysql索引失效
---
## mysql 索引失效
表：

```mysql
create table `user` (
	`id` int not null auto_increment,
	`code` varchar(30) collate utf8mb4_bin default '',
	`age` int default `0`,
	`name` varcahr(30) collate utf8mb4_bin default '',
	`height` int default `0`,
	`address` varchar(30) collate utf8mb4_bin default '',
	primary kye (`id`),
	key `idx_code_age_name` (`code`,`age`,`name`),
	key `idx_height` (`height`)
) engine=InnoDB
```

### 不满足最左匹配原则
在使用联合索引时，没注意最左前缀原则，可能导致索引失效

#### 走索引的情况

```mysql
code='100'

code='100' and age=21

code='100' and age=21 and name='小明'

# 这种情况也会走索引，虽然断层了age，但是也会走code的索引
code='100' and name='小明'
```

#### 不走索引的情况

```mysql
age=21

name='小明'

age=21 and name='小明'
```

### select *
使用 select * 时不能用到覆盖索引

```mysql
explain select * from user where name='苏三';
```

![](attachments/Pasted%20image%2020220611145242.png)

用到覆盖索引的情况

```mysql
explain select code,name from user where name='苏三';
```

![](attachments/Pasted%20image%2020220611145407.png)

### 索引列上有计算或使用函数
都不会走索引的两个例子

```mysql
explain select * from user where id+1=2;

explain select * from user  where SUBSTR(height,1,2)=17;
```

### 字段类型不同

```mysql
explain select * from user where code=101;
```

> mysql 中如果是 int 类型字段作为查询条件时，会自动讲该字段的传参进行隐式转换，把字符串转成 int 类型

### like 左边包含%

```mysql
explain select * from user where code like '%1';
```

### 列对比
两个单独建了索引的列，用来做列对比时索引会失效

```mysql
explain select * from user where id=height;
```

### 使用 or 关键字

```mysql
# 走索引
explain select * from user where id=1 or height='175';

# 不走索引
explain select * from user where id=1 or height='175' or address='成都';
```

使用 or 时只要有一个列没有索引就会导致其他字段的索引失效，所以使用 or 时需要给 or 中的所有列都加上索引

### not in 和 not exists
#### not in

```mysql
# 不走索引，虽然height有索引
explain select * from user where height not in (173,174,175,176);

# 走索引，针对主键字段时会走索引
explain select * from user where id  not in (173,174,175,176);
```

not exists 同理

### order by 的坑
#### 使用 order by 走索引的情况：
1. 遵循最左前缀的情况下，加上 limit，不加则索引会失效

```mysql
explain select * from user order by code,age,name limit 100;
```

2. 配合 where 一起使用

```mysql
explain select * from user where code='101' order by age;
```

3. 联合索引的多个排序字段，只要它们的排序规则相同也可以走索引

```mysql
explain select * from user order by code desc,age desc limit 100;
```

#### 不走索引的情况
1. 没加 where 或 limit

```mysql
explain select * from user order by code, name;
```

2. 对不同的索引进行 order by

```mysql
explain select * from user order by code, height limit 100;
```

3. 不满足最左匹配原则

```mysql
explain select * from user order by name limit 100;
```

4. 不同的排序

```mysql
explain select * from user order by code asc,age desc limit 100;
```
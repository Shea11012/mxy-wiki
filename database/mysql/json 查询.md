
```sql
create table person (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	attr JSON DEFAULT NULL
);

INSERT person (attr) VALUES ('{"sex": "F", "age": 13, "city": "beijing"}');
INSERT person (attr) VALUES ('[1, 2, 3, 4]');
```

获取 json 内的字段
**JSON_EXTRACT**
```sql
select id, JSON_EXTRACT(attr,"$.sex") sex from person;
```

`->`
```sql
select id, attr->"$.sex" from person;
```

`->>` 会去除对应字段的双引号
```sql
select id, attr->>"$.sex" from person;
```

访问json内的数组
```sql
select id, JSON_EXTRACT(attr,'$[0]') from person;
```

json 数据修改
**JSON_SET**：替换现有KEY的值，插入不存在的key的值
**JSON_INSERT**：插入不存在的key的值，已经存在的不修改
**JSON_REPLACE**：只替换已存在的key值，不存在的不做插入
```sql
update person set attr = JSON_SET(attr,"$.city","wixi","$.height",123) where id = 1;

update person set attr = JSON_INSERT(attr,"$.city","wuxi","$.height",123) where id = 1;

update person set attr = JSON_REPLACE(attr,"$.city","wuxi","$.height",123) where id = 1;
```

**JSON_REMOVE**：删除JSON对象中的指定key
```sql
update person set attr = JSON_REMOVE(attr,"$.age") where id = 1;
```

JSON 使用索引
对相关表新增一个虚拟列，给虚拟列增加一个普通索引
```sql
alter table person add column age INT as (attr->>"$.age");
alter table person add index idx_age(age);
```
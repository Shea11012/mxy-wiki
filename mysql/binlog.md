binlog 主要用途：
- 复制：如一主多从
- 恢复

binlog的查询语句
```mysql
show binlog events [IN 'log-name'] [From pos] [Limit [offset], row_count]
```

binlog 不同格式，`--binlog-format=`：
- `statement`：生成的binlog**基于语句的日志**。此时会将一条SQL语句完整的记录到binlog上，而不管该语句影响了多少记录
- `row`：生成的binlog**基于行的日志**，会将该语句所改动的全部信息都记录上
- `mixed`：生成的binlog**基于行的日志**，一般会**基于语句的日志**，特殊情况下会自动转为**基于行的日志**


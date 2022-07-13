---
date created: 2021-12-03 20:20
date modified: 2022-03-07 18:14
title: 获取自增id
---
8.0 以前执行以下命令可以**及时**获取到表的自增 id
```sql
SELECT AUTO_INCREMENT FROM information_schema.tables WHERE TABLE_SCHEMA = "XX" AND TABLE_NAME = "XX";
```

8.0 以后 information_schema 的更新由 `information_schema_stats_expiry` 变量控制，默认值为 `86400 secs(=1 day)` ,因此对表的操作不会立即反应到 `information_schema.tables` 上，需要将 `SET @@SESSION.information_schema_stats_expiry = 0` 设置为 0，在进行查询就能实时获取表信息
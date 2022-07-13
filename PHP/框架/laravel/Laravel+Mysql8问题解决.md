---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel+Mysql8问题解决
---
```mysql
Syntax error or access violation: 1231 Variable 'sql_mode'  
   can't be set to the value of 'NO_AUTO_CREATE_USER'
```

5.7 已经废除了该模式

解决如下：

将 `config/database.php` 配置文件中的 mysql 的 `strict` 的值改为 `false` 

或者在 `databse.php` 配置文件中添加 

```php
'strict' => true,
'engine' => null,
'modes' => [
    'ONLY_FULL_GROUP_BY',
    'STRICT_TRANS_TABLES',
    'NO_ZERO_IN_DATE',
    'NO_ZERO_DATE',
    'ERROR_FOR_DIVISION_BY_ZERO',
    'NO_ENGINE_SUBSTITUTION',
]
```

参考：[MySQL8](https://dev.mysql.com/doc/relnotes/mysql/8.0/en/news-8-0-11.html) [laravel/framework#23948](https://github.com/laravel/framework/pull/23948)




---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: laravel 定时任务注意事项
---
#### laravel 设置定时任务

```sh
* * * * * php /path-project/artisan schedule:run >> /dev/null 2>&1
```

##### cron 执行用户导致 laravel.log 不可写

crontab -e 命令创建的 cron 是属于 root 用户，如果定时任务在实时主动写入日志或者遇到异常未捕捉，会创建 root 权限的日志文件，最终会导致 php-fpm 的用户无法写入

```sh
crontab -u xxx -e # 需要在创建 cron 指定用户
```

cron 内容最一行回车
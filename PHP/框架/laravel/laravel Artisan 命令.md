# laravel Artisan 命令

| 命令                                             | 说明                                             |
| ------------------------------------------------ | ------------------------------------------------ |
| php artisan key:generate                         | 生成App Key                                      |
| php artisan make:controller                      | 生成控制器                                       |
| php artisan make:model -fm                       | 生成模型 -fm 参数生成 factory 和 migration 文件  |
| php artisan make:policy                          | 生成授权策略                                     |
| php artisan make:seeder                          | 生成Seeder文件                                   |
| php artisan migrate                              | 执行迁移                                         |
| php artisan migrate:rollback                     | 回滚迁移                                         |
| php artisan migrate:refresh                      | 重置数据库                                       |
| php artisan db:seed                              | 填充数据库                                       |
| php artisan tinker                               | 进入tinker环境                                   |
| php artisan route:list                           | 查看路由列表                                     |
| php artisan make：command xxxx --command=xx:xxxx | 创建一个命令 --command 是指定 artisan 调用的命令 |

可以使用`php artisan help migrate`

#### 数据库迁移

php artisan migrate：出现`SQLSTATE[42000]: Syntax error or access violation: 1071 Specified key was too long; max key length is 1000 bytes `这个错误有如下几个修改方案：

1. MySQL升级到5.5.3以上的版本
2. 手动配置迁移命令migrate生成的默认字符串长度，在`AppServiceProvider`中调用`Schmea::defaultStringLength`方法来实现配置

```php
public function boot()
{
    Schema::defaultStringLength(191)
}
```

1. 配置`config/database.php`

```php
'charset' => 'utf8',
'collation' => 'utf8_unicode_ci',
```

**创建用户对象**

```php
php artisan tinker
APP\Models\User::create(['name'=>'shea','email'=>'1872314654@qq.com','password'=>bcrybt('password')])
```

退出`ctrl+c`
##### 用户认证

`php artisan make:auth`

##### 邮箱认证

```php
use Illuminate\Contracts\Auth\MustVerifyEmail as MustVerifyEmailContract; // 作为接口类继承此类确保 User 遵守规范
use Illuminate\Auth\MustVerifyEmail as MustVerifyEmailTrait;
Class User extends Controller implements MustVerifyEmailContract
{
    use MustVerifyEmailTrait;
    // 该 trait 包含三个方法
    // hasVerifiedEmail 检测用户 email 是否认证
    // MarkEmailAsVerified 将用户标示为已认证
    // sendEmailVerificationNotification 发送 Email 认证消息通知，触发邮件的发送
}
```

###### 新建话题 观察者模式

观察者类的方法名反映 Eloquent 想监听的事件，每个方法接受 model 作为唯一参数

```php
php artisan make:observer UserObserver --model=User
```

在 `AppServiceProvider` 注册观察者

```php
class AppServiceProvider extends ServiceProvider
{
    public function boot()
    {
        User::observer(UserObserver::class);
    }
}
```

###### 使用队列

队列系统可以异步执行消耗时间的任务，比如请求一个 API 并等待返回的结果，这样可以有效降低请求响应时间

失败任务，laravel  内置了一个 failed_jobs 表的迁移文件 

```php
php artisan queue:failed-table // 生成迁移文件
```

生成一个任务类

```php
php artisan make:jobs TranslateSlug
```

注意：在模型监控器中分发任务，任务中需要避免使用  Eloquent 模型接口调用，如：`create，update，save` 等，否则会陷入死循环，模型监控器分发任务-》任务触发模型监控器-》模型监控器再次分发任务-》任务再次触发模型监控器，在这种情况下，使用 DB 类直接对数据库进行操作

在模型监控器中分发任务

```php
class TopicObserver
{
    public function saved(Topic $topic)
    {
        ...
            if(!$topic->slug){
                // 推送任务到队列
                dispatch(new TranslateSlug($topic));
            }
    }
}
```

可以使用 `php artisan queue:listen` 启动监听队列系统

使用 Horizon 方便查看和管理 redis 队列任务执行情况

```php
composer require laravel/horizon
php artisan vendor:publish --provider="Laravel\Horizon\HorizonServiceProvider"
    // 发布 horizon 相关配置文件
    // config/horizon.php 配置文件
    // public/vendor/horizon css js 页面资源文件
```

`php artisan horizon` 启动，再打开 `host/horizon/dashboard` 查看

注意：生产环境下使用队列需要注意两个问题：

1. 使用 Supervisor 进程工具进行管理
2. 每一次部署代码时，需 `artisan horizon:terminate` 然后再 `artisan horizon` 重新加载代码

##### laravel 消息通知

```php
// laravel 有一个默认的通知表
php artisan notifications:table
// 
```

生成通知类

```php
php artisan make:notification TopicReplied
```

```php
// 通知类 TopicReplied 内此方法定义通知在哪个频道发送
public function via($notifiable)
{
  	return ['mail','database']  
}

//针对不同的频道实现不同的方法
public function tomail($notifiable){}
public function toDatabase($notifiable){}
```

#### 多角色用户权限管理



#### 活跃用户

活跃用户算法：

系统每小时计算一次，统计最近 7 天所有用户发的帖子数和评论数，用户每发一个帖子得 4 分，每发一个回复得 1 分，计算出所有得分再倒序，将排名前八的用户显示在活跃用户列表里。

使用laravel 命令调度器对laravel 命令调度



#### 最近活跃用户

思路：使用 redis 记录用户访问时间，定期将 redis 数据同步到数据库中

1. 使用中间件过滤用户所有请求，记录用户访问时间到 redis 按日期区分的哈希表
2. 同步，使用计划任务每天运行一次此命令，将昨日哈希表的数据同步到数据库中，并删除
3. 读取，优先读取当日 Redis 里的数据，无数据则使用数据库的值


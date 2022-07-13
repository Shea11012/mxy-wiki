---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel 模型事件
---
## Eloquent 模型

模型会触发事件（Event），可以对模型的生命周期多个时间点进行监控：

- creating、created
- updating、updated
- saving、saved
- deleting、deleted
- restoring、restored

事件可以每当有特定的模型类在对数据库进行操作时，触发特定的事件。

观察者类的方法名对应 Eloquent 想监听的事件。每种方法接收 `model` 作为唯一的参数。

Laravel 没有为观察者提供默认目录，可以创建任意目录存放观察者类。

定义好观察者类后如 `TopicObserver` 再在 `provider\AppServiceProvider` 的 `boot` 方法内注册，如：`\App\Models\Topic::observe(\App\Observers\TopicObserver::class)`

## 模型事件的几种用法：

1. 路由中定义：

```php
Event::listen('eloquent.updated:App\Post',function(){
    dump('测试事件');
});
Route::post('post/{id}','PostController@update');
```

2. 生成事件和监听器

EventServiceProvider 定义对应关系

```php
protected $listen = [
    'App\Events\PostEvent' => [
        'App\Listeners\PostListener',
    ],
];
```

`php artisan event:generate // 生成文件`

> event 中注入要操作的类
>
> listen 中 handle 方法注入对应事件类

```php
public function handle(PostEvent $event)
{
    dump('测试事件');
}
```

最后在 post 模型中添加 events 属性

```php
protected $events = [
    'updated' => PostListener::class
];
```

3. 利用框架 boot 方法

在相关 model 类中定义

```php
public static function boot()
{
    parent::boot();
    
    static::updated(function($model){
        dump('测试事件');
    });
}

```

4. 定义 Trait

针对多个模型都需进行的事件，封装一个 Trait

```php
trait LogRecord
{
    // 必须以 boot 开头
    public static function bootLogRecord()
    {
        foreach(static::getModelEvents() as $event) {
            static::$event(function($model){
                $model->setRemind();
            });
        }
    }
    
    public static function getModelEvents()
    {
        if(isset(static::$recordEvents)){
            return static::$recordEvents;
        }
        return ['updated'];
    }
    
    public function setRemind()
    {
        dump('记录逻辑操作');
    }
}
```




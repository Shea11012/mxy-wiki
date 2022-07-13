---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: 将 laravel 的各个组件单独安装 
---
## 阅读 《laravel 框架关键技术解读》后笔记，此书参照框架版本：5.1



### 添加路由组件

`composer require illuminate/routing`

`composer require illuminate/events`

路由组件依赖 `events` 组件

```php
// index.php

// 调用自动加载文件
require __DIR__.'/../vendor/autoload.php';
// 实例化服务容器，注册事件、路由服务提供者
$app = new Illuminate\Container\Container;
with(new Illuminate\Events\EventServiceProvider($app))->register();
with(new Illuminate\Routing\RoutingService\Provider($app))->register();
// 加载路由
require __DIR__.'/../app/Http/routes.php';
// 实例化请求并分发处理请求
$request = Illuminate\Http\Request::createFromGlobals();
$response = $app['router']->dispatch($request);
//返回请求响应
$response->send();
```

控制器模块，添加路由组件时已经添加了基本控制器类 `Illuminate\Routing\Controller` 可以直接拿来使用



### 添加模型组件

`composer require illuminate/database`

```php
// database.php
return [
    'driver'    =>  'mysql',
    'host'      =>  'mysql',	// 这里使用的是 docker 所以这样写
    'database'  =>  'lara',
    'username'  =>  'zhuzi',
    'password'  =>  '123456',
    'charset'   =>  'utf8mb4',
    'collation' =>  'utf8mb4_general_ci',
    'prefix'    =>  ''
];

// index.php
use Illuminate\Database\Capsule\Manager;
require __DIR__.'/../vendor/autoload.php';
.
.
.
// 启动 Eloquent 配置
$manager = new Manager;
$manager->addConnection(require '../config/database.php');
$manager->bootEloquent();
.
.
.
```

创建 model 文件时 继承此类 `use Illuminate\Database\Eloquent\Model` 即可



### 添加视图组件

`composer require illuminate/view`

先创建视图模板文件路径：`project/resources/views`

编译视图文件存储路径：`project/storage/framework/views`

```php
//index.php
use Illuminate\Support\Fluent;
.
.
.
$app = new Illuminate\Container\Container;
Illuminate\Container\Container::setInstance($app);
.
.
.
// 视图配置和服务注册
$app->instance('config',new Fluent);
$app['config']['view.compiled'] = __DIR__."/../storage/framework/views/";
$app['config']['view.paths'] = [__DIR__."/../resources/views/"];
with(new Illuminate\View\ViewServiceProvider($app))->register();
with(new Illuminate\Filesystem\FilesystemServiceProvider($app))->register();
```

使用服务容器的 `setInstance` 静态方法将服务容器实例添加为静态属性，这样可以在任何位置获取服务容器的实例。

将服务名称 `config` 和 `illuminate\Support\Fluent` 类进行实例绑定，该类的实例主要用于存储视图模块的配置信息


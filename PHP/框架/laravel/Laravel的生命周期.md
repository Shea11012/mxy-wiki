---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel的生命周期
---
# Laravel 的生命周期

了解 Laravel 的生命周期前，先了解下 PHP 的生命周期。

### PHP 的生命周期

- PHP 两种运行模式 WEB 模式、CLI 模式
  1. 当在终端敲入 php 这个命令时，使用的是 CLI 模式
  2. 当使用 Nginx 或者其他 web 服务器时，使用的是 WEB 模式
- 生命周期
当请求一个 PHP 文件时，PHP 为了完成这次请求，会发生 5 个阶段的生命周期切换：

1. 模块初始化（MINIT），即调用 `php.ini` 中指明的扩展的初始化函数进行初始化工作，如 `mysql` 扩展
2. 请求初始化（RINIT），即初始化为执行本次脚本所需要的变量名称和变量值内容的符号表，如 `$_SESSION` 变量
3. 执行该 PHP 脚本
4. 请求处理完成（Request Shutdown），按顺序调用各个模块的 `RSHUTDOWN` 方法，对每个变量调用 `unset` 函数，如 `unset($_SESSION)` 变量
5. 关闭模块（Module Shutdown），PHP 调用每个扩展的 `MSHUTDOWN` 方法，这是各个模块最后一次释放内存的机会，这意味着没有下一个请求了

WEB 模式和 CLI 模式（命令行）很相似，区别是：

1. CLI 模式会在每次脚本执行经历完整的 5 个周期，因为你的脚本执行完成不会有下一个请求
2. WEB 模式为了应对并发，可能采用多线程，因此生命周期 1 和 5 有可能只执行一次，下次请求来时重复 2-4 生命周期，这样就节省了系统模块初始化所带来的开销

![2119758-ede5488b10aef9b0](E:\MD 笔记\PHP\框架\laravel\2119758-ede5488b10aef9b0.png)

PHP 是一种脚本语言，所有的变量只会在这一次请求中生效，下次请求之时已被重置。

### Laravel 的生命周期

![2119758-3790698ae8356013](E:\MD 笔记\PHP\框架\laravel\2119758-3790698ae8356013.png)

请求过程 `public/index.php`

```php
1. require __DIR__.'/../bootstrap/autoload.php';

2. $app = require_once __DIR__.'/..bootstrap/app.php';
   $kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

3. $response = $kernel->handle(
	   $requset = Illuminate\Http\Request::capture()
   );
   $response->send();

4. $kernel->terminate($request,$response);
```

详细解释如下：

1. 文件载入 composer 生成的自动加载设置，包括所有的 `composer require` 依赖
2. 生成容器 Container，Application 实例，并向容器注册核心组件（HttpKernel，ConsoleKernel，ExceptionHandler）
3. 处理请求，生成并发响应（99% 代码都运行在这个 handle 方法里）
4. 请求结束，进行回调

![](E:\MD 笔记\PHP\框架\laravel\2119758-3d17c34fd9c561b8.png)

更加详细的解释是：

1. 注册加载 composer 自动生成的 `class loader`：加载初始化第三方依赖
2. 生成容器 `Container`：向容器注册核心组件，是从 `bootstrap/app.php` 脚本获取 Laravel 应用实例
3. 请求被发送到 HTTP 内核或 Console 内环境，这取决于进入应用的请求类型

> **取决于是通过浏览器请求还是通过控制台请求，这里主要是通过浏览器请求**

HTTP 内核继承自 `Illuminate\Foundation\Http\Kernel` 类，该类定义一个 boostrappers 数组，这个数组中的类在请求被执行前运行，这些 boostrappers 配置了错误处理、日志、检测应用环境以及其它在请求被处理前需要执行的任务。

```php
protected $bootstrappers = [
        //注册系统环境配置 （.env）
        'Illuminate\Foundation\Bootstrap\DetectEnvironment',
        //注册系统配置（config）
        'Illuminate\Foundation\Bootstrap\LoadConfiguration',
        //注册日志配置
        'Illuminate\Foundation\Bootstrap\ConfigureLogging',
        //注册异常处理
        'Illuminate\Foundation\Bootstrap\HandleExceptions',
        //注册服务容器的门面，Facade 是个提供从容器访问对象的类。
        'Illuminate\Foundation\Bootstrap\RegisterFacades',
        //注册服务提供者
        'Illuminate\Foundation\Bootstrap\RegisterProviders',
        //注册服务提供者 `boot`
        'Illuminate\Foundation\Bootstrap\BootProviders',
    ];
```

**bootstrappers 数组中的类一定要注意顺序**

4. 将请求传递给路由

Laravel 基础服务启动之后把请求传递给路由之后，路由器会将请求分发到路由或控制器，同时运行所有路由指定的中间件。传递给路由是通过 `Pipline`（管道）来传递，`Pipline` 有一堵墙，在传递给路由之前所有请求都要经过，这堵墙定义在 `app\Http\Kernel.php` 中的 `$middleware` 数组中，默认只有一个 `CheckForMaintenanceMode` 中间件，用来检测网站是否暂时关闭，这是一个全局中间件，所有请求都要经过，也可以自定义自己的全局中间件。

然后遍历所有注册的路由，找到最先符合的第一个路由，经过它的路由中间件，进入到控制器或者闭包函数，执行具体的逻辑代码

![](E:\MD 笔记\PHP\框架\laravel\2119758-ae9730b03f6ce607.png)

[本文转载自 PHP 中文社区 - Laravel 的生命周期](https://phperzh.com/articles/3118)


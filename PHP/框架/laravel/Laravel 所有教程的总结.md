### 服务容器

几乎所有的服务容器绑定操作都是在**服务提供器**中注册

在服务提供器中，可以通过 `$this->app` 属性访问容器

### 服务提供者

Laravel 的 `config/app.php` 文件中有一个 `providers` 数组。数组中的内容是应用程序要加载的所有服务提供器类。有一些提供器并不会在每次请求时都被加载，只有当它们提供的服务实际需要时才会加载，这种被称为 『延迟』提供器。

所有服务提供器都会继承  `Illuminate\Support\ServiceProvider` 类。大多数服务提供器都包含 `register` 和 `boot` 方法。在 `register` 方法中，只需要绑定类到服务容器，而不需要尝试在 `register` 方法中注册任何事件监听器、路由或任何其他功能。

任何服务提供器方法中，可以通过 `$app` 属性来访问服务容器。

### 事件系统





### 授权策略

Laravel 使用[授权策略](http://d.laravel-china.org/docs/5.5/authorization#policies)来对用户的操作权限进行验证，在用户未经授权进行操作时返回 403 禁止访问的异常。

流程：

1.  `php artisan make:policy XxxPolicy` ：所有的授权策略文件都被放置在 `app/Polices` 下
2. 在授权策略文件中添加方法
3. 在 `AuthServiceProvider` 类中对授权策略进行设置，`AuthServiceProvider` 包含一个 `policies` 属性，该属性用于将各种模型对应到管理它们的授权策略上：

```php
protected $policies = [
    'App\Model' => 'App\Policies\ModelPolicy',
    \App\Models\User::class => \App\Policies\UserPolicy::class,
]
```

授权策略定义完成，便可以在用户控制器中使用 `authorize` 方法来验证用户授权策略

详细参考 [Laravel5,.5-Web开发入门：权限系统](https://laravel-china.org/courses/laravel-essential-training-5.5/599/permissions-system)

#### 获取当前登录用户的信息

```php
use Illuminate\Support\Facades\Auth;
use Illuninate\Http\Request;
// 获取当前认证的用户信息
Auth::user(); 或者 $request->user();
```

详见 [用户认证](https://laravel-china.org/docs/laravel/5.6/authentication/1308)



`load()` 和 `with()`  都是延迟预加载，区别在于 `load()` 是在已经查询出来的模型上调用，而 `with()` 是在 ORM 查询构造器上调用

```php
$users = User::with('commnets')->get();
// SQL 语句如下
select * from `users`
select * from `commnets` where `comments`.`user_id` in (1,2,3,4,5)
```

```php
$users = User::all();
$users = $users->load('comments');
// SQL 语句如下
select * from `users`
select * from `comments` where `comments`.`user_id` in (1,2,3,4,5)
```

参考 [with  和 load 区别](https://stackoverflow.com/questions/26005994/laravel-with-method-versus-load-method)

对于 ORM 关联模型 `$order->items()` 和 `$order->item` 的区别

`$order->items()` 是获取关联关系，这时还没有发生 SQL 查询，通常是准备做进一步的查询

`$order->items` 则是获取关联模型，SQL 已经执行完毕，已经从数据库取到了所有关联的数据

在 `Request` 模型里面 `$this->route('order')` 这样可以获取当前路由的订单对象



`where(function($query){})` 嵌套是用来生成 SQL 里的括号，保证不会因为 `or` 关键字导致查出来的结果与期望不符。
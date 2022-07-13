---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: laravel 5.5-web开发实战入门总结
---
# laravel 5.5-web 开发实战入门总结

在控制器中指定渲染某个视图，使用 view 方法，view 接受两个参数，第一个是视图的路径名称，第二个是与视图绑定的数据，第二个参数为可选

#### laravel 模板

命名方式 xxx.blade.php

```html
<html>
    <head>
        <title>@yield('title','Sample')</title>
    </head>
    <body>
        @yield('content')
    </body>
</html>
```

`@yield('content')` 表示该占位区域用于显示 `content` 区块的内容，`content` 区块内容由继承自 default 视图的子视图定义

`@yield('title','Sample')`,第一个参数是该区块的变量名称，第二个参数是默认值，表示当指定变量的值为空值时，使用 `Sample` 来作为默认值，当 `@section` 传递了第二个参数时，便不再需要通过 `@stop` 表示来告诉 Laravel 填充区块会在具体哪个位置结束了。

`@extends('layouts.default')` 使用 `@extends` 并通过传参来继承父视图 `layouts/default.blade.php` 模板

使用 `@section` 和 `@stop` 代码来填充父视图的 `content` 区块，所有包含在 `@section` 和 `@stop` 中的代码都将被插入到父视图的 `content` 区块

```html
@select('content')
<h1>主页</h1>
@stop
```

`@include('layouts._header')` 是 Blade 提供的视图引用方法，可通过传参一个具体的文件路径名称来引用视图

extends：继承后可以覆盖；include：引入后不能覆盖

`npm install` 安装 laravel 里的前端扩展 `npm run dev` 和 `npm run watch-poll` 可以及时对修改的 css 进行编译和保存

#### Laravel Mix

`mix.js('resources/assets/js/app.js','public/js')`：第二个参数是自定义生成 js 文件的输出目录

#### Laravel 路由

```php+HTML
Route::get('/help','StaticPagesController@help')->name('help')	//name方法指定路由名称
<li><a href="{{ route('help') }}">帮助</a></li>	//route方法调用help
```

隐式绑定

```php
Route::get('api/users/{user}',function(App\User $user){
    return $user->email;
})
```

上面的例子中，由于 `$user` 变量被类型提示为 Eloquent 模型 `APP\User` ，变量名称又与 URI 中的 `{user}` 匹配，因此， Laravel 会自动注入与请求 URI 中传入的 ID 匹配的用户模型实例。如果在数据库中找不到对应的模型实例，将会自动生成 404 异常。

资源路由列表信息：

| HTPP 请求 | URL                | 动作                    | 作用                   |
| -------- | ------------------ | ----------------------- | ---------------------- |
| GET      | /users             | UsersController@index   | 显示所有用户列表的月面 |
| GET      | /users/{user}      | UsersController@show    | 显示用户个人信息的页面 |
| GET      | /users/create      | UsersController@create  | 创建用户的页面         |
| POST     | /users             | UsersController@store   | 创建用户               |
| GET      | /users/{user}/edit | UsersController@edit    | 编辑用户个人资料的页面 |
| PATCH    | /users/{user}      | UsersController@update  | 更新用户               |
| DELETE   | /users/{user}      | UsersController@destroy | 删除用户               |

#### Laravel 模型关联

一对一

`user` 模型关联一个 `phone` 模型

```php
class User extends Model
{
    public function phone()
    {
       return $this->hasOne(Phone::class,'user_id','phone_id'); 
    }
}
//会默认Phone模型有一个user_id外键，如果要覆盖可以传递第二个参数
//默认会与User模型的主键值或id相匹配，如果要覆盖可以传递第三个参数
```

反向关联

通过 ``phone` 模型查到拥有该电话的 `user` 模型

```php
class Phone extends Model
{
    public function user()
    {
        return $this->belongsTo(User::class);
    }
}
//会尝试匹配phone模型上的user_id至User模型上的id，如果phone的外键不是user_id，可以传递自定义键名作为第二个参数
//如果user模型没有以id作为主键，或希望用不同的字段，可以传递第三个参数
```

### 模型文件（Eloquent）

`fillable` 属性可以过滤用户提交的字段，是为了防止用户恶意提交 `is_admin` 这样的字段，只有包含在该属性中的字段才能够被正常更新。

`guarded` 属性包含的是不想被批量赋值的属性的数，注意 `fillable` 和 `guarded` 二选一

`hidden` 属性里的值，在通过数组或 json 显示时进行隐藏

模型文件默认对应的数据表，是文件名的复数形式「蛇形命名」作为表名，例：`App\Models\User.php` 表名默认为 `users` 

全局作用于确保给定模型的每个查询都受到一定的约束

1. 首先定义一个 `Illuminate\Database\Eloquent\Scope` 接口类，实现一个 `apply` 方法

```php
namespace App\Scopes;

use Illuminate\Database\Eloquent\Scope;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Builder;

class AgeScope implements Scope
{
    public function apply(Builder $builder,Model $model)
    {
        return $builder->where('age','>',200);
    }
}
```

2. 应用全局作用域，需要重写给定的模型 `boot` 方法并使用 `addGlobalScope` 方法

```php
namespace App\Models;

use App\Scopes\AgeScope;
use Illuminate\Database\Eloquent\Model;

class User extends Model
{
    protected static function boot()
    {
        parent::boot();
        static::addGlobalScope(new AgeScope);
    }
}
```

Eloquent 模型的生命周期：`reterieved` 、`creating` 、`created` 、 `updating` 、 `updated` 、`saving` 、`saved` 、 `deleting` 、`deleted` 、 `restoring` 、`restored` 

### 邮件激活账户「本地测试」

1. 配置 `.env` 文件 `MAIL_DRIVER=log`
2. 创建一个激活邮件的路由例：`Route::get('signup/confirm/{token}','UsersController@confirmEmail')->name('confirm_email');`
3. 去 `view` 目录配置邮件模板
4. 定义一个发送邮件的方法 `sendEmailConfirmationTo`，该方法将邮件发送给指定用户

```php
//如果是在生产环境，#号注释的可以去除
use Mail;
class UsersController extend Controller
{
    public function sendEmailConfirmTo($user)
    {
        $view = 'email.email';
        $data = compact('user');
       # $from = 'test@qq.com';
       # $name = 'test';
        $to = 'fdhasfh@qq.com';
        $subject = '感谢注册！';
        //本地测试
        Mail::send($view,$data,function($message)use($from,$to,$subject){
            $message->from($from,$name)->to($to)->subject($subject);
        })   
            //线上环境
            Mail::send($view,$data,function($message) use ($subject){
                $message->to($to)->subject($subject);
            })
    }   
}
```

Mail 的 send 方法接受三个参数

- 第一个参数是包含邮件消息的视图模板
- 第二个参数是传递给该视图的数据数组
- 第三个参数是一个接受邮件消息的实例的闭包回调，可以在该回调中定义消息的发送者、接收者、邮件主题

5. 在 `__construct` 方法里允许未登录用户访问 `confirmEmail` 方法
6. 定义 `confirmEmail($token)` 方法，激活成功后 `Auth::login($user)` 自动登录，将 `activated=true;activation_token=null`

### 权限系统

Laravel 中间件（Middleware）提供了过滤机制来过滤 HTTP 应用请求，默认内置了一些中间件，如：身份验证、CSRF 保护等。所有的中间件文件都被放在项目的 `app\Http\Middleware` 文件中。

`php artisan make:middleware filename`

在 `__construct` 方法中调用 `middleware` 方法，该方法接收两个参数，第一个为中间件 名称，第二个为要进行过滤的动作

`except`：设定指定动作不使用 Auth 中间件进行过滤，意为——除了此处指定动作外，所有其他动作都必须登录用户才能访问。（一般推荐使用这个）

`only`：将只过滤指定动作

Auth 中间件过滤指定动作，默认会被重定向到 `/login`

限制登录用户的操作，使用 Laravel 的授权策略（Policy）来对用户的操作权限进行验证，对于未经授权操作时将返回 403

`php artisan make:policy filenamePolicy`

生成的授权文件都在 `app\Polices` 文件夹下

定义好授权策略文件还需要在 `AuthServiceProvider` 类中对授权策略进行设置，里面包含了一个 `policies` 属性

例如：生成一个 `UserPolicy` 再将其模型对应到管理它们的授权策略上

```php
namespace App\Policies;
use Illuminate\Auth\Access\HandlesAuthorization;
use App\Models\User;
class UserPolicy
{
    use HandlesAuthorization;
    public function update(User $currentUser,User $user)
    {
        return $currentUser->id === $user->id;
    }
}
```

`update` 方法接受两个参数，第一个参数默认为当前登录用户实例，第二个参数则为要进行授权的用户实例

授权策略需要注意以下两点：

1. 不需要检查 `$currentUser` 是不是 Null，未登录用户，框架会自动为其所有权限返回 `false`
2. 调用时，默认情况下，我们不需要传递当前登录用户至该方法内，因为框架会自动加载当前登录用户

```php
namespace App\Providers;

class AuthServiceProvider extends ServiceProvider
{
    protected $policies = [
        'App\Model'=>'App\Policies\ModelPolicy',
        \App\Models\User::class => \App\Policies\UserPolicy::class,
    ];
}
```

授权策略定义完成之后，便可以通过用户控制器使用 `authorize` 方法来验证用户授权策略。默认 `App\Http\Controllers\Controller` 类包含了 Laravel 的 `AuthorizesRequests`trait，此 trait 提供了 `authorize` 方法，当无权限运行时会抛出 `HttpException`。`authorize` 方法接受两个参数，第一个是授权策略名称，第二个是进行授权验证的数据。

`$this->authorize('update',$user);`

> update 是指授权类的 update 授权方法，$user 对应传参 update 授权方法的第二个参数，默认情况下，不需要传递第一个参数，当前登录用户至该方法内，因为框架会自动加载当前登录用户

### 友好转向

`redirect()` 实例提供一个 `intended` 方法，该方法可将页面重定向到上一次请求尝试访问的页面上，并接受一个默认跳转地址参数，当上一次请求记录为空时，跳转到默认地址上。

`redirect()->intended(route('users.show',[Auth::user()]));`



### 重置密码

1. 添加以下路由

| HTTP 请求 | URL                     | 动作                                              | 作用                       |
| -------- | ----------------------- | ------------------------------------------------- | -------------------------- |
| GET      | /password/reset         | Auth\ForgotPasswordController@showLinkRequestForm | 显示重置密码的邮箱发送页面 |
| POST     | /password/email         | Auth\ForgotPasswordController@sendResetLinkEmail  | 邮箱发送重设链接           |
| GET      | /password/reset/{token} | Auth\ResetPasswordController@showResetForm        | 密码更新页面               |
| POST     | /password/reset         | Auth\ResetPasswordController@reset                | 执行密码更新操作           |

上面的几个方法已经被 Laravel 自动定义好了，可以直接配置重置密码和更新密码的模板

2. `php artisan make:notification ResetPassword` ，生成消息通知文件，在 `app\Notifications\ResetPassword.php` 
3. 定制消息通知文件 `ResetPassword` 

```php
class ResetPassword extends Notification
{
    public $token;
    
    public function __construct($token)
    {
        $this->token = $token;
    }
    
    public function via($notifiable)
    {
        return ['mail'];
    }
    
    public function toMail($notifiable)
    {
        return (new MailMessage)
            ->subject('重置密码')
            ->line('这是一封密码重置邮件，如果是本人操作，请点击以下按钮继续：')
            ->action('重置密码',url(route('password.reset',$this->token,false)))
            ->line('如果您并没有执行此操作，您可以选择忽略此邮件。')
    }
}
```

4. User 模型里调用

```php
namespace App\Models;

use App\Notifications\ResetPassword;
class User extends Authenticatable
{
    public function sendPasswordResetNotification($token)
    {
        $this->notify(new ResetPassword($token));
    }
}
```

5. 发布密码重置 Email 视图，`php artisan vendor:publish --tag=laravel-notifications`

### 生产环境发送邮件

1. 开启任意邮箱的 POP3 和 SMTP 服务，获取授权码
2. 配置 `.env`

```env
MAIL_DRIVER = smtp
MAIL_HOST = smtp.qq.com
MAIL_PORT = 25
MAIL_USERNAME = XXXXXXX@qq.com
MAIL_PASSWORD = XXXXXX
MAIL_ENCRYPTION = tls
MAIL_FROM_ADDRESS = XXXXXXX@qq.com  //与MAIL_USERNAME一致
MAIL_FROM_NAME = xxxx
```


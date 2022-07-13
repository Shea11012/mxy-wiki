---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel web开发实战复习
---
### web 开发实战

##### laravel 解决浏览器缓存问题

laravel-mix 是 laravel 提供的一个前端解析工具

```javascript
cosnt mix = require('laravel-mix')

mix.js('resources/js/app.js','public/js')
	.sass('resources/sass/app.css','public/css').version()

# 使用了 version() 函数，在模板引用 css 时，必须使用 mix() 方法
```

```php
Route::resource('users','UsersController');
# 表示以下路由
Route::get('/users','UsersController@index')->name('users.index'); # 显示所有用户的列表页面
Route::get('/users/create','UsersController@create')->name('users.create'); # 创建用户页面
Route::get('/users/{user}','UsersController@show')
->name('users.show'); # 显示用户个人信息
Route::post('/users','UsersController@store')->name('users.store'); # 创建用户
Route::get('/users/{user}/edit','UsersController@edit')->name('users.edit'); # 编辑用户个人资料页面
Route::patch('/users/{user}','UsersController@update')->name('users.update'); # 更新用户
Route::delete('/users/{user}','UsersController@destory')->name('users.destory'); # 删除用户

```

laravel 提供全局辅助函数 `old` 帮助在 blade 模板显示旧输入数据，这样当信息填写错误，页面进行重定向访问时，输入框将自动填写上最后一次输入过的数据

`{{ csrf_field() }}` 可以防止 post 提交受到 CSRF (跨站请求伪造)

> laravel 默认会将所有的验证信息进行闪存，当检测到错误存在时，laravel 会自动将这些错误消息绑定到视图上，因此可以在所有视图上使用 `errors` 变量来显示错误信息，Tips：使用 `errors` 变量时 先使用 `count($errors)` 检查其值是否为空

```php
redirect()->route('users.show',[$user]);
// 此处这样的写法，是 laravel 的 route() 方法会自动获取 model 的主键，等同于以下
redirect()->route('users.show',[$user->id]);
```

由于 HTTP 协议是无状态的，laravel 提供了一种用于临时保存用户数据的方法 session ，并附带了支持多种会话后端驱动，可通过同意的 API 使用

使用 `session()` 方法访问会话实例，当想存入一条缓存数据，让它只在下一次请求内有效是，可以使用 `flash` 方法

```php
session()->flash('success','欢迎，您将在这里开启一段新的旅程');
# 取出方法
session()->get('success');
```

 laravel 提供了 `Auth` 的 `attempt` 方法可以很方便的完成用户的身份认证，详细的用户认证查看 [laravel - 用户认证章节](https://learnku.com/docs/laravel/5.7/authentication/2269)

```php
return redirect()->back()->withInput();
# 使用 withInput() 后模板里的 old() 方法将能获取到上一次用户提交的内容
```

`method_field('DELETE')` 由于浏览器不支持发送 DELETE 请求，而 RESTful 架构中会使用 DELETE 请求来删除一个资源，所以使用了隐藏域来伪造 DELETE 请求

```html
<input type="hidden" name="_method" value="DELETE">
```

laravel 默认配置中，如果用户登录后没有使用记住我功能，则登录状态只会被记住两个小时

[laravel 中间件](https://learnku.com/docs/laravel/5.7/middleware/2254) 为过滤进入应用程序的 HTTP 请求提供了一种方便的机制

laravel 提供了一个 Auth 中间件，在控制器的 `__construct` 方法中调用

```php
public function __construct()
{
    $this->middleware('auth',[
        'except'	=>	['show','create','store'] // 定义除了这个3个方法其余方法都需要验证
    ])
}
```

laravel 提供的 Auth 中间件在过滤指定动作时，会将未通过身份验证的用户重定向到 `/login` 登录页面

**用户只能编辑自己的资料，需要限制已登录用户的操作，当 ID 为 1 的用户尝试更新 ID 为 2 的用户信息时，应该返回一个 403 禁止访问的异常。laravel 的 [授权策略](https://learnku.com/docs/laravel/5.7/authorization/2271#policies) 可以对用户的操作权限进行验证**

```php
php artisan make:policy UserPolicy // 生成至 app/Policies 文件夹下
    
# app/Policies/UserPolicy.php
class UserPolicy
{
    .
    .
    .
    public function update(User $currentUser,User $user)
    {
        return $currentUser->id == $user->id;
    }
    
    // 此处的 $currentUser 是当前登录用户实例(等于 Auth::user())，$user 是要进行授权的用户实例（也就是隐式路由自动绑定的实例）
   // 并不需要检查 $currentUser 是不是 NULL ，未登录用户，框架会自动给为其所有权限返回 false
}

# AuthServiceProvider 类中 Policies 属性上对授权策略进行设置
class AuthServiceProvider extends ServiceProvider
{
    protected $policies = [
        'App\Model'	=>	'App\Policies\ModelPolicy',
        \App\Models\User::class => \App\Policies\UserPolicy::class
    ];
    // 这样便可以在控制器中的 authorize 方法来验证用户授权策略
}

class UsersController extends Controller
{
    public function create()
    {
        // $this->authorize('update',$user);
        // create 方法不需要验证是否是同一个用户，可以这样使用：
        $this->authorize('update',User::class);
    }
}
```
##### 友好的转向

当一个未登录用户尝试访问自己的编辑页面时，会自动跳转到登录页面，如果用户进行登录，则应该将其重定向至他尝试访问的页面，redirect() 提供了一个 intended 方法，该方法可将页面重定向到上一次请求尝试访问的页面，并接受一个默认跳转地址参数，当上一次请求为空时跳转到默认地址上。

```php
$fallback = route('users.show',Auth::user());
return redirect()->intended($fallback);
```

Laravel Auth 中间件提供 `guest` 选项，用于指定一些只允许未登录用户访问的动作，对 `guest` 属性进行设置，只让未登录用户访问登录页面和注册页面

```php
$this->middleware('guest',[
    'only'	=>	['create']
]);
```

###### 账户激活

给用户表新增两个字段用于保存用户的激活令牌和激活状态，激活令牌用于验证用户身份，激活状态用于判断用户是否已激活。


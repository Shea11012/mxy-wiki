---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel5.5-API接口开发
---
### GitHub 的 Restful HTTP API

RESTful 是一种软件设计风格，由 Roy Fielding 提出。

1. HTTPS

HTTPS 为接口的安全提供了保障，生产环境推荐使用 HTTPS（非 HTTPS 的 API 调用，不要重定向到 HTTPS，而要直接返回调用错误以禁止不安全的调用）

2. 域名

应当尽可能的将 API 与其主域名区分开：`https://api.xxx.com` 或 `https://www.xxx.com/api`

3. 版本控制

为了保证旧用户可以正常使用和开发顺利进行，需要控制 API 版本

两种做法：

- `https://api.xxx.com/v1` ：将版本号直接加入到 URL 中

- ```
  https://api.xxx.com/
  	Accept: application/prs.xxx.v1+json	//GitHub 推荐使用这种方式
  ```
4. 用 URL 定位资源，资源应当是一个名词，大部分是名词的复数，尽量不要再 URL 中出现动词

5. HTTP 动词描述操作

幂等性 ：指一次和多次请求某一个资源应该具有同样的副作用，也就是说一次访问与多次访问对这个资源带来的变化是相同的。

| 动词   | 描述                               | 是否幂等 |
| ------ | ---------------------------------- | -------- |
| GET    | 获取资源，单个或多个               | 是       |
| POST   | 创建资源                           | 否       |
| PUT    | 更新资源，客户端提供完整的资源数据 | 是       |
| PATCH  | 更新资源，客户端提供部分的资源数据 | 否       |
| DELETE | 删除资源                           | 是       |

6. 资源过滤
7. 正确使用 HTTP 状态码
8. 数据响应格式
9. 调用频率限制

### 安装 DingoApi

安装 DingAPI ，是一个 Lumen 和 Laravel 都可用的 RestFul 工具包， [文档](https://github.com/liyu001989/dingo-api-wiki-zh/blob/master/API-Blueprint-Documentation.md) 

安装时需要在 composer.json 文件中添加以下内容

```json
"config" : {
    "preferred-install" : "dist",
    "sort-packages" : true,
    "optimize-autoloader" : true
},
"minimum-stability" : "dev",	// 设定的最低稳定性的版本为 dev 也就是可以依赖开发版本的扩展
"prefer-stable" : true			// composer 优先使用更稳定的包版本
```

`php artisan vendor:publish` ，选择 `Dingo\Api\Provider\LaravelServiceProvider`

### 构建用户注册接口

`hash_equals($string1,$string2)` ：可以防止时序攻击的字符串比较

`$string1 == $string2` ：这样子比较，字符串从第一位开始逐一进行比较，发现不同就立即返回 false，那么通过计算返回的速度就知道了大概是哪一位开始不同的，这样就实现了按位破解密码的场景。而 `hash_equals` 比较两个字符串，无论字符是否相等，函数的时间消耗是恒定的，这样可以有效防止时序攻击。

### 微信登录

[OAuth2.0](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html) 第三方登录的常用协议

基本流程：

1. 客户端（app/浏览器）将用户导向第三方认证服务器
2. 用户在第三方认证服务器，选择是否给予客户端授权
3. 用户同意授权后，认证服务器将用户导向客户端实现指定的 `重定向URI` ，同时附上一个授权码
4. 客户端将授权码发送至服务器，服务器通过授权码以及 `APP_SECRET` 向第三方服务器申请 `access_token`
5. 服务器通过 `access_token` ，向第三方服务器申请用户数据，完成登录流程

使用扩展包 [socialiteproviders](https://socialiteproviders.github.io/) 

区别微信用户时，使用 `UnionID` 来区分微信用户

原因：

- `OpenID` 是针对 『 微信应用』的用户唯一值，同一个『微信开发者账号』下的不同应用中，使用同一个『微信用户』登录，此值会不一样。
- `UnionID` 是针对『微信开发者账号』的用户唯一值，同一个『微信开发者账号』下的不同应用中，使用同一个『微信用户』登录，此值是一致的。

### 修改话题

`$this->authorize(' 动作 ',' 相关模型 ')；` ：继承于 `App\Http\Controller\Controller` 这个基类控制器，和 `can` 方法类似，接收想授权的动作和相关的模型作为参数，如果动作没被授权 `authorize` 会抛出一个 `Illuminate\Auth\Access\AuthorizationException` 的异常，Laravel 默认的异常处理器会将这个异常转化成带有 `403` 状态码的 HTTP 响应

### 话题列表

DinogAPI 的 [Include 机制](https://laravel-china.org/courses/laravel-advance-training-5.5/923/list-of-posts)

当需要返回额外的资源数据时，设置 `protected $availableInclude = ['user','category']` ，数组中每一个参数都对应一个具体的方法，命名规则 `include + user` 驼峰命名，引入额外的资源方式 `?include=user,xxx`  

在 Transformer 中，可以使用：

- $this->item() ：返回单个资源
- $this->collection() ： 返回集合资源

[laravel-query-logger](https://github.com/overtrue/laravel-query-logger) 是 SQL 查询日子组件 ：`composer require overtrue/laravel-query-logger --dev`

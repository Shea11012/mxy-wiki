---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Laravel商城
---
自定义的辅助函数放在 `bootstrap/helpers.php` ，使用 composer 的 autoload 功能自动引入：

在 composer.json 文件，找到 autoload 段

```php
"autoload":{
    "classmap":[
        'database/seeds',
        'database/factories'
    ],
    "psr-4":{
        "App\\":"app/"
    },
    "files":[
        'bootstrap/helpers.php'
    ]
}
```

然后执行 `composer dumpautoload`

`url()` 全局函数，生成网站对应 URL，可以使用 Route 定义的别名，前提是要有别名

用户模块，使用 Laravel 自带的用户认证脚手架 `php artisan make:auth`

调整数据库结构，添加一个字段使用 `php artisan make:migration 表名_add_ 字段名 --table=表`

### Model 模型

 `$fillable` 属性内的字段表示可以进行批量赋值，`$guarded` 属性内的字段表示不可被批量赋值

`$casts` 属性支持插入或者更新数据时对类型转换 `$date` 属性可以自行定义哪些日期类型字段会被自动转换

参见 [Eloquent 修改器](https://laravel-china.org/docs/laravel/5.6/eloquent-mutators/1406#attribute-casting)

### 验证邮箱

当希望用户在验证邮箱之后才能正常使用功能，未验证邮箱时，访问其他页面都会被重定向到一个提示验证邮箱的页面。这种就可以使用中间件

**创建中间件** `php artisan make:middlerware CheckIfEmailVerified` ，生成文件位于 `app/Http/Middleware` 下

定义好中间件后，需要再去 `app/Http/Kernel.php` 根据用途在不同属性注册该中间件 参考 [中间件](https://laravel-china.org/docs/laravel/5.5/middleware/1294#2fb6d2) 

使用 Laravel 内置的通知模块 （Notification）来实现验证邮件的发送

**创建通知** `php artisan make:notification EmailVerificationNotification` 生成文件位于 `app/Notifications` 目录下 参见 [验证邮箱下](https://laravel-china.org/courses/laravel-shop/1584/verification-mailbox-below) [消息通知](https://laravel-china.org/docs/laravel/5.5/notifications/1322) 

激活链接 、用户主动请求发送激活邮件参见 [验证邮箱下](https://laravel-china.org/courses/laravel-shop/1584/verification-mailbox-below#f0fcb0) 

注册时触发激活邮件，可以通过 Laravel 事件系统来完成这个功能，注册完触发一个 `Illuminate\Auth\Events\Registered` 事件，可以创建一个这个事件的监听器 （listener）发送邮件。

> 监听器是 Laravel 事件系统的重要组成部分，当一个事件被触发，对应的监听器就会被执行，可以方便的解耦代码，还可以把监听器配置成异步执行，适合一些不需要获得返回值并且耗时较长的任务，比如发送邮件。

**创建事件监听器** `php artisan make:listener RegisterdListener`，生成文件位于 `app\Listeners` 下

**处理异常** 创建异常类 `php artisan make:exception InvalidRequestException` ，生成文件位于 `app/Exceptions` 下，Laravel 5.5 之后支持在异常类定义 `render()` 方法，该异常被触发时系统会调用 `render()` 方法来输出。

有一些异常不需要被记录时，可以在 `app/Exceptions/Handler.php` 的 `$dontReport` 数组中添加

`user_addresses` 字段：

| 字段名称      | 描述             | 类型          | 索引 |
| ------------- | ---------------- | ------------- | ---- |
| id            | 自增长 ID         | unsigned int  | 主键 |
| user_id       | 该地址所属的用户 | unsigned int  | 外键 |
| province      | 省               | varchar       | 无   |
| city          | 市               | varchar       | 无   |
| district      | 区               | varchar       | 无   |
| address       | 具体地址         | unsigned int  | 无   |
| zip	 |  邮编 | unsigned int | 无|
| contact_name  | 联系人姓名       | varchar       | 无   |
| contact_phone | 联系人电话       | varchar       | 无   |
| last_used_at  | 最后一次使用时间 | datetime null | 无   |

创建模型 `php artisan make:model Models/UserAddress -fm` -fm 表示创建 factory 和 migration 

使用工厂文件自动填充数据时，如果是中文填充，在 `config/app.php` 内添加一个 `'faker_locale' => 'zh_CN'` ，避免数据库字段不一致报错

**新建收货地址：**使用了 vue 实现 [3 级联动](https://laravel-china.org/courses/laravel-shop/1585/new-receiving-address#da874f) ，使用 Request 类完成数据校验 `php artisan make:request Request` ，使用扩展包中文提示扩展包 `composer require overtrue/laravel-lang` 参考 [表单验证](https://laravel-china.org/docs/laravel/5.5/validation/1302) 

**后台管理系统** 使用 laravel-admin 扩展包实现快速搭建后台，[中文文档](http://laravel-admin.org/docs/#/zh/) 

`Admin::content` ：页面内容和布局

`Admin::grid` ：模型表格

`Admin::form` ：模型表单

`Admin::show` ：模型详情

#### 商品数据结构设计

**商品 SKU 概念：**  SKU = Stock Keeping Unit（库存量单位），也可以称为『单品』。对一种商品而言，当其品牌、型号、配置、等级、花色、包装容量、单位、生产日期、保质期、用途、价格、产地等属性中任一属性与其它商品存在不同时，可称为一个单品。

两个数据表：

- `products` 表，产品信息表，对应数据模型 Product；
- `product_skus` 表，产品 SKU 表，对应数据模型 ProductSku

`products` 表：

| 字段名称     | 描述                 | 类型                                          | 索引 |
| ------------ | -------------------- | --------------------------------------------- | ---- |
| id           | 自增长 ID             | unsigned int                                  | 主键 |
| title        | 商品名称             | varchar                                       | 无   |
| description  | 商品详情             | text                                          | 无   |
| image        | 商品封面图片文件路径 | varchar                                       | 无   |
| on_sale      | 商品是否正在售卖     | tiny int，default 1                           | 无   |
| rating       | 商品平均评分         | float，default 5                              | 无   |
| sold_count   | 销量                 | unsigned int，default 0                       | 无   |
| review_count | 评价数量             | unsigned int，default 0                       | 无   |
| price        | SKU 最低价格         | decimal（数值型，不存在精度损失，28 个有效位） | 无   |

> 商品本身没有固定价格，在商品表放置 price 字段目的是方便用户搜索、排序

`product_skus` 表：

| 字段名称    | 描述        | 类型         | 索引 |
| ----------- | ----------- | ------------ | ---- |
| id          | 自增长 ID    | unsigned int | 主键 |
| title       | SKU 名称    | varchar      | 无   |
| description | SKU 描述    | varchar      | 无   |
| price       | SKU 价格    | decimal      | 无   |
| stock       | 库存        | unsigned int | 无   |
| product_id  | 所属商品 id | unsigned int | 外键 |

任何与钱相关的有小数点的字段一律使用 `decimal(a,b)` 需要两个参数 第一个是数值总精度，另一个参数则是小数位

**商品排序和筛选：** 

where 语句参数分组的解释：

```php
$builder->where(function ($query) use ($like) {
                $query->where('title', 'like', $like)
                    ->orWhere('description', 'like', $like)
                    ->orWhereHas('skus', function ($query) use ($like) {
                        $query->where('title', 'like', $like)
                            ->orWhere('description', 'like', $like);
                    });
            });
```

先使用 `$builder->where()` 传入一个匿名函数，然后才在这个匿名函数里再去添加 `like` 搜索，这样的目的是在查询条件的两边加上 `()`  ，`orwhereHas` 可以增加自定义条件至关联约束中

`{{ $products->appends($filters)->render() }}` ：appends 接受一个 key-value 形式的数组作为参数，在生成分页链接的时候会把这个数组格式化成查询参数

**收藏商品是用户和商品的多对多关联，不需要创建新的模型，只需要增加一个中间表**



### 购物车

将用户数据存入数据库，建表 `cart_items` 表：

| 字段名称       | 描述        | 类型         | 索引 |
| -------------- | ----------- | ------------ | ---- |
| id             | 自增 ID     | unsigned int | 主键 |
| user_id        | 所属用户 id | unsigned int | 外键 |
| product_sku_id | 商品 SKU ID | unsigned int | 外键 |
| amount         | 商品数量    | unsigned int | 无   |

表单验证使用闭包函数，三个参数分别是参数名、参数值和错误回调



### 订单模块

将购物车的商品提交成订单，`orders` 表：

| 字段名称       | 描述                | 类型               | 索引 |
| -------------- | ------------------- | ------------------ | ---- |
| id             | 自增 id             | unsigned int       | 主键 |
| no             | 订单流水号          | varchar            | 唯一 |
| user_id        | 下单的用户 id       | unsigned int       | 外键 |
| address        | JSON 格式的收货地址 | text               | 无   |
| total_amount   | 订单总额            | decimal            | 无   |
| remark         | 订单备注            | text               | 无   |
| paid_at        | 支付时间            | datetime，null     | 无   |
| payment_method | 支付方式            | varchar，null      | 无   |
| payment_no     | 支付平台订单号      | varchar，null      | 无   |
| refund_status  | 退款状态            | varchar            | 无   |
| refund_no      | 退款单号            | varchar，null      | 唯一 |
| closed         | 订单是否关闭        | tinyint，default 0 | 无   |
| reviewed       | 订单是否评价        | tinyint，default 0 | 无   |
| ship_status    | 物流状态            | varchar            | 无   |
| ship_data      | 物流数据            | text，null         | 无   |
| extra          | 其他额外数据        | text，null         | 无   |

`order_items` 表：

| 字段名称       | 状态            | 类型             | 索引 |
| -------------- | --------------- | ---------------- | ---- |
| ID             | 自增 id         | unsigned int     | 主键 |
| order_id       | 所属订单 id     | unsigned int     | 外键 |
| product_id     | 对应商品 id     | unsigned int     | 外键 |
| product_sku_id | 对应商品 SKU ID | unsigned int     | 外键 |
| amount         | 数量            | unsigned int     | 无   |
| price          | 单价            | decimal          | 无   |
| rating         | 用户打分        | unsigned int     | 无   |
| review         | 用户评价        | text             | 无   |
| reviewed_at    | 评价时间        | timestamp ，null | 无   |

#### 关闭未支付订单

避免恶意用户下单占用商品库存，当创建订单之后一定时间没有支付，将关闭订单并退回减去库存。

使用 Laravel 提供的延迟任务（Delayed Job）功能解决。当触发了一个延迟任务时，Laravel 会用当前时间加上任务的延迟时间计算出任务应该被执行的时间戳，然后将这个时间戳和任务信息序列化之后存入队列，Laravel 队列处理器会不断查询并执行队列中满足预计执行时间等于或早于当前时间的任务。

**创建任务** `php artisan make:job CloseOrder` ，任务类保存在 `app/jobs` 目录下，写好任务类后，在 controller 中使用 `dispatch` 辅助函数来分发

使用 reids 队列需要将 `.env` 内的驱动设置 `sync` 改为 `redis`  ，还需要引入 `composer require predis/predis`  ，启动队列处理器 `php artisan queue:work` ，队列详细用法参考 [队列](https://laravel-china.org/docs/laravel/5.5/queues/1324#supervisor-configuration)



### 封装业务代码

为了防止 controller 变得臃肿，防止以后 App 端重复写，采用 service 模式来封装代码。

封装功能需要注意 `$request` 不可以出现在控制器和中间件以外的地方，根据职责单一原则，获取数据这个任务应该由控制器来完成。



### 支付模块

使用 `yansongda/pay` 扩展 ，这个扩展封装了支付宝和微信的接口。`composer require yansongda/pay`，扩展文档 [yansongda/pay](https://yansongda.gitbooks.io/pay/)

#### 支付宝支付沙箱环境测试

支付宝的支付回调分为 前端回调 和 服务器回调

前端回调 ：指用户支付成功之后支付宝会让用户浏览器跳转回项目页面并带上支付成功的参数

服务器回调：指支付成功之后支付宝服务器会用订单相关数据作为参数请求项目的接口

因此判断支付是否成功要以服务器回调为准

#### 微信支付

微信支付没有沙箱测试，参考链接 [微信支付](https://laravel-china.org/courses/laravel-shop/1550/wechat-payment)

Laravel 验证服务器回调需要将验证路由加入 CSRF 校验白名单在 `app/Http/Middleware/VerifyCsrfToken.php` 的 `$except` 内加入

#### 完善支付后的逻辑

创建一个支付成功的事件，事件本身不需要包含逻辑，只需要包含相关信息，再创建两个监听器监听事件，一个监听销量，一个监听邮件发送，再在相应的 controller 类中触发这个事件使用 `event` 全局函数触发。详细参考 [完善支付后逻辑](https://laravel-china.org/courses/laravel-shop/1602/perfect-the-logic-after-payment)



### 完善订单模块

#### 管理后台

自定义后台的 `laravel-admin` 里的 `show` 方法对应 `disableView`

#### 评价商品

#### 同意退款（支付宝，微信）



### 优惠卷模块

`CouponCode` 表字段

| 字段       | 描述                                 | 类型                     | 索引 |
| ---------- | ------------------------------------ | ------------------------ | ---- |
| id         | 自增长 ID                             | unsigned int             | 主键 |
| name       | 优惠卷的标题                         | varchar                  | 无   |
| code       | 优惠码，用户下单时输入               | varchar                  | 唯一 |
| type       | 优惠卷类型，支持固定金额和百分比折扣 | varchar                  | 无   |
| value      | 折扣值，根据不同类型含义不同         | decimal                  | 无   |
| total      | 全站可兑换数量                       | unsigned int             | 无   |
| used       | 当前已兑换的数量                     | unsigned int ，default 0 | 无   |
| min_amount | 使用该优惠卷的最低订单金额           | decimal                  | 无   |
| not_before | 在这个时间之前不可用                 | datetime，null           | 无   |
| not_after  | 在这个时间之后不可用                 | datetime，null           | 无   |
| enabled    | 优惠卷是否生效                       | tinyint                  | 无   |



### 配置后台权限
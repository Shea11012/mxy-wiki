---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: 编写高效PHP
---
对象是通过引用传递的

```php
class A
{
    public $name;
}
$a = new A();
$a->name = 'world';
$b = $a;	// 这样是在引用（reference）
$b->name = 'hello';
echo $a->name,$b->name;	// 输出 hello hello
```

改为以下：

```php
class A
{
    public $name;
}
$a = new A();
$a->name = 'world';
$b = clone $a;	// 复制对象的一个副本
$b->name = 'hello';
echo $a->name,$b->name;	// 输出 world hello
```

**当复制一个对象是，存储在其属性中的任何对象都将是引用而不是副本，因此在处理复杂的面向对象的应用程序时必须非常小心**

**PHP 内置一个魔术方法 __clone ，当复制对象时会做些什么，也可以拒绝复制**

对象是引用传递，可以在应用程序内建立一个流畅的接口（fluent interface）

```php
class A 
{
    public function test(){return $this;}
    public function run(){return $this;}
    public function say(){return $this;}
}
$a = new A;
$a->test()->run()->say();	// 不需要按顺序调用
```



#### 魔术方法的介绍

`__construct()`：实例化一个对象

`__destruct()`：销毁一个对象

`__get()`：读取一个不存在的属性

`__set`：写入一个不存在的属性

`__clone`：复制一个对象

`__call`：当一个不存在方法时，调用

`__callStatic`：当静态方法不存在时，调用

`__toString`：当使用 echo 输出对象时，调用

`__sleep`：在序列化（serialize）时调用

`__wakeup`：在反序列化（unserialize）时调用



### 第二章 数据库

存储数据方式：

- 文本文件：对于很少更新的少量数据，是理想的选择（如配置文件，或在应用程序中记录事件或错误）
- 会话数据：对于只为下一次请求或访问持续期间所需的数据，可以在用户的会话中存储信息，避免记录过多的数据或者添加功能以清理不需要的数据
- 关系数据库：MySQL Oracle 等
- NoSQL 数据库：MongoDB CouchDB 等

#### 第三章 API

使用 API 的原因：

- 使数据用于其他系统或模块
- 以异步的方式向网站提供数据
- 构成一个面向服务架构的基础

**SOA（Service-Oriented Architecture 面向服务架构）**：它是基于一个服务处的系统，提供系统需要的所有功能，但这个服务提供的是应用层，并未链接到表现层。这样多种系统就可以使用这个相同的模块化、可重复使用的功能
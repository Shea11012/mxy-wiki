---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: PHP错误和异常处理
---
## PHP 错误和异常处理

### PHP 错误等级

```
# 系统级用户代码的一些错误类型，可由 try catch 捕获
E_PARSE 解析时错误，语法解析错误，如：少个分号，多个逗号的致命错误
E_ERROR 运行时错误，比如调用了未定义的函数或方法，致命错误

# 系统级用户代码的一些错误类型，可用 set_error_handler 捕获处理
E_WARNING	运行时警告，调用了未定义的变量
E_NOTICE	运行时提醒
E_DEPRECATED	运行时已废弃的函数或方法

# 用户级自定义错误，可由 trigger_error 触发，可由 set_error_handler 捕获
E_USER_ERROR	用户自定义错误，致命错误，未处理也会导致程序退出
E_USER_WARNING	
E_USER_NOTICE
E_USER_DEPRECATED


# zend engine 内部的一些错误，也能通过 try catch 捕获
E_CORE_ERROR
E_CORE_WARNING
E_COMPLIE_ERROR
E_COMPLIE_WARNING

# 编码标准化警告
E_STRICT	部分 try catch 部分 set_error_handler
E_RECOVERABLE_ERROR


```



#### PHP 错误相关函数

- `ini_set('display_errors',0);`	//关闭错误输出
- `error_reporting(E_ALL&~E_NOTICE);` //设置错误报告级别
- `ini_set('error_reporting',0);` // 设置错误报告级别

#### PHP 异常与错误的抛出

- 异常抛出：`throw new Exception('some error message');`
- 错误抛出：`trigger_error()`
- `trigger_error()` 触发的错误不会被 `try-catch()` 异常捕获语句捕获

#### PHP 错误处理

- `set_error_handler()`

只能处理 `Deprecated、Notice、Waring` 这三种错误，脚本继续执行后发生错误的后一行

- `register_shutdown_function()`

脚本结束前的最后一个回调函数，无论是 die，错误，异常 还是脚本正常结束都会调用

#### PHP 异常处理

- `set_exception_handler()`

设置默认的异常处理程序，有 try/catch 捕获的话这个函数就不会执行

set_exception_handler() ，不仅可以接受函数名，还可以接受类的方法（公开的静态及非静态方法）

#### PHP7 异常处理变化

PHP7 之后，出现一个异常和错误通用接口 Throwable，Exception 和 Error 类都实现了该接口，使得 Error 或 Error 的派生类的错误对象（大部分 Fatel Error）也可以像 Exception 一样被捕获

```php
try{
    // code
} catch(Exception $e){
    // 捕获异常
}catch(Error $error){
    // 捕获错误
}
//或者
try{
    // code
}catch(Throwable $e) {
    // 可以捕获异常和错误
}
```

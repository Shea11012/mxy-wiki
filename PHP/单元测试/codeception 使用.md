---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Codeception
---
# Codeception

[codeception 文档](https://codeception.com/docs)

## Acceptance Tests

配置文件：`acceptance.suite.yml`

#### PhpBrowser

codeception 使用 PhpBrowser 测试浏览器

codeception 浏览器断言：

```php
 see('内容','元素节点'); // 可以断言节点内是否含有该内容，不指定节点则是当前页
 dontsee('内容','节点'); // 反向断言
```

每一个 see 方法都有一个对应的 canSee 、dontSee、cantSee 方法。

条件断言默认关闭需要在配置文件启用

### webDriver

codeception 使用 webDriver 可以测试 JavaScript，但是需要使用 wait 等待 js 脚本执行完毕。

配置 wait 时间可以在脚本内使用：`$I-wait(3) // 等待 3 秒 `

也可以使用配置文件的方式，在 yml 文件内配置：`wait:5`

选择器也受到了限制，只能使用三种：css 选择器，xpath 选择器，`['css'=>'button']`

retry 方法默认是被禁用的，可以在 yml 文件开启：

```yml
step_decorators:
 - \Codeception\Step\Retry
```

开启之后可以引入一个 trait：

```php
use \Codeception\Lib\Actor\Shared\Retry;
```

存储登录状态，模块必须实现 `Codeception\Lib\Interfaces\SessionSnapshot`

## Functional Tests

配置文件：`functional.suite.yml`

## Api Tests

## Unit Tests

配置文件：`unit.suite.yml`

specify 可以隔离代码使得测试代码不受污染

```php
use \Codeception\Specify;

    /** @specify */ // 通过使用 annotation 方式可以恒定一个变量
    private $user;

    public function testValidation()
    {
        $this->user = User::create();

        $this->specify("username is required", function() { // 再使用 specify 方法可以使得代码被隔离
            $this->user->username = null;
            $this->assertFalse($this->user->validate(['username']));
        });

        $this->specify("username is too long", function() {
            $this->user->username = 'toolooooongnaaaaaaameeee';
            $this->assertFalse($this->user->validate(['username']));
        });

        $this->specify("username is ok", function() {
            $this->user->username = 'davert';
            $this->assertTrue($this->user->validate(['username']));
        });
    }
```


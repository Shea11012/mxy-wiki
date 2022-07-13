---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: phpunit使用
---
# phpunit 使用

### `@depends` 标注来表达依赖关系

> **`@depends` 传递对象是一个引用而不是一个副本，如果需要传递对象的副本而非引用，则应当用 `@depends clone`**

```php
use PHPUnit\Framework\TestCase;

class MultipleDependenciesTest extends TestCase
{
    public function testProducerFirst()
    {
        $this->assertTrue(true);
        return 'first';
    }

    public function testProducerSecond()
    {
        $this->assertTrue(true);
        return 'second';
    }

    /**
     * @depends testProducerFirst
     * @depends testProducerSecond
     */
    public function testConsumer()
    {
        $this->assertEquals(
            ['first', 'second'],
            func_get_args()			//会获得@depends传来的值
        );
    }
}
```

### `@dataProvider` 数据供给器

> 数据供给器方法必须声明为 `public` ，其返回值要么是一个二维数组，要么是一个实现了 `Iterator` 接口的对象

```php
use PHPUnit\Framework\TestCase;

class DataTest extends TestCase
{
    /**
     * @dataProvider additionProvider
     */
    public function testAdd($a, $b, $expected)
    {
        $this->assertEquals($expected, $a + $b);
    }

    public function additionProvider()
    {
        return [
            'adding zeros'  => [0, 0, 0],
            'zero plus one' => [0, 1, 1],
            'one plus zero' => [1, 0, 1],
            'one plus one'  => [1, 1, 3]
        ];
    }
    
    // 使用带有命名数据集的数据供给器，这样输出的信息会更加详细些
    public function additionProvider()
    {
        return [
            'adding zero' => [0,0,0],
            'zero plus zero' => [0,1,1],
            'one plus one' => [1,0,1]
        ]
    }
}
```

### 对异常进行测试

`@expectedException,@expectedExceptionCode,@expectedExceptionMessage,@expectedExceptionMessageRegExp`

### 未完成的测试与跳过的测试

`@require` 标注表达测试用例的一些常见条件

| 类型      | 可能的值                              | 范例                         |
| --------- | ------------------------------------- | ---------------------------- |
| PHP       | 任何 PHP 版本标识符                   | @requires PHP 7.1            |
| PHPUnit   | 任何 PHPUnit 版本标识符               | @requires PHPUnit 4.0        |
| OS        | 用来对 PHP_OS 进行匹配的正则表达式    | @requires OS Linux           |
| function  | 任何对 function_exists 而言有效的参数 | @requires function imap_open |
| extension | 任何扩展模块名，可以附带有版本标识符  | @requires extension PDO      |

### 数据库测试


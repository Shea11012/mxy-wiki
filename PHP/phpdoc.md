## @author

当前文件的作者

```php
语法：@author name <email>
@author mxy mxyinvoker@gmail.com
```

## @deprecated

被此标记的函数或者类方法在下个版本会被废弃，告知使用方不再推荐使用此方法，可以配合 @see 使用

```php
语法：@deprecated version description
@see apiv2()
@deprecated 1.0.0 此方法已被废弃，请使用v2
```

## @inheritdoc

文档继承，会继承父类的文档注释，在继承之后可以对指定字段进行重写

```php
class A {
	/**
	 *@param string $string
	 *@param array $arr
	 */
	public function show($string,$arr){}
}

class B extends A {
	/**
	 *@inheritdoc
	 *@param int $arr 此处对继承自父类的方法，修改了参数 $arr 的接收类型
	 */
	 public function show($string,$arr)
}
```

## @internal

被此标签标记的内部类、方法，作用范围只能限于当前文件，外部文件不可调用

通常使用场景在单元测试中，万一在引入测试类时，在 ide 的帮助下可以进行提示

```php
/**
 * @internal
 */
 class InternalTest {
	// 引入此类就会，IDE 就会提示该类只能在内部调用 
 }
```

## @link

引导到指定的外部链接

```php
语法：@link url description
@link http://github.com/Shea11012/dnmp
```

## @method

此标签告诉调用方该类有哪些方法可用

通常使用在 facade

```php
语法：@method [modifier] [return type] [name]([type parameter]....) description

/**
 *@method int incr($number,$step = 1) 递增数字
 *@method array explode(string $delimiter,string $string) 分割字符串
 *@method static void incrStatic(int $number,int $step = 1) 静态递增数字
 */
class TestFacde {
	public function __call($name,$arguments){}
}
```

## @param

定义函数、类方法的参数类型

### 变量列表

|       变量类型       |         说明         |
| :------------------: | :------------------: |
|        string        |        字符串        |
|     integer、int     |         整型         |
|    boolean、bool     |        布尔型        |
|    float、double     |        浮点型        |
|        Object        |       对象实例       |
|    specified type    |        指定类        |
|        mixed         |       任意类型       |
| array、specifiedtype | 数组、指定类型的数组 |
|       resource       |     文件资源类型     |
|         void         |       无返回值       |
|         null         |          -           |
|       callable       |   可执行的回调函数   |
|       function       |  不一定能执行的方法  |
|     self、$this      |       当前实例       |

```php
class Param {
	/**
	 *@param string[] collections 字符串类型的集合
	 *@param string $pices 碎片
	 *@return string[]
	 */
  public function join(string[] $collection,string $piece): string[] {
    $collection[] = $piece;
    return $collection;
  }
}
```

## @property

定义当前类中具有的属性方法

通常在类中包含魔术方法 `__get、__set` 时，通过此标签设置

```php
/**
 *@property int $intVar 数字
 *@property string $stringVar 字符串
 *@property mixed $any 任意类型返回值
 */
class Property {
  public function __get($name){}
  public function __set($name,$value){}
}
```

## @return

用于函数、类方法返回值信息

```php
方法：@return type description
class TestReturn {
  /**
  *	@return string 主机名称
  */
  public function getHost() {}
  
  /**
   * @return string|array
   */
  public function getConfig(){}
}
```

## @throws

抛出一个异常，告诉调用方需要处理异常

```php
// 语法：@throws type description
class TagThrows {
  /**
   *@return float|int
   *@throws \Exception
   */
  public function exec(){}
}
```

## @var

定义一个数据的类型

```php
// 语法：@var type name description
class TagVar {
  /**
   * @var array  // 在类成员变量中定义不需要指定变量名称
   */
  public $tags;
}

/**
 * @var array $tags	// 直接给具体变量定义，需要指定变量名称
 */
strlen($tags);
```


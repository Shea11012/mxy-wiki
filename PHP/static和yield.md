# static和yield

### static

1. 静态化，占用内存
2. 后期绑定，`static::` 不再被解析为定义当前方法所在的类，而是在实际运行时计算得到的

```php
class A
{
    public static function who()
    {
        echo __CLASS__;
    }
    
    public static function test()
    {
        self::who();
        static::who();	//后期静态绑定
    }
}

class B extends A
{
    public static function who()
    {
        echo __CLASS__;
    }
}

B::test();
```

当调用的时候`self::who()`输出A，`static::who()`输出B

### yield

```php
function generateData($max)
{
    for($i=0;$i<$max;$i++){
        yield $i;
    }
}
$data = generateData(1000000)	//得到的是一个迭代器 
    foreach($data as $v){
        if($v>10) break;
        echo $v;
    }
```

如果数据来源非常大，使用yield，数据来源小，可以使用for或foreach

### Trait

1. 优先级：当前类的方法会覆盖 Trait 中的方法，而 Trait 中的方法会覆盖基类的方法
2. 多个 Trait 组合：通过逗号分隔，通过 use 关键字列出多个 Trait
3. 冲突的解决：如果两个 Trait 都插入了一个同名的方法，若没有明确解决冲突将会产生一个致命错误。为了解决多个 Trait 在同一个类中的命名冲突，需要使用 insteadof 操作符来明确指定使用冲突方法中的哪一个。同时可以通过 as  操作符将其中一个冲突的方法以另一个名称来引入。
4. 修改方法的访问控制：使用 as 语法可以用来调整方法的访问控制
5. Trait 的抽象方法：在 Trait 中可以使用抽象成员，使得类中必须实现这个抽象方法
6. Trait 的静态成员：在 Trait 中可以用静态方法和静态变量
7. Trait 的属性定义：在 Trait 中同样可以定义属性
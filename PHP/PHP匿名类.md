---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: PHP匿名类
---
## PHP 匿名类

匿名类可以创建一次性的简单对象

匿名类被嵌套进普通类后，不能访问这个外部类的 private、protected 方法或属性。但如果想访问 protected 方法或属性，可以继承这个外部类，想访问 private 方法或属性，可以通过构造器

```php
class Animal{
    private $num = 1;
    protected $age = 2;
    protected function bark(){
        return 10;
    }
    public function drive(){
        return new class($this->num) extends Animal{
            protected $id;
            public function __construct($sum){
                $this->id = $sum;
            }
            public function eat(){
                return $this->id + $this->age + $this->bark();
            }
        };
    }
}
echo (new Animal)->drive()->eat();	// 13
```



匿名类的闭包实现

```php
$test = [];
for($i=0;$i<6;$i++){
    $test[] = new class($i){
        public $age;
        public function __construct($num){
            $this->age = $num;
        }
        public function getValue(){
            return $this->age;
        }
    };
}
echo $test[0]->getValue();	// 0
var_dump($test[2]);			// object 对象
```


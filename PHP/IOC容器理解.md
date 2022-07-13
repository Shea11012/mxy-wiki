---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: IOC容器理解
---
依赖实例

```php
interface Visit
{
    public function go();
}

class Leg implements Visit
{
    public function go()
    {
        echo 'Leg';
    }
}

class Car implements Visit
{
    public function go()
    {
        echo 'car';
    }
}

class Traveller
{
    protected $trafficTool;
    public function __construct()
    {
        $this->trafficTool = new Leg(); // 此处产生依赖
    }
    
    public function visitTibet()
    {
        $this->trafficTool->go();
    }
}

$traveller = new Traveller;
$traveller->visitTibet();
```

上述例子中，在程序中依赖可以理解为一个对象实现某个功能需要其他对象相关功能的支持。当用 `new` 关键字在一个组件内部实例化一个对象时就解决了一个依赖，但同时也引入了另一个严重的问题——耦合。当项目变大，需求大改时，会产生很严重的副作用。

我们不应该在类的内部初始化依赖的类，应该转由外部负责，在系统运行期间，将这种依赖关系通过动态注入的方式实现，这就是 `IOC` 模式的设计思想



上述例子使用简单工厂模式进行优化

```php
class TrafficToolFactory
{
    public function createTrafficTool($name)
    {
        switch($name) {
            case 'Leg':
                return new Leg();
                break;
            case 'Car':
                return new Car();
                break;
            default:
                exit('set TrafficTool error!!!');
                break;
        }
    }
}

class Traveller
{
    protected $trafficTool;
    public function __construct($trafficTool)
    {
        $factory = new TrafficToolFactory();
        $this->trafficTool = $factory->createTrafficTool($trafficTool);
    }
    
    public function visitTibet()
    {
        $this->trafficTool->go();
    }
}

$traveller = new Traveller('Train');
$traveller->visitTibet();
```

使用工厂模式解决了 “旅游者” 和 “交通工具” 之间的依赖关系，但是却变成了 “旅游者” 和 “交通工具工厂” 之间的依赖。当需求增加时，我们需要修改工厂模式，如果依赖增多，工厂将会十分庞大，不易于维护。

### 依赖注入模式

`IOC（Inversion of Control` 模式 【依赖注入模式 `Depe-ndency Injection` 模式】

控制反转是将组件间的依赖关系从程序内部提到外部容器来管理，而依赖注入是指组件的依赖通过外部以参数的或其他形式注入

简单的依赖容器注入实现

```php
// 容器类
class Container
{
    // 用于装提供实例的回调函数，真正的容器还会装实例等其他内容
    protected $bindings = [];

    // 绑定接口和生成相应实例的回调函数
    public function bind($abstract,$concrete = null,$shared = false)
    {
        if (!$concrete instanceof Closure) {
            // 如果提供的参数不是回调函数，则产生默认的回调函数
            $concrete = $this->getClosure($abstract,$concrete);
        }
        $this->bindings[$abstract] = compact('concrete','shared');
    }

    // 默认生成实例的回调函数
    protected function getClosure($abstract,$concrete)
    {
        // 生成实例的回调函数 ，$c 一般为 IOC 容器对象，在调用回调函数生成实例时提供
        return function($c) use ($abstract,$concrete) 
        {
            $method = ($abstract == $concrete) ? 'build' : 'make';
            // 调用的是容器的 build 或 make 方法生成实例
            return $c->$method($concrete);
        };
    }

    // 生成实例对象，首先解决接口和要实例化类之间的依赖关系
    public function make($concrete)
    {
        $concrete = $this->getConcrete($abstract);
        // 这边是递归调用直到变成可 build 的
        if ($this->isBuildable($concrete,$abstract)) {
            $object = $this->build($concrete);
        } else {
            $object = $this->make($concrete);
        }

        return $object;
    }

    protected function isBuildable($concrete,$abstract)
    {
        return $concrete === $abstract || $concrete instanceof Closure;
    }

    // 获取绑定的回调函数
    protected function getConcrete($abstract)
    {
        // 判断是否设置了回调函数
        if (!isset($this->bindings[$abstract])){
            return $abstract;
        }
        return $this->bindings[$abstract]['concrete'];
    }

    // 实例化对象
    public function build($concrete)
    {
        // 判断是否是闭包
        if ($concrete instanceof Closure){
            return $concrete($this);
        }
        // 初始化 ReflectionClass 对象
        $reflector = new ReflectionClass($concrete);
        // 判断是否可以类是否可以实例化
        if(!$reflector->isInstantiable()) {
            echo $message = "Target [$concrete] is not instantiable";
        }
        // 获取类的构造函数 返回一个 ReflectionMethod 对象
        // 如果类不存在构造函数返回 null
        $constructor = $reflector->getConstructor();
        
        if(is_null($constructor)) {
            return new $concrete;
        }
        // 返回参数列表
        $dependencies = $constructor->getParameters();
        $instance = $this->getDependencies($dependencies);
        // 创建一个类的新实例，给出的参数将传递到类的构造函数
        return $reflector->newInstanceArgs($instance);
    }
    
    // 解决通过反射机制实例化对象时的依赖
    protected function getDependencies($parameters)
    {
        $dependencies = [];
        foreach ($parameters as $parameter)
        {
            // 获取参数类型提示类
            $dependency = $parameter->getClass();
            if(is_null($dependency)) {
                $dependencies[] = NULL;
            } else {
                $dependencies[] = $this->resolveClass($parameter);
            }
        }
        return (array)$dependencies;
    }

    // 解析类
    protected function resolveClass(ReflectionParamter $parameter)
    {
        return $this->make($parameter->getClass()->name);
    }
}

class Traveller
{
    protected $trafficTool;

    public function __construct(Visit $trafficTool)
    {
        $this->trafficTool = $trafficTool;
    }

    public function visitTibet()
    {
        $this->trafficTool->go();
    }
}

$app = new Container();

$app->bind('Visit','Train');
$app->bind('traveller','Traveller');

$tra = $app->make('traveller');
$tra->visitTibet();
```




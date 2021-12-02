## swoft 注解类

| 注解标记          | 说明                                                         |
| ----------------- | ------------------------------------------------------------ |
| @Annotation       | 表明此类事一个注解类                                         |
| @Target()         | 表示该类的注解级别，可填参数 ALL(可在类、属性和方法使用)、CLASS、METHOD、PROPERTY、ANNOTATION |
| @var type         | PHPDoc 标准 tag，会根据该类型对参数进行检查                  |
| @Controller()     | 指定路由访问的模块，不填默认是类名 注意：@Controller 实现了 @Bean 的功能，不能同时使用这两个注解 |
| @RequestMapping() | 指定访问类方法的 URI                                         |
| @Bean()           | 在框架初始化的时候就实例化有 @Bean 标记的类，默认是单例      |
| @Inject()         | 注入由 @Bean 实例化的类                                      |
| @Middleware()     | 在类上则作用于整个类，在方法上则作用于方法                   |
| @Middlewares()    | 配置一组 @Middleware 按照顺序执行                            |
| @View()           | 使用视图                                                     |
| @WebSocket()      | 表明允许 ws 连接                                             |
| @Value()          | 给属性赋值，@Value("ceshi") === @Value(value="ceshi") ，赋值顺序是  env > name > value |
| @Task()           | 定义一个任务类，参数 name 默认类名且必须唯一，coroutine 是否以协程运行任务 |
| @Scheduled        | 定义一个任务的执行时间：秒 分 时 日 月 周，定时任务所在的 task 任务不能以协程方式运行 |
| @Process()        | 前置进程必须设置 boot = true 通常放置于 app/Boot 目录 自定义进程通常放置于 app/process 目录，用户自定义进程使用 ProcessBuilder::create('name')->start() 创建它 |



### 数据库注解标签

| 注解      | 说明                                                         |
| --------- | ------------------------------------------------------------ |
| @Entity() | 使用哪个数据配置，默认 default                               |
| @Table()  | 定义数据库表名                                               |
| @Column   | 参数：name 定义表字段，type 定义数据类型及更新是验证类型，所有字段属性必须设置 getter 和 setter 方法 |
| @Id       | 主键                                                         |



### RPC 注解

| 注解       | 说明                                                         |
| ---------- | ------------------------------------------------------------ |
| @Pool()    | rpc 服务连接池                                               |
| @Breaker() | rpc 熔断器名称，名称需要和连接池一致                         |
| @Reference | name 定义引用哪个服务接口，version 使用该服务版本，pool 定义使用哪个连接池，breaker 定义使用哪个熔断器，packer |
| @Fallback  | 定义 RPC 服务降级处理接口，参数：name 默认类名               |



```php
/**
*
*	@Annotation 使用关键字声明是一个注解类
*	@Target("ALL") 类注解 为 Doctrine\Common\Annotations\Annotation\Target.php，表示该注解使用的级别，类注解还是方法注解，参数可填 ALL(表示可以在类、属性和方法上使用)、CLASS、METHOD、PROPERTY、ANNOTATION
*/
class Test
{
     /**
     * @var string // @var 是 PHPDoc 标准的常用 tag，定义了属性的类型 Doctrine 会根据该类型额外对注解参数进行检查
     */
    private $name = '';
    
    /**
     *  如果注解类提供构造器，Doctrine 会调用，一般会在此处对注解类对象的 private 属性进行赋值
     * 
     * @param array $values // Doctrine 注解使用处的参数数组
     */
    public function __construct(array $value)
    {
        if(isset($value['name'])){
            $this->name = $value['name'];
        }
    }
}

# 使用注解类需要 use

use App\Module\Test\Annotation\Test;

class IndexController
{
    /**
    *
    *  @var string // @var 是 PHPDoc 标准的常用 tag，定义了属性的类型，Doctrine 会根据该类型额外对注解参数进行检查
    */
    private $name = "";
    /**
    *
    * @Test(name="777")  这里相当于 new Test([name="777"])
    *
    */
    public function __construct()
    {
        
    }
}
```


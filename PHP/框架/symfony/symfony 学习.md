# 路由

使用注解路由

`composer require annotation`

使用注解路由的控制器需要  `use Symfony\Component\Routing\Annotation\Route`

```php
/**
* @Route("/blog/{page}")
*/
// 如果给 page 定义个默认参数，使用 defaults={"page"=1}

// 给路由增加一个匹配条件，使用 condition={"contenxt.getMethod() in ['GET','HEAD'] and request.headers.get('User-Agent') matches '/firefox/i'"}

// 使用 requirements 可以对参数进行规则匹配
requirements={"page"="\d+"}
// 对路由命名使用 name 参数
name="blog_page"
// 指定路由方法
methods={"GET"}
```

在模板内可以使用 **path** 快捷函数生成路由

# 控制器

基类控制器 `Symfony\Bundle\FrameworkBundle\Controller\AbstractController`

# 模板

### 模板语法

`{{...}}` 打印变量或者一个表达式的结果

`{%...%}` 控制模板的逻辑标记，使用它执行语法

`{#...#}` 注释



#### 模板的tag

`{% extend xxxx %}` 继承某个模板

`{% block xxxx%}xxxx{% endblock %}` 继承的子模板可修改区域

`{{ parent() }}` 从父模板中获取 block  的内容

`{{ include() }}` 包含一个其他模板

`{{ path() }}` 输出要给相对路径，`{{ url() }}` 输出绝对路径

`composer require symfony/asset`  `{{ asset() }}` 引入资源路径



# 配置

使用 `%env(name)%` 解析一些全局配置

配置文件加载顺序：

- `config/packages/*.yaml` 
- `config/packages/<environment-name>/*.yaml`
- `config/packages/services.yaml`

`config/packages/framework.yaml` 全局加载配置文件

`config/packages/prod/framework.yaml` 生产环境加载的配置文件，没有则不加载

`config/packages/dev/framework.yaml` 开发环境加载配置文件，没有则不加载

`config/packages/test/framework.yaml` 测试环境配置文件

多 env 文件

- `.env.<environment>.local` 为 environment 并且只在本地机器定义或覆盖一些变量
- `.env.local` 定义或覆盖变量只在本地机器
- `.env` 定义全局的变量



# 数据库和 Doctrine ORM

使用 Doctrine 

```
composer require symfony/orm-pack
composer require --dev symfony/maker-bundle
```

针对 ORM 的扩展包：

- symfony ORM debug `composer require --dev symfony/profiler-pack`
- 自动获取一条数据 `composer require sensio/framework-extra-bundle`

创建一个数据库 `php bin/console doctrine:database:create`

创建一个表实体 `php bin/console make:entity` 如果想要添加一个属性可以再次执行这个命令，然后继续执行后面的指令，Doctrine 会管理一个 `migration_versions`

创建一个 migrate 文件 `php bin/console make:migration`

运行 migrate `php bin/console doctrine:migrations:migrate`

重新生成所有的表实例 `php bin/console make:entity --regenerate`
---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Yaf学习
---
|              常量 (启用命名空间后的常量名)               |                          说明                          |
| :-----------------------------------------------------: | :----------------------------------------------------: |
|                YAF_VERSION(Yaf\VERSION)                 |                 Yaf 框架的三位版本信息                  |
|                 YAF_ENVIRON(Yaf\ENVIRON                 | Yaf 的环境常量, 指明了要读取的配置的节, 默认的是 product |
|     YAF_ERR_STARTUP_FAILED(Yaf\ERR\STARTUP_FAILED)      |        Yaf 的错误代码常量, 表示启动失败, 值为 512        |
|       YAF_ERR_ROUTE_FAILED(Yaf\ERR\ROUTE_FAILED)        |        Yaf 的错误代码常量, 表示路由失败, 值为 513        |
|    YAF_ERR_DISPATCH_FAILED(Yaf\ERR\DISPATCH_FAILED)     |        Yaf 的错误代码常量, 表示分发失败, 值为 514        |
|     YAF_ERR_NOTFOUND_MODULE(Yaf\ERR\NOTFOUD\MODULE)     |    Yaf 的错误代码常量, 表示找不到指定的模块, 值为 515    |
| YAF_ERR_NOTFOUND_CONTROLLER(Yaf\ERR\NOTFOUD\CONTROLLER) | Yaf 的错误代码常量, 表示找不到指定的 Controller, 值为 516 |
|     YAF_ERR_NOTFOUND_ACTION(Yaf\ERR\NOTFOUD\ACTION)     |   Yaf 的错误代码常量, 表示找不到指定的 Action, 值为 517   |
|       YAF_ERR_NOTFOUND_VIEW(Yaf\ERR\NOTFOUD\VIEW)       |  Yaf 的错误代码常量, 表示找不到指定的视图文件, 值为 518  |
|        YAF_ERR_CALL_FAILED(Yaf\ERR\CALL_FAILED)         |        Yaf 的错误代码常量, 表示调用失败, 值为 519        |
|    YAF_ERR_AUTOLOAD_FAILED(Yaf\ERR\AUTOLOAD_FAILED)     |     Yaf 的错误代码常量, 表示自动加载类失败, 值为 520     |
|         YAF_ERR_TYPE_ERROR(Yaf\ERR\TYPE_ERROR)          |   Yaf 的错误代码常量, 表示关键逻辑的参数错误, 值为 521   |

Yaf 的配置选项

| 选项名称             | 默认值  | 可修改范围     | 更新记录                                                     |
| -------------------- | ------- | -------------- | ------------------------------------------------------------ |
| yaf.environ          | product | PHP_INI_ALL    | 环境名称, 当用 INI 作为 Yaf 的配置文件时, 这个指明了 Yaf 将要在 INI 配置中读取的节的名字 |
| yaf.library          | NULL    | PHP_INI_ALL    | 全局类库的目录路径                                           |
| yaf.cache_config     | 0       | PHP_INI_SYSTEM | 是否缓存配置文件 (只针对 INI 配置文件生效), 打开此选项可在复杂配置的情况下提高性能 |
| yaf.name_suffix      | 1       | PHP_INI_ALL    | 在处理 Controller, Action, Plugin, Model 的时候, 类名中关键信息是否是后缀式, 比如 UserModel, 而在前缀模式下则是 ModelUser |
| yaf.name_separator   | ""      | PHP_INI_ALL    | 在处理 Controller, Action, Plugin, Model 的时候, 前缀和名字之间的分隔符, 默认为空, 也就是 UserPlugin, 加入设置为 "_", 则判断的依据就会变成:"User_Plugin", 这个主要是为了兼容 ST 已有的命名规范 |
| yaf.forward_limit    | 5       | PHP_INI_ALL    | forward 最大嵌套深度                                          |
| yaf.use_namespace    | 0       | PHP_INI_SYSTEM | 开启的情况下, Yaf 将会使用命名空间方式注册自己的类, 比如 Yaf_Application 将会变成 Yaf\Application |
| yaf.use_spl_autoload | 0       | PHP_INI_ALL    | 开启的情况下, Yaf 在加载不成功的情况下, 会继续让 PHP 的自动加载函数加载, 从性能考虑, 除非特殊情况, 否则保持这个选项关闭 |

Yaf 的可选配置项

|                   名称                   | 值类型 |               默认值               |                             说明                             |
| :--------------------------------------: | :----: | :--------------------------------: | :----------------------------------------------------------: |
|             application.ext              | String |                php                 |                       PHP 脚本的扩展名                        |
|          application.bootstrap           | String |       Bootstrapplication.php       |                   Bootstrap 路径 (绝对路径)                    |
|           application.library            | String | application.directory + "/library" |                 本地 (自身) 类库的绝对目录地址                 |
|           application.baseUri            | String |                NULL                | 在路由中, 需要忽略的路径前缀, 一般不需要设置, Yaf 会自动判断. |
|   application.dispatcher.defaultModule   | String |               index                |                          默认的模块                          |
|  application.dispatcher.throwException   |  Bool  |                True                |                  在出错的时候, 是否抛出异常                  |
|  application.dispatcher.catchException   |  Bool  |               False                | 是否使用默认的异常捕获 Controller, 如果开启, 在有未捕获的异常的时候, 控制权会交给 ErrorController 的 errorAction 方法, 可以通过 $request->getException() 获得此异常对象 |
| application.dispatcher.defaultController | String |               index                |                         默认的控制器                         |
|   application.dispatcher.defaultAction   | String |               index                |                          默认的动作                          |
|           application.view.ext           | String |               phtml                |                        视图模板扩展名                        |
|           application.modules            | String |               Index                | 声明存在的模块名, 请注意, 如果你要定义这个值, 一定要定义 Index Module |
|           application.system.*           | String |                 *                  | 通过这个属性, 可以修改 yaf 的 runtime configure, 比如 application.system.lowcase_path, 但是请注意只有 PHP_INI_ALL 的配置项才可以在这里被修改, 此选项从 2.2.0 开始引入 |

Yaf 支持的 hook

| 触发顺序 | 名称                 | 触发时机               | 说明                                                         |
| -------- | -------------------- | ---------------------- | ------------------------------------------------------------ |
| 1        | routerStartup        | 在路由之前触发         | 这个是 7 个事件中, 最早的一个. 但是一些全局自定的工作, 还是应该放在 Bootstrap 中去完成 |
| 2        | routerShutdown       | 路由结束之后触发       | 此时路由一定正确完成, 否则这个事件不会触发                   |
| 3        | dispatchLoopStartup  | 分发循环开始之前被触发 |                                                              |
| 4        | preDispatch          | 分发之前触发           | 如果在一个请求处理过程中, 发生了 [forward](http://www.laruence.com/manual/yaf.class.controller.forward.html), 则这个事件会被触发多次 |
| 5        | postDispatch         | 分发结束之后触发       | 此时动作已经执行结束, 视图也已经渲染完成. 和 preDispatch 类似, 此事件也可能触发多次 |
| 6        | dispatchLoopShutdown | 分发循环结束之后触发   | 此时表示所有的业务逻辑都已经运行完成, 但是响应还没有发送     |


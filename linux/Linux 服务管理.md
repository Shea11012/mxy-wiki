在 linux 中，服务分为独立守护进程和超级守护进程。

独立守护进程是自行监听在后台的，基本上所有的服务都是独立守护进程类的服务。

超级守护进程专指 xinted 服务，这个服务代为管理一些特殊的服务，这类服务在被请求时才会由 xinted 通知它启动服务，服务提供完毕后就关闭服务，这类服务被称为瞬时守护进程。（超级守护进程 xinted 本身是一个常驻内存的独立守护进程，因为要监听来自外界对其管理的瞬时守护进程的请求，一般不工作时，xinted 不占用端口号，在工作时占用被请求的瞬时守护进程的端口号，并处于监听状态）。

### 管理服务的开机自启动

chkconfig ：管理 /etc/init.d 目录下存在且脚本的内容满足一定条件的服务。

要让 chkconfig 管理服务的开机是否自启动行为，只需将脚本放在 /etc/init.d 目录下，然后再脚本的前部加上 chkconfig 行和 description 行 ：

```
#!/bin/bash

# chkconfig: - 85 15
# description: The Nginx HTTP Server is an efficieent and extensible
```

chkconfig 和 description 两行必须得被注释，且必须在所有非注释行的前面。

chkconfig ： `-` 表示适用于运行级别 123456 上，85 表示开机启动时的顺序，15 表示关机停止服务时的顺序

description ： 可以随便给一点描述信息，但必须得给 description: 关键字

### 管理 xinted 及相关瞬时守护进程

该类服务不能直接使用 service 命令来启动，需要去 /etc/xinted.d/ 目录或 /etc/xinetd.conf 中配置
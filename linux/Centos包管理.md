---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: Centos包管理
---
### 包名称

在 Rhel/Centos/Fedora 上，包的名称以 rpm 结尾，分为二进制包和源码包。

源码包以 `.src.rpm` 结尾，它是未编译的包，可以自行进行编译或者用其制作自己的二进制 rpm 包。

非 `.src.rpm` 结尾的包都是二进制包，都已经编译完成。

一个 rpm 包全名的意义，如 `httpd-2.2.15-39.el6.centos.x86_64.rpm` 

httpd ：包名

2.2.15 ：版本号，版本号格式 `[主版本号.[次版本号.[修正号]]]`

39 ：软件发布次数

el6.centos ：适合的操作系统平台以及适合的操作版本

x86_64 ：适合的硬件平台，硬件平台根据 CPU 来决定，有 i386、i586、i686、x86_64、noarch 或省略（标识不区分硬件平台）

rpm ：软件包后缀扩展名

### 主包和子包

在制作 rpm 包时，很多时候都将其按功能分割成多个子包，如客户端包，服务端包等

如：MySQL 这个程序，它有一下几个包

- mysql-server.x86_64 ：服务端包
- mysql.x86_64 ：客户端包
- mysql-bench.x86_64 ： 压力测试包
- mysql-libs.x86_64 ：库文件包
- mysql-devel.x86_64 ：头文件包

### rpm 管理包

rpm 包被安装后，会在 `/var/lib/rpm` 下建立已装 rpm 数据库，以后有任何 rpm 的升级、查询、版本比较等包的操作都是从这个目录下获取信息并完成相应的操作

### 安装包后的文件分布

rpm 安装完成后，相关文件会复制到多个目录下（具体路径是在制作 rpm 包时指定）。一般来说，分布形式差不多如下

| 目录  | 说明               |
| ----- | ------------------ |
| /etc/ | 房子配置文件的目录 |
| /bin、/sbin、/usr/bin 或 /usr/sbin | 一些可执行文件 |
| /lib、/lib64、/usr/lib（/usr/lib64）|一些库文件|
|/usr/include|一些头文件|
|/usr/share/doc | 一些基本的软件使用手册与帮助文件 |
| /usr/share/man | 一些 man page 档案 |

### rpm 查询

rpm 安装需要人为的解决包的依赖关系，所以很少使用，基本上使用 yum

rpm -q ： 查询，详细命令参数，查看手册

### yum 包管理

yum 的默认配置文件 `/etc/yum.conf`

### 配置 yum 仓库

配置文件为 `/etc/yum.conf` 和 `/etc/yum.repos.d/` 中的 `.repo` 文件，其中 `/etc/yum.conf` 配置的是仓库的默认项，一般配置源都是在 `/etc/yum.repos.d/*.repo` 中配置 ，该文件夹下除了 CentOS-Base.repo，其他都可以删掉。

repo 文件配置格式如下：

> [base] ：仓库 ID，必须保证唯一性
>
> name ：仓库名称，可随意命名
>
> mirrorlist ：该地址下包含了仓库地址列表，包含一个或多个镜像站点，和 baseurl 使用一个就可以了
>
> baseurl ：仓库地址，网络上的地址则写网络地址，本地地址则写本地地址，格式为 `file://路径` ，如 `file:///mnt/cdrom`
>
> gpgcheck=1 ：指定是否需要签名，1 表示需要，0 表示不需要
>
> gpgkey= ：签名文件的路径
>
> enable ：该仓库是否生效，enable=1 表示生效，enable=0 表示不生效
>
> cost= ： 开销越高，优先级越低
>
> 【repo 配置文件中可用的宏】
>
> $releasever ：程序的版本，对 Yum 而言指的是 redhat-relrase 版本，只替换为主版本号
>
> $arch : 系统架构
>
> $basharch ：系统基本架构，如 i686，i586 等的基本架构为 i386
>
> $YUM0-9 ：在系统定义的环境变量，可与在 yum 中使用

### repo 配置示例：配置 epel 仓库

epel 由 fedora 社区维护的高质量高可靠性的安装源

1. 安装 epel-release-noarch.rpm

`rpm -ivh epel-release-latest-6.noarch.rpm`

安装完成后会在 `/etc/yum.repo.d/` 目录下生成两个 epel 相关的 repo 文件，其中一个是 epel.repo 此文件中的 epel 源设置在 fedora 的镜像站点上，会比较慢，修改为如下

```
[epel]
name=Extra Packages for Enterprise Linux 6 - $basearch
baseurl=http://mirrors.sohu.com/fedora-epel/6Server/$basearch/
#mirrorlist=https://mirrors.fedoraproject.org/metalink?repo=epel-6$arch=$basearch
failovermethod=priority
enabled=1
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM/GPG-KEY-EPEL-6
```

2. 直接增加 epel 仓库

在 `/etc/yum.repos.d/` 下任意 repo 文件中添加上 epel 仓库即可

```
[epel]
name=epel
baseurl=http://mirrors.sohu.com/fedora-epel/6Server/$basearch/
enabled=1
gpgcheck=0
```

然后清除缓存再建立缓存即可

`yum clean all;yum makecache`

### 源码编译安装程序

1. 解压进入后阅读其中 `INSTALL/README` 文件
2. 执行 `./configure` 或带有编译选项的 `./configure` ，并将定义好的安装配置写入和系统环境信息写入 Makefile 文件中
3. 执行 make 命令进行编译，编译工作主要是调用编译器将源码编译为可执行文件
4. make install 将编译的数据复制到指定目录下
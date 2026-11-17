---
title: "RPM 与 YUM 包管理"
created: 2026-07-22
updated: 2026-07-22
status: needs-review
type: synthesis
tags:
  - linux
  - linux/package-management
confidence: INFERRED
---
# RPM 与 YUM 包管理

> 本页由重复旧笔记整合而来；正文按原页面分节保留，尚未重新核验技术时效性。

## RPM 包结构与仓库基础

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

### YUM 仓库示例：配置 EPEL

> 通用 YUM 配置字段和命令见后文“YUM 命令与配置参考”。

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

## RPM 命令参考

# rpm 包命令

1. rpm 包安装

`rpm {-i | --install} [install-options] package-file` 

- options:
  - `-i，--install` ：安装
  - `-v`：verbos，输出详细信息
  - `-vv`：verbos，输出更详细信息
- `install-options`
  - `-h`：以 hash marks 格式输出进度条，每个 # 代表 2% 的进度
  - `--test`：测试安装，只做环境检查，并不真正安装
  - `--nodeps`：忽略程序依赖关系
  - `--replacepkgs`：覆盖安装，如果文件修改错误，需要将其找回，使用此方法，需要把修改错误的文件提前删除
  - `--justdb`：不执行安装操作，只更新数据库
  - `--noscripts`：不执行 rpm 自带的所有脚本
  - `--nosignature`：不检查包签名信息，即不检查来源合法性
  - `--nodigest`：不检查包完整性信息

2. RPM 包升级

- `rpm {-U | --upgrade} [install-options] package_file`
- `rpm {-F | --freshen} [install-options] package_file`
- `options`
  - `-U`：升级或安装
  - `-F`：升级但不安装
- `install-options`
  - `--oldpackage`：降级
  - `--force`：强制升级

3. RPM 卸载

- `rpm {-e | --erase} [--allmatches] [--nodeps] [--noscripts]`
- `options`
  - `-e`：删除指定程序
  - `--allmatches`：卸载所有匹配指定名称的程序包的各个版本
  - `--nodeps`：卸载时忽略依赖关系
  - `--test`：测试卸载，dry run 模式

4. RPM 包查询

`rpm {-q | --query} [select-options] [query-options]`

- `options`

  - `-qa,-all`：查询所有已安装的包
  - `-f，--file File`：查询指定的文件是由哪个包安装生成的
  - `-g，--group GROUP`：查询指定包由哪个包组提供
  - `-p，--package PACKAGE_FILE`：对未安装的程序包执行查询操作
  - `--whatprovides CAPABILITY`：查询指定的 capability 由哪个程序包提供
  - `--whatrequires CAPABILITY`：查询指定的 capability 被哪个包所依赖

- `query-options`

  - `--changelog`：查询 rpm 包的 changelog
  - `-l，--list`：列出程序包安装生成的所有文件列表
  - `-i，--info`：查询程序包的 infomation，包括版本号、大小和所属的包组等信息
  - `-c，--configfiles`：查询指定的程序包提供的配置文件
  - `-d，--docfiles`：列出指定的程序包提供的文档
  - `--provides`：列出指定的程序包提供的所有 capability
  - `-R，--requires`：查询指定程序包的依赖关系
  - `--script`：查查看程序包自带的脚本

- 校验程序包后的 format

  - 例：`rpm -V zsh`
  - `S.5....T.  d /usr/share/doc/zsh-4.3.11/README`

  说明：

  > S 文件大小改变
  >
  > M 文件权限改变
  >
  > 5 MD5 校验码改变
  >
  > D 主次设备不匹配
  >
  > L read link 变化
  >
  > U 属主改变
  >
  > G 属组改变
  >
  > T 修改时间改变
  >
  > P CAPABILITY 改变

  改变才会显示相应值，`.` 代表未修改过

5. RPM 包校验

- `options`
  - `-V package_file`：自动检查其数据的完整性及合法性
  - `-K package_file`：收到检查其数据的完整性及合法性
  - `--import key`：手动导入 GPG KEY

## YUM 命令与配置参考

# yum
#### yum 安装和卸载

`rpm -i yum-version.rpm`

`rpm -e yum`

## yum 配置文件及格式

**/etc/yum.conf**

- 格式如下：
  - [main]：主名称，固定名称
  - cachedir=：缓存目录
  - keepcache=0：是否保存缓存
  - exactarch=1：是否做精确严格平台匹配
  - gpgcheck=1：检查来源合法性和完整性
  - plugins=1：是否支持插件
  - installonly_limit：同时安装几个

**/etc/yum.repos.d/*.repo**

- 仓库的指向及其配置，格式如下：

  - [repository ID]：ID 名称，即仓库名称，不可与其他 ID 重名

  - name=：对 ID 名称说明

  - baseurl=URL1

    URL2

    URL3（如果同一个源有多个镜像，可以在此写几个，每个 URL 需换行）

  - mirrorlist=（有一台服务器在网络上，保存了多个 baseur ，使用此项，就不使用 baseurl）

  - enabled={1|0}

  - gpgcheck={1|0}

  - repo_gpgcheck=：检查仓库的元数据的签名信息

  - gpgkey=URL（gpg 密钥文件）

  - enablegroups={1|0} 是否在此仓库使用组来管理程序包

  - failovermethod=roundrobin|priority（对多个 baseurl 做优先级，roundrobin 为轮询，priority 为优先级，默认为轮询，意为随机）

  - keepalive=：如果对方是 http1.0 是否要保持连接

  - username=yum 的验证用户

  - password=yum 的验证用户密码

  - cost=默认 baseurl 都为 1000

## yum 命令

`yum [options] [command] [package..]`

- options：
  - `--nogpgcheck`：禁止进行 gpg check
  - `-y`：自动回答 Yes
  - `--enablerepo`：临时启用此处指定的 repo
  - `--disablerepo`：临时禁用此处指定的 repo
  - `--noplugins`：禁用所有插件

**yum 程序安装**

- `install package`：安装
- `reinstall package`：重新安装
- `downgrade package`：降级安装
- `localinstall package`：安装本地程序

**yum 包升级**

- `update softname`
- `localupdate package`

**yum 检查升级**

- `check-update`

**yum 程序卸载**

- `remove | erase softname`

**yum 显示程序包**

- `list {all | available | updates | installed}`
  - all：显示所有仓库的包
  - available：显示可用的软件包
  - updates：显示可用于升级的包
  - installed：显示已经安装的包
  - `yum list php*`：显示以 php 开头的所有依赖包

**yum 查看包的 information 信息**

- `info package`

**yum 查看文件是由哪个包提供**

- `provides package | file`

**yum 清理本地缓存**

- `clean [package | metadata | expire-cache | rpmdb | plugin | all ]`

**yum 生成缓存**

- `makecache`

**yum 搜索包名及 summary 信息**

- `search [string]`

**yum 显示程序包的依赖关系**

- `deplist package`

**yum 查看事务历史（事务只记录安装、升级、卸载的信息）**

- `history [info | list | packages-list | packages-info | summary | addon-info | redo | undo | rollback | new | sync | stats]`

**yum 显示仓库列表**

- `repolist [all | enabled | disabled]`
  - all：查看全部仓库
  - enabled：查看可用的仓库
  - disabled：查看不可用仓库

#### yum 组管理

**yum 组安装**

- `groupinstall`

**yum 组查看**

- `grouplist`

**yum 组的基本信息查看**

- `groupinfo`

**yum 组删除**

- `groupremove`

**yum 组更新**

- `groupupdate`

[yum 命令参考链接](

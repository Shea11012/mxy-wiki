---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: yum
---
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
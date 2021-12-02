### Ubuntu 软件安装方法

1. apt 方式
   1. 普通安装：apt-get install softname1 softname2 ....
   2. 修复安装：apt-get -f install  softname1 softname2 ....
   3. 重新安装：apt-get -reinstall install softname1 softname2 ...
2. Dpkg 方式
   1. dpkg -i package_name.deb

### Ubuntu 软件卸载方法

1. apt 方式
   1. 移除式卸载：apt-get remove softname1 softname2....
   2. 清除式安装：apt-get --purge remove softname1 softname2 ... （同时清楚配置）
2. dpkg 方式
   1. 移除式卸载：dpkg -r pkg1 pkg2 ...
   2. 清除式卸载：dpkg -P pkg1 pkg2 ...

### Ubuntu 软件包查询方法

dpkg -l softname：查看系统中软件包的状态

dpkg -s softname：查看软件包的详细信息

dpkg-query -L softname：查询系统中属于 softname 的文件

### apt 使用总结

```shell
apt-cache search # 搜索包
apt-cache show # 获取包的相关信息
apt-get install # 安装包
apt-get --reinstall install # 重新安装包
apt-get -f install # 强制安装
apt-get remove # 删除包
apt-get remove --purge # 删除包，删除配置文件
apt-get autoremove --purge # 删除包及其依赖的软件包和配置文件等
apt-get update # 更新源
apt-get upgrade # 更新已安装的包
apt-get dist-upgrade # 升级系统
apt-get dselect-upgrade # 使用 select 升级
apt-cache depends # 了解使用依赖
apt-cache rdepends # 了解，某个具体的依赖（查看该包被那些包依赖）
apt-get build-dep # 安装相关的编译环境
apt-get source # 下载该包的源代码
apt-get clean && apt-get autoclean  # 清理下载文件的存档 && 只清理过时的包
apt-get check # 检查是否有损坏的依赖
apt-file search filename # 查找 filename 属于那个软件包
apt-file list packagename # 列出软件包的内容
apt-file update # 更新 apt-file 数据库

dpkg --info "软件包名" # 列出软件包解包后的包名称
dpkg -l # 列出当前系统中所有的包
dpkg -s # 查询已安装的包的详细信息
dpkg -L # 查看系统中已安装的软件包所安装的位置
dpkg -S # 查询系统中某个文件属于哪个软件包
dpkg -I # 查询 deb 包的详细信息，在一个软件包下载到本地之后看需不需要安装
dpkg -i # 手动安装软件包，不能解决软件包之前的依赖问题
dpkg -r # 卸载软件包，配置文件任然存在
dpkg -P # 全部卸载，不能解决软件包的依赖性问题
dpkg -reconfigure # 重新配置
```

### Debian 软件包文件介绍

Debian 系统中所有包的信息都在 /var/lib/dpkg ，其中 /var/lib/dpkg/info 目录中保存了各个软件包的信息及管理文件：

- .conffiles 结尾文件记录软件包的配置列表
- .list 结尾文件记录软件包的文件列表，用户可在文件当中找到软件包文件的具体安装位置
- .md5sums 结尾的文件记录了 md5 信息，用来进行包的验证
- .config 结尾的文件是软件包的安装配置脚本
- .postinst 脚本是完成 Debian 包解开之后的配置工作，通常用来执行所安装软件包相关的命令和服务的重新启动
- .preinst 脚本实在 Debian 解包之前运行，主要作用是停止作用与即将升级的软件包服务知道软件包安装或升级完成
- .prerm 脚本负责停止与软件包关联的 daemon 服务，在删除软件包关联文件之前执行
- .postrm 脚本负责修改软件包链接或文件关联，或删除由它创建的文件

/var/lib/dpkg/available 是软件包的描述信息，包括当前系统中所欲使用的相同的 Debian 安装源中所欲的软件包，还包括当前系统中已经安装和未安装的软件包
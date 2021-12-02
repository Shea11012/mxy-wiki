# Redis安装

1. #### Linux命令安装

ubuntu安装：`sudo apt-get install redis-server`

centos安装需要先安装EPEL源，再从EPEL源安装：`yum -y install epel-release`，`yum -y install redis`

2. Linux源码安装

```
yum install -y gcc-c++ tcl wget		#centos
apt-get install gcc-c++ tcl wget	#Ubuntu
```

```
wget http://download.redis.io/releases/redis-3.0.7.tar.gz
tar xzf redis-3.0.7.tar.gz
cd redis-3.0.7
make
make install
```

3. docker 安装redis：

可以在Dockefile：

`&& pecl install redis-4.0.2 \ &&docker-php-ext-enable redis`

安装完成后，会在`/usr/local/bin`目录下生成几个可执行文件：

- redis-server：redis服务器端启动程序
- redis-cli：redis客户端操作工具
- redis-benchmark：redis性能测试工具
- redis-check-aof：数据修复工具
- redis-check-dump：检查导出工具

复制配置文件到`/etc`目录下：

```
cp redis.conf /etc/
vi /etc/redis.conf
```

**开机启动**

`echo "/usr/local/bin/redis-server /etc/redis.conf" >> /etc/rc.local`

**设置密码和后台运行**

默认情况，访问redis不需要密码，为了增加安全性，需要设置redis密码

```
sudo vi /etc/redis/redis.conf
# 取消注释 requirepass
requirepass 123456
# daemonize 改为yes
daemonize yes

# 修改完重启redis
sudo /etc/init.d/redis-server restart
```

注意不是密码也可以登录redis服务器，但无法执行命令

**允许远程访问**

默认情况下，redis服务器不允许远程访问，只允许本机访问

修改redis.conf配置文件

`# bind 127.0.0.1 将#注释掉`

修改完后redis还会处于保护模式，可以通过终端远程访问，不能通过PHP代码等第三方的方式访问，让第三方访问，需要关闭redis保护模式，或者设置绑定地址，或者设置密码

```
CONFIG SET protected-mode no
CONFIG REWRITE
```





### win下将 redis 变成服务形式

设置服务命令

`redis-server --service-install redis.windows-servicec.conf --loglevel verbose`

常用 redis 服务命令

卸载服务：`redis-server --service-uninstall`

开启服务：`redis-server --service-start`

停止服务：`redis-server --service-stop`
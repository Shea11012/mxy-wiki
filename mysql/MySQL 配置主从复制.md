## 主配置

在 `my.cnf` 内里添加配置

```
[mysqld]
server-id=100 # 此处不要配置为 1，一般配置为 ip 地址的末八位
bin-log = /var/lib/mysql/mysql-bin # 推荐使用绝对路径

### 进入数据库中配置一个用来备库的用户
create user 'repl'@'%' identified by '123456';
grant replication slave,replication client on *.* to 'repl'@'%';
### 查看二进制日志的 position
show master status;
```

## 从配置

在 `my.cnf` 内添加配置

```
[mysqld]
server-id = 102;
bin-log = /var/lib/mysql/mysql-bin
relay-log = /var/lib/mysqk/mysql-relay-bin # 指定中继日志的位置和命名
log-slave-updates = 1 # 允许备库将其重放的事件也记录到自身的二进制日志中，会给备库增加额外的工作，但一定程度上增加了安全性

### 进入数据库启动复制
change master to master_host='host',master_user='user',master_password='password',master_log_file='mysql-bin.000001',master_log_pos=0;

#启动 slave 复制
start slave;

# 查看 slave 状态
show slave status \G
```

**Tips：修改完配置文件需要重启 MySQL 来启用设置**



### 从另一个服务器开始复制

上面的配置是假定主备库均为刚刚安装好且都是默认的数据，也就是说两台服务器上的数据相同，并且知道当前主库的二进制日志。

大多数情况下有一个运行了一段时间的主库，然后用一台新安装的备库与之同步，此时备库还没有数据。

有几种办法来初始化备库或者从其他服务器克隆数据到备库：

- 使用冷备份：

  最基本的方式是关闭主库，把数据复制到备库。然后重启主库，会使用一个新的二进制文件，在备库通过执行 `change master to 指向这个文件的起始处`。这个方法的缺点很明显，在复制数据时需要关闭主库。

- 热备份只支持 MyISAM。

- 使用 mysqldump
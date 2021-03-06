---
date created: 2021-11-28 13:34
date modified: 2021-11-28 13:34
title: Linux 文件权限
---
## Linux 文件权限
### /etc/passwd 文件
/etc/passwd 文件包含了一些用户有关的信息

### /etc/shadow 文件
/etc/shadow 管理登录用户的密码

### 增删改用户
#### 添加用户
`useradd -m test`

#### 删除用户
`userdel test`

#### 修改用户
- usermod：修改用户账户字段，指定主要组以及附加组的所属关系
- passwd：修改已有用户密码
- chpasswd：从文件中读取登录名密码，并更新密码
- chage：修改密码过期日期
- chfn：修改用户账户备注信息
- chsh：修改用户账户的默认登录 shell


### Linux 用户组
#### 创建新组
`groupadd shared`

将一个组添加到用户的组列表
`usermod -G shared test`

```ad-tip
组关系的更改，必须重新登录后才会生效
```


### 文件权限
#### 文件权限符
- \- 表示文件
- d 表示目录
- l 表示链接
- c 表示字符型设备
- b 表示块设备
- n 表示网络设备

#### 默认的文件权限
umask 022

#### 隐藏文件权限属性
隐藏属性 i 设置后，可以让一个文件不能被删除、改名、设置链接，也无法写入或新增数据。该属性只能由 root 设置。
```shell
chattr +i file
lsattr file
```

具体属性可以查看 `man chattr`

#### 特殊标志位 SUID SGID SBIT
Linux 为每个文件和目录存储了 3 个额外的信息位
- **设置用户 UID**：当 s 标志出现在文件拥有者的 x 权限上，此时被称为 set uid。当文件被用户使用时，程序会以文件属主的权限运行。**SUID 仅对二进制文件生效**
- **设置组 SGID**：当 s 在群组的 x 时则称为 set gid。对文件来说，程序会以文件属主的权限运行，SGID 仅对二进制文件生效；对目录来说，目录中创建的新文件会以目录的默认属组作为默认属组。
- **粘着位 SBIT**：受限删除位。作用于其他组权限的位置，标志为 t 。当使用者在该目录下创建文件或目录时，仅有自己与 root 才能删除该文件。SBIT 目前只针对目录有效。

```shell
### 设置 SUID
chmod u+s
### SGID
chmod g+s
### SBIT
chmod o+t
```
大写的 X/S/T 含义如下：
- X：对象是目录或文件已有执行权限，则赋予执行权限。
- S/T：没有执行权限的 UID/GID 和粘置位，只有当文件的拥有者都无法执行时，才会出现。
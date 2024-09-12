---
tags: ["linux 权限"]
date created: 2021-11-28 13:34
date modified: 2023-02-08 00:31
updated: 2024-08-13
---

## 添加用户 useradd

```bash
# 简单添加
useradd test_user

# 指定家目录
useradd -d {path} test_user

# 指定uid
useradd -u 1000 test_user

# 指定gid
useradd -g 1000 test_user

# 没有家目录的用户
useradd -M test_user

# 指定到期时间
useradd -e 2020-05-30 test_user

# 指定登录shell
useradd -s /bin/sh test_user
```

## 删除用户 usedel

```bash
userdel {user}
```

## 修改用户 usermod

```bash
# 改变用户属组
usermod -g {group} test_user

# 改变用户登录名
usermod -l {new} {old}

# 锁定/解锁用户
usermod -L/U test_user

# 把用户追加入一个组
usermod -aG docker test_user
```

## 创建新组

`groupadd shared`

将一个组添加到用户的组列表
`usermod -G shared test`

> [!tip]
> 组关系的更改，必须重新登录后才会生效

## 目录权限操作

chmod ： 修改权限，能够修改权限的人只有文件所有者和超级管理员

chgrp ：更改文件和目录的所属组，要求组已存在，对于链接文件，修改组的作用对象是链接的源文件，而非链接文件本身

chown ： 修改文件所有者和所属组，对于链接文件而言，默认不会穿过链接修改源文件，而是直接修改链接文件本身

## 其余不常用命令

- passwd：修改已有用户密码
- chpasswd：从文件中读取登录名密码，并更新密码
- chage：修改密码过期日期
- chfn：修改用户账户备注信息
- chsh：修改用户账户的默认登录 shell

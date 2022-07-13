---
tags: ["linux","containerd"]
date created: 2022-05-07 14:01
date modified: 2022-05-07 17:07
title: namespace
---

# namespace
linux namespace 是一种 linux kernel 提供的资源隔离方案。系统可以为进程分配不同的 namespace，并保证不同的 namespace 资源独立分配、进程彼此隔离


## namespace种类
| namespace类型     | 隔离资源                         | kernel版本 |
| ----------------- | -------------------------------- | ---------- |
| IPC namespace     | system V IPC 和 POSIX 消息队列   | 2.6.19     |
| network namespace | 网络设备、网络协议栈、网络端口等 | 2.6.29     |
| PID namespace     | 进程                             | 2.6.14     |
| Mount namespace   | 挂载点                           | 2.4.19     |
| UTS namespace     | 主机名和域名                     | 2.6.19     |
| User namespace    | 用户和用户组                     | 3.8        |


## namespace 常用操作
查看当前系统的namespace
```shell
lsns
```

查看某进程的namespace
```shell
ls -la /proc/<pid>/ns/
```

进入某namespace运行命令
```shell
nsenter -t <pid> -n ip addr
```
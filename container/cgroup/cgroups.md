---
tags: [linux,containerd]
date created: 2022-05-07 17:46
date modified: 2022-05-07 17:51
title: cgroups
---

# cgroups

cgroups(control groups) 是 Linux 下用于对一个或一组进程进行资源控制和监控的机制。

可以对如 CPU 使用时间、内存、磁盘 IO 等进程所需的资源进行限制，不同资源的具体管理工作由相应的 cgroup 子系统来实现。针对不同类型的资源限制，只要将限制策略在不同的子系统上进行关联即可。

cgroups 在不同的系统资源管理子系统中以层级树的方式来组织管理，每个 cgroup 都可以包含其他的子 cgroup，因此子 cgroup 能使用的资源除了受本 cgroup 配置的资源参数限制，还受到父 cgroup 设置的资源限制。
---
date created: 2021-09-19 00:01
date modified: 2021-09-24 22:00
title: podman 遇到的问题
---
#### 配置 podman 的 rootless
开启 `kernel.unprivileged_userns_clone=1`
暂时开启：`sudo sysctl kernel.unprivileged_userna_clone=1`
持久化：`echo kernel.unprivileged_userna_clone=1 > /etc/sysctl.d/userns.conf`

设置 /etc/subuid 和 /etc/subgid
`sudo touch /etc/{subuid,subgid}`
`sudo usermod --add-subuids 100000-165535 --add-subgids 100000-165535 zhuzi`

配置 ping 的权限
`echo net.ipv4.ping_group_range=0 2000000`


参考：
[https://wiki.archlinux.org/title/Podman](https://wiki.archlinux.org/title/Podman)
[https://vadosware.io/post/rootless-containers-in-2020-on-arch-linux/](https://vadosware.io/post/rootless-containers-in-2020-on-arch-linux/)


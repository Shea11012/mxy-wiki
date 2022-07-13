---
tags: ["linux","ssh"]
date created: 2021-06-07 23:50
date modified: 2022-05-07 14:03
title: ssh
---

## 管理多个私钥 免密登录

在 .ssh 文件夹下创建一个 config 文件：

```bash
Host github.com # 别名
HostName github.com
User shea
IdentityFile ～/.ssh/github_id_rsa

Host gitlab.com
HostName gitlab.com
User shea
IdentityFile ~/.ssh/gitlab_id_rsa
ProxyCommand nc -X 5 -x host:port %h %p # proxycommand 可以针对指定的 ssh 实现代理

# 倒入其他目录下的config配置，可以更加灵活的配置
# include 的 config 文件里面需要指定路径的地方都需要配置绝对路径
Include config.d/*  
```

## 服务端 ssh

服务端禁用密码登录，只能使用私钥登录配置：

```bash
# /etc/ssh/sshd_config
PasswordAuthentication no
ChallengeResponseAuthentication no
UsePAM no

# 重启服务
systemcl restart sshd
```

## 端口转发

```bash
ssh -L local_addr:local_port:remote_addr:remote_port middle_host

# 示例
# 在本地指定一个由 ssh 监听的转发端口 2222,将远程host2的端口80映射为本地端口2222
# 当主机连接本地 2222 端口时，ssh 将此端口的数据包转发给中间主机 host3
# 再与远程主机的端口 host2：80 通信
# -g 表示允许外界主机连接本地转发端口，不指定 -g 则默认地址为环回地址
ssh -g -L 2222:host2:80 host3
```
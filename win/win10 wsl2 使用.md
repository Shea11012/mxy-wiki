---
tags: []
date created: 2021-11-30 21:22
date modified: 2023-01-10 04:41
title: win10 wsl2 使用
---

## 获取宿主机 ip

wsl2 因为使用了 hyper-v 的虚拟机，所以和宿主机不在同一个网段中，在使用宿主机的 vpn 时，需要先获取宿主机的 ip：

```bash
cat /etc/resolv.conf | grep 'nameserver' | awk '{print $2}'
```

为了便于使用还可以将这个设置一个脚本自动加载 proxy

## 改变 wsl 的默认登录用户

> wsl -l # 查询版本
>
> \# 如 ubuntu20.04
>
> ubuntu2004 config --default-user shea # 需要保证 shea 这个用户已经被创建在该版本中

## wsl2 与 trojan 搭配使用时遇到的端口被占用问题

启动 Trojan 时提示 socks5 的端口被占用了，使用 `netstat -ano | findstr 51837` 查看占用端口发现，该端口没有被占用，经过一番谷歌，hyper-v 虚拟机占用了该端口。

使用命令查看

```powershell
# 查看系统默认端口占用访问
netsh int ipv4 show dynamicport tcp

# 查看hyper-v 启动后的保留端口范围
netsh interface ipv4 show excludedportrange protocol=tcp  # 根据这个命令的输出，可以看到Trojan的默认socks5的端口是被占用了，所以更改Trojan默认的socks5的端口即可解决
```

## wsl2 使用 git 时每次登录都需要手动添加 key 的解决办法

在没有 key 的时候，可以将 win 上的 key 复制到 wsl 上， `cp -r /mnt/c/Users/<username>/.ssh ~/.ssh`

win 上的 key 权限在 linux 下不适用，需要修改权限

```shell
chmod 700 ~/.ssh
chmod 600 ~/.ssh/id_rsa
```

> 上面步骤如果发生 permission denid，则需要在 /etc/wsl.conf （没有则新建）添加
>
> [automount]
>		enabled = true
> 		options = "metadata,umask=22,fmask=11"

为了不再每次打开一个新 tab 时，都需要手动添加 key：

```
sudo apt install keychain

# 将下面的内容写到 bashrc 内
eval `keychain --eval --agents ssh id_rsa`
```

## wsl2 docker exit 139

```
%userprofile%\.wslconfig # 此文件中写入配置

[wsl2]
kernelCommandLine = vsyscall=emulate
```

## wsl 不能 ping 通宿主机

```shell
# 直接放开 `vEthernet (WSL)` 这张网卡的防火墙
New-NetFirewallRule -DisplayName "WSL" -Direction Inbound -InterfaceAlias "vEthernet (WSL)" -Action Allow
```



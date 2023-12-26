---
tags: []
date created: 2021-11-30 21:22
date modified: 2023-09-06 22:13
title: win10 wsl2 使用
---

## 获取宿主机 ip

wsl2 因为使用了 hyper-v 的虚拟机，所以和宿主机不在同一个网段中，在使用宿主机的 vpn 时，需要先获取宿主机的 ip：

```bash
cat /etc/resolv.conf | grep 'nameserver' | awk '{print $2}'
```

为了便于使用还可以将这个设置一个脚本自动加载 proxy

## 改变 wsl 的默认登录用户

>[!note]
> wsl -l # 查询版本
> \# 如 ubuntu20.04
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

## wsl2 固定 ip

使用 hyper-v 创建一个虚拟交换机
![[attachments/Pasted image 20230906215557.png]]

>[!tip]
> 需要重启电脑。如果使用了透明代理，需要在网络适配器中，找到该交换器修改 ipv4

在宿主机中，添加
```conf
# ~/.wslconfig
[wsl2]
networkingMode=bridged
vmSwitch=WSLBridge
```

在对应的 wsl 中，添加
```conf
# /etc/wsl.conf
[boot]
systemd=true
[network]
generateResolveConf=false
```

执行 `wsl --shutdown`

在 **Arch Linux** 中，添加
```network
# /etc/systemd/network/80-wsl-external.network
[Match]
Name=eth0
[Network]
Description=wsl bridge
Address=192.168.32.24/24
Gateway=192.168.32.100
DNS=192.168.32.100
```

执行 `systemctl enable --now systemd-networkd`

## 配置wsl2局域网可访问
### 获取wsl ip
```powershell
wsl ip -o -4 --brief -json addr show eth0 | ConvertFrom-Json | %{ $_[0].addr_info.local }
```

### 增加需要转发的端口
```powershell
netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport={port} connectaddress={wslIP} connectport={port}
```
### 添加防火墙规则
```powershell
netsh advfirewall firewall add rule name=”Open Port 22 for WSL2” dir=in action=allow protocol=TCP localport=22
```

### 清空转发端口
```powershell
netsh interface portproxy reset
```
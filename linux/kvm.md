---
tags: 
date created: 2021-12-03 20:20
date modified: 2023-02-18 02:51
title: kvm
---

# kvm

## arch 安装 kvm

### qemu 安装

```bash
sudo paru -S qemu qemu-arch-extra ovmf bridge-utils dnsmasq vde2 openbsd-netcat ebtables iptables
```

- ovmf：UEFI bios 和 Secure Boot 启动
- birdge-utils：vm 桥接网络
- vde2：qemu 模拟 Ethernet
- dnsmasq：dns 和 dhcp
- openbsd-netcat：网络测试工具
- ebtables 和 iptables：创建包路由和防火墙

## virt-manager && libvirtd service

virt-manager 是创建 vm 的可视化工具 

virt-viewer 打开 vm 实例

```bash
sudo paru -S virt-manager virt-viewer
```

## 开启服务

```bash
sudo systemctl enable libvirtd.service && sudo systemctl start libvirtd.service
```

## 配置 kvm

```bash
sudo vim /etc/libvirt/libvirtd.conf

unix_sock_group = "libvirt"

unix_sock_rw_perms = "0770"
```

## 给 vm 创建一个桥接网络

```bash
vim /tmp/br10.xml



<network>
  <name>br10</name>
  <forward mode='nat'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='br10' stp='on' delay='0'/>
  <ip address='192.168.72.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.72.50' end='192.168.72.200'/>
    </dhcp>
  </ip>
</network>
```

## 注册桥接网络

```bash
注册网络
sudo virsh net-define /tmp/br10.xml

开启网络
sudo virsh net-start br10
```

## 将当前用户加入到 libvirt 组

```bash
sudo usermod -a -G libvirt $(whoami)
```

## 参考

[kvm 安装视频](https://www.youtube.com/watch?v=itZf5FpDcV0)

[kvm virtmanager](https://boseji.com/posts/manjaro-kvm-virtmanager/)


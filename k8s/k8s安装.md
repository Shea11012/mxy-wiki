---
tags: ["k8s"]
---
# k8s 二进制安装
## 环境

-   centos 8.3
-   master 1
-   node 2

## 基础配置

1.  修改镜像源

```bash
sed -e 's|^mirrorlist=|#mirrorlist=|g' \\
         -e 's|^#baseurl=http://mirror.centos.org/$contentdir|baseurl=https://mirrors.ustc.edu.cn/centos|g' \\
         -i.bak \\
         /etc/yum.repos.d/CentOS-Linux-AppStream.repo \\
         /etc/yum.repos.d/CentOS-Linux-BaseOS.repo \\
         /etc/yum.repos.d/CentOS-Linux-Extras.repo \\
         /etc/yum.repos.d/CentOS-Linux-PowerTools.repo \\
         /etc/yum.repos.d/CentOS-Linux-Plus.repo
```

2.  关闭swap

```bash
swapoff -a && sysctl -w vm.swappiness=0 && sed -ri 's/.*swap.*/#&/' /etc/fstab
```

3.  禁用 firewalld

```bash
systemctl disable --now firewalld
```

4.  关闭selinux

```bash
setenforce 0 && sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config
sestatus # 查看 selinux 状态，centos8 修改后需要reboot
```

5.  同步节点时间

```bash
dnf update
dnf install chrony -y
sed -i'*.bak' 's#^pool .*\\.org#pool time\\.pool\\.aliyun\\.com#' /etc/chrony.conf
systemctl enable --now chronyd
chronyc -a makestep
```

6.  安装docker，master 节点可以不用安装
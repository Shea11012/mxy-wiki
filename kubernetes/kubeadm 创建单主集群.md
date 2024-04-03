---
tags: [arch, k8s]
date created: 2022-05-10 03:09
date modified: 2023-11-26 17:54
title: kubeadm 安装k8s
---

## 环境

- 系统：archlinux
- CPU：4 核
- 内存：16G 内存
- 磁盘：30G 磁盘
- 容器：containerd

# 设置主机名

```shell
sudo hostnamectl set-hostname k8s.master
```

# 修改 hosts

```shell
echo "<自己的机器的当前IP地址> k8s.master" | sudo tee /etc/hosts
```

# 网络

```shell
# /etc/modules-load.d/k8s.conf
overlay
br_netfilter

# /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward = 1

# 执行
# sysctl --system
```

# 安装 containerd

```shell
# 这个包含了containerd，nerdctl等
sudo paru -S nerdctl-full-bin
```

如需要去掉 sudo 前缀，运行
```shell
sudo chmod +s /usr/local/bin/nerdctl
```

## 配置

```shell
sudo sh -c 'containerd config default > /etc/containerd/config.toml'
```

## 修改 cgroup driver

```toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  ...
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true
```

## 启动服务

```shell
sudo systemctl enable --now containerd
```

# 安装依赖

```shell
sudo paru -S ebtables ethtool conntrack-tools socat
```

# 安装 kubeadm、kubelet、kubectl

```shell
paru -S kubeadm kubectl kubelet
```

# 初始化节点

```shell
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
```

检查节点状态和 pod 状态

```shell
kubectl get cs
kubectl get nodes
kubectl get pods -A
```

默认，pod 不会被调度到主节点，需要运行以下命令

```shell
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

安装 calico 网络插件

```shell
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml

xh -d https://docs.projectcalico.org/manifests/custom-resources.yaml
# 需要将custom-resources.yaml里的CIDR配置为与pod网络一致
kubectl create -f custom-resources.yaml
```

# 清理集群

```bash
kubectl drain <node name> --delete-emptydir-data --force --ignore-daemonsets

sudo kubeadm reset

rm $HOME/.kube/config

iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -X

# 如果安装了 ipvsadm
ipvsadm -C

kubectl delete node <node-name>
```
---
tags: [arch,k8s]
date created: 2022-05-10 03:09
date modified: 2022-05-10 17:37
title: kubeadm 安装k8s
---

## 环境
- 系统：archlinux
- CPU：4 核
- 内存：16G 内存
- 磁盘：30G 磁盘

## 前置准备
### 设置主机名

```shell
sudo hostnamectl set-hostname k8s.master
```

### 修改 hosts

```shell
echo "<自己的机器的当前IP地址> k8s.master" | sudo tee /etc/hosts
```

### 容器安装
任选一个
<font style="color: red">tips：</font>国内环境安装时，最好设置代理，从 kubeadm config 里的获取到的镜像不是全部镜像，还会有几个镜像需要下载，不走代理很可能会卡住

#### docker
安装 docker

```shell
sudo pacman -S docker
```

修改配置

```shell
mkdir /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
   },
   "storage-driver": "overlay2"
}
EOF
```

启动服务

```shell
sudo systemctl enable --now docker
```

#### containerd
安装 containerd

```shell
sudo pacman -S containerd
```

配置

```shell
sudo containerd config print default > /etc/containerd/config.toml

# 可以在 config.toml 中配置镜像源
```

启动服务

```shell
sudo systemctl enable --now containerd
```

## 安装依赖

```shell
sudo pacman -S ebtables ethtool conntrack-tools socat
```

## 安装 kubeadm、kubelet、kubectl
> aur 源中 kubeadm，会将 kubelet 作为依赖安装

```shell
yay -S kubeadm kubectl
```

## 初始化节点

```shell
# 先将默认配置导出，再修改
sudo kubeadm config print default > kubeadm.yaml
```

kubeadm.yaml

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
# 修改为本地主机ip地址
  advertiseAddress: 1.2.3.4
  bindPort: 6443
nodeRegistration:
# 根据不同的容器修改为不同的值
  criSocket: unix:///var/run/containerd/containerd.sock
  imagePullPolicy: IfNotPresent
  name: node
  taints: null
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
# 修改镜像源为 registry.aliyuncs.com/google_containers
imageRepository: k8s.gcr.io
kind: ClusterConfiguration
kubernetesVersion: 1.24.0
networking:
  # 新增一项，此项等于命令行参数 pod-network-cid
  # 这里使用的calico，所以pod网络为 192.168.0.0/16
  podSubnet: 192.168.0.0/16
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
scheduler: {}
```

执行初始化

```shell
sudo kubeadm init --config kubeadm.yaml
```

检查节点状态和 pod 状态

```shell
kubectl get cs
kubectl get nodes
kubectl get pods -A
```

可选：如果是单主单节点

```shell
kubectl taint nodes --all node-role.kubernetes.io/master-
```

安装 calico 网络插件

```shell
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml

kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
```

## 删除集群

```bash
kubectl drain <node name> --delete-emptydir-data --force --ignore-daemonsets

sudo kubeadm reset

rm $HOME/.kube/config

iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -X

# 如果安装了 ipvsadm
ipvsadm -C

kubectl delete node <node-name>
```
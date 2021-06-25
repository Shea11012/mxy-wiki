---
tags: ["k8s"]
---
# k8s 二进制安装
## 环境

-   centos 8.3 4.18.0-240.22.1.el8_3.x86_64
-   master：192.168.217.132
-   node1：192.168.217.131
-   node2：192.168.217.133
-   kubernets 1.21.0
-   etcd v3.4.13

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

2. 修改每台机器的 /etc/hosts ，添加主机名和 IP 的对应关系

```bash
cat >> /etc/hosts <<EOF
192.168.217.132 master
192.168.217.131 node-1
192.168.217.133 node-2
EOF
```

3. master生成私钥将公钥分发到node

```bash
ssh-keygen -t rsa
```

```bash
for node_ip in node-1 node-2; do
	ssh-copy-id -i .ssh/id_rsa.pub ${node_ip}
done
```

4. 安装依赖包

```bash
dnf install -y epel-release
dnf install -y conntrack ipvsadm ipset jq sysstat curl iptables libseccomp git
```
5.  关闭swap

```bash
swapoff -a && sysctl -w vm.swappiness=0 && sed -ri 's/.*swap.*/#&/' /etc/fstab
```

6.  禁用 firewalld

```bash
systemctl disable --now firewalld
```

7.  关闭selinux

```bash
setenforce 0 && sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config # reboot后生效
sestatus # 查看 selinux 状态
```

8. 加载内核模块

```bash
modprobe br_netfilter
modprobe ip_vs
```

9. 配置 ipvs

```bash
cat >> /etc/modules-load.d/ipvs.conf <<EOF
ip_vs
ip_vs_lc
ip_vs_wlc
ip_vs_rr
ip_vs_wrr
ip_vs_lblc
ip_vs_lblcr
ip_vs_dh
ip_vs_sh
ip_vs_fo
ip_vs_nq
ip_vs_sed
ip_vs_ftp
ip_vs_sh
nf_conntrack
ip_tables
ip_set
xt_set
ipt_set
ipt_rpfilter
ipt_REJECT
ipip
EOF

systemctl enable --now systemd-modules-load
```

10. 调整文件打开数

```bash
cat >> /etc/security/limits.conf <<EOF
* soft nproc 1000000
* hard nproc 1000000
* soft nofile 1000000
* hard nofile 1000000
* soft memlock unlimited
* hard memlock unlimited
EOF
```
11.  同步节点时间

```bash
dnf update
dnf install chrony -y
sed -i'*.bak' 's#^pool .*\\.org#pool time\\.pool\\.aliyun\\.com#' /etc/chrony.conf
systemctl enable --now chronyd
chronyc -a makestep
```

12. 优化k8s参数

```bash
cat >> /etc/sysctl.d/k8s.conf <<EOF
net.ipv4.ip_forward=1
net.ipv4.tcp_keepalive_time=600
net.ipv4.tcp_keepalive_probes=3
net.ipv4.tcp_keepalive_intvl=15
net.ipv4.tcp_max_tw_buckets=36000
net.ipv4.tcp_tw_reuse=1
net.ipv4.tcp_max_orphans=327680
net.ipv4.tcp_orphan_retries=3
net.ipv4.tcp_max_syn_backlog=16384
net.ipv4.tcp_timestamps=0
net.bridge_bridge-nf-call-iptables=1
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-arptables=1
net.netfilter.nf_conntrack_max=2310720
vm.overcommit_memory=1
vm.panic_on_oom=0
vm.swappiness=0
fs.inotify.max_user_watches=89100
fs.file-max=52706963
fs.nr_open=52706963
EOF

sysctl --system
```

13. 安装容器（可选）
	1.  安装docker，master 节点可以不用安装
	```bash
	dnf install -y dnf-utils device-mapper-persistent-data lvm2 # 安装docker需要的组件
	yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo # 配置源
	dnf update -y && dnf install -y docker-ce docker-ce-cli containerd.io # 安装
	mkdir -p /etc/docker
	cat >> /etc/docker/daemon.json <<EOF
	{
		"registry-mirrors": [
			"https://b0yxfmz5.mirror.aliyuncs.com",
			"https://docker.mirrors.ustc.edu.cn",
			"https://mirror.ccs.tencentyun.com"
		],
		"exec-opts": ["native.cgroupdriver=systemd"],
		"storage-driver": "overlay2",
		"log-driver": "json-file",
		"log-opts": {
			"max-size": "100m"
		}
	}
	EOF
	systemctl daemon-reload && systemctl restart docker && systemctl enable docker

	```
	2. 安装 containerd
	```bash
	dnf install -y dnf-utils device-mapper-persistent-data lvm2
	yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
	yum update -y && yum install -y containerd.io
	containerd config default > /etc/containerd/config.toml
	sed -i'*.bak' 's#sandbox_image = \".*\"#sandbox_image = \"registry.aliyuncs.com/google_containers/pause:3.2\"#' /etc/containerd/config.toml
	```
	
14. 下载二进制包

```bash
wget https://dl.k8s.io/v1.22.0-alpha.3/kubernetes-server-linux-amd64.tar.gz
wget https://github.com/etcd-io/etcd/releases/download/v3.4.13/etcd-v3.4.13-linux-amd64.tar.gz

for node_ip in node-1 node-2; do
	rsync kubernetes-server-linux-amd64.tar.gz ${node_ip}:/root
	rsync etcd-v3.4.13-linux-amd64.tar.gz ${node_ip}:/root
	ssh ${node_ip} "tar -zxvf /root/kubernetes-server-linux-amd64.tar.gz --strip-components=3 -C /usr/local/bin kubernetes/server/bin/kube{let,ctl,-apiserver,-controller-manager,-scheduler,-proxy}"
	ssh ${node_ip} "tar -zxvf /root/etcd-v3.4.13-linux-amd64.tar.gz --strip-components=1 -C /usr/local/bin etcd-v3.4.13-linux-amd64/etcd{,ctl}"
done
```

15. 生成证书
```bash
# 在master节点下载
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -O /usr/local/bin/cfssl
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -O /usr/local/bin/cfssljson
chmod +x /usr/local/bin/cfssl /usr/local/bin/cfssljson

# 在 master 节点创建目录 
mkdir -p /etc/etcd/ssl

# 在所有节点创建,存放证书
mkdir -p /etc/kubernetes/pki

# 克隆该仓库，切换到分支 1.21.0
git clone https://github.com/dotbalo/k8s-ha-install.git
cd k8s-ha-install/pki
```

 生成etcd CA 证书 和 CA 证书的 key
 ```bash
cfssl gencert -initca etcd-ca-csr.json | cfssljson -bare /etc/etcd/ssl/etcd-ca
cfssl gencert -ca=/etc/etcd/ssl/etcd-ca.pem -ca-key=/etc/etcd/ssl/etcd-ca-key.pem -config=ca-config.json -hostname=127.0.0.1,master,192.168.217.132,192.168.217.131,192.168.217.133 -profile=kubernetes etcd-csr.json | cfssljson -bare /etc/etcd/ssl/etcd

# 如果有多个master节点，可以将etcd证书复制到其他节点
for node_ip in node-1 node-2; do
	ssh ${node_ip} "mkdir -p /etc/etcd/ssl";
	for file in etcd-ca-key.pem etcd-ca.pem etcd-key.pem etcd.pem; do
		rsync /etc/etcd/ssl/${file} ${node_ip}:/etc/etcd/ssl/${file};
	done
done
 ```
 
 生成 kubernetes 根证书
 ```bash
 cfssl gencert -initca ca-csr.json | cfssljson -bare /etc/kubernetes/pki/ca
 ```
 
 生成 apiserver 证书
 ```bash
 cfssl gencert -ca=/etc/kubernetes/pki/ca.pem -ca-key=/etc/kubernetes/pki/ca-key.pem -config=ca-config.json -hostname=10.96.0.1,127.0.0.1,kubernetes,kubernetes.default,kubernetes.default.svc,kubernetes.default.svc.cluster,kubernetes.default.svc.cluster.local -profile=kubernetes apiserver-csr.json | cfssljson -bare /etc/kubernetes/pki/apiserver
 ```
 
生成 apiserver 代理证书
```bash
cfssl gencert -initca front-proxy-ca-csr.json | cfssljson -bare /etc/kubernetes/pki/front-proxy-ca

cfssl gencert -ca=/etc/kubernetes/pki/front-proxy-ca.pem -ca-key=/etc/kubernetes/pki/front-proxy-ca-key.pem -config=ca-config.json -profile=kubernetes front-proxy-client-csr.json | cfssljson -bare /etc/kubernetes/pki/front-proxy-client
```

生成controller manager 证书
```bash
cfssl gencert -ca=/etc/kubernetes/pki/ca.pem -ca-key=/etc/kubernetes/pki/ca-key.pem -config=ca-config.json -profile=kubernetes manager-csr.json | cfssljson -bare /etc/kubernetes/pki/controller-manager
```

16. kubectl 配置 controller-manger
```bash
# 设置集群参数
kubectl config set-cluster kubernetes \
--certificate-authority=/etc/kubernetes/pki/ca.pem \
--embed-certs=true \
--server=http://192.168.217.132:8843 \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置客户端参数
kubectl config set-credentials system:kube-controller-manager@kubernetes \
--client-certificate=/etc/kubernetes/pki/controller-manager.pem \
--client-key=/etc/kubernetes/pki/controller-manager-key.pem \
--embed-certs=true \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置上下文参数
kubectl config set-context system:kube-controller-manager@kubernetes \
--cluster=kubernetes \
--user=system:kube-controller-manager \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置默认上下文
kubectl config use-context system:kube-controller-manager@kubernetes \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig
```

生成 scheduler 证书
```bash
cfssl gencert \
-ca=/etc/kubernetes/pki/ca.pem \
-ca-key=/etc/kubernetes/pki/ca-key.pem \
-config=ca-config.json \
-profile=kubernetes \
scheduler-csr.json | cfssljson -bare /etc/kubernetes/pki/scheduler
```

kubectl 配置 scheduler
```bash
# 设置集群参数
kubectl config set-cluster kubernetes \
--certificate-authority=/etc/kubernetes/pki/ca.pem \
--embed-certs=true \
--server=https://192.168.217.132:8843 \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

# 设置客户端参数
kubectl config set-credentials system:kube-scheduler \
--client-certificate=/etc/kubernetes/pki/scheduler.pem \
--client-key=/etc/kubernetes/pki/scheduler-key.pem \
--embed-certs=true \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig


# 设置上下文参数
kubectl config set-context system:kube-scheduler@kubernetes \
--cluster=kubernetes \
--user=system:kube-scheduler \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

# 设置默认上下文
kubectl config use-context system:kube-scheduler@kubernetes \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

```

生成 admin 证书
```bash
cfssl gencert \
-ca=/etc/kubernetes/pki/ca.pem \
-ca-key=/etc/kubernetes/pki/ca-key.pem \
-config=ca-config.json \
-profile=kubernetes \
admin-csr.json | cfssljson -bare /etc/kubernetes/pki/admin
```

kubectl 配置 admin
```bash
# 设置集群参数
kubectl config set-cluster kubernetes \
--certificate-authority=/etc/kubernetes/pki/ca.pem \
--embed-certs=true \
--server=https://192.168.217.132:8843 \
--kubeconfig=/etc/kubernetes/admin.kubeconfig

# 设置客户端参数
kubectl config set-credentials kubernetes-admin \
--client-certificate=/etc/kubernetes/pki/admin.pem \
--client-key=/etc/kubernetes/pki/admin-key.pem \
--embed-certs=true \
--kubeconfig=/etc/kubernetes/admin.kubeconfig

# 设置上下文参数
kubectl config set-context kubernetes-admin@kubernetes \
--cluster=kubernetes \
--user=kubernetes-admin \
--kubeconfig=/etc/kubernetes/admin.kubeconfig

# 设置默认上下文
kubectl config use-context kubernetes-admin@kubernetes \
--kubeconfig=/etc/kubernetes/admin.kubeconfig

```

创建 service account key
```bash
openssl genrsa --out /etc/kubernetes/pki/sa.key 2048
openssl rsa -in /etc/kubernetes/pki/sa.key -pubout -out /etc/kubernetes/pki/sa.pub
```

将证书复制到各个 master 节点
```bash
for node_ip in ${master_nodes}; do
	for file in $(ls /etc/kubernetes/pki | grep -v etcd); do
		rsync /etc/kubernetes/pki/${file} ${node_ip}:/etc/kubernetes/pki/${file}
	done
	
	for file in admin.kubeconfig controller-manager.kubeconfig scheduler.kubeconfig; do
		rsync /etc/kubernetes/${file} ${node_ip}:/etc/kubernetes/${file}
	done
done
```

配置etcd集群
```bash
cat >> etcd.config.yml.template <<EOF
name: '~node_name~'
data-dir: /var/lib/etcd
wal-dir: /var/lib/etcd/wal
snapshot-count: 5000
heartbeat-interval: 100
election-timeout: 1000
quota-backend-bytes: 0
listen-peer-urls: 'https://~node_ip~:2380'
listen-client-urls: 'https://~node_ip~:2379,http://127.0.0.1:2379'
max-snapshots: 3
max-wals: 5
cors:
initial-advertise-peer-urls: 'https://~node_ip~:2380'
advertise-client-urls: 'https://~node_ip~:2379'
discovery:
discovery-fallback: 'proxy'
discovery-srv:
initial-cluster: '~ETCDCLUSTER~'
initial-cluster-token: 'etcd-k8s-cluster'
initial-cluster-state: 'new'
enable-v2: true
enable-pprof: true
proxy: 'off'
proxy-failure-wait: 5000
proxy-refresh-interval: 30000
proxy-dial-timeout: 1000
proxy-write-timeout: 5000
proxy-read-timeout: 0
client-transport-security:
  cert-file: '/etc/etcd/ssl/etcd.pem'
  key-file: '/etc/etcd/ssl/etcd-key.pem'
  client-cert-auth: true
  trusted-ca-file: '/etc/etcd/ssl/etcd-ca.pem'
  auto-tls: true
peer-transport-security:
  cert-file: "/etc/etcd/ssl/etcd.pem"
  key-file: "/etc/etcd/ssl/etcd-key.pem"
  client-cert-auth: true
  trusted-ca-file: "/etc/etcd/ssl/etcd-ca.pem"
  auto-tls: true
log-package-levels:
log-outputs: [default]
EOF


NODENAMES=(master node-1 node-2)
NODEIPS=(192.168.217.132 192.168.217.131 192.168.217.133)
ETCDCLUSTER="master-1=https://192.168.217.132:2380,node-1=https://192.168.217.131:2380,node-2=https://192.168.217.133:2380"
for ((i=0;i<3;i++)); do
	sed -e "s#~node_name~#${NODENAMES[i]}#" -e "s#~node_ip~#${NODEIPS[i]}#" -e "s#~ETCDCLUSTER~#${ETCDCLUSTER}#" etcd.config.yml.template > etcd-${NODEIPS[i]}.config.yml
	rsync etcd-${NODEIPS[i]}.config.yml ${NODENAMES[i]}:/etc/etcd/etcd.config.yml
done
rm etcd-*.config.yml
```

创建 service
```bash
vi /usr/lib/systemd/system/etcd.service
[Unit]
Description=Etcd service
Documentation=https://coreos.com/etcd/docs/latest/
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd --config-file=/etc/etcd/etcd.config.yml
Restart=on-failure
RestartSec=10
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
Alias=etcd3.service

# 启动 etcd 服务
systemctl daemon-reload
systemctl enable --now etcd

# 检查 etcd 服务状态
ETCDCTL_API=3 etcdctl --endpoints="192.168.217.132:2379,192.168.217.131:2379,192.168.217.133:2379" --cacert=/etc/etcd/ssl/etcd-ca.pem --cert=/etc/etcd/ssl/etcd.pem --key=/etc/etcd/ssl/etcd-key.pem endpoint status --write-out=table
---
tags: ["k8s"]
---
# k8s 二进制安装
## 环境

-   centos 8.3 4.18.0-240.22.1.el8_3.x86_64
-   master3：192.168.217.134 192.168.217.135 192.168.217.136
-   node3：192.168.217.137 192.168.217.138 192.168.217.139
-   kubernets 1.21.0
-   etcd v3.4.13

## 基础配置

1.  修改镜像源

```bash
sed -e 's|^mirrorlist=|#mirrorlist=|g' \
-e 's|^#baseurl=http://mirror.centos.org/$contentdir|baseurl=https://mirrors.ustc.edu.cn/centos|g' \
         -i.bak \
		 /etc/yum.repos.d/CentOS-Linux-AppStream.repo \
         /etc/yum.repos.d/CentOS-Linux-BaseOS.repo \
         /etc/yum.repos.d/CentOS-Linux-Extras.repo \
         /etc/yum.repos.d/CentOS-Linux-PowerTools.repo \
         /etc/yum.repos.d/CentOS-Linux-Plus.repo
```

2. 修改每台机器的 /etc/hosts ，添加主机名和 IP 的对应关系

```bash
cat >> /etc/hosts <<EOF
192.168.217.134 master-1
192.168.217.135 master-2
192.168.217.136 master-3
192.168.217.137 node-1
192.168.217.138 node-2
192.168.217.139 node-3
EOF
```

3. master生成私钥将公钥分发到node

```bash
ssh-keygen -t rsa
```

```bash
export nodes=(master-2 master-3 node-1 node-2 node-3)
for node_ip in ${nodes[@]}; do
	ssh-copy-id -i .ssh/id_rsa.pub ${node_ip}
done
```

4. 安装依赖包

```bash
dnf install -y epel-release
dnf install -y conntrack ipvsadm ipset jq sysstat curl iptables libseccomp git vim wget nc
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
modprobe overlay
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
sed -i'*.bak' 's#^pool .*\.org#pool time\.pool\.aliyun\.com#' /etc/chrony.conf
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

13. 安装容器，二选一
	1.  安装docker
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

修改镜像源

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210707000114.png)

```bash
	dnf install -y dnf-utils device-mapper-persistent-data lvm2
	yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
	yum update -y && yum install -y containerd.io
	containerd config default > /etc/containerd/config.toml
	sed -i'*.bak' 's#sandbox_image = \".*\"#sandbox_image = \"registry.aliyuncs.com/google_containers/pause:3.2\"#' /etc/containerd/config.toml
	systemctl daemon-reload && systemctl enable --now containerd.service

	# 安装命令行工具
	export version=0.10.0
	wget https://github.com/containerd/nerdctl/releases/download/v${version}/nerdctl-${version}-linux-amd64.tar.gz
	tar xvf nerdctl-${version}-linux-amd64.tar.gz -C /usr/local/bin nerdctl
	nerdctl images
```
	
14. 下载二进制包

```bash
wget https://dl.k8s.io/v1.21.0/kubernetes-server-linux-amd64.tar.gz
wget https://github.com/etcd-io/etcd/releases/download/v3.4.13/etcd-v3.4.13-linux-amd64.tar.gz

export nodes=(master-2 master-3 node-1 node-2 node-3)
for node_ip in ${nodes[@]}; do
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
```

预定义一些环境变量
```bash
# 生成 EncryptionConfig 所需的加密 key
ENCRYPTION_KEY=$(head -c 32 /dev/urandom | base64)

# worker ip
export NODES_IP=(192.168.217.137 192.168.217.138 192.168.217.139)

# master ip
export MASTERS_IP=(192.168.217.134 192.168.217.135 192.168.217.136)

# 所有ip
export ALL_IP=(192.168.217.134 192.168.217.135 192.168.217.136 192.168.217.137 192.168.217.138 192.168.217.139)

# etcd 集群服务地址列表
export ETCD_ENDPOINTS="https://192.168.217.134:2379,https://192.168.217.135:2379,https://192.168.217.136:2379"  

```

创建CA根证书和秘钥
```bash
 cat > ca-config.json <<EOF
 {
 	"signing": {
		"default": {
			"expiry": "876000h"
		},
		"profiles": {
			"kubernetes": {
				"usages": [
					"signing",
					"key encipherment",
					"server auth",
					"client auth"
				],
				"expiry": "876000h"
			}
		}
	}
 }
 EOF
 ```
 - signing: 表示该证书可用于签名其它证书
 - server auth：表示client可以用该证书对 server 提供的证书进行验证
 - client auth：表示server可以用该证书对 client 提供的证书进行验证
 - expiry：设置证书有效期
 
 
 ```bash
 cat > ca-csr.json <<EOF
 {
 	"CN": "kubernetes",
	"key": {
		"algo": "rsa",
		"size": 2048
	},
	"names": [
		{
			"C": "CN",
			"ST": "Beijing",
			"L": "Beijing",
			"O": "kubernetes",
			"OU": "Kubernetes-manual"
		}
	],
	"ca": {
		"expiry": "876000h"
	}
 }
 EOF
```
- CN（Common Name）：kube apiserver 从证书中提取该字段作为请求的用户名，浏览器使用该字段验证网站是否合法
- O（Organization）：kube apiserver 从证书中提取该字段作为请求用户所属组
- kube apiserver 将提取 user、group 作为 RBAC 授权的用户标识

生成 kubernetes 根证书
 ```bash
 cfssl gencert -initca ca-csr.json | cfssljson -bare /etc/kubernetes/pki/ca
 ```
 
 生成 admin 证书
```bash
cat > admin-csr.json <<EOF
{
  "CN": "admin",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "Beijing",
      "L": "Beijing",
      "O": "system:masters",
      "OU": "Kubernetes-manual"
    }
  ]
}
EOF

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
--server=https://192.168.217.236:8443 \
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


 生成etcd CA 证书 和 CA 证书的 key
 ```bash
 mkdir -p ssl && cd ssl
 cat > /etcd-ca-csr.json <<EOF
 {
 	"CN": "etcd",
	"key": {
		"algo": "rsa",
		"size": 2048
	},
	"names": [
		{
			"C": "CN",
			"ST": "Beijing",
			"L": "Beijing",
			"O": "etcd",
			"OU": "Etcd Security"
		}
	],
	"ca": {
		"expiry": "876000h"
	}
 }
 EOF
 
 cat > etcd-csr.json <<EOF
 {
 	"CN": "etcd",
	"key": {
		"algo": "rsa",
		"size": 2048
	},
	"names": [
		{
			"C": "CN",
			"ST": "Beijing",
			"L": "Beijing",
			"O": "etcd",
			"OU": "Etcd Security"
		}
	],
	"hosts": [
		"127.0.0.1",
		"192.168.217.134",
		"192.168.217.135",
		"192.168.217.136",
		"master-1",
		"master-2",
		"master-3"
	]
 }
 EOF
 
 
 # 生成 etcd 根证书
cfssl gencert -initca etcd-ca-csr.json | cfssljson -bare /etc/etcd/ssl/etcd-ca

# 生成 etcd 客户端证书
cfssl gencert -ca=/etc/etcd/ssl/etcd-ca.pem -ca-key=/etc/etcd/ssl/etcd-ca-key.pem -config=ca-config.json -profile=kubernetes etcd-csr.json | cfssljson -bare /etc/etcd/ssl/etcd

# 将etcd证书复制到其他 master 节点
export master_ips=(master-2 master-3)
for node_ip in ${master_ips[@]}; do
	ssh ${node_ip} "mkdir -p /etc/etcd/ssl";
	for file in etcd-ca-key.pem etcd-ca.pem etcd-key.pem etcd.pem; do
		rsync /etc/etcd/ssl/${file} ${node_ip}:/etc/etcd/ssl/${file};
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

NODENAMES=(master-1 master-2 master-3)
NODEIPS=(192.168.217.134 192.168.217.135 192.168.217.136)
ETCDCLUSTER="master-1=https://192.168.217.134:2380,master-2=https://192.168.217.135:2380,master-3=https://192.168.217.136:2380"
for ((i=0;i<3;i++)); do
	sed -e "s#~node_name~#${NODENAMES[i]}#" -e "s#~node_ip~#${NODEIPS[i]}#" -e "s#~ETCDCLUSTER~#${ETCDCLUSTER}#" etcd.config.yml.template > etcd-${NODEIPS[i]}.config.yml
	rsync etcd-${NODEIPS[i]}.config.yml ${NODENAMES[i]}:/etc/etcd/etcd.config.yml
done
rm etcd-*.config.yml
```

在master节点创建 etcd service
```bash
cat > /usr/lib/systemd/system/etcd.service <<EOF
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
EOF

export nodes=(master-2 master-3)
for node in ${nodes[@]}; do
	rsync /usr/lib/systemd/system/etcd.service ${node}:/usr/lib/systemd/system/etcd.service
done

# 启动 etcd 服务
systemctl daemon-reload && systemctl enable --now etcd

# 检查 etcd 服务状态
export master_ips="192.168.217.134:2379,192.168.217.135:2379,192.168.217.136:2379"
ETCDCTL_API=3 etcdctl --endpoints="${master_ips}" --cacert=/etc/etcd/ssl/etcd-ca.pem --cert=/etc/etcd/ssl/etcd.pem --key=/etc/etcd/ssl/etcd-key.pem endpoint status --write-out=table
```
 
 keepalived haproxy 高可用配置
 ```bash
 dnf install -y keepalived haproxy
 
 cat > /etc/haproxy/haproxy.cfg <<EOF
 global
   maxconn 2000
   ulimit-n 16384
   log 127.0.0.1 local0 err
   stats timeout 30s
   
 defaults
   log global
   mode http
   option httplog
   timeout connect 5000
   timeout client 50000
   timeout server 50000
   timeout http-request 15s
   timeout http-keep-alive 15s
   
 frontend k8s-master
   bind 0.0.0.0:8443
   bind 127.0.0.1:8443
   mode tcp
   option tcplog
   tcp-request inspect-delay 5s
   default_backend k8s-master
   
 backend k8s-master
   mode tcp
   option tcplog
   option tcp-check
   balance roundrobin
   default-server inter 10s downinter 5s rise 2 fall 2 slowstart 60s maxconn 250 maxqueue 256 weight 100
   server k8s-master-1 192.168.217.134:6443 check
   server k8s-master-2 192.168.217.135:6443 check
   server k8s-master-3 192.168.217.136:6443 check
 EOF
 
 systemctl daemon-reload && systemctl enable --now haproxy.service
 
 # keepalive 需要配置一主多从
 # 根据连接的网卡接口配置 interface ens33
 # virtual_ipaddress 需要根据当前的网段分配一个地址给 keepalive 尽量分配大一些避免冲突
 # master
 cat > /etc/keepalived/keepalived.conf <<EOF
 global_defs {
 	router_id LVS_DEVEL
 }
 
 vrrp_script chk_apiserver {
 	script "/etc/keepalived/check_apiserver.sh"
	interval 5
	weight -5
	fall 2
	rise 1
 }
 
 vrrp_instance VI_1 {
 	state MASTER
	interface ens33
	mcast_src_ip 192.168.217.134
	virtual_router_id 51
	priority 101
	nopreempt
	advert_int 2
	authentication {
		auth_type PASS
		auth_pass K8SHA_KA_AUTH
	}
	
	virtual_ipaddress {
		192.168.217.236
	}
	track_script {
		chk_apiserver
	}
 }
 EOF
 
 # backup master-2
 cat > /etc/keepalived/keepalived.conf <<EOF
 global_defs {
 	router_id LVS_DEVEL
 }
 
 vrrp_script chk_apiserver {
 	script "/etc/keepalived/check_apiserver.sh"
	interval 5
	weight -5
	fall 2
	rise 1
 }
 
 vrrp_instance VI_1 {
 	state BACKUP
	interface ens33
	mcast_src_ip 192.168.217.135
	virtual_router_id 51
	priority 101
	nopreempt
	advert_int 2
	authentication {
		auth_type PASS
		auth_pass K8SHA_KA_AUTH
	}
	
	virtual_ipaddress {
		192.168.217.236
	}
	track_script {
		chk_apiserver
	}
 }
 EOF
 
 # backup master-3
 cat > /etc/keepalived/keepalived.conf <<EOF
 global_defs {
 	router_id LVS_DEVEL
 }
 
 vrrp_script chk_apiserver {
 	script "/etc/keepalived/check_apiserver.sh"
	interval 5
	weight -5
	fall 2
	rise 1
 }
 
 vrrp_instance VI_1 {
 	state BACKUP
	interface ens33
	mcast_src_ip 192.168.217.136
	virtual_router_id 51
	priority 101
	nopreempt
	advert_int 2
	authentication {
		auth_type PASS
		auth_pass K8SHA_KA_AUTH
	}
	
	virtual_ipaddress {
		192.168.217.236
	}
	track_script {
		chk_apiserver
	}
 }
 EOF
 
 systemctl daemon-reload && systemctl enable --now keepalive.service
 ```
 
 
 生成 apiserver 证书
 ```bash
 cat > apiserver-csr.json <<EOF
 {
  "CN": "kube-apiserver",
  "hosts": [
  	"127.0.0.1",
	"192.168.217.134",
	"192.168.217.135",
	"192.168.217.136",
	"192.168.217.137",
	"192.168.217.138",
	"192.168.217.139",
	"192.168.217.236",
	"10.96.0.1",
	"kubernetes",
	"kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "Beijing",
      "L": "Beijing",
      "O": "Kubernetes",
      "OU": "Kubernetes-manual"
    }
  ]
}
EOF

 cfssl gencert -ca=/etc/kubernetes/pki/ca.pem -ca-key=/etc/kubernetes/pki/ca-key.pem -config=ca-config.json -profile=kubernetes apiserver-csr.json | cfssljson -bare /etc/kubernetes/pki/apiserver
 ```
 
生成 apiserver 代理证书
```bash
cat > front-proxy-ca-csr.json <<EOF
{
  "CN": "kubernetes",
  "key": {
     "algo": "rsa",
     "size": 2048
  }
}
EOF

cat > front-proxy-client-csr.json <<EOF
{
	"CN": "front-proxy-client",
	"key": {
		"algo": "rsa",
		"size": 2048
	}
}
EOF

cfssl gencert -initca front-proxy-ca-csr.json | cfssljson -bare /etc/kubernetes/pki/front-proxy-ca

cfssl gencert -ca=/etc/kubernetes/pki/front-proxy-ca.pem -ca-key=/etc/kubernetes/pki/front-proxy-ca-key.pem -config=ca-config.json -profile=kubernetes front-proxy-client-csr.json | cfssljson -bare /etc/kubernetes/pki/front-proxy-client
```

生成controller manager 证书
```bash
cat > controller-manager-csr.json <<EOF
{
  "CN": "system:kube-controller-manager",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "hosts": [
  	"127.0.0.1",
	"192.168.217.134",
	"192.168.217.135",
	"192.168.217.136"
  ],
  "names": [
    {
      "C": "CN",
      "ST": "Beijing",
      "L": "Beijing",
      "O": "system:kube-controller-manager",
      "OU": "Kubernetes-manual"
    }
  ]
}
EOF

cfssl gencert -ca=/etc/kubernetes/pki/ca.pem -ca-key=/etc/kubernetes/pki/ca-key.pem -config=ca-config.json -profile=kubernetes controller-manager-csr.json | cfssljson -bare /etc/kubernetes/pki/controller-manager
```

kubectl 配置 controller-manger
```bash
# 设置集群参数
kubectl config set-cluster kubernetes \
--certificate-authority=/etc/kubernetes/pki/ca.pem \
--embed-certs=true \
--server=https://192.168.217.236:8443 \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置客户端参数
kubectl config set-credentials system:kube-controller-manager \
--client-certificate=/etc/kubernetes/pki/controller-manager.pem \
--client-key=/etc/kubernetes/pki/controller-manager-key.pem \
--embed-certs=true \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置上下文参数
kubectl config set-context system:kube-controller-manager \
--cluster=kubernetes \
--user=system:kube-controller-manager \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig

# 设置默认上下文
kubectl config use-context system:kube-controller-manager \
--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig
```

生成 scheduler 证书
```bash
cat > scheduler-csr.json <<EOF
{
  "CN": "system:kube-scheduler",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "hosts": [
  	"127.0.0.1",
	"192.168.217.134",
	"192.168.217.135",
	"192.168.217.136"
  ]
  "names": [
    {
      "C": "CN",
      "ST": "Beijing",
      "L": "Beijing",
      "O": "system:kube-scheduler",
      "OU": "Kubernetes-manual"
    }
  ]
}
EOF

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
--server=https://192.168.217.236:8443 \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

# 设置客户端参数
kubectl config set-credentials system:kube-scheduler \
--client-certificate=/etc/kubernetes/pki/scheduler.pem \
--client-key=/etc/kubernetes/pki/scheduler-key.pem \
--embed-certs=true \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig


# 设置上下文参数
kubectl config set-context system:kube-scheduler \
--cluster=kubernetes \
--user=system:kube-scheduler \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

# 设置默认上下文
kubectl config use-context system:kube-scheduler \
--kubeconfig=/etc/kubernetes/scheduler.kubeconfig

```

创建 service account key
```bash
openssl genrsa --out /etc/kubernetes/pki/sa.key 2048
openssl rsa -in /etc/kubernetes/pki/sa.key -pubout -out /etc/kubernetes/pki/sa.pub
```

将证书复制到各个 master 节点
```bash
source ${HOME}/environment.sh
for node_ip in ${MASTERS_IP[@]}; do
	for file in $(ls /etc/kubernetes/pki); do
		rsync -v --progress /etc/kubernetes/pki/${file} root@${node_ip}:/etc/kubernetes/pki/${file}
	done
	
	for file in admin.kubeconfig controller-manager.kubeconfig scheduler.kubeconfig; do
		rsync -v --progress /etc/kubernetes/${file} root@${node_ip}:/etc/kubernetes/${file}
	done
done
```


16. kubernetes 组件配置
	1. 给所有节点创建目录
		```bash
		mkdir -p /etc/kubernetes/manifests /etc/systemd/system/kubelet.service.d /var/lib/kubelet /var/log/kubernetes
		```
	2. 给 master 节点创建 kube-apiserver service
		```bash
		export etcd_nodes="https://192.168.217.134:2379,https://192.168.217.135:2379,https://192.168.217.136:2379"
		cat > kube-apiserver.service.template <<EOF
		[Unit]
		Description=kubernetes API Server
		Documentation=https://github.com/kubernetes/kubernetes
		After=network.target
		
		[Service]
		ExecStart=/usr/local/bin/kube-apiserver \\
			--v=2 \\
			--logtostderr=false \\
			--log-dir=/var/log/kubernetes \\
			--allow-privileged=true \\
			--bind-address=0.0.0.0 \\
			--advertise-address=##NODE_IP## \\
			--service-cluster-ip-range=10.96.0.0/12 \\
			--service-node-port-range=30000-32767 \\
			--etcd-servers=${etcd_nodes} \\
			--etcd-cafile=/etc/etcd/ssl/etcd-ca.pem \\
			--etcd-certfile=/etc/etcd/ssl/etcd.pem \\
			--etcd-keyfile=/etc/etcd/ssl/etcd-key.pem \\
			--client-ca-file=/etc/kubernetes/pki/ca.pem \\
			--tls-cert-file=/etc/kubernetes/pki/apiserver.pem \\
			--tls-private-key-file=/etc/kubernetes/pki/apiserver-key.pem \\
			--kubelet-client-certificate=/etc/kubernetes/pki/apiserver.pem \\
			--kubelet-client-key=/etc/kubernetes/pki/apiserver-key.pem \\
			--service-account-key-file=/etc/kubernetes/pki/sa.pub \\
			--service-account-signing-key-file=/etc/kubernetes/pki/sa.key \\
			--service-account-issuer=https://kubernetes.default.svc.cluster.local \\
			--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \\
			--enable-admission-	plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota \\
			--authorization-mode=Node,RBAC \\
			--enable-bootstrap-token-auth=true \\
			--proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.pem \\
			--proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client-key.pem \\
			--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.pem \\
			--requestheader-allowed-names=aggregator \\
			--requestheader-group-headers=X-Remote-Group \\
			--requestheader-extra-headers-prefix=X-Remote-Extra- \\
			--requestheader-username-headers=X-Remote-User
		
		Restart=on-failure
		RestartSec=10s
		LimitNOFILE=65535
		
		[Install]
		WantedBy=multi-user.target
		EOF
		
		export master_ips=(192.168.217.134 192.168.217.135 192.168.217.136)
		export master_hosts=(master-1 master-2 master-3)
		for ((i=0; i<3; i++)); do
			sed -e "s/##NODE_IP##/${master_ips[i]}/" kube-apiserver.service.template > kube-apiserver-${master_ips[i]}.service
			rsync kube-apiserver-${master_ips[i]}.service ${master_hosts[i]}:/usr/lib/systemd/system/kube-apiserver.service
		done
		
		rm -f kube-apiserver-*.service
		
		# 启动 kube-apiserver 服务
		systemctl daemon-reload && systemctl enable --now kube-apiserver.service
		```
		3. 给 master 节点创建 kube-controller-manager.service
		```bash
		cat > kube-controller-manager.service <<EOF
		[Unit]
		Description=Kubernetes Controller Manager
		Documentation=https://github.com/kubernetes/kubernetes
		After=network.target
		
		[Service]
		ExecStart=/usr/local/bin/kube-controller-manager \\
		--v=2 \\
		--logtostderr=false \\
		--log-dir=/var/log/kubernetes \\
		--bind-address=127.0.0.1 \\
		--root-ca-file=/etc/kubernetes/pki/ca.pem \\
		--cluster-signing-cert-file=/etc/kubernetes/pki/ca.pem \\
		--cluster-signing-key-file=/etc/kubernetes/pki/ca-key.pem \\
		--service-account-private-key-file=/etc/kubernetes/pki/sa.key \\
		--kubeconfig=/etc/kubernetes/controller-manager.kubeconfig \\
		--use-service-account-credentials=true \\
		--node-monitor-grace-period=40s \\
		--node-monitor-period=5s \\
		--pod-eviction-timeout=2m0s \\
		--controllers=*,bootstrapsigner,tokencleaner \\
		--allocate-node-cidrs=true \\
		--cluster-cidr=172.16.0.0/12 \\
		--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.pem \\
		--node-cidr-mask-size=24
		
		Restart=on-failure
		RestartSec=10s
		
		[Install]
		WantedBy=multi-user.target
		EOF
		
		export master_hosts=(master-1 master-2 master-3)
		for ((i=0; i<3; i++)); do
			rsync kube-controller-manager.service ${master_hosts[i]}:/usr/lib/systemd/system/kube-controller-manager.service
			ssh ${master_hosts[i]} "systemctl daemon-reload && systemctl enable --now kube-controller-manager.service"
		done
		
		```
		4. 配置 scheduler 服务
		```bash
		cat > kube-scheduler.service <<EOF
		[Unit]
		Description=Kubernetes Scheduler
		Documentation=https://github.com/kubernetes/kubernetes
		After=network.target
		
		[Service]
		ExecStart=/usr/local/bin/kube-scheduler \\
		--v=2 \\
		--logtostderr=false \\
		--address=127.0.0.1 \\
		--leader-elect=true \\
		--kubeconfig=/etc/kubernetes/scheduler.kubeconfig \\
		--log-dir=/var/log/kubernetes
		
		Restart=on-failure
		RestartSec=10s
		
		[Install]
		WantedBy=multi-user.target
		EOF
		
		export master_hosts=(master-1 master-2 master-3)
		for ((i=0; i<3; i++)); do
			rsync kube-scheduler.service ${master_hosts[i]}:/usr/lib/systemd/system/kube-scheduler.service
			ssh ${master_hosts[i]} "systemctl daemon-reload && systemctl enable --now kube-scheduler.service"
		done
		```
		5. tls bootstrap 自动颁发证书
		```bash
		cat > bootstrap.secret.yaml <<EOF
		apiVersion: v1
		kind: Secret
		metadata:
		  name: bootstrap-token-c8ad9c
		  namespace: kube-system
		type: bootstrap.kubernetes.io/token
		stringData:
		  token-id: c8ad9c
		  token-secret: 2e4d610cf3e7426e
		  usage-bootstrap-authentication: "true"
		  usage-bootstrap-signing: "true"
		  auth-extra-groups:  system:bootstrappers:default-node-token,system:bootstrappers:worker,system:bootstrappers:ingress

		---
		apiVersion: rbac.authorization.k8s.io/v1
		kind: ClusterRoleBinding
		metadata:
		  name: kubelet-bootstrap
		roleRef:
		  apiGroup: rbac.authorization.k8s.io
		  kind: ClusterRole
		  name: system:node-bootstrapper
		subjects:
		- apiGroup: rbac.authorization.k8s.io
		  kind: Group
		  name: system:bootstrappers:default-node-token
		---
		apiVersion: rbac.authorization.k8s.io/v1
		kind: ClusterRoleBinding
		metadata:
		  name: node-autoapprove-bootstrap
		roleRef:
		  apiGroup: rbac.authorization.k8s.io
		  kind: ClusterRole
		  name: system:certificates.k8s.io:certificatesigningrequests:nodeclient
		subjects:
		- apiGroup: rbac.authorization.k8s.io
		  kind: Group
		  name: system:bootstrappers:default-node-token
		---
		apiVersion: rbac.authorization.k8s.io/v1
		kind: ClusterRoleBinding
		metadata:
		  name: node-autoapprove-certificate-rotation
		roleRef:
		  apiGroup: rbac.authorization.k8s.io
		  kind: ClusterRole
		  name: system:certificates.k8s.io:certificatesigningrequests:selfnodeclient
		subjects:
		- apiGroup: rbac.authorization.k8s.io
		  kind: Group
		  name: system:nodes
		---
		apiVersion: rbac.authorization.k8s.io/v1
		kind: ClusterRole
		metadata:
		  annotations:
			rbac.authorization.kubernetes.io/autoupdate: "true"
		  labels:
			kubernetes.io/bootstrapping: rbac-defaults
		  name: system:kube-apiserver-to-kubelet
		rules:
		  - apiGroups:
			  - ""
			resources:
			  - nodes/proxy
			  - nodes/stats
			  - nodes/log
			  - nodes/spec
			  - nodes/metrics
			verbs:
			  - "*"
		---
		apiVersion: rbac.authorization.k8s.io/v1
		kind: ClusterRoleBinding
		metadata:
		  name: system:kube-apiserver
		  namespace: ""
		roleRef:
		  apiGroup: rbac.authorization.k8s.io
		  kind: ClusterRole
		  name: system:kube-apiserver-to-kubelet
		subjects:
		  - apiGroup: rbac.authorization.k8s.io
			kind: User
			name: kube-apiserver
		EOF
		
		kubectl config set-cluster kubernetes --certificate-	authority=/etc/kubernetes/pki/ca.pem --embed-certs=true --server=https://192.168.217.236:8443 --kubeconfig=/etc/kubernetes/bootstrap-kubelet.kubeconfig
		
		kubectl config set-credentials tls-bootstrap-token-user --token=c8ad9c.2e4d610cf3e7426e --kubeconfig=/etc/kubernetes/bootstrap-kubelet.kubeconfig
		
		kubectl config set-context tls-bootstrap-token-user@kubernetes --cluster=kubernetes --user=tls-bootstrap-token-user --kubeconfig=/etc/kubernetes/bootstrap-kubelet.kubeconfig
		
		kubectl config use-context tls-bootstrap-token-user@kubernetes --kubeconfig=/etc/kubernetes/bootstrap-kubelet.kubeconfig
		
		mkdir -p /root/.kube && cp /etc/kubernetes/admin.kubeconfig /root/.kube/config
		
		kubectl create -f bootstrap.secret.yaml
		```
		5. 复制证书到 node 节点
		```bash
		export nodes=(master-2 master-3 node-1 node-2 node-3)
		for node in ${nodes[@]}; do
		    ssh ${node} "mkdir -p /etc/kubernetes/pki /etc/etcd/ssl"
		    for FILE in etcd-ca.pem etcd.pem etcd-key.pem; do
		        rsync -v /etc/etcd/ssl/${FILE} ${node}:/etc/etcd/ssl
		    done
		    for File in pki/ca.pem pki/ca-key.pem pki/front-proxy-ca.pem bootstrap-kubelet.kubeconfig; do
		        rsync -v --progress /etc/kubernetes/${File} ${node}:/etc/kubernetes/${File}
		    done
		done
		```
		6.  kubelet service 配置
		```bash
		cat > kubelet.service <<EOF
		[Unit]
		Description=Kubernetes Kubelet
		Documentation=https://github.com/kubernetes/kubernetes
		After=containerd.service
		
		[Service]
		ExecStart=/usr/local/bin/kubelet
		
		Restart=on-failure
		StartLimitInterval=0
		RestartSec=10
		
		[Install]
		WantedBy=multi-user.target
		EOF
		
		cat > 10-kubelet.conf <<EOF
		[Service]
		Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.kubeconfig --kubeconfig=/etc/kubernetes/kubelet.kubeconfig --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock"
		Environment="KUBELET_SYSTEM_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
		Environment="KUBELET_CONFIG_ARGS=--config=/etc/kubernetes/kubelet-conf.yaml --pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.2"
		Environment="KUBELET_EXTRA_ARGS=--node-labels=node.kubernetes.io/node=''"
		ExecStart=
		ExecStart=/usr/local/bin/kubelet \$KUBELET_KUBECONFIG_ARGS \$KUBELET_CONFIG_ARGS \$KUBELET_SYSTEM_ARGS \$KUBELET_EXTRA_ARGS
		EOF
		
		
		# 创建kubelet配置文件
		cat > kubelet-conf.yaml <<EOF
		apiVersion: "kubelet.config.k8s.io/v1beta1"
		kind: "KubeletConfiguration"
		enableServer: true
		authentication:
		  x509:
			clientCAFile: "/etc/kubernetes/pki/ca.pem"
		  webhook:
			enabled: true
			cacheTTL: "2m0s"
		  anonymous:
			enabled: true
		authorization:
		  mode: "Webhook"
		  webhook:
			cacheAuthorizedTTL: "5m0s"
			cacheUnauthorizedTTL: "30s"
		clusterDomain: "cluster.local"
		clusterDNS:
		- "10.96.0.10"
		cgroupDriver: "systemd"
		staticPodPath: "/etc/kubernetes/manifests"
		readOnlyPort: 10255
		rotateCertificates: true
		EOF
		
		export nodes=(master-2 master-3 node-1 node-2 node-3)
		for node in ${nodes[@]}; do
			ssh ${node} "mkdir -p /var/lib/kubelet /var/log/kubernetes /etc/systemd/system/kubelet.service.d /etc/kubernetes/manifests"
			rsync -v kubelet.service ${node}:/usr/lib/systemd/system/kubelet.service
			rsync -v 10-kubelet.conf ${node}:/etc/systemd/system/kubelet.service.d/10-kubelet.conf
			rsync -v kubelet-conf.yaml ${node}:/etc/kubernetes/kubelet-conf.yaml
			ssh ${node} "systemctl daemon-reload && systemctl enable --now kubelet.service"
		done
		```
		
		kube-proxy 配置
		```bash
		# 第一种方式
		cat > ssl/kube-proxy-csr.json <<EOF
		{
		  "CN": "system:kube-proxy",
		  "key": {
			"algo": "rsa",
			"size": 2048
		  },
		  "names": [
			{
			  "C": "CN",
			  "ST": "BeiJing",
			  "L": "BeiJing",
			  "O": "system:kube-proxy",
			  "OU": "Kubernetes-manual"
			}
		  ]
		}
		EOF
		
		cfssl gencert -ca=/etc/kubernetes/pki/ca.pem \
		  -ca-key=/etc/kubernetes/pki/ca-key.pem \
		  -config=ssl/ca-config.json \
		  -profile=kubernetes ssl/kube-proxy-csr.json | cfssljson -bare /etc/kubernetes/pki/kube-proxy
		  
		export PKI_DIR=/etc/kubernetes/pki
		export K8S_DIR=/etc/kubernetes
		kubectl config set-cluster kubernetes --certificate-authority=/etc/kubernetes/pki/ca.pem --embed-certs=true --server=https://192.168.217.236:8443 --kubeconfig=$K8S_DIR/kube-proxy.kubeconfig
		kubectl config set-credentials kube-proxy \
		--client-certificate=$PKI_DIR/kube-proxy.pem \
		--client-key=$PKI_DIR/kube-proxy-key.pem \
		--embed-certs=true \
		--kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig
		kubectl config set-context default --cluster=kubernetes --user=kube-proxy --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig
		kubectl config use-context default --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig
		
		
		-----
		
		# 第二种方式
		kubectl -n kube-system create serviceaccount kube-proxy
		kubectl create clusterrolebinding system:kube-proxy --clusterrole system:node-proxier --serviceaccount kube-system:kube-proxy
		export SECRET=$(kubectl -n kube-system get sa/kube-proxy --output=jsonpath='{.secrets[0].name}')
		export JWT_TOKEN=$(kubectl -n kube-system get secret/$SECRET --output=jsonpath='{.data.token}' | base64 -d)
		kubectl config set-cluster kubernetes \
		  --certificate-authority=${PKI_DIR}/ca.pem \
		  --embed-certs=true \
		  --server=https://192.168.217.236:8443 \
		  --kubeconfig=${K8S_DIR}/kube-proxy.kubeconfig
		  
		kubectl config set-credentials kubernetes --token=${JWT_TOKEN} \
		--kubeconfig=${K8S_DIR}/kube-proxy.kubeconfig
		
		kubectl config set-context kubernetes --cluster=kubernetes \
		--user=kubernetes \
		--kubeconfig=${K8S_DIR}/kube-proxy.kubeconfig
		
		kubectl config use-context kubernetes --kubeconfig=${K8S_DIR}/kube-proxy.kubeconfig
		
		```
		
		```bash
		mkdir kube-proxy
		cat > kube-proxy/kube-proxy.conf <<EOF
		apiVersion: kubeproxy.config.k8s.io/v1alpha1
		bindAddress: 0.0.0.0
		clientConnection:
		  acceptContentTypes: ""
		  burst: 10
		  contentType: application/vnd.kubernetes.protobuf
		  kubeconfig: /etc/kubernetes/kube-proxy.kubeconfig
		  qps: 5
		clusterCIDR: 172.16.0.0/12
		configSyncPeriod: 15m0s
		conntrack:
		  max: null
		  maxPerCore: 32768
		  min: 131072
		  tcpCloseWaitTimeout: 1h0m0s
		  tcpEstablishedTimeout: 24h0m0s
		enableProfiling: false
		healthzBindAddress: 0.0.0.0:10256
		hostnameOverride: ""
		iptables:
		  masqueradeAll: false
		  masqueradeBit: 14
		  minSyncPeriod: 0s
		  syncPeriod: 30s
		ipvs:
		  masqueradeAll: true
		  minSyncPeriod: 5s
		  scheduler: "rr"
		  syncPeriod: 30s
		kind: KubeProxyConfiguration
		metricsBindAddress: 127.0.0.1:10249
		mode: "ipvs"
		nodePortAddresses: null
		oomScoreAdj: -999
		portRange: ""
		udpIdleTimeout: 250ms
		EOF
		
		
		cat > kube-proxy/kube-proxy.service <<EOF
		[Unit]
		Description=Kubernetes Kube Proxy
		Documentation=https://github.com/kubernetes/kubernetes
		After=network.target
		
		[Service]
		ExecStart=/usr/local/bin/kube-proxy \\
		  --config=/etc/kubernetes/kube-proxy.conf \\
		  --logtostderr=true \\
		  --v=2
		  
		Restart=always
		RestartSec=10s
		
		[Install]
		WantedBy=multi-user.target
		EOF
		
		
		source /root/environment.sh
		for ip in ${ALL_IP[@]}; do
			rsync -v kube-proxy/kube-proxy.conf ${ip}:/etc/kubernetes/kube-proxy.conf
			rsync -v kube-proxy/kube-proxy.service ${ip}:/etc/systemd/system/kube-proxy.service
			rsync -v /etc/kubernetes/kube-proxy.kubeconfig ${ip}:/etc/kubernetes/kube-proxy.kubeconfig
			
			ssh ${ip} "systemctl daemon-reload && systemctl enable kube-proxy --now"
		done
		```
		
		配置 calico
		```bash
		# 下载calico 默认配置
		curl https://docs.projectcalico.org/manifests/calico.yaml -O
		
		# 在data段添加以下字段
		
		data:
		  etcd_endpoints: ##etcd_endpoints##
		  etcd_ca: ##etcd_ca##
		  etcd_cert: ##etcd_cert##
		  etcd_key: ##etcd_key##
		
		
		# 在 cni_network_config 段添加
		"plugins:" [
			{
				.
				.
				"log_level": "info",
				
				"etcd_endpoints": "__ETCD_ENDPOINTS__",
				"etcd_key_file": "__ETCD_KEY_FILE__",
				"etcd_cert_file": "__ETCD_CERT_FILE__",
				"etcd_ca_cert_file": "__ETCD_CA_CERT_FILE__",
				
				"log_file_path": "/var/log/calico/cni/cni.log",
				.
				.
				.
			}
		]
		
		
		source /root/environment.sh
		export ETCD_CA=`cat /etc/etcd/ssl/etcd-ca.pem | base64 | tr -d '\n'`
		export ETCD_CERT=`cat /etc/etcd/ssl/etcd.pem | base64 | tr -d '\n'`
		export ETCD_KEY=`cat /etc/etcd/ssl/etcd-key.pem | base64 | tr -d '\n'`
		export CIDR="172.16.0.0/12"
		
		sed -i -e "s@##etcd_endpoints##@${ETCD_ENDPOINTS}@g" -e "s@##etcd_ca##@${ETCD_CA}@g" -e "s@##etcd_cert##@${ETCD_CERT}@g" -e "s@##etcd_key##@${ETCD_KEY}@g" -e "s@# - name: CALICO_IPV4POOL_CIDR@- name: CALICO_IPV4POOL_CIDR@g" calico.yaml
		
		# 将 CALICO_IPV4POOL_CIDR 的 value 改为 172.16.0.0/12
		```
		
		> 如果在下载 calico 的所需镜像过慢，可以将 calico.yaml 内的镜像配置成国内的镜像源

## k8s 插件安装
安装 coredns
```bash
git clone https://github.com/coredns/deployment.git coredns-deployment

cd coredns-deployment/kubernetes
source /root/environment.sh
./deploy.sh -i ${CLUSTER_DNS_IP} | kubectl apply -f -

# 查看coredns状态
kubectl get all -n kube-system -l k8s-app=kube-dns
```

 安装 metrics server
 
 ```bash
 mkdir metrics-server
 wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml -O metrics-server/comp.yaml
 
 cd metrics-server
  
  # 在containers 段增加以下内容
  - args:
      - --kubelet-insecure-tls
	  - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.pem
	  - --requestheader-username-headers=X-Remote-User
	  - --requestheader-group-headers=X-Remote-Group
	  - --requestheader-extra-headers-prefix=X-Remote-Extra-

# 需要根据对应的 metrics-server 版本更换不同的镜像版本
# 将 metrics-server 镜像改为 registry.cn-hangzhou.aliyuncs.com/mxy/metrics-server:v0.5.0

# 在 168 行 name: tmp-dir 新增
- name: ca-ssl
  mountPath: /etc/kubernetes/pki

# 在 176 行 - emptyDir: {} 下面新增
- name: ca-ssl
  hostPath:
    path: /etc/kubernetes/pki
	
# 启动
kubectl create -f .

# 查看
kubectl top node
```
	  
安装 dashboard
```bash
mkdir /root/dashboard && cd /root/dashboard
curl -L https://raw.githubusercontent.com/kubernetes/dashboard/v2.3.1/aio/deploy/recommended.yaml -o dashboard-recommend.yaml

# 创建一个 dashboard 的管理员用户
cat > dashboard-user.yaml <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1 
kind: ClusterRoleBinding 
metadata: 
  name: admin-user
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kube-system
EOF

kubectl apply -f .

# 查看
kubectl get pods -n kubernetes-dashboard

# 修改类型将 ClusterIP 改为 NodePort
kubectl edit svc kubernetes-dashboard -n kubernetes-dashboard

# 查看暴露的端口号
kubectl get svc kubernetes-dashboard -n kubernetes-dashboard

# 访问 dashboard，必须是 https 访问,端口根据上一条命令暴露的端口更改
https://192.168.217.236:31614


# 获取token
kubectl describe secret -n kube-system $(kubectl get secrets -n kube-system  | grep admin-user | awk '{print $1}') | grep -E '^token' | awk '{print $2}'
```

验证集群是否成功
```bash
# 创建测试pod
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: dnsutils
  namespace: default
spec:
  containers:
  - name: dnsutils
    image: tutum/dnsutils
	command:
	- sleep
	- "3600"
    imagePullPolicy: IfNotPresent
  restartPolicy: Always
EOF
```
 校验
`kubectl get po`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808231920.png)
`kubectl exec dnsutils -n default -- nslookup kubernetes`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808232015.png)
`kubectl exec dnsutils -- nslookup kube-dns.kube-system`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808232321.png)
`kubectl exec dnsutils -- nslookup kubernetes.default`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808232446.png)
`nc -vz 10.96.0.1 443`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808232932.png)
`nc -vz 10.96.0.10 53`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808232946.png)
`kubectl get po -A -o wide`
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210808233341.png)
`kubectl exec dnsutils -- ping -c 172.23.119.132`

删除dnsutils
`kubectl delete po dnsutils`
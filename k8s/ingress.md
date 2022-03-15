# ingress
ingress 是 kubernetes 的一种 API 对象，将集群内部的 service 通过 HTTP/HTTPS 方式暴露到集群外部，并通过规则定义 HTTP/HTTPS 的路由。
ingress 特性：集群外部可访问的 URL、负载均衡、SSL Termination、按域名路由。

## 安装 ingress-nginx
使用 [Helm | Installing Helm](https://helm.sh/docs/intro/install/) 安装
```shell
# 添加仓库
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

# 拉取ingress-nginx
helm pull ingress-nginx/ingress-nginx
```

将拉取后的 tgz 包解压，并对里面的 values.yaml 文件做修改
```yaml
# before
dnsPolicy: ClusterFirst
# after
dnsPolicy: ClusterFirstWithHostNet

# before
hostNetwork: false
# after
hostNetwork: true

# before
kind: Deployment
# after
kind: DaemoneSet

# before
nodeSelector:
  kubernetes.io/os: linux
# after
nodeSelector:
  kubernetes.io/os: linux
  ingress: "true"
  
# before
type: LoadBalancer
# after
type: ClusterIP
```
上述修改让ingress使用宿主机网络，且可以通过label选择指定的node部署ingress

> 如果拉取不到 k8s.gcr.io 镜像，可以使用 image-syncer 工具将镜像同步至国内阿里云，再修改 values.yaml 内的 registry

安装 ingress
```shell
# 创建namespace
kubectl create ns ingress-nginx

# 选择一个节点，赋予ingress=true 标签
kubectl label node master-3 ingress=true

# 安装
helm install ingress-nginx -n ingress-nginx .
```
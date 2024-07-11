---
tags:
  - k8s
date created: 2021-06-12 17:10
date modified: 2024-04-06 01:29
---


![](attachments/Pasted%20image%2020220511040305.png)

# 主要组件

## master 节点

### ApiServer

kube-apiServer 是 kubernetes 的核心组件之一，主要提供一下功能：

- 提供集群管理的 REST API 接口
- 提供其他模块之间的数据交互和通信的枢纽（其他模块通过 APIServer 查询或修改数据，只有 APIServer 才直接操作 etcd）
- APIServer 提供 etcd 数据缓存以减少集群对 etcd 访问

### etcd

kubernetes 所有状态信息都需要持久化存储于 etcd 中，用于服务发现、共享配置、一致性保障
由于 etcd 的 watch 机制，所以其中的键值发生变化时，会通知到 APIServer。基于 watch 机制，k8s 的各组件实现高效协同。

### kubectl

[[kubectl]] 命令行工具，对集群进行管理，安装和部署容器化应用

## worker 节点

### kubelet

负责调度到对应节点的 Pod 的生命周期管理，执行任务并将 Pod 状态报告给主节点的渠道，通过容器运行时（拉取镜像、启停容器）来运行这些容器。它还会定期执行被请求的容器的健康探测程序。

- 从不同源获取 Pod 清单，并按需求启停 Pod 的核心组件
- 负责汇报当前节点的资源信息和健康状态
- 负责 Pod 的健康检查和状态汇报

### kubeproxy

负责节点网络，在主机上维护网络规则并执行连接转发；负责对正在服务的 pods 进行负载均衡。

监控集群中用户发布的服务，并完成负载均衡配置。

每个节点的 kube-proxy 都会配置相同的负载均衡策略，使得整个集群的服务发现建立在分布式负载均衡器上，服务调用无需经过额外的网络跳转。

负载均衡配置基于不同的插件实现：
- userspace
- 操作系统网络协议栈不同的 Hooks 点和插件
	- iptables
	- ipvs

## 核心附件

### KubeDNS

运行提供 DNS 服务的 pod，同一集群中的其他 pod 可使用此 DNS 服务，默认使用 coreDNS 为集群提供服务注册和服务发现的动态解析

### Dashboard

提供 web ui 管理集群

### Ingress Controller

在应用层实现的 http 负载均衡机制，它仅是一组路由规则集合，这些规则需要通过 ingress 控制器发挥作用。可用的项目有：Nginx、Envoy 等

# 资源抽象

## pod

Pod 是一组紧密关联的容器集合，它们共享 PID、IPC、Network 和 UTS namespace，是 kubernetes 调度的基本单位。

Pod 支持多个容器在一个 Pod 中共享网络和文件系统，可以通过进程间通信和文件共享方式组合完成服务。

同一个 Pod 中不同的容器可共享资源：
- 共享网络 Namespace
- 通过挂载存储卷共享存储
- 共享 Security Context

[pod](pod.md)

## service

pod 的端口和 ip 不是固定，需要通过一种方式提供一个稳定的 ip 供外部访问，所以使用 service 来对外提供服务。每个 service 会对应集群内部一个虚拟 IP，通过虚拟 ip 提供服务。

通过 service 也可以实现负载均衡。

service 常用类型：
- ClusterIP：提供一个集群内部地址，该地址只能在集群内解析和访问。是默认的服务类型
- NodePort: 每个集群节点上映射服务到一个静态的本地端口，从集群外部可以直接访问，并自动路由到内部自动创建的 ClusterIP
- LoadBalancer：使用外部的路由服务，自动路由访问到自动创建的 NodePort 和 ClusterIP
- ExternalName：将服务映射到 externalName 域指定的地址

## volume

提供数据持久化存储，支持高级的生命周期管理和参数指定，支持多种存储类型
常见的数据卷类型：
- emptyDir：当pod创建时，在节点上创建一个空挂载目录，挂载到容器内。当pod从节点离开，自动删除挂载目录内数据。节点上的挂载位置可以为物理硬盘或内存。这一类挂载适用于非持久化存储
- hostPath：将节点上某个已经存在的目录挂载到pod中，pod退出节点数据会保留
- gcePersistentDisk：使用GCE的Persistent Disk服务，pod退出数据保留
- awsElasticBlockStore：使用AWS的EBS volume服务，数据也会保留
- nfs：使用nfs协议的网络存储，持久化数据
- gitRepo：挂载一个空目录到pod，clone指定的git仓库代码，适用于直接从仓库中给定版本的代码来部署应用
- secret：传递敏感信息，基于内存的tmpfs，挂载临时密钥文件


## Namespace

通过命名空间来实现虚拟化，将同一组物理资源虚拟为不同的抽象集群，避免不同租户的资源发生命名冲突，和资源限额。

# 控制器抽象

- ReplicaSet：让集群始终维持某个 pod 的指定副本数的健康实例，副本集中的 pod 可以彼此替换

- Deployment：管理 pod 或 replicaset，支持升级操作

- StatefulSet：管理带状态的应用，可以为 pod 分配唯一的身份，确保在重新调配时也不会被替换

- Daemon：确保节点上肯定运行某个 pod，一般用来采集日志、监控节点、提供存储使用

- Job：短期任务处理场景。任务将创建若干 pod，并确保给定数目的 pod 完成任务后正常退出

- Horizontal Pod AutoScaler：自动水平扩展，根据 pod 的使用率自动调整 pod 的个数，保障服务可用性

- Ingress Controller：定义外部访问集群资源的一组规则，用来提供七层代理和负载

# 辅助概念

- Label：键值对，标记到资源对象上，用来对资源进行分类和筛选
- Selector：基于标签的使用正则表达式，筛选出一组资源
- Annotation：键值对，存放任意数据。一般用来添加对资源对象的详细说明
- Secret：存放敏感数据
- Name：提供给资源的别名，同类资源不能重名
- PersistentVolume：持久化存储，确保数据不会丢失
- Resource Context：资源限额，限制某个命名空间下对资源的使用
- Security Context：安全上下文，应用到容器上的系统安全配置
- Service Accounts：操作资源的账号
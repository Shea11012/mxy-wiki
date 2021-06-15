---
tags: ["k8s"]
---
# 基本概念
## service
- 拥有一个唯一指定的名字
- 拥有一个虚拟IP
- 能够提供某种远程服务能力
- 被映射到了提供这种服务能力的一组容器上

service 的服务进程是基于 socket 通信方式对外提供服务。一个 service 通常由多个相关的服务进程来提供服务，每个 service 都有一个独立的访问地址。

为了更好的管理服务，k8s 将每个服务都包装到相应的 pod 中，通过给每个 pod 贴上一个标签`name=mysql`，指定一个 service 作用于所有包含`name=mysql`的pod，解决了 service与pod的关联问题。

## pod 
pod 运行在 Node 环境中，Node 可以运行几百个 pod，每个 pod 中运行着一个特殊的 `pause` 的容器，其他容器则为业务容器，这些业务容器共享 `pause` 的容器的网络栈和 `volume` 挂载卷。不是每组pod和里面运行的容器都能映射到 service 上，只有提供服务的一组pod才会被映射成一个服务。

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210612171005.png)

> master 节点上运行着集群管理相关的一组进程 kube-apiserver 、kube-controller-manager 、kube-scheduler。这些进程实现了整个集群的资源管理、pod 调度、弹性伸缩、安全控制、系统监控和纠错管理等功能。

> node 管理的最小运行单元 pod。node 上运行着 kubelet、kube-proxy 服务进程。这些服务进程负责 pod 的创建、启动、监控、重启、销毁以及实现软件模式的负载均衡器。
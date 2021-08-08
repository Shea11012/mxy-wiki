---
tags: ["k8s"]
---
# 基本概念
## k8s master node
- kube-apiserver ：k8s 前端控制面板和暴露k8s API。它设计为水平扩展。
- etcd：k8s 存储所有集群信息
- kube-scheduler：监控新创建的 pods，如果没有一个 node 在上面运行则为它选择一个
- kube-controller-manager：逻辑上每一个控制器都是一个独立的进程。为了减少复杂度，讲它们编译进了一个独立的二进制文件且作为一个单一的进程运行。控制器包含：`node controller,replication controller,endpoints controller,service account,token controller`

## k8s worker node
- kubelet：运行在集群中每个node的客户端，它保证了容器运行在 pod 中。
- kube-proxy：通过设置主机上的网络开启了k8s的抽象服务，执行网络转发。
- container runtime：container runtime是一个软件负责容器运行。k8s 支持的运行时容器：`docker、rkt、runc 和任何实现了 OCI 的容器`
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
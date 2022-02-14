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


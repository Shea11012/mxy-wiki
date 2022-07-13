---
tags: ["k8s"]
date created: 2021-06-12 17:10
date modified: 2022-05-12 08:47
title: k8s
---

## k8s 架构
![](attachments/Pasted%20image%2020220511040305.png)

## 主要组件
### master 节点
- [APIServer](APIServer.md)：kubernetes 控制面板中唯一带有用户可以可访问 API 以及用户可交互的组件。API 服务会暴露 RESTful 的 kubernetes API。
- etcd：kubernetes 使用 etcd 作为存储
- [Controller Manager](Controller%20Manager.md)Controller Manager：运行着所有集群日常任务的控制器。包括节点控制器、副本控制器、端点控制器以及服务账户等。
- [Scheduler](Scheduler.md)：调度器会监控新建 Pod 并将其分配给节点。

### worker 节点
- [kubelet](kubelet.md)：负责调度到对应节点的 Pod 的生命周期管理，执行任务并将 Pod 状态报告给主节点的渠道，通过容器运行时（拉取镜像、启停容器）来运行这些容器。它还会定期执行被请求的容器的健康探测程序。
- [kube Proxy](kube%20Proxy.md)：负责节点网络，在主机上维护网络规则并执行连接转发；负责对正在服务的 pods 进行负载均衡。

## API 对象

API 对象是 kubernetes 集群中的管理操作单元，每个 API 对象都有四大类属性：

- TypeMeta
- MetaData
- Spec
- Status

## Pod

Pod 是一组紧密关联的容器集合，它们共享 PID、IPC、Network 和 UTS namespace，是 kubernetes 调度的基本单位。

Pod 支持多个容器在一个 Pod 中共享网络和文件系统，可以通过进程间通信和文件共享方式组合完成服务。

同一个 Pod 中不同的容器可共享资源：
- 共享网络 Namespace
- 通过挂载存储卷共享存储
- 共享 Security Context

[pod](pod.md)


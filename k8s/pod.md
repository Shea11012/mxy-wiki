---
date created: 2022-02-10 07:13
date modified: 2022-03-10 16:48
title: pod
---
## pod 
pod 运行在 Node 环境中，Node 可以运行几百个 pod，每个 pod 中运行着一个特殊的 `pause` 的容器，其他容器则为业务容器，这些业务容器共享 `pause` 的容器的网络栈和 `volume` 挂载卷。不是每组 pod 和里面运行的容器都能映射到 service 上，只有提供服务的一组 pod 才会被映射成一个服务。
pod 支持多个容器在一个 pod 中共享网络和文件系统

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210612171005.png)

> master 节点上运行着集群管理相关的一组进程 kube-apiserver 、kube-controller-manager 、kube-scheduler。这些进程实现了整个集群的资源管理、pod 调度、弹性伸缩、安全控制、系统监控和纠错管理等功能。

> node 管理的最小运行单元 pod。node 上运行着 kubelet、kube-proxy 服务进程。这些服务进程负责 pod 的创建、启动、监控、重启、销毁以及实现软件模式的负载均衡器。

## pod 探针
startupProbe：判断容器内应用程序是否启动，如果配置 startupProbe 则会禁止其他探测，直到它成功为止
livenessProbe：探测容器是否运行，如果失败，会根据配置的重启策略进行相应处理。默认是 success
readinessProbe：探测容器内的程序是否健康，返回 success 表示这个容器已经完成启动，并且程序可以接受流量。

### pod 探针监测方式
exec：执行命令，返回值为 0，则表示容器健康
grpc: grpc 健康监测，需要服务端实现 [GRPC Core: GRPC Health Checking Protocol](https://grpc.github.io/grpc/core/md_doc_health-checking.html),如果返回 servering 则表示成功
httpGet：通过 Get 请求访问一个指定的地址和端口，如果状态码在 200~400 之间则表示成功
tcpScoket：通过 TCP 连接一个指定的地址和端口，如果端口开放则表示容器健康

### 容器钩子
[Container Lifecycle Hooks | Kubernetes](https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks)，钩子可以通过命令或 http 请求实现
postStart：容器创建后就会被立即调用的钩子
preStop：在容器被销毁前，会执行的钩子


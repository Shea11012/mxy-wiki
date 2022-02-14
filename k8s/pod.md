## pod 
pod 运行在 Node 环境中，Node 可以运行几百个 pod，每个 pod 中运行着一个特殊的 `pause` 的容器，其他容器则为业务容器，这些业务容器共享 `pause` 的容器的网络栈和 `volume` 挂载卷。不是每组pod和里面运行的容器都能映射到 service 上，只有提供服务的一组pod才会被映射成一个服务。
pod 支持多个容器在一个pod中共享网络和文件系统

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/20210612171005.png)

> master 节点上运行着集群管理相关的一组进程 kube-apiserver 、kube-controller-manager 、kube-scheduler。这些进程实现了整个集群的资源管理、pod 调度、弹性伸缩、安全控制、系统监控和纠错管理等功能。

> node 管理的最小运行单元 pod。node 上运行着 kubelet、kube-proxy 服务进程。这些服务进程负责 pod 的创建、启动、监控、重启、销毁以及实现软件模式的负载均衡器。
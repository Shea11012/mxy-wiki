---
tags:
  - k8s
date created: 2021-06-12 17:10
date modified: 2023-12-01 16:01
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

kubernetes 使用 etcd 作为存储

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

# YAML 对象

一般使用 `kubectl create` 快速生成 yaml
```
kubectl create deployment web --image=nginx -o yaml --dry-run > web.yaml
```
API 对象是 kubernetes 集群中的管理操作单元，每个 API 对象都有四大类属性：

| 属性名称   | 介绍       |
| ---------- | ---------- |
| apiVersion | api 版本    |
| kind       | 资源类型   |
| metadata   | 资源元数据 |
| spec       | 资源规格   |
| replicas   | 副本数量   |
| selector   | 标签选择器 |
| template   | pod 模板    |
| metadata   | pod 元数据  |
| spec       | pod 规格    |
| containers | 容器配置           |

# 概念

## pod

Pod 是一组紧密关联的容器集合，它们共享 PID、IPC、Network 和 UTS namespace，是 kubernetes 调度的基本单位。

Pod 支持多个容器在一个 Pod 中共享网络和文件系统，可以通过进程间通信和文件共享方式组合完成服务。

同一个 Pod 中不同的容器可共享资源：
- 共享网络 Namespace
- 通过挂载存储卷共享存储
- 共享 Security Context

[pod](pod.md)

## controller

pod 通过 controller 实现应用的运维，如弹性伸缩、滚动升级等。
pod 和 controller 通过 label 建立关系

[[controller]]

## service

pod 的端口和 ip 不是固定，需要通过一种方式提供一个稳定的 ip 供外部访问，所以使用 service 来对外提供服务。每个 service 会对应集群内部一个虚拟 IP，通过虚拟 ip 提供服务。

通过 service 也可以实现负载均衡。

service 常用类型：
- ClusterIP：
	- 普通 service：分配一个集群内部可访问的固定虚拟 IP
	- headless service：不会分配 cluster ip，也不同过 kube-proxy 做反向代理和负载均衡，通过DNS提供稳定的
- NodePort: 对外访问
- LoadBalancer：对外访问，需指定负载均衡器
- ExternalName：主要面向运行在集群的外部服务，将外部服务映射进集群，使具备 k8s 内服务的一些特征
---
tags: [k8s]
date created: 2022-05-11 08:04
date modified: 2022-05-11 08:06
title: Kube-Proxy
---

# Kube-Proxy
监控集群中用户发布的服务，并完成负载均衡配置。
每个节点的 kube-proxy 都会配置相同的负载均衡策略，使得整个集群的服务发现建立在分布式负载均衡器上，服务调用无需经过额外的网络跳转。
负载均衡配置基于不同的插件实现：
- userspace
- 操作系统网络协议栈不同的 Hooks 点和插件
	- iptables
	- ipvs
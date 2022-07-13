---
tags: [k8s]
date created: 2022-05-11 04:10
date modified: 2022-05-11 04:12
title: APIServer
---

# APIServer

kube-apiServer 是 kubernetes 的核心组件之一，主要提供一下功能：

- 提供集群管理的 REST API 接口
- 提供其他模块之间的数据交互和通信的枢纽（其他模块通过 APIServer 查询或修改数据，只有 APIServer 才直接操作 etcd）
- APIServer 提供 etcd 数据缓存以减少集群对 etcd 访问
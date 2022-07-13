---
date created: 2021-12-03 20:20
date modified: 2022-03-18 18:32
title: 简介
---
# 简介
prometheus 通过 http 周期性抓取被监控组件状态，任意组件只要提供对应的 http 接口并且符合 prometheus 定义的数据格式，就可以接入监控。

主要特征：
- 多维度数据模型
- 灵活的查询语言
- 不依赖分布式存储，单个服务器节点是自主的
- 以 http 方式，通过 pull 模型去拉取时间序列数据
- 可以通过中间网关支持 push 模型
- 通过服务发现或者静态配置，发现目标服务对象

不适用场景：实时数据统计

![[Pasted image 20220318183244.png]]

## Prometheus Server
负责实现对监控数据的获取、存储以及查询。
Prometheus Server 本身就是一个时序数据库，将采集到的监控数据按照时间序列的方式存储在本地磁盘中。

## Exporters
将监控数据采集的端点通过 HTTP 服务形式暴露给 Prometheus Server。

## AlertManager
基于 PromQL 创建告警规则。

## PushGateway
将数据推送到 Gateway 中，Prometheus Server 从 Gateway 中 pull 数据。
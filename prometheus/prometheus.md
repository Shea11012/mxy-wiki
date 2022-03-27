# 简介
prometheus 通过http周期性抓取被监控组件状态，任意组件只要提供对应的http接口并且符合 prometheus 定义的数据格式，就可以接入监控。

主要特征：
- 多维度数据模型
- 灵活的查询语言
- 不依赖分布式存储，单个服务器节点是自主的
- 以http方式，通过pull模型去拉取时间序列数据
- 可以通过中间网关支持push模型
- 通过服务发现或者静态配置，发现目标服务对象

不适用场景：实时数据统计

![[Pasted image 20220318183244.png]]

## Prometheus Server
负责实现对监控数据的获取、存储以及查询。
Prometheus Server 本身就是一个时序数据库，将采集到的监控数据按照时间序列的方式存储在本地磁盘中。

## Exporters
将监控数据采集的端点通过HTTP服务形式暴露给Prometheus Server。

## AlertManager
基于PromQL创建告警规则。

## PushGateway
将数据推送到Gateway中，Prometheus Server 从Gateway中pull数据。
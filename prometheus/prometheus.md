---
tags: 
date created: 2021-12-03 20:20
date modified: 2024-04-02 02:04
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

- 负责实现对监控数据的获取、存储以及查询。
- Prometheus Server 本身就是一个时序数据库，将采集到的监控数据按照时间序列的方式存储在本地磁盘中。

## Exporters

将监控数据采集的端点通过 HTTP 服务形式暴露给 Prometheus Server。

## AlertManager

基于 PromQL 创建告警规则，满足 PromQL 规则，会产生一条告警，在 AlertManager 可以用邮件、slack 等内置的通知方式集成，也可以使用 WebSocket 自定义告警方式。

## PushGateway

当 exporter 无法直接与 Prometheus 直接通信时，可以将数据推送到 Gateway 中，Prometheus Server 从 Gateway 中 pull 数据。

# 指标定义

`<metric name>{<label name>=<label value>,...}`
metric name：表示指标含义
label：体现指标的维度特征

>[!example]
>api_http_requeset_total{method="post",handler="/messages"}

# 指标类型

## Counter

累积指标，代表一个单调递增的计数器，其值只能增加或在重新启动时重置为零。
通常用于统计数量。

## Gauge

表示任意上下浮动的单个数值，通常用于测量值，如内存、CPU、磁盘、并发数等。

## histogram

柱状图，在 Prometheus 的查询语言中，有三个作用：
- 对每个采样点进行统计，打倒各个桶中
- 对每个采样点值累计和
- 对采样点的次数累计和
观察桶的累积计数器：`<name>_bucket{le="upper inclusive bound>"}`
所有观察值的总和：`<name>_sum`
观察到的事件数：`<name>_count`

## Summary

采样点分位图统计，在客户端对于一段时间内的每个采样点进行统计，采样点分位图统计，用于得到数据的分布情况
- 观察到的事件的φ分位数（0≤φ≤1）分位数：`<basename>{quantile="<φ>"}`
- 统计累计和: `<name>_sum`
- 统计总数: `<name>_count`

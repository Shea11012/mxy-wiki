---
tags: 
date created: 2021-12-14 21:54
date modified: 2024-04-02 02:35
title: PromQL
---

# 查询结果类型

- 瞬时数据（Instant vector）：包含一组时序，每个时序只有一个点，如：`http_requests_total`
- 区间数据（Range vector）：包含一组时序，每个时序有多个点，如：`http_requests_total[5m]`
- 纯量数据（Scalar）：纯量只有一个数字，没有时序，如：`count(http_requests_total)`

# 查询条件

>[!example]
>http_requests_total{code="200"} 查询 code=200 的数据
> ## 查询条件支持正则匹配
>http_requests_total{code!="200"} code!=200 的数据
>http_requests_total{code=~"2.."} code 为 2xx 的数据
>http_requests_total{code!~"2.."} code 不为 2xx 的数据

# 操作符

## 算术操作符

```
+,-,*,/,%,^
```

## 比较运算符

```
==,!=,>,<,>=,<=
```

## 逻辑运算符

```
and,or,unless
```

## 聚合运算符

```
sum,min,max,avg,stddev,stdvar,count,count_values_bottomk,topk,
qunatile,
```

# 瞬时向量

在给定的时间戳中选择一组时间序列和每个时间序列的单个数值，如：`http_requests_total` 返回一个包含所有时间序列元素的瞬时向量

# 范围向量

从当前瞬间选择返回的数据范围。
>[!example]
>// 获取过去 5 分钟内记录的所有值
>http_requests_total{job="prometheus"}[5m]

# 修饰符

## offset 偏移量修饰符

offset 偏移量修饰符允许更改查询中单个瞬时向量和范围向量的时间偏移量
>[!example]
>// 相对于当前查询计算时间过去 5 分钟的值
>http_requests_total offset 5m
>// 返回一周前的前 5 分钟速率
>rate(http_requests_total[5m] offset 1w)
>// 指定负偏移，需要开启 --enale-feature=promql-negative-offset
>rate(http_requests_total[5m] offset -1w)

## @修饰符

允许在查询中更改单个瞬时向量和范围向量的计算时间。提供一个 unix 时间戳给@修饰符
>[!example]
>http_requests_total @ 1609746000

@ 修饰符在 int64 范围，可以与 offset 修饰符一起使用，offset 是相对于@修饰符的
默认是禁用的，通过设置 `--enable-feature=promql-at-modifier`。start 和 end 可以作为@修饰符。
对于范围查询，分别解析到范围查询的开始和结束
对于即时查询，都会被解析为计算时间





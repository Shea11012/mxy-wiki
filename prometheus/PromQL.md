---
date created: 2021-12-14 21:54
date modified: 2022-03-09 19:02
title: PromQL
---
## 时间序列
### 样本（sample）
Prometheus 将采集到的样本数据以时间序列（time-series）方式保存在内存数据库中，并定时持久化到硬盘。
向量：时间序列按照时间戳和值的序列顺序存放。
每条时间序列通过指标名称（metrics-name）和一组标签集（labelset）命名。
样本由三部分组成：
- 指标（metric）：mertic name 和描述当前样本特征的 labelsets。
- 时间戳（timestamp）：精确到毫秒的时间戳
- 样本值（value）：float64 的浮点型数据表示当前样本的值

```
<--------------- metric ---------------------><-timestamp -><-value->
http_request_total{status="200", method="GET"}@1434417560938 => 94355
```

### 指标（metric）
指标名称反映监控样本含义
标签反映当前样本的特征维度
Prometheus 底层是以 `__name__=<metric name>` 形式保存在数据库中的

```
api_http_requests_total{method="POST", handler="/messages"}

{__name__="api_http_requests_total"，method="POST", handler="/messages"}
```


## Metrics 类型
Prometheus 定义了 4 中类型：
- Counter（计数器）
- Gauge（仪表盘）
- Histogram（直方图）
- Summary（摘要）

### Counter 只增不减计数器
只增不减（除非系统重置）。常见监控指标，如 http_requests_total、node_cpu。
PromQL 内置的聚合操作和函数可以让用户对这些数据进一步分析
```
获取请求增长率
rate(http_requests_total[5m])

访问量前10的HTTP地址
topk(10,http_requests_total)
```

### Gauge 可增减仪表盘
Gauge 侧重反应系统当前状态。常见指标：node_memory_MemFree_bytes、node_memory_MemAvailable
```
node_memory_MemFree_bytes

计算CPU温度在两个小时内的差异
delta(cpu_temp_celsius{host="zeus"}[2h]
```

### Histogram 和 Summary
Histogram 和 Summary 主要用于统计和分析样本的分布情况
对于分位数计算，Summary 分位数是在客户端计算；Histogram 在服务端计算

## 查询时间序列
查询该指标下的所有时间序列
```
http_requests_total

http_requests_total{}
```

PromQL 支持用户根据时间序列的标签匹配模式来对时间序列进行过滤：
- `=` label = value
- `!=` label != value
- `=~`  label=~regx 
- `!~` label !~ regx

## 范围查询
以 `http_requests_total` 查询时间序列时，返回值只包含该时间序列中的最新一个样本值，这种返回结果被称为 **瞬时向量**
`http_requests_total[5m]` 获取一段范围内的样本数据，被称为 **区间向量**

## 时间位移操作
瞬时向量或区间向量，都是以当前时间为基准
```
查询5分钟前的瞬时样本
http_requests_total offset 5m

查询昨天一天前的样本
http_requests_total[1d] offset 1d
```

## 操作符
PromQL 支持数学运算、布尔运算
**集合运算符**
- and：`vector1 and vector2` 产生一个由 vector1 元素组成的新向量。该向量包含 vector1 中完全匹配 vector2 中的元素组成
- or：`vector1 or vector2` 产生一个包含 vector1 中所有样本数据，以及 vector2 中没有与 vector1 匹配到的样本数据
- unless：`verctor1 unless vector2` 产生一个由 vector1 中没有与 vector2 匹配的元素组成。

### 匹配模式
向量与向量之间运算的默认匹配规则：依次找到与左边向量元素匹配（标签完全一致）的右边向量元素进行运算，如果没找到则丢弃。
在操作符两边标签不一致时，可以使用：
- on(label list) 将匹配行为限定在某些便签之内
- ignoring 在匹配时忽略某些便签

```
vector1 operator  on(label list) vector2
vector1 operator ignoring(label list) vector2
```

#### 多对一和一对多
必须使用 group 修饰符：group_left 或 group_right 来确定哪一个向量具有更高的基数。
```
vector1 operator ignoring(label list) group_left(label list) vector2
```

## 聚合操作
聚合操作中某些函数可以传递参数，without 从计算结果中移除列举的标签；by 只保留列出的标签
```
aggregation-opeartor([parameter],vector) [without|by (label list)]
```

## HTTP API
### 响应数据类型
- vector：瞬时向量
- matrix：区间向量
- scalar：标量
- string：字符串


### 瞬时数据查询
> [!note] API
> GET /api/v1/query
>
>> [!note] URL param
>> - query: PromQL 表达式
>> - time：指定计算PromQL 时间戳。可选参数，默认当前系统时间
>> - timeout: 超时设置。可选参数，默认使用 -query.timeout 全局设置


### 区间查询
> [!note] API
> GET /api/v1/query_range
> 
>> [!tip]
>> query_range API中的查询表达式只能使用瞬时向量选择器类型表达式
>>
>>
>> [!note] URL Param
>> - query：PromQL 表达式
>> - start：起始时间
>> - end：结束时间
>> - step：查询步长
>> - timeout：超时设置

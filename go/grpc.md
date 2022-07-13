---
date created: 2021-11-30 21:22
date modified: 2022-02-28 19:45
title: grpc
---
### grpc 四种调用方式

- unary rpc：一元 rpc，客户端使用 stub 发送请求等待服务端响应
- server-side rpc：服务端流式 rpc，服务端返回流式数据
- client-side rpc：客户端流式 rpc，客户端发送流式数据
- bidirectional streaming rpc：双向流式 rpc

#### grpc 负载均衡策略

grpc 负载均衡是基于每个调用进行的，而不是每个连接

常见的负载均衡类型：

- 客户端负载：
  - 优点：高性能，去中心化，不需要借助独立的外部负载均衡组件
  - 缺点：实现成本较高，不同语言需要实现各自对应的均衡策略
- 服务端负载（代理模式）：
  - 优点：简单，透明，方便
  - 缺点：负载均衡器可能会成为性能瓶颈，且必须要保持高可用
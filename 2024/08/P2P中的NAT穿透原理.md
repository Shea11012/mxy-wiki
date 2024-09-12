---
created: 2024-08-08
tags: [network,nat,p2p]
updated: 2024-08-09
---

P2P 在实时音视频功能通常都会用到 P2P，主要场景有 IM、直播等

# NAT 的类型

### 完全锥型（Full Cone NAT）

特点：IP 和端口都不受限
将内部主机的 IP 和端口，映射到公网 IP 和 X(随机) 端口。任意外部地址通过对公网 IP 和 X 端口的访问都可以被重新定位到内部主机

### 受限锥型（Restricted Cone NAT）

特点：IP 受限，端口不受限
在映射端口后，不允许所有 IP 进行对于该端口的访问，要想通信必须内部主机对外部 IP 发起过连接，然后外部主机才能进行通信，端口不受限。如 (NatIP:NatPort -> A:P1)，A 可以通过任意端口访问 NatIP，但是 B 则不能。

### 端口受限型（Port Restricted Cone NAT）

特点：IP 和端口都受限
包含受限锥型的特性，对端口也有限制。只有内部主机发过报文给外部主机，外部主机才能以公网 IP 和指定端口进行通信

### 对称型（Symmetric NAT）

特点：对每个外部主机或端口的会话都会映射为不同的端口
只有来自同一内部 IP、PORT，且针对同一目标 IP、PORT 的请求才被 NAT 转换至同一个公网 IP、PORT，否则 NAT 将分配一个新的公网 IP、PORT，并且只有收到过内部主机请求的外部主机才能向内部主机发送数据包。

## NAT 路由类型的判断

为了实现点对点通信，有多种方案: STUN(Session Traversal Utilities for NAT)、UPNP 技术、ALG 应用层网关识别、SBC 会话边界控制、ICE 交互式连接、TURN 中继 NAT 穿透等

# STUN 协议介绍

STUN 是一种网络协议，它允许位于 NAT 后的客户端找出自己的公网地址，该协议 [RFC5389](https://tools.ietf.org/html/rfc5389) 定义

STUN 由三部分组成：
- STUN 客户端
- STUN 服务端
- NAT 路由器


![[../../Excalidraw/P2P中的NAT穿透原理 2024-08-08 22.36.59.excalidraw]]

1. 检测客户端是否能进行 UDP 通信以及客户端是否位于 NAT 后 - 测试 1
2. 检测防火墙类型 - 测试 2
3. 检测客户端 NAT 是否 Full Cone NAT - 测试 2-1
4. 检测客户端 NAT 是否 Symmetric NAT - 测试 1-1
5. 检测客户端 NAT 是否 Restricted Cone 还是 Port Restricted Cone - 测试 3
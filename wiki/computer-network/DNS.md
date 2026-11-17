---
title: "DNS 域名系统"
created: 2026-07-22
updated: 2026-07-22
status: needs-review
type: synthesis
tags:
  - computer-network
  - computer-network/dns
confidence: INFERRED
---
# DNS 域名系统

> 本页由重复旧笔记整合而来，保留独有内容；明显错误另按已列明的一手资料修正。

## DNS 协议、记录与报文

# DNS

- 由分层的 DNS 服务器实现的分布式数据库。
- 使主机能够查询分布式数据库的应用层协议
- DNS 同时使用 UDP 和 TCP。传统 DNS 的 UDP 载荷上限为 512 字节；EDNS(0) 可以通告更大的 UDP 载荷能力。响应被截断或 UDP 不适用时，客户端可改用 TCP。
- 区域传送使用 TCP；普通查询通常优先使用 UDP，但不能概括为“域名解析只使用 UDP”。

# DNS 提供的服务

- 主机名到 IP 地址转换
- 主机别名
- 邮件服务器别名
- 负载分配：一组 IP 地址和一个主机名关联，在进行 DNS 解析时，可以通过负载均衡算法分配地址。

# DNS 服务器类型

![[../../assets/计算机网络/协议/DNS/IMG-20231107221819057.png]]

- 根区配置了 13 个命名权威服务器标识（A 到 M），背后由全球数百台服务器实例提供服务；`dig . NS` 查询到的是这些命名权威服务器，而不是只有 13 台物理服务器。
- 顶级域 DNS 服务器：如 .com,.net,.cn 后缀
- 次级域名：如 baidu.com，中的 baidu
- 三级域名：如 www.baidu.com 中 www

# DNS 解析顺序

1. 本地 DNS 缓存
2. hosts 文件
3. 访问本地 DNS 服务器
4. 从根 DNS 服务器查到顶级 DNS 服务器
5. 从顶级 DNS 服务器查到次级 DNS 服务器
6. 从次级 DNS 服务器查主机名的 IP 地址

# DNS 报文格式

![[../../assets/计算机网络/协议/DNS/IMG-20231107222655661.png]]

## DNS 记录类型

域名除了关联 IP 地址外，还可以关联其他类型消息

- A：主机 IP 地址
- AAAA：主机 IPv6 地址
- ALIAS：自动解析的别名
- NS：域名服务器记录，返回下一级域名信息的服务器地址
- MX：邮件记录，返回接受电子邮件的服务器地址
- CNAME：规范名称记录，返回另一个域名，即当前查询域名是另一个域名的跳转
- PTR：逆向查询记录，只用于从 IP 地址查询域名
- TXT：描述文本

## DNS 作用与解析过程

DNS（Domain Name System）是域名系统的英文缩写，是一种组织域层次结构的计算机和网络服务命名系统，用于 TCP/IP 网络。

## DNS 的作用

一种能够进行主机名到 IP 地址转换的目录服务，域名系统作为将域名和 IP 地址相互映射的一个分布式数据库，使的用户能够更快的访问互联网。
域名解析层次通常包括：
- 根区的 13 个命名权威服务器标识（实际对应全球数百台实例）
- 第二层顶级域名服务器
- 第三层权威服务器
还有一个特殊的 DNS 服务器，叫做本地 DNS 服务器

![](attachments/Pasted%20image%2020230308183224.png)

## DNS 域名解析原理

![](attachments/Pasted%20image%2020230309052244.png)
1. 浏览器中输入 `www.baidu.com` 域名，操作系统会先检查本地的 host 文件
	1. 如果本地配置了 dns 缓存，则会在缓存中在进行查找
2. 根据连接网络时配置的 dns 服务器（这个 dns 服务器叫本地服务器），去本地 dns 服务器查询
3. 本地 dns 服务器没有的话，会去根 DNS 服务器查询，根服务器返回顶级域名服务器，顶级域名服务器返回权威域名服务器，权威域名服务器返回对应域名的 IP 地址。

## 核验依据

- [IANA Root Servers](https://www.iana.org/domains/root/servers)
- [RFC 6891: EDNS(0)](https://www.rfc-editor.org/rfc/rfc6891.html)

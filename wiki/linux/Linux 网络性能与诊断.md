---
title: "Linux 网络性能与诊断"
created: 2026-07-22
updated: 2026-07-22
status: needs-review
type: synthesis
tags:
  - linux
  - linux/network
  - linux/performance
confidence: INFERRED
---
# Linux 网络性能与诊断

> 本页由重复旧笔记整合而来；正文按原页面分节保留，尚未重新核验技术时效性。

## 网络性能指标与调优

## 性能指标

### 最大传输单元（mtu）

mtu 限制了数据链路层的传输大小，因此也限制了 IP 层大小
```shell
# 查看mtu大小
ip link
```

### 带宽

查看链路的最大传输速度
```shell
ethtool {dev}
```

### 传输速率

以网络包为单位的传输速率
```shell
sar -n DEV
```

### TCP 半连接队列大小

```
max(64,net.ipve.tcp_max_syn_backlog)
```

### TCP 全连接队列大小

```
min(backlog,net.core.somaxconn)
```

### TCP 缓冲区大小

```
# 发送缓冲区
cat /proc/sys/net/ipv4/tcp_wmem

# 接受缓冲区
cat /proc/sys/net/ipv4/tcp_rmem
```

### UDP 缓冲区大小

```
cat /proc/sys/net/ipv4/udp_mem
```

## 调优方法

### 链路层

#### 网卡驱动优化

- NAPI：中断 + 轮询
- 为网卡中断配置 CPU 亲和性（smp_affinity），将这些中断处理程序调度到不同的 CPU 上执行

#### DPDK

跳过逻辑复杂的 Linux 网络协议栈，直接由用户态进程用轮询的方式，来处理网络请求

### 传输层

#### TCP

##### TCP TIME_WAIT 优化

- 增大处于 TIME_WAIT 的连接数量：`net.ipv4.tcp_max_tw_buckets`
- 减小 `net.ipv4.tcp_fin_timeout`，让处于 fin 状态的 tcp 连接尽早释放
- 开启端口复用 `net.ipv4.tcp_tw_reuse`，这样被 TIME_WAIT 状态占用的端口能很快被用于新的连接

##### 增加可用端口号或文件数

- 增大本地端口范围 `net.ipv4.ip_local_port_range`，可以支持更多连接，提高服务并发能力
- 增大最大文件描述符数量 `ulimit -n 10240`

##### SYN 优化

- 增大 TCP 半连接数量 `net.ipv4.tcp_max_syn_backlog`
- 减少 SYN_RECV 状态的连接重传 SYN+ACK 包次数 `net.ipv4.tcp_synack_retries`

##### 其他优化

- tcp nagle 算法，启用时可提高小包场景下的网络利用率，同时也增加了网络延迟
- tcp 延迟确认

#### UDP

- 增大 UDP 缓冲区大小
- 增大本地端口范围
- 根据 MTU 大小，调整 UDP 数据包大小，避免分片发生

## 网络诊断与压测工具

## tcpdump

- -i：指定网络接口
- -nn：不解析 IP 地址和端口号
- -c {count}：指定要抓取的数据包个数
- -w {file}：抓包并保存到文件

## netstat

- -a：显示所有 socket 连接
- -l：只显示 listening 状态的 socket 连接
- -n：显示数字地址和端口
- -p：显示进程信息
- -t：显示 tcp 连接
- -u：显示 udp 连接
- -x：显示 unix 连接
- -i：显示所有网络接口的状态
- -s：汇总 ip、icmp、udp、tcp 等各种协议的收发统计信息

## ss

与 netstat 参数类似，但支持 state-filter 和 expression

### state-filter

所有的TCP标准状态
- established
- syn-sent
- syn-recv
- fin-wait-1
- fin-wait-2
- time-wait
- closed
- close-wait
- last-ack
- listening
- closing

## iperf

- -f [kmKM]：以不同的格式显示报告
- -s：server 模式启动
- -c：client模式启动
- -p：指定端口号
- -u：使用udp协议
- -t {sec}：测试多久，一般指定在client（client是发包方）
- -m：在结果中打印mss值
- -N：禁止Nagle算法
>[!note]
>```shell
># 测试TCP吞吐量
>server：iperf -s -p 10000
>clinet：iperf -c {ip} -p 10000
>```

## ab

- -c：并发数
- -n：总请求数
- -s：设置每个请求的超时时间
- -r：套节字错误时仍然继续执行
>[!tip]
>测试高并发时，需要提前修改进程能打开的最大文件描述符大小。`ulimit -n 10240`

## hping3

测试网络及主机安全
模式选择
- --rawip：RAWIP模式
- --icmp：ICMP模式
- --UDP：UDP模式
- --SCAN：SCAN模式，扫瞄对应端口
- --listen：监听模式
参数：
- -S：设置TCP协议的SYN
- -p {port}：端口号
- -c {count}：发送的总封包数
- -i {interval}：每隔interval时间发送网络帧
- -a {hostname}：源地址欺骗
- --rand-source：随机源IP
- --flood：尽快发包，无须指定-i

>[!example]
> // 端口扫描
>`hping3 ip --scan 80 -S`
> // 扫描多个端口
>`hping3 ip --scan 80,8080 -S`
> // 模拟SYN DDOS攻击
>`hping3 -S -p 80 -i u10 ip`
> // flood+随机源地址SYN攻击
>`hping3 -S -p 80 ip --flood --rand-source`
> // 伪造源IP地址
>`hping3 -S -p 80 -a ip -S ip`

## nethogs

实时统计网络带宽使用率
>[!example]
>// 显示每个进程所使用的网络带宽
>`nethogs -d 2 -s`

## iftop

套接字及进程的网络利用率
>[!example]
>`iftop -i eth0 -B -n`

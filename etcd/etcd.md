---
date created: 2021-11-30 21:22
date modified: 2022-04-15 17:13
title: etcd
---
# etcd

|            | etcd               | zookeeper          |
| ---------- | ------------------ | ------------------ |
| 一致性算法 | raft               | zab                |
| 数据模型   | 基于目录的层次模式 | 基于目录的层次模式 |
| kv 存储引擎 | 简单内存树         | concurrent hashMap |
| 部署       | 简便               | 复杂               |

#### etcd 与 redis 区别

数据复制：

- redis：主备异步复制，可能会丢失数据
- etcd：raft，为了读写一致性，读写性能比 redis 差

数据存储：

- redis：存储用户数据，可以承载上 T 数据
- etcd：低容量的关键元数据，db 大小一般不会超过 8g

数据结构：

- redis：丰富的数据结构
- etcd：仅 kv，持久化存储使用 boltdb

#### etcd 执行流程图

![](etcd.assets/457db2c506135d5d29a93ef0bd97e4bb.png)

### 搭建 etcd 集群

使用进程管理工具 [goreman](https://github.com/mattn/goreman) 可以快速创建、停止本地的多节点 etcd 集群。

下载 etcd 的二进制文件，从 [etcd 源码](https://github.com/etcd-io/etcd/blob/v3.4.9/Procfile) 中下载 goreman Profile 文件，它描述了 etcd 进程名、节点数、参数等信息。

使用 `goreman -f Profile start` 就可以快速启动 etcd 集群



### 鉴权

#### 密码认证

```bash
etcdctl user add root:root # 创建root账号
etcdctl auth enable # 启用鉴权
```

#### Simple Token

当用户身份验证通过后，生成一个随机的字符串返回给 client。etcd 默认会每 5 分钟检查 token 是否过期。

simple token 可描述性差，client 无法通过 token 获取到过期时间、用户名、签发者等信息。

#### JWT token

应当使用 jwt token 替代 simple token

#### 证书认证

etcd 中如果启用了 client 证书认证（--client-cert-auth) 。它会取 Subject 中 CN 字段作为用户名。

![[etcd.assets/55e03b4353c9a467493a3922cf68b294.png]]

#### RBAC

```bash
# 创建一个 admin role
etcdctl role add admin --user root:root
# 分配一个可读写范围数据的权限给 admin role
etcdctl role grant-permission admin readwrite hello name --user root:root
# 将用户 alice 和 admin role 关联起来
etcdctl user grant-role alice admin --user root:root
```

### 租约 Lease

Lease 一种活性检测机制。client 和 server 之间存在一个约定，内容是 etcd server 保证在约定的有效期（TTL），不会删除你关联到 Lease 上的 kv。

应用场景：leader 选举、k8s event 自动淘汰、服务发现场景故障节点自动剔除等问题。

![](etcd.assets/ac70641fa3d41c2dac31dbb551394b7c.png)

> 两个常驻 goroutine
>
> - RevokeExpiredLease：定时检查是否有过期 lease
> - CheckpointScheduleLeases：定时触发更新 lease 剩余到期时间
>
> 各个接口：
>
> - Grant：创建一个指定秒数（TTL）的 lease，lessor 会将 lease 信息持久化存储在 boltdb 中
> - Revoke：撤销 lease 并删除其关联的数据
> - LeaseTimeToLive：获取一个 lease 的剩余时间
> - LeaseKeepAlive：为 lease 续期

### 事务

![](etcd.assets/e41a4f83bda29599efcf06f6012b0bd3.png)

etcd 提供 CAS 操作，**只支持单 key，无法实现事务的各个隔离级别**

```
client.Txn(ctx).If(cmp1, cmp2, ...).Then(op1, op2, ...,).Else(op1, op2, …)
```

#### IF 检查项

- mod_revision：key 的最近一次修改版本号。通过它检查 key 最近一次被修改时的版本号是否符合预期。
- create_revision：key 的创建版本号。检查 key 是否已经存在。
- version：key 的修改次数。检查 key 的修改次数是否符合预期
- value：key 的 value 值



### 最佳实践

- 开启 etcd 数据毁坏检测功能：
  - --experimental-initial-corroupt-check：检查各个节点数据是否一致，如果不相等则无法启动
  - --experimental-corrupt-check-time：每隔一定时间检查数据一致性
- 鉴权使用证书认证；密码认证确保复用 token（不推荐）
- 做好备份
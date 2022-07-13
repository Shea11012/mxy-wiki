---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: nftables
---
# nftables
替换 iptables 的网络包过滤工具

nfttables 和 iptables 一样，由表（table）、链（chain）、规则（rule）组成，其中表包含链，链包含规则。

nftables 每个表只有一个地址簇，并且只适用于该簇的数据包。表可以指定五个簇中的一个。

- ip
- ipv6
- inet：适用于 ipv4 和 ipv6 的数据包
- arp
- bridge

## 创建表

创建表：`nft add table inet my_table`

列出所有规则：`nft list ruleset`

## 创建链

链是用来保存规则的，和表一样，链需要被显示创建，链有两种类型：

- 常规链：不需要指定钩子类型和优先级，可以用来做跳转，从逻辑上对规则分类
- 基本链：数据包的入口点，需要指定钩子类型和优先级

创建常规链：`nft add chain inet my_table my_utility_chain`

创建基本链：

```bash
nft add chain inet my_table my_filter_chain { 
	type filter hook input priority 0 \; # 此处\ 转义，priority 定义优先级，值越小越优先
}


# 列出链中的所有规则
nft list chain inet my_table my_utility_chain
nft list chain inet my_table my_filter_chain
```

## 创建规则

规则由语句或表达式构成，包含在链中

`nft add rule inet my_table my_filter_chain tcp dport ssh accept` 表示允许 ssh 登录

**add** 将规则添加到链尾，**insert**将规则添加到链头

`nft insert rule inet my_table my_filter_chain tcp dport http accept`

使用 **index** 指定规则索引，**add**表示新规则添加在索引位置的规则后面，**insert** 表示新规则添加在索引位置的规则前面

## 删除规则

单个规则只能通过句柄删除：`nft --handle list ruleset`  handle 参数列出句柄值

使用句柄值删除规则：`nft delete rule inet my_talbe my_filter_chain handle 8`

## 列出规则

列出某个表中的所有规则：`nft list table inet my_table`

列出某个链中的所有规则：`nft list chain inet my_table my_other_chain`

## 集合

集合分为 **匿名集合**与**命名集合**，匿名集合比较适用于将来不需要改变的规则。

```
nft add rule inet my_table my_filter_chain ip saddr {
	10.10.10.12, 10.10.10.231
} accept # 表示允许来自源 ip 处于 10.10.10.123 ~ 10.10.10.231 这个区间的主机流量
```

命名集合，创建集合需要指定元素类型，有如下类型：

- ipv4_addr：ipv4 地址
- ipv6_addr：ipv6 地址
- ether_addr：以太网（Ethernet）地址
- inet_proto：网络协议
- inet_service：网络服务
- mark：标记类型

创建一个空的命名集合：`nft add set inet my_table my_set { type ipve_addr \;}`

在添加规则时引用集合，可以使用 `@` 符号跟上集合的名字

`nft insert rule inet my_table my_filter_chain ip saddr @my_set drop` 表示将集合内的 ip 访问全部丢弃

向集合添加元素：`nft add element inet my_table my_set { 10.10.10.22, 10.10.10.33 }`

想在集合内添加一个区间，必须加上一个 flag interval，内核必须提前确认该集合存储的数据类型

创建一个支持区间的命名集合：

```
nft add set inet my_table my_range_set { type ipv4_addr \; flags interval }
```


date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: tcpdump 和 wireshark 使用
---
## tcpdump 和 wireshark 的区别

- tcpdump 仅支持命令行，常用在 Linux 服务器中抓取和分析网络包
- wireshark 除了抓包，还支持可视化分析

### tcpdump 常用参数

```shell
tcpdump -s 0 -i eth1 icmp and host 183.232.231.174 -nn
# -s 0 表示抓取全部数据字节，tcpdump默认只抓取68字节
# -i eth1 表示抓取 eth1 网口的数据包
# icmp 表示抓取 icmp 协议的数据包
# host 表示主机过滤，抓取对应的 ip 数据包
# -nn 表示不解析 ip 地址和端口号的名称
```

| 选项                         | 示例                                       | 说明                                          |
| ---------------------------- | ------------------------------------------ | --------------------------------------------- |
| -i                           | tcpdump -i eth0                            | 指定网络接口，默认是 0 号接口，any 表示所有接口 |
| -nn                          | tcpdump -nn                                | 不解析 ip 地址和端口号的名称                    |
| -c                           | tcpdump -c 5                               | 限制抓取的网络包个数                          |
| -w                           | tcpdump -w file.pcap                       | 保存到文件中，文件名通常以 .pcap 为后缀       |
| host、src host、dst host     | tcpdump -nn host 192.168.1.100             | 主机过滤                                      |
| port、src port、dst port     | tcpdump -nn port 80                        | 端口过滤                                      |
| ip、ip6、arp、tcp、udp、icmp | tcpdump -nn tcp                            | 协议过滤                                      |
| and、or、not                 | tcpdump -nn host 192.168.1.100 and port 80 | 逻辑表达式                                    |
| tcp[tcoflages]               | tcpdump -nn "tcp[tcpflags] & tcp-syn != 0" | 特定状态的 tcp 包                               |


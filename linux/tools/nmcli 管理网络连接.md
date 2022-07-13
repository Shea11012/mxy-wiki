---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: nmcli
---
# nmcli
| 命令                                                     | 说明                                         |
| -------------------------------------------------------- | -------------------------------------------- |
| nmcli con (connection) show                              | 显示所有连接                                 |
| nmcli con show -active                                   | 显示所有的活动连接状态                       |
| nmcli con show "ens33"                                   | 显示网络连接配置                             |
| nmcli con reload                                         | 重新加载配置                                 |
| nmcli con up {ens33}                                     | 启用 ens33 的配置                            |
| nmcli con down {ens33}                                   | 禁用 ens33 的配置                            |
| nmcli dev (device) status                                | 显示设备状态                                 |
| nmcli dev show ens33                                     | 显示网络接口属性                             |
| nmcli con add con-name default type Ethernet ifname eth0 | 创建新连接配置 default，ip 通过 dhcp 自动获取 |
| nmcli con delete default                                 | 删除连接                                     |

- nmcli con add con-name test2 ipv4.method manual ifname ens33 autoconnect no type Ethernet ipv4.addresses 192.168.1.12/24 gw4 192.168.1.1

  参数说明：

  - con-name：指定连接名字
  - ipv4.method：指定获取 IP 地址方式
  - ifname：指定网卡设备名
  - autoconnect：指定是否自动启动
  - ipv4.addresses：指定 ipv4 地址
  - gw4：指定网关

- nmcli con modify test2 connection.autoconnect yes：修改 test2 为自动启动

- nmcli con modify test2 ipv4.dns 114.114.114.114：修改 DNS

- nmcli con modify test2 +ipv4.dns 114.114.114.114：添加 DNS，删除则是将 + 改为 - 
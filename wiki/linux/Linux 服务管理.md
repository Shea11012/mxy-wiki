---
title: "Linux 服务管理"
created: 2026-07-22
updated: 2026-07-22
status: needs-review
type: synthesis
tags:
  - linux
  - linux/service
confidence: INFERRED
---
# Linux 服务管理

> 本页由重复旧笔记整合而来；正文按原页面分节保留，尚未重新核验技术时效性。

## SysV 与 xinetd 服务模型

在 linux 中，服务分为独立守护进程和超级守护进程。

独立守护进程是自行监听在后台的，基本上所有的服务都是独立守护进程类的服务。

超级守护进程专指 xinted 服务，这个服务代为管理一些特殊的服务，这类服务在被请求时才会由 xinted 通知它启动服务，服务提供完毕后就关闭服务，这类服务被称为瞬时守护进程。（超级守护进程 xinted 本身是一个常驻内存的独立守护进程，因为要监听来自外界对其管理的瞬时守护进程的请求，一般不工作时，xinted 不占用端口号，在工作时占用被请求的瞬时守护进程的端口号，并处于监听状态）。

### 管理服务的开机自启动

chkconfig ：管理 /etc/init.d 目录下存在且脚本的内容满足一定条件的服务。

要让 chkconfig 管理服务的开机是否自启动行为，只需将脚本放在 /etc/init.d 目录下，然后再脚本的前部加上 chkconfig 行和 description 行 ：

```
#!/bin/bash

# chkconfig: - 85 15
# description: The Nginx HTTP Server is an efficieent and extensible
```

chkconfig 和 description 两行必须得被注释，且必须在所有非注释行的前面。

chkconfig ： `-` 表示适用于运行级别 123456 上，85 表示开机启动时的顺序，15 表示关机停止服务时的顺序

description ： 可以随便给一点描述信息，但必须得给 description: 关键字

### 管理 xinted 及相关瞬时守护进程

该类服务不能直接使用 service 命令来启动，需要去 /etc/xinted.d/ 目录或 /etc/xinetd.conf 中配置

## systemd 与 systemctl

| 命令 | 动作 | 备注 |
| --- | --- | --- |
| systemctl |展示系统状态 | |
| systemctl list-units | 列出所有正在运行的 units | --type: 列出指定类型，大致有 service、device、socket、mount,--state：列出指定状态，有 active、load、sub |
| systemctl --failed | 列出失败的 units | |
| systemctl list-unit-files | 列出所有已安装的 unit files | |
| systemctl status pid | 展示一个 pid 的进程状态 | |
| systemctl status unit | 展示一个 unit 的状态 | |
| systemctl is-enabled unit | 检查一个 unit 是否启动 | |
| systemctl start unit | 立即启动一个 unit | |
| systemctl stop unit | 立即停止一个 unit | |
| systemctl restart unit | 立即重启一个 unit | |
| systemctl reload unit | 立即重载一个 unit 和它的配置 | |
| systemctl daemon-reload unit | 重载 systemd manager 配置 | |
| systemctl enable unit | 开机自启该 unit | |
| systemctl enable --now unit | 开机自启且立即启动该 unit | |
| systemctl disable unit | 关闭开机自启 | |
| systemctl reenable unit | 禁用和重启 | |
| systemctl mask unit | 使用掩码使 unit 不可被启动 | |
| systemctl unmask unit | 恢复 | |
| systemctl reboot | 关机且重启 | |
| systemctl poweroff | 关闭系统 | |
| systemctl suspend | 挂起系统| |
| systemctl hibernate | 休眠 | |
| systemctl hybrid-sleep | 休眠或者挂起 | |

### unit file 位置

`/usr/lib/systemd/system` : 一般是安装包写入 unit files 的位置
`/etc/systemd/system` : 由系统管理员写入
`.config/systemd/user`: 用户级别的 systemd

## journalctl 日志查询

## 查询启动时间节点

```bash
journalctl --list-boots
```

## 从指定时间查询日志

```bash
journalctl -S "xxxx-xx-xx xx:xx:xx" -U ""
```

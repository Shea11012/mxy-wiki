---
date created: 2021-09-11 14:54
date modified: 2021-09-24 23:32
title: systemd
---
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
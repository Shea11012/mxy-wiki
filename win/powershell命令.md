---
date created: 2022-03-09 19:03
date modified: 2022-07-22 03:11
title: powershell命令
---

## 定时关机

```shell
# 10点关机，只执行一次的定时任务
schtasks /create /sc once /tr "shutdown -s -t 0" /rn "shutdown" /st 10:00
```
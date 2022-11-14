---
date created: 2022-03-09 19:03
date modified: 2022-11-07 22:59
title: powershell命令
---

## 定时关机

```shell
# 10点关机，只执行一次的定时任务
schtasks /create /sc once /tr "shutdown -s -t 0" /rn "shutdown" /st 10:00
```

## 卸载win store里的APP
```shell
# 查询要卸载的app
Get-AppxPackage | Where-Object Name -like "*terminal*" | Select Name

# 卸载方法1：
Get-AppxPackage -Name "将上面命令输出的app名字放入" | Remove-AppxPackage

# 卸载方法2：
Get-AppxPackage | Where-Object Name -like "*terminal*" | Remove-AppxPackage

# 为所有用户都写在，必须在管理员模式
Get-AppxPackage -Name "{app-name}" -AllUsers | Remove-AppxPackage -AllUsers
```
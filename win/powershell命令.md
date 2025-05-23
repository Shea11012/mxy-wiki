---
tags: ['powershell']
date created: 2022-03-09 19:03
date modified: 2022-12-18 12:50
title: powershell命令
updated: 2025-03-25
---

## 定时关机

```shell
# 10点关机，只执行一次的定时任务
schtasks /create /sc once /tr "shutdown -s -t 0" /tn "shutdown" /st 10:00
```

## 卸载 win store 里的 APP

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

## 安装 powershell7

```shell
# 进入管理员shell
iex "& { $(irm https://aka.ms/install-powershell.ps1) } -UseMSI"
```

## 针对 PATH 变量的操作

### 格式化输出

```powershell
$env:PATH -split ';'
```

## sudo

 在 PowerShell 中使用 sudo 的方法
 ```powershell
sudo pwsh -c '这里填写需要sudo权限执行的命令'
```
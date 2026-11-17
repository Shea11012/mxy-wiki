---
title: "Arch Linux命令"
created: 2026-01-13
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - linux
  - linux/tool
confidence: INFERRED
---
#command #tools 

# 不常用，但有用
## qdbus 命令

### qdbus 显示 kwin

可以显示每个应用所用的合成器
```shell
qdbus6 org.kde.KWin /KWin org.kde.KWin.showDebugConsole
```

## 查看当前机器上的连接设备
为了使用makima 做remap 需要查看当前的设备
 `sudo udevadm monitor` 监控设备插拔情况，此时就可以看到是哪个设备被卸载和添加，根据输出信息，再通过 `sudo evtest` 可以测试具体是在哪个`/dev/input/event*` 上监听（有可能会是一个外接扩展器或者是一个无线接收器，所以需要通过evtest确定且获得该设备名
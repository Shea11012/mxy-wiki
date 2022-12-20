---
date created: 2021-11-30 21:22
date modified: 2022-06-27 02:47
title: win scoop 使用
tags: ["win","scoop"]
---
安装 scoop

```shell
# 指定安装目录后，所有安装的工具都会在指定目录下
$env:SCOOP = "custom_dir"
iwr -useb get.scoop.sh | iex
```

首次安装时检查环境配置

```shell
scoop checkup
```

配置 scoop 代理 

```shell
scoop 仅支持http代理
scoop config proxy 127.0.0.1:38002
```

自用 scoop 工具

```shell
scoop install 7zip dark innounp sudo sfsu
scoop install git jetbrainsMono-NF dbeaver
scoop install everything copyq obsidian quicklook snipaste sumatrapdf
```


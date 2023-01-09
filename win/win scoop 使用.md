---
date created: 2021-11-30 21:22
date modified: 2022-12-21 03:51
title: win scoop 使用
tags: ["win", "scoop"]
---

安装 scoop

首先设置 SCOOP 的环境变量

```shell
# 指定安装目录后，所有安装的工具都会在指定目录下

Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
iwr -useb get.scoop.sh | iex
```

首次安装时检查环境配置，根据提示配置环境

```shell
scoop checkup
```

配置 scoop 代理 

```shell
scoop 仅支持http代理
scoop config proxy 127.0.0.1:38002
```

添加源

```
scoop bucket add apps https://github.com/kkzzhizhou/scoop-apps
scoop update
```

安装必备工具

```shell
scoop install git # 安装git会自动安装7z
scoop install dark innounp sudo scoop-search
sudo scoop install -g jetbrainsMono-NF jetbrainsMono
```


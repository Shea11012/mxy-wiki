---
tags: []
date created: 2021-09-01 22:08
date modified: 2023-02-18 02:56
title: paru 常用命令
---

# paru 常用命令

## 安装包

- paru -S 包名：可以同时安装多个以空格分隔
- paru -Sy 包名：在同步包数据库后再执行安装
- paru -Sv 包名：显示一些操作信息后执行安装
- paru -U： 安装本地包，扩展名为 pkg.tar.gz
- paru -U url：安装一个远程包，如：`paru -U http://www.example.com/repo/example.pkg.tar.xz`
- paru -Sw 包名：只下载包，不安装

## 删除包

- paru -R 包名：只删除包，保留其全部已经安装的依赖关系
- paru -Rnc 包名：删除包时，删除所有依赖和构建文件
- paru -Rd 包名：删除包时不检查依赖

## 搜索包

- paru -Ss 关键字：在仓库搜索含关键字的包
- paru -Qs 关键字：搜索已安装的包
- paru -Qi 包名：产看包的详细信息
- paru -QI 包名：列出包文件
- paru -Qii 包名：查找出哪些软件依赖此包

## 其他

- paru -Sc：清理未安装的包文件，包文件位于 `/var/cache/paru/pkg`
- paru -Scc：清理所有缓存文件
- paru -Syy：仅同步软件包数据库
- paru -Syyu：同步数据库且更新已安装包
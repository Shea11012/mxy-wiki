---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: docker基础镜像tag的含义
---
发行版：

`stretch` 与 `jessie` 代表使用的是 `Debian` 前者表示 `Debian 9` ，后者表示 `Debian 8`

`alpine` 代表的是另一种 Linux 发行版。

`slim-jessie` 是使用 docker-slim 压缩后的 `jessie` 镜像。

#### Alpine

是 Linux 发行版的一种，这种镜像只有 5MB 大小，也有自己的包管理系统（apk），适用于做基础镜像，用于生产环境的应用。

#### docker-slim

slim-xxx ，是使用了 docker-slim 压缩后的镜像，是无损压缩
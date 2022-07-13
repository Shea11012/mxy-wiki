---
date created: 2021-09-01 22:08
date modified: 2022-03-09 19:04
title: 使用docker中遇到的问题
---
**linux 中无法使用 host.docker.internal 访问宿主机**
解决：
- 使用 docker 命令行增加 `--add-host=host.docker.internal:host-gateway`
- 使用 docker-compose 文件，给对应的 service 增加以下

```yaml
extra-hosts:
   - "host.docker.internal:host-gateway"
```
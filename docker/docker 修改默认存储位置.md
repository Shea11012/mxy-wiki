---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: docker 修改默认存储位置
---
默认情况下 docker 存放位置为：`/var/lib/docker`

### 软连接

首先停掉 docker 服务

`systemctl restart docker 或 service docker stop`

然后移动整个 `/var/lib/docker` 目录到指定路径

`mv /var/lib/docker /home/shea/docker`

`ln -s /home/shea/docker /var/lib/docker`

重新启动 docker ，`docker info` 查看存储目录，这时存储目录发生变化

### 修改或新建 daemon.json 文件

```json
vim /etc/docker/daemon.json
{
   	"graph":"/home/shea/docker"
}
```

此方法修改立即生效，不需重启服务

### system 下创建配置文件

在 `/etc/systemd/system/docker.service.d` 目录下创建一个文件 `docker.conf` 默认 `docker.service.d` 文件夹不存在

定义新的存储位置

```shell
vim /etc/systemd/system/docker.service.d/docker.conf
[Service]
ExecStart=/usr/bin/dockerd --graph="/home/shea/dockerd" --storage-driver=devicemapper
```

devicemapper 是当前 docker 所使用的存储驱动，查看当前 docker 使用的存储驱动 `docker info` 

```shell
systemctl daemon-reload
systemctl start docker
```


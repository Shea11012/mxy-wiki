---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: Ubuntu docker 安装和卸载
---
# Ubuntu docker 安装和卸载

[官方安装文档](https://docs.docker.com/install/linux/docker-ce/ubuntu/#prerequisites)

如果有旧版本先卸载旧版本

`sudo apt-get remove docker docker-engine docker.io`

使用存储库安装

设置存储库

1. 更新 `apt` 源 	

   `sudo apt-get update`

2. 安装包允许 apt 通过 HTTPS 使用存储库

   ```shell
   sudo apt-get install \
       apt-transport-https \
       ca-certificates \
       curl \
       software-properties-common
   ```

   添加 docker 官方 GPG 密钥
   `curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -`

3. 设置稳定的存储库

   ```shell
   sudo add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) \
      stable"
   ```

     stable 可以改成 edge、test

4. 再更新 apt 软件包源
   `sudo atp-get update`

   1. 安装最新版的 docker CE
      `sudo apt-get install docker-ce`

   2. 或者安装指定版本的 docker CE ,先列出可用的版本

      `apt-cache madison docker-ce`
      `sudo apt-get install docker-ce=version` 例如 docker-ce=18.03.0.ce

   3. 测试 docker 是否安装成功
      `sudo docker run hello-world`

卸载 docker CE

1. docker CE 软件包卸载

   `sudo apt-get purge docker-ce`

2. 不会自动删除主机上的图像，容器，卷或自定义配置文件。删除所有图像，容器和卷

   `sudo rm -rf /var/lib/docker`

   需要自己删除配置文件

​       

安装 docker-compose 

```shell
sudo curl -L https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
```

   下载的版本自己可以选择 [版本](https://github.com/docker/compose/releases)

分配权限

```
sudo chmod +x /usr/local/bin/docker-compose
```


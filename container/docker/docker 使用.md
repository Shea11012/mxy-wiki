---
date created: 2021-11-30 21:22
date modified: 2021-12-14 21:54
title: docker 使用
---
# docker 使用

更改镜像源，下面两个镜像源提供了更换方法

> DaoCloud: https://www.daocloud.io/mirror
> AliCloud: https://cr.console.aliyun.com/#/accelerator

[Docker 快速上手](https://segmentfault.com/a/1190000008822648)

### 镜像与容器
| 命令                                                | 说明                                                         |
| --------------------------------------------------- | ------------------------------------------------------------ |
| docker search $mirror-name                          | 搜索镜像                                                     |
| docker pull mirror-name: \$tag                       | 获取镜像，$tag 为镜像版本，省略则下载最新即:latest（若下载的镜像携带有版本标签，之后使用该镜像都需要携带版本标签，否则会因版本不同再次下载 |
| docker images                                       | 查看所有可用镜像                                             |
| docker images $mirror-name                          | 查看单个镜像                                                 |
| docker rmi $mirror-name                             | 删除镜像                                                     |
| docker ps                                           | 查看已启动的容器                                             |
| docker ps -a                                        | 查看全部容器                                                 |
| docker logs $container-name                         | 查看容器日志                                                 |
| docker run $mirror-name:tag                         | 生成容器，如果宿主机中没有该镜像则会先进行下载，注意镜像标签是否正确 |
| docker run --name \$container-name \$mirror-name     | 给\$mirror-name 指定一个名字为container-name                  |
| docker start $container-name                        | 启动容器                                                     |
| docker stop $container-name                         | 关闭容器                                                     |
| docker run -it \$mirror-name                         | 以交互方式创建容器，也可以在 $mirror-name 的前面加上 --name 来指定容器的名字 |
| docker run -d \$mirror-name                          | 以后台运行方式创建容器                                       |
| docker inspect \$container-name                      | 查看容器详情                                                 |
| docker rm \$container-name-1 ,[container-name-2] ... | 删除容器 (可以加入一个或多个批量删除)                         |
| docker rm \$(docker ps -aq)                          | 删除全部容器                                                 |
### docker 网络

| 命令                                                           | 说明                                                                                                                                                                                                                 |
| -------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| docker network ls                                              | 查看网络类别（none[表示容器为独立个体，不与外部通信],host[表示容器与宿主机（安装 docker 的机器）共享网络],bridge[容器默认网络类型，网桥模式意味着容器间可以互相通信，对外通信需要借助宿主机，通常表现为端口号的映射]） |
| docker network inspect \$network-name                          | json 格式，看到使用当前网络类型的容器列表，只会显示已经启动的容器                                                                                                                                                     |
| docker network create --driver \$network-driver \$network-name | 创建网络，自定义网络类别和网络名 ，省略 --driver \$network-driver 就默认创建 bridge                                                                                                                                     |
| docker network connect \$network-name \$container-name         | 将容器\$container-name 加入 network-name 网络中，一个容器可以属于多个网络                                                                                                                                               |
| docker network disconnect $network-name container-name         | 将容器从某一网络中移除                                                                                                                                                                                               |
| docker run --network \$network-name mirror-name                | 在容器生成时指定网络                                                                                                                                                                                                 |
| docker network rm \$network-name                               | 删除网络                                                                                                                                                                                                             |

测试网络联通情况
1. 在容器交互模式中使用 `ip addr`
2. 使用 `docker inspect $container-name`
3. 使用 `docker inspect $container-name | grep IPAddress`
  之后可以使用 ping 指令测试容器间的网络连通情况。没有 ping 指令的容器需要安装 iputils

### Docker 存储
讲容器的部分存储映射到宿主机器中，以便对配置文件、日志文件、数据文件等进行备份

| 命令                                                            | 说明                                 |
| --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- | 
| docker run --volume /my/mac/dir:/container/dir \$mirror-name    | 在创建时，将系统的某一目录指定为容器某一目录的数据卷（此时容器内部的/container/dir 与宿主机的/my/mac/dir 形成|
| docker run --volume /my/mac/file:/container/file \$mirror-name  | 将文件与文件映射起来    | 
| docker run --volume /container/dir \$mirror-name | 在指定数据卷时，可以省略宿主机目录，此时 Docker 会自动指定一个主机空间映射 |
| docker run --volume /my/mac/dir:/container/dir:ro \$mirror-name | 可以选择只读方式，这样文件或目录的修改就只能在宿主机中进行。添加 `:ro`|     
| docker volume ls                                                | 查看数据卷  |
| docker volume ls -f dangling=true                               | 当容器被删除，主机上的数据卷不会被删除，可以通过以下指令查看那些没有容器使用的数据卷，注意，这里只会显示由 Docker 自动指定的数据卷 |
| docker rm -v \$container-name                                   | 如果需要在删除容器时一并删除数据卷                                                                                               |
| docker run --volume-from \$container mirror-name                | 在创建容器时，选择该容器的数据卷与之前的某容器相同，比如在面对多容器共享项目目录空间这一需求时                                   |
| docker volume rm \$volume-id                                    | 数据卷删除                                                                                                                       |

可以通过 `docker inspect $container-name`,在 Mounts 中看到数据卷的详细情况

### Docker 端口

| 命令 | 说明 |
|----|----|
| docker run -p \$host-port : \$container-port nginx | \$host-port(宿主机端口) $container-port(容器端口) 将主机端口绑定到容器端口，例：12345:80,这样对 localhost:12345 就相当于访问容器 80 端口 |
| docker port \$container-name | 查看指定的容器端口 |

`docker ps -a` 中的 `PORTS` 栏看到端口映射情况。注意只有处于运行中的容器才会有实际的端口映射。



### docker-compose

docker-compose 依赖一个 docker-compose.yml 文件，用以指定容器数据卷、网络等

.yml 文件遵循 YAML 语法，详细可见 [YAML 语言教程](http://www.ruanyifeng.com/blog/2016/07/yaml.html?f=tt)

样例

```yaml
version: '2.0'
services: 
  # 启用一个镜像为 nginx 的容器并命名为 web1
  web1:
    image: nginx
    # 开启 80 和 443 端口，实际映射端口由 Docker 指定
    ports:
      - "80"
      - "443"
    # 将该容器加入 mynetwork 中
    networks:
      - "mynetwork"
    # 指定该容器要在 web2 容器启动之后启动，且在其停止前停止
    depends_on:
      - web2
  web2:
    image: nginx
    ports:
      - "33333:80"
    networks:
      - "mynetwork"
      - "bridge"
    volumes:
      - "/mnt"
networks:
  # 创建一个驱动为 bridge 的网络，命名为 mynetwork
  mynetwork:
    driver: bridge
```

| 命令 | 说明 |
|-----|----|
| docker-compose up | 需要一个名为 docker-compose.yml 或 docker-compose.yaml 文件，进入该文件所在目录，通过指令生成容器 |
| docker-compose up -d | 后台运行 |
| docker-compose stop | 停止容器 |
|docker-compose start | 运行容器 |
| docker-compose rm | 删除容器，但不会删除之前创建的网络 |
| docker-compose down | 既删除容器，又删除网络，数据卷的删除任需要使用删除数据卷的方法 |

### 生成 & 提交镜像

|命令 | 说明|
|----|----|
|docker commit -m $commit-msg -a author container-id namespace/mirror-name:tag| 生成镜像，例：docker commit -m 'install nginx' -a 'shea' abcd1234 shea/nginx:test |
| docker push $namespace/mirror-name:tag | 提交镜像 |
国内使用 DaoCloud 或阿里 Docker 服务，具体推送方法去服务商网站查看

### Dockerfile
Dockerfile 可以指定新镜像的原镜像来源，对原镜像的操作、环境变量、以及以此创建容器时执行的指令等。
样例
```dockerfile
# 新镜像基于的原镜像
FROM centos:centos6.8

# 指明维护者
MAINTAINER dailybird <dailybird@mail.com>

# 设置一些环境变量，使用 \ 表示连接多个设置
ENV NGINX_VERSION 1.11.11 \
    TEST_ENV hello

# 指定暴露的端口号，
EXPOSE 80 443

# 在原镜像基础上进行的修改
RUN yum install -y wget iputils \
    && wget http://nginx.org/download/nginx-1.11.11.tar.gz

# 以此镜像创建并启动时，容器执行的指令，通常用于启动服务
CMD ["echo", "hello world"]
```
比如使用以下配置可以在 centos 中安装 nginx：

```dockerfile
FROM centos:centos6.8

MAINTAINER dailybird <dailybird@mail.com>

EXPOSE 80 443

RUN cd / \
    && mkdir data \
    && cd data \
    && mkdir nmp \
    && cd nmp \
    && yum install -y wget pcre-devel gcc gcc-c++ \
       ncurses-devel perl make zlib zlib-devel \
       openssl openssl--devel iputils \
    && wget http://nginx.org/download/nginx-1.11.11.tar.gz \
    && tar zxf nginx-1.11.11.tar.gz \
    && cd nginx-1.11.11 \
    && ./configure --prefix=/usr/local/nginx \
    && make && make install && make clean \
```

更多的 Dockerfile 指令可以参考 [Dockerfile 指令](http://wiki.jikexueyuan.com/project/docker-technology-and-combat/instructions.html)

**Dockerfile 的使用**

新建一个名为 Dockerfile 的文件（没有后缀），并写入一些配置内容，然后再该文件的目录中，通过以下指令创建镜像：

`docker build --tag $namespace/$mirror-name:$tag $dockerfile-dir`

其中，$dockerfile-dir 为 Dockerfile 所在目录，比如执行：

`docker build --tag shea/nginx-demo:demo ./`

**docker-compose 中使用 Dockerfile**

当我们需要启动一个新镜像时，可以先将此镜像创建出来，然后在 `docker-compose.yml` 文件中通过 `image` 指定新镜像；也可以直接通过以下方式将这两个步骤合并： 

```yaml
version: '2.0'
services: 
  web1:
    # build 后的参数为 Dockerfile 文件所在的目录位置，替换原先的 image
    build: ./
    ports:
      - "80"
    networks:
      - "mynetwork"
      
# ...
# 其他配置
```

此后，可以通过以下指令创建容器：

`docker-compose build`

`docker-compose up`

或者直接执行：

`docker-compose up --build`


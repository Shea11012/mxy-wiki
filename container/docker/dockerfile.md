---
tags: []
date created: 2023-04-10 19:15
date modified: 2023-04-10 19:58
---

## dockerfile 指令

### CMD

格式：
- shell 格式：`CMD 命令`
- exec 格式 (推荐)：`CMD ["可执行文件", "参数1"， "参数2"...]`
指定默认的容器主进程的启动命令，在运行时可以指定新的命令来替代镜像中的默认命令

### ENTRYPOINT

格式与 CMD 一样，当指定了 ENTRYPOINT 后，CMD 就不再是直接的运行命令，而是将 CMD 的内容当作参数传给 ENTRYPOINT，变为
```
<ENTRYPOINT> "<CMD>"
```

### HEALTHCHECK

格式：
- `HEALTHCHECK [选项] CMD <命令>`：设置检查容器健康状态的命令
- `HEALTHCHECK NONE`：如果基础镜像有健康检查指令，就可以屏蔽掉其健康检查指令

选项：
- `--interval=xx`：两次健康检查的间隔，默认为 30 秒
- `--timeout=xx`：健康检查命令运行超时时间，默认 30 秒
- `--retries=x`：当连续失败指定次数后，则将容器状态视为 unhealthy，默认 3 次

### ONBUILD

ONBUILD 是一个特殊的指令，后面跟其他指令，如 run、copy 等。这些指令在当前镜像构建时并不会被执行。只有以当前镜像为基础镜像，去构建下一级镜像时才会被执行。

```dockerfile
FROM node:slim
RUN mkdir /app
WORKDIR /app
ONBUILD COPY ./package.json /app
ONBUILD RUN [ "npm", "install" ]
ONBUILD COPY . /app/
CMD [ "npm", "start" ]
```
当其他镜像以此镜像作为基础镜像时，ONBUILD 就会开始执行

### SHELL

可以指定 RUN、ENTRYPOINT、CMD 指令的 shell，Linux 中默认为 `["/bin/sh","-c"]`
```dockerfile
SHELL ["/bin/sh","-c"]

RUN ls

SHELL ["/bin/sh","-cex"]

RUN ls
```

## dockerfile 多阶段构建

```dockerfile
FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /go/src/github.com/go/helloworld/

RUN go get -d -v github.com/go-sql-driver/mysql

COPY app.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /go/src/github.com/go/helloworld/app .

CMD ["./app"]
```

### 只构建某一阶段的镜像

`docker build --target builder -t username/imagename:tag .`

### 构建时从其他镜像复制文件

`COPY --from=nginx:latest /etc/nginix/nginx.conf /nginx.conf`


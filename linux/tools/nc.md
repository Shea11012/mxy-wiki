---
tags: ["nc", "netcat"]
date created: 2022-12-21 18:20
date modified: 2022-12-21 18:48
---

参数：

- z 使用 0 io，连接成功后立即关闭，不进行数据交换
- v 详细输出
- n 不使用 DNS 反向查询 IP 地址域名

## 端口扫描
```bash
# 范围扫描
nc -zvn ip port-port

# 扫描指定端口
nc -zv ip port
```

## 文件传输
```bash
# server 将file输出到1567端口
nc -l 1567 < file.txt

# client 接受1567端口的内容
nc -n server_ip 1567 > file.txt
```

## 目录传输
```bash
# server 压缩目录然后将内容发送到1567端口
tar -cvf - dir_name | nc -l 1567

# client 接受端口内容，然后解压
nc -n server_ip 1567 | tar -xvf -
```


## 加密发送
```bash
# server 服务端发送加密后的数据
nc -n client_ip 1567 | mcrypt -flush -bare -F -q -d -m ecb > file.txt

# client 客户端监听1567，并接受加密数据
mcrypt -flush -bare -F -q -m ecb < file.txt | nc -l 1567
```

## 克隆磁盘

>[!tips]
>/dev/sda 在没有进行分区时，使用/dev/sda，如果进行了分区，则用 /dev/sda/1 之类

```bash
# server
dd if=/dev/sda | nc -l 1567

# client
nc -n server_ip 1567 | dd of=/dev/sda
```

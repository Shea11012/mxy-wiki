---
tags:
  - tools
  - linux
date created: 2023-11-02 16:53
date modified: 2023-11-02 16:55
---
linux 上一切皆文件，通过 `lsof` 可以找出文件被哪个进程打开或占用。

# 常用示例
## 列出特定用户打开的文件

```bash
lsof -u {user}
```

## 查找特定端口运行进程

```bash
lsof -i TCP:53
```

## 列出ipv4和ipv6文件

```bash
lsof -i 4

lsof -i 6
```

## 列出TCP端口范围
```bash
lsof -i TCP:1-1024
```

## 指定pid搜索
```bash
lsof -p {pid}
```

## 杀掉用户的所有进程
```bash
kill -9 `lsof -t -u {user}`
```


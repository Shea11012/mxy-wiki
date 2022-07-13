---
tags: ["linux","rsync"]
date created: 2021-06-07 23:12
date modified: 2022-04-16 00:18
title: rsync
---

## rsync
```bash
rsync /etc/test.conf /tmp   # 在本地同步
rsync -r /etc host:/tmp     # 将本地 /etc 目录拷贝到远程主机 /tmp 下，以保证 /tmp 目录和本地 /etc 保持同步
rsync -r host:/etc /tmp     # 将远程主机 /etc 目录拷贝到本地 /tmp 下，以保证本地 /tmp 目录和远程 /etc 同步
rsync /etc/ # 列出本地 /etc/ 目录下的文件列表
rsync host:/tmp/ # 列出远程主机上 /tmp/ 目录下的文件列表
rsync -R -r /var/./log/project /tmp # 使用一个点代表相对路径的起始位置，该示例表示 /tmp/log/project

```

> rsync 如果需要拷贝一个目录不需要 /，如果仅需要拷贝目录下的内容，则需要带上斜杠
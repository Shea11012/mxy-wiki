---
created: 2025-06-28
tags: [tools,backup]
updated: 2025-06-28
---

搭配 **resticprofile** 使用

![[../../assets/eef47c3a201662749b00aed608d82d29_MD5.jpeg]]
有一个这样的 snapshots，当在还原时，会出现一些问题
```bash
restic restore -t d:/tmp latest --include /c/Users/18723/workspace/mxy/mxy-wiki

# 会导致如下情况
d:/tmp/c/Users/18723/workspace/mxy/mxy-wiki
```

解决
```bash
restic restore -t d:/tmp latest:/c/Users/18723/workspace/mxy --include /mxy-wiki
```
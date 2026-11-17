---
created: 2025-11-19
tags: []
---
# 概念
jj 有change id 和 commit id
[[attachments/e76f175a0e14ed7916022931dc3f5edb_MD5.jpeg|Open: Pasted image 20251120002537.png]]
![[attachments/e76f175a0e14ed7916022931dc3f5edb_MD5.jpeg]]
change id是不会变化得id，commit id随着每次改动都会变化

jj 每次在执行命令时都会检测文件的改动创建新的commit id

## 命令
`jj git init --colocate`：已存在的git仓库
`jj bookmark track main@origin`：追踪指定远程分支
`jj new main`：从指定的bookmark创建新的工作区，如果不指定main，则是默认在@上创建一个新的gg
`jj split`：可以从多个被提交的文件中选取指定的作为另一个提交
`jj commit`：是 new + describe 的结合
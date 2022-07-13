---
date created: 2021-11-28 13:34
date modified: 2021-11-28 13:34
title: 基础shell编程
---
## 基础 shell 编程
### 重定向输入和输出
`>` 输出重定向且覆写文件
`>>` 输出重定向，输出追加至文件

### 输入重定向
`<`
**内联重定向** 
```shell
wc <<EOF
test string 1
test string 2
test string 3
EOF
```

### 管道
将一个命令的输出作为另一个命令输入，叫做管道 (piping)。
由管道串起来的命令，会有 linux 同时运行这两个命令
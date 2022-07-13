---
date created: 2021-11-28 13:34
date modified: 2021-12-03 20:20
title: shell
---
# shell
[[linux/shell/shell 基础]]

[[变量]]

[[linux/shell/Linux 文件权限]]

[基础shell编程](基础shell编程.md)

[[linux/shell/结构化命令]]

[[变量]]

[[linux/shell/获取用户输入和控制输出]]

## 创建临时文件
### 创建本地临时文件
mktemp 会在 /tmp 目录创建一个唯一的临时文件。
mktemp 需要指定一个文件名模板，只需要在文件名后面加 2 个以上 X。
```shell
mktemp testing.XXXXX
```

-d 创建临时目录 `mktemp -d dir.XXXXX`

## 控制脚本
### 捕获信号
```shell
trap commands signals
```

example
```shell
trap "echo ' Sorry! I have trapped Ctrl-C'" SIGINT
#
echo This is a test script
#
count=1
while [ $count -le 10 ]
do
    echo "Loop #$count"
    sleep 1
    count=$[ $count + 1 ]
done
#
echo "This is the end of the test script" #
```

### 删除已设置好的捕获
```shell
trap -- SIGINT
```
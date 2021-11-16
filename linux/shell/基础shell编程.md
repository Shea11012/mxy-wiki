## 基础shell编程
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
将一个命令的输出作为另一个命令输入，叫做管道(piping)。
由管道串起来的命令，会有linux同时运行这两个命令
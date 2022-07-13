---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: sed
---
## sed 编辑器
sed 编辑器根据命令来处理数据流中的数据，sed 编辑器会执行一下操作：
- 一次从输入中读取一行数据
- 根据所提供的编辑器命令匹配数据
- 按照命令修改流中的数据
- 将新的数据输出到 STDOUT

```bash
sed option script file
```
可用选项：
- -e script 在处理输入时，将 script 中指定的命令添加到已有的命令中
- -f file 在处理输入时，将 file 中指定的命令添加到已有的命令中
- -n 不产生命令输出，使用 print 命令来完成输出

### 使用地址
在 sed 编辑器中有两种形式的行寻址
- 数字形式表示行区间
- 文本模式过滤行

数字格式：
```bash
[address]command
```

```bash
address {
	command1
	command2
	command3
}
```

文本格式：
```bash
/pattern/command
```

### 命令：
s (substitute)：替换命令 `s/pattern/replacment/flags`
d (delete)：删除命令
i (insert)：插入命令
a (append)：追加命令
c (change)：修改命令，修改数据流中整行文本的内容
=：打印行号
l (list)：列出行
p (print)：打印文本行
y (transform)：转换命令，`[address]y/inchars/outchars/`，转换命令会对 inchars 和 outchars 进行一对一的映射。inchars 与 outchars 必须相等。
w (write)：向文件写入行，`[address]w filename`
r (read): 将一个独立文件读入到数据流中，`[address]r filename`

> i a 命令不能在地址区间中使用

#### 替换标记
- 数字：表明新文本将替换第几处模式匹配的地方
- g，表明新文本将替换所有匹配的文本
- p,表明打印匹配的内容
- w file,将替换的结构写到文件中

### sed 保持空间

**模式空间：sed 编辑器的工作空间，是一块活跃的缓冲区**

**保持空间：在处理模式空间中的某些行时，可以用保持空间来临时保存一些行**

**保持空间的命令**

| 命令 | 描述 |
| --- | --- |
| h | 将模式空间复制到保持空间 |
| H | 将模式空间附加到保持空间 |
| g | 将保持空间复制到模式空间 |
| G | 将保持空间附加到模式空间 |
| x | 交换模式空间和保持空间的内容 |

### 分支

分支命令格式：`[address]b label`

**如果没有定义 label 标签，跳转命令会跳转到脚本尾部。指定了标签则可以跳过地址匹配出的命令**
```shell
 cat data2.txt  
 This is the header 1 line.  
 This is the first data line.  
 This is the second data line.  
 This is the last line.  

 sed '{2,3b ; s/This is/Is this/;s/line./test?/}' data2.txt  
 # 此处没有定义 label 标签会跳转到脚本尾部，相当于不会对 2 3 行执行脚本  

 sed '{/first/b jump1;s/This is the/No jump on/; :jump1;s/This is the/Jump here on/}' data2.txt  
 # 当一行内匹配到 first 字符，就跳转到 jump1，这样就会跳过 jump1 前面的脚本命令
 ```

#### 测试

测试命令会根据替换命令的结果跳转到某个标签，如果替换命令成功匹配并替换了一个模式，测试命令就会跳转到指定的标签。如果未指定标签，则跳转到脚本结尾。
```shell
 echo "This, is, a, test, to, remove, commas." | sed -n '{:start;s/,//1p; t start}'  
 # t 命令相当于 if 判断，如果 t 成功则跳转到指定标签，否则不跳转
```
### 模式替代

& 符号可以用来代表替换命令中的匹配模式

```shell
 echo "The cat sleeps in his hat." | sed 's/.at/"&"/g'  
 # 因为点符号代表的是未知字符，无法使用 .at 这样的字符替换，sed 提供了 & 符号来代表匹配到的整个字符  
 # 也可以使用正则中的组匹配模式 (.at)，需要注意 sed 中使用组匹配模式需要给括号加转义符 \(.at\)
```

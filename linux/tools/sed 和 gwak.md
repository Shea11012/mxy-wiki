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
y (transform)：转换命令，`[address]y/inchars/outchars/`，转换命令会对 inchars 和 outchars 进行一对一的映射。inchars与outchars 必须相等。
w (write)：向文件写入行，`[address]w filename`
r (read): 将一个独立文件读入到数据流中，`[address]r filename`

> i a 命令不能在地址区间中使用

#### 替换标记
- 数字：表明新文本将替换第几处模式匹配的地方
- g，表明新文本将替换所有匹配的文本
- p,表明打印匹配的内容
- w file,将替换的结构写到文件中


## gawk
```bash
gawk options program file
```
命令行可用选项：
- -F fs 指定行中划分数据字段的字段分隔符
- -f file 从指定的文件中读取程序
- -v var=value 定义 gawk 程序中的一个变量及其默认值
- -mf N 指定要处理的数据文件中的最大字段数
- -mr N 指定数据文件中的最大数据行数
- -W keyword 指定 gawk 的兼容模式或警告等级

### 数据字段变量
- $0 代表整个文本行
- $1 代表文本行中的第一个数据字段
- $2 代表文本行中的第二个数据字段
- $n 代表文本行中的第n个数据字段

### 预处理关键字
BEGIN: 在处理数据之前会先执行，`gawk 'BEGIN {print "hellow world!"}'`
END：读完数据之后执行，`gawk 'END {print "End of File"}'`
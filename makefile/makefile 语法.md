---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: MakeFile
---
# MakeFile

```makefile
edit: main.o display.o
		cc -o edit main.o dispaly.o

main.o: main.c defs.h
		cc -c main.c

display.o: display.c defs.h
		cc -c display.c
```



make 命令执行流程：

1. make 会在当前目录下找名字叫 Makefile 或 makefile 文件
2. 找到文件后，会找文件中的第一个目标文件，并把这个文件作为最终的目标文件
3. 如果目标文件不存在，或者依赖的 `.o` 文件修改时间要比目标文件新，则会执行目标文件后面所定义的命令来生成目标文件
4. 如果依赖的 `.o` 文件不存在，make 会在当前文件中查找目标为 `.0` 文件的依赖性，在根据规则生成 `.o` 文件
5. 如果 c 文件和 H 文件存在，make 会生成 `.o` 文件，再用 `.o` 文件生成 make 的目标文件

声明变量：

```makefile
objects = main.o xx.o files.o

edit: $(objects)
		cc -o edit $(objects)
```

## 书写规则

### 文件搜寻

`VPATH` 如果没有指定此变量，make 只会在当前的目录中寻找依赖文件和目标文件，如果定义了则会在当前目录中寻找不到的情况下，会到所指定的目录中找寻文件。

```makefile
VPATH = src:../headers
```

 另一种设置文件搜索路径的方法是使用 vpath 关键字：

- `vpath pattern directories`  ：为符合模式 pattern 的文件指定搜索目录 directories，如果定义了重复的 patter 会按照先后顺序来执行搜索
- `vpath pattern`：清除符合模式 pattern 的文件搜索目录
- `vpath` 清除所有已被设置好了的文件搜索目录

### 伪目标

创建 Makefile 时应该为所有自定义的命令设置伪目标，这样可以避免有同名的文件在项目中

伪目标不是一个文件，只是一个标签。伪目标一般没有依赖的文件，但是我们可以为伪目标指定所依赖的文件。伪目标同样可以作为默认目标，只要将其放在第一个。

```makefile
.PHONY:	* # 可以使用通配符
.PHONY: $(MAKECMDGOALS) # 自动变量
.PHONY: $(filter-out vendor node_modules,$(MAKECMDGOALS)) # 过滤掉vendor和node_modules,使用自动变量给定的伪目标

prog1: prog1.o utils.o
	cc -o prog1 prog1.o utils.o

prog2: prog2.o
  cc -o prog2 prog2.o
  
prog3: prog3.o sort.o utils.o
  cc -o prog3 prog3.o sort.o utils.o
```

### 前置条件 Prerequisite

当 vendor 命令不是伪目标时，在 composer.json 或 composer.lock 发生改动时将会执行 composer install，所以在伪目标一行将 vendor 排除

```makefile
.PHONY: $(filter-out vendor,$(MAKECMDGOALS))
vendor: composer.json composer.lock
	@composer install
```



## 书写命令

make 会按顺讯一条一条的执行命令，每条命令的开头必须以 tab 键开头，除非命令是紧跟在依赖规则后面的分号后的。命令行之间的空格或空行会被忽略，如果该空格或空行是一个 tab 键开头，那么 make 会认为其是一个空命令。

### @

`@` 字符在命令行前，那么这个命令将不被 make 显示出来

```makefile
@echo 正在编译 xx 模块
# 输出
正在编译 xx 模块

#如果不使用 @则输出
echo 正在编译 xx 模块
```

`make -n` ：只显示命令，但不执行命令

### 命令执行

如果要让上一条命令的结果应用在下一条命令时，应该使用分号分隔这两条命令

```makefile
exec:
	cd /home/xxx
	pwd	 #输出当前 makefile 的目录
	
exec:
	cd /home/xx; pwd # 输出 /home/xx
```

### 命令出错

每条命令执行成功，make 都会检测其返回码，如果某条命令退出码非零，那么 make 就会终止执行当前的规则。如使用 mkdir 如果目录不存在则出错，使用 mkdir 是为了保证存在这一目录，此时的 mkdir 出错的命令可以忽略。忽略命令的出错可以在 makefile 命令行前加一个 `-` ，标记不管命令出不出错都认为是成功

```makefile
clean:
	-rm -f *.o
```

### 嵌套执行 make

传递变量到下级的 makefile 中：`export variable`

不让某些变量传递到下级：`unexport variable`

如果需要传递所有的变量，只要一个 `export` 。有两个变量 `SHELL,MAKEFLAGS` 这是两个系统级的环境变量，其总是要传递到下层的 makefile 中， 

### 定义命令包

makefile 中出现一些相同命令序列，那么可以为这些相同的命令序列定义一个变量，以 `define` 开头，以 `endef` 结束

```makefile
define run-yacc
yacc $(firstword $^)
mv y.tab.c $@
endef
```

## 使用变量

### 传递变量

make foo arg=ceshi

```makefile
foo:
	@echo $(arg) # 会输出 ceshi
```



变量不能包含 `:,#,= 空格,回车` ，变量大小敏感

申明完变量后，需要在给变量名前加上 `$` 符号，使用 `() 或 {}` 将变量名包裹起来

```makefile
objects = program.0 foo.o utils.o
program: $(objects)
```

操作符 `?=`，

```makefile
FOO ?= bar # 其含义是如果 FOO 没有被定义过，那么 FOO 的值就是 bar，否则什么也不做
```

### 变量的高级用法

`$(var:a=b)`：把变量 var 中所有以 a 字符串结尾的 a 替换成 b 字符串

```makefile
foo := a.o b.o c.o
bar := $(foo:.o=.c) # a.c b.c c.c
```

### override 指示符

如果有变量是通过 make 命令行设置的，那么 makefile 中对这个变量的赋值就会被忽略，语法是：

```makefile
override variable := value
```

## 条件判断

```makefile
ifeq ($(CC),gcc) # 判断 cc 是否是 gcc
	$(CC) -o foo $(objects) $(libs_for_gcc)
else
	$(CC) -o foo $(objects) $(normal_libs)
endif
```

`ifneq` : 判断值是否不同，不同为真

`ifdef` ：判断值是否非空 ，非空为真

`ifndef` ：判断值是否为空，为空为真

## 函数

### 字符串函数

```makefile
# subst 字符串替换函数，将字串 text 中的 from 替换成 to
$(subst from,to,text) 

# patsubst 模式字符串替换函数，查找 text 中的单词，是否符合模式 pattern，如果匹配则以 replacement 替换
$(patsubst pattern,replacement,text)

# strip 去掉空格，去掉头尾空字符
$(strip string)

# findstring 查找字符串函数，在字串 in 中查找 find 字串
$(findstring find,in)

# filter 过滤函数，以 pattern 模式过滤 text 字符串，保留符合 pattern 模式的单词
$(filter pattern,text)

# filter-out 反过滤函数，以 pattern 模式过滤 text 字符串中的单词，去除符合模式 pattern 的单词
$(filter-out pattern,text)

# sort 排序函数，给字符串 list 中的单词排序（升序），会去掉重复的单词
$(sort list)

# word 取单词函数，取字符串 text 中的第 n 个单词
$(word n,text)

# wordlist 取单词串函数，从字符串 text 中取从 ss 开始到 e 的单词串，ss 和 e 是数字
$(wordlist ss,e,text)

# words 单词个数统计，统计 text 中字符串中的单词个数
$(words text)

# firstword 首单词函数，取字符串 text 中的第一个单词
$(firstword text)
```

### 文件名操作函数

```makefile
# dir 取目录函数，从文件名序列 name 中取出目录部分
$(dir names)

# notdir 取文件函数，从文件名序列 names 中取出非目录部分
$(notdir names)

# suffix 去后缀函数，从文件名序列 names 中取出各个文件名的后缀
$(suffix names)

# basename 去前缀函数，从文件名序列 names 中取出各个文件名的前缀部分
$(basename names)

# addsuffix 加后缀函数，把后缀 suffix 加到 names 中的每个单词后面
$(addsuffix suffix,name)

# addprefix 加前缀函数，把前缀 prefix 加到 names 中的每个单词后面
$(addprefix prefix,names)

# join 连接函数，把 list2 中的单词对应的加到 list1 的单词后面，如果 list1 的单词个数要比 list2 多，那么list1 多出来的单词将保持原样，如果 list2 的单词个数要比 list1 多，则 list2 多出来的单词要被复制到 list1 中
$(join list1,list2)
```

### foreach 函数

```makefile
# foreach 把参数 list 中的单词逐一取出放到参数 var 中所指定的变量中，然后再执行 text 所包含的表达式
$(foreach var,list,text)
```

### if 函数

```makefile
# if condition，then-part
or
# if condition,then-part,else-part
```

### call 函数

```makefile
$(call expression,parm1,parm2,...,parm)
```

### origin 函数

```makefile
# 只是告诉这个变量从哪里来
$(origin variable)
```

### shell 函数

```makefile
# 使用操作系统的 shell 命令，该函数会生成一个新的 shell 程序来执行命令，需要注意执行性能
$(shell cat foo)
stop:
	# shell内的变量只能通过 $$var 的方式访问，因为是shell内的变量不是Makefile的变量
	@pid=$(shell ps -ef | grep $(execDir)/$(execName) | grep -v "grep" | awk '{print $$2}'); \
	# ;\ 使用该方式才能在下一行共享变量
	echo $$pid
```

### 控制 make 函数

```makefile
# error text 是错误信息，error 函数不会在一被使用就会产生错误，如果把其定义在某个变量中，并在后续的脚本中使用这个变量，也是可以的
$(error text)
ERR := $(error found an error)
.PHONY: err
err: $(ERR)

# warning 不会使得 make 退出，只是输出警告信息
$(warning text)
```


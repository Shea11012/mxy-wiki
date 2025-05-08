---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: template 用法
updated: 2025-04-07
---

## 两类 template：text、HTML

. 的作用域：表示当前作用域的对象

去除空白：

```go
{{- 23}} // 表示去除左边的空白
{{23 -}} // 表示去除右边的空白
```

注释：

```go
{{/* xxxxxx */}}
{{- /* xxxxx */ -}} // 注释后的内容不会被引擎替换，注释行在替换的时候也会占用行，所以应该去除前缀和后缀空白，否则会多一空行。
```

管道 pipeline：

**| 将前面命令的运算结果传递给后一个命令的最后一个位置**

```go
{{"Hello World!" | prinf "%s,%s\n" "abcd"}}
// 并非只有 | 才是 pipeline，在 Go template 中， pipeline 的概念是传递数据，只要能产生数据都是 pipeline
{{println (len "output")}}
```

变量：

```go
{{$v := "Hello"}}
{{- println $v}}
```

条件判断：

**pipeline 为 false 的情况是各种数据对象的零值：0、空指针或接口，数组、slice、map、string 长度为 0**

```go
{{if pipeline}} T1 {{else}} T0 {{end}}
{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
```

循环：

**pipeline 的值必须是：array、slice、map、channel**

```go
{{range pipeline}} T1 {{end}}
{{range pipeline}} T1 {{else}} T0 {{end}}

// 可以输出两个参数，当只接受一个值时则是val
{{ range $key,$val := . }}
{{ . }}
{{end}}
```

with...end：

```go
// 设置 . 的值
{{with pipeline}} T1 {{end}} //当 pipeline 的输出不是空时，才会执行 T1
{{with pipeline}} T1 {{else}} T0 {{end}}
```

内置函数：

```go

call // 第一个参数作为函数调用，剩余作为参数，函数必须只能有一个或两个返回，如果有第二个返回值，则必须是 error 类型
	call .X.Y 1 2 // 表示调用 .X.Y(1,2)
html // 返回转义后的 HTML，这个函数在 html/template 不可用
index // 对可索引对象取值
	index x 1 2 3 // x[1][2][3]
js // 返回转义后的 js
len // 返回参数长度
print // fmt.Sprint 别名
printf  // fmt.Sprintf 别名
println // fmt.Sprintln 别名
urlquery // 返回转义后的 urlquery，在 html/template 不可用
```

布尔运算：
```go
and // 返回第一个为空的参数或最后一个参数，可以任意多个参数
	and .x .y  // 等于 if x then y else x

or // 与 and 用法 类似
	or .x .y  // 等于 if x then x else y

not .x // 对参数取反，只接受一个参数
```

比较函数：

```go
eq // arg1 == arg2
ne // arg1 != arg2
lt // arg1 < arg2
le // arg1 <= arg2
gt // arg1 > arg2
ge // arg1 >= arg2
```

嵌套 template：

**模板之间的变量不会继承**

```go
// 注意这种方式，必须先定义个模板
{{template "name"}} // 直接执行 . 设置为 nil
{{template "name" pipeline}} // 直接执行 . 设置为 pipeline 生成的数据

{{define "name"}} T1 {{end}} // 先定义一个模板
{{template "name" pipeline}} // 再根据名字查找执行定义的模板

{{block "name" pipeline}} T1 {{end}} // 相当于上面的集合体，定义模板的同时直接执行
```



使用：

```go
t := template.Must(template.New("hh").Parse(` // 使用 Must 方法可以不用处理 New 的 err
"{{$}}"
"{{$.Name}}"
{{$.School.Name | printf "%q"}}
"{{$.School.Room}}"`))
if err := t.Execute(os.Stdout, p); err != nil {
    log.Fatalln(err)
}

t,err := template.New("xx").Parse(`xxx`)
if err != nil {
    log.Fatal(err)
}

if err:= t.Execute(os.Stdout,p);err != nil {
    log.Fatal(err)
}
```


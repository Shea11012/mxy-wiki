---
created: 2024-06-28
tags: [linux,font]
updated: 2024-06-28 01:39
---

# 字体分类

- monospace 等宽
字符宽度相同的字体，用于需要字符严格对齐的场景。用于终端和代码。中文字体不存在等宽与比例差别，所有的中文字都是等宽的。

- sans-serif 无衬线
笔画末端没有修饰的字体，用于屏幕显示

- serif 有衬线
笔画末端有修饰字体，用于打印

# fontconfig

## 通用字族名

sans-serif(sans)，serif 和 monospace(mono)。它们不是真实存在的字体，是让程序使用无衬线、衬线、等宽字体。程序会通过查 fontconfig 找到匹配的字体。

## 使用

### FC_DEBUG

fontconfig 就会打印调试信息
```shell
FC_DEBUG=4 {APP}
```
打印的信息分为：
add rule，已添加的配置文件规则

FcConfigSubstitue Pattern，包含 font pattern。s 和 w 分别代表强弱绑定，prgname 代表程序名
FcConfigSubstitue editPattern，对 font pattern 的替换操作，只有规则匹配时才执行 editPattern

FcConfigSubstitue donePattern，这是 fontconfig 执行完字体替换后的结果

### fc-list

根据 charset 找到哪些字体包含该 Unicode
```shell
fc-list ":charset=30edd" file family
```

>[!tips]
>wezterm ls-fonts --text " 𰻝𰻝面 " 也可以找到哪些字体可以正确显示

# 参考

- [[https://catcat.cc/post/2021-03-07/|fontconfig治理Linux中的字体]]
- [[https://catcat.cc/post/2020-10-31/|fontconfig字体匹配机制]]
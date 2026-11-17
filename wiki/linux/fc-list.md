---
title: "fc-list"
created: 2026-01-13
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - linux
  - linux/tool
confidence: INFERRED
---
# fc-list 常用命令

可以通过 `:` 分隔多个条件，组合过滤字体

只输出 font family
```bash
fc-list : family
```

输出指定语言的 family
```bash
fc-list :lang=zh family
```

按字体样式过滤
```bash
fc-list : style="Bold"
```

按字体文件过滤
```bash
fc-list : file="/usr/share/fonts/truetype/xx.ttf"
```

按字重过滤
```bash
fc-list : weight=200
```

按字距过滤
```bash
fc-list : spacing=100
```

显示指定字段
```bash
fc-list : family file style
```

显示字体编码范围
```bash
fc-list : encoding
```

显示字体分辨率
```bash
fc-list : dpi
```

根据码点筛选字体
```bash
fc-match -s ":charset=codepoint" --format='%{family}\n'
```
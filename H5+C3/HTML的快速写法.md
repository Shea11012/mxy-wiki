---
created: 2021-12-03 20:20
updated: 2024-07-22
title: HTML的快速写法
---

# HTML 的快速写法

1. E 代表 HTML 标签
2. E#id 代表 id 属性
3. E.class 代表 class 属性
4. `E[attrA=foo attrB]` 自定义属性
5. 
6. E{foo} 代表标签包含的内容是 foo
7. E>N：N 是 E 的子级元素
8. E+N 代表 N 是 E 的同级元素
9. `a>b^c`：c 在 b 的父级元素生成
10. \$：迭代数字，`E.class$` 生成class1、class2类推，E.class{$}在E元素内生成数字
11. `()` 分组：`div>(header>ul>li*2>a)+footer>p`

```html
连写（E*N）
li*3>a：
<li><a href=""></a></li>
<li><a href=""></a></li>
<li><a href=""></a></li>
自动编号（E$*N）
div#item$.class$$*3
<div id="item1" class="class01"></div>
<div id="item2" class="class02"></div>
<div id="item3" class="class03"></div>
```

```html
nav>ul>(li>a[href=#]{Link})*5
```


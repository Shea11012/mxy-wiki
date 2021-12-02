# HTML的快速写法

1. E 代表HTML标签
2. E#id 代表id属性
3. E.class 代表class属性
4. E[attr=foo] 代表某一个特定属性
5. E{foo} 代表标签包含的内容是foo
6. E>N 代表N是E的子元素
7. E+N  代表N是E的同级元素
8. E^N  代表N是E的上级元素

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


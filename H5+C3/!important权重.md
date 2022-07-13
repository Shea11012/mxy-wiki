---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: !important css 权重
---
# !important css 权重

尽量不要使用权重 `!important`

```html
<style>
   p { color:red !important; } //增加了权重，显示red，不加显示blue
   #test p { color: blue;}
</style>

<div id='test'>
    <p class='test1'>
        测试权重
    </p>
    <p class='test2'>
        测试权重
    </p>
</div>
```

### css 优先级顺序，从小到大

1. 通用选择器：`*`
2. 元素（类型）选择器：`span，p，form....`
3. 类选择器：`.class`
4. 属性选择器：`div[id],input[name]....`
5. 伪类：`:focus,:link`，伪类选择器，以 `:` 开头
6. ID 选择器：`#id`
7. 内联样式
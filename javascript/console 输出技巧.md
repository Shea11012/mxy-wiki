---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: console 输出技巧
---
### log 输出格式化

| 说明符   | 作用                   |
| -------- | ---------------------- |
| %s       | 元素转换为字符串       |
| %d 或 %i | 元素转换为整数         |
| %f       | 元素转换为浮点数       |
| %o 或 %O | 元素以最有效的格式显示 |
| %c       | 应用提供的 css         |

```js
console.log('%c message','font-size:36px;font-weight:bold')
```

### 输出嵌套的对象的方法

```js
console.log(JSON.stringify(obj,null,2)) // stringify 第三个参数表示设置空格缩进大小
console.dir(obj,{depth:null}); // 第二个参数表示记录对象的深度
```


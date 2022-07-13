---
date created: 2021-12-03 20:20
date modified: 2021-12-26 03:05
title: 原生JS DOM操作方法
---
### 创建节点

| 方法                                | 说明                                    |
| ----------------------------------- | --------------------------------------- |
| `document.createElement(tagname)`   | 创建一个由 `tagName` 决定的 `HTML` 元素 |
| `document.createTextNode(data)`     | 创建一个文本节点，文本内容为 `data`     |
| `document.createDocumentFragment()` | 创建一个空白的文档片段                  |

### 通过节点关系获取节点

| 方法                         | 说明                                                         |
| ---------------------------- | ------------------------------------------------------------ |
| `Node.parentNode`            | 返回指定节点在 `DOM` 树中的父节点                            |
| `Node.parentElement`         | 返回指定节点在 `DOM` 树中的父元素节点，没有父元素或者不是一个元素节点，则返回 `null` |
| `Node.childNodes`            | 返回指定节点的所有子元素的集合，包括文本节点                 |
| `Node.children`              | 返回指定节点的所有子元素的集合，只包含元素节点               |
| `Node.nextSibing`            | 返回指定节点的下一个兄弟节点，包括文本节点                   |
| `Node.nextElementSibing`     | 返回指定节点的下一个兄弟元素节点                             |
| `Node.previousSibing`        | 返回指定节点的上一个兄弟节点，包括文本节点                   |
| `Node.previousElementSibing` | 返回指定节点的上一个兄弟节点 元素节点                        |
| `Node.firstChild`            | 返回指定节点的第一个子节点，包括文本节点                     |
| `Node.firstElementChild`     | 返回指定节点的第一个子元素节点                               |
| `Node.lastChild`             | 返回指定节点的最后一个子节点，包括文本节点                   |
| `Node.lastElementChild`      | 返回指定节点的最后一个子元素节点                             |

### 节点操作

| 方法                  | 说明                                                         |
| --------------------- | ------------------------------------------------------------ |
| `Node.appendChild()`  | 将一个节点添加到指定节点的子节点列表                         |
| `Node.removeChild()`  | 将一个节点从 `DOM` 树中移除，移除后还存在内存中，可以使用变量接收 |
| `Node.insestBefore()` | 在当前节点的某个子节点之前再插入一个子节点                   |
| `Node.replaceChild()` | 用指定节点替换当前节点的一个子节点，返回被替换掉的节点       |

### 节点选择

| 方法                                   | 说明                                                         |
| -------------------------------------- | ------------------------------------------------------------ |
| `document.querySelector(selectors)`    | `selectors` 是一个字符串，包含一个或多个 `css` 选择器，返回获取到的元素 |
| `document.querySelectorAll(selectors)` | 和 `querySelector` 类似，返回值为 `NodeList` 对象            |
| `document.getElementById()`            | 根据元素 `ID` 获取元素                                       |
| `document.getElementsByTagName()`      | 根据元素标签名获取元素，返回值为 `HTMLCollection` 集合       |
| `document.getElementsByName()`         | 根据元素 `name` 属性获取元素，返回值为 `NodeList` 对象       |
| `document.getElementsByClassname()`    | 根据元素类名获取元素，返回值为 `HTMLCollection` 集合         |

### 属性操作

| 方法                                   | 说明                                                         |
| -------------------------------------- | ------------------------------------------------------------ |
| `element.setAttribute(属性名，属性值)` | 给元素设置属性，如果已存在                                   |
| `element.removeAttribute(属性名)`      | 删除元素的某个属性                                           |
| `element.getAttribute(属性名)`         | 获取元素上的属性名为 `attrName` 的属性值，如不存在则返回 null 或者 空字符串 |
| `element.hasAttribute(属性名)`         | 检测该元素上是否有该属性，返回值为 true 或 false             |

### DOM 事件

| 方法                                                   | 说明                                                         |
| ------------------------------------------------------ | ------------------------------------------------------------ |
| `element.addEventListener(type.listener,[options])`    | 给元素添加指定事件 type 以及响应该事件的回调函数             |
| `element.removeEventListener(type,listener,[options])` | 移除元素上指定事件，如果元素上分别在捕获和冒泡阶段都注册了事件，则需要分别移除 |
| `document.createEvent()`                               | 创建一个自定义事件，随后必须使用 init 进行初始化             |
| `element.dispathEvent(event)`                          | 对指定元素触发一个事件                                       |

### 元素样式尺寸

| 方法                            | ` 说明 `                                                       |
| ------------------------------- | ------------------------------------------------------------ |
| `window.getComputedStyle(elem)` | 获取 elem 所有应用了 `css` 后的属性值，返回一个实时的 `CSSStyleDeclaration` 对象 |
| `elem.getBoundingClientRect()`  | 返回元素的大小以及相对于视口的位置，返回一个 `DOMRect` 对象，包括元素的 `left,right,top,bottom,width,height,x,y` 属性值 |


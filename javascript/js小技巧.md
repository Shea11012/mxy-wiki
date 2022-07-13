---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: js小技巧
---
#### 复制到剪贴板

```javascript
awaitContent.on('click','button',function () {	// jq on 只能绑定页面已有的元素并且是需要绑定元素的父元素
		let input = this.previousSibling.previousSibling;
		input.setAttribute('readonly','readonly');
		input.focus();
		input.setSelectionRange(0,input.value.length);
		if (document.execCommand('copy')) {
			document.execCommand('copy');
			layer.msg('复制成功');	// 这是layer 的弹出层控件
		}
    });
```

#### 返回上一页并刷新

```javascript
self.location = document.referrer;
```

#### 闭包

闭包优点可以缓存数据，缺点也是缓存数据，会延长局部变量的作用域链

```javascript
function f1(){
    let num = parseInt(Math.random()*10);
    return function(){
        return num;
    }
}
f1();
f1();
f1();
// 此时调用的三次函数产生的num是一样的
```

### 使用 jquery 实现 textarea 自适应

```javascript
$('textarea').each(function () {
            this.setAttribute('style','height:'+(this.scrollHeight) + 'px;overflow-y:hidden;');
        }).on('input',function () {
            this.style.height = 'auto';
            this.style.height = (this.scrollHeight) + 'px';
        });
```

### JavaScript 获取 cookie 值 

**推荐使用 js-cookie 库**

```javascript
function getCookie(cname)
{
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++)
    {
        var c = ca[i].trim();
        if (c.indexOf(name)==0) return c.substring(name.length,c.length);
    }
    return "";
}
```

### 检查数据类型

在使用 **typeof** 判断数据类型时，没有办法判断复杂数据类型只能判断基础数据类型

```js
function checkType(data) {
  return Object.prototype.toString.call(data).slice(8,-1)
}
```


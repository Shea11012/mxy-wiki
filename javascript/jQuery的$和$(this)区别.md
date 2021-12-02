jQuery 中 this 与 $(this) 的区别

```js
$("#textbox").hover(
    function(){
        this.title = 'Test';
    }),
    function(){
    this.title = 'ok';
}
```

这里的 this 是一个 HTML 元素（textbox），具有 text 属性。

```js
错误示范
$("#textbox").hover(
    function(){
        $(this).title = 'Test';
    }),
    function(){
    $(this).title = 'ok';
}
```

`$()` 是 jQuery 的核心基础函数，功能：

- `$('#id')` ：传入一个选择器字符串，获得这个选择器对应的 dom 内容，保存在数组中
- `$(function(){....})` ：传入一个匿名函数
- `$(this)` ：将 JavaScript 对象包装成 jQuery 对象

```js
$('#id').click(function(){
    this.css('display','block');	// 报错，this是一个HTML元素，不是jQuery对象，因此this不能调用jQuery
    $(this).css();			// 成功	$(this)是一个jQuery对象，不是html元素，可以使用CSS方法
    this.style.display = 'block';	//成功 this是一个html元素，不是jQuery对象，不可以调用jQuery的方法，但是可以使用JavaScript来更改style属性
})
```


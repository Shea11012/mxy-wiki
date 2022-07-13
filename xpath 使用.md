---
date created: 2021-11-30 21:22
date modified: 2021-11-30 21:22
title: xpath 使用
---
```xpath
//tagname[@attribute='value']
```

- `//`：相对节点
- tagname：特定节点的标记名
- @：选择属性
- attribute：节点的属性名称
- value：属性的值 

### 方法

#### contains()

```
//*[contains(@name,'btn')]
```

#### or and

```
//*[@type='submit' or @name='btnReset']
//*[@type='submit' and @name='btnLogin']
```

#### starts-with()

```
//tagname[starts-with(@id,'message')]
```

#### text()

```
//td[text()='UserID']	// 文本匹配
//td[@id='name']/text() // 输出id匹配name 的文本
```

#### following

```
//*[@type='text']//following:input // 匹配 type=text 的 input 节点
//*[@type='text']//following:input[1] // 根据下标选择指定的节点
```

#### ancestor

```
//*[text()='Enterprise Testing']//ancestor::div	//获取所有与文本内容相匹配的祖先div节点，也可以根据索引选取指定节点，同 following
```

#### child

```
//*[@id='golang']/child::li[1] // 匹配 id=golang 的子节点，根据索引选择指定的节点
```

#### preceding

```
//*[@type='submit']//preceding::input // 匹配 type=submit 之前的所有 input 节点
```

#### following-sibling

```
//*[@type='submit']//following-sibling::input // 匹配 type=submit 的兄弟节点
```

#### parent

```
//*[@id='rt-feature']//parent::div // 匹配 id=rt-feature 的父节点
```

#### self

```
//*[@type='password']//self::input // 匹配自己
```

#### descendant

```
//*[@id='rt-feature']//descendant::a // 匹配 id=rt-feature 的子孙节点，且子孙为 a
```


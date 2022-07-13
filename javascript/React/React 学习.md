---
date created: 2021-12-03 20:20
date modified: 2022-03-09 19:09
title: React 学习
---
## 虚拟 DOM

- **React 虚拟 DOM：**用 JS 对象来模拟页面上的 DOM 和 DOM 嵌套，为了实现页面中 DOM 元素的高效更新

## Diff 算法


安装 React 包 

`npm install react react-dom --save-dev`

```jsx
import React,{ Component } from 'react' // 创建组件、虚拟 DOM 元素，生命周期
import ReactDOM from 'react-dom' // 把创建好的组件和虚拟 DOM 放到页面展示
```

也可以使用 `npm install -g create-react-app` 来创建 react 项目（推荐）[使用文档](https://facebook.github.io/create-react-app/docs/getting-started)

## React 生命周期

[React 生命周期介绍](https://juejin.im/post/5c9b57d65188251d081cba4a)

#### mounting 组件的挂载

- constructor 只适用 class component
  - 在构造函数中，可以初始化 state、设定默认值
- getDerivedStateFromProps 渲染前调用的最后一个方法，仅用于设置状态
  - 常见用法在 mount 期间，根据初始 props 返回状态对象
- render 渲染
- componentDidMount 第一次渲染组件之后，触发此方法
  - 加载数据在此处执行
  - 组件渲染完成可以访问 DOM 节点，添加事件侦听器

#### Updating 组件的更新

- getDerivedStateFromProps  //暂时没搞懂这个方法该在何时调用
- shouldComponentUpdate 当一个组件收到新的 state 或 props 时，它应该更新，返回一个布尔值用于控制是否更新
- render
- getSnapshotBeforeUpdate 触发时间在 render 之后，最新的渲染提交给 DOM 之前
- componentDidUpdate 所有的提交已经提交给 DOM

#### unmounting 

- componentWillUnmount 组件注销前，清除事件监听、定时器等，防止内存泄漏
- errors
  - getDerivedStateFromError 捕捉子组件树的错误，根据错误信息修改 state，同时展示出 error
  - componentDidCatch 子组件中发生错误时触发，仅捕获渲染/生命周期方法中的错误。

### JSX 语法

> JSX 语法：符合 xml 规范的 JS 语法

解析 jsx 语法，需要安装 `babel` 插件

`npm install @babel/core babel-loader @babel/preset-env @babel/preset-react --save-dev`

- babel-loader : 使用 Babel 转换 Javascript 依赖关系的 Webpack 加载器
- @babel/core ：将 ES6 代码转换为 ES5
- @babel/preset-env ：根据支持的浏览器，决定使用哪些 transformations/plugins 和 polyfills ，如为旧浏览器提供现代浏览器的新特性
- @babel/preset-react ：针对所有 React 插件的 Babel 预设，如将 JSX 转换为函数

根目录创建 `.babelrc` 文件

```babelrc
	{
        "presets": ["@babel/preset-env","@babel/preset-react"]
	}
```

在 `webpack.config.js` 文件中配置 loader

```js
module.exports = {
    module:{
        rules:[
            {
                test:/\.(js|jsx)$/,
                exclude:/node_modules/,
                use:{
                    loader:'babel-loader'
                }
            }
        ]
    },
    resolve: {
        extensions:['.js','.jsx','.json'],	// 这项配置表示在 import 包时不用写这几个文件名后缀
        alias:{
            '@':path.join(__dirname,'./src')	// 为 @ 符号做了映射，这样可以使 @ 符号表示 src 目录
        }
    }
}
```

`npm install prop-types --save-devs`

- prop-types : 为组件传入参数指定变量类型

```jsx
import React from "react";
import PropTypes from "prop-types";
const Input = ({ label, text, type, id, value, handleChange }) => (
  <div className="form-group">
    <label htmlFor={label}>{text}</label>
    <input
      type={type}
      className="form-control"
      id={id}
      value={value}
      onChange={handleChange}
      required
    />
  </div>
);
Input.propTypes = {
  label: PropTypes.string.isRequired,
  text: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  id: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired
};
export default Input;
```

```jsx
import React, { Component } from "react";
import ReactDOM from "react-dom";
import Input from "../presentational/Input.jsx";
class FormContainer extends Component {
  constructor() {
    super();
    this.state = {
      seo_title: ""
    };
    this.handleChange = this.handleChange.bind(this);
  }
  handleChange(event) {
    this.setState({ [event.target.id]: event.target.value });
  }
  render() {
    const { seo_title } = this.state;
    return (
      <form id="article-form">
        <Input
          text="SEO title"
          label="seo_title"
          type="text"
          id="seo_title"
          value={seo_title}
          handleChange={this.handleChange}
        />
      </form>
    );
  }
}
export default FormContainer;
```
### 创建组件的方式

- 使用 class 关键字创建的组件，被称作『有状态的组件』，有自己的私有数据和生命周期函数
- 使用 function 关键字创建的组件，被称作『无状态的组件』，没有私有数据和生命周期函数

创建组件时

```jsx
import React,{ Component } from 'react'

class Test extends Component {
    render(){
        return (
            // 当需要显示的子组件是同级的时候需要一个父组件包裹起来，此时会有个无用的 div
            <div> // 可以使用 <React.Fragment> 替代，这样就不会出现多余的元素
                <h1>Hello World!</h1>
                <button>按钮</button>
            </div>			</React.Fragment>
        )
    }
}
```

### 组件内绑定 this 的三种方法

```javascript
class Test extends Component {
    constructor(props){
        super(props)
        
        this.increment = this.increment.bind(this) //1
    }
    
    increment = () => { // 2
        
    }
    render(){
        return (
        	<div onClick={this.incrment.bind(this)}/> // 3
        )
    }
}
```



#### props 使用

```react
class Mouse extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            a:'1'
        }
    }
    function t(){
        
    }
    
    render(){
        <p>{this.props.id}</p>    // props 访问组件上赋的值
        {this.t.bind(this)}	// 访问方法则不必
    }
}

<Mouse id='12'></Mouse>
```


---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: webPack学习
---
在根目录创建一个文件名为 `webpack.config.js`

使用 `webpack-dev-server` ，安装此扩展需要安装 `webpack` 扩展

安装命令

```
npm i webpack webpack-dev-server webpack-cli -D
```

`webpack-dev-server` 打包生成的文件，放在了内存中，不会直接生成在物理磁盘中

`webpack-dev-server` 常用参数

```
--open 自动打开
--port xxx 指定端口号
--contentBase path 指定默认打开的目录
--hot 热更新，减少不必要的加载，可以实现浏览器的无刷新加载，针对 js 代码无效果
```

第一种配置方式

```
// package.json
scripts: {
    "dev" : "webpack-dev-server --open --contentBase src --hot"
}
```

第二种配置方式

```
// webpack.config.js
const webpack = require('webpack');
module.exports = {
    devServer : {
        open: true,
        port: 8080,
        contentBase: 'src',
        hot: true,
    },
    plugins:[
        new webpack.HotModuleReplacementPlugin()
    ]
}
```

**`html-webpack-plugin` 在内存中生成 HTML 的插件**

作用：

- 自动在内存中根据指定页面生成一个内存页面
- 自动把打包好的 `js` 文件追加到页面中去

```
const htmlWebpackPlugin = require('html-webpack-plugin');
module.exports = {
    .
    .
    .
    pulgins:[
        new htmlWebpackPlugin({
            template:path.join(__dirname,'./src/index.html'), // 指定模板页面，会根据指定路径在内存中生成
            filenamne:'index.html',	// 指定生成的名称
        })
    ]
}
```

处理 `css` 文件需要安装 `style-loader css-loader less-loader`

```
module.exports = {
    .
    .
    .
    module: { // 用于配置所有第三方模块 加载器
        rules:[ // 第三方模块匹配规则
            { test: /\.css$/,use:['style-loader','css-loader']} // use 使用过程，类似堆栈调用过程，从后往前调，数据一步一步返回
            {test: /\.less$/,use:['style-loader','css-loader','less-loader']}
        ]
    }
}
```

使用 `url-loader file-loader` 对图片路径处理

```
file-loader 是 url-loader 内部依赖
rules:[
    {test:/\.(jpg|png|jpeg|bmp|gif)$/,use:['url-loader?limit=xxx&name=[hash:8][name].[ext]']}
    // 当 limit 参数大于等于图片大小不会被转为 base64 字符串
    hash:8 表示截取 hash 位数
    name ext 保证名字和文件后缀不变
   
   	{test:/\.(ttf|ttof)$/},use:['url-loader']	// 处理字体文件
]
```

安装 babel 将高级语法转换为低级语法

```
安装以下包
npm i babel-core babel-loader babel-plugin-transform-runtime -D
npm i babel-preset-env babel-preset-stage-0 -D
在根目录新建 .babelrc 文件
{test:/\.js$/,use:'babel-loader',exclude:/node_modules/}
// 使用 exclude 派出 node_modules 文件夹
```


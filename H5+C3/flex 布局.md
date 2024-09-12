---
tags: ['css']
created: 2021-12-03 20:20
updated: 2024-07-23
title: flex 布局
---
![[../assets/H5+C3/flex 布局/IMG-20240503150917441.png]]

### 父元素属性

#### justify-content：水平对齐元素

- flex-start：元素和容器的左端对齐
- flex-end：元素和容器的右端对齐
- center：元素在容器里居中
- space-between：元素之间保持相等的距离
- space-around：元素周围保持相等的距离

#### align-items：纵向对齐元素

- flex-start：元素与容器的顶部对齐
- flex-end：元素与容器的底部对齐
- center：元素纵向居中
- baseline：元素在容器的基线位置显示
- stretch：元素被拉伸以填满整个容器

> Tips：当 flex 以列为方向时，justify-content 控制纵向对齐，align-items 控制横向对齐

#### flex-direction 元素在容器里的摆放位置

- row：元素摆放的方向和文字方向一致
- row-reverse：元素摆放的方向和文字方向相反
- column：元素从上放到下
- column-reverse：元素从下放到上

#### flex-wrap 定义 flex 元素必须在单行或自动换行成多行

- nowrap：所有元素在一行
- wrap：元素自动换成多行
- wrap-reverse：元素自动换成逆序的多行

#### flex-flow 是 flex-direction 和 flex-wrap 属性的缩写 

#### 格式为 flex-flow: flex-direction flex-wrap

#### align-content 当交叉轴有多余空间时，对齐容器内的轴线

- flex-start：多行集中在顶部
- flex-end：多行集中在底部
- center：居中
- space-between：行与行之间保持相等距离
- space-around：每行周围保持相等距离
- stretch：每一行都被拉伸以填满容器

### 子元素属性

#### order ：-1(最前) ，0 (默认值)，1(最后)

#### align-self 在交叉轴上对齐一个元素，覆盖 align-items

- flex-start：元素与容器的顶部对齐
- flex-end：元素与容器的底部对齐
- center：元素纵向居中
- baseline：元素在容器的基线位置显示
- stretch：元素被拉伸以填满整个容器
---
title: "HTTP 响应"
created: 2021-12-03
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - php
  - php/framework
  - php/laravel
confidence: INFERRED
---
所有路由和控制器都会返回一个响应并发送给用户的浏览器，Laravel 提供了几种不同的方式来返回响应。

- 返回一个字符串，会自动将字符串转换为一个完整的 HTTP 响应
- 返回一个数组（Eloquent 集合），会自动转为 json 响应
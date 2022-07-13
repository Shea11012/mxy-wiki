---
date created: 2021-12-03 20:20
date modified: 2021-12-03 20:20
title: 跨域资源共享 CORS
---
浏览器将 CORS 请求分为两类：简单请求和非简单请求

> 请求方法如下：
>
> - HEAD
> - GET
> - POST
>
> 请求头信息不超出以下字段：
>
> - Accept
> - Accept-Language
> - Content-Language
> - Last-Event-ID
> - Content-Type：只限于三个值 application/x-www-form-urlencoded multipart/form-data text/plain

如同时不满足上面两个条件，就属于非简单请求



### 简单请求流程

简单请求会在 CORS 请求时，自动在头部添加一个 Origin 字段，服务器会根据 Origin 字段值，判断是否同意这次请求。如果 Origin 指定的源，不在许可范围内，服务器会返回一个正常的 HTTP 响应，如果浏览器发现响应头信息没有 `Access-Control-Allow-Origin` 字段，则会抛出错误。此错误不能通过状态码识别，因为 HTTP 的回应状态码可能是 200。

与 CORS 请求相关的字段都以 `Access-Control-` 开头：

- `Access-Control-Allow-Origin`  ：该字段是必须的。它的值要么是请求时 Origin 字段的值，要么是一个 *，表示接受任意域名的请求
- `Access-Control-Allow-Credentials` ：该字段可选，该值是一个布尔值，表示是否允许发送 Cookie。默认情况下，Cookie 不包括在 CORS 请求之中。该值只能设为 true，如果服务器不要浏览器发送 Cookie ，删除该字段即可。
- `Access-Control-Expose-Headers` ：该字段可选，CORS 请求时，XMLHTTPRequest 对象的 getResponseHeader 方法只能拿到 6 个基本字段：Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma。如果想拿到其他字段，就必须在 Access-Control-Expose-Headers 里面指定
- `withCredentials`：如果要把 Cookie 发到服务器，指定 `Access-Control-Allow-Credentials` 值为 true，另一方面 

```javascript
let xhr = new XMLHttpRequest();
xhr.withCredentials = true; // 有些浏览器省略此设置也会发送 cookie，如果不想发送，则可以显式的关闭
```

**Tips**

> 如果需要发送 cookie Access-Control-Allow-Origin 就不能设为 * ,必须明确指定、与请求网页一致的域名。同时 Cookie 任然遵循同源策略，只有服务器域名设置的 Cookie 才会上传，其他域名的 Cookie 不会上传



### 非简单请求

#### 预检请求

非简单请求是对服务器有特殊要求的请求，如请求方法是 PUT 或 DELETE，或者 Content-Type 字段类型是 application/json

非简单请求会在正式通信之前，增加一次 HTTP 查询请求，称为预检请求。

浏览器会先询问服务器，当前网页所在的域名是否在服务器的许可名单之中，以及可以使用那些 HTTP 动词和头信息字段。只有得到肯定答复，浏览器才会发出正式的 XMLHTTPRequest 请求，否则会报错。

```javascript
var url = 'http://api.alice.com/cors';
var xhr = new XMLHttpRequest();
xhr.open('PUT', url, true);
xhr.setRequestHeader('X-Custom-Header', 'value');
xhr.send();
```

此时，浏览器发现是一个非简单请求，就会自动发出一个预检请求，预检请求头信息如下：

```javascript
OPTIONS /cors HTTP/1.1
Origin: http://api.bob.com
Access-Control-Request-Method: PUT
Access-Control-Request-Headers: X-Custom-Header
Host: api.alice.com
Accept-Language: en-US
Connection: keep-alive
User-Agent: Mozilla/5.0...
```

预检请求头信息的两个特殊字段：

`Access-Control-Request-Method`：该字段必须的，用来列出浏览器的 CORS 请求会用到哪些 HTTP 方法

`Access-Control-Request-Headers` ：该字段是一个逗号分隔的字符串，指定浏览器 CORS 请求会额外发送的头信息字段

#### 预检请求的回应

```javascript
HTTP/1.1 200 OK
Date: Mon, 01 Dec 2008 01:15:39 GMT
Server: Apache/2.0.61 (Unix)
Access-Control-Allow-Origin: http://api.bob.com
Access-Control-Allow-Methods: GET, POST, PUT
Access-Control-Allow-Headers: X-Custom-Header
Content-Type: text/html; charset=utf-8
Content-Encoding: gzip
Content-Length: 0
Keep-Alive: timeout=2, max=100
Connection: Keep-Alive
Content-Type: text/plain
```

`Access-Control-Allow-Methods` ：该字段必须，它的值是逗号分隔字符串，表明服务器支持的所有跨域请求的方法

`Access-Control-Allow-Headers`：如果浏览器请求包括 Access-Control-Request-Headers 字段，则 Access-Control-Allow-Headers 字段是必须的。

`Access-Control-Max-Age` ：该字段可选，用来指定本次预检请求的有效期，单位为秒，即允许缓存该条会议的时间，在此期间，不用发出另一条预检请求。


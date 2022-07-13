---
date created: 2021-12-03 20:20
date modified: 2022-01-23 15:09
title: Nginx location
---
### location

配置语法：

`location [ = | ~ | ~* | ^~ ] uri {...}`

`location @name {...}`

前缀含义：

```
= :  精确匹配
~ ： 大小写敏感
~*： 忽略大小写
^~： 只需匹配 uri 部分
@ ： 内部服务跳转
```

`location /uri/ {} `：表示对 /uri/ 目录及其子目录下的所有文件都匹配。所以 location / {} 的匹配范围是最大的。

`location = /uri/ {}` ：表示只对目录或文件进行匹配，不对目录中的文件和字母进行匹配。所以一般只用来做文件匹配。

`location ~ /uri/ {}` ：表示区分大小写的正则匹配

`location ~* /uri/ {}` ：表示不区分大小写的正则匹配

`location ^~ /uri/ {}`  ： 表示禁用正则匹配，即精确字符串匹配，此时正则中的元字符被解释成普通字符。

匹配优先级规则为：**nginx 先检查 URI 的前缀路径，在这些路径中找到最精确匹配请求 URI 的路径。然后 nginx 按在配置文件中出现的顺序检查正则表达式路径，匹配上某个路径后即停止匹配并使用该路径的配置，否则使用最大前缀匹配的路径的配置。**

使用 = 前缀可以定义 URI 和路径的精确匹配。如果发现匹配，则终止路径查找。如请求 `/ ` 很频繁，定义 `location = /` 可以提高这些请求的处理速度。

示例：

```
location = / {A}
location / {B}
location /documents/ {C}
location ^~ /images/ {D}
location ~* \.(gif|jpg|jpeg)$ {E}
```

请求 `/` 能匹配 A 和 B，但精确匹配为 A

请求 `/index.html` 的前缀 `/` 能匹配 A 和 B，但 A 只能匹配 `/` 自身，因此最终匹配配置 B。（前缀也能匹配 E，但文件吗无法匹配）

请求 `/documnets/document.html` 的前缀能匹配 B 和 C，但 C 更精确，因此匹配配置 C。（前缀也能匹配 E，但文件名无法匹配）

请求 `/images/1.gif` 的前缀能匹配 B、D 和 E，且 D 和 E 都是最长路径匹配，但 `^~` 优先级更高，因此匹配配置 D

请求 `/documents/1.jpg` 的前缀能匹配 B、C，同样也能匹配 E，且 E 比 B 的匹配更精确，因此最终匹配配置 E

**等号优先级最高，非正则匹配，正则匹配（它们之间有位置先后的顺序），优先级最低的是没有使用任何符号的匹配**

`location = url > location ^~ url > location *~|~ url > location url`
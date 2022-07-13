---
date created: 2021-12-03 20:20
date modified: 2022-01-23 15:09
title: Nginx URL重写
---
### Rewrite 规则

rewrite 是使用 nginx 提供的全局变量或自己设置的变量，结合正则表达式和标志位实现 URL 重写以及重定向。rewrite 只能放在 `server{},location{},if{}` 中，并且只能对域名后边的除去传递的参数外的字符串起作用。如：

`http://zhu-zi.top/a/bb/index.php?id=1&age=22` 只对 `a/bb/index.php` 重写。语法 `rewrite regex replacement [flag];` 

`rewrite` 是对同一域名内更改获取资源的路径

`location` 是对一类路径做控制访问或反向代理



### flag 标志位

- `last` ：相当于 Apache 的 [L] 标记，表示完成 rewrite
- `break`：停止执行当前虚拟主机的后续 rewrite 指令集
- `redirect`：返回 302 临时重定向，地址栏会显示跳转后的地址
- `permanent`：返回 301 永久重定向，地址栏会显示跳转后的地址

return 指令无法返回 301，302，因为还必须有重定向的 URL

last 和 break 的区别：

1.  last 一般卸载 server 和 if 中，而 break 一般使用在 location 中
2. last 不终止重写后的 url 匹配，即新的 url 会再从 server 走一遍匹配流程，而 break 终止重写后的匹配
3. break 和 last 都能阻止继续执行后面的 rewrite 指令



### if 指令与全局变量

语法为 `if(condition){...}`

`-f` 和 `!-f` 判断是否存在文件

`-d` 和 `!-d` 判断是否存在目录

`-e` 和 `!-e` 判断是否存在文件或目录

`-x` 和 `!-x` 判断文件是否可执行

```nginx
if ($http_user_agent ~ MSIE) {
    rewrite ^(.*)$ /msie/$1 break;
} //如果UA包含"MSIE"，rewrite请求到/msid/目录下

if ($http_cookie ~* "id=([^;]+)(?:;|$)") {
    set $id $1;
 } //如果cookie匹配正则，设置变量$id等于正则引用部分

if ($request_method = POST) {
    return 405;
} //如果提交方法为POST，则返回状态405（Method not allowed）。return不能返回301,302

if ($slow) {
    limit_rate 10k;
} //限速，$slow可以通过 set 指令设置

if (!-f $request_filename){
    break;
    proxy_pass  http://127.0.0.1; 
} //如果请求的文件名不存在，则反向代理到localhost 。这里的break也是停止rewrite检查

if ($args ~ post=140){
    rewrite ^ http://example.com/ permanent;
} //如果query string中包含"post=140"，永久重定向到example.com

location ~* \.(gif|jpg|png|swf|flv)$ {
    valid_referers none blocked www.jefflei.com www.leizhenfang.com;
    if ($invalid_referer) {
        return 404;
    } //防盗链
}
```



### rewrite 实例

```nginx
http {
    # 定义image日志格式
    log_format imagelog '[$time_local] ' $image_file ' ' $image_type ' ' $body_bytes_sent ' ' $status;
    # 开启重写日志
    rewrite_log on;

    server {
        root /home/www;

        location / {
                # 重写规则信息
                error_log logs/rewrite.log notice; 
                # 注意这里要用‘’单引号引起来，避免{}
                rewrite '^/images/([a-z]{2})/([a-z0-9]{5})/(.*)\.(png|jpg|gif)$' /data?file=$3.$4;
                # 注意不能在上面这条规则后面加上“last”参数，否则下面的set指令不会执行
                set $image_file $3;
                set $image_type $4;
        }

        location /data {
                # 指定针对图片的日志格式，来分析图片类型和大小
                access_log logs/images.log mian;
                root /data/images;
                # 应用前面定义的变量。判断首先文件在不在，再判断目录在不在，如果还不在就跳转到最后一个url里
                try_files /$arg_file /image404.html;
        }
        location = /image404.html {
                # 图片不存在返回特定的信息
                return 404 "image not found\n";
        }
}
```

对形如 `/images/ef/uh7b3/test.png` 的请求，重写到 `/data?file=test.png`，于是匹配到 `location /data`，先看 `/data/images/test.png` 文件存不存在，如果存在则正常响应，如果不存在则重写 tryfiles 到新的 image404 location，直接返回 404 状态码
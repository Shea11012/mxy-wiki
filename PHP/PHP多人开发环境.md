可以在访问网址 `http://wulv.test.com` 在 nginx 获取到 wulv，把 root 设置到 wulv 目录，这样就可以访问到 wulv 这个目录下的代码，nginx 配置如下：

```nginx
set $who www;
if ($http_who != "") {
    set $who $http_who;
}

root /data/gateway/$who/html;
```

参考文章 [PHP多人开发环境原理解析](https://wenda.shukaiming.com/article/433)

5.0.1 以上版本支持自动读取配置文件，需要在 `Application/extra ` 下新建自定义的配置文件。也可以使用 `ENV` 文件读取配置文件。

#### nginx.conf 配置

pathinfo 模式：`https://xxx/index.php/index/Index/index` 

路由模式：

```nginx
# thinkphp.conf
# 重写 URL 隐藏 index.php
location / {
    if (!-e $request_filename) {
        rewrite ^(.*)$ /index.php?s=/$1 last;
        break;
    }
}
location ~ \.php$ {
        fastcgi_pass php-upstream;
        fastcgi_index index.php;
        fastcgi_buffers 16 16k;
        fastcgi_buffer_size 32k;
    # 可以支持 ?s=/m/c/a 的 URL 访问模式
        fastcgi_split_path_info ^((?U).+\.php)(/?.+)$;
    	fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
    # 支持 index.php/index/Index/index 的 pathinfo 模式
    	fastcgi_param PATH_INFO $fastcgi_path_info;
        fastcgi_param PATH_TRANSLATED $document_root$fastcgi_path_info;
        fastcgi_read_timeout 600;
        include fastcgi_params;
    }
```


















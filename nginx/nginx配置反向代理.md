反向代理是客户端访问 web 服务器时，请求发送到代理服务器，由代理服务器决定将此请求转发给哪个真实的 web 服务器。

![](nginx配置反向代理.assets/733013-20171017222119959-1335404489.png)

nginx 通过反向代理实现负载均衡的效果。因为是通过反向代理实现的负载均衡，所以 nginx 实现的是七层负载均衡。它能够识别 http 协议，根据 http 报文将不同类型的请求转发到不同的后端 web 服务器上。后端的 web 服务器称为 “上游服务器”，即 upstream 服务器。

nginx 和 php-fpm 结合，指令 fastcgi_pass 实现的也是反向代理的功能，只不过这种代理是特定的只能转发给 php-fpm。

nginx 反向代理的几种实现方式：

1. 仅使用模块 ngx_http_proxy_module 实现简单的反向代理，指令为 proxy_pass
2. 使用 fastcgi 模块提供的功能，反向代理动态内容，指令为 fastcgi_pass
3. 使用 ngx_http_memcached_module 模块提供功能，反向代理 memcached 缓存内容，指令为 memcached_pass
4. 结合 upstream 模块实现分组反向代理

#### 使用 upstream 模块实现分组反向代理

upstream 默认指定权重的加权算法，还可以指定 ip_hash 算法实现同一个客户端 IP 发起的请求总是转发到同一台服务器上。最常用的还是加权算法，然后通过 session 共享的方式实现同一个客户端 IP 发起的请求转发到同一服务器上。

![](nginx配置反向代理.assets/733013-20171017222647302-344519376.png)



upstream 指令必须定义在 server 段外面，示例：

```nginx
http {
    include mime.types;
    default_type application/octet-stream;
    sendfile on;
    keepalive_timeout 65;

    # define server pool
    upstream dynamic_pool {
        server IP1:80;
    }
    upstream pic_pool {
        server IP3:80 weight=2;
        server IP4:80 weight=1;
    }
    upstream static_pool {
        server IP5:80 weight=1;
        server IP6:80 weight=1;
    }

    server {
        listen 80;
        server_name www.longshuai.com;

        # define how to proxy
        location ~ \.(php|php5)$ {
            proxy_pass http://dynamic_pool;
        }
        location ~ \.(png|jpeg|jpg|bmp|gif)$ {
            proxy_pass http://pic_pool;
        }
        location / {
            proxy_pass http://static_pool;
        }
    }
}
```

#### ngx_http_proxy_module 模块指令及其意义

| 指令                  | 意义                                                         |
| --------------------- | ------------------------------------------------------------ |
| proxy_pass            | 定义代理到哪台服务器或那个upstream                           |
| proxy_set_header      | 在代理服务器上设置 http 报头信息，如加上真实客户端地址 proxy_set_header X_Forwarded_For $remote_addr |
| proxy_connect_timeout | 反向代理连接上游服务器节点的超时时间，发起方是 proxy 方，即等待握手成功的时间 |
| proxy_send_timeout    | 上游服务器节点数据传给代理服务器的超时时间，即此时间段内，后端节点需要传完数据给代理服务器 |
| proxy_read_timeout    | 定义代理服务器何时关闭和后端服务器连接的超时时长，默认60秒，表示某次后端传输完数据给代理服务器后，如果60秒内代理服务器和后端服务器没有任何传输，则关闭此连接。目的是避免代理服务器和后端服务器一直保持连接不断开而占用资源 |

#### proxy_set_header

在代理服务器上设置了头部字段，后端服务器仅仅是获取到它，默认并没有将其记录下来，在 nginx 的访问日志上记录此项，可以在 `log_format` 指令，添加该字段如 `$http_x_forwarded_for`

#### proxy_pass 注意事项 

`proxy_pass` 指令 `/` 注意事项

不带 `/`

```nginx
location /test/ {
    proxy_pass http://zhu-zi.top;
}
```

带 `/`

```nginx
location /test/ {
    proxy_pass http://zhu-zi.top/;
}
```

第一种，访问 `http://shea.test/test/test.html` ，被 nginx 代理后，请求路径变成 `http://proxy_pass/test/test.html` ，将 test/ 作为根路径，请求 test/ 路径下的资源

第二种，访问 `http://shea.test/test/test.html` ，被 nginx 代理后，请求路径变成 `http://proxy_pass/test.html` ，直接访问 server 的根资源

当 proxy_pass 所在的 location 中使用了正则匹配时，则 proxy_pass（包括 fastcgi_pass 和 memcached_pass）定义的转发路径都不能包含任何 URI 信息。

#### ngx_http_upstream_module 模块

```nginx
upstream backend {
    server 192.168.0.24 weight=2 max_fails=2 fail_timeout=2;
    server 192.168.0.45 down;
    server 192.168.0.52 backup;
    ip_hash;	#定义此项，前面 server 附加项 weight 和 backup 功能失效
}
```

意义：

- weight ：定义该后端服务器在组中的权重，默认为 1，权重越大，代理转发到此服务器次数越多
- max_fails 和 fail_timeout ： 定义代理服务器联系后端服务器的失败重联系次数和失败后重试等待时间。默认 1 和 10
- down ：将此后端服务器定义为离线状态。常用语 ip_hash 算法，当使用 ip_hash 时，如果某台服务器坏了需要将其从服务器组中排除，直接从配置文件中直接删除该服务器将会导致此前分配到此后端服务器的客户端重新计算 ip_hash 值，而使用 down 则不会。
- backup ：指定当其它非 backup 的 server 都联系不上时，将返回 backup 所定义的服务器内容。常用于显示 sorry_page

#### nginx 代理 memcached

该模块可以将请求代理至 memcached server 上，并立即从 server 上获取响应数据，例：

```nginx
upstream memcached {
    # 使用 upstream 模块的通过 hash 一致性的算法进行负载均衡
    hash "$uri$args" consistent;
    # 也可以使用第三方模块 ngx_http_upstream_consistent_hash 
    consistent_hash consistent;
    server 127.0.0.1:11211;
    server 127.0.0.1:11212;
    server 127.0.0.1:11213;
}
```

nginx 代理 memcached 时，需要以变量 `$memcached_key` 作为 key 去访问 memcached server 上的数据。例如此处将 `$uri$args` 变量赋值给 `$memcached_key` 变量作为去访问 memcached 服务器上对应的数据






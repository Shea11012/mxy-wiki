## Nginx 简介

采用非阻塞异步的套接字，使用 epoll 方式实现事件驱动，同时采用 **一个 master + N 个 worker 进程（默认）的方式处理请求**

## Nginx 处理请求过程简单说明

master 进程用于管理 worker 进程，例如接受外界信号、向 worker 进程发送信号、销毁 worker 进程、启动 worker 进程等等。

每个 worker 进程是业务处理进程，负责监听套接字、处理请求、响应请求、代理请求至后端服务器等。

每个 worker 进程的大致流程：

- 监听套接字
- 与客户端建立连接
- 处理监听到的连接请求（加载静态文件）
- 响应数据
- 断开连接

每个 worker 进程都是平等的，它们都可以去监听套接字。Nginx 采用争抢 accept 互斥锁的方式，只有持有 accept 互斥锁的 worker 进程才有资格将连接请求接到自己的队列中并完成 TCP 连接的建立。只有 worker 进程当前建立的连接数小于 worker_connections 指定的值时（该值的7/8），才允许争抢互斥锁。除了繁忙程度限制资格，还有 epoll_wait 的 timeout 的指标，等待越久的 worker 进程争抢能力越强。**在某一时刻，一定只有一个 worker 进程监听并 accept 新的连接请求**。

## Nginx 模块及 HTTP 功能

Nginx 由一个核心和一系列的模块组成。

核心（core functionality）主要用于提供全局应用的基本功能，创建必要的运行时环境及确保不同模块之间平滑地进行交互等。

|        http 类模块名         |                         模块功能说明                         |
| :--------------------------: | :----------------------------------------------------------: |
|     ngx_http_core_module     | http 核心模块，对应配置文件中的 http 段，包含很多指令，如 location 指令 |
|    ngx_http_access_module    | 访问控制模块，控制网站用户对 nginx 的访问，对应配置文件中的 allow 和 deny 等指令 |
|  ngx_http_auth_basic_module  | 通过用户名和密码认证的访问控制，如访问站点时需要用户名和密码，指令包括 auth_basic 和 auth_basic_user_file |
|   ngx_http_charset_module    |  设置网页显示字符集，指令之一为 charset，如：charset utf-8   |
|   ngx_http_fastcgi_module    |   fastcgi 模块，和动态应用相关。该模块下有非常多的子模块。   |
|     ngx_http_flv_module      |              支持 flv 视频流的模块，如边下边播               |
|     ngx_http_mp4_module      |                         同 flv 模块                          |
|     ngx_http_gzip_module     | 压缩模块，用来压缩 Nginx 返回的响应报文。一般压缩纯文本内容，因为压缩比例非常大，而图片等不会去压缩。 |
| ngx_http_image_filter_module |  图片裁剪、缩略图相关模块，需要安装 gd-devel 才能编译该模块  |
|    ngx_http_index_module     | 定义将要被作为默认主页的文件，对应指令为 index。如：index index.html,index.php |
|  ngx_http_autoindex_module   | 当 index 指令指定的主页文件不存在时，交给 autoindex 指令，将自动列出目录中的文件 autoindex ｛on/off｝ |
|      ngx_htp_log_module      |    访问日志相关的模块，指令包括 log_format 和 access_log     |
|  ngx_http_memcached_module   | memcached 相关的模块，用于从 memcached 服务器中获取相应响应数据 |
|    ngx_http_proxy_module     |              代理相关，允许传送请求到其他服务器              |
|    ngx_http_realip_module    | 当 nginx 在反向代理的后端提供服务时，获取到真正的客户端地址，否则获取的是反向代理的 IP 地址 |
|   ngx_http_referer_module    |                     实现防盗链功能的模块                     |
|   ngx_http_rewrite_module    | URL 地址重写相关的模块，需要安装 pcre-devel 才能编译安装该模块 |
|     ngx_http_scgi_module     |    simple cgi，是 cgi 的替代品，和 fastcgi 类似，但更简单    |
|     ngx_http_ssl_module      |              提供 ssl 功能的模块，即实现 HTTPS               |
| ngx_http_stub_status_module  |                   获取 nginx 运行状态信息                    |
|      ngx_http_upstream       |                       负载均衡相关模块                       |

## Nginx 配置文件说明

worker_processes 的值和 work_connections 的值决定了最大并发数量。如果每个 worker 进程最大允许 1024 个连接，配置了 4 个 worker 进程，所以并发数量为 `worker_processes * worker_connections`。

在反向代理场景中计算方式不同，因为 nginx 既要维持和客户端的连接，又要维持和后端服务的连接，因此处理一次连接要占用 2 个连接，所以最大的并发数计算方式为：`worker_processes * worker_connections / 2` ，还需注意，除了和客户端的连接，与后端服务器的连接，nginx 可能还会打开其他的连接，这些都会占用文件描述符，从而影响并发数量的计算，同时最大并发数量还受 **允许打开的最大文件描述符数量** 限制。

在 main 段使用 worker_cpu_affinity 指令绑定 CPU 核心。nginx 通过位数识别 CPU 核心以及核心数，指定位数上占位符为 1 表示使用核心，为 0 表示不使用该核心。例如 2 核 CPU 的位数分别为 01 和 10，4核的位数分别为 1000，0100，0010，0001，8 核和 16 核 同理 例如：

```
# 每个 worker 分别对应 cpu0/cpu1/cpu2/cpu3
worker_processes	4;
worker_cpu_affinity 0001 0010 0100 1000

# 有 4 核心，但只有两个 worker，第一个 worker对应cpu0/cpu2,第二个worker对应cpu1/cpu3
worker_processes	2;
worker_cpu_affinity 0101 1010
```

#### http段

root 指令将匹配的 URI 追加在 root 路径后

```
location /i/ {
    root /data/w3;
}
实际访问路径是 /data/w3/i/
```

alias 指令会对 URI 进行替换

```
location /i/ {
    alias /data/w3/images/;
}
```

如果 alias 指令最后一部分包含了  URI 则最好使用 root 指令

它们都能使用相对路径，相对的是 prefix，如 `root html` 是 `/usr/local/nginx/html`

root 和 alias 指令相关的变量为 `$document_root`、`$realpath_root`。其中 `$document_root` 的值是 root 或者 alias 指令的值，而 `$realpath_root` 的值是对 root 或者 alias 指令进行绝对路径换算后的值。

add_header 指令添加响应头字段

```
server {
    add_header RealPath $realpath_root;
}
# 响应头将包含如下
RealPath: /usr/local/nginx-1.12.1/html
```

server_name 指令可以定义多个主机名，第一个名字为虚拟主机的首要主机名。主机名可以含有`*`，以替代名字的开始部分或结尾部分（只能是起始或结尾，如果要实现中间部分的通配，可以使用正则表达式）。

主机名使用正则表达式，就是在名字前面补个一个 `~` ：`server_name ~^www\d+\.example\.com$`

server_name 允许定义一个空主机名，这种主机名可以让虚拟主机处理没有 `Host` 首部请求

#### stub_status 指令获取nginx状态信息

ngx_http_stub_status_module 提供的功能可以获取 nginx 运行的状态信息，对应的指令只有一个，即 stub_status ，例：

```nginx
server {
    listen 80;
    server_name www.example.com;
    location / {
        root /www/;
        index index.php;
    }
    location /status {
        # 可以明确访问该信息不记录日子，且提供访问控制
        stub_status on;
        access_log off;
        allow 192.168.100.0/24;
        deny all;
    }
}
```

访问 主机名/status 获得信息如下：

```nginx
Active connections: 291
server accepts handled requests
	  16630948 16630948 31070465
Reading: 6 Writing: 179 Waitng: 106
```

状态信息意义：

1. 第一行表示当前处于活动状态的客户端连接数，包括正处于等待状态的连接

2. 第二行 accepts 的数量为 16630948 表示从服务启动开始到现在已经接收进来的总的客户端连接数；

   ​	    handled 的数量为 16630948 表示从服务启动以来已经处理过的连接数，一般 handled 的值和 accepts 的值相等，除非做了连接数限定。

   ​	    requests 的数量为服务启动以来总的客户端请求数。一个连接可以有多个请求，所以可以计算出平均每个连接发出多少个请求

3. 第四行     reading 数量为 6，表示 nginx 正在读取请求首部的数量，即正在从 socket rec buffer 中读取的数量

   ​		writing 数量为 179 表示 nginx 正在将响应数据写入 socket send buffer 以返回给客户端的连接数量

   ​		waiting 数量为 106 表示等待空闲客户端发起请求的客户端数量，包括长连接状态的连接以及已接入但 socket recv buffer 还未产生可读事件的连接，其值为 active-reading-writing




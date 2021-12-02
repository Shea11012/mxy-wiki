| 内置变量                                                     |                             意义                             |
| :----------------------------------------------------------- | :----------------------------------------------------------: |
| $arg_name                                                    | 请求中的参数名，即？后面的 arg_name=arg_value 形式的 arg_name |
| $args                                                        |                        请求中参数的值                        |
| $binary_remote_addr                                          |         客户端地址的二进制形式，固定长度为 4 个字节          |
| $body_bytes_sent                                             |            传输给客户端的字节数，响应头不计算在内            |
| $bytes_sent                                                  |                     传输给客户端的字节数                     |
| $connection                                                  |                       TCP 连接的序列号                       |
| $connection_requests                                         |                    TCP 连接当前的请求数量                    |
| $content_length                                              |                  Content-Length 请求头字段                   |
| $content_type                                                |                   Content-Type 请求头字段                    |
| $cookie_name                                                 |                         Cookie 名称                          |
| $document_root                                               |                  当前请求的文档根目录或别名                  |
| $document_uri                         |                           同 $uri |                                                              |
| $host                                                        | 优先级：HTTP请求行的主机名>HOST 请求头字段>符合请求的服务器名 |
| $hostname                                                    |                            主机名                            |
| $http_name                          | 匹配任意请求头字段，变量名中的后半部分 “name” 可以替换成任意请求头字段，如在配置文件中需要获取 http 请求头：“Accept-Language” ，可以写为 “$http_accept_language” |                                                              |
| $https                                                       |       如果开启了 SSL 安全模式，值为 on，否则为空字符串       |
| $is_args                                                     |          如果请求中有参数，值为 ？，否则为空字符串           |
| $limit_rate                                                  |                    用于设置响应的速度限制                    |
| $msec                                                        |                      当前的 Unix 时间戳                      |
| $nginx_version                                               |                          nginx 版本                          |
| $pid                                                         |                        工作进程的 PID                        |
| $pipe                                                        |           如果请求来自管道通信，值为 p，否则为 `.`           |
| $proxy_protocol_addr                                         | 获取代理访问服务器的客户端地址，如果是直接访问，该值为空字符串 |
| $query_string                         |                           同 $args |                                                              |
| $realpath_root                                               | 当前请求的文档根目录或别名的真实路径，会将所有符号连接转换为真实路径 |
| $remote_addr                                                 |                          客户端地址                          |
| $remote_port                                                 |                          客户端端口                          |
| $remote_user                                                 |                用于 HTTP 基础认证服务的用户名                |
| $request                                                     |                     代表客户端的请求地址                     |
| $request_body                                                | 客户端的请求主体，此变量可在 location 中使用，将请求主体通过 proxy_pass，fastcgi_pass 等传递给下一级的代理服务器 |
| $request_body_file                                           | 将客户端请求主体保存在临时文件中。文件处理结束后，此文件删除。如果需要一直开启此功能，需要设置 client_body_in_file_only。如果将此文件传递个后端代理服务器，需要禁用 request_body，即设置 proxy_pass_request_body off等 |
| $request_completion                                          | 如果请求成功，值为 OK，如果请求未完成或者请求不是一个范围请求的最后一部分，则为空。 |
| $request_filename                                            | 当前连接请求的文件路径，由 root 或 alias 指令与 URI 请求生成 |
| $request_length                                              |     请求的长度（包括请求的地址，http 请求头和请求主体）      |
| $request_method                                              |              HTTP 请求方法，通常为 GET 或 POST               |
| $request_time                                                |  处理客户端请求使用的时间，从读取客户端的第一个字节开始计时  |
| $request_uri                         | 此变量等于包含一些客户端请求参数的原始 URI，它无法修改，请查看 $URI 更改或重写 URI，不包含主机名，例如：“/public/test.php?arg=test” |                                                              |
| $scheme                                                      |              请求使用的 Web 协议，HTTP 或 HTTPS              |
| $sent_http_name                        | 可以设置任意 http 响应头字段，变量名中的后半部分 “name”，如需设置响应头 Content-length，可以写为：$sent_http_content_length 4096 |                                                              |
| $server_addr                                                 | 服务器端地址，需要注意：为了避免访问 Linux 系统内核，应将 ip 地址提前设置在配置文件中 |
| $server_name                                                 |                           服务器名                           |
| $server_port                                                 |                          服务器端口                          |
| $server_protocol                                             |       服务器的 HTTP 版本，通常为 HTTP/1.0 或 HTTP/1.1        |
| $status                                                      |                        HTTP 响应代码                         |
| \$tcpinfo_rtt，\$tcpinfo_rttvar，\$tcpinfo_snd_cwnd，\$tcpinfo_rcv_space |                  客户端 TCP 连接的具体信息                   |
| $time_iso8601                                                |                   服务器时间的ISO 8610格式                   |
| $time_local                                                  |                  服务器时间 LOG Format 格式                  |
| $uri                             | 请求中的当前 URI（不带请求参数，参数位于 \$args），可以不同于浏览器传递的 \$request_uri 的值，它可以通过内部重定向，或者使用 index 指令进行修改，$uri 不包含主机名，如 “/foo/bar.html” |                                                              |
| $http_user_agent                                             |                            请求头                            |


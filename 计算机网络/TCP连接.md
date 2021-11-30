每个TCP连接需要的资源：

- 内存资源

  > 每个 TCP 连接都需要占用一定的内存资源。占满会发生 OOM。

- CPU 资源

  > CPU 资源占满会导致系统卡死

- 端口号资源

  > linux 对可使用的端口范围有具体限制
  >
  > 通过指令查看 `cat /proc/sys/net/ipv4/ip_local_port_range` 输出端口范围限制
  >
  > ```bash
  > # vim /etc/sysctl.conf
  > # 新增一行
  > net.ipv4.ip_local_port_range = 60000 60009
  > 
  > # sysctl -p /etc/sysctl.conf 使其生效
  > 
  > ```
  >
  > 根据以上配置，在目标地址和端口不变的情况下，只能有 10 个连接。
  >
  > 只要保证四元组不重复，理论上 TCP 的连接数是 $(2^{32})^4$
  >
  > 端口号超出限制，会报 cannot assign requested address

  

- 文件描述符资源

  > linux 对可打开的文件描述符做了限制
  >
  > - 系统级：当前系统可打开的最大数量 cat /proc/sys/fs/file-max 
  > - 用户级：指定用户可打开的最大数量 cat /etc/security/limit.conf
  > - 进程级：单个进程可打开的最大数量 cat /proc/sys/fs/nr_open
  >
  > 修改单个进程可打开的最大文件描述符 `echo 100000 > /proc/sys/fs/nr_open`
  >
  > 文件描述符超出限制会报 too may open files

- 线程资源

  > 每一个TCP连接都需要消耗一个线程，如果有N个连接，则CPU需要不停的进行上下文切换，会导致TCP连接建立变长。
  >
  > **C10K** 问题
  >
  > 当服务器连接数到达 1 万且每个连接都需要消耗一个线程资源时，操作系统就会不停的忙于线程的上下文切换，导致系统崩溃。每个连接都创建一个线程是最早的 **多线程并发模型**，可以使用 **IO多路复用**解决。


### TCP
TCP 是一个可靠的(reliable)、面向连接(connection-oriented)、基于字节流(byte-stream)、全双工(full-duplex)的协议。发送端在发送数据后启动一个定时器，如果超时没有收到对端确认会进行重传，接收端利用序列号对收到的包进行排序、丢弃重复的数据，TCP还提供了流量控制、拥塞控制等机制保证稳定性。

![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/202110211708548.webp)


![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/202110211652294.webp)


#### TCP 四元组
源IP、源端口、目标IP、目标端口构成了 TCP 连接的四元组。一个四元组可以唯一标识一个连接。

#### 序列号
序列号是一个32位无符号的整数，达到 2^32-1 后循环到0。
TCP 传输的字节流的每个字节都分配了序列号，序列号指本报文段第一个字节的序列号。
序列号加上报文长度，就可以确定传输的是哪一段数据。
> 在 SYN 报文中，序列号用于交换彼此的初始序列号。在其他报文中，序列号用于保证包的顺序。

#### 确认号
TCP 使用确认好来告知对方下一个期望接受的序列号，小于此确认号的所有字节都已收到。
> ACK 包不需要被确认

#### TCP flags
TCP flag 定义为 8 位
常用位：
- SYN（Synchronize）：用于发起连接数据包同步双方的初始序列号
- ACK（Acknowledge）：确认数据包
- RST（Reset）：强制断开连接
- FIN（Finish）：发完所有数据，准备断开连接
- PSH(Push)：这些数据包收到以后应该马上交给上层应用，不能缓存

#### 窗口大小（Window Size）
窗口大小为 16位，最大窗口大小为 2^16-1 即 65535 字节。
因为设计的缺陷，TCP引入了 **TCP 窗口缩放**比例因子，比例因子值的范围 0 - 14。0 表示不缩放。比例因子可以将窗口扩大到原来的 2^n 。如：窗口大小 1050，缩放因子 7，则真正的窗口大小位 1050 * 2^7 大小

#### 可选项
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/202110261115291.webp)


#### 最大传输单元(Maximum Transmission Unit,MTU)
![](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/202110261145806.webp)

IEEE 802.3 以太网标准规定的有效载荷最大为 1500。
同时也存在[巨型帧（jumbo frames)](https://zh.wikipedia.org/wiki/%E5%B7%A8%E5%9E%8B%E5%B8%A7) ，巨型帧需要由硬件和软件一同提供支持。

#### IP 分段
![ip|740](https://mxy-imgs.oss-cn-hangzhou.aliyuncs.com/imgs/202110261153584.webp)

ipv4 数据包最大为 65535 字节。当一个IP数据包大于MTU时，IP会把数据报文进行切割为多个小的片段，使得这些小的报文可以通过链路层进行传输。

#### TCP 最大段大小（Max Segment Size,MSS）
TCP 为了避免被发送方分片，会主动把数据分割成小段再交给网络层，最大分段大小被称为MSS。
MSS 计算公式：`MSS= MTU - IP Header Size - TCP header size`

#### TCP 头部时间戳可选项（TCP Timestamps Option,TSopt）
TSopt 由四部分构成：
- 类别(kind)，kind = 8
- 长度(length)，length = 10
- 发送方时间戳(Ts value)
- 回显时间戳(TS Echo Reply)

发送方发送数据时，将发送时间戳放在 TSval 中
接受方接受到数据包时，将收到的时间戳放在 TSecr 中，将自己的时间戳放在 TSval 中。

TSopt 的作用：
- 测量出两端往返所消耗的时间（RTTM）
- 序列号回绕（PAWS）

#### 半连接、全连接队列
- 半连接队列（Incomplete connection queue），SYN 队列
- 全连接队列（Completed connection queue），Accept 队列

#### 检验和（checksum）
每个TCP包首部都有两个字节来表示检验和，防止在传输过程中有损坏。如果收到一个校验和有差错的报文，TCP 不会发送任何确认直接丢弃它，等待发送端重传。
> 包的序列号保证接收数据的乱序和重复问题

#### 超时重传 
TCP 发送数据后会启动一个定时器，等待对端确认收到这个数据包。如果在指定时间内没有收到 ACK 确认，就会重传数据包，然后等待更长时间。在多次重传仍然失败后，TCP 会放弃这个包。

#### 流量控制

#### 拥塞控制
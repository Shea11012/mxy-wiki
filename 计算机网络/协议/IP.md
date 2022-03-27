## IPV4
### IPV4 数据包格式
![[Pasted image 20220317195329.png]]
- 版本号（Version）：占4bit，通信双方必须使用一致的版本
- 首部长度（Internete Header Length，IHL）：占4bit，由于IPv4首部可能包含数目不定的选项，这个字段也用来确定数据的偏移量
- 区分服务（Differentitated Service，DS）：
- 显示拥塞通告（Explicit Congestion Notification，ECN）
- 全长（Total Length）：报文总长，包含首部和数据
- 标识符（Identification）：用来唯一的标识一个报文的所有分片，因为分片不一定按序到达，所以在重组是需要知道分片所属的报文
- 标志（flags）：控制和识别分片：
	- 位0：保留
	- 位1：禁止分片
	- 位2：更多分片
- 分片偏移（Fragment Offset）：指明每个分片相对于原始报文开头的偏移量
- 存活时间（Time To Live，TTL）：每个报文的存活时间，在报文经过每个路由器时此字段减1，当此字段等于0时，报文被丢弃。
- 协议（Protocol）：仅当IP数据报到达最终目的地时才会用到，协议号指定了IP数据报的数据部分会交给哪个特定的运输层协议。
- 首部校验和（Header Checksum）：只对首部检验，在每一跳中路由器都要对此进行检验。
- 源地址（Source address）
- 目的地址（Destination address）
- 选项（Options）：附加的首部字段（不经常使用）


## IPV6
### IPV6 数据报格式
![[Pasted image 20220317204546.png]]
- 版本（Version）：IPV6值为6
- 流量类型（Traffic Class）
- 流标签（Flow Label）：用于标识一条数据报的流，能够对一条流中的某些数据报给出优先权。
- 有效载荷长度（Payload Length)：指定跟在IPV6中定长的40字节数据报首部后面的字节数量。
- 下一个首部（Next Header）：标识数据报中的数据部分需要交给哪个协议，与IPV4的 Protocol 字段相同
- 跳限制（Hop Limit）：转发数据报的路由器对该字段减1，为0则丢弃。
- 源地址
- 目的地址

## IPV4 vs IPV6
### 不同
- 分片/重新组装：ipv6 不允许在中间路由器上进行分片与重新组装，如果路由器收到的数据报太大则直接丢弃并发回一个 “分组太大”的 ICMP 差错报文，让发送发重新以一个较小的包发送，这样加快了网络的转发速度。
- 首部检验和：将检验操作交给了运输层去处理，加快转发速度。
- 选项：由“下一个首部”代替
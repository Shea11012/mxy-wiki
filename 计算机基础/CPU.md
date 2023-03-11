---
tags: []
date created: 2023-03-12 06:14
date modified: 2023-03-12 07:18
---

## CPU Cache

![](attachments/Pasted%20image%2020230312061926.png)

CPU 访问不同 Cache 的时钟周期
Cache | 时钟周期
:-----: | :-------:
L1 | 2~4
L2 | 10 ~ 20
L3 | 20 ~ 60
内存 | 200 ~ 300
>[!info]
>1GHZ 的 CPU，一个时钟周期为 1 ns；2GHZ 的 CPU，一个时钟周期 0.5
>时钟周期 = 1/CPU 的主频

## MESI 协议

mesi 协议解决 CPU 缓存一致性问题
- Modified：已修改
- Exclusive：独占
- Shared：共享
- Invalidated：已失效

---
date created: 2022-03-09 19:24
date modified: 2022-04-13 03:37
title: GC
---
## 三色标记法
三色标记将对象分为黑、灰、白三种颜色：
- 黑：该对象已经标记过了，且该对象下的属性也全部都被标记过了（活着的对象）
- 灰：对象已经标记过了，但该对象下的属性没有全部标记完（等待 GC 继续扫描）
- 白：没有被标记（垃圾对象）

GC 以广度优先遍历的方式进行，主要工作步骤：
1. GC 从 roots 开始遍历访问，根对象标记为灰色，放入灰色集合
2. 从灰色集合中获取对象，如有将该对象引用到的对象标记为灰色，自身标记为黑色
3. 重复步骤 2，直至灰色集合为空
4. 结束后，剩下的白色对象就会被 GC 回收

三色标记存在的问题：
- 多标 - 浮动垃圾
- 漏标 - 悬挂指针

## 内存屏障
内存屏障，是一种屏障指令，它能使 CPU 或编译器对在该屏障指令之前和之后发出的内存操作强制执行排序约束，在内存屏障前执行的操作一定会先于内存屏障后执行的操作

内存屏障分为：
- 读屏障
- 写屏障

## GC 阶段
1. sweep termination（清理终止）
	1. 触发 STW，所有的 P 都会进入安全点（safe-point）
	2. 清理未被清理的 mspan
2. the mark phase（标记阶段）
	1. 将 GC`_GCoff` 状态改为 `_GCmark`，开启写屏障、协助程序并将根对象入队
	2. 恢复执行程序，标记进程和协助线程会开始并发标记内存中对象。写屏障会将被覆盖的指针和新指针都标记成灰色，而所有新创建的对象都会被直接标记为黑色
	3. 开始扫描根对象、包括**所有 G 的栈、全局对象以及不在堆中的运行时数据结构**，扫描 G 的栈期间会暂停当前 G
	4. 依次处理灰色集合中的对象，将对象标记成黑色并将指向它们的对象标记成灰色
	5. 使用分布式终止算法检测何时不再有根标记作业或者灰色对象，完成后进入标记终止阶段
3. mark termination（标记终止）
	1. STW，GC 状态变为 `_GCmarktermination` 并关闭辅助标记的用户程序
	2. 清理处理器上的线程缓存
4. the sweep phase（清理阶段）
	1. GC 状态变为 `_GCoff` ，初始化清理状态并关闭写屏障
	2. 恢复程序执行，新创建的对象会标记为白色
	3. 后台并发清理所有的 mspan，当 G 申请新的 mspan 时就会触发清理

## 触发方式
- 辅助 GC，在分配内存时会判断当前的 heap 内存分配是否达到了触发一轮 GC 的阈值（每轮 GC 完成后，该阈值会被动态设置），如果超过阈值则启动一轮 GC
- 调用 `runtime.GC` 强制 GC
- sysmon 守护线程，当超过 `forcegcperiod`（2 分钟）没有运行 GC 会启动一轮 GC


refs:
[Garbage Collection In Go : Part I - Semantics (ardanlabs.com)](https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html)
[Go 语言 GC 实现原理及源码分析 - luozhiyun`s Blog](https://www.luozhiyun.com/archives/475)
[Go 语言垃圾收集器的实现原理 | Go 语言设计与实现 (draveness.me)](https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/#72-%E5%9E%83%E5%9C%BE%E6%94%B6%E9%9B%86%E5%99%A8)
[Go 垃圾回收 1：历史和原理 | Go 语言充电站 (lessisbetter.site)](https://lessisbetter.site/2019/10/20/go-gc-1-history-and-priciple/)
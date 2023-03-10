---
date created: 2022-03-26 14:50
date modified: 2022-04-04 04:53
title: map实现
---
## 结构
```go
// runtime/map.go

type hmap struct {
	count     int // 表示当前hash表中的元素数量
	flags     uint8	// map状态标识，是否正在被写或者迁移
	B         uint8  // 表示当前hash表持有的buckets数量，因为桶的数量都是2的倍数，所以存储对数，len(buckets) == 2^B
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash种子，为哈希函数引入随机性，在创建hash表时确定，作为hash函数的参数传入

	buckets    unsafe.Pointer // 指向buckets数组，如果元素个数为0，则是nil.
	oldbuckets unsafe.Pointer // 在扩容时用于保存之前buckets字段，大小是当前buckets的一半
	nevacuate  uintptr        // 扩容进度，小于此地址的buckets迁移完成

	extra *mapextra // 用于扩容的指针
}

// 存储key和value都不是指针类型的map，并且大小都小于128字节，可以避免扫描整个map
type mapextra struct {
	// 为了保证所有的溢出桶存活，将所有的溢出桶指针都存储在 overflow 和 oldoverflow 中
	// overflow 和 oldoverflow 仅在key和elem都不包含指针的时候才会被使用到
	// overflow 包含 hmap.buckets 的溢出存储桶
	// oldoverflow 包含 hmap.oldbuckets 的溢出桶
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow 包含指向可用溢出存储桶的指针
	nextOverflow *bmap
}


type bmap struct {
	// tophash 是hash值的高8位
	tophash [bucketCnt]uint8
	// 以下字段会在编译时由编译器自动添加
	// keys 	每个桶最多可以装8个key
	// values 	
	// overflow pointer 发生hash碰撞后创建overflow bucket
}
```

hash 函数根据 key 生成一个哈希值，其中低位哈希用来判断桶位置，高位哈希用来确定在哪个桶中。

根据 `hmap.B` 值确定该 map 中有多少个桶，如：`B=5` 则有 `2^5=32` 个桶，只需要取 hash 值的低 5 位就可以确定当前落在哪个桶中；高位 8 位的 hash 确定在桶中的位置，每个桶可以存储 8 个 kv，存储结构为 `key/key/key...value/value/value` 这样可以避免字节对齐时的 padding，节省内存空间。

当桶溢出时，将 kv 存储在 overflow bucket 中，通过 `bmap.overflow` 指针将它们串连起来

## 扩容
1. 装载因子\>6.5，装载因子=元素个数/桶个数
2. overflow bucket 过多时，`>=1<<(B&15)`

第一种情况，map 是翻倍扩容，将旧桶里的值根据新的**B**重新分配到不同的新桶中
第二种情况，map 是等量扩容，将原有的数据拷贝到新的桶中，清除老的桶并释放内存

map 的数据转移是扩容后逐步进行的，在访问或者删除时做至少一次迁移工作

## 写入
依旧是根据 key 的 hash，判断是否正在扩容，是则执行一次扩容操作，再根据 hash 计算桶和桶中位置

## 访问
根据 key 计算出 hash，根据 hash 找到对应的桶和桶中的位置，判断是否正在扩容，对应的桶是否已经完成扩容，没有则去旧桶中找，否则遍历所有的桶包括溢出桶

## 删除
获取 key 的 hash，判断是否在扩容，是则执行一次扩容操作，根据 hash 找桶，执行清除操作、计数减 1 等


## ref:
[Golang 底层实现系列——map 的底层实现 - SegmentFault 思否](https://segmentfault.com/a/1190000040269520)
[理解 Golang 哈希表 Map 的原理 | Go 语言设计与实现 (draveness.me)](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-hashmap/#33-%E5%93%88%E5%B8%8C%E8%A1%A8)
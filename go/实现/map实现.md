---
tags: []
date created: 2022-03-26 14:50
date modified: 2023-04-13 00:31
go version: 1.20
---

## 结构

```go
// runtime/map.go

type hmap struct {
	count     int // 表示当前hash表中的元素数量
	flags     uint8	// map状态标识，是否正在被写或者迁移
	B         uint8  // 表示当前hash表持有的buckets数量，因为桶的数量都是2的倍数，所以存储对数，len(buckets) == 2^B
	noverflow uint16 // overflow 的 bucket 近似数
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

当 map 的 key 和 value 都不是指针，并且 size 都小于 128 字节，会把 bmap 标记为不含指针，这样可以避免 GC 时扫描整个 hmap，提升效率。但是 bmap 里面包含一个 overflow 字段是指针类型，这破坏了 bmap 不包含指针的设想，这时会把 overflow 移动到 extra 来，启用 Overflow 和 oldoverflow

hash 函数根据 key 生成一个哈希值，其中低位哈希用来判断桶位置，高位哈希用来确定在哪个桶中。

根据 `hmap.B` 值确定该 map 中有多少个桶，如：`B=5` 则有 `2^5=32` 个桶，只需要取 hash 值的低 5 位就可以确定当前落在哪个桶中；高位 8 位的 hash 确定在桶中的位置，每个桶可以存储 8 个 kv，存储结构为 `key/key/key...value/value/value` 这样可以避免字节对齐时的 padding，节省内存空间。

当桶溢出时，将 kv 存储在 overflow bucket 中，通过 `bmap.overflow` 指针将它们串连起来

## 扩容

1. 装载因子\>6.5，装载因子=元素个数/桶个数，表示元素过多 bucket 快要被装满了，需要增加 bucket
2. overflow bucket 过多时，`B>15`，表示 bucket 过多，元素太分散，需要减少 bucket

第一种情况，map 是翻倍扩容，将旧桶里的值根据新的 **B** 重新分配到不同的新桶中
第二种情况，map 是等量扩容，将原有的数据拷贝到新的桶中，清除老的桶并释放内存

map 的数据转移是扩容后逐步进行的，每次最多只会搬迁 2 个 bucket，在插入、修改、删除时，都会检查 odlbucket 是否搬迁完毕

map 不是线程安全的，在每次查找、赋值、遍历、删除时都会检测写标记，发现写标记则 panic

## ref:

[理解 Golang 哈希表 Map 的原理 | Go 语言设计与实现 (draveness.me)](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-hashmap/#33-%E5%93%88%E5%B8%8C%E8%A1%A8)
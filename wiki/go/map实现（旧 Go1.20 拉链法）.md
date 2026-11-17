---
title: "Map 实现（旧 - 拉链法）"
created: 2021-11-30
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - go
  - go/runtime
confidence: INFERRED
---
> [!warning] 已过时
> Go 1.24+ 起 map 实现已从拉链法改为 Swiss Table 开放寻址法，见 [[map实现（Swiss Table）]]

## 哈希函数

哈希函数的选择在很大程度上能够决定哈希表的读写性能，在理想情况下，哈希函数应该能够将不同的键映射到不同的索引上，这要求**哈希函数输出范围大于输入范围**，由于键的数量远远大于映射范围，所以实际使用时，理想结果不可能实现。

**不均匀哈希函数**：较均匀的哈希函数中，哈希增删改查都是 O(1) 的时间复杂度，不均匀的哈希函数会导致所有的操作都会占用最差 O(n) 的复杂度。

## 冲突解决

### 开放寻址法

> 核心思想是对数组中的元素依次探测和比较以判断目标键值对是否存在于哈希表中。

写入数据时，当一个 key1 的位置已经被 key2 占用时，会在 key2 的后面寻找空闲的位置放置 key1；读取时，会从索引的位置开始对数组进行线性探测，找到目标键值对或者空内存就意味着查找结束。

开放寻址法中性能影响最大的就是装载因子（数组中元素数量与数组大小的比值），装载率超过 70% 后性能急剧下降，达到 100% 则完全失效。

### 拉链法

绝大部分的编程语言都采用拉链法实现哈希表，实现比开放寻址复杂，但平均查找长度比较短，存储的节点的内存是动态申请的。

拉链法一般使用数组加上链表，有些语言会采用红黑树进行优化。

写入时，先对键求哈希值取模选择桶，遍历链表：找到相同键则更新，未找到则在链表末尾追加。一个性能较好的哈希表中，每个桶应该有 0~3 个元素，装载因子 = 元素数量 / 桶数量，一般不超过 1。

## Go Map 实现（拉链法）

### hmap 结构体

```go
// runtime/map.go
type hmap struct {
    count     int            // 当前哈希表中的键值对数量
    flags     uint8          // map 状态标识，是否正在被写或者迁移
    B         uint8          // 桶数量的对数，len(buckets) == 2^B
    noverflow uint16         // overflow bucket 近似数
    hash0     uint32         // hash 种子，引入随机性

    buckets    unsafe.Pointer // 指向 buckets 数组
    oldbuckets unsafe.Pointer // 扩容时保存之前 buckets，大小是当前一半
    nevacuate  uintptr        // 扩容进度，小于此地址的 buckets 迁移完成

    extra *mapextra           // 扩容相关指针
}
```

### mapextra

当 map 的 key 和 value 都不是指针，并且 size 都小于 128 字节，会把 bmap 标记为不含指针，避免 GC 扫描整个 hmap。但 bmap 的 overflow 字段是指针类型，破坏了这一设想，所以把 overflow 移到 extra 中：

```go
type mapextra struct {
    overflow    *[]*bmap  // hmap.buckets 的溢出桶
    oldoverflow *[]*bmap  // hmap.oldbuckets 的溢出桶
    nextOverflow *bmap    // 下一个可用溢出桶
}
```

### bmap 桶结构

```go
type bmap struct {
    tophash [bucketCnt]uint8  // hash 值的高 8 位
    // 以下字段由编译器自动添加
    // keys     每个桶最多 8 个 key
    // values
    // overflow pointer  发生 hash 碰撞后创建 overflow bucket
}
```

hash 函数根据 key 生成哈希值，**低位哈希用来判断桶位置，高位哈希用来确定在桶中的位置**。

根据 `hmap.B` 值确定桶数：`B=5` → `2^5=32` 个桶，取 hash 值低 5 位确定桶；高 8 位确定桶内位置。每个桶存储 8 个 kv，结构为 `key/key/key...value/value/value...`，避免字节对齐 padding 节省内存。

当桶溢出时，kv 存储在 overflow bucket 中，通过 `bmap.overflow` 指针串联。

tophash 作为一级缓存帮助快速遍历桶中元素，减少访问键值对次数以提高性能。

### 扩容

1. **装载因子 > 6.5**（count / 2^B > 6.5），翻倍扩容，旧桶里的值根据新的 B 重新分配
2. **overflow bucket 过多**：
   - B ≤ 15：noverflow ≥ 2^B
   - B > 15：noverflow ≥ 2^15
   - 等量扩容，将数据拷贝到新桶，清除老桶释放内存

map 的数据转移是渐进式的，每次最多搬迁 2 个 bucket，在插入、修改、删除时检查 oldbucket 是否搬迁完毕。

### 并发安全

map 不是线程安全的，每次查找、赋值、遍历、删除时都会检测写标记，发现写标记则 panic。

## 遍历 Map 的实现

编译器使用 `runtime.mapiterinit` 和 `runtime.mapiternext` 重写 for/range 循环：

```go
ha := a
hit := hiter(n.Type)
th := hit.Type
mapiterinit(typename(t), ha, &hit)
for ; hit.key != nil; mapiternext(&hit) {
    key := *hit.key
    val := *hit.val
}
```

`mapiterinit` 使用 `fastrand` 生成随机数随机选择起始桶，Go 故意不提供固定遍历顺序：

```go
func mapiterinit(t *maptype, h *hmap, it *hiter) {
    it.t = t
    it.h = h
    it.B = h.B
    it.buckets = h.buckets
    r := uintptr(fastrand())
    it.startBucket = r & bucketMask(h.B)
    it.offset = uint8(r >> h.B & (bucketCnt - 1))
    it.bucket = it.startBucket
    mapiternext(it)
}
```

遍历过程：选出一个正常桶开始遍历 → 遍历对应的所有溢出桶 → 按索引顺序遍历其余桶 → 所有桶遍历完成。

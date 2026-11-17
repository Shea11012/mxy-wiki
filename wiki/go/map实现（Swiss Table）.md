---
title: "Map 实现（Swiss Table）"
created: 2025-04-23
updated: 2026-10-04
status: needs-review
type: concept
tags:
  - go/runtime
  - go/performance
source: "raw/2026/10/04/How Go Maps Work Internally - Swiss Tables.md"
sources:
  - "raw/2026/10/04/How Go Maps Work Internally - Swiss Tables.md"
  - "wiki/go/Go map Swiss Tables 视频学习笔记.md"
source_kind: synthesis
as_of: 2026-10-04
confidence: AMBIGUOUS
contested: false
contradictions: []
related:
  - "Go map Swiss Tables 视频学习笔记"
  - "map实现（旧 Go1.20 拉链法）"
---
# Map 实现（Swiss Table）

## 阅读范围与证据分层

本页整合原实现笔记与 Phuong Le 的视频学习笔记，以概念为主入口。视频的全部实质主题、算例、性能条件和限制保留在主体，原页独有的源码细节放在文末待核验区。

- 主体中的视频陈述：依据英文人工字幕，可追溯到 [[How Go Maps Work Internally - Swiss Tables]]；摘录为 raw 的精确子串。记录作者说法不等于独立核验当前实现。
- 原页中的源码片段：原笔记只提供浮动 master 链接，没有固定 revision，因此不当作已审计源码。
- 视频基准：仅适用于作者描述的版本、开关、硬件和操作分布。
- `mapsplitgroup`：布局算例保留；版本状态与未来默认行为未核验。

> [!warning] UNVERIFIED — 旧笔记来源边界
> 文末“待核验实现细节”中的结构体字段、符号、阈值和行为来自原笔记，不能以本次视频作为它们的证据。保留是为了后续核验，不代表本次确认正确。

## 术语速查

| 术语 | 角色 |
|---|---|
| Slot | 一个 key/value 的存储位置 |
| Group | 8 个 slot 与对应控制字节 |
| H1 | hash 去掉低 7 位后的部分，参与寻址 |
| H2 | hash 的低 7 位，作为候选指纹 |
| Table | 一组 group，独立增长或分裂 |
| Directory | table 指针目录，使用 hash 高位索引 |
| Tombstone | 删除标记；不会像 empty 一样终止探测 |

### H1/H2 与探测公式

原页曾写 `H1 = hash >> 57`，这只保留最高 7 位，与“64-bit hash 的高 57 位”不一致。按照视频的 57/7 拆分，表示应为：

```text
H1 = hash >> 7
H2 = hash & 0x7F
start = H1 mod groupCount
probe(i) = (start + i*(i+1)/2) mod groupCount, i = 0,1,2,...
```

> [!note] INFERRED — 公式表达
> 位移与探测公式是依据视频的哈希拆分和逐次增加步长所作的数学表达，不是本次读取的 Go 源码摘录。directory 的 table 选择仍需单独使用 hash 高位，不能把 table 内的 `start` 当作目录索引。

批量匹配指的是一个 group 内的 8 个控制字节筛选，不是同时探测 8 个不同 group。单个 table 的搬迁可以一次完成，但整个 map 通过多个 table 局部增长；不能把它与旧版 bucket evacuation 的“渐进式迁移”混为一谈。

## 一句话总结

视频从 runtime 视角解释 Go map 的 Swiss Table 布局：通过 `directory → table → group` 分层、H1/H2 哈希拆分、控制字节和批量候选筛选降低查找成本；当容量增长时以 table 为单位分裂。它同时强调，这种设计并非所有场景都更快：墓碑、table split 和特定填充模式会增加写入或尾延迟，微基准的收益也不能直接等同于真实应用收益。

> 摘录：Go 1.24 replaced the old map implementation

## 1. 从开发者视角到 runtime 视角

对开发者而言，`map[string]int` 是键值对集合；视频将其展开为 runtime 的内部对象。map 变量指向一个 `maps.Map`，其中保存 map 的状态以及找到实际存储的路径。把 map 赋给另一个变量，或把它作为参数传给函数，并不会复制所有键值对，而是复制指向同一 map 的指针，因此通过任一变量写入都能从另一变量看到。

视频还介绍了 `used` 与 `seed` 的角色：`used` 表示 map 中的 live entries，视频用它解释 `len` 如何直接得到数量；每个 map 在创建时获得自己的随机 seed，hash 实际结合 key 与该 seed。因此相同 key 在不同 map 中不应被假定拥有稳定的内部位置。这里讨论的是视频示例，不是对可观察 API 行为的承诺。

> 摘录：Every map gets its own random seed

## 2. 存储层次：directory、table 与 group

Swiss Table 的存储层次为：

```text
directory → table → group → 8 个 slot（key/value）
```

- **group**：包含 8 个 key/value slot 和 8 个控制字节。
- **table**：包含一组 group，并记录 `used`、容量和增长相关状态。
- **directory**：保存 table 指针；hash 的一部分先选择 directory entry，再在 table 内选择 group。

小 map 可以让 `dirPtr` 直接指向一个 group；较大的 map 才需要 directory 和 table。视频用 4、100、1000 个 key 的布局作为教学例子，但特别提醒：实现不会在恰好 100 或 1000 个 key 时机械切换布局。

> 摘录：room for up to 8 pairs

### group 的控制字节

8 个控制字节与 8 个 slot 一一对应，并可以放在一个 64-bit word 中。控制字节描述 slot 是 live、empty 还是 deleted/tombstone；key 和 value 的具体大小取决于 map 类型，因此视频中的结构体是解释布局的伪代码，而非可以直接复制的完整 runtime 定义。

## 3. H1/H2：先筛选，再比较完整 key

在视频的 64-bit 示例中，hash 被拆为两部分：

- **H1**：顶部 57 bits，用于确定起始 group（以及在相应层次中参与 table/directory 路由）。
- **H2**：底部 7 bits，写入控制字节，用于快速筛选候选 slot。

控制字节最高位用来区分 live 与非 live。视频给出 `empty = 128`、`deleted = 254` 的示例编码。查找遇到 empty 时可以停止，因为后续探测路径上没有更早被插入的 key；遇到 tombstone 则不能停止，因为删除留下的空洞后面仍可能有目标 key。H2 相同只说明可能是候选，最终仍要比较完整 key。

删除不会让 map 缩小，也不会把内存返还给系统。一个 slot 可能被复用；增长或重排时，墓碑会被清理。墓碑的代价会在后面的 delete/insert 基准中体现。

> 摘录：top 57 bits become H1
>
> 摘录：the bottom 7 bits become H2
>
> 摘录：A deleted slot is also called a tombstone

## 4. 一次检查一个 group 的 8 个控制字节

更新或查找 `cow` 时，runtime 先用 H2 在控制字节中找候选，再读取并比较候选 key。因为 8 个控制字节连续放在一个 word 中，硬件可以批量筛选，避免逐一读取和比较全部 key。视频称 AMD64 使用 SIMD 指令；其他架构则可以把 word 当作整数，用算术和位运算生成候选 bitmap。这里的关键优化是“先用便宜的指纹缩小候选集，再做昂贵的完整 key 比较”，而不是 H2 本身保证唯一性。

> 摘录：control bytes in parallel

## 5. group 满后的增长与三角探测

第 9 个 key 使一个 group 满时，视频示例把 group 数从 1 增为 2，搬迁原有 8 个 entry，再插入新 key。更大的 table 会按 group 数扩展，重新依据 hash 将 live key 放到新的 group 中。H1 的一部分选择起始 group；如果起始 group 没有空位，就要继续探测。

视频介绍的探测序列不是简单线性递增，而是 **triangular probing（三角探测）**：第一次前进 1 个 group，之后前进 2、3、4……个 group。group 数为 2 的幂时，这个序列可以遍历所有 group 且不重复。它的目标是减少线性探测容易形成的相邻 cluster；视频中的示例将查找范围从 4 个 group 减到 2 个，但这不是所有负载下的固定比例收益。

> 摘录：triangular probe sequence
>
> 摘录：each step grows by 1

例如 8 个 group、起点为 3 时，前几步为 `3 → 4 → 6 → 1`：分别加 1、2、3，并对 group 数取模。插入或查找必须使用相同的探测规则，才能在冲突后找到原先放置的 key。

## 6. load factor、1024-slot 上限与 table split

视频将 load factor 定义为：

```text
(live entries + deleted slots) / total slots
```

教学示例使用 7/8，也就是 87.5% 作为安全上限。因此 16 个 slot 中即使当前只有 10 个 live entry，也不能把剩下的 6 个 empty 全部当作可用容量，因为增长判断还要为 tombstone 和探测路径保留空间。

单个 table 的 group 数逐步加倍，最多 128 个 group、1024 个 slot。达到这个上限后，不能再把同一 table 加倍，而是把它 **split** 成两个各有 1024 个 slot 的 table。split 会迁移被分裂 table 的 entry；这与“整个 directory 加倍”是两个不同动作。

> 摘录：entries and deleted slots
>
> 摘录：7 out of 8 slots, or 87.5%
>
> 摘录：This is a split, not a doubling

## 7. directory 如何增长

directory 使用 hash 的高位选择 table，table 内部再使用 H1 的低位选择 group，最后使用 H2 筛选 8 个 slot。directory 中多个 index 可以指向同一个 table，所以一个 table split 不必搬迁所有 table，只处理被分裂的那个 table；视频据此把单次迁移的规模限制在最多 1024 个 slot。

这里有两个深度：

- **global depth**：directory 使用多少位来选择 index。
- **local depth**：识别某个 table 所需的位数。

当 local depth 小于 global depth 时，可以在不扩展 directory 的情况下 split；当两者相等时，需要先将 directory 加倍并提高 global depth，再进行 split。这个机制把 directory 的增长和 table 的局部迁移解耦。

## 8. `make(map, hint)` 的 hint 不等于精确容量

视频把 `make(map[K]V, n)` 中的 `n` 称为 hint，而不是“精确预留 n 个 slot”。示例中：

- hint 为 4 时，起始布局仍是一个 8-slot group。
- hint 为 100 时，13 groups 有 104 slots，按 87.5% 上限只能容纳 91 个 key，因此还不够。需要至少 15 groups，再将 group 数向上取为 2 的幂，得到 16 groups、128 slots。
- hint 为 1000 时，一个 128-group table 按该负载上限只能容纳 896 个 key，因此示例使用两个 table，共 2048 slots。

因此 hint 的工程意义是减少达到目标规模前的增长概率，而不是保证 map 内部恰好分配 n 个位置，也不是跨版本稳定的内存承诺。

> 摘录：exactly that many slots
>
> 摘录：one table can only hold 896

## 9. 与旧版 overflow bucket map 的对比

视频把 Go 1.24 前的实现概括为 bucket + overflow bucket 链：每个 bucket 能放 8 个键值对，冲突后增加 overflow bucket。查找沿 overflow 链走时会增加 pointer load；overflow bucket 也可能与主 bucket 相距较远，导致额外的内存访问等待。

旧实现扩容时建立两倍大小的新 bucket array，并渐进地把旧 bucket 搬过去。扩容期间新旧数组共存，lookup 还需要判断目标 bucket 是否已经搬迁。视频用这一点对比 Swiss Table：Swiss Table 以一个最多 1024-slot 的 table 为局部迁移单位，在一次 split 中完成该 table 的重排。

旧实现笔记 [[map实现（旧 Go1.20 拉链法）]] 适合解释旧版 overflow bucket 机制，本页整合 Swiss Table 的教学路径、内部机制与性能边界。两者是不同版本的互补资料。

> 摘录：Go adds an overflow bucket
>
> 摘录：adds one more pointer load

## 10. 视频中的 benchmark：收益必须连同条件阅读

视频没有直接拿 Go 1.23 与 Go 1.24 做简单对比，因为版本之间还可能混入其他变化。它使用 Go 1.25.7，在同一个 benchmark 中分别测试正常 Swiss Table 和通过 experiment flag 恢复的旧 map；硬件是 Apple M4 Pro（ARM64）。因此该实验不包含 AMD64 上 SIMD 控制字节比较的效果。视频先重复运行 Swiss Table，两次结果差异不超过 3%，然后才把大于该噪声范围的差异当作值得讨论的变化。

> 摘录：use Go 1.25.7
>
> 摘录：Apple M4 Pro, which
>
> 摘录：Those 2 runs never differ by more than 3%

### 视频报告的正向结果

以下数字都是视频特定的版本、平台、数据规模和操作模式下的观察：

| 场景 | 视频报告的观察 | 阅读边界 |
|---|---|---|
| 约 65,000 keys 的 lookup | 多数 lookup 快 58%–68% | 是该 benchmark 的查找场景，不是所有 map 操作的保证 |
| Go team microbenchmark（视频引用） | 某些操作最高快约 60% | 是微基准，不能替代应用级测量 |
| full application benchmark（视频引用） | CPU time 约少 1.5% | 更接近应用级收益，但仍是视频引用的外部结果 |
| 约 4 million keys、增长过程 | 旧实现峰值约为最终占用的 1.6 倍；Swiss 峰值接近最终占用 | 只说明该测试的扩容空间峰值 |
| 最慢单次 insert | 旧实现约 8ms，Swiss 约 0.1ms | 关注极端搬迁停顿的一个测试 |
| 约 100 万或 400 万 keys、超出 cache 的 lookup | 仅快约 15%–45% | cache 状态改变后收益下降 |

视频也给出一个重要解释：real application 的整体收益通常不会等于 map 独立微基准的收益，因为应用还有其他计算、I/O、锁和调度成本。

> 摘录：most lookups are 58 to 68% faster

### 视频报告的不利结果

Swiss Table 会把一次很大的搬迁拆成多次 table split，但 split 仍会制造中等大小的暂停。视频称最慢的 1/10,000 次 insert 约慢 5 倍，这是一个尾部区间的比较，不是平均插入时间。持续删除一个 key、再插入另一个 key 的 1-million-key 场景中，约 219,000 个墓碑使插入常常需要多探测一两个 group，Swiss 比对照慢 17%–38%。

插入不能在遇到第一个墓碑时立即结束：还需要确认同一个 key 是否已在后续 group 中，避免将更新误判为新增。查到 empty 后才知道 key 不存在，此时可以回头复用之前记录的墓碑位置。

另一个“预留 1 million keys 后填满”的测试中，视频报告 Swiss 慢 30%–52%，并将其与等待主存和分支预测联系起来。作者解释，CPU 等待主存时可以基于预测继续执行，但批量控制字节匹配后的路径依赖返回的数据；旧实现逐一检查 slot 的分支在该场景中反而更容易被预测，能够提前推进后续工作。视频还展示了一个刻意改为逐一检查 control byte 的 runtime 变体，在这个填充测试中追平旧实现。这是作者对特定测试的解释，不是所有硬件上的普遍优势。

> 摘录：inserts get about 5 times slower
>
> 摘录：it's 17 to 38% slower
>
> 摘录：the Swiss Table is 30 to 52% slower

## 11. 视频末尾的 `mapsplitgroup` 说法

> [!warning] AMBIGUOUS — 版本状态未核验
> 以下只记录视频说法，未以固定 Go 源码或官方文档核验，不把它写成当前 Go 的确定行为。

视频称 Go 1.27 有一个名为 `mapsplitgroup` 的实验布局：原布局让每个 key 与 value 相邻；split 布局把 8 个 key 放入一个数组、8 个 value 放入另一个数组，但二者仍使用相同索引。查找时，control bytes 先找到候选 key，只有 key 匹配后才需要读取 value，因此 key 在内存中相邻。

视频用 `map[int64]struct{}` 做空间算例：一个 int64 key 占 8 bytes，空结构体 value 占零 bytes。原布局中零大小 value 位于每个 slot 末尾；作者解释，为使指向该字段的指针不越出 slot，需要补 1 byte，再按 8-byte 对齐，使每个 slot 达到 16 bytes。因此 8 个 slot 的 key 加 padding 占 128 bytes，再加 8-byte control word，共 136 bytes。

split 布局把 8 个 key 集中到 64 bytes，零大小 value 数组放在整个 group 末尾，只需一次 padding。key 加 padding 为 72 bytes，加上不变的控制字共 80 bytes，视频据此称每个 group 节省 56 bytes。节省来自消除逐 slot 重复的 padding，不是减少 key 本身。这个数字是布局算例，不是对所有 key/value 类型的节省比例。

视频还称该布局当时仍需通过 `GOEXPERIMENT` 启用，并预测未来版本可能成为默认。预测部分必须与事实分开；若要把这些版本状态写入稳定知识页，应另行固定 Go 源码 revision 或官方文档作为来源。

> 摘录：Go 1.27 has an experimental
>
> 摘录：All 8 keys fit in 64 bytes
>
> 摘录：we go from 128 bytes to 72
>
> 摘录：is still an experiment

## 12. 工程阅读要点

1. 关注 map 的**操作分布**，不要只看平均 lookup：高频 delete/insert、增长、填满和尾延迟可能决定结果。
2. 使用 `make` 的 hint 可以降低增长概率，但不能把 hint 当作精确容量或固定内存布局。
3. 看到 benchmark 数字时同时记录 Go 版本、实验开关、CPU 架构、数据规模、是否超 cache、操作序列和噪声控制。
4. 新旧 map 笔记是有版本边界的互补材料，不应把旧版 overflow bucket 的描述当作当前实现，也不应把视频中的简化结构当作稳定源码 API。
5. 对 `mapsplitgroup`、视频引用的 Go team 数字以及视频对未来版本的预测，进一步使用官方源码、release notes 或可复现 benchmark 核验后再升级置信度。

## 视频回看定位

| 时间 | 内容 |
|---|---|
| 00:00 / 02:19 | map 对象、复制与传参、seed |
| 03:04 / 03:53 | 存储层次与 group 控制字节 |
| 04:38 / 07:09 | H1/H2、墓碑与批量筛选 |
| 08:22 / 10:21 | table 增长与三角探测 |
| 11:56 / 13:26 | 负载因子、1024-slot 上限与目录增长 |
| 17:43 / 19:32 | 容量 hint 与旧版 overflow bucket |
| 22:06 | 基准设置、查找、空间峰值、尾延迟与写入退化 |
| 27:35 | 视频所述的 mapsplitgroup 布局 |

## 待核验实现细节

> [!warning] UNVERIFIED — 原笔记保留内容
> 以下片段来自整合前的实现笔记。原文没有固定源码 revision，本次没有读取官方源码，因此字段、符号、阈值和行为均须逐项核验；不能把它们当作视频直接陈述。

### 原笔记结构体快照

#### Map（替代旧 hmap）

```go
// internal/runtime/maps/map.go
type Map struct {
    used            uint64         // 所有 table 中的元素总数（替代旧 count）
    seed            uintptr        // 哈希种子（替代旧 hash0，64位系统使用64位hash）
    dirPtr          unsafe.Pointer // 指向 table 指针数组，或指向单个 group（小map优化）
    dirLen          int            // directory 长度，0 表示小 map
    globalDepth     uint8          // 目录查找使用的位数
    globalShift     uint8          // 64 - globalDepth，用于右移取高位
    writing         uint8          // 并发写检测标志，XOR toggle
    tombstonePossible bool         // 是否存在墓碑
    clearSeq        uint64         // Clear 调用计数器，迭代时检测
}
```

**与旧 hmap 的关键区别**：

| 旧 hmap | 新 Map |
|---------|--------|
| `B uint8` 桶数对数 | 无 B，由 table.capacity 描述 |
| `oldbuckets` 渐进迁移 | 无，整个 table 一次性 rehash |
| `noverflow uint16` | 原笔记称无此字段；墓碑与 overflow 的作用不可混同 |
| `extra *mapextra` | 原笔记称无此字段；GC 扫描行为未核验 |
| `flags uint8` 位标记 | 原笔记记录 `writing uint8` XOR toggle；可靠性比较未核验 |

#### Table（替代旧 buckets 数组）

```go
// internal/runtime/maps/table.go
type table struct {
    used       uint16  // table 中已填充的 slot 数
    capacity   uint16  // 总 slot 数（2^N），最大 maxTableCapacity(1024)
    growthLeft uint16  // 剩余可插入数（含墓碑折算）
    localDepth uint8   // 该 table 创建时的目录深度
    index      int     // 在 directory 中首次出现的索引，-1 表示已过期
    groups     groupsReference // group 数组
}
```

#### Group + 控制字（替代旧 bmap）

```go
// 控制字节编码：
//   empty:   1 0 0 0 0 0 0 0  (0x80)
//   deleted: 1 1 1 1 1 1 1 0  (0xFE)  墓碑
//   full:    0 h h h h h h h           低7位 = H2

type ctrlGroup uint64  // 8 字节，每个字节对应一个 slot 的状态

type groupReference struct {
    data unsafe.Pointer // 指向 group 内存（ctrl + keys + elems）
}
```

Group 的内存布局（默认交错存储）：

```
| ctrls(8B) | key0 | elem0 | key1 | elem1 | ... | key7 | elem7 |
```

也可通过 `GOEXPERIMENT=mapsplitgroup` 切换为分离存储（类似旧的 key/key/.../value/value/...）。

### 原笔记插入符号与回收路径

```
1. 查找 key 是否已存在 → 存在则更新 elem
2. 不存在时，优先使用探测路径上的墓碑 slot（不消耗 growthLeft）
3. 无墓碑则使用 empty slot，消耗 growthLeft
4. growthLeft == 0 时，先尝试 pruneTombstones 回收墓碑
5. 回收后仍无空间 → rehash
```

### 原笔记墓碑回收阈值

旧拉链法删除只需清空 tophash，无需墓碑。新开放寻址法**必须使用墓碑**：

- 如果 group 中仍有 empty slot → 直接标记为 empty
- 如果 group 完全满（可能在探测序列中间） → 标记为 deleted（墓碑）

墓碑不影响正确性，但会消耗 growthLeft。`pruneTombstones` 会尽量回收不需要的墓碑（当墓碑 ≥ 10% capacity 时触发）。

### 小 Map 优化

原笔记记录的小 map 路径（不据此断言旧版没有小 map 优化）：

- 元素 ≤ 8 时（`dirLen == 0`），直接使用**单个 group** 存储，不分配 table 和 directory
- `dirPtr` 直接指向 group 内存
- 小 map 不存在墓碑（单 group 无需维护探测序列，删除直接标记 empty）
- 超过 8 个元素时调用 `growToTable` 升级为完整 map

### 并发检测

```go
// 写操作前
m.writing ^= 1  // toggle
// 写操作后
m.writing ^= 1  // toggle 恢复

// 读操作检查
if m.writing != 0 {
    fatal("concurrent map read and map write")
}
```

原笔记将 XOR toggle 与旧实现 `flags` 位标记作对比，但“检测更可靠”的论断尚未核验。该代码是检测逻辑的示意，不是同步机制；不能将其视为并发访问的安全保证。

### 遍历

迭代器 `Iter` 替代旧 `hiter`，需要处理更复杂的场景：

1. **随机起始**：`entryOffset` 和 `dirOffset` 随机化，保证遍历顺序不确定
2. **table 增长/分裂**：迭代器持有旧 table 引用，用旧 table 决定访问哪些 key，再用新 table 查找当前值（满足更新/删除语义）
3. **directory 增长**：通过 `globalDepth` 差值调整 `dirIdx`（`dirIdx <<= orders`）
4. **Clear 检测**：通过 `clearSeq` 检测迭代期间是否发生过 Clear
5. **跳过重复 table**：同一 table 在 directory 中占多个索引时，用 `localDepth` 计算跳过

### 清空与克隆

- **Clear**：遍历所有 table 调用 `table.Clear`，重置 `used` 和 `tombstonePossible`，递增 `clearSeq`，重置 `seed`（防碰撞攻击）
- **Clone**：浅拷贝 Map 结构体，深拷贝 dirPtr（逐 table clone），`maps.Clone()` 使用此功能

### 需谨慎保留的旧结论

- 原页称 `ctrlGroup` 空标记承担 GC 优化，缺少证据，不能据此断言 GC 扫描机制。
- 原页称旧版没有小 map 优化，新版独有；本次没有核验旧源码，不沿用该绝对比较。
- 原页将墓碑描述为替代 overflow 概念，两者用途不同：墓碑保持删除后的探测正确性，并非独立的 overflow 存储。
- 原页称旧实现每次搬迁两个 bucket，这种固定数量描述尚未核验。
- 原页给出的旧负载因子 6.5 是每 bucket 的阈值；若比较比例，应与 bucket 的 8 slots 相除，不能直接与新实现 7/8 比较。
- `pruneTombstones` 的存在、触发阈值及实验开关必须固定 Go revision 后核验。

## 参考来源

- 本地原始视频字幕：[[How Go Maps Work Internally - Swiss Tables]]。
- 视频入口：[[Go map Swiss Tables 视频学习笔记]]。
- 原笔记参考链接：[Abseil Swiss Tables 设计](https://abseil.io/about/design/swisstables)，本次未重新读取，不作为新增 claim 的证据。
- 原笔记参考链接：[Go maps/map.go 浮动 master](https://go.googlesource.com/go/+/refs/heads/master/src/internal/runtime/maps/map.go)，非固定版本，不能据此认定片段与当前源码一致。

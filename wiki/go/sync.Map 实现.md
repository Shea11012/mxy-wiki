---
title: "sync.Map 实现"
created: 2026-04-23
updated: 2026-07-22
status: needs-review
type: concept
tags:
  - go
  - go/runtime
confidence: INFERRED
---
## 结构

Go 1.24 起，`sync.Map` 的底层实现从"读写分离双 map"彻底改为**并发 Hash Trie**（`internal/sync.HashTrieMap`）。`sync.Map` 本身只是 `HashTrieMap[any, any]` 的薄包装。

```go
// sync/map.go
type Map struct {
    _ noCopy
    m isync.HashTrieMap[any, any]
}

// internal/sync/hashtriemap.go
type HashTrieMap[K comparable, V any] struct {
    inited   atomic.Uint32
    initMu   Mutex
    root     atomic.Pointer[indirect[K, V]]
    keyHash  hashFunc    // 从 map 类型元信息中获取的哈希函数
    valEqual equalFunc   // 值的相等比较函数（用于 CompareAndSwap/Delete）
    seed     uintptr     // 随机种子，防止 hash 攻击
}
```

核心数据结构是一个 **16 叉 Hash Trie**，每个内部节点（indirect）最多有 16 个子节点：

```go
const (
    nChildrenLog2 = 4           // 每层用 4 bit 索引
    nChildren     = 1 << 4      // 16 个子节点
    nChildrenMask = nChildren - 1
)

type indirect[K comparable, V any] struct {
    node[K, V]
    dead     atomic.Bool                         // 标记节点是否已被移除
    mu       Mutex                               // 保护 children 的变更
    parent   *indirect[K, V]                     // 父节点指针，用于删除时向上修剪
    children [nChildren]atomic.Pointer[node[K, V]]
}

type entry[K comparable, V any] struct {
    node[K, V]
    overflow atomic.Pointer[entry[K, V]]         // 哈希碰撞时的溢出链表
    key      K
    value    V
}

type node[K comparable, V any] struct {
    isEntry bool  // 区分 indirect 节点和 entry 节点
}
```

**设计思路**：通过哈希值的每 4 bit 逐层向下定位，形成一个树高最多 `PtrSize * 8 / 4` 层的 trie。读操作无需加锁，沿 atomic pointer 逐层查找；写操作只锁定目标 indirect 节点，不同子树上的写操作互不阻塞。

## 初始化

`HashTrieMap` 使用延迟初始化（lazy init），首次操作时才创建 root 节点和哈希函数：

```go
func (ht *HashTrieMap[K, V]) init() {
    if ht.inited.Load() == 0 {
        ht.initSlow()
    }
}

func (ht *HashTrieMap[K, V]) initSlow() {
    ht.initMu.Lock()
    defer ht.initMu.Unlock()
    if ht.inited.Load() != 0 {
        return
    }
    var m map[K]V
    mapType := abi.TypeOf(m).MapType()
    ht.root.Store(newIndirectNode[K, V](nil))
    ht.keyHash = mapType.Hasher   // 复用 runtime map 的哈希函数
    ht.valEqual = mapType.Elem.Equal
    ht.seed = uintptr(runtime_rand())
    ht.inited.Store(1)
}
```

## Load

读操作完全无锁，沿 trie 树按哈希逐层遍历：

```go
func (ht *HashTrieMap[K, V]) Load(key K) (value V, ok bool) {
    ht.init()
    hash := ht.keyHash(abi.NoEscape(unsafe.Pointer(&key)), ht.seed)

    i := ht.root.Load()
    hashShift := 8 * goarch.PtrSize
    for hashShift != 0 {
        hashShift -= nChildrenLog2

        n := i.children[(hash>>hashShift)&nChildrenMask].Load()
        if n == nil {
            return *new(V), false
        }
        if n.isEntry {
            return n.entry().lookup(key)   // 在 entry 的 overflow 链上查找
        }
        i = n.indirect()                   // 继续向下
    }
    panic("ran out of hash bits while iterating")
}
```

**流程**：从 root 开始，每次取哈希值的高 4 bit 作为索引访问 children，遇到 nil 直接返回未找到，遇到 entry 节点则在 overflow 链上查找 key。

## Store / Swap

Store 底层调用 Swap，Swap 的核心逻辑：

```go
func (ht *HashTrieMap[K, V]) Swap(key K, new V) (previous V, loaded bool) {
    for {
        // 无锁遍历找到插入位置（entry 节点或 nil 槽位）
        i, hashShift, slot, n := /* 遍历 trie */

        // 加锁并 double-check
        i.mu.Lock()
        n = slot.Load()
        if (n == nil || n.isEntry) && !i.dead.Load() {
            break  // 确认状态未变，继续操作
        }
        i.mu.Unlock()  // 状态已变，重新遍历
    }
    defer i.mu.Unlock()

    if n != nil {
        // entry 已存在，检查 key 是否相同
        oldEntry := n.entry()
        newEntry, old, swapped := oldEntry.swap(key, new)
        if swapped {
            slot.Store(&newEntry.node)  // 原子替换 entry
            return old, true
        }
        // key 不同，需要 expand
    }
    // 插入新 entry
    newEntry := newEntryNode(key, new)
    if oldEntry == nil {
        slot.Store(&newEntry.node)
    } else {
        slot.Store(ht.expand(oldEntry, newEntry, hash, hashShift, i))
    }
    return zero, false
}
```

**流程**：
1. 无锁遍历 trie 找到目标位置（nil 槽或 entry 节点）
2. 锁定该 indirect 节点并 double-check（防止遍历期间结构已变）
3. 如果 entry 存在且 key 匹配，替换 value；如果 key 不匹配（哈希碰撞到同一槽位但 key 不同），调用 `expand` 在两者之间插入新的 indirect 节点

### expand — 处理插入冲突

当两个不同 key 的哈希在高位相同、需要共存于同一路径时，`expand` 递归创建 indirect 节点直到两者的哈希在某层分叉：

```go
func (ht *HashTrieMap[K, V]) expand(oldEntry, newEntry *entry[K, V], ...) *node[K, V] {
    oldHash := ht.keyHash(unsafe.Pointer(&oldEntry.key), ht.seed)
    if oldHash == newHash {
        // 真正的哈希碰撞：用 overflow 链表解决
        newEntry.overflow.Store(oldEntry)
        return &newEntry.node
    }
    // 沿哈希分叉位逐层创建 indirect 节点
    for {
        hashShift -= nChildrenLog2
        oi := (oldHash >> hashShift) & nChildrenMask
        ni := (newHash >> hashShift) & nChildrenMask
        if oi != ni {
            // 分叉点：两个 entry 分别放入不同子槽
            newIndirect.children[oi].Store(&oldEntry.node)
            newIndirect.children[ni].Store(&newEntry.node)
            break
        }
        // 还未分叉，继续创建更深层的 indirect 节点
        nextIndirect := newIndirectNode(newIndirect)
        newIndirect.children[oi].Store(&nextIndirect.node)
        newIndirect = nextIndirect
    }
    return &top.node
}
```

## LoadAndDelete

删除操作需要加锁，删除后向上修剪空的 indirect 节点：

```go
func (ht *HashTrieMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
    ht.init()
    hash := ht.keyHash(...)

    // find：无锁遍历 + 加锁 double-check
    i, hashShift, slot, n := ht.find(key, hash, nil, *new(V))
    if n == nil {
        return *new(V), false
    }

    v, e, loaded := n.entry().loadAndDelete(key)
    if !loaded {
        return *new(V), false
    }
    if e != nil {
        // 仅删除了 overflow 链中的一个元素，替换 entry 头
        slot.Store(&e.node)
        i.mu.Unlock()
        return v, true
    }
    // 删除整个 entry
    slot.Store(nil)

    // 向上修剪空的 indirect 节点
    for i.parent != nil && i.empty() {
        hashShift += nChildrenLog2
        parent := i.parent
        parent.mu.Lock()
        i.dead.Store(true)
        parent.children[(hash>>hashShift)&nChildrenMask].Store(nil)
        i.mu.Unlock()
        i = parent
    }
    i.mu.Unlock()
    return v, true
}
```

**流程**：
1. 找到目标 entry，加锁
2. 从 overflow 链中移除匹配的 key
3. 如果 entry 被完全移除，将槽位置 nil
4. 从当前 indirect 节点向上检查，若节点变空则一并移除（设置 `dead` 标记防止并发访问已删除节点）

## 与旧版的核心区别

| 方面 | 旧版（Go ≤1.23） | 新版（Go 1.24+） |
|------|------------------|-------------------|
| 数据结构 | 双 map 读写分离（`read` + `dirty`） | 16 叉并发 Hash Trie |
| 删除方式 | 标记 `expunged`，延迟清理 | 直接移除，向上修剪空节点 |
| entry 三态 | nil / expunged / 有效指针 | 直接存 key-value，无三态 |
| 锁粒度 | 全局一把 `Mutex` | 每个 indirect 节点独立 `Mutex` |
| miss 机制 | `misses` 计数触发 dirty → read 提升 | 无 miss 机制，直接在 trie 上操作 |
| 哈希冲突 | 依赖 Go map 原生处理 | overflow 链表 + trie 分叉 |
| 泛型 | 无 | `HashTrieMap[K comparable, V any]` |
| 新增方法 | — | `Clear`, `CompareAndSwap`, `CompareAndDelete` |

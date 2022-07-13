---
date created: 2021-11-30 21:22
date modified: 2022-03-27 18:56
title: lru
---
## LRU(Least Recently Used)
LRU，如果数据最近被访问过，那么将来访问的频率也会变高。
实现：
- 使用 map 存储键和值的映射关系，这样查找效率为 O(1)
- 实现一个双向链表，当访问到某个值是将其移到队尾，删除元素则从队头开始，这样时间复杂度为 O(1)
```go
package lru

type Node struct {
	prev *Node
	next *Node
	key  string
	val  interface{}
}

type DoubleList struct {
	head *Node
	tail *Node
}

func NewDoubleList() *DoubleList {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return &DoubleList{
		head: head,
		tail: tail,
	}
}

func (l *DoubleList) AddLast(node *Node) {
	node.next = l.tail
	node.prev = l.tail.prev
	l.tail.prev.next = node
	l.tail.prev = node
}

func (l *DoubleList) Remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (l *DoubleList) RemoveFirst() *Node {
	if l.head.next == l.tail {
		return nil
	}

	first := l.head.next
	l.Remove(first)
	return first
}

type LRUCache struct {
	cache *DoubleList
	m     map[string]*Node
	size  int
}

func NewCache(size int) *LRUCache {
	return &LRUCache{
		size:  size,
		m:     make(map[string]*Node, size),
		cache: NewDoubleList(),
	}
}

func (c *LRUCache) Put(key string, val interface{}) {
	n, ok := c.m[key]
	if ok {
		c.makeRecently(n)
		return
	}

	if len(c.m) == c.size {
		c.removeFirst()
	}

	n = &Node{
		key: key,
		val: val,
	}

	c.addLast(n)
}

func (c *LRUCache) Get(key string) interface{} {
	v, ok := c.m[key]
	if !ok {
		return nil
	}

	c.makeRecently(v)
	return v
}

// 删除任意节点
func (c *LRUCache) remove(n *Node) {
	c.cache.Remove(n)
	delete(c.m, n.key)
}

// 删除头部节点
func (c *LRUCache) removeFirst() {
	n := c.cache.RemoveFirst()
	delete(c.m, n.key)
}

// 将节点追加至尾部
func (c *LRUCache) addLast(n *Node) {
	c.cache.AddLast(n)
	c.m[n.key] = n
}

// 将最近访问的节点放至尾部
func (c *LRUCache) makeRecently(n *Node) {
	c.remove(n)
	c.addLast(n)
}
```
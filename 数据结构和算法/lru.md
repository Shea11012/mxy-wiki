## lru
```go
type Node struct {
	key  string
	val  int
	pre  *Node
	next *Node
}

func (n Node) String() string {
	return fmt.Sprintf("key: %s, val: %d\n",n.key,n.val)
}

type LRUCache struct {
	head  *Node
	end   *Node
	limit int
	m     map[string]*Node
}

func NewLRUCache(limit int) *LRUCache {
	return &LRUCache{
		limit: limit,
		m:     map[string]*Node{},
	}
}

func (l *LRUCache) Get(key string) *Node {
	node, ok := l.m[key]
	if !ok {
		return nil
	}

	l.refreshNode(node)
	return node
}

func (l *LRUCache) Put(key string, value int) {
	node, ok := l.m[key]
	if !ok {
		if len(l.m) >= l.limit {
			oldKey := l.removeNode(l.head)
			delete(l.m, oldKey.key)
		}

		node = &Node{
			key: key,
			val: value,
		}

		l.addNode(node)
		l.m[key] = node
	}

	node.val = value
	l.refreshNode(node)
}

func (l *LRUCache) Remove(key string) {
	node, ok := l.m[key]
	if !ok {
		return
	}

	l.removeNode(node)
	delete(l.m, key)
}

func (l *LRUCache) refreshNode(node *Node) {
	if node == l.end {
		return
	}

	// 移除节点
	l.removeNode(node)
	// 重新插入节点
	l.addNode(node)
}

func (l *LRUCache) removeNode(node *Node) *Node {
	if node == l.head && node == l.end {
		l.head = nil
		l.end = nil
	} else if node == l.end {
		l.end = l.end.pre
		l.end.next = nil
	} else if node == l.head {
		l.head = l.head.next
		l.head.pre = nil
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
	}

	return node
}

func (l *LRUCache) addNode(node *Node) {
	if l.end != nil {
		l.end.next = node
		node.pre = l.end
	}

	l.end = node
	if l.head == nil {
		l.head = node
	}
}
```
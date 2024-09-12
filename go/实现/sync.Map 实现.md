---
tags: []
date created: 2022-04-04 04:53
date modified: 2023-04-13 02:33
go version: 1.20
---

## 结构

```go
// sync/map.go

type Map struct {
	mu Mutex

	// 只读的map，通过原子更新，删除entry时需要加锁
	read atomic.Value // readOnly

	// 需要加锁才能访问，包含了所有在read中的但未被删除的元素以及新加的元素
	dirty map[interface{}]*entry

	// 记录从read中读取的miss次数，一旦miss数和dirty长度一样，就会把dirty提升为read
	misses int
}

type readOnly struct {
	m       map[interface{}]*entry
	amended bool // 当dirty中包含read没有的数据时为true
}

// 标识此项是已经删掉的指针，当map中的一个项目被删除了，只是把它的值标记为expunged，以后才有机会真正删除
var expunged = unsafe.Pointer(new(interface{}))

// 表示一个值
type entry struct {
	p unsafe.Pointer // *interface{}
}
```

## store

```go
func (m *Map) Store(key, value interface{}) {
	read, _ := m.read.Load().(readOnly)
	// 如果read中包含此项，则表示更新，cas更新值
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	// read中不存在，或者值是处于expunged，则加锁访问dirty
	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	// 再次检查read
	if e, ok := read.m[key]; ok {
		// 将expunged状态转为nil
		if e.unexpungeLocked() {
			// 将值存入dirty中
			m.dirty[key] = e
		}
		// 更新
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {	// dirty 中存在直接更新
		e.storeLocked(&value)
	} else {	// 表示是一个新key
		if !read.amended {
			// dirty为nil时，将read中不为expunged的值复制到dirty中
			m.dirtyLocked()
			// 
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		// 将新值放到dirty中
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}
```

## load

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// key不存在且dirty包含read不存在的值
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		// 再次检查read中的key
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 从dirty中读取，miss数加1，如果miss数>= dirty数量则将dirty提升为read，dirty置为nil
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	
	return e.load()
}
```

## delete

```go
func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		// 双检查逻辑与load一致
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete()
	}
	return nil, false
}

// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	m.LoadAndDelete(key)
}

func (e *entry) delete() (value interface{}, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return nil, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*interface{})(p), true
		}
	}
}
```
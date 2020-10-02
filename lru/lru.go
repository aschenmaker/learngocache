package lru

import (
	"container/list"
)

// 使用字典map+双向链表来实现队列。
// 对于并发来说并不是安全的。
// Cache is made by LRU
type Cache struct {
	// 最大允许使用内存
	maxBytes int64
	// 已经使用内存
	nBytes int64
	ll     *list.List
	cache  map[string]*list.Element
	// 可选，
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		// 执行回收时，毁掉函数
		OnEvicted: onEvicted,
	}
}

// k-v 将访问过的节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// Remove the oldest item
func (c *Cache) RemoveOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 存在回调函数，则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add a value to the cache.
func (c *Cache) Add(key string, value Value) {
	// 如果键值对存在则更新，并移动到队尾
	// 如果不存在则，添加新的节点，建立映射
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := c.ll.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 队列满则清理缓存
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// return the len of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

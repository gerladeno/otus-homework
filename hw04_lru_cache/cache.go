package hw04_lru_cache //nolint:golint,stylecheck
import (
	"sync"
)

type Key string

var mx sync.Mutex

type Cache interface {
	Get(Key) (interface{}, bool)
	Set(Key, interface{}) bool
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = map[Key]*ListItem{}
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	mx.Lock()
	defer mx.Unlock()
	if val, ok := c.items[key]; ok {
		c.queue.MoveToFront(val)
		return val.Value.(cacheItem).Value, true
	}
	return nil, false
}

func (c *lruCache) Set(key Key, i interface{}) bool {
	mx.Lock()
	defer mx.Unlock()
	if _, ok := c.items[key]; ok {
		c.items[key].Value = cacheItem{
			Key:   key,
			Value: i,
		}
		c.queue.MoveToFront(c.items[key])
		return true
	}
	if c.queue.Len() >= c.capacity {
		delete(c.items, c.queue.Back().Value.(cacheItem).Key)
		c.queue.Remove(c.queue.Back())
	}
	c.items[key] = c.queue.PushFront(cacheItem{
		Key:   key,
		Value: i,
	})
	return false
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    map[Key]*ListItem{},
	}
}

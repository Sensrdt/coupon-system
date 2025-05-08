package cache

import (
	"container/list"
	"sync"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

type lruCache struct {
	mu    sync.RWMutex
	cache map[string]interface{}
	queue *list.List
}

func NewLRU(capacity int) Cache {
	return &lruCache{
		cache: make(map[string]interface{}),
		queue: list.New(),
	}
}

func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.cache[key]
	return value, ok
}

func (c *lruCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
	c.queue.PushFront(key)

	if c.queue.Len() > 100 {
		last := c.queue.Back()
		delete(c.cache, last.Value.(string))
		c.queue.Remove(last)
	}
}

func (c *lruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

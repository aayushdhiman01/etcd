package cache

import (
	"container/list"
	"sync"
	"time"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
)

// CacheEntry represents a cached watch response
type CacheEntry struct {
	key       string
	timestamp time.Time
	events    []*pb.Event
}

// LRUCache implements an LRU cache for etcd watch responses
type LRUCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*list.Element
	lruList  *list.List
}

// NewLRUCache initializes the cache
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		lruList:  list.New(),
	}
}

// Get retrieves cached watch events for a key
func (c *LRUCache) Get(key string) ([]*pb.Event, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.lruList.MoveToFront(elem)
		entry := elem.Value.(*CacheEntry)
		return entry.events, true
	}
	return nil, false
}

// Put stores new watch events in the cache
func (c *LRUCache) Put(key string, events []*pb.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.items[key]; found {
		c.lruList.MoveToFront(elem)
		elem.Value.(*CacheEntry).events = events
		return
	}

	entry := &CacheEntry{key: key, timestamp: time.Now(), events: events}
	elem := c.lruList.PushFront(entry)
	c.items[key] = elem

	if len(c.items) > c.capacity {
		oldest := c.lruList.Back()
		if oldest != nil {
			delete(c.items, oldest.Value.(*CacheEntry).key)
			c.lruList.Remove(oldest)
		}
	}
}

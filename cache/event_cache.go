package cache

import (
	"log"
	"sync"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
)

// EventCache stores watch events for keys
type EventCache struct {
	mu    sync.Mutex
	cache *LRUCache
}

// NewEventCache initializes the event cache
func NewEventCache(capacity int) *EventCache {
	return &EventCache{
		cache: NewLRUCache(capacity),
	}
}

// StoreWatchEvents caches events for a given key
func (ec *EventCache) StoreWatchEvents(key string, events []*pb.Event) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	log.Printf("Caching events for key: %s\n", key)
	ec.cache.Put(key, events)
}

// GetCachedEvents retrieves cached watch events
func (ec *EventCache) GetCachedEvents(key string) ([]*pb.Event, bool) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	return ec.cache.Get(key)
}

package cache

import "sync"

// KVCache caches key-value pairs
type KVCache struct {
	mu    sync.Mutex
	store map[string]string
}

// NewKVCache initializes a new KV cache
func NewKVCache() *KVCache {
	return &KVCache{store: make(map[string]string)}
}

// Get retrieves a value
func (kvc *KVCache) Get(key string) (string, bool) {
	kvc.mu.Lock()
	defer kvc.mu.Unlock()
	val, found := kvc.store[key]
	return val, found
}

// Put stores a key-value pair
func (kvc *KVCache) Put(key, value string) {
	kvc.mu.Lock()
	defer kvc.mu.Unlock()
	kvc.store[key] = value
}

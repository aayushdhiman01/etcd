package cache

import pb "go.etcd.io/etcd/api/v3/etcdserverpb"

// WatchCacheAPI defines the interface for watch cache interactions
type WatchCacheAPI interface {
	GetCachedEvents(key string) ([]*pb.Event, bool)
	StoreWatchEvents(key string, events []*pb.Event)
}

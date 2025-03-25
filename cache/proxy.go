package cache

import (
	"context"
	"log"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
)

// WatchProxy intercepts and caches etcd watch requests
type WatchProxy struct {
	eventCache *EventCache
	watchC     pb.WatchClient
}

// NewWatchProxy initializes the proxy
func NewWatchProxy(client pb.WatchClient, cacheSize int) *WatchProxy {
	return &WatchProxy{
		eventCache: NewEventCache(cacheSize),
		watchC:     client,
	}
}

// Watch intercepts watch requests, serves from cache if available
func (wp *WatchProxy) Watch(ctx context.Context, req *pb.WatchRequest) (pb.Watch_WatchClient, error) {
	key := string(req.Key)

	if events, found := wp.eventCache.GetCachedEvents(key); found {
		log.Println("Serving watch request from cache:", key)
		return newFakeWatchStream(events), nil
	}

	log.Println("Forwarding watch request to etcd:", key)
	stream, err := wp.watchC.Watch(ctx, req)
	if err != nil {
		return nil, err
	}

	go wp.storeWatchEvents(key, stream)
	return stream, nil
}

// storeWatchEvents caches new watch events
func (wp *WatchProxy) storeWatchEvents(key string, stream pb.Watch_WatchClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Println("Watch error:", err)
			return
		}
		wp.eventCache.StoreWatchEvents(key, resp.Events)
	}
}

package client

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// CacheClient represents a wrapper around etcd client
type CacheClient struct {
	etcdClient *clientv3.Client
}

// NewCacheClient initializes a new CacheClient
func NewCacheClient(endpoints []string) (*CacheClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &CacheClient{etcdClient: cli}, nil
}

// Put sets a key-value pair in etcd
func (c *CacheClient) Put(key, value string) error {
	_, err := c.etcdClient.Put(context.Background(), key, value)
	return err
}

// Get retrieves a value from etcd
func (c *CacheClient) Get(key string) (string, error) {
	resp, err := c.etcdClient.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// Watch starts watching a key and prints events
func (c *CacheClient) Watch(key string) {
	watchChan := c.etcdClient.Watch(context.Background(), key)
	log.Println("Watching key:", key)

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			log.Printf("Event received: %s %q : %q\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}

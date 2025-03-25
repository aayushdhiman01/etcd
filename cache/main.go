package cache

import (
	"fmt"
	"log"

	"go.etcd.io/cache/client"
)

func main() {
	cacheClient, err := client.NewCacheClient([]string{"localhost:2379"})
	if err != nil {
		log.Fatal(err)
	}

	err = cacheClient.Put("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	value, err := cacheClient.Get("foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Value:", value)
}

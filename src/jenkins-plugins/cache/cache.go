package cache

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

//Cache is a key value store.
type Cache struct {
	cache map[string]*cacheShard
}
type cacheShard struct {
	items map[string][]byte
	lock  *sync.RWMutex
}

//New creates a empty cache, important to premake right size of cache base. 16^3=4096, 16^2=256
func New() Cache {
	c := make(map[string]*cacheShard, 4096)
	for i := 0; i < 4096; i++ {
		c[fmt.Sprintf("%03x", i)] = &cacheShard{
			items: make(map[string][]byte, 2048),
			lock:  new(sync.RWMutex),
		}
	}

	return Cache{cache: c}
}

//Get object from cache by using key.
func (c Cache) Get(key string) []byte {
	shard := c.getShard(key)
	shard.lock.RLock()
	defer shard.lock.RUnlock()
	return shard.items[key]
}

//Set a object in cache by using key, data.
func (c Cache) Set(key string, data []byte) {
	shard := c.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	shard.items[key] = data
}

func (c Cache) getShard(key string) (shard *cacheShard) {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	shardKey := fmt.Sprintf("%x", hasher.Sum(nil))[0:3]
	return c.cache[shardKey]
}

//Size returns the number of items in the cache.
func (c Cache) Size() (size int) {
	size = 0
	for _, shard := range c.cache {
		size = size + len(shard.items)
	}
	return size
}

//GetAll returns all keys in cache.
func (c Cache) GetAll() []string {
	var list []string
	for _, shard := range c.cache {
		for k := range shard.items {
			list = append(list, k)
		}
	}
	return list
}

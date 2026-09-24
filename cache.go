package uaparser

import (
	"container/list"
	"hash/fnv"
	"sync"
)

// Cacher defines the cache interface for parsed User-Agent results.
type Cacher interface {
	Get(key string) (*Result, bool)
	Set(key string, val *Result)
}

// ShardedLRU is a concurrent, high-performance sharded LRU cache
// designed to minimize lock contention under heavy multi-threaded workloads.
type ShardedLRU struct {
	numShards uint32
	shards    []*lruShard
}

type lruShard struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]*list.Element
	evict    *list.List
}

type cacheEntry struct {
	key string
	val *Result
}

// NewShardedLRU creates a new sharded LRU cache with the specified total capacity.
// Shards default to 32 to ensure parallel throughput across multiple CPU cores.
func NewShardedLRU(totalCapacity int) *ShardedLRU {
	const shardsCount = 32
	if totalCapacity < shardsCount {
		totalCapacity = shardsCount
	}
	shardCapacity := totalCapacity / shardsCount

	s := &ShardedLRU{
		numShards: uint32(shardsCount),
		shards:    make([]*lruShard, shardsCount),
	}

	for i := 0; i < shardsCount; i++ {
		s.shards[i] = &lruShard{
			capacity: shardCapacity,
			items:    make(map[string]*list.Element, shardCapacity),
			evict:    list.New(),
		}
	}

	return s
}

func (s *ShardedLRU) getShard(key string) *lruShard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	shardIdx := h.Sum32() % s.numShards
	return s.shards[shardIdx]
}

// Get retrieves a result from the cache.
func (s *ShardedLRU) Get(key string) (*Result, bool) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if elem, ok := shard.items[key]; ok {
		shard.evict.MoveToFront(elem)
		return elem.Value.(*cacheEntry).val, true
	}
	return nil, false
}

// Set stores a result in the cache.
func (s *ShardedLRU) Set(key string, val *Result) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if elem, ok := shard.items[key]; ok {
		shard.evict.MoveToFront(elem)
		elem.Value.(*cacheEntry).val = val
		return
	}

	if shard.evict.Len() >= shard.capacity {
		oldest := shard.evict.Back()
		if oldest != nil {
			shard.evict.Remove(oldest)
			kv := oldest.Value.(*cacheEntry)
			delete(shard.items, kv.key)
		}
	}

	entry := &cacheEntry{key: key, val: val}
	elem := shard.evict.PushFront(entry)
	shard.items[key] = elem
}

// NoOpCache is a no-op implementation of Cacher when caching is disabled.
type NoOpCache struct{}

func (NoOpCache) Get(key string) (*Result, bool) { return nil, false }
func (NoOpCache) Set(key string, val *Result)    {}

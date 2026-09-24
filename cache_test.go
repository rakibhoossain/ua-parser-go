package uaparser

import (
	"fmt"
	"sync"
	"testing"
)

func TestShardedLRU(t *testing.T) {
	cache := NewShardedLRU(64)

	res1 := &Result{UA: "test-ua-1", Browser: Browser{Name: "Chrome"}}
	cache.Set("test-ua-1", res1)

	got, ok := cache.Get("test-ua-1")
	if !ok || got.Browser.Name != "Chrome" {
		t.Fatalf("expected Chrome, got %v (ok=%v)", got, ok)
	}

	_, ok = cache.Get("non-existent")
	if ok {
		t.Fatalf("expected not found for non-existent key")
	}
}

func TestShardedLRUConcurrent(t *testing.T) {
	cache := NewShardedLRU(128)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("ua-%d", id%10)
			cache.Set(key, &Result{UA: key})
			val, _ := cache.Get(key)
			if val != nil && val.UA != key {
				t.Errorf("mismatched key: expected %s, got %s", key, val.UA)
			}
		}(i)
	}
	wg.Wait()
}

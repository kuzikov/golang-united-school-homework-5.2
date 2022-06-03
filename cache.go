package cache

import (
	"sync"
	"time"
)

// TODO: clear expired entries
type storage struct {
	val    string
	expire *time.Time
}

// Cache contain concurrent-safe storage for key-value pairs.
type Cache struct {
	mu   sync.Mutex
	data map[string]storage
}

// NewCache create initialized empty Cache.
func NewCache() Cache {
	return Cache{
		mu:   sync.Mutex{},
		data: make(map[string]storage),
	}
}

// Get return data
func (cache *Cache) Get(key string) (string, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	storage, ok := cache.data[key]
	if !ok {
		return "", false
	}

	switch storage.expire {
	case &time.Time{}:
		t := storage.expire.Sub(time.Now())
		if t > time.Nanosecond {
			return storage.val, true
		}
	case nil:
		return storage.val, true
	}

	return "", false
}

// Put ...
func (cache *Cache) Put(key, value string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = storage{
		val:    value,
		expire: nil,
	}
	return
}

// Keys ...
func (cache *Cache) Keys() []string {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	keys := make([]string, 0, len(cache.data))
	now := time.Now()
	for key, entry := range cache.data {

		if entry.expire == nil || entry.expire.Sub(now) > time.Nanosecond {
			keys = append(keys, key)
			continue
		}

		// notExpired :=
		// if notExpired {
		// 	keys = append(keys, key)
		// 	continue
		// }

	}

	return keys
}

// PutTill ...
func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	timeCopy := deadline
	cache.data[key] = storage{
		val:    value,
		expire: &timeCopy,
	}
	return

}

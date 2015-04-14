package physical

import "github.com/hashicorp/golang-lru"

const (
	// DefaultCacheSize is used if no cache size is specified
	// for NewPhysicalCache
	DefaultCacheSize = 32 * 1024
)

// PhysicalCache is used to wrap an underlying physical backend
// and provide an LRU cache layer on top. Most of the reads done by
// Vault are for policy objects so there is a large read reduction
// by using a simple write-through cache.
type PhysicalCache struct {
	backend Backend
	lru     *lru.Cache
}

// NewPhysicalCache returns a physical cache of the given size.
// If no size is provided, the default size is used.
func NewPhysicalCache(b Backend, size int) *PhysicalCache {
	if size <= 0 {
		size = DefaultCacheSize
	}
	cache, _ := lru.New(size)
	c := &PhysicalCache{
		backend: b,
		lru:     cache,
	}
	return c
}

// Purge is used to clear the cache
func (c *PhysicalCache) Purge() {
	c.lru.Purge()
}

func (c *PhysicalCache) Put(entry *Entry) error {
	err := c.backend.Put(entry)
	c.lru.Add(entry.Key, entry)
	return err
}

func (c *PhysicalCache) Get(key string) (*Entry, error) {
	// Check the LRU first
	if raw, ok := c.lru.Get(key); ok {
		if raw == nil {
			return nil, nil
		} else {
			return raw.(*Entry), nil
		}
	}

	// Read from the underlying backend
	ent, err := c.backend.Get(key)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.lru.Add(key, ent)
	return ent, err
}

func (c *PhysicalCache) Delete(key string) error {
	err := c.backend.Delete(key)
	c.lru.Remove(key)
	return err
}

func (c *PhysicalCache) List(prefix string) ([]string, error) {
	// Always pass-through as this would be difficult to cache.
	return c.backend.List(prefix)
}

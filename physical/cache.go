package physical

import (
	"strings"

	"github.com/hashicorp/golang-lru"
	log "github.com/mgutz/logxi/v1"
)

const (
	// DefaultCacheSize is used if no cache size is specified for NewCache
	DefaultCacheSize = 32 * 1024
)

// Cache is used to wrap an underlying physical backend
// and provide an LRU cache layer on top. Most of the reads done by
// Vault are for policy objects so there is a large read reduction
// by using a simple write-through cache.
type Cache struct {
	backend Backend
	lru     *lru.TwoQueueCache
}

// NewCache returns a physical cache of the given size.
// If no size is provided, the default size is used.
func NewCache(b Backend, size int, logger log.Logger) *Cache {
	if size <= 0 {
		size = DefaultCacheSize
	}
	if logger.IsTrace() {
		logger.Trace("physical/cache: creating LRU cache", "size", size)
	}
	cache, _ := lru.New2Q(size)
	c := &Cache{
		backend: b,
		lru:     cache,
	}
	return c
}

// Purge is used to clear the cache
func (c *Cache) Purge() {
	c.lru.Purge()
}

func (c *Cache) Put(entry *Entry) error {
	err := c.backend.Put(entry)
	if err == nil {
		c.lru.Add(entry.Key, entry)
	}
	return err
}

func (c *Cache) Get(key string) (*Entry, error) {
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

	// Cache the result. We do NOT cache negative results
	// for keys in the 'core/' prefix otherwise we risk certain
	// race conditions upstream. The primary issue is with the HA mode,
	// we could potentially negatively cache the leader entry and cause
	// leader discovery to fail.
	if ent != nil || !strings.HasPrefix(key, "core/") {
		c.lru.Add(key, ent)
	}
	return ent, err
}

func (c *Cache) Delete(key string) error {
	err := c.backend.Delete(key)
	if err == nil {
		c.lru.Remove(key)
	}
	return err
}

func (c *Cache) List(prefix string) ([]string, error) {
	// Always pass-through as this would be difficult to cache.
	return c.backend.List(prefix)
}

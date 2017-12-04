package physical

import (
	"strings"

	"github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/helper/locksutil"
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
	locks   []*locksutil.LockEntry
	logger  log.Logger
}

// TransactionalCache is a Cache that wraps the physical that is transactional
type TransactionalCache struct {
	*Cache
	Transactional
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
		locks:   locksutil.CreateLocks(),
		logger:  logger,
	}

	return c
}

func NewTransactionalCache(b Backend, size int, logger log.Logger) *TransactionalCache {
	c := &TransactionalCache{
		Cache:         NewCache(b, size, logger),
		Transactional: b.(Transactional),
	}
	return c
}

// Purge is used to clear the cache
func (c *Cache) Purge() {
	// Lock the world
	for _, lock := range c.locks {
		lock.Lock()
		defer lock.Unlock()
	}

	c.lru.Purge()
}

func (c *Cache) Put(entry *Entry) error {
	lock := locksutil.LockForKey(c.locks, entry.Key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Put(entry)
	if err == nil && !strings.HasPrefix(entry.Key, "core/") {
		c.lru.Add(entry.Key, entry)
	}
	return err
}

func (c *Cache) Get(key string) (*Entry, error) {
	lock := locksutil.LockForKey(c.locks, key)
	lock.RLock()
	defer lock.RUnlock()

	// We do NOT cache negative results for keys in the 'core/' prefix
	// otherwise we risk certain race conditions upstream. The primary issue is
	// with the HA mode, we could potentially negatively cache the leader entry
	// and cause leader discovery to fail.
	if strings.HasPrefix(key, "core/") {
		return c.backend.Get(key)
	}

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
	if ent != nil {
		c.lru.Add(key, ent)
	}

	return ent, nil
}

func (c *Cache) Delete(key string) error {
	lock := locksutil.LockForKey(c.locks, key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Delete(key)
	if err == nil && !strings.HasPrefix(key, "core/") {
		c.lru.Remove(key)
	}
	return err
}

func (c *Cache) List(prefix string) ([]string, error) {
	// Always pass-through as this would be difficult to cache. For the same
	// reason we don't lock as we can't reasonably know which locks to readlock
	// ahead of time.
	return c.backend.List(prefix)
}

func (c *TransactionalCache) Transaction(txns []*TxnEntry) error {
	// Lock the world
	for _, lock := range c.locks {
		lock.Lock()
		defer lock.Unlock()
	}

	if err := c.Transactional.Transaction(txns); err != nil {
		return err
	}

	for _, txn := range txns {
		switch txn.Operation {
		case PutOperation:
			c.lru.Add(txn.Entry.Key, txn.Entry)
		case DeleteOperation:
			c.lru.Remove(txn.Entry.Key)
		}
	}

	return nil
}

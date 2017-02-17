package physical

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/strutil"
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
	backend       Backend
	transactional Transactional
	lru           *lru.TwoQueueCache
	locks         map[string]*sync.RWMutex
	logger        log.Logger
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
		locks:   make(map[string]*sync.RWMutex, 256),
		logger:  logger,
	}
	if err := locksutil.CreateLocks(c.locks, 256); err != nil {
		logger.Error("physical/cache: error creating locks", "error", err)
		return nil
	}

	if txnl, ok := c.backend.(Transactional); ok {
		c.transactional = txnl
	}

	return c
}

func (c *Cache) lockHashForKey(key string) string {
	hf := sha1.New()
	hf.Write([]byte(key))
	return strings.ToLower(hex.EncodeToString(hf.Sum(nil))[:2])
}

func (c *Cache) lockForKey(key string) *sync.RWMutex {
	return c.locks[c.lockHashForKey(key)]
}

// Purge is used to clear the cache
func (c *Cache) Purge() {
	// Lock the world
	lockHashes := make([]string, 0, len(c.locks))
	for hash := range c.locks {
		lockHashes = append(lockHashes, hash)
	}

	// Sort and deduplicate. This ensures we don't try to grab the same lock
	// twice, and enforcing a sort means we'll not have multiple goroutines
	// deadlock by acquiring in different orders.
	lockHashes = strutil.RemoveDuplicates(lockHashes)

	for _, lockHash := range lockHashes {
		lock := c.locks[lockHash]
		lock.Lock()
		defer lock.Unlock()
	}

	c.lru.Purge()
}

func (c *Cache) Put(entry *Entry) error {
	lock := c.lockForKey(entry.Key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Put(entry)
	if err == nil {
		c.lru.Add(entry.Key, entry)
	}
	return err
}

func (c *Cache) Get(key string) (*Entry, error) {
	lock := c.lockForKey(key)
	lock.RLock()
	defer lock.RUnlock()

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
	lock := c.lockForKey(key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Delete(key)
	if err == nil {
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

func (c *Cache) Transaction(txns []TxnEntry) error {
	if c.transactional == nil {
		return fmt.Errorf("physical/cache: underlying backend does not support transactions")
	}

	var lockHashes []string
	for _, txn := range txns {
		lockHashes = append(lockHashes, c.lockHashForKey(txn.Entry.Key))
	}

	// Sort and deduplicate. This ensures we don't try to grab the same lock
	// twice, and enforcing a sort means we'll not have multiple goroutines
	// deadlock by acquiring in different orders.
	lockHashes = strutil.RemoveDuplicates(lockHashes)

	for _, lockHash := range lockHashes {
		lock := c.locks[lockHash]
		lock.Lock()
		defer lock.Unlock()
	}

	if err := c.transactional.Transaction(txns); err != nil {
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

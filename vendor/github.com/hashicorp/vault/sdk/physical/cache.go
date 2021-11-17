package physical

import (
	"context"
	"sync/atomic"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/pathmanager"
)

const (
	// DefaultCacheSize is used if no cache size is specified for NewCache
	DefaultCacheSize = 128 * 1024

	// refreshCacheCtxKey is a ctx value that denotes the cache should be
	// refreshed during a Get call.
	refreshCacheCtxKey = "refresh_cache"
)

// These paths don't need to be cached by the LRU cache. This should
// particularly help memory pressure when unsealing.
var cacheExceptionsPaths = []string{
	"wal/logs/",
	"index/pages/",
	"index-dr/pages/",
	"sys/expire/",
	"core/poison-pill",
	"core/raft/tls",
	"core/license",
}

// CacheRefreshContext returns a context with an added value denoting if the
// cache should attempt a refresh.
func CacheRefreshContext(ctx context.Context, r bool) context.Context {
	return context.WithValue(ctx, refreshCacheCtxKey, r)
}

// cacheRefreshFromContext is a helper to look up if the provided context is
// requesting a cache refresh.
func cacheRefreshFromContext(ctx context.Context) bool {
	r, ok := ctx.Value(refreshCacheCtxKey).(bool)
	if !ok {
		return false
	}
	return r
}

// Cache is used to wrap an underlying physical backend
// and provide an LRU cache layer on top. Most of the reads done by
// Vault are for policy objects so there is a large read reduction
// by using a simple write-through cache.
type Cache struct {
	backend         Backend
	lru             *lru.TwoQueueCache
	locks           []*locksutil.LockEntry
	logger          log.Logger
	enabled         *uint32
	cacheExceptions *pathmanager.PathManager
	metricSink      metrics.MetricSink
}

// TransactionalCache is a Cache that wraps the physical that is transactional
type TransactionalCache struct {
	*Cache
	Transactional
}

// Verify Cache satisfies the correct interfaces
var (
	_ ToggleablePurgemonster = (*Cache)(nil)
	_ ToggleablePurgemonster = (*TransactionalCache)(nil)
	_ Backend                = (*Cache)(nil)
	_ Transactional          = (*TransactionalCache)(nil)
)

// NewCache returns a physical cache of the given size.
// If no size is provided, the default size is used.
func NewCache(b Backend, size int, logger log.Logger, metricSink metrics.MetricSink) *Cache {
	if logger.IsDebug() {
		logger.Debug("creating LRU cache", "size", size)
	}
	if size <= 0 {
		size = DefaultCacheSize
	}

	pm := pathmanager.New()
	pm.AddPaths(cacheExceptionsPaths)

	cache, _ := lru.New2Q(size)
	c := &Cache{
		backend: b,
		lru:     cache,
		locks:   locksutil.CreateLocks(),
		logger:  logger,
		// This fails safe.
		enabled:         new(uint32),
		cacheExceptions: pm,
		metricSink:      metricSink,
	}
	return c
}

func NewTransactionalCache(b Backend, size int, logger log.Logger, metricSink metrics.MetricSink) *TransactionalCache {
	c := &TransactionalCache{
		Cache:         NewCache(b, size, logger, metricSink),
		Transactional: b.(Transactional),
	}
	return c
}

func (c *Cache) ShouldCache(key string) bool {
	if atomic.LoadUint32(c.enabled) == 0 {
		return false
	}

	return !c.cacheExceptions.HasPath(key)
}

// SetEnabled is used to toggle whether the cache is on or off. It must be
// called with true to actually activate the cache after creation.
func (c *Cache) SetEnabled(enabled bool) {
	if enabled {
		atomic.StoreUint32(c.enabled, 1)
		return
	}
	atomic.StoreUint32(c.enabled, 0)
}

// Purge is used to clear the cache
func (c *Cache) Purge(ctx context.Context) {
	// Lock the world
	for _, lock := range c.locks {
		lock.Lock()
		defer lock.Unlock()
	}

	c.lru.Purge()
}

func (c *Cache) Put(ctx context.Context, entry *Entry) error {
	if entry != nil && !c.ShouldCache(entry.Key) {
		return c.backend.Put(ctx, entry)
	}

	lock := locksutil.LockForKey(c.locks, entry.Key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Put(ctx, entry)
	if err == nil {
		c.lru.Add(entry.Key, entry)
		c.metricSink.IncrCounter([]string{"cache", "write"}, 1)
	}
	return err
}

func (c *Cache) Get(ctx context.Context, key string) (*Entry, error) {
	if !c.ShouldCache(key) {
		return c.backend.Get(ctx, key)
	}

	lock := locksutil.LockForKey(c.locks, key)
	lock.RLock()
	defer lock.RUnlock()

	// Check the LRU first
	if !cacheRefreshFromContext(ctx) {
		if raw, ok := c.lru.Get(key); ok {
			if raw == nil {
				return nil, nil
			}
			c.metricSink.IncrCounter([]string{"cache", "hit"}, 1)
			return raw.(*Entry), nil
		}
	}

	c.metricSink.IncrCounter([]string{"cache", "miss"}, 1)
	// Read from the underlying backend
	ent, err := c.backend.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// Cache the result, even if nil
	c.lru.Add(key, ent)

	return ent, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if !c.ShouldCache(key) {
		return c.backend.Delete(ctx, key)
	}

	lock := locksutil.LockForKey(c.locks, key)
	lock.Lock()
	defer lock.Unlock()

	err := c.backend.Delete(ctx, key)
	if err == nil {
		c.lru.Remove(key)
	}
	return err
}

func (c *Cache) List(ctx context.Context, prefix string) ([]string, error) {
	// Always pass-through as this would be difficult to cache. For the same
	// reason we don't lock as we can't reasonably know which locks to readlock
	// ahead of time.
	return c.backend.List(ctx, prefix)
}

func (c *TransactionalCache) Locks() []*locksutil.LockEntry {
	return c.locks
}

func (c *TransactionalCache) LRU() *lru.TwoQueueCache {
	return c.lru
}

func (c *TransactionalCache) Transaction(ctx context.Context, txns []*TxnEntry) error {
	// Bypass the locking below
	if atomic.LoadUint32(c.enabled) == 0 {
		return c.Transactional.Transaction(ctx, txns)
	}

	// Collect keys that need to be locked
	var keys []string
	for _, curr := range txns {
		keys = append(keys, curr.Entry.Key)
	}
	// Lock the keys
	for _, l := range locksutil.LocksForKeys(c.locks, keys) {
		l.Lock()
		defer l.Unlock()
	}

	if err := c.Transactional.Transaction(ctx, txns); err != nil {
		return err
	}

	for _, txn := range txns {
		if !c.ShouldCache(txn.Entry.Key) {
			continue
		}

		switch txn.Operation {
		case PutOperation:
			c.lru.Add(txn.Entry.Key, txn.Entry)
			c.metricSink.IncrCounter([]string{"cache", "write"}, 1)
		case DeleteOperation:
			c.lru.Remove(txn.Entry.Key)
			c.metricSink.IncrCounter([]string{"cache", "delete"}, 1)
		}
	}

	return nil
}

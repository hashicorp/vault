// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"sync"
	"time"
)

// New creates a cacher.
func New() *Cache {
	return &Cache{
		data: map[string]*cacheEntry{},
	}
}

// Func is the signature for a cache function.
type Func func() (interface{}, error)

// Cache is the internal cache implementation.
type Cache struct {
	lock sync.RWMutex
	data map[string]*cacheEntry
}

// cacheEntry represents an item in the cache with an expiration and lifetime.
type cacheEntry struct {
	result   interface{}
	created  time.Time
	lifetime time.Duration
}

// Fetch retrieves an item from the cache. If the item exists in the cache and
// is within its lifetime, it is returned. If the item does not exist, or if the
// item exists but has exceeded its lifetime, the function f is invoked and the
// result is updated in the cache.
func (c *Cache) Fetch(name string, t time.Duration, f Func) (interface{}, error) {
	// Attempt to read from the cache, returning the cached value if it's still
	// valid.
	c.lock.RLock()
	e, ok := c.data[name]
	if ok && e.result != nil && time.Now().Sub(e.created) < e.lifetime {
		c.lock.RUnlock()
		return e.result, nil
	}
	c.lock.RUnlock()

	// Either no cached value exists, or the cached item has exceeded its lifetime.
	c.lock.Lock()

	// Go doesn't have the ability to "upgrade" a lock, so it's possible that
	// another concurrent invocation sized the lock between our RLock and Lock,
	// thus we have to check again.
	e, ok = c.data[name]
	if ok && e.result != nil && time.Now().Sub(e.created) < e.lifetime {
		c.lock.Unlock()
		return e.result, nil
	}

	result, err := f()
	if err != nil {
		c.lock.Unlock()
		return nil, err
	}

	c.data[name] = &cacheEntry{
		result:   result,
		created:  time.Now(),
		lifetime: t,
	}

	c.lock.Unlock()

	return result, nil
}

// Expire removes the given item from the cache, if it exists.
func (c *Cache) Expire(name string) {
	c.lock.Lock()
	delete(c.data, name)
	c.lock.Unlock()
}

// Clear empties the cache for all values.
func (c *Cache) Clear() {
	c.lock.Lock()
	c.data = map[string]*cacheEntry{}
	c.lock.Unlock()
}

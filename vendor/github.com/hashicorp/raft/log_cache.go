// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"fmt"
	"sync"
)

// LogCache wraps any LogStore implementation to provide an
// in-memory ring buffer. This is used to cache access to
// the recently written entries. For implementations that do not
// cache themselves, this can provide a substantial boost by
// avoiding disk I/O on recent entries.
type LogCache struct {
	store LogStore

	cache []*Log
	l     sync.RWMutex
}

// NewLogCache is used to create a new LogCache with the
// given capacity and backend store.
func NewLogCache(capacity int, store LogStore) (*LogCache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be positive")
	}
	c := &LogCache{
		store: store,
		cache: make([]*Log, capacity),
	}
	return c, nil
}

// IsMonotonic implements the MonotonicLogStore interface. This is a shim to
// expose the underyling store as monotonically indexed or not.
func (c *LogCache) IsMonotonic() bool {
	if store, ok := c.store.(MonotonicLogStore); ok {
		return store.IsMonotonic()
	}

	return false
}

func (c *LogCache) GetLog(idx uint64, log *Log) error {
	// Check the buffer for an entry
	c.l.RLock()
	cached := c.cache[idx%uint64(len(c.cache))]
	c.l.RUnlock()

	// Check if entry is valid
	if cached != nil && cached.Index == idx {
		*log = *cached
		return nil
	}

	// Forward request on cache miss
	return c.store.GetLog(idx, log)
}

func (c *LogCache) StoreLog(log *Log) error {
	return c.StoreLogs([]*Log{log})
}

func (c *LogCache) StoreLogs(logs []*Log) error {
	err := c.store.StoreLogs(logs)
	// Insert the logs into the ring buffer, but only on success
	if err != nil {
		return fmt.Errorf("unable to store logs within log store, err: %q", err)
	}
	c.l.Lock()
	for _, l := range logs {
		c.cache[l.Index%uint64(len(c.cache))] = l
	}
	c.l.Unlock()
	return nil
}

func (c *LogCache) FirstIndex() (uint64, error) {
	return c.store.FirstIndex()
}

func (c *LogCache) LastIndex() (uint64, error) {
	return c.store.LastIndex()
}

func (c *LogCache) DeleteRange(min, max uint64) error {
	// Invalidate the cache on deletes
	c.l.Lock()
	c.cache = make([]*Log, len(c.cache))
	c.l.Unlock()

	return c.store.DeleteRange(min, max)
}

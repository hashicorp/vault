// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"os"
	"sync"
	"time"
)

// cachingFileReader reads a file and keeps an in-memory copy of it, until the
// copy is considered stale. Next ReadFile() after expiry will re-read the file from disk.
type cachingFileReader struct {
	// path is the file path to the cached file.
	path string

	// ttl is the time-to-live duration when cached file is considered stale
	ttl time.Duration

	// cache is the buffer holding the in-memory copy of the file.
	cache cachedFile

	l sync.RWMutex

	// currentTime is a function that returns the current local time.
	// Normally set to time.Now but it can be overwritten by test cases to manipulate time.
	currentTime func() time.Time
}

type cachedFile struct {
	// buf is the buffer holding the in-memory copy of the file.
	buf string

	// expiry is the time when the cached copy is considered stale and must be re-read.
	expiry time.Time
}

func newCachingFileReader(path string, ttl time.Duration, currentTime func() time.Time) *cachingFileReader {
	return &cachingFileReader{
		path:        path,
		ttl:         ttl,
		currentTime: currentTime,
	}
}

func (r *cachingFileReader) ReadFile() (string, error) {
	// Fast path requiring read lock only: file is already in memory and not stale.
	r.l.RLock()
	now := r.currentTime()
	cache := r.cache
	r.l.RUnlock()
	if now.Before(cache.expiry) {
		return cache.buf, nil
	}

	// Slow path: read the file from disk.
	r.l.Lock()
	defer r.l.Unlock()

	buf, err := os.ReadFile(r.path)
	if err != nil {
		return "", err
	}
	r.cache = cachedFile{
		buf:    string(buf),
		expiry: now.Add(r.ttl),
	}

	return r.cache.buf, nil
}

package fileutil

import (
	"os"
	"sync"
	"time"
)

// CachingFileReader reads a file and keeps an in-memory copy of it, until the
// copy is considered stale. Next ReadFile() after expiry will re-read the file from disk.
type CachingFileReader struct {
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
	buf []byte

	// expiry is the time when the cached copy is considered stale and must be re-read.
	expiry time.Time
}

func NewCachingFileReader(path string, ttl time.Duration) *CachingFileReader {
	return &CachingFileReader{
		path:        path,
		ttl:         ttl,
		currentTime: time.Now,
	}
}

func (r *CachingFileReader) ReadFile() ([]byte, error) {
	// Fast path requiring read lock only: file is already in memory and not stale.
	r.l.RLock()
	now := r.currentTime()
	cache := r.cache
	r.l.RUnlock()
	if now.Before(cache.expiry) {
		newBuf := make([]byte, len(cache.buf))
		copy(newBuf, cache.buf)
		return newBuf, nil
	}

	// Slow path: read the file from disk.
	r.l.Lock()
	defer r.l.Unlock()

	buf, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	r.cache = cachedFile{
		buf:    buf,
		expiry: r.currentTime().Add(r.ttl),
	}

	newBuf := make([]byte, len(r.cache.buf))
	copy(newBuf, r.cache.buf)
	return newBuf, nil
}

func (r *CachingFileReader) setStaticTime(staticTime time.Time) {
	r.l.Lock()
	defer r.l.Unlock()
	r.currentTime = func() time.Time {
		return staticTime
	}
}

// Package memorystore defines an in-memory storage system for limiting.
package memorystore

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/internal/fasttime"
)

var _ limiter.Store = (*store)(nil)

type store struct {
	tokens   uint64
	interval time.Duration
	rate     float64

	sweepInterval time.Duration
	sweepMinTTL   uint64

	data     map[string]*bucket
	dataLock sync.RWMutex

	stopped uint32
	stopCh  chan struct{}
}

// Config is used as input to New. It defines the behavior of the storage
// system.
type Config struct {
	// Tokens is the number of tokens to allow per interval. The default value is
	// 1.
	Tokens uint64

	// Interval is the time interval upon which to enforce rate limiting. The
	// default value is 1 second.
	Interval time.Duration

	// SweepInterval is the rate at which to run the garabage collection on stale
	// entries. Setting this to a low value will optimize memory consumption, but
	// will likely reduce performance and increase lock contention. Setting this
	// to a high value will maximum throughput, but will increase the memory
	// footprint. This can be tuned in combination with SweepMinTTL to control how
	// long stale entires are kept. The default value is 6 hours.
	SweepInterval time.Duration

	// SweepMinTTL is the minimum amount of time a session must be inactive before
	// clearing it from the entries. There's no validation, but this should be at
	// least as high as your rate limit, or else the data store will purge records
	// before they limit is applied. The default value is 12 hours.
	SweepMinTTL time.Duration

	// InitialAlloc is the size to use for the in-memory map. Go will
	// automatically expand the buffer, but choosing higher number can trade
	// memory consumption for performance as it limits the number of times the map
	// needs to expand. The default value is 4096.
	InitialAlloc int
}

// New creates an in-memory rate limiter that uses a bucketing model to limit
// the number of permitted events over an interval. It's optimized for runtime
// and memory efficiency.
func New(c *Config) (limiter.Store, error) {
	if c == nil {
		c = new(Config)
	}

	tokens := uint64(1)
	if c.Tokens > 0 {
		tokens = c.Tokens
	}

	interval := 1 * time.Second
	if c.Interval > 0 {
		interval = c.Interval
	}

	sweepInterval := 6 * time.Hour
	if c.SweepInterval > 0 {
		sweepInterval = c.SweepInterval
	}

	sweepMinTTL := 12 * time.Hour
	if c.SweepMinTTL > 0 {
		sweepMinTTL = c.SweepMinTTL
	}

	initialAlloc := 4096
	if c.InitialAlloc > 0 {
		initialAlloc = c.InitialAlloc
	}

	s := &store{
		tokens:   tokens,
		interval: interval,
		rate:     float64(interval) / float64(tokens),

		sweepInterval: sweepInterval,
		sweepMinTTL:   uint64(sweepMinTTL),

		data:   make(map[string]*bucket, initialAlloc),
		stopCh: make(chan struct{}),
	}
	go s.purge()
	return s, nil
}

// Take attempts to remove a token from the named key. If the take is
// successful, it returns true, otherwise false. It also returns the configured
// limit, remaining tokens, and reset time.
func (s *store) Take(key string) (uint64, uint64, uint64, bool) {
	// If the store is stopped, all requests are rejected.
	if atomic.LoadUint32(&s.stopped) == 1 {
		return 0, 0, 0, false
	}

	// Acquire a read lock first - this allows other to concurrently check limits
	// without taking a full lock.
	s.dataLock.RLock()
	if b, ok := s.data[key]; ok {
		s.dataLock.RUnlock()
		return b.take()
	}
	s.dataLock.RUnlock()

	// Unfortunately we did not find the key in the map. Take out a full lock. We
	// have to check if the key exists again, because it's possible another
	// goroutine created it between our shared lock and exclusive lock.
	s.dataLock.Lock()
	if b, ok := s.data[key]; ok {
		s.dataLock.Unlock()
		return b.take()
	}

	// This is the first time we've seen this entry (or it's been garbage
	// collected), so create the bucket and take an initial request.
	b := newBucket(s.tokens, s.interval, s.rate)

	// Add it to the map and take.
	s.data[key] = b
	s.dataLock.Unlock()
	return b.take()
}

// Close stops the memory limiter and cleans up any outstanding sessions. You
// should absolutely always call Close() as it releases the memory consumed by
// the map AND releases the tickers.
func (s *store) Close() error {
	if !atomic.CompareAndSwapUint32(&s.stopped, 0, 1) {
		return nil
	}

	// Close the channel to prevent future purging.
	close(s.stopCh)

	// Delete all the things.
	s.dataLock.Lock()
	for k := range s.data {
		delete(s.data, k)
	}
	s.dataLock.Unlock()
	return nil
}

// purge continually iterates over the map and purges old values on the provided
// sweep interval. Earlier designs used a go-function-per-item expiration, but
// it actually generated *more* lock contention under normal use. The most
// performant option with real-world data was a global garbage collection on a
// fixed interval.
func (s *store) purge() {
	ticker := time.NewTicker(s.sweepInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
		}

		s.dataLock.Lock()
		now := fasttime.Now()
		for k, b := range s.data {
			lastTick := (*bucketState)(atomic.LoadPointer(&b.bucketState)).lastTick
			lastTime := b.startTime + (lastTick * uint64(b.interval))

			if now-lastTime > s.sweepMinTTL {
				delete(s.data, k)
			}
		}
		s.dataLock.Unlock()
	}
}

// bucket is an internal wrapper around a taker.
type bucket struct {
	// startTime is the number of nanoseconds from unix epoch when this bucket was
	// initially created.
	startTime uint64

	// maxTokens is the maximum number of tokens permitted on the bucket at any
	// time. The number of available tokens will never exceed this value.
	maxTokens uint64

	// interval is the time at which ticking should occur.
	interval time.Duration

	// bucketState is the mutable internal state of the event. It includes the
	// current number of available tokens and the last time the clock ticked. It
	// should always be loaded with atomic as it is not concurrent safe.
	bucketState unsafe.Pointer

	// fillRate is the number of tokens to add per nanosecond. It is calculated
	// based on the provided maxTokens and interval.
	fillRate float64
}

// bucketState represents the internal bucket state.
type bucketState struct {
	// availableTokens is the current point-in-time number of tokens remaining.
	// This value changes frequently and must be guarded by an atomic read/write.
	availableTokens uint64

	// lastTick is the last clock tick, used to re-calculate the number of tokens
	// on the bucket.
	lastTick uint64
}

// newBucket creates a new bucket from the given tokens and interval.
func newBucket(tokens uint64, interval time.Duration, rate float64) *bucket {
	b := &bucket{
		startTime: fasttime.Now(),
		maxTokens: tokens,
		interval:  interval,
		fillRate:  rate,

		bucketState: unsafe.Pointer(&bucketState{
			availableTokens: tokens,
		}),
	}
	return b
}

// take attempts to remove a token from the bucket. If there are no tokens
// available and the clock has ticked forward, it recalculates the number of
// tokens and retries. It returns the limit, remaining tokens, time until
// refresh, and whether the take was successful.
func (b *bucket) take() (uint64, uint64, uint64, bool) {
	// Capture the current request time, current tick, and amount of time until
	// the bucket resets.
	now := fasttime.Now()
	currTick := tick(b.startTime, now, b.interval)
	next := b.startTime + ((currTick + 1) * uint64(b.interval))

	for {
		curr := atomic.LoadPointer(&b.bucketState)
		currState := (*bucketState)(curr)
		lastTick := currState.lastTick
		tokens := currState.availableTokens

		if lastTick < currTick {
			tokens = availableTokens(currState.lastTick, currTick, b.maxTokens, b.fillRate)
			lastTick = currTick

			if !atomic.CompareAndSwapPointer(&b.bucketState, curr, unsafe.Pointer(&bucketState{
				availableTokens: tokens,
				lastTick:        lastTick,
			})) {
				// Someone else modified the value
				continue
			}
		}

		if tokens > 0 {
			tokens--
			if !atomic.CompareAndSwapPointer(&b.bucketState, curr, unsafe.Pointer(&bucketState{
				availableTokens: tokens,
				lastTick:        lastTick,
			})) {
				// There were tokens left, but someone took them :(
				continue
			}

			return b.maxTokens, tokens, next, true
		}

		// Returning the TTL until next tick.
		return b.maxTokens, 0, next, false
	}
}

// availableTokens returns the number of available tokens, up to max, between
// the two ticks.
func availableTokens(last, curr, max uint64, fillRate float64) uint64 {
	delta := curr - last

	available := uint64(float64(delta) * fillRate)
	if available > max {
		available = max
	}

	return available
}

// tick is the total number of times the current interval has occurred between
// when the time started (start) and the current time (curr). For example, if
// the start time was 12:30pm and it's currently 1:00pm, and the interval was 5
// minutes, tick would return 6 because 1:00pm is the 6th 5-minute tick. Note
// that tick would return 5 at 12:59pm, because it hasn't reached the 6th tick
// yet.
func tick(start, curr uint64, interval time.Duration) uint64 {
	return (curr - start) / uint64(interval)
}

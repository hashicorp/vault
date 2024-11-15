// Package cache provides generic cache types.
package cache

import (
	"sync"
)

// Comparer is an interface defining a generic compare function.
type Comparer[E any] interface {
	Compare(e E) bool
}

// List is a generic cache list.
type List[K Comparer[K], V any] struct {
	maxEntries int
	valueFn    func(k K) V
	mu         sync.RWMutex
	idx        int
	keys       []K
	values     []V
}

// NewList returns a new cache list.
func NewList[K Comparer[K], V any](maxEntries int, valueFn func(k K) V) *List[K, V] {
	return &List[K, V]{
		maxEntries: maxEntries,
		valueFn:    valueFn,
		keys:       make([]K, 0, maxEntries),
		values:     make([]V, 0, maxEntries),
	}
}

func (l *List[K, V]) find(k K) (v V, ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i, k1 := range l.keys {
		if k1.Compare(k) {
			return l.values[i], true
		}
	}
	return
}

// Get returns the value for the given key.
func (l *List[K, V]) Get(k K) V {
	if v, ok := l.find(k); ok {
		return v
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	v := l.valueFn(k)
	l.idx %= l.maxEntries
	if l.idx >= len(l.keys) {
		l.keys = l.keys[:l.idx+1]
		l.values = l.values[:l.idx+1]
	}
	l.keys[l.idx], l.values[l.idx] = k, v
	l.idx++
	return v
}

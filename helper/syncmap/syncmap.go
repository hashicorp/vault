// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package syncmap

import "sync"

// SyncMap implements a map similar to sync.Map, but with generics and with an equality
// in the values specified by an "ID()" method.
type SyncMap[K comparable, V IDer] struct {
	// lock is used to synchronize access to the map
	lock sync.RWMutex
	// data holds the actual data
	data map[K]V
}

// NewSyncMap returns a new, empty SyncMap.
func NewSyncMap[K comparable, V IDer]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		data: make(map[K]V),
	}
}

// Get returns the value for the given key.
func (m *SyncMap[K, V]) Get(k K) V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data[k]
}

// Pop deletes and returns the value for the given key, if it exists.
func (m *SyncMap[K, V]) Pop(k K) V {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.data[k]
	if ok {
		delete(m.data, k)
	}
	return v
}

// PopIfEqual deletes and returns the value for the given key, if it exists
// and only if the ID is equal to the provided string.
func (m *SyncMap[K, V]) PopIfEqual(k K, id string) V {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.data[k]
	if ok && v.ID() == id {
		delete(m.data, k)
		return v
	}
	var zero V
	return zero
}

// Put adds the given key-value pair to the map and returns the previous value, if any.
func (m *SyncMap[K, V]) Put(k K, v V) V {
	m.lock.Lock()
	defer m.lock.Unlock()
	oldV := m.data[k]
	m.data[k] = v
	return oldV
}

// Clear deletes all entries from the map, and returns the previous map.
func (m *SyncMap[K, V]) Clear() map[K]V {
	m.lock.Lock()
	defer m.lock.Unlock()
	old := m.data
	m.data = make(map[K]V)
	return old
}

// Values returns a copy of all values in the map.
func (m *SyncMap[K, V]) Values() []V {
	m.lock.RLock()
	defer m.lock.RUnlock()

	values := make([]V, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, v)
	}
	return values
}

// IDer is used to extract an ID that SyncMap uses for equality checking.
type IDer interface {
	ID() string
}

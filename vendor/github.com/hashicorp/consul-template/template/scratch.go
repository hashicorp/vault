package template

import (
	"fmt"
	"sort"
	"sync"
)

// Scratch is a wrapper around a map which is used by the template.
type Scratch struct {
	once sync.Once
	sync.RWMutex
	values map[string]interface{}
}

// Key returns a boolean indicating whether the given key exists in the map.
func (s *Scratch) Key(k string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.values[k]
	return ok
}

// Get returns a value previously set by Add or Set
func (s *Scratch) Get(k string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.values[k]
}

// Set stores the value v at the key k. It will overwrite an existing value
// if present.
func (s *Scratch) Set(k string, v interface{}) string {
	s.init()

	s.Lock()
	defer s.Unlock()
	s.values[k] = v
	return ""
}

// SetX behaves the same as Set, except it will not overwrite existing keys if
// already present.
func (s *Scratch) SetX(k string, v interface{}) string {
	s.init()

	s.Lock()
	defer s.Unlock()
	if _, ok := s.values[k]; !ok {
		s.values[k] = v
	}
	return ""
}

// MapSet stores the value v into a key mk in the map named k.
func (s *Scratch) MapSet(k, mk string, v interface{}) (string, error) {
	s.init()

	s.Lock()
	defer s.Unlock()
	return s.mapSet(k, mk, v, true)
}

// MapSetX behaves the same as MapSet, except it will not overwrite the map
// key if it already exists.
func (s *Scratch) MapSetX(k, mk string, v interface{}) (string, error) {
	s.init()

	s.Lock()
	defer s.Unlock()
	return s.mapSet(k, mk, v, false)
}

// mapSet is sets the value in the map, overwriting if o is true. This function
// does not perform locking; callers should lock before invoking.
func (s *Scratch) mapSet(k, mk string, v interface{}, o bool) (string, error) {
	if _, ok := s.values[k]; !ok {
		s.values[k] = make(map[string]interface{})
	}

	typed, ok := s.values[k].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("%q is not a map", k)
	}

	if _, ok := typed[mk]; o || !ok {
		typed[mk] = v
	}
	return "", nil
}

// MapValues returns the list of values in the map sorted by key.
func (s *Scratch) MapValues(k string) ([]interface{}, error) {
	s.init()

	s.Lock()
	defer s.Unlock()
	if s.values == nil {
		return nil, nil
	}

	typed, ok := s.values[k].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	keys := make([]string, 0, len(typed))
	for k := range typed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make([]interface{}, len(keys))
	for i, k := range keys {
		sorted[i] = typed[k]
	}
	return sorted, nil
}

// init initializes the scratch.
func (s *Scratch) init() {
	if s.values == nil {
		s.values = make(map[string]interface{})
	}
}

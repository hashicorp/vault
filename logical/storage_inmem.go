package logical

import (
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
)

// InmemStorage implements Storage and stores all data in memory. It is
// basically a straight copy of physical.Inmem, but it prevents backends from
// having to load all of physical's dependencies (which are legion) just to
// have some testing storage.
type InmemStorage struct {
	sync.RWMutex
	root *radix.Tree
	once sync.Once
}

func (s *InmemStorage) Get(key string) (*StorageEntry, error) {
	s.once.Do(s.init)

	s.RLock()
	defer s.RUnlock()

	if raw, ok := s.root.Get(key); ok {
		se := raw.(*StorageEntry)
		return &StorageEntry{
			Key:   se.Key,
			Value: se.Value,
		}, nil
	}

	return nil, nil
}

func (s *InmemStorage) Put(entry *StorageEntry) error {
	s.once.Do(s.init)

	s.Lock()
	defer s.Unlock()

	s.root.Insert(entry.Key, &StorageEntry{
		Key:   entry.Key,
		Value: entry.Value,
	})
	return nil
}

func (s *InmemStorage) Delete(key string) error {
	s.once.Do(s.init)

	s.Lock()
	defer s.Unlock()

	s.root.Delete(key)
	return nil
}

func (s *InmemStorage) List(prefix string) ([]string, error) {
	s.once.Do(s.init)

	s.RLock()
	defer s.RUnlock()

	var out []string
	seen := make(map[string]interface{})
	walkFn := func(s string, v interface{}) bool {
		trimmed := strings.TrimPrefix(s, prefix)
		sep := strings.Index(trimmed, "/")
		if sep == -1 {
			out = append(out, trimmed)
		} else {
			trimmed = trimmed[:sep+1]
			if _, ok := seen[trimmed]; !ok {
				out = append(out, trimmed)
				seen[trimmed] = struct{}{}
			}
		}
		return false
	}
	s.root.WalkPrefix(prefix, walkFn)

	return out, nil

}

func (s *InmemStorage) init() {
	s.root = radix.New()
}

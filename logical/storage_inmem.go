package logical

import (
	"sync"

	"github.com/hashicorp/vault/physical"
)

// InmemStorage implements Storage and stores all data in memory.
type InmemStorage struct {
	phys *physical.InmemBackend

	once sync.Once
}

func (s *InmemStorage) List(prefix string) ([]string, error) {
	s.once.Do(s.init)

	return s.phys.List(prefix)
}

func (s *InmemStorage) Get(key string) (*StorageEntry, error) {
	s.once.Do(s.init)
	entry, err := s.phys.Get(key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &StorageEntry{
		Key:   entry.Key,
		Value: entry.Value,
	}, nil
}

func (s *InmemStorage) Put(entry *StorageEntry) error {
	s.once.Do(s.init)
	physEntry := &physical.Entry{
		Key:   entry.Key,
		Value: entry.Value,
	}
	return s.phys.Put(physEntry)
}

func (s *InmemStorage) Delete(k string) error {
	s.once.Do(s.init)
	return s.phys.Delete(k)
}

func (s *InmemStorage) init() {
	s.phys = physical.NewInmem(nil)
}

package logical

import (
	"strings"
	"sync"
)

// InmemStorage implements Storage and stores all data in memory.
type InmemStorage struct {
	Data map[string]*StorageEntry

	once sync.Once
}

func (s *InmemStorage) List(prefix string) ([]string, error) {
	s.once.Do(s.init)

	var result []string
	for k, _ := range s.Data {
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}

	return result, nil
}

func (s *InmemStorage) Get(key string) (*StorageEntry, error) {
	s.once.Do(s.init)
	return s.Data[key], nil
}

func (s *InmemStorage) Put(entry *StorageEntry) error {
	s.once.Do(s.init)
	s.Data[entry.Key] = entry
	return nil
}

func (s *InmemStorage) Delete(k string) error {
	s.once.Do(s.init)
	delete(s.Data, k)
	return nil
}

func (s *InmemStorage) init() {
	s.Data = make(map[string]*StorageEntry)
}

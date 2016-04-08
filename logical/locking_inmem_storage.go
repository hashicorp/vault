package logical

import (
	"strings"
	"sync"
)

// LockingInmemStorage implements Storage and stores all data in memory.
type LockingInmemStorage struct {
	sync.RWMutex

	Data map[string]*StorageEntry

	once sync.Once
}

func (s *LockingInmemStorage) List(prefix string) ([]string, error) {
	s.once.Do(s.init)

	s.RLock()
	defer s.RUnlock()

	var result []string
	for k, _ := range s.Data {
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}

	return result, nil
}

func (s *LockingInmemStorage) Get(key string) (*StorageEntry, error) {
	s.once.Do(s.init)
	s.RLock()
	defer s.RUnlock()
	return s.Data[key], nil
}

func (s *LockingInmemStorage) Put(entry *StorageEntry) error {
	s.once.Do(s.init)
	s.Lock()
	defer s.Unlock()
	s.Data[entry.Key] = entry
	return nil
}

func (s *LockingInmemStorage) Delete(k string) error {
	s.once.Do(s.init)
	s.Lock()
	defer s.Unlock()
	delete(s.Data, k)
	return nil
}

func (s *LockingInmemStorage) init() {
	s.Data = make(map[string]*StorageEntry)
}

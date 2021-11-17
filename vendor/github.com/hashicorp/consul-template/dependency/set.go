package dependency

import (
	"strings"
	"sync"
)

// Set is a dependency-specific set implementation. Relative ordering is
// preserved.
type Set struct {
	once sync.Once
	sync.RWMutex
	list []string
	set  map[string]Dependency
}

// Add adds a new element to the set if it does not already exist.
func (s *Set) Add(d Dependency) bool {
	s.init()
	s.Lock()
	defer s.Unlock()
	if _, ok := s.set[d.String()]; !ok {
		s.list = append(s.list, d.String())
		s.set[d.String()] = d
		return true
	}
	return false
}

// Get retrieves a single element from the set by name.
func (s *Set) Get(v string) Dependency {
	s.RLock()
	defer s.RUnlock()
	return s.set[v]
}

// List returns the insertion-ordered list of dependencies.
func (s *Set) List() []Dependency {
	s.RLock()
	defer s.RUnlock()
	r := make([]Dependency, len(s.list))
	for i, k := range s.list {
		r[i] = s.set[k]
	}
	return r
}

// Len is the size of the set.
func (s *Set) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.list)
}

// String is a string representation of the set.
func (s *Set) String() string {
	s.RLock()
	defer s.RUnlock()
	return strings.Join(s.list, ", ")
}

func (s *Set) init() {
	s.once.Do(func() {
		if s.list == nil {
			s.list = make([]string, 0, 8)
		}

		if s.set == nil {
			s.set = make(map[string]Dependency)
		}
	})
}

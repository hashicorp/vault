package locksutil

import (
	"fmt"
	"sync"
)

// Takes in a map, indexed by string and creates new 'sync.RWMutex' items.
// This utility creates 'count' number of mutexes (with a cap of 256) and
// places them in the map. The indices will be 2 character hexadecimal
// string values from 0 to count.
func CreateLocks(p map[string]*sync.RWMutex, count int64) error {
	// Since the indices of the map entries are based on 2 character
	// hex values, this utility can only create upto 256 locks.
	if count <= 0 || count > 256 {
		return fmt.Errorf("invalid count: %d", count)
	}

	for i := int64(0); i < count; i++ {
		p[fmt.Sprintf("%02x", i)] = &sync.RWMutex{}
	}

	return nil
}

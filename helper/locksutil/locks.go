package locksutil

import (
	"fmt"
	"sync"
)

func CreateLocks(p map[string]*sync.RWMutex, count int64) {
	for i := int64(0); i < count; i++ {
		p[fmt.Sprintf("%02x", i)] = &sync.RWMutex{}
	}
}

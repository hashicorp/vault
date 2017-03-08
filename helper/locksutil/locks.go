package locksutil

import (
	"crypto/md5"
	"sync"
)

const (
	LockCount = 256
)

type LockEntry struct {
	sync.RWMutex
}

func CreateLocks() []*LockEntry {
	ret := make([]*LockEntry, LockCount)
	for i := range ret {
		ret[i] = new(LockEntry)
	}
	return ret
}

func LockIndexForKey(key string) uint8 {
	hf := md5.New()
	hf.Write([]byte(key))
	return uint8(hf.Sum(nil)[0])
}

func LockForKey(locks []*LockEntry, key string) *LockEntry {
	return locks[LockIndexForKey(key)]
}

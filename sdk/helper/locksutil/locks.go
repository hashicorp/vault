// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package locksutil

import (
	"sync"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/sasha-s/go-deadlock"
)

const (
	LockCount = 256
)

// DeadlockRWMutex is the RW version of DeadlockMutex.
type DeadlockRWMutex struct {
	deadlock.RWMutex
}

type LockEntry struct {
	sync.RWMutex
}

// CreateLocks returns an array so that the locks can be iterated over in
// order.
//
// This is only threadsafe if a process is using a single lock, or iterating
// over the entire lock slice in order. Using a consistent order avoids
// deadlocks because you can never have the following:
//
// Lock A, Lock B
// Lock B, Lock A
//
// Where process 1 is now deadlocked trying to lock B, and process 2 deadlocked trying to lock A
func CreateLocks() []*LockEntry {
	ret := make([]*LockEntry, LockCount)
	for i := range ret {
		ret[i] = new(LockEntry)
	}
	return ret
}

func CreateLocksWithDeadlockDetection() []*DeadlockRWMutex {
	ret := make([]*DeadlockRWMutex, LockCount)
	for i := range ret {
		ret[i] = new(DeadlockRWMutex)
	}
	return ret
}

func LockIndexForKey(key string) uint8 {
	return uint8(cryptoutil.Blake2b256Hash(key)[0])
}

func LockForKey(locks []*LockEntry, key string) *LockEntry {
	return locks[LockIndexForKey(key)]
}

func LocksForKeys(locks []*LockEntry, keys []string) []*LockEntry {
	lockIndexes := make(map[uint8]struct{}, len(keys))
	for _, k := range keys {
		lockIndexes[LockIndexForKey(k)] = struct{}{}
	}

	locksToReturn := make([]*LockEntry, 0, len(keys))
	for i, l := range locks {
		if _, ok := lockIndexes[uint8(i)]; ok {
			locksToReturn = append(locksToReturn, l)
		}
	}

	return locksToReturn
}

func LockForKeyWithDeadLockDetection(locks []*DeadlockRWMutex, key string) *DeadlockRWMutex {
	return locks[LockIndexForKey(key)]
}

func LocksForKeysWithDeadLockDetection(locks []*DeadlockRWMutex, keys []string) []*DeadlockRWMutex {
	lockIndexes := make(map[uint8]struct{}, len(keys))
	for _, k := range keys {
		lockIndexes[LockIndexForKey(k)] = struct{}{}
	}

	locksToReturn := make([]*DeadlockRWMutex, 0, len(keys))
	for i, l := range locks {
		if _, ok := lockIndexes[uint8(i)]; ok {
			locksToReturn = append(locksToReturn, l)
		}
	}

	return locksToReturn
}

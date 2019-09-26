package locksutil

import (
	"strconv"
	"sync"
	"testing"
)

func Test_CreateLocks(t *testing.T) {
	locks := CreateLocks()
	if len(locks) != 256 {
		t.Fatalf("bad: len(locks): expected:256 actual:%d", len(locks))
	}
}

func Test_LockRetrieval(t *testing.T) {
	locks := CreateLocks()

	l1 := LockForKey(locks, "key1")
	l2 := LockForKey(locks, "key2")

	if l1 != LockForKey(locks, "key1") {
		t.Fatal("received different lock for the same key")
	}

	if l2 != LockForKey(locks, "key2") {
		t.Fatal("received different lock for the same key")
	}

	// This is valid for the specific case of "key1" and "key2", which are known to return
	// different locks, and is testing for a regression that always returns the same lock.
	if l1 == l2 {
		t.Fatal("expected different locks")
	}
}

func Test_LockUsage(t *testing.T) {
	const (
		numGoroutines = 400
		numAdditions  = 200
		numCounters   = 10
	)

	// many goroutines will compete to make additions to a limited number of counters
	sharedCounters := make([]int, numCounters)

	locks := CreateLocks()

	var wg sync.WaitGroup
	incrementer := func(counterID int) {
		for i := 0; i < numAdditions; i++ {
			l := LockForKey(locks, strconv.Itoa(counterID))
			l.Lock()
			sharedCounters[counterID]++
			l.Unlock()
		}
		wg.Done()
	}

	// start goroutines, assigning each one to a shared counter
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go incrementer(i % numCounters)
	}
	wg.Wait()

	// confirm that all additions were made
	for _, v := range sharedCounters {
		if exp := numAdditions * numGoroutines / 10; v != exp {
			t.Fatalf("expected counter value %d, got %d", exp, v)
		}
	}
}

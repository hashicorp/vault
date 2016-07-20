package locksutil

import (
	"sync"
	"testing"
)

func Test_CreateLocks(t *testing.T) {
	locks := map[string]*sync.RWMutex{}

	// Invalid argument
	if err := CreateLocks(locks, -1); err == nil {
		t.Fatal("expected an error")
	}

	// Invalid argument
	if err := CreateLocks(locks, 0); err == nil {
		t.Fatal("expected an error")
	}

	// Invalid argument
	if err := CreateLocks(locks, 300); err == nil {
		t.Fatal("expected an error")
	}

	// Maximum number of locks
	if err := CreateLocks(locks, 256); err != nil {
		t.Fatal("err: %v", err)
	}
	if len(locks) != 256 {
		t.Fatal("bad: len(locks): expected:256 actual:%d", len(locks))
	}

	// Clear out the locks for testing the next case
	for k, _ := range locks {
		delete(locks, k)
	}

	// General case
	if err := CreateLocks(locks, 10); err != nil {
		t.Fatal("err: %v", err)
	}
	if len(locks) != 10 {
		t.Fatal("bad: len(locks): expected:10 actual:%d", len(locks))
	}

}

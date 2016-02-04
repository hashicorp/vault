package logical

import (
	"testing"
)

// Note: This uses the normal TestStorage, but the best way to exercise this is
// to run transit's unit tests, which spawn 1000 goroutines to hammer the
// backend for 10 seconds with this as the storage.

func TestLockingInmemStorage(t *testing.T) {
	TestStorage(t, new(LockingInmemStorage))
}

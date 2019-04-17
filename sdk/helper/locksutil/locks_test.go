package locksutil

import "testing"

func Test_CreateLocks(t *testing.T) {
	locks := CreateLocks()
	if len(locks) != 256 {
		t.Fatalf("bad: len(locks): expected:256 actual:%d", len(locks))
	}
}

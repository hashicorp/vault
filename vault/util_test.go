package vault

import "testing"

func TestMemZero(t *testing.T) {
	b := []byte{1, 2, 3, 4}
	memzero(b)
	if b[0] != 0 || b[1] != 0 || b[2] != 0 || b[3] != 0 {
		t.Fatalf("bad: %v", b)
	}
}

func TestRandBytes(t *testing.T) {
	b := randbytes(12)
	if len(b) != 12 {
		t.Fatalf("bad: %v", b)
	}
}

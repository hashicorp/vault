package vault

import (
	"regexp"
	"testing"
)

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

func TestGenerateUUID(t *testing.T) {
	prev := generateUUID()
	for i := 0; i < 100; i++ {
		id := generateUUID()
		if prev == id {
			t.Fatalf("Should get a new ID!")
		}

		matched, err := regexp.MatchString(
			"[\\da-f]{8}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{12}", id)
		if !matched || err != nil {
			t.Fatalf("expected match %s %v %s", id, matched, err)
		}
	}
}

func TestStrListContains(t *testing.T) {
	haystack := []string{
		"dev",
		"ops",
		"prod",
		"root",
	}
	if strListContains(haystack, "tubez") {
		t.Fatalf("Bad")
	}
	if !strListContains(haystack, "root") {
		t.Fatalf("Bad")
	}
}

func TestStrListSubset(t *testing.T) {
	parent := []string{
		"dev",
		"ops",
		"prod",
		"root",
	}
	child := []string{
		"prod",
		"ops",
	}
	if !strListSubset(parent, child) {
		t.Fatalf("Bad")
	}
	if !strListSubset(parent, parent) {
		t.Fatalf("Bad")
	}
	if !strListSubset(child, child) {
		t.Fatalf("Bad")
	}
	if !strListSubset(child, nil) {
		t.Fatalf("Bad")
	}
	if strListSubset(child, parent) {
		t.Fatalf("Bad")
	}
	if strListSubset(nil, child) {
		t.Fatalf("Bad")
	}
}

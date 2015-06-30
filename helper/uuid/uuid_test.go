package uuid

import (
	"regexp"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	prev := GenerateUUID()
	for i := 0; i < 100; i++ {
		id := GenerateUUID()
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

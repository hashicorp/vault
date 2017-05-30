package physical

import "testing"

func TestAppendIfMissing(t *testing.T) {
	keys := []string{}

	keys = appendIfMissing(keys, "foo")

	if len(keys) != 1 {
		t.Fatalf("expected slice to be length of 1: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to contain key 'foo': %v", keys)
	}

	keys = appendIfMissing(keys, "bar")

	if len(keys) != 2 {
		t.Fatalf("expected slice to be length of 2: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to contain key 'foo': %v", keys)
	}
	if keys[1] != "bar" {
		t.Fatalf("expected slice to contain key 'bar': %v", keys)
	}

	keys = appendIfMissing(keys, "foo")

	if len(keys) != 2 {
		t.Fatalf("expected slice to still be length of 2: %v", keys)
	}
	if keys[0] != "foo" {
		t.Fatalf("expected slice to still contain key 'foo': %v", keys)
	}
	if keys[1] != "bar" {
		t.Fatalf("expected slice to still contain key 'bar': %v", keys)
	}
}

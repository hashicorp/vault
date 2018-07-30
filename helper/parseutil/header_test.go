package parseutil

import (
	"testing"
)

func TestHeaderGetAll(t *testing.T) {
	h := NewHeader()
	h.Add("hello", "world")
	h.Add("hello", "monde")
	results := h.GetAll("hello")
	if len(results) != 2 {
		t.Fatal("expected 2 results")
	}
	if results[0] != "world" {
		t.Fatal("expected the first result to be world")
	}
	if results[1] != "monde" {
		t.Fatal("expected the second result to be monde")
	}
}

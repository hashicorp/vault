package credsutil

import (
	"bytes"
	"testing"
)

// RandomAlphaNumericOfLen returns a random byte slice of characters [A-Za-z0-9]
// of the provided length.
func TestRandomAlphaNumericOfLen(t *testing.T) {
	s, err := RandomAlphaNumericOfLen(1)
	if err != nil {
		t.Fatal("Unexpected error: %s", err)
	}
	if len(s) != 1 {
		t.Fatal("Unexpected length of string, expected 1, got string: %s", s)
	}

	s, err = RandomAlphaNumericOfLen(10)
	if err != nil {
		t.Fatal("Unexpected error: %s", err)
	}
	if len(s) != 10 {
		t.Fatal("Unexpected length of string, expected 10, got string: %s", s)
	}

	if bytes.Equal(s, make([]byte, 10)) {
		t.Fatal("returned byte slice is empty")
	}
}

package vault

import (
	"crypto/rand"
	"fmt"
)

// memzero is used to zero out a byte buffer. This specific format is optimized
// by the compiler to use memclr to improve performance. See this code review:
// https://codereview.appspot.com/137880043
func memzero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// randbytes is used to create a buffer of size n filled with random bytes
func randbytes(n int) []byte {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("failed to generate %d random bytes: %v", n, err))
	}
	return buf
}

// generateUUID is used to generate a random UUID
func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

// strListContains looks for a string in a list of strings.
func strListContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// strListSubset checks if a given list is a subset
// of another set
func strListSubset(super, sub []string) bool {
	for _, item := range sub {
		if !strListContains(super, item) {
			return false
		}
	}
	return true
}

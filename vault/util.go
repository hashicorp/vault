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

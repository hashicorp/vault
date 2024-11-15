// Package alphanum implements functions for randomized alphanum content.
package alphanum

import (
	"crypto/rand"

	"github.com/SAP/go-hdb/driver/internal/unsafe"
)

const csAlphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // alphanumeric character set.
var numAlphanum = byte(len(csAlphanum))                                             // len character sets <= max(byte)

// Read fills p with random alphanumeric characters and returns the number of read bytes and a potential error.
func Read(p []byte) (n int, err error) {
	if n, err = rand.Read(p); err != nil {
		return n, err
	}
	for i, b := range p {
		p[i] = csAlphanum[b%numAlphanum]
	}
	return n, nil
}

// ReadString returns a random string of alphanumeric characters and panics if crypto random reader returns an error.
func ReadString(n int) string {
	b := make([]byte, n)
	if _, err := Read(b); err != nil {
		panic(err) // rand should never fail
	}
	return unsafe.ByteSlice2String(b)
}

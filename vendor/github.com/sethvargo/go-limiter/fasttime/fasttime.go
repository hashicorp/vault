// Package fasttime provides a fast clock implementation, roughly 2.5x faster
// than time.Now.
package fasttime

import (
	_ "unsafe"
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// Now returns the current unix time.
func Now() uint64 {
	return uint64(nanotime())
}

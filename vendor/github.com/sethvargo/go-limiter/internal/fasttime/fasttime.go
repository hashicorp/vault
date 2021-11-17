//go:build !windows
// +build !windows

// Package fasttime gets wallclock time, but super fast.
package fasttime

import (
	_ "unsafe"
)

//go:noescape
//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)

// Now returns a monotonic clock value. The actual value will differ across
// systems, but that's okay because we generally only care about the deltas.
func Now() uint64 {
	sec, nsec, _ := now()
	return uint64(sec)*1e9 + uint64(nsec)
}

// +build !windows

// Package fasttime gets wallclock time, but super fast.
package fasttime

import (
	_ "unsafe"
)

//go:noescape
//go:linkname walltime runtime.walltime
func walltime() (int64, int32)

// Now returns a monotonic clock value. The actual value will differ across
// systems, but that's okay because we generally only care about the deltas.
func Now() uint64 {
	x, y := walltime()
	return uint64(x)*1e9 + uint64(y)
}

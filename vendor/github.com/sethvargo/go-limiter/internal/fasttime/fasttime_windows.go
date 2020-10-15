// +build windows

package fasttime

import "time"

// Now returns a monotonic clock value. On Windows, no such clock exists, so we
// fallback to time.Now().
func Now() uint64 {
	return uint64(time.Now().UnixNano())
}

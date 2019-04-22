package consul

import (
	"math/rand"
	"time"
)

// DurationMinusBuffer returns a duration, minus a buffer and jitter
// subtracted from the duration.  This function is used primarily for
// servicing Consul TTL Checks in advance of the TTL.
func DurationMinusBuffer(intv time.Duration, buffer time.Duration, jitter int64) time.Duration {
	d := intv - buffer
	if jitter == 0 {
		d -= RandomStagger(d)
	} else {
		d -= RandomStagger(time.Duration(int64(d) / jitter))
	}
	return d
}

// DurationMinusBufferDomain returns the domain of valid durations from a
// call to DurationMinusBuffer.  This function is used to check user
// specified input values to DurationMinusBuffer.
func DurationMinusBufferDomain(intv time.Duration, buffer time.Duration, jitter int64) (min time.Duration, max time.Duration) {
	max = intv - buffer
	if jitter == 0 {
		min = max
	} else {
		min = max - time.Duration(int64(max)/jitter)
	}
	return min, max
}

// RandomStagger returns an interval between 0 and the duration
func RandomStagger(intv time.Duration) time.Duration {
	if intv == 0 {
		return 0
	}
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}

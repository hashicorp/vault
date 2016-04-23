package lib

import (
	"math/rand"
	"time"
)

// DurationMinusBuffer returns a duration, minus a buffer and jitter
// subtracted from the duration.  This function is used primarily for
// servicing Consul TTL Checks in advance of the TTL.
func DurationMinusBuffer(intv time.Duration, buffer time.Duration, jitter int64) time.Duration {
	d := intv - buffer
	d -= RandomStagger(time.Duration(int64(d) / jitter))
	return d
}

// Returns a random stagger interval between 0 and the duration
func RandomStagger(intv time.Duration) time.Duration {
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}

// RateScaledInterval is used to choose an interval to perform an action in
// order to target an aggregate number of actions per second across the whole
// cluster.
func RateScaledInterval(rate float64, min time.Duration, n int) time.Duration {
	const minRate = 1 / 86400 // 1/(1 * time.Day)
	if rate <= minRate {
		return min
	}
	interval := time.Duration(float64(time.Second) * float64(n) / rate)
	if interval < min {
		return min
	}

	return interval
}

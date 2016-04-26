package lib

import (
	"math/rand"
	"sync"
	"time"
)

var (
	once sync.Once
)

// SeedMathRand provides weak, but guaranteed seeding, which is better than
// running with Go's default seed of 1.  A call to SeedMathRand() is expected
// to be called via init(), but never a second time.
func SeedMathRand() {
	once.Do(func() { rand.Seed(time.Now().UTC().UnixNano()) })
}

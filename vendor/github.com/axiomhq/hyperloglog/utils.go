package hyperloglog

import (
	"math"
	"math/bits"

	metro "github.com/dgryski/go-metro"
)

var hash = hashFunc

func alpha(m float64) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	}
	return 0.7213 / (1 + 1.079/m)
}

func getPosVal(x uint64, p uint8) (uint64, uint8) {
	i := bextr(x, 64-p, p) // {x63,...,x64-p}
	w := x<<p | 1<<(p-1)   // {x63-p,...,x0}
	rho := uint8(bits.LeadingZeros64(w)) + 1
	return i, rho
}

func linearCount(m uint32, v uint32) float64 {
	fm := float64(m)
	return fm * math.Log(fm/float64(v))
}

func bextr(v uint64, start, length uint8) uint64 {
	return (v >> start) & ((1 << length) - 1)
}

func bextr32(v uint32, start, length uint8) uint32 {
	return (v >> start) & ((1 << length) - 1)
}

func hashFunc(e []byte) uint64 {
	return metro.Hash64(e, 1337)
}

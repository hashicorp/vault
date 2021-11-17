// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package decimal128 // import "github.com/apache/arrow/go/arrow/decimal128"

import (
	"math/big"
)

var (
	MaxDecimal128 = New(542101086242752217, 687399551400673280-1)
)

// Num represents a signed 128-bit integer in two's complement.
// Calculations wrap around and overflow is ignored.
//
// For a discussion of the algorithms, look at Knuth's volume 2,
// Semi-numerical Algorithms section 4.3.1.
//
// Adapted from the Apache ORC C++ implementation
type Num struct {
	lo uint64 // low bits
	hi int64  // high bits
}

// New returns a new signed 128-bit integer value.
func New(hi int64, lo uint64) Num {
	return Num{lo: lo, hi: hi}
}

// FromU64 returns a new signed 128-bit integer value from the provided uint64 one.
func FromU64(v uint64) Num {
	return New(0, v)
}

// FromI64 returns a new signed 128-bit integer value from the provided int64 one.
func FromI64(v int64) Num {
	switch {
	case v > 0:
		return New(0, uint64(v))
	case v < 0:
		return New(-1, uint64(v))
	default:
		return Num{}
	}
}

// FromBigInt will convert a big.Int to a Num, if the value in v has a
// BitLen > 128, this will panic.
func FromBigInt(v *big.Int) (n Num) {
	bitlen := v.BitLen()
	if bitlen > 128 {
		panic("arrow/decimal128: cannot represent value larger than 128bits")
	} else if bitlen == 0 {
		// if bitlen is 0, then the value is 0 so return the default zeroed
		// out n
		return
	}

	// if the value is negative, then get the high and low bytes from
	// v, and then negate it. this is because Num uses a two's compliment
	// representation of values and big.Int stores the value as a bool for
	// the sign and the absolute value of the integer. This means that the
	// raw bytes are *always* the absolute value.
	b := v.Bits()
	n.lo = uint64(b[0])
	if len(b) > 1 {
		n.hi = int64(b[1])
	}
	if v.Sign() < 0 {
		return n.negated()
	}
	return
}

func (n Num) negated() Num {
	n.lo = ^n.lo + 1
	n.hi = ^n.hi
	if n.lo == 0 {
		n.hi += 1
	}
	return n
}

// LowBits returns the low bits of the two's complement representation of the number.
func (n Num) LowBits() uint64 { return n.lo }

// HighBits returns the high bits of the two's complement representation of the number.
func (n Num) HighBits() int64 { return n.hi }

// Sign returns:
//
// -1 if x <  0
//  0 if x == 0
// +1 if x >  0
func (n Num) Sign() int {
	if n == (Num{}) {
		return 0
	}
	return int(1 | (n.hi >> 63))
}

func toBigIntPositive(n Num) *big.Int {
	return (&big.Int{}).SetBits([]big.Word{big.Word(n.lo), big.Word(n.hi)})
}

// while the code would be simpler to just do lsh/rsh and add
// it turns out from benchmarking that calling SetBits passing
// in the words and negating ends up being >2x faster
func (n Num) BigInt() *big.Int {
	if n.Sign() < 0 {
		b := toBigIntPositive(n.negated())
		return b.Neg(b)
	}
	return toBigIntPositive(n)
}

// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rand

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
)

const (
	poolSize = 512
)

// random number generator pool
var (
	pool    = make([]*Xor128Rand, poolSize)
	pos     uint64
	rndInit = &Xor128Rand{src: [2]uint64{uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())}}
)

func init() {
	for i := range pool {
		pool[i] = NewXorRand()
	}
}

// Int64 returns a random int64 number. It can be negative.
// This function uses a pool and is lockless.
func Int64() int64 {
	apos := int(atomic.AddUint64(&pos, 1) % poolSize)
	return pool[apos].Int64()
}

// Uint64 returns a random uint64 number.
// This function uses a pool and is lockless.
func Uint64() uint64 {
	apos := int(atomic.AddUint64(&pos, 1) % poolSize)
	return pool[apos].Uint64()
}

// Xor128Rand is a random number generator
type Xor128Rand struct {
	src [2]uint64
	l   sync.Mutex
}

// NewXorRand creates a XOR Shift random number generator.
func NewXorRand() *Xor128Rand {
	t := time.Now().UnixNano() + rndInit.Int64()
	return &Xor128Rand{src: [2]uint64{uint64(t), uint64(t)}}
}

// Int64 returns a random int64 number. It can be negative.
func (r *Xor128Rand) Int64() int64 {
	return int64(r.Uint64())
}

// Uint64 returns a random uint64 number.
func (r *Xor128Rand) Uint64() uint64 {
	r.l.Lock()
	s1 := r.src[0]
	s0 := r.src[1]
	r.src[0] = s0
	s1 ^= s1 << 23
	r.src[1] = (s1 ^ s0 ^ (s1 >> 17) ^ (s0 >> 26))
	res := r.src[1] + s0
	r.l.Unlock()
	return res
}

// Read will fill the argument slice with random bytes.
// Implements the Reader interface.
func (r *Xor128Rand) Read(p []byte) (n int, err error) {
	l := len(p) / 8
	for i := 0; i < l; i += 8 {
		binary.PutUvarint(p[i:], r.Uint64())
	}
	return len(p), nil
}

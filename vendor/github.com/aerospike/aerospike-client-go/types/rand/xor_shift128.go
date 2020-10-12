// Copyright 2013-2020 Aerospike, Inc.
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
var pool = make([]*Xor128Rand, poolSize)
var pos uint64

func init() {
	for i := range pool {
		pool[i] = NewXorRand()
		// to guarantee a different number on less accurate system clocks
		time.Sleep(time.Microsecond + 31*time.Nanosecond)
	}
}

func Int64() int64 {
	apos := int(atomic.AddUint64(&pos, 1) % poolSize)
	return pool[apos].Int64()
}

type Xor128Rand struct {
	src [2]uint64
	l   sync.Mutex
}

func NewXorRand() *Xor128Rand {
	t := time.Now().UnixNano()
	return &Xor128Rand{src: [2]uint64{uint64(t), uint64(t)}}
}

func (r *Xor128Rand) Int64() int64 {
	return int64(r.Uint64())
}

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

func (r *Xor128Rand) Read(p []byte) (n int, err error) {
	l := len(p) / 8
	for i := 0; i < l; i += 8 {
		binary.PutUvarint(p[i:], r.Uint64())
	}
	return len(p), nil
}

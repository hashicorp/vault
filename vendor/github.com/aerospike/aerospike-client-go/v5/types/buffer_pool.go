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

package types

import "sync"

// BufferPool implements a specialized buffer pool.
// Pool size will be limited, and each buffer size will be
// constrained to the init and max buffer sizes.
type BufferPool struct {
	pool     [][]byte
	poolSize int

	pos int64

	maxBufSize  int
	initBufSize int

	mutex sync.Mutex
}

// NewBufferPool creates a new buffer pool.
// New buffers will be created with size and capacity of initBufferSize.
// If  cap(buffer) is larger than maxBufferSize when it is put back in the buffer,
// it will be thrown away. This will prevent unwanted memory bloat and
// set a deterministic maximum-size for the pool which will not be exceeded.
func NewBufferPool(poolSize, initBufferSize, maxBufferSize int) *BufferPool {
	return &BufferPool{
		pool:        make([][]byte, poolSize),
		pos:         -1,
		poolSize:    poolSize,
		maxBufSize:  maxBufferSize,
		initBufSize: initBufferSize,
	}
}

// Get returns a buffer from the pool. If pool is empty, a new buffer of
// size initBufSize will be created and returned.
func (bp *BufferPool) Get() (res []byte) {
	bp.mutex.Lock()
	if bp.pos >= 0 {
		res = bp.pool[bp.pos]
		bp.pos--
	} else {
		res = make([]byte, bp.initBufSize)
	}

	bp.mutex.Unlock()
	return res
}

// Put will put the buffer back in the pool, unless cap(buf) is bigger than
// initBufSize, in which case it will be thrown away
func (bp *BufferPool) Put(buf []byte) {
	if len(buf) <= bp.maxBufSize {
		bp.mutex.Lock()
		if bp.pos < int64(bp.poolSize-1) {
			bp.pos++
			bp.pool[bp.pos] = buf
		}
		bp.mutex.Unlock()
	}
}

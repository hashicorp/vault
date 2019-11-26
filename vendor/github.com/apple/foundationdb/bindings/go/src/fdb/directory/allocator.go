/*
 * allocator.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go Directory Layer

package directory

import (
	"bytes"
	"encoding/binary"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"math/rand"
	"sync"
)

var oneBytes = []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var allocatorMutex = sync.Mutex{}

type highContentionAllocator struct {
	counters, recent subspace.Subspace
}

func newHCA(s subspace.Subspace) highContentionAllocator {
	var hca highContentionAllocator

	hca.counters = s.Sub(0)
	hca.recent = s.Sub(1)

	return hca
}

func windowSize(start int64) int64 {
	// Larger window sizes are better for high contention, smaller sizes for
	// keeping the keys small.  But if there are many allocations, the keys
	// can't be too small.  So start small and scale up.  We don't want this to
	// ever get *too* big because we have to store about window_size/2 recent
	// items.
	if start < 255 {
		return 64
	}
	if start < 65535 {
		return 1024
	}
	return 8192
}

func (hca highContentionAllocator) allocate(tr fdb.Transaction, s subspace.Subspace) (subspace.Subspace, error) {
	for {
		rr := tr.Snapshot().GetRange(hca.counters, fdb.RangeOptions{Limit: 1, Reverse: true})
		kvs, e := rr.GetSliceWithError()
		if e != nil {
			return nil, e
		}

		var start int64
		var window int64

		if len(kvs) == 1 {
			t, e := hca.counters.Unpack(kvs[0].Key)
			if e != nil {
				return nil, e
			}
			start = t[0].(int64)
		}

		windowAdvanced := false
		for {
			allocatorMutex.Lock()

			if windowAdvanced {
				tr.ClearRange(fdb.KeyRange{hca.counters, hca.counters.Sub(start)})
				tr.Options().SetNextWriteNoWriteConflictRange()
				tr.ClearRange(fdb.KeyRange{hca.recent, hca.recent.Sub(start)})
			}

			// Increment the allocation count for the current window
			tr.Add(hca.counters.Sub(start), oneBytes)
			countFuture := tr.Snapshot().Get(hca.counters.Sub(start))

			allocatorMutex.Unlock()

			countStr, e := countFuture.Get()
			if e != nil {
				return nil, e
			}

			var count int64
			if countStr == nil {
				count = 0
			} else {
				e = binary.Read(bytes.NewBuffer(countStr), binary.LittleEndian, &count)
				if e != nil {
					return nil, e
				}
			}

			window = windowSize(start)
			if count*2 < window {
				break
			}

			start += window
			windowAdvanced = true
		}

		for {
			// As of the snapshot being read from, the window is less than half
			// full, so this should be expected to take 2 tries.  Under high
			// contention (and when the window advances), there is an additional
			// subsequent risk of conflict for this transaction.
			candidate := rand.Int63n(window) + start
			key := hca.recent.Sub(candidate)

			allocatorMutex.Lock()

			latestCounter := tr.Snapshot().GetRange(hca.counters, fdb.RangeOptions{Limit: 1, Reverse: true})
			candidateValue := tr.Get(key)
			tr.Options().SetNextWriteNoWriteConflictRange()
			tr.Set(key, []byte(""))

			allocatorMutex.Unlock()

			kvs, e = latestCounter.GetSliceWithError()
			if e != nil {
				return nil, e
			}
			if len(kvs) > 0 {
				t, e := hca.counters.Unpack(kvs[0].Key)
				if e != nil {
					return nil, e
				}
				currentStart := t[0].(int64)
				if currentStart > start {
					break
				}
			}

			v, e := candidateValue.Get()
			if e != nil {
				return nil, e
			}
			if v == nil {
				tr.AddWriteConflictKey(key)
				return s.Sub(candidate), nil
			}
		}
	}
}

/*
 * snapshot.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2024 Apple Inc. and the FoundationDB project authors
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

// FoundationDB Go API

package fdb

// Snapshot is a handle to a FoundationDB transaction snapshot, suitable for
// performing snapshot reads. Snapshot reads offer a more relaxed isolation
// level than FoundationDB's default serializable isolation, reducing
// transaction conflicts but making it harder to reason about concurrency.
//
// For more information on snapshot reads, see
// https://apple.github.io/foundationdb/developer-guide.html#snapshot-reads.
type Snapshot struct {
	*transaction
}

// ReadTransact executes the caller-provided function, passing it the Snapshot
// receiver object (as a ReadTransaction).
//
// A panic of type Error during execution of the function will be recovered and
// returned to the caller as an error, but ReadTransact will not retry the
// function.
//
// By satisfying the ReadTransactor interface, Snapshot may be passed to a
// read-only transactional function from another (possibly read-only)
// transactional function, allowing composition.
//
// See the ReadTransactor interface for an example of using ReadTransact with
// Transaction, Snapshot and Database objects.
func (s Snapshot) ReadTransact(f func(ReadTransaction) (interface{}, error)) (r interface{}, e error) {
	defer panicToError(&e)

	r, e = f(s)
	return
}

// Cancel cancels the underlying transaction of the snapshot.
// See Transaction.Cancel() for more information.
func (s Snapshot) Cancel() {
	s.transaction.cancel()
}

// Snapshot returns the receiver and allows Snapshot to satisfy the
// ReadTransaction interface.
func (s Snapshot) Snapshot() Snapshot {
	return s
}

// Get is equivalent to (Transaction).Get, performed as a snapshot read.
func (s Snapshot) Get(key KeyConvertible) FutureByteSlice {
	return s.get(key.FDBKey(), 1)
}

// GetKey is equivalent to (Transaction).GetKey, performed as a snapshot read.
func (s Snapshot) GetKey(sel Selectable) FutureKey {
	return s.getKey(sel.FDBKeySelector(), 1)
}

// GetRange is equivalent to (Transaction).GetRange, performed as a snapshot
// read.
func (s Snapshot) GetRange(r Range, options RangeOptions) RangeResult {
	return s.getRange(r, options, true)
}

// GetReadVersion is equivalent to (Transaction).GetReadVersion, performed as
// a snapshot read.
func (s Snapshot) GetReadVersion() FutureInt64 {
	return s.getReadVersion()
}

// GetDatabase returns a handle to the database with which this snapshot is
// interacting.
func (s Snapshot) GetDatabase() Database {
	return s.transaction.db
}

// GetEstimatedRangeSizeBytes returns an estimate for the number of bytes
// stored in the given range.
func (s Snapshot) GetEstimatedRangeSizeBytes(r ExactRange) FutureInt64 {
	beginKey, endKey := r.FDBRangeKeys()
	return s.getEstimatedRangeSizeBytes(
		beginKey.FDBKey(),
		endKey.FDBKey(),
	)
}

// GetRangeSplitPoints returns a list of keys that can split the given range
// into (roughly) equally sized chunks based on chunkSize.
// Note: the returned split points contain the start key and end key of the given range.
func (s Snapshot) GetRangeSplitPoints(r ExactRange, chunkSize int64) FutureKeyArray {
	beginKey, endKey := r.FDBRangeKeys()
	return s.getRangeSplitPoints(
		beginKey.FDBKey(),
		endKey.FDBKey(),
		chunkSize,
	)
}

// Snapshot returns the receiver and allows Snapshot to satisfy the
// ReadTransaction interface.
func (s Snapshot) Options() TransactionOptions {
	return s.Options()
}

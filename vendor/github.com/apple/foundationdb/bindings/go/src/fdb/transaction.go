/*
 * transaction.go
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

// #define FDB_API_VERSION 740
// #include <foundationdb/fdb_c.h>
import "C"

// A ReadTransaction can asynchronously read from a FoundationDB
// database. Transaction and Snapshot both satisfy the ReadTransaction
// interface.
//
// All ReadTransactions satisfy the ReadTransactor interface and may be used
// with read-only transactional functions.
type ReadTransaction interface {
	Get(key KeyConvertible) FutureByteSlice
	GetKey(sel Selectable) FutureKey
	GetRange(r Range, options RangeOptions) RangeResult
	GetReadVersion() FutureInt64
	GetDatabase() Database
	Snapshot() Snapshot
	GetEstimatedRangeSizeBytes(r ExactRange) FutureInt64
	GetRangeSplitPoints(r ExactRange, chunkSize int64) FutureKeyArray
	Options() TransactionOptions
	Cancel()

	ReadTransactor
}

// Transaction is a handle to a FoundationDB transaction. Transaction is a
// lightweight object that may be efficiently copied, and is safe for concurrent
// use by multiple goroutines.
//
// In FoundationDB, a transaction is a mutable snapshot of a database. All read
// and write operations on a transaction see and modify an otherwise-unchanging
// version of the database and only change the underlying database if and when
// the transaction is committed. Read operations do see the effects of previous
// write operations on the same transaction. Committing a transaction usually
// succeeds in the absence of conflicts.
//
// Transactions group operations into a unit with the properties of atomicity,
// isolation, and durability. Transactions also provide the ability to maintain
// an applications invariants or integrity constraints, supporting the property
// of consistency. Together these properties are known as ACID.
//
// Transactions are also causally consistent: once a transaction has been
// successfully committed, all subsequently created transactions will see the
// modifications made by it.
type Transaction struct {
	*transaction
}

type transaction struct {
	ptr *C.FDBTransaction
	db  Database
}

// TransactionOptions is a handle with which to set options that affect a
// Transaction object. A TransactionOptions instance should be obtained with the
// (Transaction).Options method.
type TransactionOptions struct {
	transaction *transaction
}

func (opt TransactionOptions) setOpt(code int, param []byte) error {
	return setOpt(func(p *C.uint8_t, pl C.int) C.fdb_error_t {
		return C.fdb_transaction_set_option(opt.transaction.ptr, C.FDBTransactionOption(code), p, pl)
	}, param)
}

func (t *transaction) destroy() {
	C.fdb_transaction_destroy(t.ptr)
}

func (t *transaction) cancel() {
	C.fdb_transaction_cancel(t.ptr)
}

// GetDatabase returns a handle to the database with which this transaction is
// interacting.
func (t Transaction) GetDatabase() Database {
	return t.transaction.db
}

// Transact executes the caller-provided function, passing it the Transaction
// receiver object.
//
// A panic of type Error during execution of the function will be recovered and
// returned to the caller as an error, but Transact will not retry the function
// or commit the Transaction after the caller-provided function completes.
//
// By satisfying the Transactor interface, Transaction may be passed to a
// transactional function from another transactional function, allowing
// composition. The outermost transactional function must have been provided a
// Database, or else the transaction will never be committed.
//
// See the Transactor interface for an example of using Transact with
// Transaction and Database objects.
func (t Transaction) Transact(f func(Transaction) (interface{}, error)) (r interface{}, e error) {
	defer panicToError(&e)

	r, e = f(t)
	return
}

// ReadTransact executes the caller-provided function, passing it the
// Transaction receiver object (as a ReadTransaction).
//
// A panic of type Error during execution of the function will be recovered and
// returned to the caller as an error, but ReadTransact will not retry the
// function.
//
// By satisfying the ReadTransactor interface, Transaction may be passed to a
// read-only transactional function from another (possibly read-only)
// transactional function, allowing composition.
//
// See the ReadTransactor interface for an example of using ReadTransact with
// Transaction, Snapshot and Database objects.
func (t Transaction) ReadTransact(f func(ReadTransaction) (interface{}, error)) (r interface{}, e error) {
	defer panicToError(&e)

	r, e = f(t)
	return
}

// Cancel cancels a transaction. All pending or future uses of the transaction
// will encounter an error. The Transaction object may be reused after calling
// (Transaction).Reset.
//
// Be careful if you are using (Transaction).Reset and (Transaction).Cancel
// concurrently with the same transaction. Since they negate each others
// effects, a race condition between these calls will leave the transaction in
// an unknown state.
//
// If your program attempts to cancel a transaction after (Transaction).Commit
// has been called but before it returns, unpredictable behavior will
// result. While it is guaranteed that the transaction will eventually end up in
// a cancelled state, the commit may or may not occur. Moreover, even if the
// call to (Transaction).Commit appears to return a transaction_cancelled
// error, the commit may have occurred or may occur in the future. This can make
// it more difficult to reason about the order in which transactions occur.
func (t Transaction) Cancel() {
	t.transaction.cancel()
}

// (Infrequently used) SetReadVersion sets the database version that the transaction will read from
// the database. The database cannot guarantee causal consistency if this method
// is used (the transactions reads will be causally consistent only if the
// provided read version has that property).
func (t Transaction) SetReadVersion(version int64) {
	C.fdb_transaction_set_read_version(t.ptr, C.int64_t(version))
}

// Snapshot returns a Snapshot object, suitable for performing snapshot
// reads. Snapshot reads offer a more relaxed isolation level than
// FoundationDB's default serializable isolation, reducing transaction conflicts
// but making it harder to reason about concurrency.
//
// For more information on snapshot reads, see
// https://apple.github.io/foundationdb/developer-guide.html#snapshot-reads.
func (t Transaction) Snapshot() Snapshot {
	return Snapshot{t.transaction}
}

// OnError determines whether an error returned by a Transaction method is
// retryable. Waiting on the returned future will return the same error when
// fatal, or return nil (after blocking the calling goroutine for a suitable
// delay) for retryable errors.
//
// Typical code will not use OnError directly. (Database).Transact uses
// OnError internally to implement a correct retry loop.
func (t Transaction) OnError(e Error) FutureNil {
	return &futureNil{
		future: newFuture(t.transaction, C.fdb_transaction_on_error(t.ptr, C.fdb_error_t(e.Code))),
	}
}

// Commit attempts to commit the modifications made in the transaction to the
// database. Waiting on the returned future will block the calling goroutine
// until the transaction has either been committed successfully or an error is
// encountered. Any error should be passed to (Transaction).OnError to determine
// if the error is retryable or not.
//
// As with other client/server databases, in some failure scenarios a client may
// be unable to determine whether a transaction succeeded. For more information,
// see
// https://apple.github.io/foundationdb/developer-guide.html#transactions-with-unknown-results.
func (t Transaction) Commit() FutureNil {
	return &futureNil{
		future: newFuture(t.transaction, C.fdb_transaction_commit(t.ptr)),
	}
}

// Watch creates a watch and returns a FutureNil that will become ready when the
// watch reports a change to the value of the specified key.
//
// A watchs behavior is relative to the transaction that created it. A watch
// will report a change in relation to the keys value as readable by that
// transaction. The initial value used for comparison is either that of the
// transactions read version or the value as modified by the transaction itself
// prior to the creation of the watch. If the value changes and then changes
// back to its initial value, the watch might not report the change.
//
// Until the transaction that created it has been committed, a watch will not
// report changes made by other transactions. In contrast, a watch will
// immediately report changes made by the transaction itself. Watches cannot be
// created if the transaction has called SetReadYourWritesDisable on the
// Transaction options, and an attempt to do so will return a watches_disabled
// error.
//
// If the transaction used to create a watch encounters an error during commit,
// then the watch will be set with that error. A transaction whose commit
// result is unknown will set all of its watches with the commit_unknown_result
// error. If an uncommitted transaction is reset or destroyed, then any watches
// it created will be set with the transaction_cancelled error.
//
// By default, each database connection can have no more than 10,000 watches
// that have not yet reported a change. When this number is exceeded, an attempt
// to create a watch will return a too_many_watches error. This limit can be
// changed using SetMaxWatches on the Database. Because a watch outlives the
// transaction that creates it, any watch that is no longer needed should be
// cancelled by calling (FutureNil).Cancel on its returned future.
func (t Transaction) Watch(key KeyConvertible) FutureNil {
	kb := key.FDBKey()
	return &futureNil{
		future: newFuture(t.transaction, C.fdb_transaction_watch(t.ptr, byteSliceToPtr(kb), C.int(len(kb)))),
	}
}

func (t *transaction) get(key []byte, snapshot int) FutureByteSlice {
	return &futureByteSlice{
		future: newFuture(t, C.fdb_transaction_get(
			t.ptr,
			byteSliceToPtr(key),
			C.int(len(key)),
			C.fdb_bool_t(snapshot),
		)),
	}
}

// Get returns the (future) value associated with the specified key. The read is
// performed asynchronously and does not block the calling goroutine. The future
// will become ready when the read is complete.
func (t Transaction) Get(key KeyConvertible) FutureByteSlice {
	return t.get(key.FDBKey(), 0)
}

func (t *transaction) doGetRange(r Range, options RangeOptions, snapshot bool, iteration int) futureKeyValueArray {
	begin, end := r.FDBRangeKeySelectors()
	bsel := begin.FDBKeySelector()
	esel := end.FDBKeySelector()
	bkey := bsel.Key.FDBKey()
	ekey := esel.Key.FDBKey()

	return futureKeyValueArray{
		future: newFuture(t, C.fdb_transaction_get_range(
			t.ptr,
			byteSliceToPtr(bkey),
			C.int(len(bkey)),
			C.fdb_bool_t(boolToInt(bsel.OrEqual)),
			C.int(bsel.Offset),
			byteSliceToPtr(ekey),
			C.int(len(ekey)),
			C.fdb_bool_t(boolToInt(esel.OrEqual)),
			C.int(esel.Offset),
			C.int(options.Limit),
			C.int(0),
			C.FDBStreamingMode(options.Mode-1),
			C.int(iteration),
			C.fdb_bool_t(boolToInt(snapshot)),
			C.fdb_bool_t(boolToInt(options.Reverse)),
		))}
}

func (t *transaction) getRange(r Range, options RangeOptions, snapshot bool) RangeResult {
	f := t.doGetRange(r, options, snapshot, 1)
	begin, end := r.FDBRangeKeySelectors()
	return RangeResult{
		t:        t,
		sr:       SelectorRange{begin, end},
		options:  options,
		snapshot: snapshot,
		f:        &f,
	}
}

// GetRange performs a range read. The returned RangeResult represents all
// KeyValue objects kv where beginKey <= kv.Key < endKey, ordered by kv.Key
// (where beginKey and endKey are the keys described by the key selectors
// returned by r.FDBKeySelectors). All reads performed as a result of GetRange
// are asynchronous and do not block the calling goroutine.
func (t Transaction) GetRange(r Range, options RangeOptions) RangeResult {
	return t.getRange(r, options, false)
}

func (t *transaction) getEstimatedRangeSizeBytes(beginKey Key, endKey Key) FutureInt64 {
	return &futureInt64{
		future: newFuture(t, C.fdb_transaction_get_estimated_range_size_bytes(
			t.ptr,
			byteSliceToPtr(beginKey),
			C.int(len(beginKey)),
			byteSliceToPtr(endKey),
			C.int(len(endKey)),
		)),
	}
}

// GetEstimatedRangeSizeBytes returns an estimate for the number of bytes
// stored in the given range.
// Note: the estimated size is calculated based on the sampling done by FDB server. The sampling
// algorithm works roughly in this way: the larger the key-value pair is, the more likely it would
// be sampled and the more accurate its sampled size would be. And due to
// that reason it is recommended to use this API to query against large ranges for accuracy considerations.
// For a rough reference, if the returned size is larger than 3MB, one can consider the size to be
// accurate.
func (t Transaction) GetEstimatedRangeSizeBytes(r ExactRange) FutureInt64 {
	beginKey, endKey := r.FDBRangeKeys()
	return t.getEstimatedRangeSizeBytes(
		beginKey.FDBKey(),
		endKey.FDBKey(),
	)
}

func (t *transaction) getRangeSplitPoints(beginKey Key, endKey Key, chunkSize int64) FutureKeyArray {
	return &futureKeyArray{
		future: newFuture(t, C.fdb_transaction_get_range_split_points(
			t.ptr,
			byteSliceToPtr(beginKey),
			C.int(len(beginKey)),
			byteSliceToPtr(endKey),
			C.int(len(endKey)),
			C.int64_t(chunkSize),
		)),
	}
}

// GetRangeSplitPoints returns a list of keys that can split the given range
// into (roughly) equally sized chunks based on chunkSize.
// Note: the returned split points contain the start key and end key of the given range.
func (t Transaction) GetRangeSplitPoints(r ExactRange, chunkSize int64) FutureKeyArray {
	beginKey, endKey := r.FDBRangeKeys()
	return t.getRangeSplitPoints(
		beginKey.FDBKey(),
		endKey.FDBKey(),
		chunkSize,
	)
}

func (t *transaction) getReadVersion() FutureInt64 {
	return &futureInt64{
		future: newFuture(t, C.fdb_transaction_get_read_version(t.ptr)),
	}
}

// (Infrequently used) GetReadVersion returns the (future) transaction read version. The read is
// performed asynchronously and does not block the calling goroutine. The future
// will become ready when the read version is available.
func (t Transaction) GetReadVersion() FutureInt64 {
	return t.getReadVersion()
}

// Set associated the given key and value, overwriting any previous association
// with key. Set returns immediately, having modified the snapshot of the
// database represented by the transaction.
func (t Transaction) Set(key KeyConvertible, value []byte) {
	kb := key.FDBKey()
	C.fdb_transaction_set(t.ptr, byteSliceToPtr(kb), C.int(len(kb)), byteSliceToPtr(value), C.int(len(value)))
}

// Clear removes the specified key (and any associated value), if it
// exists. Clear returns immediately, having modified the snapshot of the
// database represented by the transaction.
func (t Transaction) Clear(key KeyConvertible) {
	kb := key.FDBKey()
	C.fdb_transaction_clear(t.ptr, byteSliceToPtr(kb), C.int(len(kb)))
}

// ClearRange removes all keys k such that begin <= k < end, and their
// associated values. ClearRange returns immediately, having modified the
// snapshot of the database represented by the transaction.
// Range clears are efficient with FoundationDB -- clearing large amounts of data
// will be fast. However, this will not immediately free up disk -
// data for the deleted range is cleaned up in the background.
// For purposes of computing the transaction size, only the begin and end keys of a clear range are counted.
// The size of the data stored in the range does not count against the transaction size limit.
func (t Transaction) ClearRange(er ExactRange) {
	begin, end := er.FDBRangeKeys()
	bkb := begin.FDBKey()
	ekb := end.FDBKey()
	C.fdb_transaction_clear_range(t.ptr, byteSliceToPtr(bkb), C.int(len(bkb)), byteSliceToPtr(ekb), C.int(len(ekb)))
}

// (Infrequently used) GetCommittedVersion returns the version number at which a
// successful commit modified the database. This must be called only after the
// successful (non-error) completion of a call to Commit on this Transaction, or
// the behavior is undefined. Read-only transactions do not modify the database
// when committed and will have a committed version of -1. Keep in mind that a
// transaction which reads keys and then sets them to their current values may
// be optimized to a read-only transaction.
func (t Transaction) GetCommittedVersion() (int64, error) {
	var version C.int64_t

	if err := C.fdb_transaction_get_committed_version(t.ptr, &version); err != 0 {
		return 0, Error{int(err)}
	}

	return int64(version), nil
}

// (Infrequently used) Returns a future which will contain the versionstamp
// which was used by any versionstamp operations in this transaction. The
// future will be ready only after the successful completion of a call to Commit
// on this Transaction. Read-only transactions do not modify the database when
// committed and will result in the future completing with an error. Keep in
// mind that a transaction which reads keys and then sets them to their current
// values may be optimized to a read-only transaction.
func (t Transaction) GetVersionstamp() FutureKey {
	return &futureKey{future: newFuture(t.transaction, C.fdb_transaction_get_versionstamp(t.ptr))}
}

func (t *transaction) getApproximateSize() FutureInt64 {
	return &futureInt64{
		future: newFuture(t, C.fdb_transaction_get_approximate_size(t.ptr)),
	}
}

// Returns a future that is the approximate transaction size so far in this
// transaction, which is the summation of the estimated size of mutations,
// read conflict ranges, and write conflict ranges.
func (t Transaction) GetApproximateSize() FutureInt64 {
	return t.getApproximateSize()
}

// Reset rolls back a transaction, completely resetting it to its initial
// state. This is logically equivalent to destroying the transaction and
// creating a new one.
func (t Transaction) Reset() {
	C.fdb_transaction_reset(t.ptr)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (t *transaction) getKey(sel KeySelector, snapshot int) FutureKey {
	key := sel.Key.FDBKey()
	return &futureKey{
		future: newFuture(t, C.fdb_transaction_get_key(
			t.ptr,
			byteSliceToPtr(key),
			C.int(len(key)),
			C.fdb_bool_t(boolToInt(sel.OrEqual)),
			C.int(sel.Offset),
			C.fdb_bool_t(snapshot),
		)),
	}
}

// GetKey returns the future key referenced by the provided key selector. The
// read is performed asynchronously and does not block the calling
// goroutine. The future will become ready when the read version is available.
//
// By default, the key is cached for the duration of the transaction, providing
// a potential performance benefit. However, the value of the key is also
// retrieved, using network bandwidth. Invoking
// (TransactionOptions).SetReadYourWritesDisable will avoid both the caching and
// the increased network bandwidth.
func (t Transaction) GetKey(sel Selectable) FutureKey {
	return t.getKey(sel.FDBKeySelector(), 0)
}

func (t Transaction) atomicOp(key []byte, param []byte, code int) {
	C.fdb_transaction_atomic_op(
		t.ptr,
		byteSliceToPtr(key),
		C.int(len(key)),
		byteSliceToPtr(param),
		C.int(len(param)),
		C.FDBMutationType(code),
	)
}

func addConflictRange(t *transaction, er ExactRange, crtype conflictRangeType) error {
	begin, end := er.FDBRangeKeys()
	bkb := begin.FDBKey()
	ekb := end.FDBKey()
	if err := C.fdb_transaction_add_conflict_range(
		t.ptr,
		byteSliceToPtr(bkb),
		C.int(len(bkb)),
		byteSliceToPtr(ekb),
		C.int(len(ekb)),
		C.FDBConflictRangeType(crtype),
	); err != 0 {
		return Error{int(err)}
	}

	return nil
}

// AddReadConflictRange adds a range of keys to the transactions read conflict
// ranges as if you had read the range. As a result, other transactions that
// write a key in this range could cause the transaction to fail with a
// conflict.
//
// For more information on conflict ranges, see
// https://apple.github.io/foundationdb/developer-guide.html#conflict-ranges.
func (t Transaction) AddReadConflictRange(er ExactRange) error {
	return addConflictRange(t.transaction, er, conflictRangeTypeRead)
}

func copyAndAppend(orig []byte, b byte) []byte {
	ret := make([]byte, len(orig)+1)
	copy(ret, orig)
	ret[len(orig)] = b
	return ret
}

// AddReadConflictKey adds a key to the transactions read conflict ranges as if
// you had read the key. As a result, other transactions that concurrently write
// this key could cause the transaction to fail with a conflict.
//
// For more information on conflict ranges, see
// https://apple.github.io/foundationdb/developer-guide.html#conflict-ranges.
func (t Transaction) AddReadConflictKey(key KeyConvertible) error {
	return addConflictRange(
		t.transaction,
		KeyRange{key, Key(copyAndAppend(key.FDBKey(), 0x00))},
		conflictRangeTypeRead,
	)
}

// AddWriteConflictRange adds a range of keys to the transactions write
// conflict ranges as if you had cleared the range. As a result, other
// transactions that concurrently read a key in this range could fail with a
// conflict.
//
// For more information on conflict ranges, see
// https://apple.github.io/foundationdb/developer-guide.html#conflict-ranges.
func (t Transaction) AddWriteConflictRange(er ExactRange) error {
	return addConflictRange(t.transaction, er, conflictRangeTypeWrite)
}

// AddWriteConflictKey adds a key to the transactions write conflict ranges as
// if you had written the key. As a result, other transactions that concurrently
// read this key could fail with a conflict.
//
// For more information on conflict ranges, see
// https://apple.github.io/foundationdb/developer-guide.html#conflict-ranges.
func (t Transaction) AddWriteConflictKey(key KeyConvertible) error {
	return addConflictRange(
		t.transaction,
		KeyRange{key, Key(copyAndAppend(key.FDBKey(), 0x00))},
		conflictRangeTypeWrite,
	)
}

// Options returns a TransactionOptions instance suitable for setting options
// specific to this transaction.
func (t Transaction) Options() TransactionOptions {
	return TransactionOptions{t.transaction}
}

func localityGetAddressesForKey(t *transaction, key KeyConvertible) FutureStringSlice {
	kb := key.FDBKey()
	return &futureStringSlice{
		future: newFuture(t, C.fdb_transaction_get_addresses_for_key(
			t.ptr,
			byteSliceToPtr(kb),
			C.int(len(kb)),
		)),
	}
}

// LocalityGetAddressesForKey returns the (future) public network addresses of
// each of the storage servers responsible for storing key and its associated
// value. The read is performed asynchronously and does not block the calling
// goroutine. The future will become ready when the read is complete.
func (t Transaction) LocalityGetAddressesForKey(key KeyConvertible) FutureStringSlice {
	return localityGetAddressesForKey(t.transaction, key)
}

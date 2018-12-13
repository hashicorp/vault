/*
 * Copyright 2017 Dgraph Labs, Inc. and Contributors
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

package badger

import (
	"bytes"
	"context"
	"encoding/hex"
	"math"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/y"
	farm "github.com/dgryski/go-farm"
	"github.com/pkg/errors"
)

type oracle struct {
	// A 64-bit integer must be at the top for memory alignment. See issue #311.
	refCount  int64
	isManaged bool // Does not change value, so no locking required.

	sync.Mutex // For nextTxnTs and commits.
	// writeChLock lock is for ensuring that transactions go to the write
	// channel in the same order as their commit timestamps.
	writeChLock sync.Mutex
	nextTxnTs   uint64

	// Used to block NewTransaction, so all previous commits are visible to a new read.
	txnMark *y.WaterMark

	// Either of these is used to determine which versions can be permanently
	// discarded during compaction.
	discardTs uint64       // Used by ManagedDB.
	readMark  *y.WaterMark // Used by DB.

	// commits stores a key fingerprint and latest commit counter for it.
	// refCount is used to clear out commits map to avoid a memory blowup.
	commits map[uint64]uint64
}

func newOracle(opt Options) *oracle {
	orc := &oracle{
		isManaged: opt.managedTxns,
		commits:   make(map[uint64]uint64),
		// We're not initializing nextTxnTs and readOnlyTs. It would be done after replay in Open.
		//
		// WaterMarks must be 64-bit aligned for atomic package, hence we must use pointers here.
		// See https://golang.org/pkg/sync/atomic/#pkg-note-BUG.
		readMark: &y.WaterMark{Name: "badger.PendingReads"},
		txnMark:  &y.WaterMark{Name: "badger.TxnTimestamp"},
	}
	orc.readMark.Init()
	orc.txnMark.Init()
	return orc
}

func (o *oracle) addRef() {
	atomic.AddInt64(&o.refCount, 1)
}

func (o *oracle) decrRef() {
	if atomic.AddInt64(&o.refCount, -1) != 0 {
		return
	}

	// Clear out commits maps to release memory.
	o.Lock()
	defer o.Unlock()
	// Avoids the race where something new is added to commitsMap
	// after we check refCount and before we take Lock.
	if atomic.LoadInt64(&o.refCount) != 0 {
		return
	}
	if len(o.commits) >= 1000 { // If the map is still small, let it slide.
		o.commits = make(map[uint64]uint64)
	}
}

func (o *oracle) readTs() uint64 {
	if o.isManaged {
		panic("ReadTs should not be retrieved for managed DB")
	}

	var readTs uint64
	o.Lock()
	readTs = o.nextTxnTs - 1
	o.readMark.Begin(readTs)
	o.Unlock()

	// Wait for all txns which have no conflicts, have been assigned a commit
	// timestamp and are going through the write to value log and LSM tree
	// process. Not waiting here could mean that some txns which have been
	// committed would not be read.
	y.Check(o.txnMark.WaitForMark(context.Background(), readTs))
	return readTs
}

func (o *oracle) nextTs() uint64 {
	o.Lock()
	defer o.Unlock()
	return o.nextTxnTs
}

// Any deleted or invalid versions at or below ts would be discarded during
// compaction to reclaim disk space in LSM tree and thence value log.
func (o *oracle) setDiscardTs(ts uint64) {
	o.Lock()
	defer o.Unlock()
	o.discardTs = ts
}

func (o *oracle) discardAtOrBelow() uint64 {
	if o.isManaged {
		o.Lock()
		defer o.Unlock()
		return o.discardTs
	}
	return o.readMark.DoneUntil()
}

// hasConflict must be called while having a lock.
func (o *oracle) hasConflict(txn *Txn) bool {
	if len(txn.reads) == 0 {
		return false
	}
	for _, ro := range txn.reads {
		// A commit at the read timestamp is expected.
		// But, any commit after the read timestamp should cause a conflict.
		if ts, has := o.commits[ro]; has && ts > txn.readTs {
			return true
		}
	}
	return false
}

func (o *oracle) newCommitTs(txn *Txn) uint64 {
	o.Lock()
	defer o.Unlock()

	if o.hasConflict(txn) {
		return 0
	}

	var ts uint64
	if !o.isManaged {
		// This is the general case, when user doesn't specify the read and commit ts.
		ts = o.nextTxnTs
		o.nextTxnTs++
		o.txnMark.Begin(ts)

	} else {
		// If commitTs is set, use it instead.
		ts = txn.commitTs
	}

	for _, w := range txn.writes {
		o.commits[w] = ts // Update the commitTs.
	}
	return ts
}

func (o *oracle) doneCommit(cts uint64) {
	if o.isManaged {
		// No need to update anything.
		return
	}
	o.txnMark.Done(cts)
}

// Txn represents a Badger transaction.
type Txn struct {
	readTs   uint64
	commitTs uint64

	update bool     // update is used to conditionally keep track of reads.
	reads  []uint64 // contains fingerprints of keys read.
	writes []uint64 // contains fingerprints of keys written.

	pendingWrites map[string]*Entry // cache stores any writes done by txn.

	db        *DB
	discarded bool

	size         int64
	count        int64
	numIterators int32
}

type pendingWritesIterator struct {
	entries  []*Entry
	nextIdx  int
	readTs   uint64
	reversed bool
}

func (pi *pendingWritesIterator) Next() {
	pi.nextIdx++
}

func (pi *pendingWritesIterator) Rewind() {
	pi.nextIdx = 0
}

func (pi *pendingWritesIterator) Seek(key []byte) {
	key = y.ParseKey(key)
	pi.nextIdx = sort.Search(len(pi.entries), func(idx int) bool {
		cmp := bytes.Compare(pi.entries[idx].Key, key)
		if !pi.reversed {
			return cmp >= 0
		}
		return cmp <= 0
	})
}

func (pi *pendingWritesIterator) Key() []byte {
	y.AssertTrue(pi.Valid())
	entry := pi.entries[pi.nextIdx]
	return y.KeyWithTs(entry.Key, pi.readTs)
}

func (pi *pendingWritesIterator) Value() y.ValueStruct {
	y.AssertTrue(pi.Valid())
	entry := pi.entries[pi.nextIdx]
	return y.ValueStruct{
		Value:     entry.Value,
		Meta:      entry.meta,
		UserMeta:  entry.UserMeta,
		ExpiresAt: entry.ExpiresAt,
		Version:   pi.readTs,
	}
}

func (pi *pendingWritesIterator) Valid() bool {
	return pi.nextIdx < len(pi.entries)
}

func (pi *pendingWritesIterator) Close() error {
	return nil
}

func (txn *Txn) newPendingWritesIterator(reversed bool) *pendingWritesIterator {
	if !txn.update || len(txn.pendingWrites) == 0 {
		return nil
	}
	entries := make([]*Entry, 0, len(txn.pendingWrites))
	for _, e := range txn.pendingWrites {
		entries = append(entries, e)
	}
	// Number of pending writes per transaction shouldn't be too big in general.
	sort.Slice(entries, func(i, j int) bool {
		cmp := bytes.Compare(entries[i].Key, entries[j].Key)
		if !reversed {
			return cmp < 0
		}
		return cmp > 0
	})
	return &pendingWritesIterator{
		readTs:   txn.readTs,
		entries:  entries,
		reversed: reversed,
	}
}

func (txn *Txn) checkSize(e *Entry) error {
	count := txn.count + 1
	// Extra bytes for version in key.
	size := txn.size + int64(e.estimateSize(txn.db.opt.ValueThreshold)) + 10
	if count >= txn.db.opt.maxBatchCount || size >= txn.db.opt.maxBatchSize {
		return ErrTxnTooBig
	}
	txn.count, txn.size = count, size
	return nil
}

// Set adds a key-value pair to the database.
//
// It will return ErrReadOnlyTxn if update flag was set to false when creating the
// transaction.
//
// The current transaction keeps a reference to the key and val byte slice
// arguments. Users must not modify key and val until the end of the transaction.
func (txn *Txn) Set(key, val []byte) error {
	e := &Entry{
		Key:   key,
		Value: val,
	}
	return txn.SetEntry(e)
}

// SetWithMeta adds a key-value pair to the database, along with a metadata
// byte.
//
// This byte is stored alongside the key, and can be used as an aid to
// interpret the value or store other contextual bits corresponding to the
// key-value pair.
//
// The current transaction keeps a reference to the key and val byte slice
// arguments. Users must not modify key and val until the end of the transaction.
func (txn *Txn) SetWithMeta(key, val []byte, meta byte) error {
	e := &Entry{Key: key, Value: val, UserMeta: meta}
	return txn.SetEntry(e)
}

// SetWithDiscard acts like SetWithMeta, but adds a marker to discard earlier
// versions of the key.
//
// This method is only useful if you have set a higher limit for
// options.NumVersionsToKeep. The default setting is 1, in which case, this
// function doesn't add any more benefit than just calling the normal
// SetWithMeta (or Set) function. If however, you have a higher setting for
// NumVersionsToKeep (in Dgraph, we set it to infinity), you can use this method
// to indicate that all the older versions can be discarded and removed during
// compactions.
//
// The current transaction keeps a reference to the key and val byte slice
// arguments. Users must not modify key and val until the end of the
// transaction.
func (txn *Txn) SetWithDiscard(key, val []byte, meta byte) error {
	e := &Entry{
		Key:      key,
		Value:    val,
		UserMeta: meta,
		meta:     bitDiscardEarlierVersions,
	}
	return txn.SetEntry(e)
}

// SetWithTTL adds a key-value pair to the database, along with a time-to-live
// (TTL) setting. A key stored with a TTL would automatically expire after the
// time has elapsed , and be eligible for garbage collection.
//
// The current transaction keeps a reference to the key and val byte slice
// arguments. Users must not modify key and val until the end of the
// transaction.
func (txn *Txn) SetWithTTL(key, val []byte, dur time.Duration) error {
	expire := time.Now().Add(dur).Unix()
	e := &Entry{Key: key, Value: val, ExpiresAt: uint64(expire)}
	return txn.SetEntry(e)
}

func exceedsSize(prefix string, max int64, key []byte) error {
	return errors.Errorf("%s with size %d exceeded %d limit. %s:\n%s",
		prefix, len(key), max, prefix, hex.Dump(key[:1<<10]))
}

func (txn *Txn) modify(e *Entry) error {
	const maxKeySize = 65000

	switch {
	case !txn.update:
		return ErrReadOnlyTxn
	case txn.discarded:
		return ErrDiscardedTxn
	case len(e.Key) == 0:
		return ErrEmptyKey
	case bytes.HasPrefix(e.Key, badgerPrefix):
		return ErrInvalidKey
	case len(e.Key) > maxKeySize:
		// Key length can't be more than uint16, as determined by table::header.  To
		// keep things safe and allow badger move prefix and a timestamp suffix, let's
		// cut it down to 65000, instead of using 65536.
		return exceedsSize("Key", maxKeySize, e.Key)
	case int64(len(e.Value)) > txn.db.opt.ValueLogFileSize:
		return exceedsSize("Value", txn.db.opt.ValueLogFileSize, e.Value)
	}

	if err := txn.checkSize(e); err != nil {
		return err
	}
	fp := farm.Fingerprint64(e.Key) // Avoid dealing with byte arrays.
	txn.writes = append(txn.writes, fp)
	txn.pendingWrites[string(e.Key)] = e
	return nil
}

// SetEntry takes an Entry struct and adds the key-value pair in the struct,
// along with other metadata to the database.
//
// The current transaction keeps a reference to the entry passed in argument.
// Users must not modify the entry until the end of the transaction.
func (txn *Txn) SetEntry(e *Entry) error {
	return txn.modify(e)
}

// Delete deletes a key.
//
// This is done by adding a delete marker for the key at commit timestamp.  Any
// reads happening before this timestamp would be unaffected. Any reads after
// this commit would see the deletion.
//
// The current transaction keeps a reference to the key byte slice argument.
// Users must not modify the key until the end of the transaction.
func (txn *Txn) Delete(key []byte) error {
	e := &Entry{
		Key:  key,
		meta: bitDelete,
	}
	return txn.modify(e)
}

// Get looks for key and returns corresponding Item.
// If key is not found, ErrKeyNotFound is returned.
func (txn *Txn) Get(key []byte) (item *Item, rerr error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	} else if txn.discarded {
		return nil, ErrDiscardedTxn
	}

	item = new(Item)
	if txn.update {
		if e, has := txn.pendingWrites[string(key)]; has && bytes.Equal(key, e.Key) {
			if isDeletedOrExpired(e.meta, e.ExpiresAt) {
				return nil, ErrKeyNotFound
			}
			// Fulfill from cache.
			item.meta = e.meta
			item.val = e.Value
			item.userMeta = e.UserMeta
			item.key = key
			item.status = prefetched
			item.version = txn.readTs
			item.expiresAt = e.ExpiresAt
			// We probably don't need to set db on item here.
			return item, nil
		}
		// Only track reads if this is update txn. No need to track read if txn serviced it
		// internally.
		txn.addReadKey(key)
	}

	seek := y.KeyWithTs(key, txn.readTs)
	vs, err := txn.db.get(seek)
	if err != nil {
		return nil, errors.Wrapf(err, "DB::Get key: %q", key)
	}
	if vs.Value == nil && vs.Meta == 0 {
		return nil, ErrKeyNotFound
	}
	if isDeletedOrExpired(vs.Meta, vs.ExpiresAt) {
		return nil, ErrKeyNotFound
	}

	item.key = key
	item.version = vs.Version
	item.meta = vs.Meta
	item.userMeta = vs.UserMeta
	item.db = txn.db
	item.vptr = vs.Value // TODO: Do we need to copy this over?
	item.txn = txn
	item.expiresAt = vs.ExpiresAt
	return item, nil
}

func (txn *Txn) addReadKey(key []byte) {
	if txn.update {
		fp := farm.Fingerprint64(key)
		txn.reads = append(txn.reads, fp)
	}
}

// Discard discards a created transaction. This method is very important and must be called. Commit
// method calls this internally, however, calling this multiple times doesn't cause any issues. So,
// this can safely be called via a defer right when transaction is created.
//
// NOTE: If any operations are run on a discarded transaction, ErrDiscardedTxn is returned.
func (txn *Txn) Discard() {
	if txn.discarded { // Avoid a re-run.
		return
	}
	if atomic.LoadInt32(&txn.numIterators) > 0 {
		panic("Unclosed iterator at time of Txn.Discard.")
	}
	txn.discarded = true
	if !txn.db.orc.isManaged {
		txn.db.orc.readMark.Done(txn.readTs)
	}
	if txn.update {
		txn.db.orc.decrRef()
	}
}

func (txn *Txn) commitAndSend() (func() error, error) {
	orc := txn.db.orc
	// Ensure that the order in which we get the commit timestamp is the same as
	// the order in which we push these updates to the write channel. So, we
	// acquire a writeChLock before getting a commit timestamp, and only release
	// it after pushing the entries to it.
	orc.writeChLock.Lock()
	defer orc.writeChLock.Unlock()

	commitTs := orc.newCommitTs(txn)
	if commitTs == 0 {
		return nil, ErrConflict
	}

	// The following debug information is what led to determining the cause of
	// bank txn violation bug, and it took a whole bunch of effort to narrow it
	// down to here. So, keep this around for at least a couple of months.
	// var b strings.Builder
	// fmt.Fprintf(&b, "Read: %d. Commit: %d. reads: %v. writes: %v. Keys: ",
	// 	txn.readTs, commitTs, txn.reads, txn.writes)
	entries := make([]*Entry, 0, len(txn.pendingWrites)+1)
	for _, e := range txn.pendingWrites {
		// fmt.Fprintf(&b, "[%q : %q], ", e.Key, e.Value)

		// Suffix the keys with commit ts, so the key versions are sorted in
		// descending order of commit timestamp.
		e.Key = y.KeyWithTs(e.Key, commitTs)
		e.meta |= bitTxn
		entries = append(entries, e)
	}
	// log.Printf("%s\n", b.String())
	e := &Entry{
		Key:   y.KeyWithTs(txnKey, commitTs),
		Value: []byte(strconv.FormatUint(commitTs, 10)),
		meta:  bitFinTxn,
	}
	entries = append(entries, e)

	req, err := txn.db.sendToWriteCh(entries)
	if err != nil {
		orc.doneCommit(commitTs)
		return nil, err
	}
	ret := func() error {
		err := req.Wait()
		// Wait before marking commitTs as done.
		// We can't defer doneCommit above, because it is being called from a
		// callback here.
		orc.doneCommit(commitTs)
		return err
	}
	return ret, nil
}

func (txn *Txn) commitPrecheck() {
	if txn.commitTs == 0 && txn.db.opt.managedTxns {
		panic("Commit cannot be called with managedDB=true. Use CommitAt.")
	}
	if txn.discarded {
		panic("Trying to commit a discarded txn")
	}
}

// Commit commits the transaction, following these steps:
//
// 1. If there are no writes, return immediately.
//
// 2. Check if read rows were updated since txn started. If so, return ErrConflict.
//
// 3. If no conflict, generate a commit timestamp and update written rows' commit ts.
//
// 4. Batch up all writes, write them to value log and LSM tree.
//
// 5. If callback is provided, Badger will return immediately after checking
// for conflicts. Writes to the database will happen in the background.  If
// there is a conflict, an error will be returned and the callback will not
// run. If there are no conflicts, the callback will be called in the
// background upon successful completion of writes or any error during write.
//
// If error is nil, the transaction is successfully committed. In case of a non-nil error, the LSM
// tree won't be updated, so there's no need for any rollback.
func (txn *Txn) Commit() error {
	txn.commitPrecheck() // Precheck before discarding txn.
	defer txn.Discard()

	if len(txn.writes) == 0 {
		return nil // Nothing to do.
	}

	txnCb, err := txn.commitAndSend()
	if err != nil {
		return err
	}
	// If batchSet failed, LSM would not have been updated. So, no need to rollback anything.

	// TODO: What if some of the txns successfully make it to value log, but others fail.
	// Nothing gets updated to LSM, until a restart happens.
	return txnCb()
}

type txnCb struct {
	commit func() error
	user   func(error)
	err    error
}

func runTxnCallback(cb *txnCb) {
	switch {
	case cb == nil:
		panic("txn callback is nil")
	case cb.err != nil:
		cb.user(cb.err)
	case cb.commit != nil:
		err := cb.commit()
		cb.user(err)
	case cb.user == nil:
		panic("Must have caught a nil callback for txn.CommitWith")
	default:
		cb.user(nil)
	}
}

// CommitWith acts like Commit, but takes a callback, which gets run via a
// goroutine to avoid blocking this function. The callback is guaranteed to run,
// so it is safe to increment sync.WaitGroup before calling CommitWith, and
// decrementing it in the callback; to block until all callbacks are run.
func (txn *Txn) CommitWith(cb func(error)) {
	txn.commitPrecheck() // Precheck before discarding txn.
	defer txn.Discard()

	if cb == nil {
		panic("Nil callback provided to CommitWith")
	}

	if len(txn.writes) == 0 {
		// Do not run these callbacks from here, because the CommitWith and the
		// callback might be acquiring the same locks. Instead run the callback
		// from another goroutine.
		go runTxnCallback(&txnCb{user: cb, err: nil})
		return
	}

	commitCb, err := txn.commitAndSend()
	if err != nil {
		go runTxnCallback(&txnCb{user: cb, err: err})
		return
	}

	go runTxnCallback(&txnCb{user: cb, commit: commitCb})
}

// ReadTs returns the read timestamp of the transaction.
func (txn *Txn) ReadTs() uint64 {
	return txn.readTs
}

// NewTransaction creates a new transaction. Badger supports concurrent execution of transactions,
// providing serializable snapshot isolation, avoiding write skews. Badger achieves this by tracking
// the keys read and at Commit time, ensuring that these read keys weren't concurrently modified by
// another transaction.
//
// For read-only transactions, set update to false. In this mode, we don't track the rows read for
// any changes. Thus, any long running iterations done in this mode wouldn't pay this overhead.
//
// Running transactions concurrently is OK. However, a transaction itself isn't thread safe, and
// should only be run serially. It doesn't matter if a transaction is created by one goroutine and
// passed down to other, as long as the Txn APIs are called serially.
//
// When you create a new transaction, it is absolutely essential to call
// Discard(). This should be done irrespective of what the update param is set
// to. Commit API internally runs Discard, but running it twice wouldn't cause
// any issues.
//
//  txn := db.NewTransaction(false)
//  defer txn.Discard()
//  // Call various APIs.
func (db *DB) NewTransaction(update bool) *Txn {
	return db.newTransaction(update, false)
}

func (db *DB) newTransaction(update, isManaged bool) *Txn {
	if db.opt.ReadOnly && update {
		// DB is read-only, force read-only transaction.
		update = false
	}

	txn := &Txn{
		update: update,
		db:     db,
		count:  1,                       // One extra entry for BitFin.
		size:   int64(len(txnKey) + 10), // Some buffer for the extra entry.
	}
	if update {
		txn.pendingWrites = make(map[string]*Entry)
		txn.db.orc.addRef()
	}
	// It is important that the oracle addRef happens BEFORE we retrieve a read
	// timestamp. Otherwise, it is possible that the oracle commit map would
	// become nil after we get the read timestamp.
	// The sequence of events can be:
	// 1. This txn gets a read timestamp.
	// 2. Another txn working on the same keyset commits them, and decrements
	//    the reference to oracle.
	// 3. Oracle ref reaches zero, resetting commit map.
	// 4. This txn increments the oracle reference.
	// 5. Now this txn would go on to commit the keyset, and no conflicts
	//    would be detected.
	// See issue: https://github.com/dgraph-io/badger/issues/574
	if !isManaged {
		txn.readTs = db.orc.readTs()
	}
	return txn
}

// View executes a function creating and managing a read-only transaction for the user. Error
// returned by the function is relayed by the View method.
// If View is used with managed transactions, it would assume a read timestamp of MaxUint64.
func (db *DB) View(fn func(txn *Txn) error) error {
	var txn *Txn
	if db.opt.managedTxns {
		txn = db.NewTransactionAt(math.MaxUint64, false)
	} else {
		txn = db.NewTransaction(false)
	}
	defer txn.Discard()

	return fn(txn)
}

// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// Update cannot be used with managed transactions.
func (db *DB) Update(fn func(txn *Txn) error) error {
	if db.opt.managedTxns {
		panic("Update can only be used with managedDB=false.")
	}
	txn := db.NewTransaction(true)
	defer txn.Discard()

	if err := fn(txn); err != nil {
		return err
	}

	return txn.Commit()
}

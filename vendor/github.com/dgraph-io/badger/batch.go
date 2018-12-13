/*
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
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
	"sync"
)

type WriteBatch struct {
	sync.Mutex
	txn *Txn
	db  *DB
	wg  sync.WaitGroup
	err error
}

// NewWriteBatch creates a new WriteBatch. This provides a way to conveniently do a lot of writes,
// batching them up as tightly as possible in a single transaction and using callbacks to avoid
// waiting for them to commit, thus achieving good performance. This API hides away the logic of
// creating and committing transactions. Due to the nature of SSI guaratees provided by Badger,
// blind writes can never encounter transaction conflicts (ErrConflict).
func (db *DB) NewWriteBatch() *WriteBatch {
	return &WriteBatch{db: db, txn: db.newTransaction(true, true)}
}

// Cancel function must be called if there's a chance that Flush might not get
// called. If neither Flush or Cancel is called, the transaction oracle would
// never get a chance to clear out the row commit timestamp map, thus causing an
// unbounded memory consumption. Typically, you can call Cancel as a defer
// statement right after NewWriteBatch is called.
//
// Note that any committed writes would still go through despite calling Cancel.
func (wb *WriteBatch) Cancel() {
	wb.wg.Wait()
	wb.txn.Discard()
}

func (wb *WriteBatch) callback(err error) {
	// sync.WaitGroup is thread-safe, so it doesn't need to be run inside wb.Lock.
	defer wb.wg.Done()
	if err == nil {
		return
	}

	wb.Lock()
	defer wb.Unlock()
	if wb.err != nil {
		return
	}
	wb.err = err
}

// Set is equivalent of Txn.SetWithMeta.
func (wb *WriteBatch) Set(k, v []byte, meta byte) error {
	wb.Lock()
	defer wb.Unlock()

	if err := wb.txn.SetWithMeta(k, v, meta); err != ErrTxnTooBig {
		return err
	}
	// Txn has reached it's zenith. Commit now.
	if cerr := wb.commit(); cerr != nil {
		return cerr
	}
	// This time the error must not be ErrTxnTooBig, otherwise, we make the
	// error permanent.
	if err := wb.txn.SetWithMeta(k, v, meta); err != nil {
		wb.err = err
		return err
	}
	return nil
}

// Delete is equivalent of Txn.Delete.
func (wb *WriteBatch) Delete(k []byte) error {
	wb.Lock()
	defer wb.Unlock()

	if err := wb.txn.Delete(k); err != ErrTxnTooBig {
		return err
	}
	if err := wb.commit(); err != nil {
		return err
	}
	if err := wb.txn.Delete(k); err != nil {
		wb.err = err
		return err
	}
	return nil
}

// Caller to commit must hold a write lock.
func (wb *WriteBatch) commit() error {
	if wb.err != nil {
		return wb.err
	}
	// Get a new txn before we commit this one. So, the new txn doesn't need
	// to wait for this one to commit.
	wb.wg.Add(1)
	wb.txn.CommitWith(wb.callback)
	wb.txn = wb.db.newTransaction(true, true)
	wb.txn.readTs = 0 // We're not reading anything.
	return wb.err
}

// Flush must be called at the end to ensure that any pending writes get committed to Badger. Flush
// returns any error stored by WriteBatch.
func (wb *WriteBatch) Flush() error {
	wb.Lock()
	_ = wb.commit()
	wb.txn.Discard()
	wb.Unlock()

	wb.wg.Wait()
	// Safe to access error without any synchronization here.
	return wb.err
}

// Error returns any errors encountered so far. No commits would be run once an error is detected.
func (wb *WriteBatch) Error() error {
	wb.Lock()
	defer wb.Unlock()
	return wb.err
}

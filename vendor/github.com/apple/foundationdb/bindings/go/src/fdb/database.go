/*
 * database.go
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

import (
	"errors"
	"runtime"
)

// Database is a handle to a FoundationDB database. Database is a lightweight
// object that may be efficiently copied, and is safe for concurrent use by
// multiple goroutines.
//
// Although Database provides convenience methods for reading and writing data,
// modifications to a database are usually made via transactions, which are
// usually created and committed automatically by the (Database).Transact
// method.
type Database struct {
	// String reference to the cluster file.
	clusterFile string
	// This variable is to track if we have to remove the database from the cached
	// database structs. We can't use clusterFile alone, since the default clusterFile
	// would be an empty string.
	isCached bool
	*database
}

type database struct {
	ptr *C.FDBDatabase
}

// DatabaseOptions is a handle with which to set options that affect a Database
// object. A DatabaseOptions instance should be obtained with the
// (Database).Options method.
type DatabaseOptions struct {
	d *database
}

// Close will close the Database and clean up all resources.
// It must be called exactly once for each created database.
// You have to ensure that you're not reusing this database
// after it has been closed.
func (d Database) Close() {
	// Remove database object from the cached databases
	if d.isCached {
		openDatabases.Delete(d.clusterFile)
	}

	if d.ptr == nil {
		return
	}

	// Destroy the database
	C.fdb_database_destroy(d.ptr)
}

func (opt DatabaseOptions) setOpt(code int, param []byte) error {
	return setOpt(func(p *C.uint8_t, pl C.int) C.fdb_error_t {
		return C.fdb_database_set_option(opt.d.ptr, C.FDBDatabaseOption(code), p, pl)
	}, param)
}

// CreateTransaction returns a new FoundationDB transaction. It is generally
// preferable to use the (Database).Transact method, which handles
// automatically creating and committing a transaction with appropriate retry
// behavior.
func (d Database) CreateTransaction() (Transaction, error) {
	var outt *C.FDBTransaction

	if err := C.fdb_database_create_transaction(d.ptr, &outt); err != 0 {
		return Transaction{}, Error{int(err)}
	}

	t := &transaction{outt, d}
	// transactions cannot be destroyed explicitly if any future is still potentially used
	// thus the GC is used to figure out when all Go wrapper objects for futures have gone out of scope,
	// making the transaction ready to be garbage-collected.
	runtime.SetFinalizer(t, (*transaction).destroy)

	return Transaction{t}, nil
}

// RebootWorker is a wrapper around fdb_database_reboot_worker and allows to reboot processes
// from the go bindings. If a suspendDuration > 0 is provided the rebooted process will be
// suspended for suspendDuration seconds. If checkFile is set to true the process will check
// if the data directory is writeable by creating a validation file. The address must be a
// process address is the form of IP:Port pair.
func (d Database) RebootWorker(address string, checkFile bool, suspendDuration int) error {
	t := &futureInt64{
		future: newFutureWithDb(
			d.database,
			nil,
			C.fdb_database_reboot_worker(
				d.ptr,
				byteSliceToPtr([]byte(address)),
				C.int(len(address)),
				C.fdb_bool_t(boolToInt(checkFile)),
				C.int(suspendDuration),
			),
		),
	}

	dbVersion, err := t.Get()

	if dbVersion == 0 {
		return errors.New("failed to send reboot process request")
	}

	return err
}

func retryable(wrapped func() (interface{}, error), onError func(Error) FutureNil) (ret interface{}, e error) {
	for {
		ret, e = wrapped()

		// No error means success!
		if e == nil {
			return
		}

		// Check if the error chain contains an
		// fdb.Error
		var ep Error
		if errors.As(e, &ep) {
			e = onError(ep).Get()
		}

		// If OnError returns an error, then it's not
		// retryable; otherwise take another pass at things
		if e != nil {
			return
		}
	}
}

// Transact runs a caller-provided function inside a retry loop, providing it
// with a newly created Transaction. After the function returns, the Transaction
// will be committed automatically. Any error during execution of the function
// (by panic or return) or the commit will cause the function and commit to be
// retried or, if fatal, return the error to the caller.
//
// When working with Future objects in a transactional function, you may either
// explicitly check and return error values using Get, or call MustGet. Transact
// will recover a panicked Error and either retry the transaction or return the
// error.
//
// The transaction is retried if the error is or wraps a retryable Error.
// The error is unwrapped.
//
// Do not return Future objects from the function provided to Transact. The
// Transaction created by Transact may be finalized at any point after Transact
// returns, resulting in the cancellation of any outstanding
// reads. Additionally, any errors returned or panicked by the Future will no
// longer be able to trigger a retry of the caller-provided function.
//
// See the Transactor interface for an example of using Transact with
// Transaction and Database objects.
func (d Database) Transact(f func(Transaction) (interface{}, error)) (interface{}, error) {
	tr, e := d.CreateTransaction()
	// Any error here is non-retryable
	if e != nil {
		return nil, e
	}

	wrapped := func() (ret interface{}, e error) {
		defer panicToError(&e)

		ret, e = f(tr)

		if e == nil {
			e = tr.Commit().Get()
		}

		return
	}

	return retryable(wrapped, tr.OnError)
}

// ReadTransact runs a caller-provided function inside a retry loop, providing
// it with a newly created Transaction (as a ReadTransaction). Any error during
// execution of the function (by panic or return) will cause the function to be
// retried or, if fatal, return the error to the caller.
//
// When working with Future objects in a read-only transactional function, you
// may either explicitly check and return error values using Get, or call
// MustGet. ReadTransact will recover a panicked Error and either retry the
// transaction or return the error.
//
// The transaction is retried if the error is or wraps a retryable Error.
// The error is unwrapped.
// Read transactions are never committed and destroyed before returning to caller.
//
// Do not return Future objects from the function provided to ReadTransact. The
// Transaction created by ReadTransact may be finalized at any point after
// ReadTransact returns, resulting in the cancellation of any outstanding
// reads. Additionally, any errors returned or panicked by the Future will no
// longer be able to trigger a retry of the caller-provided function.
//
// See the ReadTransactor interface for an example of using ReadTransact with
// Transaction, Snapshot and Database objects.
func (d Database) ReadTransact(f func(ReadTransaction) (interface{}, error)) (interface{}, error) {
	tr, e := d.CreateTransaction()
	if e != nil {
		// Any error here is non-retryable
		return nil, e
	}

	wrapped := func() (ret interface{}, e error) {
		defer panicToError(&e)

		ret, e = f(tr)

		// read-only transactions are not committed and will be destroyed automatically via GC,
		// once all the futures go out of scope

		return
	}

	return retryable(wrapped, tr.OnError)
}

// Options returns a DatabaseOptions instance suitable for setting options
// specific to this database.
func (d Database) Options() DatabaseOptions {
	return DatabaseOptions{d.database}
}

// LocalityGetBoundaryKeys returns a slice of keys that fall within the provided
// range. Each key is located at the start of a contiguous range stored on a
// single server.
//
// If limit is non-zero, only the first limit keys will be returned. In large
// databases, the number of boundary keys may be large. In these cases, a
// non-zero limit should be used, along with multiple calls to
// LocalityGetBoundaryKeys.
//
// If readVersion is non-zero, the boundary keys as of readVersion will be
// returned.
func (d Database) LocalityGetBoundaryKeys(er ExactRange, limit int, readVersion int64) ([]Key, error) {
	tr, e := d.CreateTransaction()
	if e != nil {
		return nil, e
	}

	if readVersion != 0 {
		tr.SetReadVersion(readVersion)
	}

	tr.Options().SetReadSystemKeys()
	tr.Options().SetLockAware()

	bk, ek := er.FDBRangeKeys()
	ffer := KeyRange{
		append(Key("\xFF/keyServers/"), bk.FDBKey()...),
		append(Key("\xFF/keyServers/"), ek.FDBKey()...),
	}

	kvs, e := tr.Snapshot().GetRange(ffer, RangeOptions{Limit: limit}).GetSliceWithError()
	if e != nil {
		return nil, e
	}

	size := len(kvs)
	if limit != 0 && limit < size {
		size = limit
	}

	boundaries := make([]Key, size)

	for i := 0; i < size; i++ {
		boundaries[i] = kvs[i].Key[13:]
	}

	return boundaries, nil
}

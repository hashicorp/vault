/*
 * fdb.go
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
// #include <stdlib.h>
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var (
	// ErrNetworkAlreadyStopped for multiple calls to StopNetwork().
	ErrNetworkAlreadyStopped = errors.New("network has already been stopped")

	// ErrNetworkIsStopped is returned when attempting to execute a function which needs to interact
	// with the network thread while the network thread is no more running.
	ErrNetworkIsStopped = errors.New("network is stopped")

	// ErrNetworkAlreadyStopped for a too early call to StopNetwork().
	ErrNetworkNotStarted = errors.New("network has not been started")
)

// Would put this in futures.go but for the documented issue with
// exports and functions in preamble
// (https://code.google.com/p/go-wiki/wiki/cgo#Global_functions)
//
//export unlockMutex
func unlockMutex(p unsafe.Pointer) {
	m := (*sync.Mutex)(p)
	m.Unlock()
}

// A Transactor can execute a function that requires a Transaction. Functions
// written to accept a Transactor are called transactional functions, and may be
// called with either a Database or a Transaction.
type Transactor interface {
	// Transact executes the caller-provided function, providing it with a
	// Transaction (itself a Transactor, allowing composition of transactional
	// functions).
	Transact(func(Transaction) (interface{}, error)) (interface{}, error)

	// All Transactors are also ReadTransactors, allowing them to be used with
	// read-only transactional functions.
	ReadTransactor
}

// A ReadTransactor can execute a function that requires a
// ReadTransaction. Functions written to accept a ReadTransactor are called
// read-only transactional functions, and may be called with a Database,
// Transaction or Snapshot.
type ReadTransactor interface {
	// ReadTransact executes the caller-provided function, providing it with a
	// ReadTransaction (itself a ReadTransactor, allowing composition of
	// read-only transactional functions).
	ReadTransact(func(ReadTransaction) (interface{}, error)) (interface{}, error)
}

func setOpt(setter func(*C.uint8_t, C.int) C.fdb_error_t, param []byte) error {
	if err := setter(byteSliceToPtr(param), C.int(len(param))); err != 0 {
		return Error{int(err)}
	}

	return nil
}

// NetworkOptions is a handle with which to set options that affect the entire
// FoundationDB client. A NetworkOptions instance should be obtained with the
// fdb.Options function.
type NetworkOptions struct {
}

// Options returns a NetworkOptions instance suitable for setting options that
// affect the entire FoundationDB client.
func Options() NetworkOptions {
	return NetworkOptions{}
}

func (opt NetworkOptions) setOpt(code int, param []byte) error {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion == 0 {
		return errAPIVersionUnset
	}

	return setOpt(func(p *C.uint8_t, pl C.int) C.fdb_error_t {
		return C.fdb_network_set_option(C.FDBNetworkOption(code), p, pl)
	}, param)
}

// APIVersion determines the runtime behavior the fdb package. If the requested
// version is not supported by both the fdb package and the FoundationDB C
// library, an error will be returned. APIVersion must be called prior to any
// other functions in the fdb package.
// This function is safe to be called from multiple goroutines.
//
// Currently, this package supports API versions 200 through 740.
//
// Warning: When using the multi-version client API, setting an API version that
// is not supported by a particular client library will prevent that client from
// being used to connect to the cluster. In particular, you should not advance
// the API version of your application after upgrading your client until the
// cluster has also been upgraded.
func APIVersion(version int) error {
	headerVersion := 740

	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion != 0 {
		if apiVersion == version {
			return nil
		}
		return errAPIVersionAlreadySet
	}

	if version < 200 || version > headerVersion {
		return errAPIVersionNotSupported
	}

	if e := C.fdb_select_api_version_impl(C.int(version), C.int(headerVersion)); e != 0 {
		if e != 0 {
			if e == 2203 {
				maxSupportedVersion := C.fdb_get_max_api_version()
				if headerVersion > int(maxSupportedVersion) {
					return fmt.Errorf("This version of the FoundationDB Go binding is "+
						"not supported by the installed FoundationDB C library. "+
						"The binding requires a library that supports API version %d, "+
						"but the installed library supports a maximum version of %d.", headerVersion, maxSupportedVersion)
				}
				return fmt.Errorf("API version %d is not supported by the installed FoundationDB C library.", version)
			}
			return Error{int(e)}
		}
	}

	apiVersion = version

	return nil
}

// Determines if an API version has already been selected, i.e., if
// APIVersion or MustAPIVersion have already been called.
func IsAPIVersionSelected() bool {
	return apiVersion != 0
}

// Returns the API version that has been selected through APIVersion
// or MustAPIVersion. If the version has already been selected, then
// the first value returned is the API version and the error is
// nil. If the API version has not yet been set, then the error
// will be non-nil.
func GetAPIVersion() (int, error) {
	if IsAPIVersionSelected() {
		return apiVersion, nil
	}
	return 0, errAPIVersionUnset
}

// MustAPIVersion is like APIVersion but panics if the API version is not
// supported.
func MustAPIVersion(version int) {
	err := APIVersion(version)
	if err != nil {
		panic(err)
	}
}

// MustGetAPIVersion is like GetAPIVersion but panics if the API version
// has not yet been set.
func MustGetAPIVersion() int {
	apiVersion, err := GetAPIVersion()
	if err != nil {
		panic(err)
	}
	return apiVersion
}

var apiVersion int
var networkStarted, networkStopped bool
var networkMutex sync.RWMutex
var networkRunning sync.WaitGroup

var openDatabases sync.Map

// executeWithRunningNetworkThread starts the internal network event loop, if not already done,
// then runs the provided function while network thread is running.
// This function is safe to be called from multiple goroutines.
func executeWithRunningNetworkThread(f func()) error {
	networkMutex.RLock()
	if networkStopped {
		networkMutex.RUnlock()

		return ErrNetworkIsStopped
	}

	if networkStarted {

		// network thread is guaranteed to be running while this user-provided function runs
		f()

		networkMutex.RUnlock()
		return nil
	}
	// release read lock and acquire write lock
	networkMutex.RUnlock()
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if networkStopped {
		return ErrNetworkIsStopped
	}

	// check if meanwhile another goroutine started the network thread
	if !networkStarted {
		if e := C.fdb_setup_network(); e != 0 {
			return Error{int(e)}
		}

		networkRunning.Add(1)
		go func() {
			e := C.fdb_run_network()
			networkRunning.Done()
			if e != 0 {
				panic(fmt.Sprintf("Unhandled error in FoundationDB network thread: %v (%v)\n", C.GoString(C.fdb_get_error(e)), e))
			}
		}()

		networkStarted = true
	}

	// network thread is guaranteed to be running while this user-provided function runs
	f()

	return nil
}

// Deprecated: the network is started automatically when a database is opened.
// StartNetwork does nothing, but it will ensure that the API version is set and return an error otherwise.
func StartNetwork() error {
	if apiVersion == 0 {
		return errAPIVersionUnset
	}

	return nil
}

// StopNetwork signals the internal network event loop to terminate and waits for its termination.
// This function is safe to be called from multiple goroutines.
// This function returns an error if network has not yet started or if network has already been stopped.
// See also: https://github.com/apple/foundationdb/issues/3015
func StopNetwork() error {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if !networkStarted {
		return ErrNetworkNotStarted
	}

	if networkStopped {
		return ErrNetworkAlreadyStopped
	}

	C.fdb_stop_network()
	networkRunning.Wait()

	networkStopped = true

	return nil
}

// DefaultClusterFile should be passed to fdb.Open to allow the FoundationDB C
// library to select the platform-appropriate default cluster file on the current machine.
const DefaultClusterFile string = ""

// OpenDefault returns a database handle to the FoundationDB cluster identified
// by the DefaultClusterFile on the current machine.
//
// A single client can use this function multiple times to connect to different
// clusters simultaneously, with each invocation requiring its own cluster file.
// To connect to multiple clusters running at different, incompatible versions,
// the multi-version client API must be used.
func OpenDefault() (Database, error) {
	return OpenDatabase(DefaultClusterFile)
}

// MustOpenDefault is like OpenDefault but panics if the default database cannot
// be opened.
func MustOpenDefault() Database {
	db, err := OpenDefault()
	if err != nil {
		panic(err)
	}
	return db
}

// Open returns a database handle to the FoundationDB cluster identified
// by the provided cluster file and database name.
//
// A single client can use this function multiple times to connect to different
// clusters simultaneously, with each invocation requiring its own cluster file.
// To connect to multiple clusters running at different, incompatible versions,
// the multi-version client API must be used.
// Caller must call Close() to release resources.
func OpenDatabase(clusterFile string) (Database, error) {
	var db Database
	var okDb bool
	anyy, exist := openDatabases.Load(clusterFile)
	if db, okDb = anyy.(Database); !exist || !okDb {
		var e error
		db, e = createDatabase(clusterFile)
		if e != nil {
			return Database{}, e
		}
		openDatabases.Store(clusterFile, db)
	}

	return db, nil
}

// MustOpenDatabase is like OpenDatabase but panics if the default database cannot
// be opened.
func MustOpenDatabase(clusterFile string) Database {
	db, err := OpenDatabase(clusterFile)
	if err != nil {
		panic(err)
	}
	return db
}

// Deprecated: Use OpenDatabase instead.
// The database name must be []byte("DB").
func Open(clusterFile string, dbName []byte) (Database, error) {
	if bytes.Compare(dbName, []byte("DB")) != 0 {
		return Database{}, Error{2013} // invalid_database_name
	}
	return OpenDatabase(clusterFile)
}

// Deprecated: Use MustOpenDatabase instead.
// MustOpen is like Open but panics if the database cannot be opened.
func MustOpen(clusterFile string, dbName []byte) Database {
	db, err := Open(clusterFile, dbName)
	if err != nil {
		panic(err)
	}
	return db
}

// createDatabase is the internal function used to create a database.
// Caller must call Close() to release resources.
func createDatabase(clusterFile string) (Database, error) {
	if apiVersion == 0 {
		return Database{}, errAPIVersionUnset
	}

	var cf *C.char

	if len(clusterFile) != 0 {
		cf = C.CString(clusterFile)
		defer C.free(unsafe.Pointer(cf))
	}

	var outdb *C.FDBDatabase
	var createErr error
	if err := executeWithRunningNetworkThread(func() {
		if err := C.fdb_create_database(cf, &outdb); err != 0 {
			createErr = Error{int(err)}
		}
	}); err != nil {
		return Database{}, err
	}
	if createErr != nil {
		return Database{}, createErr
	}

	db := &database{outdb}

	return Database{clusterFile: clusterFile, isCached: true, database: db}, nil
}

// OpenWithConnectionString returns a database handle to the FoundationDB cluster identified
// by the provided connection string. This method can be useful for scenarios where you want to connect
// to the database only for a short time e.g. to test different connection strings.
// Caller must call Close() to release resources.
func OpenWithConnectionString(connectionString string) (Database, error) {
	if apiVersion == 0 {
		return Database{}, errAPIVersionUnset
	}

	var cf *C.char

	if connectionString == "" {
		return Database{}, errors.New("connection string must be a non-empty string")
	}

	cf = C.CString(connectionString)
	defer C.free(unsafe.Pointer(cf))

	var outdb *C.FDBDatabase
	var createErr error
	if err := executeWithRunningNetworkThread(func() {
		if err := C.fdb_create_database_from_connection_string(cf, &outdb); err != 0 {
			createErr = Error{int(err)}
		}
	}); err != nil {
		return Database{}, err
	}
	if createErr != nil {
		return Database{}, createErr
	}

	db := &database{outdb}

	return Database{"", false, db}, nil
}

// Deprecated: Use OpenDatabase instead.
// CreateCluster returns a cluster handle to the FoundationDB cluster identified
// by the provided cluster file.
// This function is safe to be called from multiple goroutines.
func CreateCluster(clusterFile string) (Cluster, error) {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion == 0 {
		return Cluster{}, errAPIVersionUnset
	}

	if !networkStarted {
		return Cluster{}, errNetworkNotSetup
	}

	if networkStopped {
		return Cluster{}, ErrNetworkIsStopped
	}

	return Cluster{clusterFile}, nil
}

func byteSliceToPtr(b []byte) *C.uint8_t {
	if len(b) > 0 {
		return (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nil
}

// A KeyConvertible can be converted to a FoundationDB Key. All functions in the
// FoundationDB API that address a specific key accept a KeyConvertible.
type KeyConvertible interface {
	FDBKey() Key
}

// Key represents a FoundationDB key, a lexicographically-ordered sequence of
// bytes. Key implements the KeyConvertible interface.
type Key []byte

// FDBKey allows Key to (trivially) satisfy the KeyConvertible interface.
func (k Key) FDBKey() Key {
	return k
}

// String describes the key as a human readable string.
func (k Key) String() string {
	return Printable(k)
}

// Printable returns a human readable version of a byte array. The bytes that correspond with
// ASCII printable characters [32-127) are passed through. Other bytes are
// replaced with \x followed by a two character zero-padded hex code for byte.
func Printable(d []byte) string {
	buf := new(bytes.Buffer)
	for _, b := range d {
		if b >= 32 && b < 127 && b != '\\' {
			buf.WriteByte(b)
			continue
		}
		if b == '\\' {
			buf.WriteString("\\\\")
			continue
		}
		buf.WriteString(fmt.Sprintf("\\x%02x", b))
	}
	return buf.String()
}

func panicToError(e *error) {
	if r := recover(); r != nil {
		fe, ok := r.(Error)
		if ok {
			*e = fe
		} else {
			panic(r)
		}
	}
}

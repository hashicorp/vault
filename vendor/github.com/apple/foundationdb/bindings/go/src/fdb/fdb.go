/*
 * fdb.go
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

// FoundationDB Go API

package fdb

/*
 #define FDB_API_VERSION 600
 #include <foundationdb/fdb_c.h>
 #include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"unsafe"
)

/* Would put this in futures.go but for the documented issue with
/* exports and functions in preamble
/* (https://code.google.com/p/go-wiki/wiki/cgo#Global_functions) */
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
//
// Currently, this package supports API versions 200 through 600.
//
// Warning: When using the multi-version client API, setting an API version that
// is not supported by a particular client library will prevent that client from
// being used to connect to the cluster. In particular, you should not advance
// the API version of your application after upgrading your client until the
// cluster has also been upgraded.
func APIVersion(version int) error {
	headerVersion := 600

	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion != 0 {
		if apiVersion == version {
			return nil
		}
		return errAPIVersionAlreadySet
	}

	if version < 200 || version > 600 {
		return errAPIVersionNotSupported
	}

	if e := C.fdb_select_api_version_impl(C.int(version), C.int(headerVersion)); e != 0 {
		if e != 0 {
			if e == 2203 {
				maxSupportedVersion := C.fdb_get_max_api_version()
				if headerVersion > int(maxSupportedVersion) {
					return fmt.Errorf("This version of the FoundationDB Go binding is not supported by the installed FoundationDB C library. The binding requires a library that supports API version %d, but the installed library supports a maximum version of %d.", version, maxSupportedVersion)
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
var networkStarted bool
var networkMutex sync.Mutex

type DatabaseId struct {
	clusterFile string
	dbName      string
}

var openClusters map[string]Cluster
var openDatabases map[DatabaseId]Database

func init() {
	openClusters = make(map[string]Cluster)
	openDatabases = make(map[DatabaseId]Database)
}

func startNetwork() error {
	if e := C.fdb_setup_network(); e != 0 {
		return Error{int(e)}
	}

	go func() {
		e := C.fdb_run_network()
		if e != 0 {
			log.Printf("Unhandled error in FoundationDB network thread: %v (%v)\n", C.GoString(C.fdb_get_error(e)), e)
		}
	}()

	networkStarted = true

	return nil
}

// StartNetwork initializes the FoundationDB client networking engine. It is not
// necessary to call StartNetwork when using the fdb.Open or fdb.OpenDefault
// functions to obtain a database handle. StartNetwork must not be called more
// than once.
func StartNetwork() error {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion == 0 {
		return errAPIVersionUnset
	}

	return startNetwork()
}

// DefaultClusterFile should be passed to fdb.Open or fdb.CreateCluster to allow
// the FoundationDB C library to select the platform-appropriate default cluster
// file on the current machine.
const DefaultClusterFile string = ""

// OpenDefault returns a database handle to the default database from the
// FoundationDB cluster identified by the DefaultClusterFile on the current
// machine. The FoundationDB client networking engine will be initialized first,
// if necessary.
func OpenDefault() (Database, error) {
	return Open(DefaultClusterFile, []byte("DB"))
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

// Open returns a database handle to the named database from the FoundationDB
// cluster identified by the provided cluster file and database name. The
// FoundationDB client networking engine will be initialized first, if
// necessary.
//
// In the current release, the database name must be []byte("DB").
func Open(clusterFile string, dbName []byte) (Database, error) {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion == 0 {
		return Database{}, errAPIVersionUnset
	}

	var e error

	if !networkStarted {
		e = startNetwork()
		if e != nil {
			return Database{}, e
		}
	}

	cluster, ok := openClusters[clusterFile]
	if !ok {
		cluster, e = createCluster(clusterFile)
		if e != nil {
			return Database{}, e
		}
		openClusters[clusterFile] = cluster
	}

	db, ok := openDatabases[DatabaseId{clusterFile, string(dbName)}]
	if !ok {
		db, e = cluster.OpenDatabase(dbName)
		if e != nil {
			return Database{}, e
		}
		openDatabases[DatabaseId{clusterFile, string(dbName)}] = db
	}

	return db, nil
}

// MustOpen is like Open but panics if the database cannot be opened.
func MustOpen(clusterFile string, dbName []byte) Database {
	db, err := Open(clusterFile, dbName)
	if err != nil {
		panic(err)
	}
	return db
}

func createCluster(clusterFile string) (Cluster, error) {
	var cf *C.char

	if len(clusterFile) != 0 {
		cf = C.CString(clusterFile)
		defer C.free(unsafe.Pointer(cf))
	}

	f := C.fdb_create_cluster(cf)
	fdb_future_block_until_ready(f)

	var outc *C.FDBCluster

	if err := C.fdb_future_get_cluster(f, &outc); err != 0 {
		return Cluster{}, Error{int(err)}
	}

	C.fdb_future_destroy(f)

	c := &cluster{outc}
	runtime.SetFinalizer(c, (*cluster).destroy)

	return Cluster{c}, nil
}

// CreateCluster returns a cluster handle to the FoundationDB cluster identified
// by the provided cluster file.
func CreateCluster(clusterFile string) (Cluster, error) {
	networkMutex.Lock()
	defer networkMutex.Unlock()

	if apiVersion == 0 {
		return Cluster{}, errAPIVersionUnset
	}

	if !networkStarted {
		return Cluster{}, errNetworkNotSetup
	}

	return createCluster(clusterFile)
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

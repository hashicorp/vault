/*
 * cluster.go
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
*/
import "C"

import (
	"runtime"
)

// Cluster is a handle to a FoundationDB cluster. Cluster is a lightweight
// object that may be efficiently copied, and is safe for concurrent use by
// multiple goroutines.
//
// It is generally preferable to use Open or OpenDefault to obtain a database
// handle directly.
type Cluster struct {
	*cluster
}

type cluster struct {
	ptr *C.FDBCluster
}

func (c *cluster) destroy() {
	C.fdb_cluster_destroy(c.ptr)
}

// OpenDatabase returns a database handle from the FoundationDB cluster. It is
// generally preferable to use Open or OpenDefault to obtain a database handle
// directly.
//
// In the current release, the database name must be []byte("DB").
func (c Cluster) OpenDatabase(dbName []byte) (Database, error) {
	f := C.fdb_cluster_create_database(c.ptr, byteSliceToPtr(dbName), C.int(len(dbName)))
	fdb_future_block_until_ready(f)

	var outd *C.FDBDatabase

	if err := C.fdb_future_get_database(f, &outd); err != 0 {
		return Database{}, Error{int(err)}
	}

	C.fdb_future_destroy(f)

	d := &database{outd}
	runtime.SetFinalizer(d, (*database).destroy)

	return Database{d}, nil
}

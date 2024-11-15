/*
 * errors.go
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
	"fmt"
)

// Error represents a low-level error returned by the FoundationDB C library. An
// Error may be returned by any FoundationDB API function that returns error, or
// as a panic from any FoundationDB API function whose name ends with OrPanic.
//
// You may compare the Code field of an Error against the list of FoundationDB
// error codes at https://apple.github.io/foundationdb/api-error-codes.html,
// but generally an Error should be passed to (Transaction).OnError. When using
// (Database).Transact, non-fatal errors will be retried automatically.
type Error struct {
	Code int
}

func (e Error) Error() string {
	return fmt.Sprintf("FoundationDB error code %d (%s)", e.Code, C.GoString(C.fdb_get_error(C.fdb_error_t(e.Code))))
}

// SOMEDAY: these (along with others) should be coming from fdb.options?

var (
	errNetworkNotSetup          = Error{2008}
	errNetworkAlreadySetup      = Error{2009} // currently unused
	errNetworkCannotBeRestarted = Error{2025} // currently unused

	errAPIVersionUnset        = Error{2200}
	errAPIVersionAlreadySet   = Error{2201}
	errAPIVersionNotSupported = Error{2203}
)

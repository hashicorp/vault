// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

// RecordExistsAction determines how to handle writes when
// the record already exists.
type RecordExistsAction int

const (

	// UPDATE means: Create or update record.
	// Merge write command bins with existing bins.
	UPDATE RecordExistsAction = iota

	// UPDATE_ONLY means: Update record only. Fail if record does not exist.
	// Merge write command bins with existing bins.
	UPDATE_ONLY

	// REPLACE means: Create or replace record.
	// Delete existing bins not referenced by write command bins.
	// Supported by Aerospike 2 server versions >= 2.7.5 and
	// Aerospike 3 server versions >= 3.1.6 and later.
	REPLACE

	// REPLACE_ONLY means: Replace record only. Fail if record does not exist.
	// Delete existing bins not referenced by write command bins.
	// Supported by Aerospike 2 server versions >= 2.7.5 and
	// Aerospike 3 server versions >= 3.1.6 and later.
	REPLACE_ONLY

	// CREATE_ONLY means: Create only. Fail if record exists.
	CREATE_ONLY
)

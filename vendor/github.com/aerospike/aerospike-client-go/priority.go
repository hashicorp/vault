// Copyright 2013-2020 Aerospike, Inc.
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

// Priority of operations on database server.
// Currently, this only affects Scan operations.
type Priority int

const (

	// DEFAULT determines that the server defines the priority.
	DEFAULT Priority = iota

	// LOW determines that the server should run the operation in a background thread.
	LOW

	// MEDIUM determines that the server should run the operation at medium priority.
	MEDIUM

	// HIGH determines that the server should run the operation at the highest priority.
	HIGH
)

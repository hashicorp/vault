// Copyright 2014-2021 Aerospike, Inc.
//
// Portions may be licensed to Aerospike, Inc. under one or more contributor
// license agreements.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package aerospike

// UserRoles contains information about a user.
type UserRoles struct {
	// User name.
	User string

	// Roles is a list of assigned roles.
	Roles []string

	// ReadInfo is the list of read statistics. List may be nil.
	// Current statistics by offset are:
	//
	// 0: read quota in records per second
	// 1: single record read transaction rate (TPS)
	// 2: read scan/query record per second rate (RPS)
	// 3: number of limitless read scans/queries
	//
	// Future server releases may add additional statistics.
	ReadInfo []int

	// WriteInfo is the list of write statistics. List may be nil.
	// Current statistics by offset are:
	//
	// 0: write quota in records per second
	// 1: single record write transaction rate (TPS)
	// 2: write scan/query record per second rate (RPS)
	// 3: number of limitless write scans/queries
	//
	// Future server releases may add additional statistics.
	WriteInfo []int

	// ConnsInUse is the number of currently open connections for the user
	ConnsInUse int
}

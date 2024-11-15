/*
 * Copyright 2014-2021 Aerospike, Inc.
 *
 * Portions may be licensed to Aerospike, Inc. under one or more contributor
 * license agreements.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package aerospike

// ReadModeSC is the read policy in SC (strong consistency) mode namespaces.
// Determines SC read consistency options.
type ReadModeSC int

const (
	// ReadModeSCSession ensures this client will only see an increasing sequence of record versions.
	// Server only reads from master.  This is the default.
	ReadModeSCSession ReadModeSC = iota

	// ReadModeSCLinearize ensures ALL clients will only see an increasing sequence of record versions.
	// Server only reads from master.
	ReadModeSCLinearize

	// ReadModeSCAllowReplica indicates that the server may read from master or any full (non-migrating) replica.
	// Increasing sequence of record versions is not guaranteed.
	ReadModeSCAllowReplica

	// ReadModeSCAllowUnavailable indicates that the server may read from master or any full (non-migrating) replica or from unavailable
	// partitions.  Increasing sequence of record versions is not guaranteed.
	ReadModeSCAllowUnavailable
)

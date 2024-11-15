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

// CommitLevel indicates the desired consistency guarantee when committing a transaction on the server.
type CommitLevel int

const (
	// COMMIT_ALL indicates the server should wait until successfully committing master and all replicas.
	COMMIT_ALL CommitLevel = iota

	// COMMIT_MASTER indicates the server should wait until successfully committing master only.
	COMMIT_MASTER
)

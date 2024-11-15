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

// ReplicaPolicy defines type of node partition targeted by read commands.
type ReplicaPolicy int

const (
	// MASTER reads from node containing key's master partition.
	// This is the default behavior.
	MASTER ReplicaPolicy = iota

	// MASTER_PROLES Distributes reads across nodes containing key's master and replicated partitions
	// in round-robin fashion.
	MASTER_PROLES

	// RANDOM Distribute reads across all nodes in cluster in round-robin fashion.
	// This option is useful when the replication factor equals the number
	// of nodes in the cluster and the overhead of requesting proles is not desired.
	RANDOM

	// SEQUENCE Tries node containing master partition first.
	// If connection fails, all commands try nodes containing replicated partitions.
	// If socketTimeout is reached, reads also try nodes containing replicated partitions,
	// but writes remain on master node.
	SEQUENCE

	// PREFER_RACK Tries nodes on the same rack first.
	//
	// This option requires ClientPolicy.Rackaware to be enabled
	// in order to function properly.
	PREFER_RACK
)

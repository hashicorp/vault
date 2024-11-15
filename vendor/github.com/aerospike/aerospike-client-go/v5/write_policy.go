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

import "math"

const (
	// TTLServerDefault will default to namespace configuration variable "default-ttl" on the server.
	TTLServerDefault = 0
	// TTLDontExpire will never expire for Aerospike 2 server versions >= 2.7.2 and Aerospike 3+ server.
	TTLDontExpire = math.MaxUint32
	// TTLDontUpdate will not change the record's ttl when record is written. Supported by Aerospike server versions >= 3.10.1
	TTLDontUpdate = math.MaxUint32 - 1
)

// WritePolicy encapsulates parameters for policy attributes used in write operations.
// This object is passed into methods where database writes can occur.
type WritePolicy struct {
	BasePolicy

	// RecordExistsAction qualifies how to handle writes where the record already exists.
	RecordExistsAction RecordExistsAction //= RecordExistsAction.UPDATE;

	// GenerationPolicy qualifies how to handle record writes based on record generation. The default (NONE)
	// indicates that the generation is not used to restrict writes.
	GenerationPolicy GenerationPolicy //= GenerationPolicy.NONE;

	// Desired consistency guarantee when committing a transaction on the server. The default
	// (COMMIT_ALL) indicates that the server should wait for master and all replica commits to
	// be successful before returning success to the client.
	CommitLevel CommitLevel //= COMMIT_ALL

	// Generation determines expected generation.
	// Generation is the number of times a record has been
	// modified (including creation) on the server.
	// If a write operation is creating a record, the expected generation would be 0.
	Generation uint32

	// Expiration determines record expiration in seconds. Also known as TTL (Time-To-Live).
	// Seconds record will live before being removed by the server.
	// Expiration values:
	// TTLServerDefault (0): Default to namespace configuration variable "default-ttl" on the server.
	// TTLDontExpire (MaxUint32): Never expire for Aerospike 2 server versions >= 2.7.2 and Aerospike 3+ server
	// TTLDontUpdate (MaxUint32 - 1): Do not change ttl when record is written. Supported by Aerospike server versions >= 3.10.1
	// > 0: Actual expiration in seconds.
	Expiration uint32

	// RespondPerEachOp defines for client.Operate() method, return a result for every operation.
	// Some list operations do not return results by default (ListClearOp() for example).
	// This can sometimes make it difficult to determine the desired result offset in the returned
	// bin's result list.
	//
	// Setting RespondPerEachOp to true makes it easier to identify the desired result offset
	// (result offset equals bin's operate sequence). This only makes sense when multiple list
	// operations are used in one operate call and some of those operations do not return results
	// by default.
	RespondPerEachOp bool

	// DurableDelete leaves a tombstone for the record if the transaction results in a record deletion.
	// This prevents deleted records from reappearing after node failures.
	// Valid for Aerospike Server Enterprise Edition 3.10+ only.
	DurableDelete bool
}

// NewWritePolicy initializes a new WritePolicy instance with default parameters.
func NewWritePolicy(generation, expiration uint32) *WritePolicy {
	res := &WritePolicy{
		BasePolicy:         *NewPolicy(),
		RecordExistsAction: UPDATE,
		GenerationPolicy:   NONE,
		CommitLevel:        COMMIT_ALL,
		Generation:         generation,
		Expiration:         expiration,
	}

	// Writes may not be idempotent.
	// do not allow retries on writes by default.
	res.MaxRetries = 0

	return res
}

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

import "fmt"

// Record is the container struct for database records.
// Records are equivalent to rows.
type Record struct {
	// Key is the record's key.
	// Might be empty, or may only consist of digest value.
	Key *Key

	// Node from which the Record is originating from.
	Node *Node

	// Bins is the map of requested name/value bins.
	Bins BinMap

	// Generation shows record modification count.
	Generation uint32

	// Expiration is TTL (Time-To-Live).
	// Number of seconds until record expires.
	Expiration uint32
}

func newRecord(node *Node, key *Key, bins BinMap, generation, expiration uint32) *Record {
	r := &Record{
		Node:       node,
		Key:        key,
		Bins:       bins,
		Generation: generation,
		Expiration: expiration,
	}

	// always assign a map of length zero if Bins is nil
	if r.Bins == nil {
		r.Bins = make(BinMap)
	}

	return r
}

// String implements the Stringer interface.
// Returns string representation of record.
func (rc *Record) String() string {
	return fmt.Sprintf("%s %v", rc.Key, rc.Bins)
}

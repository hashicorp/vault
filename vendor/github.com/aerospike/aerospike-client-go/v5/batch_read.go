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

import (
	"fmt"
)

// BatchRead specifies the Key and bin names used in batch read commands
// where variable bins are needed for each key.
type BatchRead struct {
	// Key specifies the key to retrieve.
	Key *Key

	// BinNames specifies the Bins to retrieve for this key.
	// BinNames are mutually exclusive with Ops.
	BinNames []string

	// ReadAllBins defines what data should be read from the record.
	// If true, ignore binNames and read all bins.
	// If false and binNames are set, read specified binNames.
	// If false and binNames are not set, read record header (generation, expiration) only.
	ReadAllBins bool //= false

	// Record result after batch command has completed.  Will be null if record was not found.
	Record *Record

	// Ops specifies the operations to perform for every key.
	// Ops are mutually exclusive with BinNames.
	// A binName can be emulated with `GetOp(binName)`
	// Supported by server v5.6.0+.
	Ops []*Operation
}

// NewBatchRead defines a key and bins to retrieve in a batch operation.
func NewBatchRead(key *Key, binNames []string) *BatchRead {
	res := &BatchRead{
		Key:      key,
		BinNames: binNames,
	}

	if len(binNames) == 0 {
		res.ReadAllBins = true
	}

	return res
}

// NewBatchReadOps defines a key and bins to retrieve in a batch operation, including expressions.
func NewBatchReadOps(key *Key, binNames []string, ops []*Operation) *BatchRead {
	res := &BatchRead{
		Key:      key,
		BinNames: binNames,
		Ops:      ops,
	}

	if len(binNames) == 0 {
		res.ReadAllBins = true
	}

	return res
}

// NewBatchReadHeader defines a key to retrieve the record headers only in a batch operation.
func NewBatchReadHeader(key *Key) *BatchRead {
	return &BatchRead{
		Key:         key,
		ReadAllBins: false,
	}
}

// String implements the Stringer interface.
func (br *BatchRead) String() string {
	return fmt.Sprintf("%s: %v", br.Key, br.BinNames)
}

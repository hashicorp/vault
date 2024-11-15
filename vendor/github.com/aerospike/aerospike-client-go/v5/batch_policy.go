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

// BatchPolicy encapsulates parameters for policy attributes used in write operations.
// This object is passed into methods where database writes can occur.
type BatchPolicy struct {
	BasePolicy

	// Maximum number of concurrent batch request goroutines to server nodes at any point in time.
	// If there are 16 node/namespace combinations requested and ConcurrentNodes is 8,
	// then batch requests will be made for 8 node/namespace combinations in concurrent goroutines.
	// When a request completes, a new request will be issued until all 16 goroutines are complete.
	//
	// Values:
	// 1: Issue batch requests sequentially.  This mode has a performance advantage for small
	// to medium sized batch sizes because requests can be issued in the main transaction goroutine.
	// This is the default.
	// 0: Issue all batch requests in concurrent goroutines.  This mode has a performance
	// advantage for extremely large batch sizes because each node can process the request
	// immediately.  The downside is extra goroutines will need to be created (or taken from
	// a goroutine pool).
	// > 0: Issue up to ConcurrentNodes batch requests in concurrent goroutines.  When a request
	// completes, a new request will be issued until all goroutines are complete.  This mode
	// prevents too many concurrent goroutines being created for large cluster implementations.
	// The downside is extra goroutines will still need to be created (or taken from a goroutine pool).
	ConcurrentNodes int // = 1

	// Allow batch to be processed immediately in the server's receiving thread when the server
	// deems it to be appropriate.  If false, the batch will always be processed in separate
	// transaction goroutines.  This field is only relevant for the new batch index protocol.
	//
	// For batch exists or batch reads of smaller sized records (<= 1K per record), inline
	// processing will be significantly faster on "in memory" namespaces.  The server disables
	// inline processing on disk based namespaces regardless of this policy field.
	//
	// Inline processing can introduce the possibility of unfairness because the server
	// can process the entire batch before moving onto the next command.
	AllowInline bool //= true

	// AllowPartialResults determines if the results for some nodes should be returned in case
	// some nodes encounter an error. The result for the unreceived records will be nil.
	// The returned records will be safe to use, since only fully received data will be parsed
	// and set.
	//
	// This flag is only supported for BatchGet and BatchGetHeader methods. BatchGetComplex always returns
	// partial results by design.
	AllowPartialResults bool //= false

	// Send set name field to server for every key in the batch for batch index protocol.
	// This is only necessary when authentication is enabled and security roles are defined
	// on a per set basis.
	SendSetName bool //= false
}

// NewBatchPolicy initializes a new BatchPolicy instance with default parameters.
func NewBatchPolicy() *BatchPolicy {
	return &BatchPolicy{
		BasePolicy:          *NewPolicy(),
		ConcurrentNodes:     1,
		AllowInline:         true,
		AllowPartialResults: false,
		SendSetName:         false,
	}
}

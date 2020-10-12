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

import (
	"fmt"

	. "github.com/aerospike/aerospike-client-go/types"
)

// NodeError is a type to encapsulate the node that the error occurred in.
type NodeError struct {
	error

	node *Node
}

func newNodeError(node *Node, err error) *NodeError {
	return &NodeError{
		error: err,
		node:  node,
	}
}

func newAerospikeNodeError(node *Node, code ResultCode, messages ...string) *NodeError {
	return &NodeError{
		error: NewAerospikeError(code, messages...),
		node:  node,
	}
}

// Node returns the node where the error occurred.
func (ne *NodeError) Node() *Node { return ne.node }

// Err returns the error
func (ne *NodeError) Err() error { return ne.error }

// Err returns the error
func (ne *NodeError) Error() string {
	return fmt.Sprintf("Node %s: %s", ne.node.String(), ne.error.Error())
}

func newInvalidNodeError(clusterSize int, partition *Partition) error {
	// important to check for clusterSize first, since partition may be nil sometimes
	if clusterSize == 0 {
		return NewAerospikeError(INVALID_NODE_ERROR, "Cluster is empty.")
	}
	return NewAerospikeError(INVALID_NODE_ERROR, "Node not found for partition "+partition.String()+" in partition table.")
}

// BatchError is a type to encapsulate the node that the error occurred in.
type BatchError struct {
	Errors map[*Node]error
}

func newBatchError() *BatchError {
	return &BatchError{
		Errors: make(map[*Node]error, 4),
	}
}

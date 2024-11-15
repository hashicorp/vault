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
	"sync/atomic"

	"github.com/aerospike/aerospike-client-go/v5/types"
)

// Partition encapsulates partition information.
type Partition struct {
	// Namespace of the partition
	Namespace string
	// PartitionId of the partition
	PartitionId int
	partitions  *Partitions
	replica     ReplicaPolicy
	prevNode    *Node
	sequence    int
	linearize   bool
}

// NewPartition returns a partition representation
func NewPartition(partitions *Partitions, key *Key, replica ReplicaPolicy, linearize bool) *Partition {
	return &Partition{
		partitions:  partitions,
		Namespace:   key.Namespace(),
		replica:     replica,
		linearize:   linearize,
		PartitionId: key.PartitionId(),
	}
}

// PartitionForWrite returns a partition for write purposes
func PartitionForWrite(cluster *Cluster, policy *BasePolicy, key *Key) (*Partition, Error) {
	// Must copy hashmap reference for copy on write semantics to work.
	pmap := cluster.getPartitions()
	partitions := pmap[key.namespace]

	if partitions == nil {
		return nil, newInvalidNamespaceError(key.namespace, len(pmap))
	}

	return NewPartition(partitions, key, policy.ReplicaPolicy, false), nil
}

// PartitionForRead returns a partition for read purposes
func PartitionForRead(cluster *Cluster, policy *BasePolicy, key *Key) (*Partition, Error) {
	// Must copy hashmap reference for copy on write semantics to work.
	pmap := cluster.getPartitions()
	partitions := pmap[key.namespace]

	if partitions == nil {
		return nil, newInvalidNamespaceError(key.namespace, len(pmap))
	}

	var replica ReplicaPolicy
	var linearize bool

	if partitions.SCMode {
		switch policy.ReadModeSC {
		case ReadModeSCSession:
			replica = MASTER
			linearize = false

		case ReadModeSCLinearize:
			replica = policy.ReplicaPolicy
			if policy.ReplicaPolicy == PREFER_RACK {
				replica = SEQUENCE
			}
			linearize = true

		default:
			replica = policy.ReplicaPolicy
			linearize = false
		}
	} else {
		replica = policy.ReplicaPolicy
		linearize = false
	}
	return NewPartition(partitions, key, replica, linearize), nil
}

// GetReplicaPolicySC returns a ReplicaPolicy based on different variables in SC mode
func GetReplicaPolicySC(policy *BasePolicy) ReplicaPolicy {
	switch policy.ReadModeSC {
	case ReadModeSCSession:
		return MASTER

	case ReadModeSCLinearize:
		if policy.ReplicaPolicy == PREFER_RACK {
			return SEQUENCE
		}
		return policy.ReplicaPolicy

	default:
		return policy.ReplicaPolicy
	}
}

// GetNodeBatchRead returns a node for batch reads
func GetNodeBatchRead(cluster *Cluster, key *Key, replica ReplicaPolicy, replicaSC ReplicaPolicy, sequence int, sequenceSC int) (*Node, Error) {
	// Must copy hashmap reference for copy on write semantics to work.
	pmap := cluster.getPartitions()
	partitions := pmap[key.namespace]

	if partitions == nil {
		return nil, newInvalidNamespaceError(key.namespace, len(pmap))
	}

	if partitions.SCMode {
		replica = replicaSC
		sequence = sequenceSC
	}

	p := NewPartition(partitions, key, replica, false)
	p.sequence = sequence
	return p.GetNodeRead(cluster)
}

// GetNodeRead returns a node for read operations
func (ptn *Partition) GetNodeRead(cluster *Cluster) (*Node, Error) {
	switch ptn.replica {
	default:
		fallthrough
	case SEQUENCE:
		return ptn.getSequenceNode(cluster)

	case PREFER_RACK:
		return ptn.getRackNode(cluster)

	case MASTER:
		return ptn.getMasterNode(cluster)

	case MASTER_PROLES:
		return ptn.getMasterProlesNode(cluster)

	case RANDOM:
		return cluster.GetRandomNode()
	}
}

// GetNodeWrite returns a node for write operations
func (ptn *Partition) GetNodeWrite(cluster *Cluster) (*Node, Error) {
	switch ptn.replica {
	default:
		fallthrough
	case SEQUENCE:
		fallthrough
	case PREFER_RACK:
		return ptn.getSequenceNode(cluster)

	case MASTER:
		fallthrough
	case MASTER_PROLES:
		fallthrough
	case RANDOM:
		return ptn.getMasterNode(cluster)
	}
}

// PrepareRetryRead increases sequence number before read retries
func (ptn *Partition) PrepareRetryRead(isClientTimeout bool) {
	if !isClientTimeout || !ptn.linearize {
		ptn.sequence++
	}
}

// PrepareRetryWrite increases sequence number before write retries
func (ptn *Partition) PrepareRetryWrite(isClientTimeout bool) {
	if !isClientTimeout {
		ptn.sequence++
	}
}

func (ptn *Partition) getSequenceNode(cluster *Cluster) (*Node, Error) {
	replicas := ptn.partitions.Replicas

	for range replicas {
		index := ptn.sequence % len(replicas)
		node := replicas[index][ptn.PartitionId]

		if node != nil && node.IsActive() {
			return node, nil
		}
		ptn.sequence++
	}
	nodeArray := cluster.GetNodes()
	return nil, newInvalidNodeError(len(nodeArray), ptn)
}

func (ptn *Partition) getRackNode(cluster *Cluster) (*Node, Error) {
	replicas := ptn.partitions.Replicas

	for _, rackId := range cluster.clientPolicy.RackIds {
		seq := ptn.sequence
		for range replicas {
			index := ptn.sequence % len(replicas)
			node := replicas[index][ptn.PartitionId]

			if node != nil && node != ptn.prevNode && node.hasRack(ptn.Namespace, rackId) && node.IsActive() {
				ptn.prevNode = node
				ptn.sequence = seq
				return node, nil
			}
			seq++
		}
	}

	for range replicas {
		index := ptn.sequence % len(replicas)
		node := replicas[index][ptn.PartitionId]

		if node != nil && node.IsActive() {
			ptn.prevNode = node
			return node, nil
		}
		ptn.sequence++
	}

	nodeArray := cluster.GetNodes()
	return nil, newInvalidNodeError(len(nodeArray), ptn)
}

func (ptn *Partition) getMasterNode(cluster *Cluster) (*Node, Error) {
	node := ptn.partitions.Replicas[0][ptn.PartitionId]

	if node != nil && node.IsActive() {
		return node, nil
	}
	nodeArray := cluster.GetNodes()
	return nil, newInvalidNodeError(len(nodeArray), ptn)
}

func (ptn *Partition) getMasterProlesNode(cluster *Cluster) (*Node, Error) {
	replicas := ptn.partitions.Replicas

	for range replicas {
		index := int(atomic.AddUint64(&cluster.replicaIndex, 1) % uint64(len(replicas)))
		node := replicas[index][ptn.PartitionId]

		if node != nil && node.IsActive() {
			return node, nil
		}
	}
	nodeArray := cluster.GetNodes()
	return nil, newInvalidNodeError(len(nodeArray), ptn)
}

// String implements the Stringer interface.
func (ptn *Partition) String() string {
	return fmt.Sprintf("%s:%d", ptn.Namespace, ptn.PartitionId)
}

// Equals checks equality of two partitions.
func (ptn *Partition) Equals(other *Partition) bool {
	return ptn.PartitionId == other.PartitionId && ptn.Namespace == other.Namespace
}

// newnewInvalidNamespaceError creates an AerospikeError with Resultcode INVALID_NAMESPACE
// and a corresponding message.
func newInvalidNamespaceError(ns string, mapSize int) Error {
	s := "Partition map empty"
	if mapSize != 0 {
		s = "Namespace not found in partition map: " + ns
	}
	return newError(types.INVALID_NAMESPACE, s)
}

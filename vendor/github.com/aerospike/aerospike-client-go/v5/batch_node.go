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

type batchNode struct {
	Node    *Node
	offsets []int
}

func newBatchNodeList(cluster *Cluster, policy *BatchPolicy, keys []*Key) ([]*batchNode, Error) {
	nodes := cluster.GetNodes()

	if len(nodes) == 0 {
		return nil, ErrClusterIsEmpty.err()
	}

	// Create initial key capacity for each node as average + 25%.
	keysPerNode := len(keys) / len(nodes)
	keysPerNode += keysPerNode / 2

	// The minimum key capacity is 10.
	if keysPerNode < 10 {
		keysPerNode = 10
	}

	replicaPolicy := policy.ReplicaPolicy
	replicaPolicySC := GetReplicaPolicySC(policy.GetBasePolicy())

	// Split keys by server node.
	batchNodes := make([]*batchNode, 0, len(nodes))

	for i := range keys {
		node, err := GetNodeBatchRead(cluster, keys[i], replicaPolicy, replicaPolicySC, 0, 0)
		if err != nil {
			return nil, err
		}

		if batchNode := findBatchNode(batchNodes, node); batchNode == nil {
			batchNodes = append(batchNodes, newBatchNode(node, keysPerNode, i))
		} else {
			batchNode.AddKey(i)
		}
	}
	return batchNodes, nil
}

func newBatchNodeListKeys(cluster *Cluster, policy *BatchPolicy, keys []*Key, sequenceAP, sequenceSC int, batchSeed *batchNode) ([]*batchNode, Error) {
	nodes := cluster.GetNodes()

	if len(nodes) == 0 {
		return nil, ErrClusterIsEmpty.err()
	}

	// Create initial key capacity for each node as average + 25%.
	keysPerNode := len(keys) / len(nodes)
	keysPerNode += keysPerNode / 2

	// The minimum key capacity is 10.
	if keysPerNode < 10 {
		keysPerNode = 10
	}

	replicaPolicy := policy.ReplicaPolicy
	replicaPolicySC := GetReplicaPolicySC(policy.GetBasePolicy())

	// Split keys by server node.
	batchNodes := make([]*batchNode, 0, len(nodes))

	for _, offset := range batchSeed.offsets {
		node, err := GetNodeBatchRead(cluster, keys[offset], replicaPolicy, replicaPolicySC, sequenceAP, sequenceSC)
		if err != nil {
			return nil, err
		}

		if batchNode := findBatchNode(batchNodes, node); batchNode == nil {
			batchNodes = append(batchNodes, newBatchNode(node, keysPerNode, offset))
		} else {
			batchNode.AddKey(offset)
		}
	}
	return batchNodes, nil
}

func newBatchNodeListRecords(cluster *Cluster, policy *BatchPolicy, records []*BatchRead, sequenceAP, sequenceSC int, batchSeed *batchNode) ([]*batchNode, Error) {
	nodes := cluster.GetNodes()

	if len(nodes) == 0 {
		return nil, ErrClusterIsEmpty.err()
	}

	// Create initial key capacity for each node as average + 25%.
	keysPerNode := len(batchSeed.offsets) / len(nodes)
	keysPerNode += keysPerNode / 2

	// The minimum key capacity is 10.
	if keysPerNode < 10 {
		keysPerNode = 10
	}

	replicaPolicy := policy.ReplicaPolicy
	replicaPolicySC := GetReplicaPolicySC(policy.GetBasePolicy())

	// Split keys by server node.
	batchNodes := make([]*batchNode, 0, len(nodes))

	for _, offset := range batchSeed.offsets {
		node, err := GetNodeBatchRead(cluster, records[offset].Key, replicaPolicy, replicaPolicySC, sequenceAP, sequenceSC)
		if err != nil {
			return nil, err
		}

		if batchNode := findBatchNode(batchNodes, node); batchNode == nil {
			batchNodes = append(batchNodes, newBatchNode(node, keysPerNode, offset))
		} else {
			batchNode.AddKey(offset)
		}
	}
	return batchNodes, nil
}

func newBatchIndexNodeList(cluster *Cluster, policy *BatchPolicy, records []*BatchRead) ([]*batchNode, Error) {
	nodes := cluster.GetNodes()

	if len(nodes) == 0 {
		return nil, ErrClusterIsEmpty.err()
	}

	// Create initial key capacity for each node as average + 25%.
	keysPerNode := len(records) / len(nodes)
	keysPerNode += keysPerNode / 2

	// The minimum key capacity is 10.
	if keysPerNode < 10 {
		keysPerNode = 10
	}

	replicaPolicy := policy.ReplicaPolicy
	replicaPolicySC := GetReplicaPolicySC(policy.GetBasePolicy())

	// Split keys by server node.
	batchNodes := make([]*batchNode, 0, len(nodes))

	for i := range records {
		node, err := GetNodeBatchRead(cluster, records[i].Key, replicaPolicy, replicaPolicySC, 0, 0)
		if err != nil {
			return nil, err
		}

		if batchNode := findBatchNode(batchNodes, node); batchNode == nil {
			batchNodes = append(batchNodes, newBatchNode(node, keysPerNode, i))
		} else {
			batchNode.AddKey(i)
		}
	}
	return batchNodes, nil
}

func newBatchNode(node *Node, capacity int, offset int) *batchNode {
	res := &batchNode{
		Node:    node,
		offsets: make([]int, 1, capacity),
	}

	res.offsets[0] = offset
	return res
}

func (bn *batchNode) AddKey(offset int) {
	bn.offsets = append(bn.offsets, offset)
}

func findBatchNode(nodes []*batchNode, node *Node) *batchNode {
	for i := range nodes {
		// Note: using pointer equality for performance.
		if nodes[i].Node == node {
			return nodes[i]
		}
	}
	return nil
}

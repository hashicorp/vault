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
	"bytes"
	"fmt"
	"strconv"

	"github.com/aerospike/aerospike-client-go/v5/types"
)

// Partitions represents a list of partitions
type Partitions struct {
	Replicas [][]*Node
	SCMode   bool
	regimes  []int
}

func newPartitions(partitionCount int, replicaCount int, cpMode bool) *Partitions {
	replicas := make([][]*Node, replicaCount)
	for i := range replicas {
		replicas[i] = make([]*Node, partitionCount)
	}

	return &Partitions{
		Replicas: replicas,
		SCMode:   cpMode,
		regimes:  make([]int, partitionCount),
	}
}

func (p *Partitions) removeNodes(nodes []*Node) {
	for i := range p.Replicas {
		for j := range p.Replicas[i] {
			for k := range nodes {
				if p.Replicas[i][j] == nodes[k] {
					p.Replicas[i][j] = nil
				}
			}
		}
	}
}

func (p *Partitions) setReplicaCount(replicaCount int) {
	if len(p.Replicas) < replicaCount {
		i := len(p.Replicas)

		// Extend the size
		for ; i < replicaCount; i++ {
			p.Replicas = append(p.Replicas, make([]*Node, _PARTITIONS))
		}
	} else {
		// Reduce the size
		p.Replicas = p.Replicas[:replicaCount]
	}
}

// Copy partition map while reserving space for a new replica count.
func (p *Partitions) clone() *Partitions {
	replicas := make([][]*Node, len(p.Replicas))

	for i := range p.Replicas {
		r := make([]*Node, len(p.Replicas[i]))
		copy(r, p.Replicas[i])
		replicas[i] = r
	}

	regimes := make([]int, len(p.regimes))
	copy(regimes, p.regimes)

	return &Partitions{
		Replicas: replicas,
		SCMode:   p.SCMode,
		regimes:  regimes,
	}
}

/*

	partitionMap

*/

type partitionMap map[string]*Partitions

// removeNodes removes all the referenced nodes stored in the lists
// when they are removed from the cluster
func (pm partitionMap) removeNodes(nodes []*Node) {
	if len(nodes) == 0 {
		return
	}

	for _, partitions := range pm {
		partitions.removeNodes(nodes)
	}
}

// cleanup removes all the references stored in the lists
// to help the GC identify the unused pointers.
func (pm partitionMap) cleanup() {
	for ns, partitions := range pm {
		for i := range partitions.Replicas {
			for j := range partitions.Replicas[i] {
				partitions.Replicas[i][j] = nil
			}
			partitions.Replicas[i] = nil
		}

		partitions.Replicas = nil
		partitions.regimes = nil

		delete(pm, ns)
	}
}

// String implements stringer interface for partitionMap
func (pm partitionMap) clone() partitionMap {
	// Make deep copy of map.
	pmap := make(partitionMap, len(pm))
	for ns := range pm {
		pmap[ns] = pm[ns].clone()
	}
	return pmap
}

// String implements stringer interface for partitionMap
func (pm partitionMap) String() string {
	res := bytes.Buffer{}
	for ns, partitions := range pm {
		res.WriteString("-----------------------------------------------------------------------\n")
		res.WriteString("Namespace: " + ns + "\n")
		res.WriteString(fmt.Sprintf("Regimes: %v\n", partitions.regimes))
		res.WriteString(fmt.Sprintf("SCMode: %v\n", partitions.SCMode))
		replicaArray := partitions.Replicas
		for i, nodeArray := range replicaArray {
			if i == 0 {
				res.WriteString("\nMASTER:")
			} else {
				res.WriteString(fmt.Sprintf("\nReplica %d: ", i))
			}
			for partitionID, node := range nodeArray {
				res.WriteString(strconv.Itoa(partitionID) + "/")
				if node != nil {
					res.WriteString(node.host.String())
					res.WriteString(", ")
				} else {
					res.WriteString("nil, ")
				}
			}
			res.WriteString("\n")
		}
	}
	res.WriteString("\n")
	return res.String()
}

// naively validates the partition map
func (pm partitionMap) validate() Error {
	masterNodePartitionNotDefined := map[string][]int{}
	replicaNodePartitionNotDefined := map[string][]int{}
	var errs Error

	for nsName, partition := range pm {
		if len(partition.regimes) != _PARTITIONS {
			errs = chainErrors(newError(types.COMMON_ERROR, fmt.Sprintf("Wrong number of regimes for namespace `%s`. Must be %d, but found %d.", nsName, _PARTITIONS, len(partition.regimes))), errs)
		}

		for replica, partitionNodes := range partition.Replicas {
			if len(partitionNodes) != _PARTITIONS {
				errs = chainErrors(newError(types.COMMON_ERROR, fmt.Sprintf("Wrong number of partitions for namespace `%s`, replica `%d`. Must be %d, but found %d.", nsName, replica, _PARTITIONS, len(partitionNodes))), errs)
			}

			for pIndex, node := range partitionNodes {
				if node == nil {
					if replica == 0 {
						masterNodePartitionNotDefined[nsName] = append(masterNodePartitionNotDefined[nsName], pIndex)
					} else {
						replicaNodePartitionNotDefined[nsName] = append(replicaNodePartitionNotDefined[nsName], pIndex)
					}
				}
			}
		}
	}

	if errs != nil || len(masterNodePartitionNotDefined) > 0 || len(replicaNodePartitionNotDefined) > 0 {
		for nsName, partitionList := range masterNodePartitionNotDefined {
			errs = chainErrors(newError(types.COMMON_ERROR, fmt.Sprintf("Master partition nodes not defined for namespace `%s`: %d out of %d", nsName, len(partitionList), _PARTITIONS)), errs)
		}

		for nsName, partitionList := range replicaNodePartitionNotDefined {
			errs = chainErrors(newError(types.COMMON_ERROR, fmt.Sprintf("Replica partition nodes not defined for namespace `%s`: %d", nsName, len(partitionList))), errs)
		}

		errs = chainErrors(ErrInvalidPartitionMap.err(), errs)
		return errs
	}

	return nil
}

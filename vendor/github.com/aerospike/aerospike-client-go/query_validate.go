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
	"strconv"

	. "github.com/aerospike/aerospike-client-go/types"
)

func queryValidateBegin(node *Node, namespace string) (int64, error) {
	if !node.supportsClusterStable.Get() {
		return 0, nil
	}

	// Fail when cluster is in migration.
	cmd := "cluster-stable:namespace=" + namespace
	result, err := node.requestInfo(node.cluster.clientPolicy.Timeout, cmd)
	if err != nil {
		return -1, err
	}

	i, err := strconv.ParseInt(result[cmd], 16, 64)
	if err == nil {
		return i, nil
	}

	// Yes, even scans return QUERY_ABORTED.
	return -1, newAerospikeNodeError(node, QUERY_ABORTED, "Cluster is in migration:", result[cmd])
}

func queryValidate(node *Node, namespace string, expectedKey int64) error {
	if expectedKey == 0 || !node.supportsClusterStable.Get() {
		return nil
	}

	// Fail when cluster is in migration.
	clusterKey, err := queryValidateBegin(node, namespace)
	if err != nil {
		return err
	}

	if clusterKey != expectedKey {
		return newAerospikeNodeError(node, QUERY_ABORTED, fmt.Sprintf("Cluster is in migration: %d != %d", expectedKey, clusterKey))
	}

	return nil
}

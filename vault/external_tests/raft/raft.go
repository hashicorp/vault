// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rafttests

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
)

func Raft_Configuration_Test(t *testing.T, cluster testcluster.VaultCluster) {
	client := cluster.Nodes()[0].APIClient()
	secret, err := client.Logical().Read("sys/storage/raft/configuration")
	if err != nil {
		t.Fatal(err)
	}
	servers := secret.Data["config"].(map[string]interface{})["servers"].([]interface{})
	found := make(map[string]struct{})
	for _, s := range servers {
		server := s.(map[string]interface{})
		nodeID := server["node_id"].(string)
		leader := server["leader"].(bool)
		switch nodeID {
		case "core-0":
			if !leader {
				t.Fatalf("expected server to be leader: %#v", server)
			}
		default:
			if leader {
				t.Fatalf("expected server to not be leader: %#v", server)
			}
		}

		found[nodeID] = struct{}{}
	}
	expected := map[string]struct{}{}
	for i := range cluster.Nodes() {
		expected[fmt.Sprintf("core-%d", i)] = struct{}{}
	}
	if diff := deep.Equal(expected, found); len(diff) > 0 {
		t.Fatalf("configuration mismatch, diff: %v", diff)
	}
}

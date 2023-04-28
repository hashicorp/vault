// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package replication_binary

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

// TestStandardPerfReplication_Docker tests that we can create two 3-node
// clusters of docker containers and connect them using perf replication.
func TestStandardPerfReplication_Docker(t *testing.T) {
	if !constants.IsEnterprise {
		// Disable on OSS since this needs an ent binary (or docker image) to work
		t.Skip()
	}

	r, err := docker.NewReplicationSetDocker(t)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cleanup()

	err = r.StandardPerfReplication(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

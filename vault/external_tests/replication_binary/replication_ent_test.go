// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package replication_binary

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func TestStandardPerfReplication_Docker(t *testing.T) {
	// Disable for now since this needs an ent binary to work
	t.Skip()

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

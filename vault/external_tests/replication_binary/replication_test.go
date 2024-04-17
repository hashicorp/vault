// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package replication_binary

/*
Example of how to use docker.NewReplicationSetDocker(t), assuming
you point VAULT_BINARY to an Enterprise Vault binary:

import (




)

// TestStandardPerfReplication_Docker tests that we can create two 3-node
// clusters of docker containers and connect them using perf replication.
func TestStandardPerfReplication_Docker(t *testing.T) {
	r, err := docker.NewReplicationSetDocker(t)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = r.StandardPerfReplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

*/

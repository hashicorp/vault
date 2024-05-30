// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package misc

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/mitchellh/mapstructure"
)

// TestRecovery_Docker exercises recovery mode.  It starts a single node raft
// cluster, writes some data, then restarts it and makes sure that we can read
// the data (that's mostly to make sure that our framework is properly handling
// a volume that persists across runs.)  It then starts the node in recovery mode
// and deletes the data via sys/raw, and finally restarts it in normal mode and
// makes sure the data has been deleted.
func TestRecovery_Docker(t *testing.T) {
	ctx := context.TODO()

	t.Parallel()
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}
	opts := &docker.DockerClusterOptions{
		ImageRepo: "hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:    "latest",
		VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			NumCores: 1,
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
				// If you want the test to run faster locally, you could
				// uncomment this performance_multiplier change.
				//StorageOptions: map[string]string{
				//	"performance_multiplier": "1",
				//},
			},
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	var secretUUID string
	{
		client := cluster.Nodes()[0].APIClient()
		if err := client.Sys().Mount("secret/", &api.MountInput{
			Type: "kv-v1",
		}); err != nil {
			t.Fatal(err)
		}

		fooVal := map[string]interface{}{"bar": 1.0}
		_, err := client.Logical().Write("secret/foo", fooVal)
		if err != nil {
			t.Fatal(err)
		}
		secret, err := client.Logical().List("secret/")
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(secret.Data["keys"], []interface{}{"foo"}); len(diff) > 0 {
			t.Fatalf("got=%v, want=%v, diff: %v", secret.Data["keys"], []string{"foo"}, diff)
		}
		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}
		secretMount := mounts["secret/"]
		if secretMount == nil {
			t.Fatalf("secret mount not found, mounts: %v", mounts)
		}
		secretUUID = secretMount.UUID
	}

	listSecrets := func() []string {
		client := cluster.Nodes()[0].APIClient()
		secret, err := client.Logical().List("secret/")
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			return nil
		}
		var result []string
		err = mapstructure.Decode(secret.Data["keys"], &result)
		return result
	}

	restart := func() {
		cluster.Nodes()[0].(*docker.DockerClusterNode).Stop()

		err := cluster.Nodes()[0].(*docker.DockerClusterNode).Start(ctx, opts)
		if err != nil {
			t.Fatalf("node restart post-recovery failed: %v", err)
		}

		err = testcluster.UnsealAllNodes(ctx, cluster)
		if err != nil {
			t.Fatalf("node unseal post-recovery failed: %v", err)
		}

		_, err = testcluster.WaitForActiveNode(ctx, cluster)
		if err != nil {
			t.Fatalf("node didn't become active: %v", err)
		}
	}

	restart()
	if len(listSecrets()) == 0 {
		t.Fatal("expected secret to still be there")
	}

	// Now bring it up in recovery mode.
	{
		cluster.Nodes()[0].(*docker.DockerClusterNode).Stop()

		newOpts := *opts
		opts := &newOpts
		opts.Args = []string{"-recovery"}
		opts.StartProbe = func(client *api.Client) error {
			// In recovery mode almost no paths are supported, and pretty much
			// the only ones that don't require a recovery token are the ones used
			// to generate a recovery token.
			_, err := client.Sys().GenerateRecoveryOperationTokenStatusWithContext(ctx)
			return err
		}
		err := cluster.Nodes()[0].(*docker.DockerClusterNode).Start(ctx, opts)
		if err != nil {
			t.Fatalf("node restart with -recovery failed: %v", err)
		}
		client := cluster.Nodes()[0].APIClient()

		recoveryToken, err := testcluster.GenerateRoot(cluster, testcluster.GenerateRecovery)
		if err != nil {
			t.Fatalf("recovery token generation failed: %v", err)
		}
		_, err = testcluster.GenerateRoot(cluster, testcluster.GenerateRecovery)
		if err == nil {
			t.Fatal("expected second generate-root to fail")
		}
		client.SetToken(recoveryToken)

		secret, err := client.Logical().List(path.Join("sys/raw/logical", secretUUID))
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(secret.Data["keys"], []interface{}{"foo"}); len(diff) > 0 {
			t.Fatalf("got=%v, want=%v, diff: %v", secret.Data, []string{"foo"}, diff)
		}

		_, err = client.Logical().Delete(path.Join("sys/raw/logical", secretUUID, "foo"))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Now go back to regular mode and verify that our changes are present
	restart()
	if len(listSecrets()) != 0 {
		t.Fatal("expected secret to still be gone")
	}
}

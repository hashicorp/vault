// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package misc

import (
	"context"
	"net/http"
	"path"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func TestRecovery_Docker(t *testing.T) {
	ctx := context.TODO()

	t.Parallel()
	//binary := os.Getenv("VAULT_BINARY")
	//if binary == "" {
	//	t.Skip("only running docker test when $VAULT_BINARY present")
	//}
	opts := &docker.DockerClusterOptions{
		ImageRepo: "hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag: "1.12.1",
		// VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			NumCores: 1,
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
				// If you want the test to run faster locally, you could
				// uncomment this performance_multiplier change.
				StorageOptions: map[string]string{
					"performance_multiplier": "1",
				},
			},
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)
	// defer cluster.Cleanup()

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

	cluster.Nodes()[0].(*docker.DockerClusterNode).Stop()

	getSecret := func() *api.Secret {
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

		client := cluster.Nodes()[0].APIClient()
		secret, err := client.Logical().List("secret/")
		if err != nil {
			t.Fatal(err)
		}
		return secret
	}

	if getSecret() == nil {
		t.Fatal("expected secret to still be there")
	}
	cluster.Nodes()[0].(*docker.DockerClusterNode).Stop()

	// Now bring it up in recovery mode.
	{
		newOpts := *opts
		opts := &newOpts
		opts.Args = []string{"-recovery"}
		opts.StartProbe = func(client *api.Client) error {
			// In recovery mode almost no paths are supported, and pretty much
			// the only ones that don't require a recovery token are the ones used
			// to generate a recovery token.
			r := client.NewRequest(http.MethodGet, "/v1/sys/generate-recovery-token/attempt")
			_, err := client.RawRequestWithContext(ctx, r)
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

	cluster.Nodes()[0].(*docker.DockerClusterNode).Stop()

	// Now go back to regular mode and verify that our changes are present
	if getSecret() != nil {
		t.Fatal("expected secret to still be gone")
	}
}

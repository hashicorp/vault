package api

import (
	"context"
	"reflect"
	"testing"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestKVHelpers(t *testing.T) {
	t.Parallel()

	// initialize test cluster
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv":    logicalKv.Factory,
			"kv-v2": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	core := cores[0].Core
	client := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// mount the KVv2 backend
	// (the test cluster has already mounted the KVv1 backend at "secret")
	err := client.Sys().MountWithContext(context.Background(), "secret-v2", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	secretData := map[string]interface{}{
		"foo": "bar",
	}

	//// v1 ////
	t.Run("kv v1 helpers", func(t *testing.T) {
		if err := client.KVv1("secret").Put(context.Background(), "my-secret", secretData); err != nil {
			t.Fatal(err)
		}

		secret, err := client.KVv1("secret").Get(context.Background(), "my-secret")
		if err != nil {
			t.Fatal(err)
		}

		if secret.Data["foo"] != "bar" {
			t.Fatalf("kv v1 secret did not contain expected value")
		}

		if err := client.KVv1("secret").Delete(context.Background(), "my-secret"); err != nil {
			t.Fatal(err)
		}
	})

	//// v2 ////
	t.Run("kv v2 helpers", func(t *testing.T) {
		// create a secret
		writtenSecret, err := client.KVv2("secret-v2").Put(context.Background(), "my-secret", secretData)
		if err != nil {
			t.Fatal(err)
		}
		if writtenSecret == nil || writtenSecret.VersionMetadata == nil {
			t.Fatal("kv v2 secret did not have expected contents")
		}

		secret, err := client.KVv2("secret-v2").Get(context.Background(), "my-secret")
		if err != nil {
			t.Fatal(err)
		}
		if secret.Data["foo"] != "bar" {
			t.Fatal("kv v2 secret did not contain expected value")
		}
		if secret.VersionMetadata.CreatedTime != writtenSecret.VersionMetadata.CreatedTime {
			t.Fatal("the created_time on the secret did not match the response from when it was created")
		}

		// get its full metadata
		fullMetadata, err := client.KVv2("secret-v2").GetMetadata(context.Background(), "my-secret")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(secret.CustomMetadata, fullMetadata.CustomMetadata) {
			t.Fatalf("custom metadata on the secret does not match the custom metadata in the full metadata")
		}

		// create a second version
		_, err = client.KVv2("secret-v2").Put(context.Background(), "my-secret", map[string]interface{}{
			"foo": "baz",
		})
		if err != nil {
			t.Fatal(err)
		}

		s2, err := client.KVv2("secret-v2").Get(context.Background(), "my-secret")
		if err != nil {
			t.Fatal(err)
		}
		if s2.Data["foo"] != "baz" {
			t.Fatalf("second version of secret did not have expected contents")
		}
		if s2.VersionMetadata.Version != 2 {
			t.Fatalf("wrong version of kv v2 secret was read, expected 2 but got %d", s2.VersionMetadata.Version)
		}

		// get a specific past version
		s1, err := client.KVv2("secret-v2").GetVersion(context.Background(), "my-secret", 1)
		if err != nil {
			t.Fatal(err)
		}
		if s1.VersionMetadata.Version != 1 {
			t.Fatalf("wrong version of kv v2 secret was read, expected 1 but got %d", s1.VersionMetadata.Version)
		}

		// delete that version
		if err = client.KVv2("secret-v2").DeleteVersions(context.Background(), "my-secret", []int{1}); err != nil {
			t.Fatal(err)
		}

		s1AfterDelete, err := client.KVv2("secret-v2").GetVersion(context.Background(), "my-secret", 1)
		if err != nil {
			t.Fatal(err)
		}

		if s1AfterDelete.VersionMetadata.DeletionTime.IsZero() {
			t.Fatalf("the deletion_time in the first version of the secret was not updated")
		}

		if s1AfterDelete.Data != nil {
			t.Fatalf("data still exists on the first version of the secret despite this version being deleted")
		}

		// check that KVOption works
		_, err = client.KVv2("secret-v2").Put(context.Background(), "my-secret", map[string]interface{}{
			"meow": "woof",
		}, api.WithCheckAndSet(99))
		if err == nil {
			t.Fatalf("expected error from trying to update different version from check-and-set value")
		}

		versions, err := client.KVv2("secret-v2").GetVersionsAsList(context.Background(), "my-secret")
		if err != nil {
			t.Fatal(err)
		}

		if len(versions) != 2 {
			t.Fatalf("expected there to be 2 versions of the secret but got %d", len(versions))
		}

		if versions[0].Version != 1 {
			t.Fatalf("incorrect value for version; expected 1 but got %d", versions[0].Version)
		}
	})
}

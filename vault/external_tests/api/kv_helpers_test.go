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
		mountPath := "secret"
		secretPath := "my-secret"
		if err := client.KVv1(mountPath).Put(context.Background(), secretPath, secretData); err != nil {
			t.Fatal(err)
		}

		secret, err := client.KVv1(mountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		if secret.Data["foo"] != "bar" {
			t.Fatalf("kv v1 secret did not contain expected value")
		}

		if err := client.KVv1(mountPath).Delete(context.Background(), secretPath); err != nil {
			t.Fatal(err)
		}
	})

	//// v2 ////
	t.Run("kv v2 helpers", func(t *testing.T) {
		mountPath := "secret-v2"
		secretPath := "my-secret"
		// create a secret
		writtenSecret, err := client.KVv2(mountPath).Put(context.Background(), secretPath, secretData)
		if err != nil {
			t.Fatal(err)
		}
		if writtenSecret == nil || writtenSecret.VersionMetadata == nil {
			t.Fatal("kv v2 secret did not have expected contents")
		}

		secret, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
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
		fullMetadata, err := client.KVv2(mountPath).GetMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(secret.CustomMetadata, fullMetadata.CustomMetadata) {
			t.Fatalf("custom metadata on the secret does not match the custom metadata in the full metadata")
		}

		// create a second version
		_, err = client.KVv2(mountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"foo": "baz",
		})
		if err != nil {
			t.Fatal(err)
		}

		s2, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
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
		s1, err := client.KVv2(mountPath).GetVersion(context.Background(), secretPath, 1)
		if err != nil {
			t.Fatal(err)
		}
		if s1.VersionMetadata.Version != 1 {
			t.Fatalf("wrong version of kv v2 secret was read, expected 1 but got %d", s1.VersionMetadata.Version)
		}

		// delete that version
		if err = client.KVv2(mountPath).DeleteVersions(context.Background(), secretPath, []int{1}); err != nil {
			t.Fatal(err)
		}

		s1AfterDelete, err := client.KVv2(mountPath).GetVersion(context.Background(), secretPath, 1)
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
		////
		// WithCheckAndSet
		_, err = client.KVv2(mountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"meow": "woof",
		}, api.WithCheckAndSet(99))
		// should fail
		if err == nil {
			t.Fatalf("expected error from trying to update different version from check-and-set value using WithCheckAndSet")
		}

		// WithOption (generic)
		_, err = client.KVv2(mountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"bow": "wow",
		}, api.WithOption("cas", 99))
		// should fail
		if err == nil {
			t.Fatalf("expected error from trying to update different version from check-and-set value using generic WithOption")
		}

		// WithMergeMethod Patch (implicit)
		patch, err := client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "cat",
		})
		if err != nil {
			t.Fatal(err)
		}
		if patch.VersionMetadata.Version != 3 {
			t.Fatalf("incorrect version %d, expected 3", patch.VersionMetadata.Version)
		}

		// WithMergeMethod Patch (explicit)
		patchExp, err := client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"rat": "mouse",
		}, api.WithMergeMethod(api.KVMergeMethodPatch))
		if err != nil {
			t.Fatal(err)
		}
		if patchExp.VersionMetadata.Version != 4 {
			t.Fatalf("incorrect version %d, expected 4", patchExp.VersionMetadata.Version)
		}

		// WithMergeMethod RW
		patchRW, err := client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"bird": "tweet",
		}, api.WithMergeMethod(api.KVMergeMethodReadWrite))
		if err != nil {
			t.Fatal(err)
		}
		if patchRW.VersionMetadata.Version != 5 {
			t.Fatalf("incorrect version %d, expected 5", patchRW.VersionMetadata.Version)
		}

		// patch something that doesn't exist
		_, err = client.KVv2(mountPath).Patch(context.Background(), "nonexistent-secret", map[string]interface{}{
			"no": "nope",
		})
		if err == nil {
			t.Fatal("expected error from trying to patch something that doesn't exist")
		}

		secretAfterPatches, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		_, ok := secretAfterPatches.Data["dog"]
		if !ok {
			t.Fatalf("secret did not contain data patched with implicit Patch method")
		}
		_, ok = secretAfterPatches.Data["rat"]
		if !ok {
			t.Fatalf("secret did not contain data patched with explicit Patch method")
		}
		_, ok = secretAfterPatches.Data["bird"]
		if !ok {
			t.Fatalf("secret did not contain data patched with RW method")
		}
		value, ok := secretAfterPatches.Data["foo"]
		if !ok || value != "baz" {
			t.Fatalf("secret did not keep original data after patch")
		}

		// patch an existing field
		_, err = client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "pug",
		})
		if err != nil {
			t.Fatal(err)
		}
		patchedFieldKV, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		v, ok := patchedFieldKV.Data["dog"]
		if !ok || v != "pug" {
			t.Fatalf("secret's data was not replaced by patch")
		}

		// delete a key in a secret via patch
		_, err = client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": nil,
		})
		if err != nil {
			t.Fatal(err)
		}
		deletedFieldKV, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		_, ok = deletedFieldKV.Data["dog"]
		if ok {
			t.Fatalf("secret key \"dog\" should have been removed by nil patch")
		}

		// set a key to an empty string via patch
		_, err = client.KVv2(mountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "",
		})
		if err != nil {
			t.Fatal(err)
		}
		emptyValueKV, err := client.KVv2(mountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		v, ok = emptyValueKV.Data["dog"]
		if !ok || v != "" {
			t.Fatalf("secret key \"dog\" should have an empty string value")
		}

		////

		// get versions as list
		versions, err := client.KVv2(mountPath).GetVersionsAsList(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		expectedLength := 8
		if len(versions) != expectedLength {
			t.Fatalf("expected there to be %d versions of the secret but got %d", expectedLength, len(versions))
		}

		if versions[0].Version != 1 || versions[len(versions)-1].Version != expectedLength {
			t.Fatalf("versions list is not ordered as expected")
		}
	})
}

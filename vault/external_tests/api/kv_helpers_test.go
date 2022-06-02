package api

import (
	"context"
	"reflect"
	"testing"
	"time"

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

		// put and patch metadata
		////
		noDataSecretPath := "empty"

		// create a secret with metadata but no data
		deleteVersionAfter := 5 * time.Hour
		maxVersions := 5
		customMetadata := map[string]interface{}{"ape": "gorilla"}
		err = client.KVv2(mountPath).PutMetadata(context.Background(), noDataSecretPath, api.KVMetadataInput{
			DeleteVersionAfter: &deleteVersionAfter,
			MaxVersions:        &maxVersions,
			CustomMetadata:     customMetadata,
		})
		if err != nil {
			t.Fatal(err)
		}

		// get its metadata to make sure it was created successfully
		md, err := client.KVv2(mountPath).GetMetadata(context.Background(), noDataSecretPath)
		if err != nil {
			t.Fatal(err)
		}
		if md.CreatedTime.IsZero() {
			t.Fatalf("secret metadata was not populated as expected: %v", err)
		}
		if len(md.Versions) > 0 {
			t.Fatalf("no secret versions should have been created since only metadata was populated")
		}

		// replace all modifiable metadata fields
		casRequired := true
		deleteVersionAfter = 6 * time.Hour
		maxVersions = 6
		customMetadata = map[string]interface{}{"ape": "orangutan", "cat": "tabby"}
		err = client.KVv2(mountPath).PutMetadata(context.Background(), noDataSecretPath, api.KVMetadataInput{
			CASRequired:        &casRequired,
			DeleteVersionAfter: &deleteVersionAfter,
			MaxVersions:        &maxVersions,
			CustomMetadata:     customMetadata,
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that metadata was replaced
		md2, err := client.KVv2(mountPath).GetMetadata(context.Background(), noDataSecretPath)
		if err != nil {
			t.Fatal(err)
		}
		if md.CASRequired == md2.CASRequired || md.MaxVersions == md2.MaxVersions || md.DeleteVersionAfter == md2.DeleteVersionAfter || reflect.DeepEqual(md.CustomMetadata, md2.CustomMetadata) {
			t.Fatalf("metadata fields should have been updated by PutMetadata")
		}

		// now let's try a patch
		maxVersions = 7
		err = client.KVv2(mountPath).PatchMetadata(context.Background(), noDataSecretPath, api.KVMetadataInput{
			MaxVersions:    &maxVersions,
			CustomMetadata: map[string]interface{}{"ape": nil, "rat": "brown"},
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that the metadata was only partially replaced
		md3, err := client.KVv2(mountPath).GetMetadata(context.Background(), noDataSecretPath)
		if err != nil {
			t.Fatal(err)
		}
		if md2.CASRequired != md3.CASRequired || md2.DeleteVersionAfter != md3.DeleteVersionAfter {
			t.Fatalf("expected fields to remain unchanged but they were updated")
		}
		if md3.MaxVersions == 0 {
			t.Fatalf("field was reset to its zero value when it should not have been")
		}
		if md2.MaxVersions == md3.MaxVersions {
			t.Fatalf("expected field to be updated but it remained unchanged")
		}

		// let's check the custom metadata was updated correctly
		if r, ok := md3.CustomMetadata["rat"]; ok {
			if r != "brown" {
				t.Fatalf("expected value to be \"brown\"")
			}
		} else {
			t.Fatalf("expected there to be a new \"rat\" key")
		}

		if _, ok := md3.CustomMetadata["ape"]; ok {
			t.Fatalf("expected \"ape\" key to be removed")
		}

		if _, ok := md3.CustomMetadata["cat"]; !ok {
			t.Fatalf("did not expect \"cat\" key to be removed")
		}

		// now let's do another patch to test the "explicit zero value" use case
		explicitFalse := false
		explicitZero := 0
		explicitTimeZero, err := time.ParseDuration("0s")
		if err != nil {
			t.Fatal(err)
		}
		err = client.KVv2(mountPath).PatchMetadata(context.Background(), noDataSecretPath, api.KVMetadataInput{
			CASRequired:        &explicitFalse,
			MaxVersions:        &explicitZero,
			DeleteVersionAfter: &explicitTimeZero,
			CustomMetadata:     map[string]interface{}{},
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that those fields were reset to their zero value
		md4, err := client.KVv2(mountPath).GetMetadata(context.Background(), noDataSecretPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(md4.CustomMetadata) > 0 {
			t.Fatalf("expected empty map to cause deletion of all custom metadata")
		}

		if md4.MaxVersions != 0 || md4.CASRequired != false {
			t.Fatalf("expected fields to be reset to their zero values but they were %d and %t instead", md4.MaxVersions, md4.CASRequired)
		}

		// TODO: Bring this test back to life once the bug is fixed where we can't actually reset delete-version-after to zero...
		// if md4.DeleteVersionAfter.String() != "0s" {
		// 	t.Fatalf("expected delete-version-after to be reset to its zero value but instead it was %d", md4.DeleteVersionAfter.String())
		// }

		////
	})
}

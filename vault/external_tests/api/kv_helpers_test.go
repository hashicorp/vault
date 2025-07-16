// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	v1MountPath = "secret"
	v2MountPath = "secret-v2"
	secretPath  = "my-secret"
)

var secretData = map[string]interface{}{
	"foo": "bar",
}

// setupKVv2Test creates the secret that will be used in each KV v2 subtest. It
// returns a function (that should be deferred whenever setupKVv2Test is called)
// which will perform the cleanup of all existing versions of the secret, as
// well as the secret that was written for comparison.
func setupKVv2Test(t *testing.T, client *api.Client) (func(t *testing.T), *api.KVSecret) {
	writtenSecret, err := client.KVv2(v2MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}
	if writtenSecret == nil || writtenSecret.VersionMetadata == nil {
		t.Fatal("secret created during kv v2 subtest setup did not have expected contents")
	}

	return func(t *testing.T) {
		err := client.KVv2(v2MountPath).DeleteMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
	}, writtenSecret
}

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

	//// v1 ////
	t.Run("kv v1: put, get, and delete data", func(t *testing.T) {
		if err := client.KVv1(v1MountPath).Put(context.Background(), secretPath, secretData); err != nil {
			t.Fatal(err)
		}

		secret, err := client.KVv1(v1MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		if secret.Data["foo"] != "bar" {
			t.Fatalf("kv v1 secret did not contain expected value")
		}

		if err := client.KVv1(v1MountPath).Delete(context.Background(), secretPath); err != nil {
			t.Fatal(err)
		}

		_, err = client.KVv1(v1MountPath).Get(context.Background(), secretPath)
		if !errors.Is(err, api.ErrSecretNotFound) {
			t.Fatalf("KVv1.Get is expected to return an api.ErrSecretNotFound wrapped error after secret had been deleted; got %v", err)
		}
	})

	t.Run("kv v1: get secret that does not exist", func(t *testing.T) {
		_, err = client.KVv1(v1MountPath).Get(context.Background(), "does/not/exist")
		if err == nil {
			t.Fatalf("KVv1.Get is expected to return an error for a missing secret")
		}
		if !errors.Is(err, api.ErrSecretNotFound) {
			t.Fatalf("KVv1.Get is expected to return an api.ErrSecretNotFound wrapped error for a missing secret; got %v", err)
		}
	})

	//// v2 ////
	t.Run("kv v2: get data and full metadata", func(t *testing.T) {
		teardownTest, originalSecret := setupKVv2Test(t, client)
		defer teardownTest(t)

		secret, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if secret.Data["foo"] != "bar" {
			t.Fatal("kv v2 secret did not contain expected value")
		}
		if secret.VersionMetadata.CreatedTime != originalSecret.VersionMetadata.CreatedTime {
			t.Fatal("the created_time on the secret did not match the response from when it was created")
		}

		// get its full metadata
		fullMetadata, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(secret.CustomMetadata, fullMetadata.CustomMetadata) {
			t.Fatalf("custom metadata on the secret does not match the custom metadata in the full metadata")
		}
	})

	t.Run("kv v2: get secret that does not exist", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		_, err = client.KVv2(v2MountPath).Get(context.Background(), "does/not/exist")
		if !errors.Is(err, api.ErrSecretNotFound) {
			t.Fatalf("KVv2.Get is expected to return an api.ErrSecretNotFound wrapped error for a missing secret; got %v", err)
		}

		_, err = client.KVv2(v2MountPath).GetMetadata(context.Background(), "does/not/exist")
		if !errors.Is(err, api.ErrSecretNotFound) {
			t.Fatalf("KVv2.GetMetadata is expected to return an api.ErrSecretNotFound wrapped error for a missing secret; got %v", err)
		}

		_, err = client.KVv2(v2MountPath).GetVersion(context.Background(), secretPath, 99)
		if !errors.Is(err, api.ErrSecretNotFound) {
			t.Fatalf("KVv2.GetVersion is expected to return an api.ErrSecretNotFound wrapped error for a missing secret version; got %v", err)
		}
	})

	t.Run("kv v2: multiple versions", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// create a second version
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"foo": "baz",
		})
		if err != nil {
			t.Fatal(err)
		}

		s2, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if s2.Data["foo"] != "baz" {
			t.Fatalf("second version of secret did not have expected contents")
		}
		if s2.VersionMetadata.Version != 2 {
			t.Fatalf("wrong version of kv v2 secret was read, expected 2 but got %d", s2.VersionMetadata.Version)
		}
	})

	t.Run("kv v2: delete and undelete", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// create a second version
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"foo": "baz",
		})
		if err != nil {
			t.Fatal(err)
		}

		// get a specific past version
		s1, err := client.KVv2(v2MountPath).GetVersion(context.Background(), secretPath, 1)
		if err != nil {
			t.Fatal(err)
		}
		if s1.VersionMetadata.Version != 1 {
			t.Fatalf("wrong version of kv v2 secret was read, expected 1 but got %d", s1.VersionMetadata.Version)
		}

		// delete that version
		if err = client.KVv2(v2MountPath).DeleteVersions(context.Background(), secretPath, []int{1}); err != nil {
			t.Fatal(err)
		}

		s1AfterDelete, err := client.KVv2(v2MountPath).GetVersion(context.Background(), secretPath, 1)
		if err != nil {
			t.Fatal(err)
		}

		if s1AfterDelete.VersionMetadata.DeletionTime.IsZero() {
			t.Fatalf("the deletion_time in the first version of the secret was not updated")
		}

		if s1AfterDelete.Data != nil {
			t.Fatalf("data still exists on the first version of the secret despite this version being deleted")
		}

		// undelete it
		err = client.KVv2(v2MountPath).Undelete(context.Background(), secretPath, []int{1})
		if err != nil {
			t.Fatal(err)
		}

		s1AfterUndelete, err := client.KVv2(v2MountPath).GetVersion(context.Background(), secretPath, 1)
		if err != nil {
			t.Fatal(err)
		}

		if s1AfterUndelete.Data == nil {
			t.Fatalf("data is empty for the first version of the secret despite this version being undeleted")
		}
	})

	t.Run("kv v2: destroy", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		err = client.KVv2(v2MountPath).Destroy(context.Background(), secretPath, []int{1})
		if err != nil {
			t.Fatal(err)
		}

		destroyedSecret, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		if !destroyedSecret.VersionMetadata.Destroyed {
			t.Fatalf("expected secret to be destroyed but it wasn't")
		}
	})

	t.Run("kv v2: use named functional options and generic WithOption", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// check that KVOption works
		// WithCheckAndSet
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"meow": "woof",
		}, api.WithCheckAndSet(99))
		// should fail
		if err == nil {
			t.Fatalf("expected error from trying to update different version from check-and-set value using WithCheckAndSet")
		}

		// WithOption (generic)
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"bow": "wow",
		}, api.WithOption("cas", 99))
		// should fail
		if err == nil {
			t.Fatalf("expected error from trying to update different version from check-and-set value using generic WithOption")
		}
	})

	t.Run("kv v2: patch", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// WithMergeMethod Patch (implicit)
		patch, err := client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "cat",
		})
		if err != nil {
			t.Fatal(err)
		}
		if patch.VersionMetadata.Version != 2 {
			t.Fatalf("incorrect version %d, expected 2", patch.VersionMetadata.Version)
		}

		// WithMergeMethod Patch (explicit)
		patchExp, err := client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"rat": "mouse",
		}, api.WithMergeMethod(api.KVMergeMethodPatch))
		if err != nil {
			t.Fatal(err)
		}
		if patchExp.VersionMetadata.Version != 3 {
			t.Fatalf("incorrect version %d, expected 3", patchExp.VersionMetadata.Version)
		}

		// WithMergeMethod RW
		patchRW, err := client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"bird": "tweet",
		}, api.WithMergeMethod(api.KVMergeMethodReadWrite))
		if err != nil {
			t.Fatal(err)
		}
		if patchRW.VersionMetadata.Version != 4 {
			t.Fatalf("incorrect version %d, expected 4", patchRW.VersionMetadata.Version)
		}

		secretAfterPatches, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
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
		if !ok || value != "bar" {
			t.Fatalf("secret did not keep original data after patch")
		}

		// patch an existing field
		_, err = client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "pug",
		})
		if err != nil {
			t.Fatal(err)
		}
		patchedFieldKV, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		v, ok := patchedFieldKV.Data["dog"]
		if !ok || v != "pug" {
			t.Fatalf("secret's data was not replaced by patch")
		}

		// delete a key in a secret via patch
		_, err = client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": nil,
		})
		if err != nil {
			t.Fatal(err)
		}
		deletedFieldKV, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		_, ok = deletedFieldKV.Data["dog"]
		if ok {
			t.Fatalf("secret key \"dog\" should have been removed by nil patch")
		}

		// set a key to an empty string via patch
		_, err = client.KVv2(v2MountPath).Patch(context.Background(), secretPath, map[string]interface{}{
			"dog": "",
		})
		if err != nil {
			t.Fatal(err)
		}
		emptyValueKV, err := client.KVv2(v2MountPath).Get(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		v, ok = emptyValueKV.Data["dog"]
		if !ok || v != "" {
			t.Fatalf("secret key \"dog\" should have an empty string value")
		}
	})

	t.Run("kv v2: patch a secret that does not exist", func(t *testing.T) {
		for _, method := range [][]api.KVOption{
			{},
			{api.WithMergeMethod(api.KVMergeMethodPatch)},
			{api.WithMergeMethod(api.KVMergeMethodReadWrite)},
		} {
			_, err = client.KVv2(v2MountPath).Patch(
				context.Background(),
				"does/not/exist",
				map[string]interface{}{"no": "nope"},
				method...,
			)
			if !errors.Is(err, api.ErrSecretNotFound) {
				t.Fatalf("expected an api.ErrSecretNotFound wrapped error from trying to patch something that doesn't exist for %v method; got: %v", method, err)
			}
		}
	})

	t.Run("kv v2: roll back to an old version", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// create a second version
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"color": "yellow",
		})
		if err != nil {
			t.Fatal(err)
		}

		// get versions as list
		versions, err := client.KVv2(v2MountPath).GetVersionsAsList(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		expectedLength := 2
		if len(versions) != expectedLength {
			t.Fatalf("expected there to be %d versions of the secret but got %d", expectedLength, len(versions))
		}

		if versions[0].Version != 1 || versions[len(versions)-1].Version != expectedLength {
			t.Fatalf("versions list is not ordered as expected")
		}

		// roll back to version 1
		rb, err := client.KVv2(v2MountPath).Rollback(context.Background(), secretPath, 1)
		if err != nil {
			t.Fatal(err)
		}
		if rb.VersionMetadata.Version != 3 {
			t.Fatalf("expected returned secret's version %d to be the latest version, which should be 3", rb.VersionMetadata.Version)
		}

		// destroy version 1
		err = client.KVv2(v2MountPath).Destroy(context.Background(), secretPath, []int{1})
		if err != nil {
			t.Fatal(err)
		}

		// roll back but fail
		_, err = client.KVv2(v2MountPath).Rollback(context.Background(), secretPath, 1)
		if err == nil {
			t.Fatalf("expected error from trying to rollback to destroyed version")
		}
	})

	t.Run("kv v2: delete all versions of a secret", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// create a second version
		_, err = client.KVv2(v2MountPath).Put(context.Background(), secretPath, map[string]interface{}{
			"color": "yellow",
		})
		if err != nil {
			t.Fatal(err)
		}

		// delete it all
		err = client.KVv2(v2MountPath).DeleteMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		versions, err := client.KVv2(v2MountPath).GetVersionsAsList(context.Background(), secretPath)
		if err == nil {
			t.Fatalf("expected to be unable to get list of versions since all metadata was destroyed")
		}
		if len(versions) > 0 {
			t.Fatalf("expected no versions of secret after deleting all metadata")
		}
	})

	t.Run("kv v2: create a secret with metadata but no data", func(t *testing.T) {
		// put and patch metadata
		////
		noDataSecretPath := "empty"

		// create a secret with metadata but no data
		err = client.KVv2(v2MountPath).PutMetadata(context.Background(), noDataSecretPath, api.KVMetadataPutInput{
			DeleteVersionAfter: 5 * time.Hour,
			MaxVersions:        5,
			CustomMetadata:     map[string]interface{}{"ape": "gorilla"},
		})
		if err != nil {
			t.Fatal(err)
		}

		// get its metadata to make sure it was created successfully
		md, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), noDataSecretPath)
		if err != nil {
			t.Fatal(err)
		}
		if md.CreatedTime.IsZero() {
			t.Fatalf("secret metadata was not populated as expected: %v", err)
		}
		if len(md.Versions) > 0 {
			t.Fatalf("no secret versions should have been created since only metadata was populated")
		}
	})

	t.Run("kv v2: put and patch metadata", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		md, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}

		// replace all modifiable metadata fields
		err = client.KVv2(v2MountPath).PutMetadata(context.Background(), secretPath, api.KVMetadataPutInput{
			CASRequired:        true,
			DeleteVersionAfter: 6 * time.Hour,
			MaxVersions:        6,
			CustomMetadata:     map[string]interface{}{"foo": "fwah", "cat": "tabby"},
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that metadata was replaced
		md2, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if md.CASRequired == md2.CASRequired || md.MaxVersions == md2.MaxVersions || md.DeleteVersionAfter == md2.DeleteVersionAfter || reflect.DeepEqual(md.CustomMetadata, md2.CustomMetadata) {
			t.Fatalf("metadata fields should have been updated by PutMetadata")
		}

		// now let's try a patch
		maxVersions := 7
		err = client.KVv2(v2MountPath).PatchMetadata(context.Background(), secretPath, api.KVMetadataPatchInput{
			MaxVersions:    &maxVersions,
			CustomMetadata: map[string]interface{}{"foo": nil, "rat": "brown"},
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that the metadata was only partially replaced
		md3, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), secretPath)
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

		if _, ok := md3.CustomMetadata["foo"]; ok {
			t.Fatalf("expected \"foo\" key to be removed")
		}

		if _, ok := md3.CustomMetadata["cat"]; !ok {
			t.Fatalf("did not expect \"cat\" key to be removed")
		}
	})

	t.Run("kv v2: patch with explicit zero values", func(t *testing.T) {
		teardownTest, _ := setupKVv2Test(t, client)
		defer teardownTest(t)

		// now let's do another patch to test the "explicit zero value" use case
		var (
			explicitFalse    bool
			explicitZero     int
			explicitTimeZero time.Duration
		)
		err = client.KVv2(v2MountPath).PatchMetadata(context.Background(), secretPath, api.KVMetadataPatchInput{
			CASRequired:        &explicitFalse,
			MaxVersions:        &explicitZero,
			DeleteVersionAfter: &explicitTimeZero,
			CustomMetadata:     map[string]interface{}{},
		})
		if err != nil {
			t.Fatal(err)
		}

		// check that those fields were reset to their zero value
		md4, err := client.KVv2(v2MountPath).GetMetadata(context.Background(), secretPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(md4.CustomMetadata) > 0 {
			t.Fatalf("expected empty map to cause deletion of all custom metadata")
		}

		if md4.MaxVersions != 0 || md4.CASRequired != false {
			t.Fatalf("expected fields to be reset to their zero values but they were %d and %t instead", md4.MaxVersions, md4.CASRequired)
		}

		if md4.DeleteVersionAfter.String() != "0s" {
			t.Fatalf("expected delete-version-after to be reset to its zero value but instead it was %s", md4.DeleteVersionAfter.String())
		}
	})
}

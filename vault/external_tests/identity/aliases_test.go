// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

func TestIdentityStore_ListAlias(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{
		Type: "github",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	var githubAccessor string
	for k, v := range mounts {
		if k == "github/" {
			githubAccessor = v.Accessor
			break
		}
	}
	if githubAccessor == "" {
		t.Fatal("did not find github accessor")
	}

	resp, err := client.Logical().Write("identity/entity", nil)
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatalf("expected a non-nil response")
	}

	entityID := resp.Data["id"].(string)

	// Create an alias
	resp, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testaliasname",
		"mount_accessor": githubAccessor,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	testAliasCanonicalID := resp.Data["canonical_id"].(string)
	testAliasAliasID := resp.Data["id"].(string)

	resp, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "entityalias",
		"mount_accessor": githubAccessor,
		"canonical_id":   entityID,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	entityAliasAliasID := resp.Data["id"].(string)

	resp, err = client.Logical().List("identity/entity-alias/id")
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys := resp.Data["keys"].([]interface{})
	if len(keys) != 2 {
		t.Fatalf("bad: length of alias IDs listed; expected: 2, actual: %d", len(keys))
	}

	// Do some due diligence on the key info
	aliasInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}
	aliasInfo := aliasInfoRaw.(map[string]interface{})
	for _, keyRaw := range keys {
		key := keyRaw.(string)
		infoRaw, ok := aliasInfo[key]
		if !ok {
			t.Fatal("expected key info")
		}
		info := infoRaw.(map[string]interface{})
		currName := "entityalias"
		if info["canonical_id"].(string) == testAliasCanonicalID {
			currName = "testaliasname"
		}
		switch {
		case info["name"].(string) != currName:
			t.Fatalf("bad name: %v", info["name"].(string))
		case info["mount_accessor"].(string) != githubAccessor:
			t.Fatalf("bad mount_path: %v", info["mount_accessor"].(string))
		}
	}

	// Now do the same with entity info
	resp, err = client.Logical().List("identity/entity/id")
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	keys = resp.Data["keys"].([]interface{})
	if len(keys) != 2 {
		t.Fatalf("bad: length of entity IDs listed; expected: 2, actual: %d", len(keys))
	}

	expectedEntityIDs := []string{entityID, testAliasCanonicalID}
	actualEntityIDs := make([]string, 0, len(keys))
	for _, keyRaw := range keys {
		actualEntityIDs = append(actualEntityIDs, keyRaw.(string))
	}
	sort.Strings(expectedEntityIDs)
	sort.Strings(actualEntityIDs)
	if !reflect.DeepEqual(expectedEntityIDs, actualEntityIDs) {
		t.Fatalf("bad: listed entity IDs; expected: %#v\n actual: %#v\n", expectedEntityIDs, actualEntityIDs)
	}

	entityInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}

	// This is basically verifying that the entity has the alias in key_info
	// that we expect to be tied to it, plus validates nested alias metadata.
	entityInfo := entityInfoRaw.(map[string]interface{})
	for _, keyRaw := range keys {
		key := keyRaw.(string)
		infoRaw, ok := entityInfo[key]
		if !ok {
			t.Fatal("expected key info")
		}
		info := infoRaw.(map[string]interface{})

		switch {
		case info["creation_time"].(string) == "":
			t.Fatal("expected entity creation_time")
		case info["last_update_time"].(string) == "":
			t.Fatal("expected entity last_update_time")
		case info["disabled"].(bool):
			t.Fatal("expected entity disabled to be false")
		}

		currAliasID := entityAliasAliasID
		if key == testAliasCanonicalID {
			currAliasID = testAliasAliasID
		}
		currAliases := info["aliases"].([]interface{})
		if len(currAliases) != 1 {
			t.Fatal("bad aliases length")
		}
		for _, v := range currAliases {
			curr := v.(map[string]interface{})
			switch {
			case curr["id"].(string) != currAliasID:
				t.Fatalf("bad alias id: %v", curr["id"])
			case curr["mount_accessor"].(string) != githubAccessor:
				t.Fatalf("bad mount accessor: %v", curr["mount_accessor"])
			case curr["mount_path"].(string) != "auth/github/":
				t.Fatalf("bad mount path: %v", curr["mount_path"])
			case curr["mount_type"].(string) != "github":
				t.Fatalf("bad mount type: %v", curr["mount_type"])
			case curr["creation_time"].(string) == "":
				t.Fatal("expected alias creation_time")
			case curr["last_update_time"].(string) == "":
				t.Fatal("expected alias last_update_time")
			}
		}
	}
}

// TestIdentityStore_DeprecatedIdentityAliasEndpoint_NoPanic verifies that updating
// an alias's name through the deprecated identity/alias/id/{id} endpoint
// succeeds instead of panicking and returning a 503 error.
func TestIdentityStore_DeprecatedIdentityAliasEndpoint_NoPanic(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	require.NoError(t, err)

	mounts, err := client.Sys().ListAuth()
	require.NoError(t, err)
	mountAccessor := mounts["userpass/"].Accessor

	entityResp, err := client.Logical().Write("identity/entity", nil)
	require.NoError(t, err)
	require.NotNil(t, entityResp)
	entityID := entityResp.Data["id"].(string)

	// creating the alias so we can later test the deprecated endpoint
	aliasResp, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testUser",
		"mount_accessor": mountAccessor,
		"canonical_id":   entityID,
	})
	require.NoError(t, err)
	require.NotNil(t, aliasResp)
	aliasID := aliasResp.Data["id"].(string)

	// Prior to the fix, this write panicked (returned as a 500) because
	// handleAliasCreateUpdate read external_id/issuer via d.Get regardless of
	// which path's schema was matched.
	updateResp, err := client.Logical().Write("identity/alias/id/"+aliasID, map[string]interface{}{
		"name": "testUser_2",
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, aliasID, updateResp.Data["id"])

	readResp, err := client.Logical().Read("identity/alias/id/" + aliasID)
	require.NoError(t, err)
	require.NotNil(t, readResp)
	require.Equal(t, "testUser_2", readResp.Data["name"])
}

// TestIdentityStore_RenameAlias_CannotMergeEntity verifies that an error is
// returned on an attempt to rename an alias to match another alias with the
// same mount accessor.  This used to result in a merge entity.
func TestIdentityStore_RenameAlias_CannotMergeEntity(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	// Now create a new unrelated entity and alias
	entityResp, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "bob-smith",
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, entityResp)
	}
	if entityResp == nil {
		t.Fatalf("expected a non-nil response")
	}

	aliasResp, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "bob",
		"mount_accessor": mountAccessor,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, aliasResp)
	}
	aliasID2 := aliasResp.Data["id"].(string)

	// Rename this new alias to have the same name as the one implicitly created by our login as bsmith
	_, err = client.Logical().Write("identity/entity-alias/id/"+aliasID2, map[string]interface{}{
		"name": "bsmith",
	})
	if err == nil {
		t.Fatal("expected rename over existing entity to fail")
	}
}

func TestIdentityStore_MergeEntities_FailsDueToClash(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	_, entityIdBob, aliasIdBob := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")

	// Create userpass login for alice
	_, err = client.Logical().Write("auth/userpass/users/alice", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, entityIdAlice, aliasIdAlice := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "alice-smith", "alice")

	// Perform entity merge
	mergeResp, err := client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    entityIdBob,
		"from_entity_ids": entityIdAlice,
	})
	if err == nil {
		t.Fatalf("Expected error upon merge. Resp:%#v", mergeResp)
	}
	if !strings.Contains(err.Error(), "toEntity and at least one fromEntity have aliases with the same mount accessor") {
		t.Fatalf("Error was not due to conflicting alias mount accessors. Error: %v", err)
	}
	if !strings.Contains(err.Error(), entityIdAlice) {
		t.Fatalf("Did not identify alice's entity (%s) as conflicting. Error: %v", entityIdAlice, err)
	}
	if !strings.Contains(err.Error(), entityIdBob) {
		t.Fatalf("Did not identify bob's entity (%s) as conflicting. Error: %v", entityIdBob, err)
	}
	if !strings.Contains(err.Error(), aliasIdAlice) {
		t.Fatalf("Did not identify alice's alias (%s) as conflicting. Error: %v", aliasIdAlice, err)
	}
	if !strings.Contains(err.Error(), aliasIdBob) {
		t.Fatalf("Did not identify bob's alias (%s) as conflicting. Error: %v", aliasIdBob, err)
	}
	if !strings.Contains(err.Error(), mountAccessor) {
		t.Fatalf("Did not identify mount accessor %s as being reason for conflict. Error: %v", mountAccessor, err)
	}
}

func TestIdentityStore_MergeEntities_FailsDueToClashInFromEntities(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{
		Type: "github",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	var mountAccessorGitHub string
	for k, v := range mounts {
		if k == "github/" {
			mountAccessorGitHub = v.Accessor
			break
		}
	}
	if mountAccessorGitHub == "" {
		t.Fatal("did not find github accessor")
	}

	_, entityIdBob, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")
	_, entityIdAlice, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessorGitHub, "alice-smith", "alice")
	_, entityIdClara, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessorGitHub, "clara-smith", "clara")

	// Perform entity merge
	mergeResp, err := client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    entityIdBob,
		"from_entity_ids": []string{entityIdAlice, entityIdClara},
	})
	if err == nil {
		t.Fatalf("Expected error upon merge. Resp:%#v", mergeResp)
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("mount accessor %s found in multiple fromEntities, merge should be done with one fromEntity at a time", mountAccessorGitHub)) {
		t.Fatalf("Error was not due to conflicting alias mount accessors in fromEntities. Error: %v", err)
	}
}

func TestIdentityStore_MergeEntities_FailsDueToDoubleClash(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{
		Type: "github",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob-github", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	var mountAccessorGitHub string
	for k, v := range mounts {
		if k == "github/" {
			mountAccessorGitHub = v.Accessor
			break
		}
	}
	if mountAccessorGitHub == "" {
		t.Fatal("did not find github accessor")
	}

	_, entityIdBob, aliasIdBob := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")

	aliasResp, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "bob-github",
		"canonical_id":   entityIdBob,
		"mount_accessor": mountAccessorGitHub,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, aliasResp)
	}

	aliasIdBobGitHub := aliasResp.Data["id"].(string)
	if aliasIdBobGitHub == "" {
		t.Fatal("Alias ID not present in response")
	}

	// Create userpass login for alice
	_, err = client.Logical().Write("auth/userpass/users/alice", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, entityIdAlice, aliasIdAlice := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "alice-smith", "alice")
	_, entityIdClara, aliasIdClara := testhelpers.CreateEntityAndAlias(t, client, mountAccessorGitHub, "clara-smith", "clara")

	// Perform entity merge
	mergeResp, err := client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    entityIdBob,
		"from_entity_ids": []string{entityIdAlice, entityIdClara},
	})
	if err == nil {
		t.Fatalf("Expected error upon merge. Resp:%#v", mergeResp)
	}
	if mergeResp != nil {
		t.Fatalf("Response was non-nil. Resp:%#v", mergeResp)
	}
	if !strings.Contains(err.Error(), "toEntity and at least one fromEntity have aliases with the same mount accessor") {
		t.Fatalf("Error was not due to conflicting alias mount accessors. Error: %v", err)
	}
	if !strings.Contains(err.Error(), entityIdAlice) {
		t.Fatalf("Did not identify alice's entity (%s) as conflicting. Error: %v", entityIdAlice, err)
	}
	if !strings.Contains(err.Error(), entityIdBob) {
		t.Fatalf("Did not identify bob's entity (%s) as conflicting. Error: %v", entityIdBob, err)
	}
	if !strings.Contains(err.Error(), entityIdClara) {
		t.Fatalf("Did not identify clara's alias (%s) as conflicting. Error: %v", entityIdClara, err)
	}
	if !strings.Contains(err.Error(), aliasIdAlice) {
		t.Fatalf("Did not identify alice's alias (%s) as conflicting. Error: %v", aliasIdAlice, err)
	}
	if !strings.Contains(err.Error(), aliasIdBob) {
		t.Fatalf("Did not identify bob's alias (%s) as conflicting. Error: %v", aliasIdBob, err)
	}
	if !strings.Contains(err.Error(), aliasIdClara) {
		t.Fatalf("Did not identify bob's alias (%s) as conflicting. Error: %v", aliasIdClara, err)
	}
	if !strings.Contains(err.Error(), mountAccessor) {
		t.Fatalf("Did not identify mount accessor %s as being reason for conflict. Error: %v", mountAccessor, err)
	}
	if !strings.Contains(err.Error(), mountAccessorGitHub) {
		t.Fatalf("Did not identify mount accessor %s as being reason for conflict. Error: %v", mountAccessorGitHub, err)
	}
}

func TestIdentityStore_MergeEntities_FailsDueToClashInFromEntities_CheckRawRequest(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	_, entityIdBob, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")

	// Create userpass login for alice
	_, err = client.Logical().Write("auth/userpass/users/alice", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, entityIdAlice, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "alice-smith", "alice")

	// Perform entity merge as a Raw Request so we can investigate the response body
	req := client.NewRequest("POST", "/v1/identity/entity/merge")
	req.SetJSONBody(map[string]interface{}{
		"to_entity_id":    entityIdBob,
		"from_entity_ids": []string{entityIdAlice},
	})

	resp, err := client.RawRequest(req)
	if err == nil {
		t.Fatalf("Expected error but did not get one. Response: %v", resp)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(bodyBytes)

	if resp.StatusCode != 400 {
		t.Fatal("Incorrect status code for response")
	}

	var mapOutput map[string]interface{}
	if err = json.Unmarshal([]byte(bodyString), &mapOutput); err != nil {
		t.Fatal(err)
	}

	errorStrings, ok := mapOutput["errors"].([]interface{})
	if !ok {
		t.Fatalf("error not present in response - full response: %s", bodyString)
	}

	if len(errorStrings) != 1 {
		t.Fatalf("Incorrect number of errors in response - full response: %s", bodyString)
	}

	errorString, ok := errorStrings[0].(string)
	if !ok {
		t.Fatalf("error not present in response - full response: %s", bodyString)
	}

	if !strings.Contains(errorString, "toEntity and at least one fromEntity have aliases with the same mount accessor") {
		t.Fatalf("Error was not due to conflicting alias mount accessors. Error: %s", errorString)
	}

	dataArray, ok := mapOutput["data"].([]interface{})
	if !ok {
		t.Fatalf("data not present in response - full response: %s", bodyString)
	}

	if len(dataArray) != 2 {
		t.Fatalf("Incorrect amount of clash data in response - full response: %s", bodyString)
	}

	for _, data := range dataArray {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			t.Fatalf("data could not be understood - full response: %s", bodyString)
		}

		entityId, ok := dataMap["entity_id"].(string)
		if !ok {
			t.Fatalf("entity_id not present in data - full response: %s", bodyString)
		}

		if entityId != entityIdBob && entityId != entityIdAlice {
			t.Fatalf("entityId not bob or alice - full response: %s", bodyString)
		}

		entity, ok := dataMap["entity"].(string)
		if !ok {
			t.Fatalf("entity not present in data - full response: %s", bodyString)
		}

		if entity != "bob-smith" && entity != "alice-smith" {
			t.Fatalf("entity not bob or alice - full response: %s", bodyString)
		}

		alias, ok := dataMap["alias"].(string)
		if !ok {
			t.Fatalf("alias not present in data - full response: %s", bodyString)
		}

		if alias != "bob" && alias != "alice" {
			t.Fatalf("alias not bob or alice - full response: %s", bodyString)
		}

		mountPath, ok := dataMap["mount_path"].(string)
		if !ok {
			t.Fatalf("mountPath not present in data - full response: %s", bodyString)
		}

		if mountPath != "auth/userpass/" {
			t.Fatalf("mountPath not auth/userpass/ - full response: %s", bodyString)
		}

		mount, ok := dataMap["mount"].(string)
		if !ok {
			t.Fatalf("mount not present in data - full response: %s", bodyString)
		}

		if mount != "userpass" {
			t.Fatalf("mount not userpass - full response: %s", bodyString)
		}
	}
}

func TestIdentityStore_MergeEntities_SameMountAccessor_ThenUseAlias(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("auth/userpass/login/bob", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	_, entityIdBob, aliasIdBob := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")

	// Create userpass login for alice
	_, err = client.Logical().Write("auth/userpass/users/alice", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("auth/userpass/login/alice", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, entityIdAlice, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "alice-smith", "alice")

	// Try and login with alias 2 (alice) pre-merge
	userpassAuth, err := auth.NewUserpassAuth("alice", &auth.Password{FromString: "testpassword"})
	if err != nil {
		t.Fatal(err)
	}
	loginResp, err := client.Logical().Write("auth/userpass/login/alice", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, loginResp)
	}
	if loginResp.Auth == nil {
		t.Fatalf("Request auth is nil, something has gone wrong - resp:%#v", loginResp)
	}
	loginEntityId := loginResp.Auth.EntityID
	if loginEntityId != entityIdAlice {
		t.Fatalf("Login entity ID is not Alice. loginEntityId:%s aliceEntityId:%s", loginEntityId, entityIdAlice)
	}

	// Perform entity merge
	mergeResp, err := client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":                  entityIdBob,
		"from_entity_ids":               entityIdAlice,
		"conflicting_alias_ids_to_keep": aliasIdBob,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, mergeResp)
	}

	// Delete entity id 1 (bob)
	deleteResp, err := client.Logical().Delete(fmt.Sprintf("identity/entity/id/%s", entityIdBob))
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, deleteResp)
	}

	// Try and login with alias 2 (alice) post-merge
	// Notably, this login method sets the client token, which is why we didn't use it above
	loginResp, err = client.Auth().Login(context.Background(), userpassAuth)
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, loginResp)
	}
	if loginResp.Auth == nil {
		t.Fatalf("Request auth is nil, something has gone wrong - resp:%#v", loginResp)
	}
	if loginEntityId != entityIdAlice {
		t.Fatalf("Login entity ID is not Alice. loginEntityId:%s aliceEntityId:%s", loginEntityId, entityIdAlice)
	}
}

func TestIdentityStore_MergeEntities_FailsDueToMultipleClashMergesAttempted(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{
		Type: "github",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/userpass/users/bob-github", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	mounts, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}

	var mountAccessor string
	for k, v := range mounts {
		if k == "userpass/" {
			mountAccessor = v.Accessor
			break
		}
	}
	if mountAccessor == "" {
		t.Fatal("did not find userpass accessor")
	}

	var mountAccessorGitHub string
	for k, v := range mounts {
		if k == "github/" {
			mountAccessorGitHub = v.Accessor
			break
		}
	}
	if mountAccessorGitHub == "" {
		t.Fatal("did not find github accessor")
	}

	_, entityIdBob, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "bob-smith", "bob")
	aliasResp, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "bob-github",
		"canonical_id":   entityIdBob,
		"mount_accessor": mountAccessorGitHub,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, aliasResp)
	}

	aliasIdBobGitHub := aliasResp.Data["id"].(string)
	if aliasIdBobGitHub == "" {
		t.Fatal("Alias ID not present in response")
	}

	// Create userpass login for alice
	_, err = client.Logical().Write("auth/userpass/users/alice", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, entityIdAlice, aliasIdAlice := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "alice-smith", "alice")
	_, entityIdClara, aliasIdClara := testhelpers.CreateEntityAndAlias(t, client, mountAccessorGitHub, "clara-smith", "alice")

	// Perform entity merge
	mergeResp, err := client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":                  entityIdBob,
		"from_entity_ids":               []string{entityIdAlice, entityIdClara},
		"conflicting_alias_ids_to_keep": []string{aliasIdAlice, aliasIdClara},
	})
	if err == nil {
		t.Fatalf("Expected error upon merge. Resp:%#v", mergeResp)
	}
	if !strings.Contains(err.Error(), "merge one entity at a time") {
		t.Fatalf("did not error for the right reason. Error: %v", err)
	}
}

// aliasSpec describes a alias
type aliasSpec struct {
	mountPath string // name of the mount accessor
	name      string
}

// aliasExistsByID returns true when the alias with aliasID is still present in
// the identity alias list, false when it has been deleted
func aliasExistsByID(t *testing.T, client *api.Client, aliasID string) bool {
	t.Helper()
	resp, err := client.Logical().Read("identity/entity-alias/id/" + aliasID)
	require.NoError(t, err)
	return resp != nil
}

// entityAliasIDs returns the sorted set of alias IDs currently attached to an entity
func entityAliasIDs(t *testing.T, client *api.Client, entityID string) []string {
	t.Helper()
	resp, err := client.Logical().Read("identity/entity/id/" + entityID)
	require.NoError(t, err)
	require.NotNil(t, resp, "entity %s not found", entityID)

	aliases, _ := resp.Data["aliases"].([]interface{})
	ids := make([]string, 0, len(aliases))
	for _, raw := range aliases {
		m := raw.(map[string]interface{})
		ids = append(ids, m["id"].(string))
	}
	sort.Strings(ids)
	return ids
}

// TestEntityMerge_ConflictResolution verifies that the caller can resolve alias conflicts that arise
// when merging two entities that share one or more mount accessors.
// Each test case declares the full alias topology for both entities via toAliases / fromAliases:
// any mount path that appears in both lists is a conflict and requires a resolution entry.
// The merged entity should have non-conflicting aliases and alias decided to kept,
// and shouldn't have alias that get deleted during the entityMerge.
func TestEntityMerge_ConflictResolution(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		toAliases   []aliasSpec // aliases to create on toEntity
		fromAliases []aliasSpec // aliases to create on fromEntity
		// conflictResolution receives maps of aliasName -> aliasID for each entity
		// and returns the slice of alias IDs to pass as conflicting_alias_ids_to_keep.
		conflictResolution func(toAliasIDs, fromAliasIDs map[string]string) []string
	}{
		{
			name:        "one conflict mount: keep toEntity alias",
			toAliases:   []aliasSpec{{mountPath: "jwt1", name: "to-alias"}},
			fromAliases: []aliasSpec{{mountPath: "jwt1", name: "from-alias"}},
			conflictResolution: func(toAliasIDs, _ map[string]string) []string {
				return []string{toAliasIDs["to-alias"]}
			},
		},
		{
			name:        "one conflict mount: keep fromEntity alias",
			toAliases:   []aliasSpec{{mountPath: "jwt1", name: "to-alias"}},
			fromAliases: []aliasSpec{{mountPath: "jwt1", name: "from-alias"}},
			conflictResolution: func(_, fromAliasIDs map[string]string) []string {
				return []string{fromAliasIDs["from-alias"]}
			},
		},
		{
			// two conflicting mounts: keep one alias from each entity
			name: "two conflict mounts: keep one alias from each entity",
			toAliases: []aliasSpec{
				{mountPath: "jwt1", name: "to-alias-jwt1"},
				{mountPath: "jwt2", name: "to-alias-jwt2"},
			},
			fromAliases: []aliasSpec{
				{mountPath: "jwt1", name: "from-alias-jwt1"},
				{mountPath: "jwt2", name: "from-alias-jwt2"},
			},
			conflictResolution: func(toAliasIDs, fromAliasIDs map[string]string) []string {
				return []string{toAliasIDs["to-alias-jwt1"], fromAliasIDs["from-alias-jwt2"]}
			},
		},
		{
			// one conflicting mount and one non conflicting alias on toEntity only
			// The non conflicting and kept alias should survive
			name: "one conflict plus one non-conflicting toEntity alias",
			toAliases: []aliasSpec{
				{mountPath: "jwt1", name: "to-alias-conflict"},
				{mountPath: "jwt2", name: "to-alias-unique"},
			},
			fromAliases: []aliasSpec{
				{mountPath: "jwt1", name: "from-alias-conflict"},
			},
			conflictResolution: func(toAliasIDs, _ map[string]string) []string {
				return []string{toAliasIDs["to-alias-conflict"]}
			},
		},
		{
			// One conflicting mount and one non conflicting alias on fromEntity only
			// The non conflicting alias must be transferred to toEntity
			name: "one conflict plus one non-conflicting fromEntity alias",
			toAliases: []aliasSpec{
				{mountPath: "jwt1", name: "to-alias-conflict"},
			},
			fromAliases: []aliasSpec{
				{mountPath: "jwt1", name: "from-alias-conflict"},
				{mountPath: "jwt2", name: "from-alias-unique"},
			},
			conflictResolution: func(_, fromAliasIDs map[string]string) []string {
				return []string{fromAliasIDs["from-alias-conflict"]}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client

			// Get and enable mount paths across both entities
			seenMounts := make(map[string]bool)
			for _, a := range append(tc.toAliases, tc.fromAliases...) {
				if !seenMounts[a.mountPath] {
					seenMounts[a.mountPath] = true
					require.NoError(t, client.Sys().EnableAuthWithOptions(a.mountPath, &api.EnableAuthOptions{
						Type: "jwt",
					}))
				}
			}

			authMounts, err := client.Sys().ListAuth()
			require.NoError(t, err)
			accessorOf := make(map[string]string, len(seenMounts))
			for path := range seenMounts {
				accessorOf[path] = authMounts[path+"/"].Accessor
			}

			// Create toEntity and all its aliases and map aliasName -> aliasID
			toResp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "to-entity"})
			require.NoError(t, err)
			toEntityID := toResp.Data["id"].(string)

			toAliasIDs := make(map[string]string, len(tc.toAliases))
			for _, a := range tc.toAliases {
				r, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
					"name":           a.name,
					"mount_accessor": accessorOf[a.mountPath],
					"canonical_id":   toEntityID,
				})
				require.NoError(t, err)
				toAliasIDs[a.name] = r.Data["id"].(string)
			}

			// Create fromEntity and all its aliases and map aliasName -> aliasID
			fromResp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "from-entity"})
			require.NoError(t, err)
			fromEntityID := fromResp.Data["id"].(string)

			fromAliasIDs := make(map[string]string, len(tc.fromAliases))
			for _, a := range tc.fromAliases {
				r, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
					"name":           a.name,
					"mount_accessor": accessorOf[a.mountPath],
					"canonical_id":   fromEntityID,
				})
				require.NoError(t, err)
				fromAliasIDs[a.name] = r.Data["id"].(string)
			}

			// Determine which aliases are kept and which are discarded
			keptAliasIDs := tc.conflictResolution(toAliasIDs, fromAliasIDs)
			keptSet := make(map[string]bool, len(keptAliasIDs))
			for _, id := range keptAliasIDs {
				keptSet[id] = true
			}

			// Perform the merge with conflict resolution
			_, err = client.Logical().Write("identity/entity/merge", map[string]interface{}{
				"to_entity_id":                  toEntityID,
				"from_entity_ids":               []string{fromEntityID},
				"conflicting_alias_ids_to_keep": keptAliasIDs,
			})
			require.NoError(t, err, "merge with conflict resolution should succeed")

			// fromEntity should be deleted after the merge
			resp, err := client.Logical().Read("identity/entity/id/" + fromEntityID)
			require.NoError(t, err)
			require.Nil(t, resp, "fromEntity should have been deleted after merge")

			// Determine conflicting mount paths
			toMountPaths := make(map[string]bool, len(tc.toAliases))
			for _, a := range tc.toAliases {
				toMountPaths[a.mountPath] = true
			}
			conflictingMounts := make(map[string]bool)
			for _, a := range tc.fromAliases {
				if toMountPaths[a.mountPath] {
					conflictingMounts[a.mountPath] = true
				}
			}

			// Build the complete set of alias IDs that survive on toEntity after the merge:
			// all non conflicting aliases and keptAliasIDs from both entities
			expectedSurvivors := make(map[string]bool)
			for _, a := range tc.toAliases {
				id := toAliasIDs[a.name]
				if conflictingMounts[a.mountPath] {
					if keptSet[id] {
						expectedSurvivors[id] = true
					}
				} else {
					expectedSurvivors[id] = true
				}
			}
			for _, a := range tc.fromAliases {
				id := fromAliasIDs[a.name]
				if conflictingMounts[a.mountPath] {
					if keptSet[id] {
						expectedSurvivors[id] = true
					}
				} else {
					expectedSurvivors[id] = true
				}
			}

			// Every alias that should survive should exist and be attached to toEntity.
			for id := range expectedSurvivors {
				require.True(t, aliasExistsByID(t, client, id),
					"alias %s should still exist after merge", id)
			}

			// Every discarded alias should be deleted
			allAliasIDs := make(map[string]struct{})
			for _, id := range toAliasIDs {
				allAliasIDs[id] = struct{}{}
			}
			for _, id := range fromAliasIDs {
				allAliasIDs[id] = struct{}{}
			}
			for id := range allAliasIDs {
				if !expectedSurvivors[id] {
					require.False(t, aliasExistsByID(t, client, id),
						"alias %s should have been deleted after merge", id)
				}
			}

			// toEntity should only contain the surviving aliases: non conflicting and
			// explicitly kept via conflicting_alias_ids_to_keep
			wantIDs := make([]string, 0, len(expectedSurvivors))
			for id := range expectedSurvivors {
				wantIDs = append(wantIDs, id)
			}
			sort.Strings(wantIDs)

			gotIDs := entityAliasIDs(t, client, toEntityID)
			require.Equal(t, wantIDs, gotIDs,
				"toEntity should carry exactly the expected aliases after merge")
		})
	}
}

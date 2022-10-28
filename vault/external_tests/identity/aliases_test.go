package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestIdentityStore_ListAlias(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"github": github.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
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
		t.Logf("key: %v\nmount: %#v", k, *v)
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
		t.Logf("alias info: %#v", info)
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

	entityInfoRaw, ok := resp.Data["key_info"]
	if !ok {
		t.Fatal("expected key_info map in response")
	}

	// This is basically verifying that the entity has the alias in key_info
	// that we expect to be tied to it, plus tests a value further down in it
	// for fun
	entityInfo := entityInfoRaw.(map[string]interface{})
	for _, keyRaw := range keys {
		key := keyRaw.(string)
		infoRaw, ok := entityInfo[key]
		if !ok {
			t.Fatal("expected key info")
		}
		info := infoRaw.(map[string]interface{})
		t.Logf("entity info: %#v", info)
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
			}
		}
	}
}

// TestIdentityStore_RenameAlias_CannotMergeEntity verifies that an error is
// returned on an attempt to rename an alias to match another alias with the
// same mount accessor.  This used to result in a merge entity.
func TestIdentityStore_RenameAlias_CannotMergeEntity(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
			"github":   github.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
			"github":   github.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
			"github":   github.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

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

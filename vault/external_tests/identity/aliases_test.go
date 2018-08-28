package identity

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
)

func TestIdentityStore_EntityAliasLocalMount(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": credLdap.Factory,
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

	// Create a local auth mount
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type:  "ldap",
		Local: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Extract out the mount accessor for LDAP auth
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor := auths["ldap/"].Accessor

	// Create an entity
	secret, err := client.Logical().Write("identity/entity", nil)
	if err != nil {
		t.Fatal(err)
	}
	entityID := secret.Data["id"].(string)

	// Attempt to create an entity alias against a local mount should fail
	secret, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser",
		"mount_accessor": ldapMountAccessor,
		"canonical_id":   entityID,
	})
	if err == nil {
		t.Fatalf("expected error since mount is local")
	}
}

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

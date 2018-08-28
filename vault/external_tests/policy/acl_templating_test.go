package policy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPolicyTemplating(t *testing.T) {
	goodPolicy1 := `
path "secret/{{ identity.entity.name}}/*" {
	capabilities = ["read", "create", "update"]

}

path "secret/{{ identity.entity.aliases.%s.name}}/*" {
	capabilities = ["read", "create", "update"]

}
`

	goodPolicy2 := `
path "secret/{{ identity.groups.ids.%s.name}}/*" {
	capabilities = ["read", "create", "update"]

}

path "secret/{{ identity.groups.names.group_name.id}}/*" {
	capabilities = ["read", "create", "update"]

}
`

	badPolicy1 := `
path "secret/{{ identity.groups.names.foobar.name}}/*" {
	capabilities = ["read", "create", "update"]

}
`

	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
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

	resp, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "entity_name",
		"policies": []string{
			"goodPolicy1",
			"badPolicy1",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	entityID := resp.Data["id"].(string)

	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"policies": []string{
			"goodPolicy2",
		},
		"member_entity_ids": []string{
			entityID,
		},
		"name": "group_name",
	})
	if err != nil {
		t.Fatal(err)
	}
	groupID := resp.Data["id"]

	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"name": "foobar",
	})
	if err != nil {
		t.Fatal(err)
	}
	foobarGroupID := resp.Data["id"]

	// Enable userpass auth
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create an external group and renew the token. This should add external
	// group policies to the token.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	userpassAccessor := auths["userpass/"].Accessor

	// Create an alias
	resp, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser",
		"mount_accessor": userpassAccessor,
		"canonical_id":   entityID,
	})
	if err != nil {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Add a user to userpass backend
	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write in policies
	goodPolicy1 = fmt.Sprintf(goodPolicy1, userpassAccessor)
	goodPolicy2 = fmt.Sprintf(goodPolicy2, groupID)
	err = client.Sys().PutPolicy("goodPolicy1", goodPolicy1)
	if err != nil {
		t.Fatal(err)
	}
	err = client.Sys().PutPolicy("goodPolicy2", goodPolicy2)
	if err != nil {
		t.Fatal(err)
	}

	// Authenticate
	secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	clientToken := secret.Auth.ClientToken

	var tests = []struct {
		name string
		path string
		fail bool
	}{
		{
			name: "entity name",
			path: "secret/entity_name/foo",
		},
		{
			name: "bad entity name",
			path: "secret/entityname/foo",
			fail: true,
		},
		{
			name: "group name",
			path: "secret/group_name/foo",
		},
		{
			name: "group id",
			path: fmt.Sprintf("secret/%s/foo", groupID),
		},
		{
			name: "alias name",
			path: "secret/testuser/foo",
		},
		{
			name: "bad group name",
			path: "secret/foobar/foo",
		},
	}

	runTests := func(failGroupName bool) {
		for _, test := range tests {
			resp, err := client.Logical().Write(test.path, map[string]interface{}{"zip": "zap"})
			fail := test.fail
			if test.name == "bad group name" {
				fail = failGroupName
			}
			if err != nil && !fail {
				if resp.Data["error"].(string) != "permission denied" {
					t.Fatalf("unexpected status %v", resp.Data["error"])
				}
				t.Fatalf("%s: got unexpected error: %v", test.name, err)
			}
			if err == nil && fail {
				t.Fatalf("%s: expected error", test.name)
			}
		}
	}

	rootToken := client.Token()
	client.SetToken(clientToken)
	runTests(true)

	client.SetToken(rootToken)
	// Test that a policy with bad group membership doesn't kill the other paths
	err = client.Sys().PutPolicy("badPolicy1", badPolicy1)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(clientToken)
	runTests(true)

	// Test that adding group membership now allows access
	client.SetToken(rootToken)
	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"id": foobarGroupID,
		"member_entity_ids": []string{
			entityID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(clientToken)
	runTests(false)
}

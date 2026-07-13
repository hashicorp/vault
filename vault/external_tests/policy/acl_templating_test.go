// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package policy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
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

	cluster := minimal.NewTestSoloCluster(t, nil)
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

	tests := []struct {
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

// TestPolicyTemplating_DenySlashInTemplatedPaths exercises the
// deny_slash_in_templated_path config option
func TestPolicyTemplating_DenySlashInTemplatedPaths(t *testing.T) {
	tests := []struct {
		name                string
		denySlashEnabled    bool
		customMetadataValue string
		expectAccessDenied  bool
	}{
		{
			name:                "deny_slash_enabled_with_slashes_in_metadata",
			denySlashEnabled:    true,
			customMetadataValue: "path/with/slashes",
			expectAccessDenied:  true,
		},
		{
			name:                "deny_slash_disabled_with_slashes_in_metadata",
			denySlashEnabled:    false,
			customMetadataValue: "path/with/slashes",
			expectAccessDenied:  false,
		},
		{
			name:                "deny_slash_enabled_without_slashes_in_metadata",
			denySlashEnabled:    true,
			customMetadataValue: "pathwithoutslashes",
			expectAccessDenied:  false,
		},
		{
			name:                "deny_slash_disabled_without_slashes_in_metadata",
			denySlashEnabled:    false,
			customMetadataValue: "pathwithoutslashes",
			expectAccessDenied:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
				DenySlashInTemplatedPolicyPaths: tc.denySlashEnabled,
			})
			client := cluster.Cores[0].Client

			// Create entity
			resp, err := client.Logical().Write("identity/entity", map[string]interface{}{
				"name": "test_entity",
			})
			require.NoError(t, err)
			entityID := resp.Data["id"].(string)

			// Enable userpass auth
			err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
				Type: "userpass",
			})
			require.NoError(t, err)

			auths, err := client.Sys().ListAuth()
			require.NoError(t, err)
			userpassAccessor := auths["userpass/"].Accessor

			// Create an alias with custom metadata containing the test value
			resp, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
				"name":            "testuser",
				"mount_accessor":  userpassAccessor,
				"canonical_id":    entityID,
				"custom_metadata": map[string]string{"test_path": tc.customMetadataValue},
			})
			require.NoError(t, err)

			// Add a user to userpass backend
			_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
				"password": "testpassword",
			})
			require.NoError(t, err)

			// Create a policy that uses the alias custom metadata in the path
			policy := fmt.Sprintf(`
path "secret/{{identity.entity.aliases.%s.custom_metadata.test_path}}/*" {
	capabilities = ["read", "create", "update"]
}
`, userpassAccessor)

			err = client.Sys().PutPolicy("testPolicy", policy)
			require.NoError(t, err)

			// Update entity to use the policy
			_, err = client.Logical().Write("identity/entity/id/"+entityID, map[string]interface{}{
				"policies": []string{"testPolicy"},
			})
			require.NoError(t, err)

			// Authenticate as the user
			secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
				"password": "testpassword",
			})
			require.NoError(t, err)
			clientToken := secret.Auth.ClientToken

			// Try to access a path using the templated custom metadata value
			client.SetToken(clientToken)

			testPath := fmt.Sprintf("secret/%s/data", tc.customMetadataValue)
			_, err = client.Logical().Write(testPath, map[string]interface{}{"key": "value"})

			if tc.expectAccessDenied {
				if err == nil {
					t.Fatalf("expected access denied when deny_slash=%v and custom_metadata=%q, but got success",
						tc.denySlashEnabled, tc.customMetadataValue)
				}
			} else {
				if err != nil {
					t.Fatalf("expected success when deny_slash=%v and custom_metadata=%q, but got error: %v",
						tc.denySlashEnabled, tc.customMetadataValue, err)
				}
			}
		})
	}
}

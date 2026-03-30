// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package secrets

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// testIdentitySecretsCreate tests Identity secrets engine creation
func testIdentitySecretsCreate(t *testing.T, v *blackbox.Session) {
	// Identity is a built-in secrets engine, no need to enable it
	// Create an entity
	entityResp := v.MustWrite("identity/entity", map[string]any{
		"name":     "test-entity",
		"policies": []string{"default"},
		"metadata": map[string]string{
			"team": "engineering",
			"env":  "test",
		},
	})

	if entityResp.Data == nil {
		t.Fatal("Expected entity creation response")
	}

	// Verify entity was created
	assertions := v.AssertSecret(entityResp)
	assertions.Data().HasKeyExists("id")

	entityID := entityResp.Data["id"].(string)

	// Create an entity alias (requires an auth mount)
	v.MustEnableAuth("userpass-identity", &api.EnableAuthOptions{Type: "userpass"})

	// Get the accessor for the auth mount
	authsResp := v.MustRead("sys/auth")
	if authsResp.Data == nil {
		t.Fatal("Expected to read auth mounts")
	}

	var accessor string
	if authData, ok := authsResp.Data["userpass-identity/"]; ok {
		if authMap, ok := authData.(map[string]any); ok {
			accessor = authMap["accessor"].(string)
		}
	}

	if accessor == "" {
		t.Fatal("Failed to get auth mount accessor")
	}

	// Create entity alias
	aliasResp := v.MustWrite("identity/entity-alias", map[string]any{
		"name":           "test-user",
		"canonical_id":   entityID,
		"mount_accessor": accessor,
	})

	if aliasResp.Data == nil {
		t.Fatal("Expected entity alias creation response")
	}

	aliasAssertions := v.AssertSecret(aliasResp)
	aliasAssertions.Data().HasKeyExists("id")

	t.Log("Successfully created identity entity and alias")
}

// testIdentitySecretsRead tests Identity secrets engine read operations
func testIdentitySecretsRead(t *testing.T, v *blackbox.Session) {
	// Create an entity
	entityResp := v.MustWrite("identity/entity", map[string]any{
		"name":     "read-entity",
		"policies": []string{"default"},
		"metadata": map[string]string{
			"purpose": "testing",
		},
	})

	entityID := entityResp.Data["id"].(string)

	// Read the entity by ID
	readResp := v.MustRead("identity/entity/id/" + entityID)
	if readResp.Data == nil {
		t.Fatal("Expected to read entity")
	}

	// Verify entity properties
	assertions := v.AssertSecret(readResp)
	assertions.Data().
		HasKey("name", "read-entity").
		HasKeyExists("id").
		HasKeyExists("policies")

	// Read entity by name
	nameResp := v.MustRead("identity/entity/name/read-entity")
	if nameResp.Data == nil {
		t.Fatal("Expected to read entity by name")
	}

	nameAssertions := v.AssertSecret(nameResp)
	nameAssertions.Data().
		HasKey("name", "read-entity").
		HasKey("id", entityID)

	t.Log("Successfully read identity entity by ID and name")
}

// testIdentitySecretsDelete tests Identity secrets engine delete operations
func testIdentitySecretsDelete(t *testing.T, v *blackbox.Session) {
	entityResp := v.MustWrite("identity/entity", map[string]any{
		"name": "delete-entity",
	})
	entityID := entityResp.Data["id"].(string)
	readResp := v.MustRead("identity/entity/id/" + entityID)
	if readResp.Data == nil {
		t.Fatal("Expected entity to exist before deletion")
	}
	_, err := v.Client.Logical().Delete("identity/entity/id/" + entityID)
	if err != nil {
		t.Fatalf("Failed to delete entity: %v", err)
	}
	deletedResp, err := v.Client.Logical().Read("identity/entity/id/" + entityID)
	if err == nil && deletedResp != nil {
		t.Fatal("Expected entity to be deleted, but it still exists")
	}
	t.Logf("Successfully deleted identity entity: %s", entityID)
}

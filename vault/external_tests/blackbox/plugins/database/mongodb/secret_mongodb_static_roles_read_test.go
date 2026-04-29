// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

const (
	readTestUsername = "readtestuser"
	readTestRoleName = "read-test-role"
)

// TestMongoDBStaticRoleReadWorkflows runs all MongoDB static role read workflow tests.
func TestMongoDBStaticRoleReadWorkflows(t *testing.T) {
	t.Run("ReadReturnsConfiguration", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleReadReturnsConfiguration(t, v)
	})

	t.Run("ReadNonExistentRole", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleReadNonExistentRole(t, v)
	})

	t.Run("ReadAfterUpdate", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleReadAfterUpdate(t, v)
	})

	t.Run("ListMultipleRoles", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleListMultipleRoles(t, v)
	})
}

// testMongoDBStaticRoleReadReturnsConfiguration verifies reading a static role
// returns its configuration without sensitive data.
func testMongoDBStaticRoleReadReturnsConfiguration(t *testing.T, v *blackbox.Session) {
	mount, connURL := setupMongoDBTest(t, v)

	createMongoDBUser(t, connURL, readTestUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+readTestRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        readTestUsername,
		"rotation_period": testRotationPeriod,
	})

	role := v.MustReadRequired(mount + "/static-roles/" + readTestRoleName)
	v.AssertSecret(role).
		Data().
		HasKey("username", readTestUsername).
		HasKey("rotation_period", float64(testRotationPeriod)).
		HasKey("db_name", testConnectionName)

	// Verify password is not returned in role configuration
	if _, ok := role.Data["password"]; ok {
		t.Fatal("password should not be returned in role read")
	}

	// Verify last_vault_rotation is present
	v.AssertSecret(role).Data().HasKey("last_vault_rotation")
}

// testMongoDBStaticRoleReadNonExistentRole verifies reading a non-existent
// static role returns an appropriate error.
func testMongoDBStaticRoleReadNonExistentRole(t *testing.T, v *blackbox.Session) {
	mount, _ := setupMongoDBTest(t, v)

	_, err := v.Read(mount + "/static-roles/nonexistent-role")
	if err == nil {
		t.Fatal("expected error when reading non-existent role")
	}
}

// testMongoDBStaticRoleReadAfterUpdate verifies reading a static role after
// updating it returns the updated configuration.
func testMongoDBStaticRoleReadAfterUpdate(t *testing.T, v *blackbox.Session) {
	mount, connURL := setupMongoDBTest(t, v)

	createMongoDBUser(t, connURL, readTestUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+readTestRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        readTestUsername,
		"rotation_period": testRotationPeriod,
	})

	role1 := v.MustReadRequired(mount + "/static-roles/" + readTestRoleName)
	v.AssertSecret(role1).Data().HasKey("rotation_period", float64(testRotationPeriod))

	// Update rotation period
	newRotationPeriod := 7200
	v.MustWrite(mount+"/static-roles/"+readTestRoleName, map[string]any{
		"rotation_period": newRotationPeriod,
	})

	role2 := v.MustReadRequired(mount + "/static-roles/" + readTestRoleName)
	v.AssertSecret(role2).Data().HasKey("rotation_period", float64(newRotationPeriod))

	// Verify username is still the same
	v.AssertSecret(role2).Data().HasKey("username", readTestUsername)
}

// testMongoDBStaticRoleListMultipleRoles verifies LIST /static-roles
// returns all configured role names.
func testMongoDBStaticRoleListMultipleRoles(t *testing.T, v *blackbox.Session) {
	mount, connURL := setupMongoDBTest(t, v)

	// Create multiple static roles
	roles := []string{"read-role-1", "read-role-2", "read-role-3"}
	for i, roleName := range roles {
		username := fmt.Sprintf("listuser%d", i+1)
		createMongoDBUser(t, connURL, username, testInitialPassword)

		v.MustWrite(mount+"/static-roles/"+roleName, map[string]any{
			"db_name":         testConnectionName,
			"username":        username,
			"rotation_period": testRotationPeriod,
		})
	}

	list := v.MustList(mount + "/static-roles")
	v.AssertSecret(list).Data().HasKey("keys")

	keys := list.Data["keys"].([]interface{})
	if len(keys) != 3 {
		t.Fatalf("expected 3 roles, got %d", len(keys))
	}

	// Verify all role names are present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k.(string)] = true
	}
	for _, roleName := range roles {
		if !keyMap[roleName] {
			t.Fatalf("expected role %s in list", roleName)
		}
	}

	// Verify list doesn't contain sensitive data
	for _, k := range keys {
		roleName := k.(string)
		if roleName == "password" || roleName == "connection_url" {
			t.Fatalf("list should not contain sensitive field: %s", roleName)
		}
	}
}

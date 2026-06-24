//go:build isolated
// +build isolated

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

const (
	deleteTestUsername = "deletetestuser"
	deleteTestRoleName = "delete-test-role"
)

// TestMongoDBStaticRoleDeleteWorkflows runs all MongoDB static role delete workflow tests.
func TestMongoDBStaticRoleDeleteWorkflows(t *testing.T) {
	t.Run("DeleteExistingRole", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleDeleteExistingRole(t, v)
	})

	t.Run("DeletePreventsCredentialAccess", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleDeletePreventsCredentialAccess(t, v)
	})

	t.Run("DeleteIsIdempotent", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleDeleteIsIdempotent(t, v)
	})

	t.Run("DeleteNonExistentRole", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleDeleteNonExistentRole(t, v)
	})
}

// testMongoDBStaticRoleDeleteExistingRole verifies deleting an existing
// static role succeeds and removes it from the list.
func testMongoDBStaticRoleDeleteExistingRole(t *testing.T, v *blackbox.Session) {
	mount, _, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, deleteTestUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+deleteTestRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        deleteTestUsername,
		"rotation_period": testRotationPeriod,
	})

	// Verify role exists
	role := v.MustReadRequired(mount + "/static-roles/" + deleteTestRoleName)
	v.AssertSecret(role).Data().HasKey("username", deleteTestUsername)

	// Delete the role
	v.MustDelete(mount + "/static-roles/" + deleteTestRoleName)

	// Verify role is gone - Read returns nil secret for non-existent paths
	secret, err := v.Read(mount + "/static-roles/" + deleteTestRoleName)
	if err != nil {
		// Error is acceptable
	} else if secret != nil {
		t.Fatal("expected nil secret when reading deleted role")
	}

	// Verify role is not in list
	list := v.MustList(mount + "/static-roles")
	if list == nil || list.Data == nil || list.Data["keys"] == nil {
		return
	}

	keys := list.Data["keys"].([]interface{})
	for _, k := range keys {
		if k.(string) == deleteTestRoleName {
			t.Fatalf("deleted role %s should not appear in list", deleteTestRoleName)
		}
	}
}

// testMongoDBStaticRoleDeletePreventsCredentialAccess verifies that after
// deleting a static role, credentials can no longer be read.
func testMongoDBStaticRoleDeletePreventsCredentialAccess(t *testing.T, v *blackbox.Session) {
	mount, _, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, deleteTestUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+deleteTestRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        deleteTestUsername,
		"rotation_period": testRotationPeriod,
	})

	// Verify credentials can be read before deletion
	creds := v.MustReadRequired(mount + "/static-creds/" + deleteTestRoleName)
	v.AssertSecret(creds).Data().HasKeyExists("username")
	v.AssertSecret(creds).Data().HasKeyExists("password")

	// Delete the role
	v.MustDelete(mount + "/static-roles/" + deleteTestRoleName)

	// Verify credentials can no longer be read - Read returns nil secret for deleted roles
	secret, err := v.Read(mount + "/static-creds/" + deleteTestRoleName)
	if err != nil {
		// Error is acceptable
	} else if secret != nil {
		t.Fatal("expected nil secret when reading credentials for deleted role")
	}
}

// testMongoDBStaticRoleDeleteIsIdempotent verifies that deleting a role
// multiple times succeeds without error.
func testMongoDBStaticRoleDeleteIsIdempotent(t *testing.T, v *blackbox.Session) {
	mount, _, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, deleteTestUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+deleteTestRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        deleteTestUsername,
		"rotation_period": testRotationPeriod,
	})

	// Delete the role first time
	v.MustDelete(mount + "/static-roles/" + deleteTestRoleName)

	// Verify role is gone - Read returns nil secret for non-existent paths
	secret, err := v.Read(mount + "/static-roles/" + deleteTestRoleName)
	if err != nil {
		// Error is acceptable
	} else if secret != nil {
		t.Fatal("expected nil secret when reading deleted role")
	}

	// Delete the role again - should succeed (idempotent)
	v.MustDelete(mount + "/static-roles/" + deleteTestRoleName)

	// Delete a third time to ensure consistent behavior
	v.MustDelete(mount + "/static-roles/" + deleteTestRoleName)
}

// testMongoDBStaticRoleDeleteNonExistentRole verifies that attempting to
// delete a non-existent role succeeds (idempotent behavior).
func testMongoDBStaticRoleDeleteNonExistentRole(t *testing.T, v *blackbox.Session) {
	mount, _, _, _ := setupMongoDBTest(t, v)

	// Attempt to delete a role that was never created
	v.MustDelete(mount + "/static-roles/never-existed-role")

	// Verify it's still not there - Read returns nil secret for non-existent paths
	secret, err := v.Read(mount + "/static-roles/never-existed-role")
	if err != nil {
		// Error is acceptable
		return
	}
	if secret != nil {
		t.Fatal("expected nil secret when reading non-existent role")
	}
}

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

const (
	testStaticRoleName  = "my-static-role"
	testUsername        = "staticuser1"
	mongoConnectTimeout = 10 * time.Second
)

// TestMongoDBStaticRoleWorkflows runs all MongoDB static role workflow tests.
func TestMongoDBStaticRoleWorkflows(t *testing.T) {
	t.Parallel()

	t.Run("CreateBasic", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleCreateBasic(t, v)
	})

	t.Run("ReadCredentials", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleReadCredentials(t, v)
	})

	t.Run("ManualRotation", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBStaticRoleManualRotation(t, v)
	})
}

// TODO: TestMongoDBStaticRoleValidationWorkflows
// Future implementation should test:
// - RequiresUsername: Verify error when username is missing
// - RequiresRotationPeriodOrSchedule: Verify error without rotation config
// - RejectsInvalidRotationPeriod: Verify error with period < 5 seconds
// - RejectsMutuallyExclusiveFields: Verify error with both period and schedule

// testMongoDBStaticRoleCreateBasic verifies that creating a basic static role
// succeeds with required fields.
func testMongoDBStaticRoleCreateBasic(t *testing.T, v *blackbox.Session) {
	mount, _, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, testUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+testStaticRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        testUsername,
		"rotation_period": testRotationPeriod,
	})

	role := v.MustReadRequired(mount + "/static-roles/" + testStaticRoleName)
	v.AssertSecret(role).
		Data().
		HasKey("username", testUsername).
		HasKey("rotation_period", float64(testRotationPeriod)).
		HasKey("db_name", testConnectionName)
}

// TODO: Additional test implementations for future teams
// - testMongoDBStaticRoleCreateWithRotationSchedule: Test rotation_schedule instead of rotation_period
// - testMongoDBStaticRoleListReturnsNamesOnly: Test listing multiple static roles
// - testMongoDBStaticRoleReadReturnsConfiguration: Test reading role config without sensitive data
// - testMongoDBStaticRoleUpdateRotationPeriod: Test updating rotation period
// - testMongoDBStaticRoleDeleteRole: Test deleting a static role

// testMongoDBStaticRoleManualRotation verifies manually rotating a static
// role's credentials succeeds.
func testMongoDBStaticRoleManualRotation(t *testing.T, v *blackbox.Session) {
	mount, credVerifyURL, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, testUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+testStaticRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        testUsername,
		"rotation_period": testRotationPeriod,
	})

	creds1 := v.MustReadRequired(mount + "/static-creds/" + testStaticRoleName)
	password1 := creds1.Data["password"].(string)

	if password1 == testInitialPassword {
		t.Fatal("expected password to be rotated on role creation")
	}

	v.MustWrite(mount+"/rotate-role/"+testStaticRoleName, nil)

	creds2 := v.MustReadRequired(mount + "/static-creds/" + testStaticRoleName)
	password2 := creds2.Data["password"].(string)

	if password1 == password2 {
		t.Fatal("expected password to change after rotation")
	}

	verifyMongoDBCredentials(t, credVerifyURL, testUsername, password2)
}

// TODO: testMongoDBStaticRoleAutomaticRotation - Test automatic rotation with short period
// This test requires waiting for rotation to occur, which may be time-consuming
// Pattern: Create role with short rotation_period, wait, verify password changed

// testMongoDBStaticRoleReadCredentials verifies reading static credentials
// returns the current password.
func testMongoDBStaticRoleReadCredentials(t *testing.T, v *blackbox.Session) {
	mount, credVerifyURL, dbName, client := setupMongoDBTest(t, v)

	createMongoDBUser(t, client, dbName, testUsername, testInitialPassword)

	v.MustWrite(mount+"/static-roles/"+testStaticRoleName, map[string]any{
		"db_name":         testConnectionName,
		"username":        testUsername,
		"rotation_period": testRotationPeriod,
	})

	creds := v.MustReadRequired(mount + "/static-creds/" + testStaticRoleName)
	v.AssertSecret(creds).
		Data().
		HasKey("username", testUsername).
		HasKey("password", creds.Data["password"])

	password := creds.Data["password"].(string)
	if password == "" {
		t.Fatal("expected non-empty password")
	}

	verifyMongoDBCredentials(t, credVerifyURL, testUsername, password)
}

// testMongoDBStaticRoleRequiresUsername verifies creating a static role
// without username fails.
func testMongoDBStaticRoleRequiresUsername(t *testing.T, v *blackbox.Session) {
	mount, _, _, _ := setupMongoDBTest(t, v)

	_, err := v.Client.Logical().Write(mount+"/static-roles/"+testStaticRoleName, map[string]any{
		"db_name":         testConnectionName,
		"rotation_period": testRotationPeriod,
	})

	if err == nil {
		t.Fatal("expected error when creating static role without username")
	}

	if !strings.Contains(err.Error(), "username") {
		t.Fatalf("expected error message to mention 'username', got: %v", err)
	}
}

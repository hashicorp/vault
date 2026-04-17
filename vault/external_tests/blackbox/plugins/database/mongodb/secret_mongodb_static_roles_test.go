// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testConnectionName  = "my-mongodb-db"
	testStaticRoleName  = "my-static-role"
	testUsername        = "staticuser1"
	testInitialPassword = "initialpass"
	testRotationPeriod  = 86400 // 24 hours in seconds
	mongoConnectTimeout = 10 * time.Second
)

// TestMongoDBStaticRoleWorkflows runs all MongoDB static role workflow tests.
func TestMongoDBStaticRoleWorkflows(t *testing.T) {
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
	mount, connURL := setupMongoDBTest(t, v)

	createMongoDBUser(t, connURL, testUsername, testInitialPassword)

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
	mount, connURL := setupMongoDBTest(t, v)

	createMongoDBUser(t, connURL, testUsername, testInitialPassword)

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

	verifyMongoDBCredentials(t, connURL, testUsername, password2)
}

// TODO: testMongoDBStaticRoleAutomaticRotation - Test automatic rotation with short period
// This test requires waiting for rotation to occur, which may be time-consuming
// Pattern: Create role with short rotation_period, wait, verify password changed

// testMongoDBStaticRoleReadCredentials verifies reading static credentials
// returns the current password.
func testMongoDBStaticRoleReadCredentials(t *testing.T, v *blackbox.Session) {
	mount, connURL := setupMongoDBTest(t, v)

	createMongoDBUser(t, connURL, testUsername, testInitialPassword)

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

	verifyMongoDBCredentials(t, connURL, testUsername, password)
}

// testMongoDBStaticRoleRequiresUsername verifies creating a static role
// without username fails.
func testMongoDBStaticRoleRequiresUsername(t *testing.T, v *blackbox.Session) {
	mount, _ := setupMongoDBTest(t, v)

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

// setupMongoDBTest performs common test setup: creates container, enables mount, configures connection.
// Returns mount path and connection URL.
func setupMongoDBTest(t *testing.T, v *blackbox.Session) (string, string) {
	t.Helper()

	requireVaultEnv(t)
	cleanup, connURL := PrepareTestContainer(t)
	t.Cleanup(cleanup)

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	v.MustWrite(
		mount+"/config/"+testConnectionName,
		mongoConnectionConfigPayload(connURL, "*", false),
	)

	return mount, connURL
}

// createMongoDBUser creates a MongoDB user for testing static roles.
func createMongoDBUser(t *testing.T, connURL, username, password string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("admin")
	err = db.RunCommand(ctx, bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password},
		{Key: "roles", Value: bson.A{
			bson.D{
				{Key: "role", Value: "readWrite"},
				{Key: "db", Value: "admin"},
			},
		}},
	}).Err()
	if err != nil {
		t.Fatalf("failed to create MongoDB user: %v", err)
	}

	t.Logf("Created MongoDB user: %s", username)
}

// verifyMongoDBCredentials verifies that the given credentials work for MongoDB.
func verifyMongoDBCredentials(t *testing.T, connURL, username, password string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	// Replace credentials in connection URL
	u, err := parseMongoURL(connURL)
	if err != nil {
		t.Fatalf("failed to parse connection URL: %v", err)
	}

	u.User = username
	u.Password = password
	testURL := buildMongoURL(u)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testURL))
	if err != nil {
		t.Fatalf("failed to connect with credentials: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("failed to ping with credentials: %v", err)
	}

	t.Logf("Verified MongoDB credentials for user: %s", username)
}

// mongoURL represents a parsed MongoDB connection URL.
type mongoURL struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Database string
	Options  string
}

// parseMongoURL parses a MongoDB connection URL into components.
func parseMongoURL(connURL string) (*mongoURL, error) {
	// Simple parser for mongodb:// URLs
	// Format: mongodb://user:pass@host/database?options
	u := &mongoURL{Scheme: "mongodb"}

	// Remove scheme
	rest := connURL
	if len(rest) > 10 && rest[:10] == "mongodb://" {
		rest = rest[10:]
	}

	// Extract user:pass if present
	atIdx := -1
	for i, c := range rest {
		if c == '@' {
			atIdx = i
			break
		}
	}

	if atIdx > 0 {
		userPass := rest[:atIdx]
		rest = rest[atIdx+1:]

		colonIdx := -1
		for i, c := range userPass {
			if c == ':' {
				colonIdx = i
				break
			}
		}

		if colonIdx > 0 {
			u.User = userPass[:colonIdx]
			u.Password = userPass[colonIdx+1:]
		}
	}

	// Extract host and database
	slashIdx := -1
	for i, c := range rest {
		if c == '/' {
			slashIdx = i
			break
		}
	}

	if slashIdx > 0 {
		u.Host = rest[:slashIdx]
		rest = rest[slashIdx+1:]

		// Extract database and options
		qIdx := -1
		for i, c := range rest {
			if c == '?' {
				qIdx = i
				break
			}
		}

		if qIdx > 0 {
			u.Database = rest[:qIdx]
			u.Options = rest[qIdx+1:]
		} else {
			u.Database = rest
		}
	} else {
		u.Host = rest
	}

	return u, nil
}

// buildMongoURL builds a MongoDB connection URL from components.
func buildMongoURL(u *mongoURL) string {
	url := u.Scheme + "://"

	if u.User != "" {
		url += u.User
		if u.Password != "" {
			url += ":" + u.Password
		}
		url += "@"
	}

	url += u.Host

	if u.Database != "" {
		url += "/" + u.Database
	}

	if u.Options != "" {
		url += "?" + u.Options
	}

	return url
}

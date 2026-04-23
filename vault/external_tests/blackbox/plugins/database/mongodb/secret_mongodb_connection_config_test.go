// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestMongoDBConnectionConfigCRUDWorkflows runs all MongoDB connection
// config CRUD workflow tests.
func TestMongoDBConnectionConfigCRUDWorkflows(t *testing.T) {
	t.Run("CreateBasic", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testMongoDBConnectionConfigCreateBasic(t, v)
	})

	// TODO: Additional tests for future implementation:
	// - ListReturnsNamesOnly: Verify listing returns only connection names
	// - ReadRedactsSensitiveFields: Verify passwords are redacted in read operations
	// - ResetConnection: Verify connection reset preserves configuration
	// - DeleteConnection: Verify connection deletion and idempotency
}

// TODO: TestMongoDBConnectionConfigValidationWorkflows
// Future implementation should test:
// - VerifyConnectionValid: Test with verify_connection=true and valid credentials
// - VerifyConnectionInvalid: Test with verify_connection=true and invalid credentials

// testMongoDBConnectionConfigCreateBasic verifies that configuring a
// MongoDB database connection succeeds at database/config/{name}.
func testMongoDBConnectionConfigCreateBasic(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)
	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	v.MustWrite(
		mount+"/config/my-mongodb-db",
		mongoConnectionConfigPayload(connURL, "*", false),
	)

	config := v.MustReadRequired(mount + "/config/my-mongodb-db")
	v.AssertSecret(config).Data().HasKey("plugin_name", "mongodb-database-plugin")
}

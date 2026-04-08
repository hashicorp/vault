// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package postgres

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestPostgreSQLConnectionConfigCRUDWorkflows runs all PostgreSQL connection
// config CRUD workflow tests.
func TestPostgreSQLConnectionConfigCRUDWorkflows(t *testing.T) {
	t.Run("CreateBasic", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigCreateBasic(t, v)
	})

	t.Run("ListReturnsNamesOnly", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigListReturnsNamesOnly(t, v)
	})

	t.Run("ReadRedactsSensitiveFields", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigReadRedactsSensitiveFields(t, v)
	})

	t.Run("ResetConnectionSuccess", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigResetConnectionSuccess(t, v)
	})

	t.Run("DeleteConnectionSuccess", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigDeleteConnectionSuccess(t, v)
	})
}

// TestPostgreSQLConnectionConfigValidationWorkflows runs all PostgreSQL
// connection config validation workflow tests.
func TestPostgreSQLConnectionConfigValidationWorkflows(t *testing.T) {
	t.Run("CreateWithDSN", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigCreateWithDSN(t, v)
	})

	t.Run("ConnectionVerificationValid", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigUpdateVerifyConnectionValidConnection(t, v)
	})

	t.Run("ConnectionVerificationInvalid", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigUpdateVerifyConnectionInvalidCredentials(t, v)
	})
}

// testPostgreSQLConnectionConfigCreateBasic verifies that configuring a
// PostgreSQL database connection succeeds at database/config/{name}.
func testPostgreSQLConnectionConfigCreateBasic(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	path := mount + "/config/my-postgres-db"

	v.MustWrite(path, postgresConnectionConfigPayload(
		templatedConnectionURL(connURL),
		"my-role",
		"secret",
		false,
	))

	config := v.MustReadRequired(path)
	v.AssertSecret(config).Data().
		HasKey("plugin_name", "postgresql-database-plugin").
		HasKeyExists("allowed_roles")
}

// testPostgreSQLConnectionConfigListReturnsNamesOnly verifies LIST
// /database/config returns configured connection names only, without
// sensitive config details.
func testPostgreSQLConnectionConfigListReturnsNamesOnly(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	configPath := mount + "/config"

	firstName := "postgres-list-config-one"
	secondName := "postgres-list-config-two"

	writeConfig := func(name string) {
		t.Helper()

		v.MustWrite(configPath+"/"+name, postgresConnectionConfigPayload(connURL, "*", "secret", false))
	}

	writeConfig(firstName)
	writeConfig(secondName)

	listResp, err := v.Client.Logical().List(configPath)
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.NotNil(t, listResp.Data)

	keysAny, ok := listResp.Data["keys"].([]any)
	require.True(t, ok, "expected keys list in LIST response")

	keys := make([]string, 0, len(keysAny))
	for _, raw := range keysAny {
		k, ok := raw.(string)
		require.True(t, ok, "expected string key in LIST response")
		keys = append(keys, k)
	}

	require.True(t, slices.Contains(keys, firstName), "expected first connection name in list")
	require.True(t, slices.Contains(keys, secondName), "expected second connection name in list")

	require.NotContains(t, listResp.Data, "connection_details")
	require.NotContains(t, listResp.Data, "connection_url")
	require.NotContains(t, listResp.Data, "username")
	require.NotContains(t, listResp.Data, "password")
	require.NotContains(t, listResp.Data, "private_key")
	require.NotContains(t, listResp.Data, "service_account_json")
}

// testPostgreSQLConnectionConfigReadRedactsSensitiveFields verifies that reading
// a PostgreSQL connection config returns sanitized connection details.
func testPostgreSQLConnectionConfigReadRedactsSensitiveFields(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	path := mount + "/config/postgres-read-config-redaction"
	v.MustWrite(path, postgresConnectionConfigPayload(connURL, "*", "secret", false))

	config := v.MustReadRequired(path)
	v.AssertSecret(config).Data().
		HasKey("plugin_name", "postgresql-database-plugin").
		HasKeyExists("connection_details")

	connectionDetails, ok := config.Data["connection_details"].(map[string]any)
	require.True(t, ok, "expected connection_details in response")

	if _, exists := connectionDetails["password"]; exists {
		t.Fatal("password should NOT be found in the returned config")
	}
	if _, exists := connectionDetails["private_key"]; exists {
		t.Fatal("private_key should NOT be found in the returned config")
	}
	if _, exists := connectionDetails["service_account_json"]; exists {
		t.Fatal("service_account_json should NOT be found in the returned config")
	}

	returnedConnURL, ok := connectionDetails["connection_url"].(string)
	require.True(t, ok, "expected connection_url in connection_details")
	require.NotContains(t, returnedConnURL, "secret", "connection_url should be redacted")
}

// testPostgreSQLConnectionConfigUpdateVerifyConnectionValidConnection verifies
// that configuration succeeds with verify_connection=true when credentials are
// valid.
func testPostgreSQLConnectionConfigUpdateVerifyConnectionValidConnection(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	templatedConnURL := strings.Replace(connURL, "postgres:secret@", "{{username}}:{{password}}@", 1)
	if templatedConnURL == connURL {
		t.Fatalf("failed to templatize postgres connection URL: %q", connURL)
	}

	path := mount + "/config/postgres-verify-enabled"
	v.MustWrite(path, postgresConnectionConfigPayload(templatedConnURL, "*", "secret", true))

	config := v.MustReadRequired(path)
	v.AssertSecret(config).Data().
		HasKey("plugin_name", "postgresql-database-plugin").
		HasKeyExists("allowed_roles")
}

// testPostgreSQLConnectionConfigUpdateVerifyConnectionInvalidCredentials
// verifies that configuration fails with verify_connection=true when
// credentials are intentionally invalid.
func testPostgreSQLConnectionConfigUpdateVerifyConnectionInvalidCredentials(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	templatedConnURL := strings.Replace(connURL, "postgres:secret@", "{{username}}:{{password}}@", 1)
	if templatedConnURL == connURL {
		t.Fatalf("failed to templatize postgres connection URL: %q", connURL)
	}

	path := mount + "/config/postgres-verify-enabled-invalid"
	secret, err := v.Client.Logical().Write(path, postgresConnectionConfigPayload(templatedConnURL, "*", "intentionally-invalid-password", true))

	require.Error(t, err)
	require.Nil(t, secret)
}

// testPostgreSQLConnectionConfigResetConnectionSuccess verifies that resetting a
// PostgreSQL database connection succeeds and preserves the stored connection
// configuration.
func testPostgreSQLConnectionConfigResetConnectionSuccess(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	connectionName := "my-postgres-db"
	configPath := fmt.Sprintf("%s/config/%s", mount, connectionName)
	resetPath := fmt.Sprintf("%s/reset/%s", mount, connectionName)

	templatedConnURL := strings.Replace(connURL, "postgres:secret@", "{{username}}:{{password}}@", 1)
	if templatedConnURL == connURL {
		t.Fatalf("failed to templatize postgres connection URL: %q", connURL)
	}

	v.MustWrite(configPath, postgresConnectionConfigPayload(templatedConnURL, "*", "secret", true))

	configBeforeReset := v.MustReadRequired(configPath)

	_, err := v.Client.Logical().Write(resetPath, nil)
	require.NoError(t, err)

	configAfterReset := v.MustReadRequired(configPath)
	require.Equal(t, configBeforeReset.Data, configAfterReset.Data)
}

// testPostgreSQLConnectionConfigDeleteConnectionSuccess verifies that deleting
// a PostgreSQL database connection removes it, prevents new credential
// generation,
// and remains idempotent when deleted again.
func testPostgreSQLConnectionConfigDeleteConnectionSuccess(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	connectionName := "postgres-secondary"
	configPath := fmt.Sprintf("%s/config/%s", mount, connectionName)
	roleName := "postgres-delete-connection-role"

	v.MustWrite(configPath, postgresConnectionConfigPayload(connURL, "*", "secret", true))

	creationSQL := `CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';`
	v.MustCreateDBRole(mount, roleName, connectionName, creationSQL)

	baselineCredsPath := fmt.Sprintf("%s/creds/%s", mount, roleName)
	baselineCreds, err := v.Client.Logical().Read(baselineCredsPath)
	require.NoError(t, err)
	require.NotNil(t, baselineCreds)

	_, err = v.Client.Logical().Delete(configPath)
	require.NoError(t, err)

	configAfterDelete, err := v.Client.Logical().Read(configPath)
	require.NoError(t, err)
	require.Nil(t, configAfterDelete, "config should be deleted")

	credsAfterDelete, err := v.Client.Logical().Read(baselineCredsPath)
	require.Error(t, err)
	if credsAfterDelete != nil {
		require.Empty(t, credsAfterDelete.Data)
	}

	_, err = v.Client.Logical().Delete(configPath)
	require.NoError(t, err)
}

func postgresConnectionConfigPayload(connURL, allowedRoles, password string, verifyConnection bool) map[string]any {
	return map[string]any{
		"plugin_name":       "postgresql-database-plugin",
		"connection_url":    connURL,
		"username":          "postgres",
		"password":          password,
		"allowed_roles":     allowedRoles,
		"verify_connection": verifyConnection,
	}
}

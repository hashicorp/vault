// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package postgres

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// testPostgreSQLConnectionConfigCreateWithDSN verifies that configuring a
// PostgreSQL database connection succeeds using DSN key-value format at
// database/config/{name}.
func testPostgreSQLConnectionConfigCreateWithDSN(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)
	cleanup, connURL := PrepareTestContainer(t)
	defer cleanup()

	dsnConnectionURL := dsnConnectionURLFromConnURL(t, connURL)

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})
	path := mount + "/config/postgres-dsn-basic"

	v.MustWrite(path, postgresConnectionConfigPayload(dsnConnectionURL, "*", "secret", true))

	config := v.MustReadRequired(path)
	v.AssertSecret(config).Data().
		HasKey("plugin_name", "postgresql-database-plugin").
		HasKeyExists("allowed_roles")
}

func dsnConnectionURLFromConnURL(t *testing.T, connURL string) string {
	t.Helper()

	parsedConnURL, err := url.Parse(connURL)
	require.NoError(t, err)
	require.NotEmpty(t, parsedConnURL.Host)

	host := parsedConnURL.Hostname()
	require.NotEmpty(t, host)

	port := parsedConnURL.Port()
	if port == "" {
		_, parsedPort, splitErr := net.SplitHostPort(parsedConnURL.Host)
		require.NoError(t, splitErr)
		port = parsedPort
	}

	dbName := strings.TrimPrefix(parsedConnURL.Path, "/")
	if dbName == "" {
		dbName = "postgres"
	}

	sslMode := parsedConnURL.Query().Get("sslmode")
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user={{username}} password={{password}} sslmode=%s",
		host,
		port,
		dbName,
		sslMode,
	)
}

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package postgres

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestPostgreSQLConnectionConfigFailoverWorkflows runs all PostgreSQL
// connection failover workflow tests.
func TestPostgreSQLConnectionConfigFailoverWorkflows(t *testing.T) {
	t.Run("MultiHostPrimaryAvailable", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigCreateMultiHostPrimaryAvailable(t, v)
	})

	t.Run("MultiHostPrimaryUnavailable", func(t *testing.T) {
		t.Parallel()
		v := blackbox.New(t)
		testPostgreSQLConnectionConfigCreateMultiHostPrimaryUnavailable(t, v)
	})
}

// testPostgreSQLConnectionConfigCreateMultiHostPrimaryAvailable verifies
// that a multi-host PostgreSQL connection succeeds when the primary host is
// available.
func testPostgreSQLConnectionConfigCreateMultiHostPrimaryAvailable(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanupPrimary, primaryConnURL := PrepareTestContainer(t)
	defer cleanupPrimary()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	primaryPreferredURL := strings.Replace(
		primaryConnURL,
		"/postgres?sslmode=disable",
		",localhost:55/postgres?sslmode=disable",
		1,
	)
	if primaryPreferredURL == primaryConnURL {
		t.Fatalf("failed to construct primary-preferred multi-host URL from %q", primaryConnURL)
	}

	path := mount + "/config/postgres-multihost-primary"
	writeAndAssertPostgresConfig(t, v, path, templatedConnectionURL(primaryPreferredURL))
}

// testPostgreSQLConnectionConfigCreateMultiHostPrimaryUnavailable verifies
// that a multi-host PostgreSQL connection succeeds when the primary host is
// down and the driver falls back to a secondary host.
func testPostgreSQLConnectionConfigCreateMultiHostPrimaryUnavailable(t *testing.T, v *blackbox.Session) {
	requireVaultEnv(t)

	cleanupFailover, failoverConnURL := PrepareTestContainerMultiHost(t)
	defer cleanupFailover()

	mount := fmt.Sprintf("database-%s", sanitize(t.Name()))
	v.MustEnableSecretsEngine(mount, &api.MountInput{Type: "database"})

	if !strings.Contains(failoverConnURL, "localhost:55,") {
		t.Fatalf("expected failover URL with unreachable first host, got %q", failoverConnURL)
	}

	path := mount + "/config/postgres-multihost-ports"
	writeAndAssertPostgresConfig(t, v, path, templatedConnectionURL(failoverConnURL))
}

func templatedConnectionURL(connURL string) string {
	return strings.Replace(connURL, "postgres:secret@", "{{username}}:{{password}}@", 1)
}

func writeAndAssertPostgresConfig(t *testing.T, v *blackbox.Session, path, connURL string) {
	t.Helper()

	v.MustWrite(path, postgresConnectionConfigPayload(connURL, "*", "secret", true))

	config := v.MustReadRequired(path)
	v.AssertSecret(config).Data().
		HasKey("plugin_name", "postgresql-database-plugin").
		HasKeyExists("allowed_roles")
}

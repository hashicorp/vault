// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestPostgresDynamicSecrets verifies the database dynamic secrets engine functionality
// by configuring a PostgreSQL connection, creating a role, generating credentials,
// and testing the full lifecycle including credential revocation.
func TestPostgresDynamicSecrets(t *testing.T) {
	v := blackbox.New(t)

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	connURL := fmt.Sprintf("postgres://{{username}}:{{password}}@localhost:5432/%s?sslmode=disable", db)

	v.MustEnableSecretsEngine("database", &api.MountInput{Type: "database"})
	v.MustConfigureDBConnection(
		"database",
		"my-postgres",
		"postgresql-database-plugin",
		connURL,
		map[string]any{
			"username": user,
			"password": pass,
		},
	)

	creationSQL := `CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';`
	v.MustCreateDBRole("database", "readonly-role", "my-postgres", creationSQL)

	creds := v.MustGenerateCreds("database/creds/readonly-role")
	t.Logf("generated DB user/pass: %s / %s", creds.Username, creds.Password)

	v.AssertLeaseExists(creds.LeaseID)
	v.MustCheckCreds(creds.Username, creds.Password, true)
	v.MustRevokeLease(creds.LeaseID)
	v.AssertLeaseRevoked(creds.LeaseID)
	v.MustCheckCreds(creds.Username, creds.Password, false)
}

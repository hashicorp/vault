// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
)

type DynamicSecret struct {
	Secret   *api.Secret
	LeaseID  string
	Username string
	Password string
}

func (s *Session) MustGenerateCreds(path string) *DynamicSecret {
	s.t.Helper()

	secret := s.MustReadRequired(path)

	ds := &DynamicSecret{
		Secret:  secret,
		LeaseID: secret.LeaseID,
	}

	// usually the creds are in the 'data' map
	if val, ok := secret.Data["username"]; ok {
		if username, ok := val.(string); ok {
			ds.Username = username
		} else {
			s.t.Fatalf("username field is not a string, got type %T with value %v", val, val)
		}
	}
	if val, ok := secret.Data["password"]; ok {
		if password, ok := val.(string); ok {
			ds.Password = password
		} else {
			s.t.Fatalf("password field is not a string, got type %T with value %v", val, val)
		}
	}

	if ds.Username == "" || ds.Password == "" {
		s.t.Fatal("expected username and password to be populated")
	}

	return ds
}

func (s *Session) MustRevokeLease(leaseID string) {
	s.t.Helper()

	err := s.Client.Sys().Revoke(leaseID)
	require.NoError(s.t, err)
}

func (s *Session) AssertLeaseExists(leaseID string) {
	s.t.Helper()

	_, err := s.Client.Sys().Lookup(leaseID)
	require.NoError(s.t, err)
}

func (s *Session) AssertLeaseRevoked(leaseID string) {
	s.t.Helper()

	// when a lease is revoked, Lookup returns an error, so we expect one here
	_, err := s.Client.Sys().Lookup(leaseID)
	require.Error(s.t, err)
}

func (s *Session) MustConfigureDBConnection(mountPath, name, plugin, connectionURL string, extraConfig map[string]any) {
	s.t.Helper()

	path := fmt.Sprintf("%s/config/%s", mountPath, name)
	payload := map[string]any{
		"plugin_name":    plugin,
		"connection_url": connectionURL,
		"allowed_roles":  "*",
	}

	// merge any extras
	for k, v := range extraConfig {
		payload[k] = v
	}

	s.MustWrite(path, payload)
}

func (s *Session) MustCreateDBRole(mountPath, roleName, dbName, creationSQL string) {
	s.t.Helper()

	path := fmt.Sprintf("%s/roles/%s", mountPath, roleName)
	payload := map[string]any{
		"db_name":             dbName,
		"creation_statements": creationSQL,
		"default_ttl":         "1h",
		"max_ttl":             "24h",
	}

	s.MustWrite(path, payload)
}

// MustCheckCreds verifies database credentials work (or don't work) against PostgreSQL.
// Uses POSTGRES_HOST, POSTGRES_PORT, and POSTGRES_DB environment variables,
// defaulting to localhost:5432/vault if not set.
func (s *Session) MustCheckCreds(username, password string, shouldBeValid bool) {
	s.t.Helper()

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "vault"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbName)
	db, err := sql.Open("pgx", connStr)
	require.NoError(s.t, err)
	defer func() { _ = db.Close() }()

	err = db.Ping()
	if shouldBeValid {
		require.NoError(s.t, err)
	} else {
		require.Error(s.t, err)
	}
}

// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mssql

import (
	"database/sql"
	"net/url"
	"os"
	"strings"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	mssqlhelper "github.com/hashicorp/vault/helper/testhelpers/mssql"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/require"
)

// TestInvalidIdentifier checks validity of an identifier
func TestInvalidIdentifier(t *testing.T) {
	testcases := map[string]bool{
		"name":             true,
		"_name":            true,
		"Name":             true,
		"#name":            false,
		"?Name":            false,
		"9name":            false,
		"@name":            false,
		"$name":            false,
		" name":            false,
		"n ame":            false,
		"n4444444":         true,
		"_4321098765":      true,
		"_##$$@@__":        true,
		"_123name#@":       true,
		"name!":            false,
		"name%":            false,
		"name^":            false,
		"name&":            false,
		"name*":            false,
		"name(":            false,
		"name)":            false,
		"nåame":            true,
		"åname":            true,
		"name'":            false,
		"nam`e":            false,
		"пример":           true,
		"_#Āā@#$_ĂĄąćĈĉĊċ": true,
		"ÛÜÝÞßàáâ":         true,
		"豈更滑a23$#@":        true,
	}

	for i, expected := range testcases {
		if !isInvalidIdentifier(i) != expected {
			t.Fatalf("unexpected identifier %s: expected validity %v", i, expected)
		}
	}
}

// TestQuoteIdentifier verifies that delimiters and SQL-looking text remain
// inside a single identifier, independently of configuration validation.
func TestQuoteIdentifier(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"Vault":              "[Vault]",
		"va]ult":             "[va]]ult]",
		"[vault]":            "[[vault]]]",
		"vault.schema":       "[vault.schema]",
		"vault'quoted":       "[vault'quoted]",
		"vault]; SELECT 1--": "[vault]]; SELECT 1--]",
		"\u7528\u6237":       "[\u7528\u6237]",
	}
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, want, quoteIdentifier(name), "identifier must be delimited and escaped")
		})
	}
}

// TestNewMSSQLBackendRejectsBracketInConfig preserves the existing validation
// contract without connecting to a database.
func TestNewMSSQLBackendRejectsBracketInConfig(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"database", "schema", "table"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			conf := map[string]string{"server": "localhost", field: "vault]"}
			_, err := NewMSSQLBackend(conf, logging.NewVaultLogger(log.Debug))
			require.EqualError(t, err, "invalid "+field+" name", "invalid identifiers must fail before connecting")
		})
	}
}

// TestMSSQLBackendQuotedIdentifiers exercises schema creation and CRUD with
// reserved-word identifiers, which require quoting despite passing validation.
func TestMSSQLBackendQuotedIdentifiers(t *testing.T) {
	t.Parallel()
	cleanup, connURL := mssqlhelper.PrepareMSSQLTestContainer(t)
	t.Cleanup(cleanup)
	conn, err := url.Parse(connURL)
	require.True(t, err == nil, "parse test connection (error not printed: it would include the password)")
	setupDB, err := sql.Open("sqlserver", connURL)
	require.NoError(t, err, "open setup connection")
	t.Cleanup(func() { require.NoError(t, setupDB.Close(), "close setup connection") })
	id, err := uuid.GenerateUUID()
	require.NoError(t, err, "generate isolated database name")
	database := "vault_quote_" + strings.ReplaceAll(id, "-", "")
	t.Cleanup(func() {
		_, err := setupDB.Exec("IF DB_ID(@p1) IS NOT NULL BEGIN ALTER DATABASE "+quoteIdentifier(database)+
			" SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE "+quoteIdentifier(database)+"; END", database)
		require.NoError(t, err, "drop isolated database")
	})
	password, _ := conn.User.Password()
	b, err := NewMSSQLBackend(map[string]string{
		"server":   conn.Hostname(),
		"port":     conn.Port(),
		"username": conn.User.Username(),
		"password": password,
		"database": database,
		"schema":   "select",
		"table":    "table",
	}, logging.NewVaultLogger(log.Debug))
	require.NoError(t, err, "create backend with quoted schema and table")
	backend := b.(*MSSQLBackend)
	t.Cleanup(func() { require.NoError(t, backend.client.Close(), "close backend") })
	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestMSSQLBackend(t *testing.T) {
	server := os.Getenv("MSSQL_SERVER")
	if server == "" {
		t.SkipNow()
	}

	database := os.Getenv("MSSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MSSQL_TABLE")
	if table == "" {
		table = "test"
	}

	schema := os.Getenv("MSSQL_SCHEMA")
	if schema == "" {
		schema = "test"
	}

	username := os.Getenv("MSSQL_USERNAME")
	password := os.Getenv("MSSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMSSQLBackend(map[string]string{
		"server":   server,
		"database": database,
		"table":    table,
		"schema":   schema,
		"username": username,
		"password": password,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mssql := b.(*MSSQLBackend)
		_, err := mssql.client.Exec("DROP TABLE " + mssql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestMSSQLBackend_schema(t *testing.T) {
	server := os.Getenv("MSSQL_SERVER")
	if server == "" {
		t.SkipNow()
	}

	database := os.Getenv("MSSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MSSQL_TABLE")
	if table == "" {
		table = "test"
	}

	schema := os.Getenv("MSSQL_SCHEMA")
	if schema == "" {
		schema = "test"
	}

	username := os.Getenv("MSSQL_USERNAME")
	password := os.Getenv("MSSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMSSQLBackend(map[string]string{
		"server":   server,
		"database": database,
		"schema":   schema,
		"table":    table,
		"username": username,
		"password": password,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mssql := b.(*MSSQLBackend)
		_, err := mssql.client.Exec("DROP TABLE " + mssql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

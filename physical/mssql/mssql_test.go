// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
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

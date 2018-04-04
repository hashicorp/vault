package mssql

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	_ "github.com/denisenkom/go-mssqldb"
)

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

	username := os.Getenv("MSSQL_USERNAME")
	password := os.Getenv("MSSQL_PASSWORD")

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewMSSQLBackend(map[string]string{
		"server":   server,
		"database": database,
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

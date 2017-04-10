package physical

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestMsSQLBackend(t *testing.T) {
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
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("mssql", logger, map[string]string{
		"server":   server,
		"database": database,
		"table":    table,
		"username": username,
		"password": password,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		mssql := b.(*MsSQLBackend)
		_, err := mssql.client.Exec("DROP TABLE " + mssql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}

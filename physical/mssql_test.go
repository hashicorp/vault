package physical

import (
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestMsSQLBackend(t *testing.T) {
	address := os.Getenv("MSSQL_ADDR")
	if address == "" {
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
	b, err := NewBackend("mssql", map[string]string{
		"address":  address,
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

package physical

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgreSQLBackend(t *testing.T) {
	address := os.Getenv("PGSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("PGSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("PGSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("PGSQL_USERNAME")
	password := os.Getenv("PGSQL_PASSWORD")

	// Run vault tests
	b, err := NewBackend("postgres", map[string]string{
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
		pgsql := b.(*PostgreSQLBackend)
		_, err := pgsql.client.Exec("DROP TABLE " + pgsql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}

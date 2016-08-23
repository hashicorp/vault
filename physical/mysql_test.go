package physical

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLBackend(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	database := os.Getenv("MYSQL_DB")
	if database == "" {
		database = "test"
	}

	table := os.Getenv("MYSQL_TABLE")
	if table == "" {
		table = "test"
	}

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")

	// Run vault tests
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("mysql", logger, map[string]string{
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
		mysql := b.(*MySQLBackend)
		_, err := mysql.client.Exec("DROP TABLE " + mysql.dbTable)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}

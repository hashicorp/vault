package postgresql

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	_ "github.com/lib/pq"
)

func TestPostgreSQLBackend(t *testing.T) {
	connURL := os.Getenv("PGURL")
	if connURL == "" {
		t.SkipNow()
	}

	table := os.Getenv("PGTABLE")
	if table == "" {
		table = "vault_kv_store"
	}

	// Run vault tests
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		pg := b.(*PostgreSQLBackend)
		_, err := pg.client.Exec("TRUNCATE TABLE " + pg.table)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}

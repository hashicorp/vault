package postgresql

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

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

	hae := os.Getenv("PGHAENABLED")
	if hae == "" {
		hae = "false"
	}

	// Run vault tests
	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
		"ha_enabled":     hae,
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	b2, err := NewPostgreSQLBackend(map[string]string{
		"connection_url": connURL,
		"table":          table,
	}, logger)

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		pg := b.(*PostgreSQLBackend)
		_, err := pg.client.Exec("DELETE FROM " + table)
		if err != nil {
			t.Fatalf("Failed to delete table: %v", err)
		}
	}()

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)

	ha1, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("PostgreSQLDB does not implement HABackend")
	}

	ha2, ok := b2.(physical.HABackend)
	if !ok {
		t.Fatalf("PostgreSQLDB does not implement HABackend")
	}

	if ha1.HAEnabled() && ha1.HAEnabled() {
		physical.ExerciseHABackend(t, ha1, ha2)
	}
}

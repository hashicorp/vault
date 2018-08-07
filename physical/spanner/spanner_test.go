package spanner

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func testCleanup(t testing.TB, client *spanner.Client, table string) {
	t.Helper()

	// Delete all data in the table
	ctx := context.Background()
	m := spanner.Delete(table, spanner.AllKeys())
	if _, err := client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
		t.Fatal(err)
	}
}

func TestBackend(t *testing.T) {
	database := os.Getenv("GOOGLE_SPANNER_DATABASE")
	if database == "" {
		t.Skip("GOOGLE_SPANNER_DATABASE not set")
	}

	table := os.Getenv("GOOGLE_SPANNER_TABLE")
	if table == "" {
		t.Skip("GOOGLE_SPANNER_TABLE not set")
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		t.Fatal(err)
	}

	testCleanup(t, client, table)
	defer testCleanup(t, client, table)

	backend, err := NewBackend(map[string]string{
		"database":   database,
		"table":      table,
		"ha_enabled": "false",
	}, logging.NewVaultLogger(log.Debug))
	if err != nil {
		t.Fatal(err)
	}

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
	physical.ExerciseTransactionalBackend(t, backend)
}

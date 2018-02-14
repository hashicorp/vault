package spanner

import (
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
)

func TestHABackend(t *testing.T) {
	database := os.Getenv("GOOGLE_SPANNER_DATABASE")
	if database == "" {
		t.Skip("GOOGLE_SPANNER_DATABASE not set")
	}

	table := os.Getenv("GOOGLE_SPANNER_TABLE")
	if table == "" {
		t.Skip("GOOGLE_SPANNER_TABLE not set")
	}

	haTable := os.Getenv("GOOGLE_SPANNER_HA_TABLE")
	if haTable == "" {
		t.Skip("GOOGLE_SPANNER_HA_TABLE not set")
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		t.Fatal(err)
	}

	testCleanup(t, client, table)
	defer testCleanup(t, client, table)
	testCleanup(t, client, haTable)
	defer testCleanup(t, client, haTable)

	backend, err := NewBackend(map[string]string{
		"database":   database,
		"table":      table,
		"ha_table":   haTable,
		"ha_enabled": "true",
	}, logformat.NewVaultLogger(log.LevelTrace))
	if err != nil {
		t.Fatal(err)
	}

	ha, ok := backend.(physical.HABackend)
	if !ok {
		t.Fatalf("does not implement")
	}

	physical.ExerciseHABackend(t, ha, ha)
}

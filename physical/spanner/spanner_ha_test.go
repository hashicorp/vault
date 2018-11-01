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

	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"database":   database,
		"table":      table,
		"ha_table":   haTable,
		"ha_enabled": "true",
	}

	b, err := NewBackend(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := NewBackend(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
}

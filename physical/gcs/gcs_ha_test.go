package gcs

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func TestHABackend(t *testing.T) {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOOGLE_PROJECT_ID not set")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-gcs-testacc-%d", r)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testCleanup(t, client, bucket)
	defer testCleanup(t, client, bucket)

	bh := client.Bucket(bucket)
	if err := bh.Create(context.Background(), projectID, nil); err != nil {
		t.Fatal(err)
	}

	logger := logging.NewVaultLogger(log.Trace)
	config := map[string]string{
		"bucket":     bucket,
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

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
}

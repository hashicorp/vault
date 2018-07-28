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
	"google.golang.org/api/googleapi"
)

func testCleanup(t testing.TB, client *storage.Client, bucket string) {
	t.Helper()

	ctx := context.Background()
	if err := client.Bucket(bucket).Delete(ctx); err != nil {
		if terr, ok := err.(*googleapi.Error); !ok || terr.Code != 404 {
			t.Fatal(err)
		}
	}
}

func TestBackend(t *testing.T) {
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

	b := client.Bucket(bucket)
	if err := b.Create(context.Background(), projectID, nil); err != nil {
		t.Fatal(err)
	}

	backend, err := NewBackend(map[string]string{
		"bucket":     bucket,
		"ha_enabled": "false",
	}, logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatal(err)
	}

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

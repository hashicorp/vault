package gcs

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
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
	}, logformat.NewVaultLogger(log.LevelTrace))
	if err != nil {
		t.Fatal(err)
	}

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

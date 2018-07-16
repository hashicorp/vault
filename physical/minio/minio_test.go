package minio

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func TestBackend(t *testing.T) {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		t.Skip("MINIO_ENDPOINT not set")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-minio-testacc-%d", r)

	backend, err := NewMinioBackend(map[string]string{
		"bucket":      bucket,
		"disable_ssl": "false",
	}, logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatal(err)
	}

	minioBackend := backend.(*MinioBackend)
	defer minioBackend.removeBucket()

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

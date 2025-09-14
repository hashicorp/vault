// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package gcs

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
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

	universeDomain := os.Getenv("GOOGLE_UNIVERSE_DOMAIN")

	r := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-gcs-testacc-%d", r)

	ctx := context.Background()
	// Build client options: if a custom universe domain is provided in env, use it.
	clientOpts := []option.ClientOption{}
	if universeDomain != "" {
		clientOpts = append(clientOpts, option.WithUniverseDomain(universeDomain))
	}
	client, err := storage.NewClient(ctx, clientOpts...)
	if err != nil {
		t.Fatal(err)
	}

	testCleanup(t, client, bucket)
	defer testCleanup(t, client, bucket)

	b := client.Bucket(bucket)
	// Support minimal for providers that require an explicit Location.
	bucketLocation := os.Getenv("GOOGLE_BUCKET_LOCATION")

	if bucketLocation != "" {
		// Create the bucket with the explicit location required by some universe domains.
		if err := b.Create(ctx, projectID, &storage.BucketAttrs{
			Location: bucketLocation,
		}); err != nil {
			t.Fatalf("failed to create bucket %q with location %q: %v", bucket, bucketLocation, err)
		}
	} else {
		// Default behaviour (no explicit location).
		if err := b.Create(ctx, projectID, nil); err != nil {
			t.Fatalf("failed to create bucket %q: %v", bucket, err)
		}
	}

	backend, err := NewBackend(map[string]string{
		"bucket":          bucket,
		"ha_enabled":      "false",
		"universe_domain": universeDomain,
	}, logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatal(err)
	}

	// Verify chunkSize is set correctly on the Backend
	be := backend.(*Backend)
	expectedChunkSize, err := strconv.Atoi(defaultChunkSize)
	if err != nil {
		t.Fatalf("failed to convert defaultChunkSize to int: %s", err)
	}
	expectedChunkSize = expectedChunkSize * 1024
	if be.chunkSize != expectedChunkSize {
		t.Fatalf("expected chunkSize to be %d. got=%d", expectedChunkSize, be.chunkSize)
	}

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

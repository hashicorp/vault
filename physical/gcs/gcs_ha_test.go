// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"google.golang.org/api/option"
)

func TestHABackend(t *testing.T) {
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

	testCleanup(t, client, bucket)
	defer testCleanup(t, client, bucket)

	bh := client.Bucket(bucket)
	// Support minimal for providers that require an explicit Location.
	bucketLocation := os.Getenv("GOOGLE_BUCKET_LOCATION")

	if bucketLocation != "" {
		// Create the bucket with the explicit location required by some universe domains.
		if err := bh.Create(ctx, projectID, &storage.BucketAttrs{
			Location: bucketLocation,
		}); err != nil {
			t.Fatalf("failed to create bucket %q with location %q: %v", bucket, bucketLocation, err)
		}
	} else {
		// Default behaviour (no explicit location).
		if err := bh.Create(ctx, projectID, nil); err != nil {
			t.Fatalf("failed to create bucket %q: %v", bucket, err)
		}
	}

	logger := logging.NewVaultLogger(log.Trace)
	config := map[string]string{
		"bucket":          bucket,
		"ha_enabled":      "true",
		"universe_domain": universeDomain,
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

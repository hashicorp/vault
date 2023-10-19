// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agentproxyshared

import (
	"context"
	"os"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func testNewLeaseCache(t *testing.T, responses []*cache.SendResponse) *cache.LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	lc, err := cache.NewLeaseCache(&cache.LeaseCacheConfig{
		Client:      client,
		BaseContext: context.Background(),
		Proxier:     cache.NewMockProxier(responses),
		Logger:      logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return lc
}

func populateTempFile(t *testing.T, name, contents string) *os.File {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.WriteString(contents)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	return file
}

// Test_AddPersistentStorageToLeaseCache Tests that AddPersistentStorageToLeaseCache() correctly
// adds persistent storage to a lease cache
func Test_AddPersistentStorageToLeaseCache(t *testing.T) {
	tempDir := t.TempDir()
	serviceAccountTokenFile := populateTempFile(t, "proxy-config.hcl", "token")

	persistConfig := &PersistConfig{
		Type:                    "kubernetes",
		Path:                    tempDir,
		KeepAfterImport:         false,
		ExitOnErr:               false,
		ServiceAccountTokenFile: serviceAccountTokenFile.Name(),
	}

	leaseCache := testNewLeaseCache(t, nil)
	if leaseCache.PersistentStorage() != nil {
		t.Fatal("persistent storage was available before ours was added")
	}

	deferFunc, token, err := AddPersistentStorageToLeaseCache(context.Background(), leaseCache, persistConfig, logging.NewVaultLogger(hclog.Info))
	if err != nil {
		t.Fatal(err)
	}

	if leaseCache.PersistentStorage() == nil {
		t.Fatal("persistent storage was not added")
	}

	if token != "" {
		t.Fatal("expected token to be empty")
	}

	if deferFunc == nil {
		t.Fatal("expected deferFunc to not be nil")
	}
}

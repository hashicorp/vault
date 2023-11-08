// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agentproxyshared

import (
	"context"
	"os"
	"testing"

	cache2 "github.com/hashicorp/vault/command-server/agentproxyshared/cache"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func testNewLeaseCache(t *testing.T, responses []*cache2.SendResponse) *cache2.LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	lc, err := cache2.NewLeaseCache(&cache2.LeaseCacheConfig{
		Client:         client,
		BaseContext:    context.Background(),
		Proxier:        cache2.NewMockProxier(responses),
		Logger:         logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
		UserAgentToUse: "test",
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

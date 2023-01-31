package approle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestAppRole_TidyDanglingAccessors_Normal(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := []*framework.Path{pathTidySecretID(b)}

	// Create a role
	createRole(t, b, storage, "role1", "a,b,c")

	// Create a secret-id
	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	accessorHashes, err := storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 1 {
		t.Fatalf("bad: len(accessorHashes); expect 1, got %d", len(accessorHashes))
	}

	entry1, err := logical.StorageEntryJSON(
		"accessor/invalid1",
		&secretIDAccessorStorageEntry{
			SecretIDHMAC: "samplesecretidhmac",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := storage.Put(context.Background(), entry1); err != nil {
		t.Fatal(err)
	}

	entry2, err := logical.StorageEntryJSON(
		"accessor/invalid2",
		&secretIDAccessorStorageEntry{
			SecretIDHMAC: "samplesecretidhmac2",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry2); err != nil {
		t.Fatal(err)
	}

	accessorHashes, err = storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 3 {
		t.Fatalf("bad: len(accessorHashes); expect 3, got %d", len(accessorHashes))
	}

	secret, err := b.tidySecretID(context.Background(), &logical.Request{
		Storage: storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	schema.ValidateResponse(
		t,
		schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation),
		secret,
		true,
	)

	// It runs async so we give it a bit of time to run
	time.Sleep(10 * time.Second)

	accessorHashes, err = storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 1 {
		t.Fatalf("bad: len(accessorHashes); expect 1, got %d", len(accessorHashes))
	}
}

func TestAppRole_TidyDanglingAccessors_RaceTest(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	paths := []*framework.Path{pathTidySecretID(b)}

	// Create a role
	createRole(t, b, storage, "role1", "a,b,c")

	// Create an initial entry
	roleSecretIDReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/role1/secret-id",
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), roleSecretIDReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}
	count := 1

	wg := &sync.WaitGroup{}
	start := time.Now()
	for time.Now().Sub(start) < 10*time.Second {
		if time.Now().Sub(start) > 100*time.Millisecond && atomic.LoadUint32(b.tidySecretIDCASGuard) == 0 {
			secret, err := b.tidySecretID(context.Background(), &logical.Request{
				Storage: storage,
			})
			if err != nil {
				t.Fatal(err)
			}
			schema.ValidateResponse(
				t,
				schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation),
				secret,
				true,
			)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			roleSecretIDReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "role/role1/secret-id",
				Storage:   storage,
			}
			resp, err := b.HandleRequest(context.Background(), roleSecretIDReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%v resp:%#v", err, resp)
			}
		}()

		entry, err := logical.StorageEntryJSON(
			fmt.Sprintf("accessor/invalid%d", count),
			&secretIDAccessorStorageEntry{
				SecretIDHMAC: "samplesecretidhmac",
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := storage.Put(context.Background(), entry); err != nil {
			t.Fatal(err)
		}

		count++
		time.Sleep(100 * time.Microsecond)
	}

	logger := b.Logger().Named(t.Name())
	logger.Info("wrote entries", "count", count)

	wg.Wait()
	// Let tidy finish
	for atomic.LoadUint32(b.tidySecretIDCASGuard) != 0 {
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("running tidy again")

	// Run tidy again
	secret, err := b.tidySecretID(context.Background(), &logical.Request{
		Storage: storage,
	})
	if err != nil || len(secret.Warnings) > 0 {
		t.Fatal(err, secret.Warnings)
	}
	schema.ValidateResponse(
		t,
		schema.FindResponseSchema(t, paths, 0, logical.UpdateOperation),
		secret,
		true,
	)

	// Wait for tidy to start
	for atomic.LoadUint32(b.tidySecretIDCASGuard) == 0 {
		time.Sleep(100 * time.Millisecond)
	}

	// Let tidy finish
	for atomic.LoadUint32(b.tidySecretIDCASGuard) != 0 {
		time.Sleep(100 * time.Millisecond)
	}

	accessorHashes, err := storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != count {
		t.Fatalf("bad: len(accessorHashes); expect %d, got %d", count, len(accessorHashes))
	}

	roleHMACs, err := storage.List(context.Background(), secretIDPrefix)
	if err != nil {
		t.Fatal(err)
	}
	secretIDs, err := storage.List(context.Background(), fmt.Sprintf("%s%s", secretIDPrefix, roleHMACs[0]))
	if err != nil {
		t.Fatal(err)
	}
	if len(secretIDs) != count {
		t.Fatalf("bad: len(secretIDs); expect %d, got %d", count, len(secretIDs))
	}
}

package approle

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestAppRole_TidyDanglingAccessors_Normal(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	err = storage.Put(context.Background(), entry1)
	if err != nil {
		t.Fatal(err)
	}

	entry2, err := logical.StorageEntryJSON(
		"accessor/invalid2",
		&secretIDAccessorStorageEntry{
			SecretIDHMAC: "samplesecretidhmac2",
		},
	)
	err = storage.Put(context.Background(), entry2)
	if err != nil {
		t.Fatal(err)
	}

	accessorHashes, err = storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 3 {
		t.Fatalf("bad: len(accessorHashes); expect 3, got %d", len(accessorHashes))
	}

	_, err = b.tidySecretID(context.Background(), &logical.Request{
		Storage: storage,
	})
	if err != nil {
		t.Fatal(err)
	}

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

	b.testTidyDelay = 300 * time.Millisecond

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
	now := time.Now()
	started := false
	for {
		if time.Now().Sub(now) > 700*time.Millisecond {
			break
		}
		if time.Now().Sub(now) > 100*time.Millisecond && !started {
			started = true
			_, err = b.tidySecretID(context.Background(), &logical.Request{
				Storage: storage,
			})
			if err != nil {
				t.Fatal(err)
			}
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
		count++
	}

	t.Logf("wrote %d entries", count)

	wg.Wait()
	// Let tidy finish
	time.Sleep(1 * time.Second)

	// Run tidy again
	_, err = b.tidySecretID(context.Background(), &logical.Request{
		Storage: storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

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

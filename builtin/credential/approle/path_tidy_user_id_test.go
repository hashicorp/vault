package approle

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestAppRole_TidyDanglingAccessors(t *testing.T) {
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

	err = b.tidySecretID(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}

	accessorHashes, err = storage.List(context.Background(), "accessor/")
	if err != nil {
		t.Fatal(err)
	}
	if len(accessorHashes) != 1 {
		t.Fatalf("bad: len(accessorHashes); expect 1, got %d", len(accessorHashes))
	}
}

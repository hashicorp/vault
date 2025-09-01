// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/snapshots"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestCubbyholeBackend_Write(t *testing.T) {
	b := testCubbyholeBackend()
	req := logical.TestRequest(t, logical.UpdateOperation, "foo")
	clientToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	req.ClientToken = clientToken
	storage := req.Storage
	req.Data["raw"] = "test"

	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage
	req.ClientToken = clientToken
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestCubbyholeBackend_Read(t *testing.T) {
	b := testCubbyholeBackend()
	req := logical.TestRequest(t, logical.UpdateOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage
	clientToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	req.ClientToken = clientToken

	if _, err := b.HandleRequest(context.Background(), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage
	req.ClientToken = clientToken

	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		Data: map[string]interface{}{
			"raw": "test",
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}

func TestCubbyholeBackend_Delete(t *testing.T) {
	b := testCubbyholeBackend()
	req := logical.TestRequest(t, logical.UpdateOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage
	clientToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	req.ClientToken = clientToken

	if _, err := b.HandleRequest(context.Background(), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.DeleteOperation, "foo")
	req.Storage = storage
	req.ClientToken = clientToken
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage
	req.ClientToken = clientToken
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestCubbyholeBackend_List(t *testing.T) {
	b := testCubbyholeBackend()
	req := logical.TestRequest(t, logical.UpdateOperation, "foo")
	clientToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	req.Data["raw"] = "test"
	req.ClientToken = clientToken
	storage := req.Storage

	if _, err := b.HandleRequest(context.Background(), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "bar")
	req.Data["raw"] = "baz"
	req.ClientToken = clientToken
	req.Storage = storage

	if _, err := b.HandleRequest(context.Background(), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ListOperation, "")
	req.Storage = storage
	req.ClientToken = clientToken
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expKeys := []string{"foo", "bar"}
	respKeys := resp.Data["keys"].([]string)
	sort.Strings(expKeys)
	sort.Strings(respKeys)
	if !reflect.DeepEqual(respKeys, expKeys) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expKeys, respKeys)
	}
}

func TestCubbyholeIsolation(t *testing.T) {
	b := testCubbyholeBackend()

	clientTokenA, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	clientTokenB, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	var storageA logical.Storage
	var storageB logical.Storage

	// Populate and test A entries
	req := logical.TestRequest(t, logical.UpdateOperation, "foo")
	req.ClientToken = clientTokenA
	storageA = req.Storage
	req.Data["raw"] = "test"

	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storageA
	req.ClientToken = clientTokenA
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		Data: map[string]interface{}{
			"raw": "test",
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}

	// Populate and test B entries
	req = logical.TestRequest(t, logical.UpdateOperation, "bar")
	req.ClientToken = clientTokenB
	storageB = req.Storage
	req.Data["raw"] = "baz"

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "bar")
	req.Storage = storageB
	req.ClientToken = clientTokenB
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected = &logical.Response{
		Data: map[string]interface{}{
			"raw": "baz",
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}

	// We shouldn't be able to read A from B and vice versa
	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storageB
	req.ClientToken = clientTokenB
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("err: was able to read from other user's cubbyhole")
	}

	req = logical.TestRequest(t, logical.ReadOperation, "bar")
	req.Storage = storageA
	req.ClientToken = clientTokenA
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("err: was able to read from other user's cubbyhole")
	}
}

func testCubbyholeBackend() logical.Backend {
	b, _ := CubbyholeBackendFactory(context.Background(), &logical.BackendConfig{
		Logger: nil,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	})
	return b
}

// TestSnapshotOperations verifies that snapshot operations on the cubbyhole
// backend succeed. It tests reading, listing, recovering, and recovering as a
// copy from a snapshot.
func TestSnapshotOperations(t *testing.T) {
	backend := testCubbyholeBackend()
	tc := snapshots.NewSnapshotTestCase(t, backend)
	path := "secret"
	_, err := backend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   tc.SnapshotStorage(),
		Path:      path,
		Data: map[string]interface{}{
			"data": "old data",
		},
		ClientToken: "a",
	})
	require.NoError(t, err)

	_, err = backend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   tc.RegularStorage(),
		Path:      path,
		Data: map[string]interface{}{
			"data": "new data",
		},
		ClientToken: "a",
	})

	checkContent := func(t *testing.T, path string) {
		readResp, err := backend.HandleRequest(context.Background(), &logical.Request{
			Operation:   logical.ReadOperation,
			Storage:     tc.RegularStorage(),
			ClientToken: "a",
			Path:        path,
		})
		require.NoError(t, err)
		require.NotNil(t, readResp)
		require.Equal(t, "old data", readResp.Data["data"])
	}
	require.NoError(t, err)
	addToken := func(req *logical.Request) {
		req.ClientToken = "a"
	}
	t.Run("read no side effects", func(t *testing.T) {
		tc.RunRead(t, path, snapshots.WithModifyRequests(func(req *logical.Request) {
			req.ClientToken = "token"
		}))
	})
	t.Run("list no side effects", func(t *testing.T) {
		tc.RunList(t, path, snapshots.WithModifyRequests(addToken))
	})
	t.Run("recover works", func(t *testing.T) {
		_, err := tc.DoRecover(t, path, snapshots.WithModifyRequests(addToken))
		require.NoError(t, err)
		checkContent(t, path)
	})
	t.Run("recover new path works", func(t *testing.T) {
		_, err := tc.DoRecover(t, "different-path", snapshots.WithModifyRequests(addToken), snapshots.WithRecoverSourcePath(path))
		require.NoError(t, err)
		checkContent(t, "different-path")
	})
}

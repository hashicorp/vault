// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func testCountActiveEntities(t *testing.T, c *Core, root string, expectedEntities int) {
	t.Helper()

	rootCtx := namespace.RootContext(nil)
	resp, err := c.HandleRequest(rootCtx, &logical.Request{
		ClientToken: root,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/counters/entities",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
	}

	if diff := deep.Equal(resp.Data, map[string]interface{}{
		"counters": &ActiveEntities{
			Entities: EntityCounter{
				Total: expectedEntities,
			},
		},
	}); diff != nil {
		t.Fatal(diff)
	}
}

func TestIdentityStore_CountActiveEntities(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	rootCtx := namespace.RootContext(nil)

	// Count the root token
	testCountActiveEntities(t, c, root, 0)

	// Create some entities
	req := &logical.Request{
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Path:        "entity",
	}
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		resp, err := c.identityStore.HandleRequest(rootCtx, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
		}
		ids[i] = resp.Data["id"].(string)

		testCountActiveEntities(t, c, root, i+1)
	}

	req.Operation = logical.DeleteOperation
	for i := 0; i < 10; i++ {
		req.Path = "entity/id/" + ids[i]
		resp, err := c.identityStore.HandleRequest(rootCtx, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
		}

		testCountActiveEntities(t, c, root, 9-i)
	}
}

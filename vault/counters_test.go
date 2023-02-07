// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// noinspection SpellCheckingInspection
func testParseTime(t *testing.T, format, timeval string) time.Time {
	t.Helper()
	tm, err := time.Parse(format, timeval)
	if err != nil {
		t.Fatalf("Error parsing time %q: %v", timeval, err)
	}
	return tm
}

func testCountActiveTokens(t *testing.T, c *Core, root string) int {
	t.Helper()

	rootCtx := namespace.RootContext(nil)
	resp, err := c.HandleRequest(rootCtx, &logical.Request{
		ClientToken: root,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/counters/tokens",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
	}

	activeTokens := resp.Data["counters"].(*ActiveTokens)
	return activeTokens.ServiceTokens.Total
}

func TestTokenStore_CountActiveTokens(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	rootCtx := namespace.RootContext(nil)

	// Count the root token
	count := testCountActiveTokens(t, c, root)
	if count != 1 {
		t.Fatalf("expected %d tokens, not %d", 1, count)
	}

	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {

		// Create some service tokens
		req := &logical.Request{
			ClientToken: root,
			Operation:   logical.UpdateOperation,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"ttl": "1h",
			},
		}

		resp, err := c.HandleRequest(rootCtx, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
		}

		tokens[i] = resp.Auth.ClientToken

		count = testCountActiveTokens(t, c, root)
		if count != i+2 {
			t.Fatalf("expected %d tokens, not %d", i+2, count)
		}
	}

	// Revoke the service tokens
	for i := 0; i < 10; i++ {

		req := &logical.Request{
			ClientToken: root,
			Operation:   logical.UpdateOperation,
			Path:        "auth/token/revoke",
			Data: map[string]interface{}{
				"token": tokens[i],
			},
		}

		resp, err := c.HandleRequest(rootCtx, req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
		}
	}

	// We should now have only 1 token (the root token).  However, because
	// token deletion works by setting the TTL of the token to 0 and waiting
	// for it to get cleaned up by the expiration manager, occasionally we will
	// have to wait briefly for all the tokens to actually get deleted.
	for i := 0; i < 10; i++ {
		count = testCountActiveTokens(t, c, root)
		if count == 1 {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("expected %d tokens, not %d", 1, count)
}

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

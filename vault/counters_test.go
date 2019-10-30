package vault

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

//noinspection SpellCheckingInspection
func testParseTime(t *testing.T, format, timeval string) time.Time {
	t.Helper()
	tm, err := time.Parse(format, timeval)
	if err != nil {
		t.Fatalf("Error parsing time '%s': %v", timeval, err)
	}
	return tm
}

// TestRequestCounterLoadCurrent exercises the code that primes the in-mem
// request counters from persistent storage.
func TestRequestCounterLoadCurrent(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	december2018 := testParseTime(t, time.RFC3339, "2018-12-05T09:44:12-05:00")
	decemberRequests := uint64(555)

	// It's December, and we got some requests.  Persist the counter.
	atomic.StoreUint64(c.counters.requests, decemberRequests)
	err := c.saveCurrentRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}

	// It's still December, simulate being restarted.  At startup the counter is
	// zero initially, until we read the counter from storage post-unseal via
	// loadCurrentRequestCounters.
	atomic.StoreUint64(c.counters.requests, 0)
	err = c.loadCurrentRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadUint64(c.counters.requests); got != decemberRequests {
		t.Fatalf("expected=%d, got=%d", decemberRequests, got)
	}

	// Now simulate being restarted in January. We never wrote anything out during
	// January, so the in-mem counter should remain zero.
	january2019 := testParseTime(t, time.RFC3339, "2019-01-02T08:21:11-05:00")
	atomic.StoreUint64(c.counters.requests, 0)
	err = c.loadCurrentRequestCounters(context.Background(), january2019)
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadUint64(c.counters.requests); got != 0 {
		t.Fatalf("expected=%d, got=%d", 0, got)
	}
}

// TestRequestCounterSaveCurrent exercises the code that saves the in-mem
// request counters to persistent storage.
func TestRequestCounterSaveCurrent(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// storeSaveLoad stores newValue in the in-mem counter, saves it to storage,
	// then verifies in-mem counter has value expectedPostLoad.
	storeSaveLoad := func(newValue, expectedPostLoad uint64, now time.Time) {
		t.Helper()
		atomic.StoreUint64(c.counters.requests, newValue)
		err := c.saveCurrentRequestCounters(context.Background(), now)
		if err != nil {
			t.Fatal(err)
		}
		if got := atomic.LoadUint64(c.counters.requests); got != expectedPostLoad {
			t.Fatalf("expected=%d, got=%d", expectedPostLoad, got)
		}
	}

	// Start in December. The first write ever should persist the current in-mem value.
	december2018 := testParseTime(t, time.RFC3339, "2018-12-05T09:44:12-05:00")
	decemberRequests := uint64(555)
	storeSaveLoad(decemberRequests, decemberRequests, december2018)

	// Update request count.
	decemberRequests++
	storeSaveLoad(decemberRequests, decemberRequests, december2018)

	decemberStartTime := testParseTime(t, requestCounterDatePathFormat, december2018.Format(requestCounterDatePathFormat))
	expected2018 := []DatedRequestCounter{
		{StartTime: decemberStartTime, RequestCounter: RequestCounter{Total: &decemberRequests}},
	}
	all, err := c.loadAllRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(all, expected2018); len(diff) != 0 {
		t.Errorf("Expected=%v, got=%v, diff=%v", expected2018, all, diff)
	}

	// Now it's January. Saving after transition to new month should reset in-mem
	// counter to zero, and also write zero to storage for the new month.
	january2019 := testParseTime(t, time.RFC3339, "2019-01-02T08:21:11-05:00")
	decemberRequests += 5
	storeSaveLoad(decemberRequests, 0, january2019)

	januaryRequests := uint64(333)
	storeSaveLoad(januaryRequests, januaryRequests, january2019)

	all, err = c.loadAllRequestCounters(context.Background(), january2019)
	if err != nil {
		t.Fatal(err)
	}

	januaryStartTime := testParseTime(t, requestCounterDatePathFormat, january2019.Format(requestCounterDatePathFormat))
	expected2019 := expected2018
	expected2019 = append(expected2019,
		DatedRequestCounter{januaryStartTime, RequestCounter{&januaryRequests}})
	if diff := deep.Equal(all, expected2019); len(diff) != 0 {
		t.Errorf("Expected=%v, got=%v, diff=%v", expected2019, all, diff)
	}
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

	// Create some service tokens
	req := &logical.Request{
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Path:        "create",
		Data: map[string]interface{}{
			"ttl": "1h",
		},
	}
	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {
		resp, err := c.tokenStore.HandleRequest(rootCtx, req)
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
	req.Path = "revoke"
	req.Data = make(map[string]interface{})
	for i := 0; i < 10; i++ {
		req.Data["token"] = tokens[i]
		resp, err := c.tokenStore.HandleRequest(rootCtx, req)
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

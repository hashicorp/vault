// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/queue"
)

// some tests rely on the ordering of items from this method
func testCases() (tc []*MFACachedAuthResponse) {
	// create a slice of items with times offest by these seconds
	for _, m := range []time.Duration{
		5,
		183600,  // 51 hours
		15,      // 15 seconds
		45,      // 45 seconds
		900,     // 15 minutes
		360,     // 6 minutes
		7200,    // 2 hours
		183600,  // 51 hours
		7201,    // 2 hours, 1 second
		115200,  // 32 hours
		1209600, // 2 weeks
	} {
		n := time.Now()
		ft := n.Add(time.Second * m)
		uid, err := uuid.GenerateUUID()
		if err != nil {
			continue
		}
		tc = append(tc, &MFACachedAuthResponse{
			TimeOfStorage: ft,
			RequestID:     uid,
		})
	}
	return
}

func TestLoginMFAPriorityQueue_PushPopByKey(t *testing.T) {
	pq := NewLoginMFAPriorityQueue()

	if pq.Len() != 0 {
		t.Fatalf("expected new queue to have zero size, got (%d)", pq.Len())
	}

	tc := testCases()
	tcl := len(tc)
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	if pq.Len() != tcl {
		t.Fatalf("error adding items, expected (%d) items, got (%d)", tcl, pq.Len())
	}

	item, err := pq.PopByKey(tc[0].RequestID)
	if err != nil {
		t.Fatalf("error popping item: %s", err)
	}
	if tc[0].TimeOfStorage != item.TimeOfStorage {
		t.Fatalf("expected tc[0] and popped item to match, got (%v) and (%v)", tc[0].TimeOfStorage, item.TimeOfStorage)
	}

	// push item with duplicate key
	dErr := pq.Push(tc[1])
	if dErr != queue.ErrDuplicateItem {
		t.Fatal(err)
	}
	// push item with no key
	tc[2].RequestID = ""
	kErr := pq.Push(tc[2])
	if kErr != nil && kErr.Error() != "error adding item: Item Key is required" {
		t.Fatal(kErr)
	}

	// check nil,nil error for not found
	i, err := pq.PopByKey("empty")
	if err != nil && i != nil {
		t.Fatalf("expected nil error for PopByKey of non-existing key, got: %s", err)
	}
}

func TestLoginMFARemoveStaleEntries(t *testing.T) {
	pq := NewLoginMFAPriorityQueue()

	tc := testCases()
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	cutoffTime := time.Now().Add(371 * time.Second)
	timeout := time.Now().Add(5 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		pq.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, cutoffTime)
	}

	if pq.Len() != 8 {
		t.Fatalf("failed to remove %d stale entries", pq.Len())
	}
}

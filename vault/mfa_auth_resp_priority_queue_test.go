// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

// TestLoginMFAPriorityQueue_PeekByKey tests the PeekByKey method of
// LoginMFAPriorityQueue. It verifies that PeekByKey returns the correct item
// without removing it from the queue, returns errors on unhappy paths, handles
// non-existing keys appropriately, works correctly on empty queues, and
// properly handles empty key strings.
func TestLoginMFAPriorityQueue_PeekByKey(t *testing.T) {
	pq := NewLoginMFAPriorityQueue()
	tc := testCases()
	expectedLength := len(tc)

	// Peek from empty queue
	peekedItem, err := pq.PeekByKey("item-2")
	if peekedItem != nil {
		t.Fatal("expected nil when peeking from empty queue, got item")
	}
	if err == nil {
		t.Fatal("expected an error when peeking from empty queue, got nil")
	}
	if pq.Len() != 0 {
		t.Fatalf("expected empty queue to remain size 0, got %d", pq.Len())
	}

	// Push test items
	for _, item := range tc {
		if err := pq.Push(item); err != nil {
			t.Fatal(err)
		}
	}

	// Peek with empty key
	peekedItem, err = pq.PeekByKey("")
	if peekedItem != nil {
		t.Fatal("expected nil for empty key, got item")
	}
	if err == nil {
		t.Fatal("expected error when peeking with empty key , got nil")
	}
	// Verify queue size unchanged
	if pq.Len() != expectedLength {
		t.Fatalf("expected queue size to remain %d, got %d", expectedLength, pq.Len())
	}

	// Peek at non-existing item
	peekedItem, err = pq.PeekByKey("non-existing-key")
	if peekedItem != nil {
		t.Fatal("expected nil for non-existing key, got item")
	}
	if err == nil {
		t.Fatal("expected error when peeking with non-existing key, got nil")
	}
	// Verify queue size unchanged
	if pq.Len() != expectedLength {
		t.Fatalf("expected queue size to remain %d, got %d", expectedLength, pq.Len())
	}

	// Peek at a specific item
	peekedItem, err = pq.PeekByKey(tc[2].RequestID)
	if peekedItem == nil {
		t.Fatal("expected to peek item-2, got nil")
	}
	if err != nil {
		t.Fatal("expected no error when peeking existing key, got", err)
	}
	if peekedItem.RequestID != tc[2].RequestID {
		t.Fatal("expected the same item on subsequent peeks, got different items")
	}
	// Verify queue size unchanged
	if pq.Len() != expectedLength {
		t.Fatalf("expected queue size to remain %d, got %d", expectedLength, pq.Len())
	}
	// Verify item still exists in queue
	stillExists, err := pq.PeekByKey(tc[2].RequestID)
	if stillExists == nil {
		t.Fatal("item should still exist after peek")
	}
	if err != nil {
		t.Fatal("expected no error when peeking existing key for the second time, got", err)
	}
	if stillExists.RequestID != tc[2].RequestID {
		t.Fatal("expected the same item on subsequent peeks, got different items")
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

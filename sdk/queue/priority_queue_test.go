// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package queue

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

// Ensure we satisfy the heap.Interface
var _ heap.Interface = &queue{}

// some tests rely on the ordering of items from this method
func testCases() (tc []*Item) {
	// create a slice of items with priority / times offest by these seconds
	for i, m := range []time.Duration{
		5,
		183600,  // 51 hours
		15,      // 15 seconds
		45,      // 45 seconds
		900,     // 15 minutes
		300,     // 5 minutes
		7200,    // 2 hours
		183600,  // 51 hours
		7201,    // 2 hours, 1 second
		115200,  // 32 hours
		1209600, // 2 weeks
	} {
		n := time.Now()
		ft := n.Add(time.Second * m)
		tc = append(tc, &Item{
			Key:      fmt.Sprintf("item-%d", i),
			Value:    1,
			Priority: ft.Unix(),
		})
	}
	return
}

func TestPriorityQueue_New(t *testing.T) {
	pq := New()

	if len(pq.data) != len(pq.dataMap) || len(pq.data) != 0 {
		t.Fatalf("error in queue/map size, expected data and map to be initialized, got (%d) and (%d)", len(pq.data), len(pq.dataMap))
	}

	if pq.Len() != 0 {
		t.Fatalf("expected new queue to have zero size, got (%d)", pq.Len())
	}
}

func TestPriorityQueue_Push(t *testing.T) {
	pq := New()

	// don't allow nil pushing
	if err := pq.Push(nil); err == nil {
		t.Fatal("Expected error on pushing nil")
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

	testValidateInternalData(t, pq, len(tc), false)

	item, err := pq.Pop()
	if err != nil {
		t.Fatalf("error popping item: %s", err)
	}
	if tc[0].Priority != item.Priority {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], item.Priority)
	}
	if tc[0].Key != item.Key {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], item.Priority)
	}

	testValidateInternalData(t, pq, len(tc)-1, false)

	// push item with no key
	dErr := pq.Push(tc[1])
	if dErr != ErrDuplicateItem {
		t.Fatal(err)
	}
	// push item with no key
	tc[2].Key = ""
	kErr := pq.Push(tc[2])
	if kErr != nil && kErr.Error() != "error adding item: Item Key is required" {
		t.Fatal(kErr)
	}

	testValidateInternalData(t, pq, len(tc)-1, true)

	// check nil,nil error for not found
	i, err := pq.PopByKey("empty")
	if err != nil && i != nil {
		t.Fatalf("expected nil error for PopByKey of non-existing key, got: %s", err)
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	pq := New()

	tc := testCases()
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	topItem, err := pq.Pop()
	if err != nil {
		t.Fatalf("error calling pop: %s", err)
	}
	if tc[0].Priority != topItem.Priority {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], topItem.Priority)
	}
	if tc[0].Key != topItem.Key {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], topItem.Priority)
	}

	var items []*Item
	items = append(items, topItem)
	// pop the remaining items, compare size of input and output
	i, _ := pq.Pop()
	for ; i != nil; i, _ = pq.Pop() {
		items = append(items, i)
	}
	if len(items) != len(tc) {
		t.Fatalf("expected popped item count to match test cases, got (%d)", len(items))
	}
}

func TestPriorityQueue_PopByKey(t *testing.T) {
	pq := New()

	tc := testCases()
	for _, i := range tc {
		if err := pq.Push(i); err != nil {
			t.Fatal(err)
		}
	}

	// grab the top priority item, to capture it's value for checking later
	item, _ := pq.Pop()
	oldPriority := item.Priority
	oldKey := item.Key

	// push the item back on, so it gets removed with PopByKey and we verify
	// the top item has changed later
	err := pq.Push(item)
	if err != nil {
		t.Fatalf("error re-pushing top item: %s", err)
	}

	popKeys := []int{2, 4, 7, 1, 0}
	for _, i := range popKeys {
		item, err := pq.PopByKey(fmt.Sprintf("item-%d", i))
		if err != nil {
			t.Fatalf("failed to pop item-%d, \n\terr: %s\n\titem: %#v", i, err, item)
		}
	}

	testValidateInternalData(t, pq, len(tc)-len(popKeys), false)

	// grab the top priority item again, to compare with the top item priority
	// from above
	item, _ = pq.Pop()
	newPriority := item.Priority
	newKey := item.Key

	if oldPriority == newPriority || oldKey == newKey {
		t.Fatalf("expected old/new key and priority to differ, got (%s/%s) and (%d/%d)", oldKey, newKey, oldPriority, newPriority)
	}

	testValidateInternalData(t, pq, len(tc)-len(popKeys)-1, true)
}

// testValidateInternalData checks the internal data structure of the PriorityQueue
// and verifies that items are in-sync. Use drain only at the end of a test,
// because it will mutate the input queue
func testValidateInternalData(t *testing.T, pq *PriorityQueue, expectedSize int, drain bool) {
	actualSize := pq.Len()
	if actualSize != expectedSize {
		t.Fatalf("expected new queue size to be (%d), got (%d)", expectedSize, actualSize)
	}

	if len(pq.data) != len(pq.dataMap) || len(pq.data) != expectedSize {
		t.Fatalf("error in queue/map size, expected data and map to be (%d), got (%d) and (%d)", expectedSize, len(pq.data), len(pq.dataMap))
	}

	if drain && pq.Len() > 0 {
		// pop all the items, verify lengths
		i, _ := pq.Pop()
		for ; i != nil; i, _ = pq.Pop() {
			expectedSize--
			if len(pq.data) != len(pq.dataMap) || len(pq.data) != expectedSize {
				t.Fatalf("error in queue/map size, expected data and map to be (%d), got (%d) and (%d)", expectedSize, len(pq.data), len(pq.dataMap))
			}
		}
	}
}

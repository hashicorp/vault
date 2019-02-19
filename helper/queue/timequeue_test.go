package queue

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// Compile time test to enforce TimeQueue satisfies the heap.Interface interface
var _ PriorityQueue = &TimeQueue{}

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

func TestNewTimeQueue(t *testing.T) {
	tq := NewTimeQueue()

	tqi := tq.(*TimeQueue)
	if len(tqi.data) != len(tqi.dataMap) || len(tqi.data) != 0 {
		t.Fatalf("error in queue/map size, expected data and map to be initialized, got (%d) and (%d)", len(tqi.data), len(tqi.dataMap))
	}

	if tq.Len() != 0 {
		t.Fatalf("expected new queue to have zero size, got (%d)", tq.Len())
	}
}

func TestTimeQueue_PushItem(t *testing.T) {
	tq := NewTimeQueue()

	tc := testCases()
	tcl := len(tc)
	for _, i := range tc {
		if err := tq.PushItem(i); err != nil {
			t.Fatal(err)
		}
	}

	if tq.Len() != tcl {
		t.Fatalf("error adding items, expected (%d) items, got (%d)", tcl, tq.Len())
	}

	testValidateInternalData(t, tq, len(tc), false)

	item, err := tq.PopItem()
	if err != nil {
		t.Fatalf("error popping item: %s", err)
	}
	if tc[0].Priority != item.Priority {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], item.Priority)
	}
	if !reflect.DeepEqual(tc[0], item) {
		t.Fatal("expected test case and popped item match")
	}

	testValidateInternalData(t, tq, len(tc)-1, true)
}

func TestTimeQueue_PopItem(t *testing.T) {
	tq := NewTimeQueue()

	tc := testCases()
	for _, i := range tc {
		if err := tq.PushItem(i); err != nil {
			t.Fatal(err)
		}
	}

	topItem, err := tq.PopItem()
	if err != nil {
		t.Fatalf("error calling pop: %s", err)
	}
	if tc[0].Priority != topItem.Priority {
		t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], topItem.Priority)
	}
	if !reflect.DeepEqual(tc[0], topItem) {
		t.Fatal("expected test case and popped item match")
	}

	var items []*Item
	items = append(items, topItem)
	// pop the remaining items, compare size of input and output
	it, _ := tq.PopItem()
	for ; it != nil; it, _ = tq.PopItem() {
		items = append(items, it)
	}
	if len(items) != len(tc) {
		t.Fatalf("expected popped item count to match test cases, got (%d)", len(items))
	}

	testValidateInternalData(t, tq, 0, true)
}

func TestTimeQueue_PopItemByKey(t *testing.T) {
	tq := NewTimeQueue()

	tc := testCases()
	for _, i := range tc {
		if err := tq.PushItem(i); err != nil {
			t.Fatal(err)
		}
	}

	// grab the top priority item, to capture it's value for checking later
	item, _ := tq.PopItem()
	oldPriority := item.Priority
	oldKey := item.Key

	// push the item back on, so it gets removed with PopItemByKey and we verify
	// the top item has changed later
	err := tq.PushItem(item)
	if err != nil {
		t.Fatalf("error re-pushing top item: %s", err)
	}

	popKeys := []int{2, 4, 7, 1, 0}
	for _, i := range popKeys {
		item, err := tq.PopItemByKey(fmt.Sprintf("item-%d", i))
		if err != nil || item == nil {
			t.Fatalf("failed to pop item-%d, \n\terr: %s\n\titem: %#v", i, err, item)
		}
	}

	testValidateInternalData(t, tq, len(tc)-len(popKeys), false)

	// grab the top priority item again, to compare with the top item priority
	// from above
	item, _ = tq.PopItem()
	newPriority := item.Priority
	newKey := item.Key

	if oldPriority == newPriority || oldKey == newKey {
		t.Fatalf("expected old/new key and priority to differ, got (%s/%s) and (%d/%d)", oldKey, newKey, oldPriority, newPriority)
	}

	testValidateInternalData(t, tq, len(tc)-len(popKeys)-1, true)
}

// testValidateInternalData checks the internal data stucture of the TimeQueue
// and verifies that items are in-sync. Use drain only at the end of a test,
// because it will mutate the input queue
func testValidateInternalData(t *testing.T, pq PriorityQueue, expectedSize int, drain bool) {
	actualSize := pq.Len()
	if actualSize != expectedSize {
		t.Fatalf("expected new queue size to be (%d), got (%d)", expectedSize, actualSize)
	}

	tq := pq.(*TimeQueue)
	if len(tq.data) != len(tq.dataMap) || len(tq.data) != expectedSize {
		t.Fatalf("error in queue/map size, expected data and map to be (%d), got (%d) and (%d)", expectedSize, len(tq.data), len(tq.dataMap))
	}

	if drain && tq.Len() > 0 {
		// pop all the items, verify lengths
		i, _ := tq.PopItem()
		for ; i != nil; i, _ = tq.PopItem() {
			expectedSize--
			if len(tq.data) != len(tq.dataMap) || len(tq.data) != expectedSize {
				t.Fatalf("error in queue/map size, expected data and map to be (%d), got (%d) and (%d)", expectedSize, len(tq.data), len(tq.dataMap))
			}
		}
	}
}

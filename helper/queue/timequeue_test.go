package queue

import (
        "container/heap"
        "fmt"
        "reflect"
        "testing"
        "time"
)

// Compile time test to enforce TimeQueue satisfies the heap.Interface interface
var _ PriorityQueue = &TimeQueue{}

var secondOffsets = []time.Duration{
        0,       // some tests rely on the first test case being the highest priority
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
}

func testCases() []*Item {
        tc := make([]*Item, len(secondOffsets))
        // create a slice of items with priority / times offest by these seconds
        for i, m := range secondOffsets {
                n := time.Now()
                ft := n.Add(time.Second * m)
                tc[i] = &Item{
                        Key:      fmt.Sprintf("item-%d", i),
                        Value:    1,
                        Priority: ft.Unix(),
                }
        }
        // fmt.Println("test cases")
        // for i, t := range tc {
        // 	fmt.Printf("\t %d - %s\n", i, time.Unix(t.Priority, 0).String())
        // }
        return tc
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

        // check the internal data structures
        tqi := tq.(*TimeQueue)
        if len(tqi.data) != len(tqi.dataMap) {
                t.Fatalf("error in queue/map size, expected data and map to be initialized, got (%d) and (%d)", len(tqi.data), len(tqi.dataMap))
        }

        item := heap.Pop(tq).(*Item)
        if tc[0].Priority != item.Priority {
                t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], item.Priority)
        }
        if !reflect.DeepEqual(tc[0], item) {
                t.Fatal("expected test case and popped item match")
        }

}

func TestTimeQueue_PopItem(t *testing.T) {
        tq := NewTimeQueue()

        tc := testCases()
        for _, i := range tc {
                if err := tq.PushItem(i); err != nil {
                        t.Fatal(err)
                }
        }

        topItem := tq.PopItem()
        if tc[0].Priority != topItem.Priority {
                t.Fatalf("expected tc[0] and popped item to match, got (%q) and (%q)", tc[0], topItem.Priority)
        }
        if !reflect.DeepEqual(tc[0], topItem) {
                t.Fatal("expected test case and popped item match")
        }

        var items []*Item
        items = append(items, topItem)
        // pop the remaining items, compare size of input and output
        it := tq.PopItem()
        for ; it != nil; it = tq.PopItem() {
                items = append(items, it)
        }
        if len(items) != len(tc) {
                t.Fatalf("expected popped item count to match test cases, got (%d)", len(items))
        }

        tqi := tq.(*TimeQueue)
        if len(tqi.data) != len(tqi.dataMap) || len(tqi.data) != 0 {
                t.Fatalf("error in queue/map size, expected data and map to be zero, got (%d) and (%d)", len(tqi.data), len(tqi.dataMap))
        }
}

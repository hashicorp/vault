package queue

import (
        "testing"
        "time"
)

// Compile time test to enforce TimeQueue satisfies the heap.Interface interface
var _ PriorityQueue = &TimeQueue{}

func TestNew(t *testing.T) {

        secondOffsets := []time.Duration{
                0,
                15,
                45,
                300,  // 5 minutes
                900,  // 15 minutes
                7200, // 2 hours
        }

        testCases := make([]*Item, len(secondOffsets))

        // create a slice of items with priority / times offest by these seconds
        for i, m := range secondOffsets {
                n := time.Now()
                ft := n.Add(time.Second * m)
                testCases[i] = &Item{
                        value:    1,
                        priority: ft.Unix(),
                }
        }

        tq := NewTimeQueue(testCases)
        if tq.Size() != len(testCases) {
                t.Fatalf("error in queue size, expected (%d), got (%d)", len(testCases), tq.Size())
        }
}

package queue

import (
        "testing"
)

// Compile time test to enforce TimeQueue satisfies the heap.Interface interface
var _ PriorityQueue = &TimeQueue{}

func TestHelloWorld(t *testing.T) {
        // t.Fatal("not implemented")
}

package priorityqueue

import (
        "testing"
)

// Compile time test to enforce TimeQueue satisfies the heap.Interface interface
var _ Queue = &TimeQueue{}

func TestHelloWorld(t *testing.T) {
        // t.Fatal("not implemented")
}

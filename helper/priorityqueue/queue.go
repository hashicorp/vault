package priorityqueue

import "container/heap"

// Queue interface defines what a Queue must include, which is satisfying
// heap.Interface, and a few additional methods
type Queue interface {
        heap.Interface
        Peek()
        Pluck()
        Find()
}

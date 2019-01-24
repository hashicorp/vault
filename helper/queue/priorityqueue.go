package queue

import "container/heap"

// PriorityQueue interface defines what a Queue must include, which is satisfying
// heap.Interface, and a few additional methods
// TODO: refactor to be just Queue, or add a generic Queue type that implements
// the methods except Less() so each Queue type can do it's own sorting of
// priority
type PriorityQueue interface {
        heap.Interface
        Peek()
        Pluck()
        Find()
}

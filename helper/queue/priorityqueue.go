package queue

import (
	"container/heap"
	"fmt"
)

// PriorityQueue interface defines what a Queue must include, which is satisfying
// heap.Interface, and a few additional methods
// TODO: refactor to be just Queue, or add a generic Queue type that implements
// the methods except Less() so each Queue type can do it's own sorting of
// priority
type PriorityQueue interface {
	heap.Interface

	// PopItem pops the highest priority item from the queue. This is a
	// wrapper/convienence method that calls heap.Pop, so consumers do not need to
	// invoke heap functions directly
	PopItem() (*Item, error)

	// PushItem pushes an item on to the queue. This is a wrapper/convienence
	// method that calls heap.Push, so consumers do not need to invoke heap
	// functions directly
	PushItem(*Item) error

	// PopItemByKey searchs the queue for an item with the given key and removes it
	// from the queue if found. Returns ErrItemNotFound(key) if not found. This
	// method must fix the queue after removal.
	PopItemByKey(string) (*Item, error)

	// // Updates an item in the queue. This must call heap.Fix()
	// UpdateItem(*Item) error

	// // Peek returns the highest priority item, but does not remove it from the
	// // queue
	// Peek() (*Item, error)

	// // Find searches and returns item from the queue, if found. This does not
	// // remove the item. If no item is found, returns ErrItemNotFound
	// Find(string) (*Item, error)
}

// ErrItemNotFound creates a "not found" error for the given key
func ErrItemNotFound(key string) error {
	return fmt.Errorf("queue item with key (%s) not found", key)
}

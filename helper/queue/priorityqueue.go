package queue

import (
        "container/heap"
        "errors"
)

// ErrNoItemFound is used in Pluck and Find, to indicate a matching item was not
// found
var ErrNoItemFound = errors.New("item not found in queue")

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
        PopItem() *Item

        // PushItem pushes an item on to the queue. This is a wrapper/convienence
        // method that calls heap.Push, so consumers do not need to invoke heap
        // functions directly
        PushItem(*Item) error

        // Peek returns the top item from the queue, but does not remove it
        // Peek()

        // Pluck removes an item from the queue by the given Key. Pluck must fix the
        // queue after removal. If no item is removed, returns ErrNoItemFound
        Pluck(string) (*Item, error)

        // Find searches and returns item from the queue, if found. This does not
        // remove the item. If no item is found, returns ErrNoItemFound
        Find(string) (*Item, error)

        // Size reports the number of items in the queue
        // Size() int
}

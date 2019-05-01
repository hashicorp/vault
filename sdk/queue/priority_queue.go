package queue

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"

	"github.com/mitchellh/copystructure"
)

// ErrEmpty is returned for queues with no items
var ErrEmpty = errors.New("queue is empty")

// ErrDuplicateItem is returned when the queue attmepts to push an item to a key that
// already exists. The queue does not attempt to update, instead returns this
// error. If an Item needs to be updated or replaced, pop the item first.
var ErrDuplicateItem = errors.New("duplicate item")

// ErrItemNotFound is a struct that implements the error interface
var _ error = &ErrItemNotFound{}

type ErrItemNotFound struct {
	Key string
}

func (e *ErrItemNotFound) Error() string {
	return fmt.Sprintf("queue item with key (%s) not found", e.Key)
}

// NewErrItemNotFound creates a "not found" error for the given key
func NewErrItemNotFound(key string) *ErrItemNotFound {
	return &ErrItemNotFound{
		Key: key,
	}
}

// Item is something managed in the priority queue
type Item struct {
	// Key is a unique string used to identify items in the internal data map
	Key string
	// Value is an unspecified type that implementations can use to store
	// information
	Value interface{}

	// Priority determines ordering in the queue, with the lowest value being the
	// highest priority
	Priority int64

	// index is an internal value used by the heap package, and should not be
	// modified by any consumer of the priority queue
	index int
}

// NewPriorityQueue initializes the internal data structures and returns a new
// PriorityQueue
func NewPriorityQueue() *PriorityQueue {
	pq := PriorityQueue{
		data:    make([]*Item, 0),
		dataMap: make(map[string]*Item),
	}
	heap.Init(&pq)
	return &pq
}

// PriorityQueue satisfies heap.Interface and sort.Interface, and adds a few
// convenience methods. The ordering (priority) is an int64 value with the
// smallest value is the highest priority.
// See https://golang.org/pkg/container/heap/#example__priorityQueue
type PriorityQueue struct {
	// data is the internal structure that holds the queue, and is operated on by
	// heap functions
	data []*Item

	// dataMap represents all the items in the queue, with unique indexes, used
	// for finding specific items. dataMap is kept in sync with the data slice
	dataMap map[string]*Item

	// mapMutex is used to facilitate read/write locks on the dataMap
	dataMutex sync.Mutex
}

// PopItem pops the highest priority item from the queue. This is a
// wrapper/convenience method that calls heap.Pop, so consumers do not need to
// invoke heap functions directly
func (pq *PriorityQueue) PopItem() (*Item, error) {
	pq.dataMutex.Lock()
	defer pq.dataMutex.Unlock()

	if pq.Len() == 0 {
		return nil, ErrEmpty
	}
	item := heap.Pop(pq).(*Item)
	delete(pq.dataMap, item.Key)
	return item, nil
}

// PushItem pushes an item on to the queue. This is a wrapper/convenience
// method that calls heap.Push, so consumers do not need to invoke heap
// functions directly. Items must have unique Keys, and Items in the queue
// cannot be updated. To modify an Item, users must first remove it and re-push
// it after modifications
func (pq *PriorityQueue) PushItem(i *Item) error {
	if i.Key == "" {
		return errors.New("error adding item: Item Key is required")
	}

	pq.dataMutex.Lock()
	defer pq.dataMutex.Unlock()

	if _, ok := pq.dataMap[i.Key]; ok {
		return ErrDuplicateItem
	}

	// copy the item value(s) so that modifications to the source item does not
	// affect the item on the queue
	clone, err := copystructure.Copy(i)
	if err != nil {
		return err
	}

	pq.dataMap[i.Key] = clone.(*Item)
	heap.Push(pq, clone)
	return nil
}

// PopItemByKey searches the queue for an item with the given key and removes it
// from the queue if found. Returns ErrItemNotFound(key) if not found. This
// method must fix the queue after removal.
func (pq *PriorityQueue) PopItemByKey(key string) (*Item, error) {
	pq.dataMutex.Lock()
	defer pq.dataMutex.Unlock()

	item, ok := pq.dataMap[key]
	if !ok {
		return nil, NewErrItemNotFound(key)
	}

	// remove the item the heap and delete it from the dataMap
	itemRaw := heap.Remove(pq, item.index)
	delete(pq.dataMap, key)

	if i, ok := itemRaw.(*Item); ok {
		return i, nil
	}

	return nil, NewErrItemNotFound(key)
}

// Len returns the count of items in the queue.
func (pq *PriorityQueue) Len() int { return len(pq.data) }

// Less returns the less of two things, which in this case, we return the
// highest thing, because the priority ordering is "lowest value, highest
// priority"
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.data[i].Priority < pq.data[j].Priority
}

// Swap swaps things in-place
func (pq *PriorityQueue) Swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.data[i].index = i
	pq.data[j].index = j
}

// Push is used by heap.Interface to push items onto the heap. Do not use this
// method to add items to a queue: use PushItem instead.
// See: https://golang.org/pkg/container/heap/#Interface
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(pq.data)
	item := x.(*Item)
	item.index = n
	pq.data = append(pq.data, item)
}

// Pop is used by heap.Interface to pop items off of the heap. Do not use this
// method to remove items from a queue: use PopItem instead.
// See: https://golang.org/pkg/container/heap/#Interface
func (pq *PriorityQueue) Pop() interface{} {
	old := pq.data
	n := len(old)
	item := old[n-1]
	item.index = -1 //for saftey
	pq.data = old[0 : n-1]
	return item
}

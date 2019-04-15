package queue

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
)

// PriorityQueue interface defines what a Queue must include, which is satisfying
// heap.Interface, and a few additional methods
// TODO: refactor to be just Queue, or add a generic Queue type that implements
// the methods except Less() so each Queue type can do it's own sorting of
// priority
// type PriorityQueue interface {
// 	heap.Interface

// 	// PopItem pops the highest priority item from the queue. This is a
// 	// wrapper/convienence method that calls heap.Pop, so consumers do not need to
// 	// invoke heap functions directly
// 	PopItem() (*Item, error)

// 	// PushItem pushes an item on to the queue. This is a wrapper/convienence
// 	// method that calls heap.Push, so consumers do not need to invoke heap
// 	// functions directly
// 	PushItem(*Item) error

// 	// PopItemByKey searchs the queue for an item with the given key and removes it
// 	// from the queue if found. Returns ErrItemNotFound(key) if not found. This
// 	// method must fix the queue after removal.
// 	PopItemByKey(string) (*Item, error)

// 	// // Updates an item in the queue. This must call heap.Fix()
// 	// UpdateItem(*Item) error

// 	// // Peek returns the highest priority item, but does not remove it from the
// 	// // queue
// 	// Peek() (*Item, error)

// 	// // Find searches and returns item from the queue, if found. This does not
// 	// // remove the item. If no item is found, returns NewErrItemNotFound
// 	// Find(string) (*Item, error)
// }

var ErrEmpty = errors.New("queue is empty")
var ErrDuplicateItem = errors.New("duplicate item")

// NewPriorityQueue initializes a PriorityQueue struct for use as a PriorityQueue
func NewPriorityQueue() *PriorityQueue {
	tq := PriorityQueue{
		data:    make([]*Item, 0),
		dataMap: make(map[string]*Item),
	}
	heap.Init(&tq)
	return &tq
}

// PriorityQueue is a priority queue who's ordering is determined by the time in
// Unix of the item (nanoseconds elapsed since January 1, 1970 UTC)
type PriorityQueue struct {
	// data is the internal structure that holds the queue, and is operated on by
	// heap functions
	data []*Item

	// dataMap represents all the items in the queue, with unique indexes, used
	// for finding specific items. dataMap must be kept in sync with data
	dataMap map[string]*Item

	// mapMutex is used to facilitate read/write locks on the dataMap
	dataMutex sync.Mutex
}

// PopItem pops the highest priority item from the queue. This is a
// wrapper/convenience method that calls heap.Pop, so consumers do not need to
// invoke heap functions directly
func (tq *PriorityQueue) PopItem() (*Item, error) {
	tq.dataMutex.Lock()
	defer tq.dataMutex.Unlock()

	if tq.Len() == 0 {
		return nil, ErrEmpty
	}

	item := heap.Pop(tq).(*Item)
	delete(tq.dataMap, item.Key)
	return item, nil
}

// PushItem pushes an item on to the queue. This is a wrapper/convenience
// method that calls heap.Push, so consumers do not need to invoke heap
// functions directly
func (tq *PriorityQueue) PushItem(i *Item) error {
	if i.Key == "" {
		return errors.New("error adding item: Item Key is required")
	}

	tq.dataMutex.Lock()
	defer tq.dataMutex.Unlock()

	if _, ok := tq.dataMap[i.Key]; ok {
		return ErrDuplicateItem
	}

	tq.dataMap[i.Key] = i
	heap.Push(tq, i)
	return nil
}

// PopItemByKey removes an item from the queue by key, if found
func (tq *PriorityQueue) PopItemByKey(key string) (*Item, error) {
	tq.dataMutex.Lock()
	defer tq.dataMutex.Unlock()

	item, ok := tq.dataMap[key]
	if !ok {
		return nil, NewErrItemNotFound(key)
	}

	// remove the item the heap and delete it from the dataMap
	itemRaw := heap.Remove(tq, item.index)
	delete(tq.dataMap, key)

	if i, ok := itemRaw.(*Item); ok {
		return i, nil
	}

	return nil, NewErrItemNotFound(key)
}

//////
// begin heap.Interface methods
// TODO: extract into a generic Queue type, and other structs can embed new type
// and either override the Less method or supply a LessFunc or something
//////

// Len returns the count of items in the queue.
func (tq *PriorityQueue) Len() int { return len(tq.data) }

// Less returns the less of two things, which in this case, we return the
// highest thing.
func (tq *PriorityQueue) Less(i, j int) bool {
	return tq.data[i].Priority < tq.data[j].Priority
}

// Swap swaps things in-place
func (tq *PriorityQueue) Swap(i, j int) {
	tq.data[i], tq.data[j] = tq.data[j], tq.data[i]
	tq.data[i].index = i
	tq.data[j].index = j
}

// Push is used by heap.Interface to push items onto the heap. Do not use this
// to add items to a queue: use PushItem instead
func (tq *PriorityQueue) Push(x interface{}) {
	n := len(tq.data)
	item := x.(*Item)
	item.index = n
	tq.data = append(tq.data, item)
}

// Pop removes the highest priority thing
func (tq *PriorityQueue) Pop() interface{} {
	old := tq.data
	n := len(old)
	item := old[n-1]
	item.index = -1 //for saftey
	tq.data = old[0 : n-1]
	return item
}

//////
// end heap.Interface methods
//////

// NewErrItemNotFound creates a "not found" error for the given key
func NewErrItemNotFound(key string) error {
	return &ErrItemNotFound{
		Key: key,
	}
}

type ErrItemNotFound struct {
	Key string
}

func (e *ErrItemNotFound) Error() string {
	return fmt.Sprintf("queue item with key (%s) not found", e.Key)
}

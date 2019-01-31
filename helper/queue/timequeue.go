package queue

import (
        "container/heap"
        "errors"
        "fmt"
        "sync"
)

// NewTimeQueue initializes a TimeQueue struct for use as a PriorityQueue
func NewTimeQueue() PriorityQueue {
        tq := TimeQueue{
                data:    make([]*Item, 0),
                dataMap: make(map[string]*Item),
        }
        heap.Init(&tq)
        return &tq
}

// TimeQueue is a priority queue who's ordering is determined by the time in
// Unix of the item (nanoseconds elapsed since January 1, 1970 UTC)
type TimeQueue struct {
        // data is the internal structure that holds the queue, and is operated on by
        // heap functions
        data []*Item

        // dataMap represents all the items in the queue, with unique indexes, used
        // for finding specific items. dataMap must be kept in sync with data
        dataMap map[string]*Item

        // mapMutex is used to facilitate read/write locks on the dataMap
        dataMutex sync.RWMutex
}

// Peek returns the top priority item without removing it from the queue
func (tq *TimeQueue) Peek() {}

// Pluck removes an item from the queue by index. Plux must "fix" the heap when
// it's done
func (tq *TimeQueue) Pluck() {}

// Find searches the queue for an item by index and returns it if found,
// otherwise ErrNotFound
func (tq *TimeQueue) Find() {}

// Size reports the size of the queue, e.g. number of items in data
// TODO: duplicate of Len()?
func (tq *TimeQueue) Size() int {
        return len(tq.data)
}

// PopItem pops the highest priority item from the queue. This is a
// wrapper/convienence method that calls heap.Pop, so consumers do not need to
// invoke heap functions directly
func (tq *TimeQueue) PopItem() {
}

// PushItem pushes an item on to the queue. This is a wrapper/convienence
// method that calls heap.Push, so consumers do not need to invoke heap
// functions directly
func (tq *TimeQueue) PushItem(i *Item) error {
        if i.Key == "" {
                return errors.New("error adding item: Item Key is required")
        }

        tq.dataMutex.RLock()
        if _, ok := tq.dataMap[i.Key]; ok {
                tq.dataMutex.RUnlock()
                return fmt.Errorf("error adding item: Item already in queue. Use UpdateItem() instead")
        }
        tq.dataMutex.RUnlock()
        tq.dataMutex.Lock()
        defer tq.dataMutex.Unlock()
        tq.dataMap[i.Key] = i
        heap.Push(tq, i)
        return nil
}

// Update modifies the priority and value of an Item
func (tq *TimeQueue) Update(item *Item, value string, priority int64) {
        item.Value = value
        item.Priority = priority
        heap.Fix(tq, item.index)
}

//////
// begin heap.Interface methods
// TODO: extract into a generic Queue type, and other structs can embed new type
// and either override the Less method or supply a LessFunc or something
//////

// Len returns the count of items in the queue.
func (tq *TimeQueue) Len() int { return len(tq.data) }

// Less returns the less of two things, which in this case, we return the
// highest thing.
func (tq *TimeQueue) Less(i, j int) bool {
        return tq.data[i].Priority < tq.data[j].Priority
}

// Swap swaps things in-place
func (tq *TimeQueue) Swap(i, j int) {
        tq.data[i], tq.data[j] = tq.data[j], tq.data[i]
        tq.data[i].index = i
        tq.data[j].index = j
}

// Push is used by heap.Interface to push items onto the heap. Do not use this
// to add items to a queue: use PushItem instead
func (tq *TimeQueue) Push(x interface{}) {
        n := len(tq.data)
        item := x.(*Item)
        item.index = n
        tq.data = append(tq.data, item)
}

// Pop removes the highest priority thing
func (tq *TimeQueue) Pop() interface{} {
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

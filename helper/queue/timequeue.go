package queue

import (
        "container/heap"
        "time"
)

// TimeQueue is a priority queue who's ordering is determined by the time in
// nanoseconds of the item
type TimeQueue struct {
        // data is the internal structure that holds the queue, and is operated on by
        // heap functions
        data []*Item

        // dataMap represents all the items in the queue, with unique indexes, used
        // for finding specific items. dataMap must be kept in sync with data
        dataMap map[string]*Item
}

// Peek returns the top priority item without removing it from the queue
func (tq *TimeQueue) Peek() {}

// Pluck removes an item from the queue by index. Plux must "fix" the heap when
// it's done
func (tq *TimeQueue) Pluck() {}

// Find searches the queue for an item by index and returns it if found,
// otherwise ErrNotFound
func (tq *TimeQueue) Find() {}

// Item is something managed in the priority queue
type Item struct {
        value     string
        priority  int // priority of item in queue
        index     int // index is needed by update and maintained by heap package
        createdAt time.Time
}

//////
// begin heap.Interface methods
// TODO: extract into a generic Queue type
//////

// Len returns the count of items in the queue.
func (tq *TimeQueue) Len() int { return len(tq.data) }

// Less returns the less of two things, which in this case, we return the
// highest thing.
func (tq *TimeQueue) Less(i, j int) bool {
        // we want pop to give the highest, not lowest, priority
        // TODO: same priority?
        if tq.data[i].priority == tq.data[j].priority {
                return tq.data[j].createdAt.After(tq.data[i].createdAt)
        }
        return tq.data[i].priority > tq.data[j].priority
}

// Swap swaps things in-place
func (tq *TimeQueue) Swap(i, j int) {
        tq.data[i], tq.data[j] = tq.data[j], tq.data[i]
        tq.data[i].index = i
        tq.data[j].index = j
}

// Push pushes things
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

// update modifies the priority and value of an Item
func (tq *TimeQueue) update(item *Item, value string, priority int) {
        item.value = value
        item.priority = priority
        heap.Fix(tq, item.index)
}

package kheap

import (
	"container/heap"
)

// HeapItem is an interface type implemented by objects stored in the ScoreHeap
type HeapItem interface {
	Data() interface{} // The data object
	Score() float64    // Score to use as the sort criteria
}

// A ScoreHeap implements heap.Interface and is a min heap
// that keeps the top K elements by Score. Push can be called
// with an arbitrary number of values but only the top K are stored
type ScoreHeap struct {
	items    []HeapItem
	capacity int
}

func NewScoreHeap(capacity uint32) *ScoreHeap {
	return &ScoreHeap{capacity: int(capacity)}
}

func (pq ScoreHeap) Len() int { return len(pq.items) }

func (pq ScoreHeap) Less(i, j int) bool {
	return pq.items[i].Score() < pq.items[j].Score()
}

func (pq ScoreHeap) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push implements heap.Interface and only stores
// the top K elements by Score
func (pq *ScoreHeap) Push(x interface{}) {
	item := x.(HeapItem)
	if len(pq.items) < pq.capacity {
		pq.items = append(pq.items, item)
	} else {
		// Pop the lowest scoring element if this item's Score is
		// greater than the min Score so far
		minIndex := 0
		min := pq.items[minIndex]
		if item.Score() > min.Score() {
			// Replace min and heapify
			pq.items[minIndex] = item
			heap.Fix(pq, minIndex)
		}
	}
}

// Push implements heap.Interface and returns the top K scoring
// elements in increasing order of Score. Callers must reverse the order
// of returned elements to get the top K scoring elements in descending order
func (pq *ScoreHeap) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

// GetItemsReverse returns the items in this min heap in reverse order
// sorted by score descending
func (pq *ScoreHeap) GetItemsReverse() []interface{} {
	ret := make([]interface{}, pq.Len())
	i := pq.Len() - 1
	for pq.Len() > 0 {
		item := heap.Pop(pq)
		ret[i] = item
		i--
	}
	return ret
}

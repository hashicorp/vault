package clock

import (
	"container/heap"
	"time"
)

type mockTimer struct {
	deadline  time.Time
	fire      func() time.Duration
	mock      *Mock
	heapIndex int
}

const removed = -1

func newMockTimer(m *Mock, d time.Time) *mockTimer {
	return &mockTimer{
		deadline:  d,
		mock:      m,
		heapIndex: removed,
	}
}

func (t mockTimer) stopped() bool {
	return t.heapIndex == removed
}

// timerHeap implements mockTimers with a heap.
type timerHeap []*mockTimer

func (h timerHeap) Len() int { return len(h) }

func (h timerHeap) Less(i, j int) bool {
	return h[i].deadline.Before(h[j].deadline)
}

func (h timerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *timerHeap) Push(x interface{}) {
	n := len(*h)
	t := x.(*mockTimer)
	t.heapIndex = n
	*h = append(*h, t)
}

func (h *timerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	t := old[n-1]
	t.heapIndex = removed
	*h = old[0 : n-1]
	return t
}

func (h *timerHeap) start(t *mockTimer) {
	heap.Push(h, t)
}

func (h *timerHeap) stop(t *mockTimer) {
	if !t.stopped() {
		heap.Remove(h, t.heapIndex)
	}
}

func (h *timerHeap) reset(t *mockTimer) {
	if !t.stopped() {
		heap.Fix(h, t.heapIndex)
	} else {
		heap.Push(h, t)
	}
}

func (h timerHeap) next() *mockTimer {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}

func (h timerHeap) len() int {
	return h.Len()
}

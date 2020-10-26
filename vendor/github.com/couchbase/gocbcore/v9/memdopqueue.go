package gocbcore

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	errOpQueueClosed = errors.New("queue is closed")
	errOpQueueFull   = errors.New("queue is full")
	errAlreadyQueued = errors.New("request was already queued somewhere else")
)

type memdOpConsumer struct {
	parent   *memdOpQueue
	isClosed bool
}

func (c *memdOpConsumer) Queue() *memdOpQueue {
	return c.parent
}

func (c *memdOpConsumer) Pop() *memdQRequest {
	return c.parent.pop(c)
}

func (c *memdOpConsumer) Close() {
	c.parent.closeConsumer(c)
}

type memdOpQueue struct {
	lock   sync.Mutex
	signal *sync.Cond
	items  *list.List
	isOpen bool
}

func newMemdOpQueue() *memdOpQueue {
	q := memdOpQueue{
		isOpen: true,
		items:  list.New(),
	}
	q.signal = sync.NewCond(&q.lock)
	return &q
}

// nolint: unused
func (q *memdOpQueue) debugString() string {
	var outStr string
	q.lock.Lock()

	outStr += fmt.Sprintf("Num Items: %d\n", q.items.Len())
	outStr += fmt.Sprintf("Is Open: %t", q.isOpen)

	q.lock.Unlock()
	return outStr
}

func (q *memdOpQueue) Remove(req *memdQRequest) bool {
	q.lock.Lock()

	if !atomic.CompareAndSwapPointer(&req.queuedWith, unsafe.Pointer(q), nil) {
		q.lock.Unlock()
		return false
	}

	for e := q.items.Front(); e != nil; e = e.Next() {
		if e.Value.(*memdQRequest) == req {
			q.items.Remove(e)
			break
		}
	}

	q.lock.Unlock()

	return true
}

func (q *memdOpQueue) Push(req *memdQRequest, maxItems int) error {
	q.lock.Lock()
	if !q.isOpen {
		q.lock.Unlock()
		return errOpQueueClosed
	}

	if maxItems > 0 && q.items.Len() >= maxItems {
		q.lock.Unlock()
		return errOpQueueFull
	}

	if !atomic.CompareAndSwapPointer(&req.queuedWith, nil, unsafe.Pointer(q)) {
		q.lock.Unlock()
		return errAlreadyQueued
	}

	if req.isCancelled() {
		atomic.CompareAndSwapPointer(&req.queuedWith, unsafe.Pointer(q), nil)
		q.lock.Unlock()

		return errRequestCanceled
	}

	q.items.PushBack(req)
	q.lock.Unlock()

	q.signal.Broadcast()
	return nil
}

func (q *memdOpQueue) Consumer() *memdOpConsumer {
	return &memdOpConsumer{
		parent:   q,
		isClosed: false,
	}
}

func (q *memdOpQueue) closeConsumer(c *memdOpConsumer) {
	q.lock.Lock()
	c.isClosed = true
	q.lock.Unlock()

	q.signal.Broadcast()
}

func (q *memdOpQueue) pop(c *memdOpConsumer) *memdQRequest {
	q.lock.Lock()

	for q.isOpen && !c.isClosed && q.items.Len() == 0 {
		q.signal.Wait()
	}

	if !q.isOpen || c.isClosed {
		q.lock.Unlock()
		return nil
	}

	e := q.items.Front()
	q.items.Remove(e)

	req, ok := e.Value.(*memdQRequest)
	if !ok {
		logErrorf("Encountered incorrect type in memdOpQueue")
		return q.pop(c)
	}

	atomic.CompareAndSwapPointer(&req.queuedWith, unsafe.Pointer(q), nil)

	q.lock.Unlock()

	return req
}

type drainCallback func(*memdQRequest)

func (q *memdOpQueue) Drain(cb drainCallback) {
	q.lock.Lock()

	if q.isOpen {
		logErrorf("Attempted to Drain open memdOpQueue, ignoring")
		q.lock.Unlock()
		return
	}

	for e := q.items.Front(); e != nil; e = e.Next() {
		req, ok := e.Value.(*memdQRequest)
		if !ok {
			logErrorf("Encountered incorrect type in memdOpQueue")
			continue
		}

		atomic.CompareAndSwapPointer(&req.queuedWith, unsafe.Pointer(q), nil)

		cb(req)
	}

	q.lock.Unlock()
}

func (q *memdOpQueue) Close() {
	q.lock.Lock()
	q.isOpen = false
	q.lock.Unlock()

	q.signal.Broadcast()
}

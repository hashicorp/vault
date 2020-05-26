package gocbcore

import (
	"sync/atomic"

	"github.com/couchbase/gocbcore/v9/memd"
)

type memdOpMapItem struct {
	value *memdQRequest
	next  *memdOpMapItem
}

// This is used to store operations while they are pending
//  a response from the server to allow mapping of a response
//  opaque back to the originating request.  This queue takes
//  advantage of the monotonic nature of the opaque values
//  and synchronous responses from the server to nearly always
//  return the request without needing to iterate at all.
type memdOpMap struct {
	opIndex uint32

	first *memdOpMapItem
	last  *memdOpMapItem
}

// Add a new request to the bottom of the op queue.
func (q *memdOpMap) Add(req *memdQRequest) {
	q.opIndex++
	atomic.StoreUint32(&req.Opaque, q.opIndex)

	item := &memdOpMapItem{
		value: req,
		next:  nil,
	}

	if q.last == nil {
		q.first = item
		q.last = item
	} else {
		q.last.next = item
		q.last = item
	}
}

// Removes a request from the op queue.  Expects to be passed
//  the request to remove, along with the request that
//  immediately precedes it in the queue.
func (q *memdOpMap) remove(prev *memdOpMapItem, req *memdOpMapItem) {
	if prev == nil {
		q.first = req.next
		if q.first == nil {
			q.last = nil
		}
		return
	}
	prev.next = req.next
	if prev.next == nil {
		q.last = prev
	}
}

// Removes a specific request from the op queue.
func (q *memdOpMap) Remove(req *memdQRequest) bool {
	cur := q.first
	var prev *memdOpMapItem
	for cur != nil {
		if cur.value == req {
			q.remove(prev, cur)
			return true
		}
		prev = cur
		cur = cur.next
	}

	return false
}

// This allows searching through the list of requests for a specific
// request.  This is only used by the DCP server bug fix for MB-26363.
func (q *memdOpMap) FindOpenStream(vbID uint16) *memdQRequest {
	cur := q.first
	for cur != nil {
		if cur.value.Magic == memd.CmdMagicReq &&
			cur.value.Command == memd.CmdDcpStreamReq &&
			cur.value.Vbucket == vbID {
			return cur.value
		}
		cur = cur.next
	}

	return nil
}

// Locates a request (searching FIFO-style) in the op queue using
// the opaque value that was assigned to it when it was dispatched.
func (q *memdOpMap) Find(opaque uint32) *memdQRequest {
	cur := q.first
	for cur != nil {
		if cur.value.Opaque == opaque {
			return cur.value
		}
		cur = cur.next
	}

	return nil
}

// Locates a request (searching FIFO-style) in the op queue using
// the opaque value that was assigned to it when it was dispatched.
// It then removes the request from the queue if it is not persistent
// or if alwaysRemove is set to true.
func (q *memdOpMap) FindAndMaybeRemove(opaque uint32, force bool) *memdQRequest {
	cur := q.first
	var prev *memdOpMapItem
	for cur != nil {
		if cur.value.Opaque == opaque {
			if !cur.value.Persistent || force {
				q.remove(prev, cur)
			}

			return cur.value
		}
		prev = cur
		cur = cur.next
	}

	return nil
}

// Clears the queue of all requests and calls the passed function
// once for each request found in the queue.
func (q *memdOpMap) Drain(cb func(*memdQRequest)) {
	for cur := q.first; cur != nil; cur = cur.next {
		cb(cur.value)
	}
	q.first = nil
	q.last = nil
}

package gocbcore

import (
	"sync/atomic"

	"github.com/couchbase/gocbcore/v10/memd"
)

// memdOpMap - Uses the requests opaque to map requests to responses. Note that this structure is not thread safe, and
// uses should be guarded by a mutex.
type memdOpMap struct {
	opaque   uint32
	requests map[uint32]*memdQRequest
}

// newMemdOpMap - Creates a new empty 'memdOpMap' initializing any internal structures. Note that the requests opaque
// will begin at one and monotonically increase from there.
func newMemdOpMap() *memdOpMap {
	return &memdOpMap{requests: make(map[uint32]*memdQRequest)}
}

// Add - Add a new request to the map, the provided requests opaque value will be updated atomically.
func (m *memdOpMap) Add(req *memdQRequest) {
	m.opaque++
	atomic.StoreUint32(&req.Opaque, m.opaque)
	m.requests[m.opaque] = req
}

// Remove - Remove the provided request from the map.
func (m *memdOpMap) Remove(req *memdQRequest) bool {
	_, ok := m.requests[req.Opaque]
	delete(m.requests, req.Opaque)
	return ok
}

// FindOpenStream - This allows searching through the list of requests for a specific request. This is only used to fix
// the DCP server bug MB-26363.
func (m *memdOpMap) FindOpenStream(vbID uint16) *memdQRequest {
	for _, req := range m.requests {
		if req.Magic == memd.CmdMagicReq && req.Command == memd.CmdDcpStreamReq && req.Vbucket == vbID {
			return req
		}
	}

	return nil
}

// FindAndRemoveAllPersistent - Find all persistent requests, removing them from the map and returning them all.
func (m *memdOpMap) FindAndRemoveAllPersistent() []*memdQRequest {
	var reqs []*memdQRequest
	for _, req := range m.requests {
		if req.Persistent {
			reqs = append(reqs, req)
			delete(m.requests, req.Opaque)
		}
	}

	return reqs
}

// Find - Lookup a request using its opaque, note that this function by return a <nil> pointer.
func (m *memdOpMap) Find(opaque uint32) *memdQRequest {
	return m.requests[opaque]
}

// FindAndMaybeRemove - Lookup a request using its opaque and then remove it from the map if it's not persistent or the
// 'force' argument is true.
func (m *memdOpMap) FindAndMaybeRemove(opaque uint32, force bool) *memdQRequest {
	req, ok := m.requests[opaque]
	if !ok {
		return nil
	}

	if force || !req.Persistent {
		delete(m.requests, opaque)
	}

	return req
}

func (m *memdOpMap) Size() int {
	return len(m.requests)
}

// Drain - Remove all the requests from the map whilst running the provided callback for each request.
func (m *memdOpMap) Drain(callback func(req *memdQRequest)) {
	for _, req := range m.requests {
		callback(req)
	}

	m.requests = make(map[uint32]*memdQRequest)
}

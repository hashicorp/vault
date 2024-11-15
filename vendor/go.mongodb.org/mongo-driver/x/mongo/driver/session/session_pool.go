// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Node represents a server session in a linked list
type Node struct {
	*Server
	next *Node
	prev *Node
}

// topologyDescription is used to track a subset of the fields present in a description.Topology instance that are
// relevant for determining session expiration.
type topologyDescription struct {
	kind           description.TopologyKind
	timeoutMinutes *int64
}

// Pool is a pool of server sessions that can be reused.
type Pool struct {
	// number of sessions checked out of pool (accessed atomically)
	checkedOut int64

	descChan       <-chan description.Topology
	head           *Node
	tail           *Node
	latestTopology topologyDescription
	mutex          sync.Mutex // mutex to protect list and sessionTimeout
}

func (p *Pool) createServerSession() (*Server, error) {
	s, err := newServerSession()
	if err != nil {
		return nil, err
	}

	atomic.AddInt64(&p.checkedOut, 1)
	return s, nil
}

// NewPool creates a new server session pool
func NewPool(descChan <-chan description.Topology) *Pool {
	p := &Pool{
		descChan: descChan,
	}

	return p
}

// assumes caller has mutex to protect the pool
func (p *Pool) updateTimeout() {
	select {
	case newDesc := <-p.descChan:
		p.latestTopology = topologyDescription{
			kind:           newDesc.Kind,
			timeoutMinutes: newDesc.SessionTimeoutMinutesPtr,
		}
	default:
		// no new description waiting
	}
}

// GetSession retrieves an unexpired session from the pool.
func (p *Pool) GetSession() (*Server, error) {
	p.mutex.Lock() // prevent changing the linked list while seeing if sessions have expired
	defer p.mutex.Unlock()

	// empty pool
	if p.head == nil && p.tail == nil {
		return p.createServerSession()
	}

	p.updateTimeout()
	for p.head != nil {
		// pull session from head of queue and return if it is valid for at least 1 more minute
		if p.head.expired(p.latestTopology) {
			p.head = p.head.next
			continue
		}

		// found unexpired session
		session := p.head.Server
		if p.head.next != nil {
			p.head.next.prev = nil
		}
		if p.tail == p.head {
			p.tail = nil
			p.head = nil
		} else {
			p.head = p.head.next
		}

		atomic.AddInt64(&p.checkedOut, 1)
		return session, nil
	}

	// no valid session found
	p.tail = nil // empty list
	return p.createServerSession()
}

// ReturnSession returns a session to the pool if it has not expired.
func (p *Pool) ReturnSession(ss *Server) {
	if ss == nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	atomic.AddInt64(&p.checkedOut, -1)
	p.updateTimeout()
	// check sessions at end of queue for expired
	// stop checking after hitting the first valid session
	for p.tail != nil && p.tail.expired(p.latestTopology) {
		if p.tail.prev != nil {
			p.tail.prev.next = nil
		}
		p.tail = p.tail.prev
	}

	// session expired
	if ss.expired(p.latestTopology) {
		return
	}

	// session is dirty
	if ss.Dirty {
		return
	}

	newNode := &Node{
		Server: ss,
		next:   nil,
		prev:   nil,
	}

	// empty list
	if p.tail == nil {
		p.head = newNode
		p.tail = newNode
		return
	}

	// at least 1 valid session in list
	newNode.next = p.head
	p.head.prev = newNode
	p.head = newNode
}

// IDSlice returns a slice of session IDs for each session in the pool
func (p *Pool) IDSlice() []bsoncore.Document {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var ids []bsoncore.Document
	for node := p.head; node != nil; node = node.next {
		ids = append(ids, node.SessionID)
	}

	return ids
}

// String implements the Stringer interface
func (p *Pool) String() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	s := ""
	for head := p.head; head != nil; head = head.next {
		s += head.SessionID.String() + "\n"
	}

	return s
}

// CheckedOut returns number of sessions checked out from pool.
func (p *Pool) CheckedOut() int64 {
	return atomic.LoadInt64(&p.checkedOut)
}

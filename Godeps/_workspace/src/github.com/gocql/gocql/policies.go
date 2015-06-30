// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//This file will be the future home for more policies
package gocql

import (
	"log"
	"sync"
	"sync/atomic"
)

//RetryableQuery is an interface that represents a query or batch statement that
//exposes the correct functions for the retry policy logic to evaluate correctly.
type RetryableQuery interface {
	Attempts() int
	GetConsistency() Consistency
}

// RetryPolicy interface is used by gocql to determine if a query can be attempted
// again after a retryable error has been received. The interface allows gocql
// users to implement their own logic to determine if a query can be attempted
// again.
//
// See SimpleRetryPolicy as an example of implementing and using a RetryPolicy
// interface.
type RetryPolicy interface {
	Attempt(RetryableQuery) bool
}

// SimpleRetryPolicy has simple logic for attempting a query a fixed number of times.
//
// See below for examples of usage:
//
//     //Assign to the cluster
//     cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
//
//     //Assign to a query
//     query.RetryPolicy(&gocql.SimpleRetryPolicy{NumRetries: 1})
//
type SimpleRetryPolicy struct {
	NumRetries int //Number of times to retry a query
}

// Attempt tells gocql to attempt the query again based on query.Attempts being less
// than the NumRetries defined in the policy.
func (s *SimpleRetryPolicy) Attempt(q RetryableQuery) bool {
	return q.Attempts() <= s.NumRetries
}

//HostSelectionPolicy is an interface for selecting
//the most appropriate host to execute a given query.
type HostSelectionPolicy interface {
	SetHosts
	SetPartitioner
	//Pick returns an iteration function over selected hosts
	Pick(*Query) NextHost
}

//NextHost is an iteration function over picked hosts
type NextHost func() *HostInfo

//NewRoundRobinHostPolicy is a round-robin load balancing policy
func NewRoundRobinHostPolicy() HostSelectionPolicy {
	return &roundRobinHostPolicy{hosts: []HostInfo{}}
}

type roundRobinHostPolicy struct {
	hosts []HostInfo
	pos   uint32
	mu    sync.RWMutex
}

func (r *roundRobinHostPolicy) SetHosts(hosts []HostInfo) {
	r.mu.Lock()
	r.hosts = hosts
	r.mu.Unlock()
}

func (r *roundRobinHostPolicy) SetPartitioner(partitioner string) {
	// noop
}

func (r *roundRobinHostPolicy) Pick(qry *Query) NextHost {
	// i is used to limit the number of attempts to find a host
	// to the number of hosts known to this policy
	var i uint32 = 0
	return func() *HostInfo {
		r.mu.RLock()
		if len(r.hosts) == 0 {
			r.mu.RUnlock()
			return nil
		}

		var host *HostInfo
		// always increment pos to evenly distribute traffic in case of
		// failures
		pos := atomic.AddUint32(&r.pos, 1)
		if int(i) < len(r.hosts) {
			host = &r.hosts[(pos)%uint32(len(r.hosts))]
			i++
		}
		r.mu.RUnlock()
		return host
	}
}

//NewTokenAwareHostPolicy is a token aware host selection policy
func NewTokenAwareHostPolicy(fallback HostSelectionPolicy) HostSelectionPolicy {
	return &tokenAwareHostPolicy{fallback: fallback, hosts: []HostInfo{}}
}

type tokenAwareHostPolicy struct {
	mu          sync.RWMutex
	hosts       []HostInfo
	partitioner string
	tokenRing   *tokenRing
	fallback    HostSelectionPolicy
}

func (t *tokenAwareHostPolicy) SetHosts(hosts []HostInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// always update the fallback
	t.fallback.SetHosts(hosts)
	t.hosts = hosts

	t.resetTokenRing()
}

func (t *tokenAwareHostPolicy) SetPartitioner(partitioner string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.partitioner != partitioner {
		t.fallback.SetPartitioner(partitioner)
		t.partitioner = partitioner

		t.resetTokenRing()
	}
}

func (t *tokenAwareHostPolicy) resetTokenRing() {
	if t.partitioner == "" {
		// partitioner not yet set
		return
	}

	// create a new token ring
	tokenRing, err := newTokenRing(t.partitioner, t.hosts)
	if err != nil {
		log.Printf("Unable to update the token ring due to error: %s", err)
		return
	}

	// replace the token ring
	t.tokenRing = tokenRing
}

func (t *tokenAwareHostPolicy) Pick(qry *Query) NextHost {
	if qry == nil {
		return t.fallback.Pick(qry)
	}

	routingKey, err := qry.GetRoutingKey()
	if err != nil {
		return t.fallback.Pick(qry)
	}
	if routingKey == nil {
		return t.fallback.Pick(qry)
	}

	var host *HostInfo

	t.mu.RLock()
	// TODO retrieve a list of hosts based on the replication strategy
	host = t.tokenRing.GetHostForPartitionKey(routingKey)
	t.mu.RUnlock()

	if host == nil {
		return t.fallback.Pick(qry)
	}

	// scope these variables for the same lifetime as the iterator function
	var (
		hostReturned bool
		fallbackIter NextHost
	)
	return func() *HostInfo {
		if !hostReturned {
			hostReturned = true
			return host
		}

		// fallback
		if fallbackIter == nil {
			fallbackIter = t.fallback.Pick(qry)
		}

		fallbackHost := fallbackIter()

		// filter the token aware selected hosts from the fallback hosts
		if fallbackHost == host {
			fallbackHost = fallbackIter()
		}

		return fallbackHost
	}
}

//ConnSelectionPolicy is an interface for selecting an
//appropriate connection for executing a query
type ConnSelectionPolicy interface {
	SetConns(conns []*Conn)
	Pick(*Query) *Conn
}

type roundRobinConnPolicy struct {
	conns []*Conn
	pos   uint32
	mu    sync.RWMutex
}

func NewRoundRobinConnPolicy() ConnSelectionPolicy {
	return &roundRobinConnPolicy{}
}

func (r *roundRobinConnPolicy) SetConns(conns []*Conn) {
	r.mu.Lock()
	r.conns = conns
	r.mu.Unlock()
}

func (r *roundRobinConnPolicy) Pick(qry *Query) *Conn {
	pos := atomic.AddUint32(&r.pos, 1)
	var conn *Conn
	r.mu.RLock()
	if len(r.conns) > 0 {
		conn = r.conns[pos%uint32(len(r.conns))]
	}
	r.mu.RUnlock()
	return conn
}

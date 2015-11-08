// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//This file will be the future home for more policies
package gocql

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/hailocab/go-hostpool"
)

// RetryableQuery is an interface that represents a query or batch statement that
// exposes the correct functions for the retry policy logic to evaluate correctly.
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

// HostSelectionPolicy is an interface for selecting
// the most appropriate host to execute a given query.
type HostSelectionPolicy interface {
	SetHosts
	SetPartitioner
	//Pick returns an iteration function over selected hosts
	Pick(*Query) NextHost
}

// SelectedHost is an interface returned when picking a host from a host
// selection policy.
type SelectedHost interface {
	Info() *HostInfo
	Mark(error)
}

// NextHost is an iteration function over picked hosts
type NextHost func() SelectedHost

// RoundRobinHostPolicy is a round-robin load balancing policy, where each host
// is tried sequentially for each query.
func RoundRobinHostPolicy() HostSelectionPolicy {
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
	var i uint32
	return func() SelectedHost {
		r.mu.RLock()
		defer r.mu.RUnlock()
		if len(r.hosts) == 0 {
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
		return selectedRoundRobinHost{host}
	}
}

// selectedRoundRobinHost is a host returned by the roundRobinHostPolicy and
// implements the SelectedHost interface
type selectedRoundRobinHost struct {
	info *HostInfo
}

func (host selectedRoundRobinHost) Info() *HostInfo {
	return host.info
}

func (host selectedRoundRobinHost) Mark(err error) {
	// noop
}

// TokenAwareHostPolicy is a token aware host selection policy, where hosts are
// selected based on the partition key, so queries are sent to the host which
// owns the partition. Fallback is used when routing information is not available.
func TokenAwareHostPolicy(fallback HostSelectionPolicy) HostSelectionPolicy {
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
	} else if qry.binding != nil && len(qry.values) == 0 {
		// If this query was created using session.Bind we wont have the query
		// values yet, so we have to pass down to the next policy.
		// TODO: Remove this and handle this case
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
	return func() SelectedHost {
		if !hostReturned {
			hostReturned = true
			return selectedTokenAwareHost{host}
		}

		// fallback
		if fallbackIter == nil {
			fallbackIter = t.fallback.Pick(qry)
		}

		fallbackHost := fallbackIter()

		// filter the token aware selected hosts from the fallback hosts
		if fallbackHost.Info() == host {
			fallbackHost = fallbackIter()
		}

		return fallbackHost
	}
}

// selectedTokenAwareHost is a host returned by the tokenAwareHostPolicy and
// implements the SelectedHost interface
type selectedTokenAwareHost struct {
	info *HostInfo
}

func (host selectedTokenAwareHost) Info() *HostInfo {
	return host.info
}

func (host selectedTokenAwareHost) Mark(err error) {
	// noop
}

// HostPoolHostPolicy is a host policy which uses the bitly/go-hostpool library
// to distribute queries between hosts and prevent sending queries to
// unresponsive hosts. When creating the host pool that is passed to the policy
// use an empty slice of hosts as the hostpool will be populated later by gocql.
// See below for examples of usage:
//
//     // Create host selection policy using a simple host pool
//     cluster.PoolConfig.HostSelectionPolicy = HostPoolHostPolicy(hostpool.New(nil))
//
//     // Create host selection policy using an epsilon greddy pool
//     cluster.PoolConfig.HostSelectionPolicy = HostPoolHostPolicy(
//         hostpool.NewEpsilonGreedy(nil, 0, &hostpool.LinearEpsilonValueCalculator{}),
//     )
//
func HostPoolHostPolicy(hp hostpool.HostPool) HostSelectionPolicy {
	return &hostPoolHostPolicy{hostMap: map[string]HostInfo{}, hp: hp}
}

type hostPoolHostPolicy struct {
	hp      hostpool.HostPool
	hostMap map[string]HostInfo
	mu      sync.RWMutex
}

func (r *hostPoolHostPolicy) SetHosts(hosts []HostInfo) {
	peers := make([]string, len(hosts))
	hostMap := make(map[string]HostInfo, len(hosts))

	for i, host := range hosts {
		peers[i] = host.Peer
		hostMap[host.Peer] = host
	}

	r.mu.Lock()
	r.hp.SetHosts(peers)
	r.hostMap = hostMap
	r.mu.Unlock()
}

func (r *hostPoolHostPolicy) SetPartitioner(partitioner string) {
	// noop
}

func (r *hostPoolHostPolicy) Pick(qry *Query) NextHost {
	return func() SelectedHost {
		r.mu.RLock()
		defer r.mu.RUnlock()

		if len(r.hostMap) == 0 {
			return nil
		}

		hostR := r.hp.Get()
		host, ok := r.hostMap[hostR.Host()]
		if !ok {
			return nil
		}

		return selectedHostPoolHost{&host, hostR}
	}
}

// selectedHostPoolHost is a host returned by the hostPoolHostPolicy and
// implements the SelectedHost interface
type selectedHostPoolHost struct {
	info  *HostInfo
	hostR hostpool.HostPoolResponse
}

func (host selectedHostPoolHost) Info() *HostInfo {
	return host.info
}

func (host selectedHostPoolHost) Mark(err error) {
	host.hostR.Mark(err)
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

func RoundRobinConnPolicy() func() ConnSelectionPolicy {
	return func() ConnSelectionPolicy {
		return &roundRobinConnPolicy{}
	}
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

package raft

import (
	"sync/atomic"
	"time"
)

// Observation is sent along the given channel to observers when an event occurs.
type Observation struct {
	// Raft holds the Raft instance generating the observation.
	Raft *Raft
	// Data holds observation-specific data. Possible types are
	// *RequestVoteRequest
	// RaftState
	// PeerObservation
	// LeaderObservation
	Data interface{}
}

// LeaderObservation is used for the data when leadership changes.
type LeaderObservation struct {
	Leader ServerAddress
}

// PeerObservation is sent to observers when peers change.
type PeerObservation struct {
	Removed bool
	Peer    Server
}

// FailedHeartbeatObservation is sent when a node fails to heartbeat with the leader
type FailedHeartbeatObservation struct {
	PeerID      ServerID
	LastContact time.Time
}

// ResumedHeartbeatObservation is sent when a node resumes to heartbeat with the leader following failures
type ResumedHeartbeatObservation struct {
	PeerID ServerID
}

// nextObserverId is used to provide a unique ID for each observer to aid in
// deregistration.
var nextObserverID uint64

// FilterFn is a function that can be registered in order to filter observations.
// The function reports whether the observation should be included - if
// it returns false, the observation will be filtered out.
type FilterFn func(o *Observation) bool

// Observer describes what to do with a given observation.
type Observer struct {
	// numObserved and numDropped are performance counters for this observer.
	// 64 bit types must be 64 bit aligned to use with atomic operations on
	// 32 bit platforms, so keep them at the top of the struct.
	numObserved uint64
	numDropped  uint64

	// channel receives observations.
	channel chan Observation

	// blocking, if true, will cause Raft to block when sending an observation
	// to this observer. This should generally be set to false.
	blocking bool

	// filter will be called to determine if an observation should be sent to
	// the channel.
	filter FilterFn

	// id is the ID of this observer in the Raft map.
	id uint64
}

// NewObserver creates a new observer that can be registered
// to make observations on a Raft instance. Observations
// will be sent on the given channel if they satisfy the
// given filter.
//
// If blocking is true, the observer will block when it can't
// send on the channel, otherwise it may discard events.
func NewObserver(channel chan Observation, blocking bool, filter FilterFn) *Observer {
	return &Observer{
		channel:  channel,
		blocking: blocking,
		filter:   filter,
		id:       atomic.AddUint64(&nextObserverID, 1),
	}
}

// GetNumObserved returns the number of observations.
func (or *Observer) GetNumObserved() uint64 {
	return atomic.LoadUint64(&or.numObserved)
}

// GetNumDropped returns the number of dropped observations due to blocking.
func (or *Observer) GetNumDropped() uint64 {
	return atomic.LoadUint64(&or.numDropped)
}

// RegisterObserver registers a new observer.
func (r *Raft) RegisterObserver(or *Observer) {
	r.observersLock.Lock()
	defer r.observersLock.Unlock()
	r.observers[or.id] = or
}

// DeregisterObserver deregisters an observer.
func (r *Raft) DeregisterObserver(or *Observer) {
	r.observersLock.Lock()
	defer r.observersLock.Unlock()
	delete(r.observers, or.id)
}

// observe sends an observation to every observer.
func (r *Raft) observe(o interface{}) {
	// In general observers should not block. But in any case this isn't
	// disastrous as we only hold a read lock, which merely prevents
	// registration / deregistration of observers.
	r.observersLock.RLock()
	defer r.observersLock.RUnlock()
	for _, or := range r.observers {
		// It's wasteful to do this in the loop, but for the common case
		// where there are no observers we won't create any objects.
		ob := Observation{Raft: r, Data: o}
		if or.filter != nil && !or.filter(&ob) {
			continue
		}
		if or.channel == nil {
			continue
		}
		if or.blocking {
			or.channel <- ob
			atomic.AddUint64(&or.numObserved, 1)
		} else {
			select {
			case or.channel <- ob:
				atomic.AddUint64(&or.numObserved, 1)
			default:
				atomic.AddUint64(&or.numDropped, 1)
			}
		}
	}
}

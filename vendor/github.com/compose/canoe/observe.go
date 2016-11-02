package canoe

import (
	"sync/atomic"
)

// Observation is sent out to each observer.
// An obeservation can have many different types.
// It is currently used to detect the successful addition of a node to
// a cluster during the cluster join or bootstrap phase
type Observation interface{}

// FilterFn is a function used to filter what events an Observer gets piped
type FilterFn func(o Observation) bool

var nextObserverID uint64

// Observer is a struct responsible for monitoring raft's internal operations if one
// wants to perform unpredicted operations
type Observer struct {
	channel chan Observation
	filter  FilterFn
	id      uint64
}

// NewObserver gets an observer. Note, if you aren't actively consuming the observer,
// the observations will get lost
func NewObserver(channel chan Observation, filter FilterFn) *Observer {
	return &Observer{
		channel: channel,
		filter:  filter,
		id:      atomic.AddUint64(&nextObserverID, 1),
	}
}

func (rn *Node) observe(data Observation) {
	rn.observersLock.RLock()
	defer rn.observersLock.RUnlock()
	for _, observer := range rn.observers {
		if observer.filter != nil && !observer.filter(interface{}(data).(Observation)) {
			continue
		}
		if observer.channel == nil {
			continue
		}

		// make sure we don't block if consumer isn't consuming fast enough
		select {
		case observer.channel <- data:
			continue
		default:
			continue
		}
	}
}

// RegisterObserver registers and begins to send observations down an Observer
func (rn *Node) RegisterObserver(o *Observer) {
	rn.observersLock.Lock()
	defer rn.observersLock.Unlock()
	rn.observers[o.id] = o
}

// UnregisterObserver is called when one no longer needs to look for a particular raft event occuring
func (rn *Node) UnregisterObserver(o *Observer) {
	rn.observersLock.Lock()
	defer rn.observersLock.Unlock()
	delete(rn.observers, o.id)
}

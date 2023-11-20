// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventbus

import (
	"slices"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
)

// Filters keeps track of all the event patterns that each node is interested in.
type Filters struct {
	lock     sync.RWMutex
	parallel bool
	self     nodeID
	filters  map[nodeID]*NodeFilter
}

// nodeID is used to syntactically indicate that the string is a node's name identifier.
type nodeID string

// pattern is used to represent one or more combinations of patterns
type pattern struct {
	eventTypePattern  string
	namespacePatterns []string
}

// NodeFilter keeps track of all patterns that a particular node is interested in.
type NodeFilter struct {
	patterns []pattern
}

func (nf *NodeFilter) match(nsPath string, eventType logical.EventType) bool {
	if nf == nil {
		return false
	}
	for _, p := range nf.patterns {
		if glob.Glob(p.eventTypePattern, string(eventType)) {
			for _, nsp := range p.namespacePatterns {
				if glob.Glob(nsp, nsPath) {
					return true
				}
			}
		}
	}
	return false
}

// NewFilters creates an empty set of filters to keep track of each node's pattern interests.
func NewFilters(self string) *Filters {
	return &Filters{
		self:    nodeID(self),
		filters: map[nodeID]*NodeFilter{},
	}
}

// addPattern adds a pattern to a node's list.
func (f *Filters) addPattern(node nodeID, namespacePatterns []string, eventTypePattern string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if _, ok := f.filters[node]; !ok {
		f.filters[node] = &NodeFilter{}
	}
	f.filters[node].patterns = append(f.filters[node].patterns, pattern{eventTypePattern: eventTypePattern, namespacePatterns: namespacePatterns})
}

// removePattern removes a pattern from a node's list.
func (f *Filters) removePattern(node nodeID, namespacePatterns []string, eventTypePattern string) {
	check := pattern{eventTypePattern: eventTypePattern, namespacePatterns: namespacePatterns}
	f.lock.Lock()
	defer f.lock.Unlock()
	filters, ok := f.filters[node]
	if !ok {
		return
	}
	filters.patterns = slices.DeleteFunc(filters.patterns, func(m pattern) bool {
		return m.eventTypePattern == check.eventTypePattern &&
			slices.Equal(m.namespacePatterns, check.namespacePatterns)
	})
}

// anyMatch returns true if any node's pattern list matches the arguments.
func (f *Filters) anyMatch(nsPath string, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if f.parallel {
		wg := sync.WaitGroup{}
		anyMatched := atomic.Bool{}
		for _, nf := range f.filters {
			wg.Add(1)
			go func(nf *NodeFilter) {
				if nf.match(nsPath, eventType) {
					anyMatched.Store(true)
				}
				wg.Done()
			}(nf)
		}
		wg.Wait()
		return anyMatched.Load()
	} else {
		for _, nf := range f.filters {
			if nf.match(nsPath, eventType) {
				return true
			}
		}
		return false
	}
}

// nodeMatch returns true if the given node's pattern list matches the arguments.
func (f *Filters) nodeMatch(node nodeID, nsPath string, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.filters[node].match(nsPath, eventType)
}

// localMatch returns true if the local node's pattern list matches the arguments.
func (f *Filters) localMatch(nsPath string, eventType logical.EventType) bool {
	return f.nodeMatch(f.self, nsPath, eventType)
}

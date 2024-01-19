// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"slices"
	"sort"
	"sync"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
)

// Filters keeps track of all the event patterns that each node is interested in.
type Filters struct {
	lock    sync.RWMutex
	self    nodeID
	filters map[nodeID]*NodeFilter
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

func (nf *NodeFilter) match(ns *namespace.Namespace, eventType logical.EventType) bool {
	if nf == nil {
		return false
	}
	for _, p := range nf.patterns {
		if glob.Glob(p.eventTypePattern, string(eventType)) {
			for _, nsp := range p.namespacePatterns {
				if glob.Glob(nsp, ns.Path) {
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
	nsPatterns := slices.Clone(namespacePatterns)
	sort.Strings(nsPatterns)
	f.filters[node].patterns = append(f.filters[node].patterns, pattern{eventTypePattern: eventTypePattern, namespacePatterns: nsPatterns})
}

func (f *Filters) addNsPattern(node nodeID, ns *namespace.Namespace, eventTypePattern string) {
	f.addPattern(node, []string{ns.Path}, eventTypePattern)
}

// removePattern removes a pattern from a node's list.
func (f *Filters) removePattern(node nodeID, namespacePatterns []string, eventTypePattern string) {
	nsPatterns := slices.Clone(namespacePatterns)
	sort.Strings(nsPatterns)
	check := pattern{eventTypePattern: eventTypePattern, namespacePatterns: nsPatterns}
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

func (f *Filters) removeNsPattern(node nodeID, ns *namespace.Namespace, eventTypePattern string) {
	f.removePattern(node, []string{ns.Path}, eventTypePattern)
}

// anyMatch returns true if any node's pattern list matches the arguments.
func (f *Filters) anyMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	for _, nf := range f.filters {
		if nf.match(ns, eventType) {
			return true
		}
	}
	return false
}

// nodeMatch returns true if the given node's pattern list matches the arguments.
func (f *Filters) nodeMatch(node nodeID, ns *namespace.Namespace, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.filters[node].match(ns, eventType)
}

// localMatch returns true if the local node's pattern list matches the arguments.
func (f *Filters) localMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	return f.nodeMatch(f.self, ns, eventType)
}

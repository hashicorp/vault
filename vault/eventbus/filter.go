// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
	"k8s.io/apimachinery/pkg/util/sets"
)

const globalCluster = ""

// Filters keeps track of all the event patterns that each cluster node is interested in.
type Filters struct {
	lock    sync.RWMutex
	self    clusterNodeID
	filters map[clusterNodeID]*ClusterNodeFilter

	// notifyChanges is used to notify about changes to filters. The condition variables are tied to single lock above.
	notifyChanges map[clusterNodeID]*sync.Cond
}

// clusterNodeID is used to syntactically indicate that the string is a cluster nodes's identifier.
type clusterNodeID string

// pattern is used to represent one or more combinations of patterns
type pattern struct {
	namespacePatterns string // space-separated (spaces are not allowed in namespaces, and slices are not comparable)
	eventTypePattern  string
}

func (p pattern) String() string {
	return fmt.Sprintf("{ns=%s,ev=%s}", p.namespacePatterns, p.eventTypePattern)
}

func (p pattern) isEmpty() bool {
	return p.namespacePatterns == "" && p.eventTypePattern == ""
}

// ClusterNodeFilter keeps track of all patterns that a particular cluster node is interested in.
type ClusterNodeFilter struct {
	patterns sets.Set[pattern]
}

// match checks if the given ns and eventType matches any pattern in the cluster node's filter.
// Must be called while holding a (read) lock for the filter.
func (nf *ClusterNodeFilter) match(ns *namespace.Namespace, eventType logical.EventType) bool {
	if nf == nil {
		return false
	}
	for p := range nf.patterns {
		if glob.Glob(p.eventTypePattern, string(eventType)) {
			for _, nsp := range strings.Split(p.namespacePatterns, " ") {
				if glob.Glob(nsp, ns.Path) {
					return true
				}
			}
		}
	}
	return false
}

// NewFilters creates an empty set of filters to keep track of each cluster node's pattern interests.
func NewFilters(self string) *Filters {
	f := &Filters{
		self:          clusterNodeID(self),
		filters:       map[clusterNodeID]*ClusterNodeFilter{},
		notifyChanges: map[clusterNodeID]*sync.Cond{},
	}
	f.notifyChanges[clusterNodeID(self)] = sync.NewCond(&f.lock)
	f.notifyChanges[globalCluster] = sync.NewCond(&f.lock)
	return f
}

func (f *Filters) String() string {
	x := "Filters {\n"
	for k, v := range f.filters {
		x += fmt.Sprintf("  %s: {%s}\n", k, v)
	}
	return x
}

func (nf *ClusterNodeFilter) String() string {
	var x []string
	l := nf.patterns.UnsortedList()
	for _, v := range l {
		x = append(x, v.String())
	}
	return strings.Join(x, ",")
}

func (f *Filters) addGlobalPattern(namespacePatterns []string, eventTypePattern string) {
	f.addPattern(globalCluster, namespacePatterns, eventTypePattern)
}

func (f *Filters) removeGlobalPattern(namespacePatterns []string, eventTypePattern string) {
	f.removePattern(globalCluster, namespacePatterns, eventTypePattern)
}

func (f *Filters) clearGlobalPatterns() {
	defer f.notify(globalCluster)
	f.lock.Lock()
	defer f.lock.Unlock()
	delete(f.filters, globalCluster)
}

func (f *Filters) getOrCreateNotify(c clusterNodeID) *sync.Cond {
	// fast check when we don't need to create the Cond
	f.lock.RLock()
	n, ok := f.notifyChanges[c]
	f.lock.RUnlock()
	if ok {
		return n
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	// check again to avoid race condition
	n, ok = f.notifyChanges[c]
	if ok {
		return n
	}
	n = sync.NewCond(&f.lock)
	f.notifyChanges[c] = n
	return n
}

func (f *Filters) notify(c clusterNodeID) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if notifier, ok := f.notifyChanges[c]; ok {
		notifier.Broadcast()
	}
}

func (f *Filters) clearClusterNodePatterns(c clusterNodeID) {
	defer f.notify(c)
	f.lock.Lock()
	defer f.lock.Unlock()
	delete(f.filters, c)
}

// copyPatternWithLock gets a copy of a cluster node's filters
func (f *Filters) copyPatternWithLock(c clusterNodeID) *ClusterNodeFilter {
	filters := &ClusterNodeFilter{}
	if got, ok := f.filters[c]; ok {
		filters.patterns = got.patterns.Clone()
	} else {
		filters.patterns = sets.New[pattern]()
	}
	return filters
}

// applyChanges applies the changes in the given list, atomically.
func (f *Filters) applyChanges(c clusterNodeID, changes []FilterChange) {
	defer f.notify(c)
	f.lock.Lock()
	defer f.lock.Unlock()
	var newPatterns sets.Set[pattern]
	if existing, ok := f.filters[c]; ok {
		newPatterns = existing.patterns
	} else {
		newPatterns = sets.New[pattern]()
	}
	for _, change := range changes {
		applyChange(newPatterns, &change)
	}
	f.filters[c] = &ClusterNodeFilter{patterns: newPatterns}
}

// applyChange applies a single filter change to the given set.
func applyChange(s sets.Set[pattern], change *FilterChange) {
	switch change.Operation {
	case FilterChangeAdd:
		nsPatterns := slices.Clone(change.NamespacePatterns)
		sort.Strings(nsPatterns)
		p := pattern{eventTypePattern: change.EventTypePattern, namespacePatterns: cleanJoinNamespaces(nsPatterns)}
		s.Insert(p)
	case FilterChangeRemove:
		nsPatterns := slices.Clone(change.NamespacePatterns)
		sort.Strings(nsPatterns)
		check := pattern{eventTypePattern: change.EventTypePattern, namespacePatterns: cleanJoinNamespaces(nsPatterns)}
		s.Delete(check)
	case FilterChangeClear:
		s.Clear()
	}
}

func cleanJoinNamespaces(nsPatterns []string) string {
	trimmed := make([]string, len(nsPatterns))
	for i := 0; i < len(nsPatterns); i++ {
		trimmed[i] = strings.TrimSpace(nsPatterns[i])
	}
	// sort and uniq
	trimmed = sets.NewString(trimmed...).List()
	return strings.Join(trimmed, " ")
}

// addPattern adds a pattern to a cluster node's list.
func (f *Filters) addPattern(c clusterNodeID, namespacePatterns []string, eventTypePattern string) {
	defer f.notify(c)
	f.lock.Lock()
	defer f.lock.Unlock()
	if _, ok := f.filters[c]; !ok {
		f.filters[c] = &ClusterNodeFilter{
			patterns: sets.New[pattern](),
		}
	}
	nsPatterns := slices.Clone(namespacePatterns)
	sort.Strings(nsPatterns)
	p := pattern{eventTypePattern: eventTypePattern, namespacePatterns: cleanJoinNamespaces(namespacePatterns)}
	f.filters[c].patterns.Insert(p)
}

// removePattern removes a pattern from a cluster node's list.
func (f *Filters) removePattern(c clusterNodeID, namespacePatterns []string, eventTypePattern string) {
	defer f.notify(c)
	nsPatterns := slices.Clone(namespacePatterns)
	sort.Strings(nsPatterns)
	check := pattern{eventTypePattern: eventTypePattern, namespacePatterns: cleanJoinNamespaces(nsPatterns)}
	f.lock.Lock()
	defer f.lock.Unlock()
	filters, ok := f.filters[c]
	if !ok {
		return
	}
	filters.patterns.Delete(check)
}

// anyMatch returns true if any cluster node's pattern list matches the arguments.
func (f *Filters) anyMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	for _, cf := range f.filters {
		if cf.match(ns, eventType) {
			return true
		}
	}
	return false
}

// globalMatch returns true if the global cluster's pattern list matches the arguments.
func (f *Filters) globalMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	return f.clusterNodeMatch(globalCluster, ns, eventType)
}

// clusterNodeMatch returns true if the given cluster node's pattern list matches the arguments.
func (f *Filters) clusterNodeMatch(c clusterNodeID, ns *namespace.Namespace, eventType logical.EventType) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.filters[c].match(ns, eventType)
}

// localMatch returns true if the local cluster node's pattern list matches the arguments.
func (f *Filters) localMatch(ns *namespace.Namespace, eventType logical.EventType) bool {
	return f.clusterNodeMatch(f.self, ns, eventType)
}

// watch creates a notification channel that receives changes for the given cluster node.
func (f *Filters) watch(ctx context.Context, clusterNode clusterNodeID) (<-chan []FilterChange, context.CancelFunc, error) {
	notify := f.getOrCreateNotify(clusterNode)
	ctx, cancelFunc := context.WithCancel(ctx)
	doneCh := ctx.Done()
	ch := make(chan []FilterChange)

	// ensure that the sleeping goroutine wakes up if the channel is closed
	go func() {
		select {
		case <-doneCh:
			notify.Broadcast()
		}
	}()

	sendToNotify := make(chan *ClusterNodeFilter)

	// goroutine for polling the condition variable.
	// it's necessary to hold the lock the entire time to ensure there are no race conditions.
	go func() {
		// use a WG to ensure we don't try to send to a closed channel
		senders := sync.WaitGroup{}
		defer func() {
			senders.Wait()
			close(sendToNotify)
		}()
		f.lock.Lock()
		defer f.lock.Unlock()
		for {
			select {
			case <-doneCh:
				return
			default:
			}
			next := f.copyPatternWithLock(clusterNode)
			senders.Add(1)
			// don't block to send since we hold the lock
			go func() {
				sendToNotify <- next
				senders.Done()
			}()
			notify.Wait()
		}
	}()

	// calculate changes and forward to notification channel
	go func() {
		defer close(ch)
		var current *ClusterNodeFilter
		for {
			next, ok := <-sendToNotify
			if !ok {
				return
			}
			changes := calculateChanges(current, next)
			current = next
			// check if the context is finished before sending
			select {
			case <-doneCh:
				return
			default:
				ch <- changes
			}
		}
	}()

	return ch, cancelFunc, nil
}

// FilterChange represents a change to a cluster node's filters.
type FilterChange struct {
	Operation         int
	NamespacePatterns []string
	EventTypePattern  string
}

const (
	FilterChangeAdd    = 0
	FilterChangeRemove = 1
	FilterChangeClear  = 2
)

// calculateChanges calculates a set of changes necessary to transform from into to.
func calculateChanges(from *ClusterNodeFilter, to *ClusterNodeFilter) []FilterChange {
	var changes []FilterChange
	if to == nil {
		changes = append(changes, FilterChange{
			Operation: FilterChangeClear,
		})
	} else if from == nil {
		changes = append(changes, FilterChange{
			Operation: FilterChangeClear,
		})
		for pattern := range to.patterns {
			if !pattern.isEmpty() {
				changes = append(changes, FilterChange{
					Operation:         FilterChangeAdd,
					NamespacePatterns: strings.Split(pattern.namespacePatterns, " "),
					EventTypePattern:  pattern.eventTypePattern,
				})
			}
		}
	} else {
		additions := to.patterns.Difference(from.patterns)
		subtractions := from.patterns.Difference(to.patterns)
		for add := range additions {
			if !add.isEmpty() {
				changes = append(changes, FilterChange{
					Operation:         FilterChangeAdd,
					NamespacePatterns: strings.Split(add.namespacePatterns, " "),
					EventTypePattern:  add.eventTypePattern,
				})
			}
		}
		for sub := range subtractions {
			if !sub.isEmpty() {
				changes = append(changes, FilterChange{
					Operation:         FilterChangeRemove,
					NamespacePatterns: strings.Split(sub.namespacePatterns, " "),
					EventTypePattern:  sub.eventTypePattern,
				})
			}
		}
	}
	return changes
}

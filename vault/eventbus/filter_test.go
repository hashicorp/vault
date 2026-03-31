// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/stretchr/testify/assert"
)

// TestCleanJoinNamespaces tests some cases in cleanJoinNamespaces.
func TestCleanJoinNamespaces(t *testing.T) {
	assert.Equal(t, "", cleanJoinNamespaces([]string{""}))
	assert.Equal(t, " abc", cleanJoinNamespaces([]string{"", "abc"}))
	// just checking that inverting works as expected
	assert.Equal(t, []string{"", "abc"}, strings.Split(" abc", " "))
	assert.Equal(t, "abc", cleanJoinNamespaces([]string{"abc"}))
	assert.Equal(t, "abc def", cleanJoinNamespaces([]string{"def", "abc"}))
}

// TestFilters_AddRemoveMatchLocal checks that basic matching, adding, and removing of patterns all work.
func TestFilters_AddRemoveMatchLocal(t *testing.T) {
	f := NewFilters("self")
	ns := &namespace.Namespace{
		ID:   "ns1",
		Path: "ns1",
	}

	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
	f.addPattern("self", []string{ns.Path}, "abc", "uuid")
	assert.True(t, f.localMatch(ns, "abc"))
	assert.False(t, f.localMatch(ns, "abcd"))
	assert.True(t, f.anyMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abcd"))
	f.removePattern("self", []string{ns.Path}, "abc", "uuid")
	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
}

// TestFilters_Watch checks that adding a watch for a cluster node will send a
// notification when the patterns are modified.
func TestFilters_Watch(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	f := NewFilters("self")
	f.addPattern("self", []string{"ns1"}, "e3", "uuid")
	ch, cancelFunc2 := f.watch(ctx, "self")
	t.Cleanup(cancelFunc2)
	initial := <-ch // we always get one immediately for the current state
	assert.Len(t, initial, 2)
	assert.Equal(t, FilterChangeClear, initial[0].Operation)
	assert.Equal(t, FilterChangeAdd, initial[1].Operation)
	assert.Equal(t, []string{"ns1"}, initial[1].NamespacePatterns)
	assert.Equal(t, "e3", initial[1].EventTypePattern)
	assert.Equal(t, "uuid", initial[1].SubscriberID)

	go func() {
		f.addPattern("self", []string{"ns1"}, "e2", "uuid1")
	}()
	changes := waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeAdd,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e2",
		SubscriberID:      "uuid1",
	}}, changes)
	go func() {
		f.addPattern("self", []string{"ns1"}, "e2", "uuid2")
	}()
	changes = waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeAdd,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e2",
		SubscriberID:      "uuid2",
	}}, changes)
	go func() {
		f.removePattern("self", []string{"ns1"}, "e3", "uuid")
	}()
	changes = waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeRemove,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e3",
		SubscriberID:      "uuid",
	}}, changes)

	// Remove and add one of the e2 patterns to test the copyPatternWithLock
	// logic
	go func() {
		f.removePattern("self", []string{"ns1"}, "e2", "uuid1")
	}()
	changes = waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeRemove,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e2",
		SubscriberID:      "uuid1",
	}}, changes)
	go func() {
		f.addPattern("self", []string{"ns1"}, "e2", "uuid1")
	}()
	changes = waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeAdd,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e2",
		SubscriberID:      "uuid1",
	}}, changes)
}

func waitForChanges(t *testing.T, ch <-chan []FilterChange) []FilterChange {
	t.Helper()
	timeout := time.After(2000 * time.Millisecond)
	var changes []FilterChange
	select {
	case changes = <-ch:
	case <-timeout:
		fmt.Println("Timeout waiting for changes")
	}
	return changes
}

// TestFilters_WatchCancel checks that calling the cancel function will clean up the channel.
func TestFilters_WatchCancel(t *testing.T) {
	f := NewFilters("self")
	f.addPattern("self", []string{"ns1"}, "e3", "uuid")
	ch, cancelFunc := f.watch(context.Background(), "self")
	t.Cleanup(cancelFunc)
	initial := <-ch // we always get one immediately for the current state
	assert.Len(t, initial, 2)
	assert.Equal(t, FilterChangeClear, initial[0].Operation)
	assert.Equal(t, FilterChangeAdd, initial[1].Operation)
	assert.Equal(t, []string{"ns1"}, initial[1].NamespacePatterns)
	assert.Equal(t, "e3", initial[1].EventTypePattern)

	var changes []FilterChange
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		changes = waitForChanges(t, ch)
		wg.Done()
	}()

	cancelFunc()
	wg.Wait()
	assert.Nil(t, changes)
	select {
	case _, ok := <-ch:
		assert.False(t, ok)
	default:
		t.Fatal("Channel should be closed")
	}
}

// TestFilters_AddRemoveClear tests that add/remove/clear works as expected.
func TestFilters_AddRemoveClear(t *testing.T) {
	f := NewFilters("self")
	f.addPattern("somecluster", []string{"ns1"}, "abc", "uuid")
	f.removePattern("somecluster", []string{"ns1"}, "abcd", "uuid")
	assert.Equal(t, "{ns=ns1,ev=abc}: [uuid]", f.filters["somecluster"].String())
	f.removePattern("somecluster", []string{"ns1"}, "abc", "uuid")
	assert.Equal(t, "", f.filters["somecluster"].String())
	f.addPattern("somecluster", []string{"ns1"}, "abc", "uuid")
	f.clearClusterNodePatterns("somecluster")
	assert.NotContains(t, f.filters, "somecluster")

	f.addClusterWidePattern([]string{"ns1"}, "abc")
	f.removeClusterWidePattern([]string{"ns1"}, "abcd")
	assert.Equal(t, "{ns=ns1,ev=abc}: [__cluster_wide__]", f.filters[clusterWide].String())
	f.removeClusterWidePattern([]string{"ns1"}, "abc")
	assert.Equal(t, "", f.filters[clusterWide].String())
	f.addClusterWidePattern([]string{"ns1"}, "abc")
	f.clearClusterWidePatterns()
	assert.NotContains(t, f.filters, clusterWide)
}

// TestFilters_makeClusterWideFilters tests that making cluster-wide filters
// works as expected.
func TestFilters_makeClusterWideFilters(t *testing.T) {
	f := NewFilters("self")
	f.addPattern("node1", []string{"ns1"}, "abc", "uuid")
	f.addPattern("node2", []string{"ns1"}, "abc", "uuid")
	f.makeClusterWideFilters()
	assert.Equal(t, "{ns=ns1,ev=abc}: [__cluster_wide__]", f.filters[clusterWide].String())

	f.addPattern("node3", []string{"ns3"}, "def", "uuid")
	f.makeClusterWideFilters()
	assert.Equal(t, "{ns=ns1,ev=abc}: [__cluster_wide__],{ns=ns3,ev=def}: [__cluster_wide__]", f.filters[clusterWide].String())
}

// TestFilters_duplicate_filters tests the underpinnings of the scenario where
// two subscribers subscribe to the same standby node using the same event
// filters, then one disconnects. The filter should still be present for the
// remaining subscriber (so that the active node knows to keep forwarding
// matching events, for example).
func TestFilters_duplicate_filters(t *testing.T) {
	f := NewFilters("self")
	expectedPattern := pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}

	f.addPattern("node1", []string{expectedPattern.namespacePatterns}, expectedPattern.eventTypePattern, "uuid1")
	f.addPattern("node1", []string{expectedPattern.namespacePatterns}, expectedPattern.eventTypePattern, "uuid2")
	assert.Equal(t, 1, len(f.filters["node1"].patterns))
	assert.Equal(t, []string{"uuid1", "uuid2"}, f.filters["node1"].patterns[expectedPattern])

	f.removePattern("node1", []string{"ns1"}, "abc", "uuid1")
	assert.Equal(t, 1, len(f.filters["node1"].patterns))
	assert.Equal(t, []string{"uuid2"}, f.filters["node1"].patterns[expectedPattern])

	f.removePattern("node1", []string{"ns1"}, "abc", "uuid2")
	assert.Empty(t, f.filters["node1"].patterns)
}

// TestPatternSet_basics tests the basic functionality of the patternSet type.
func TestPatternSet_basics(t *testing.T) {
	ps := newPatternSet()
	ps.Insert(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_def")
	ps.Insert(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_abc")
	ps.Insert(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_ghi")
	// The subscriber id's should be sorted
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_abc", "uuid_def", "uuid_ghi"},
		},
		ps,
	)

	// Duplicate uuids should be ignored
	ps.Insert(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_abc")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_abc", "uuid_def", "uuid_ghi"},
		},
		ps,
	)

	ps.Insert(pattern{eventTypePattern: "def", namespacePatterns: "ns2"}, "uuid")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_abc", "uuid_def", "uuid_ghi"},
			pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid"},
		},
		ps,
	)

	ps.Delete(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_abc")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_def", "uuid_ghi"},
			pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid"},
		},
		ps,
	)

	ps.Delete(pattern{eventTypePattern: "def", namespacePatterns: "ns2"}, "uuid")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_def", "uuid_ghi"},
		},
		ps,
	)

	// subscriber id isn't present, nothing deleted
	ps.Delete(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_abc")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_def", "uuid_ghi"},
		},
		ps,
	)

	// Delete all patterns
	ps.Delete(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_def")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid_ghi"},
		},
		ps,
	)
	ps.Delete(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid_ghi")
	assert.Empty(t, ps)

	// Insert again
	ps.Insert(pattern{eventTypePattern: "ghi", namespacePatterns: "ns1"}, "uuid")
	assert.Len(t, ps, 1)
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "ghi", namespacePatterns: "ns1"}: []string{"uuid"},
		},
		ps,
	)
	ps.Insert(pattern{eventTypePattern: "jkl", namespacePatterns: "ns2"}, "uuid")
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "ghi", namespacePatterns: "ns1"}: []string{"uuid"},
			pattern{eventTypePattern: "jkl", namespacePatterns: "ns2"}: []string{"uuid"},
		},
		ps,
	)
	ps.Delete(pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}, "uuid")
	// should have deleted nothing
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "ghi", namespacePatterns: "ns1"}: []string{"uuid"},
			pattern{eventTypePattern: "jkl", namespacePatterns: "ns2"}: []string{"uuid"},
		},
		ps,
	)
	ps.Delete(pattern{eventTypePattern: "ghi", namespacePatterns: "ns1"}, "uuid")
	assert.Len(t, ps, 1)
	assert.Equal(t,
		patternSet{
			pattern{eventTypePattern: "jkl", namespacePatterns: "ns2"}: []string{"uuid"},
		},
		ps,
	)

	ps.Clear()
	assert.Empty(t, ps)
}

// TestPatternSetDelete does more in-depth testing of the patternSet Delete()
func TestPatternSetDelete(t *testing.T) {
	testCases := map[string]struct {
		ps             patternSet
		pattern        pattern
		subscriptionID string
		expected       patternSet
	}{
		"empty ps": {
			ps:             newPatternSet(),
			pattern:        pattern{eventTypePattern: "abc", namespacePatterns: "ns1"},
			subscriptionID: "uuid",
			expected:       patternSet{},
		},
		"remove one full pattern": {
			ps: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1", "uuid2"},
			},
			pattern:        pattern{eventTypePattern: "abc", namespacePatterns: "ns1"},
			subscriptionID: "uuid1",
			expected: patternSet{
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1", "uuid2"},
			},
		},
		"remove first subscription id": {
			ps: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2", "uuid3"},
			},
			pattern:        pattern{eventTypePattern: "abc", namespacePatterns: "ns1"},
			subscriptionID: "uuid1",
			expected: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid2", "uuid3"},
			},
		},
		"remove last subscription id": {
			ps: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2"},
			},
			pattern:        pattern{eventTypePattern: "abc", namespacePatterns: "ns1"},
			subscriptionID: "uuid2",
			expected: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
			},
		},
		"remove only pattern": {
			ps: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
			},
			pattern:        pattern{eventTypePattern: "abc", namespacePatterns: "ns1"},
			subscriptionID: "uuid1",
			expected:       patternSet{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.ps.Delete(tc.pattern, tc.subscriptionID)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestPatternSetDifference tests the Difference method of the patternSet type.
func TestPatternSetDifference(t *testing.T) {
	testCases := map[string]struct {
		ps1          patternSet
		ps2          patternSet
		expectedDiff patternSet
	}{
		"empty ps1": {
			ps1: newPatternSet(),
			ps2: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1", "uuid2"},
			},
			expectedDiff: newPatternSet(),
		},
		"empty ps2": {
			ps1: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1", "uuid2"},
			},
			ps2: newPatternSet(),
			expectedDiff: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1", "uuid2"},
			},
		},
		"both empty": {
			ps1:          newPatternSet(),
			ps2:          newPatternSet(),
			expectedDiff: newPatternSet(),
		},
		"regular diff": {
			ps1: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			ps2: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "xyz", namespacePatterns: ""}:    []string{"uuid1"},
			},
			expectedDiff: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid2"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
		},
	}

	for _, tc := range testCases {
		diff := tc.ps1.Difference(tc.ps2)
		assert.Equal(t, tc.expectedDiff, diff)
	}
}

// Test_calculateChanges exercises the logic for calculating changes in the form
// of []FilterChange between two patternSets
func Test_calculateChanges(t *testing.T) {
	type testCase struct {
		from            patternSet
		to              patternSet
		expectedChanges []FilterChange
	}
	testCases := map[string]testCase{
		"remove one pattern": {
			from: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			to: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
			},
			expectedChanges: []FilterChange{
				{
					Operation:         FilterChangeRemove,
					NamespacePatterns: []string{"ns2"},
					EventTypePattern:  "def",
					SubscriberID:      "uuid1",
				},
			},
		},
		"remove one uuid": {
			from: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			to: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid2"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			expectedChanges: []FilterChange{
				{
					Operation:         FilterChangeRemove,
					NamespacePatterns: []string{"ns1"},
					EventTypePattern:  "abc",
					SubscriberID:      "uuid1",
				},
			},
		},
		"remove all uuids": {
			from: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			to: patternSet{},
			expectedChanges: []FilterChange{
				{
					Operation:         FilterChangeRemove,
					NamespacePatterns: []string{"ns1"},
					EventTypePattern:  "abc",
					SubscriberID:      "uuid1",
				},
				{
					Operation:         FilterChangeRemove,
					NamespacePatterns: []string{"ns1"},
					EventTypePattern:  "abc",
					SubscriberID:      "uuid2",
				},
				{
					Operation:         FilterChangeRemove,
					NamespacePatterns: []string{"ns2"},
					EventTypePattern:  "def",
					SubscriberID:      "uuid1",
				},
			},
		},
		"add one pattern": {
			from: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
			},
			to: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
				pattern{eventTypePattern: "def", namespacePatterns: "ns2"}: []string{"uuid1"},
			},
			expectedChanges: []FilterChange{
				{
					Operation:         FilterChangeAdd,
					NamespacePatterns: []string{"ns2"},
					EventTypePattern:  "def",
					SubscriberID:      "uuid1",
				},
			},
		},
		"add one uuid": {
			from: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1"},
			},
			to: patternSet{
				pattern{eventTypePattern: "abc", namespacePatterns: "ns1"}: []string{"uuid1", "uuid2"},
			},
			expectedChanges: []FilterChange{
				{
					Operation:         FilterChangeAdd,
					NamespacePatterns: []string{"ns1"},
					EventTypePattern:  "abc",
					SubscriberID:      "uuid2",
				},
			},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			fromCNF := ClusterNodeFilter{
				patterns: tt.from,
			}
			toCNF := ClusterNodeFilter{
				patterns: tt.to,
			}
			got := calculateChanges(&fromCNF, &toCNF)
			assert.ElementsMatch(t, tt.expectedChanges, got)
		})
	}
}

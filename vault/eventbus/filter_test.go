// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/stretchr/testify/assert"
)

// TestFilters_AddRemoveMatchLocal checks that basic matching, adding, and removing of patterns all work.
func TestFilters_AddRemoveMatchLocal(t *testing.T) {
	f := NewFilters("self")
	ns := &namespace.Namespace{
		ID:   "ns1",
		Path: "ns1",
	}

	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
	f.addPattern("self", []string{ns.Path}, "abc")
	assert.True(t, f.localMatch(ns, "abc"))
	assert.False(t, f.localMatch(ns, "abcd"))
	assert.True(t, f.anyMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abcd"))
	f.removePattern("self", []string{ns.Path}, "abc")
	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
}

// TestFilters_Watch checks that adding a watch for a cluster will send a notification when the patterns are modified.
func TestFilters_Watch(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	f := NewFilters("self")
	f.addPattern("self", []string{"ns1"}, "e3")
	ch, cancelFunc2, err := f.watch(ctx, "self")
	assert.Nil(t, err)
	t.Cleanup(cancelFunc2)
	initial := <-ch // we always get one immediately for the current state
	assert.Len(t, initial, 2)
	assert.Equal(t, FilterChangeClear, initial[0].Operation)
	assert.Equal(t, FilterChangeAdd, initial[1].Operation)
	assert.Equal(t, []string{"ns1"}, initial[1].NamespacePatterns)
	assert.Equal(t, "e3", initial[1].EventTypePattern)

	go func() {
		f.addPattern("self", []string{"ns1"}, "e2")
	}()
	changes := waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeAdd,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e2",
	}}, changes)
	go func() {
		f.removePattern("self", []string{"ns1"}, "e3")
	}()
	changes = waitForChanges(t, ch)
	assert.Equal(t, []FilterChange{{
		Operation:         FilterChangeRemove,
		NamespacePatterns: []string{"ns1"},
		EventTypePattern:  "e3",
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
	f.addPattern("self", []string{"ns1"}, "e3")
	ch, cancelFunc, err := f.watch(context.Background(), "self")
	assert.Nil(t, err)
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

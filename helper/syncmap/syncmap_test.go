// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package syncmap

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringID struct {
	val string
	id  string
}

func (s stringID) ID() string {
	return s.id
}

var _ IDer = stringID{"", ""}

// TestSyncMap_Get tests that basic getting and putting works.
func TestSyncMap_Get(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "b"})
	assert.Equal(t, stringID{"b", "b"}, m.Get("a"))
	assert.Equal(t, stringID{"", ""}, m.Get("c"))
}

// TestSyncMap_Pop tests that basic Pop operations work.
func TestSyncMap_Pop(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "b"})
	assert.Equal(t, stringID{"b", "b"}, m.Pop("a"))
	assert.Equal(t, stringID{"", ""}, m.Pop("a"))
	assert.Equal(t, stringID{"", ""}, m.Pop("c"))
}

// TestSyncMap_PopIfEqual tests that basic PopIfEqual operations pop only if the IDs are equal.
func TestSyncMap_PopIfEqual(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "c"})
	assert.Equal(t, stringID{"", ""}, m.PopIfEqual("a", "b"))
	assert.Equal(t, stringID{"b", "c"}, m.PopIfEqual("a", "c"))
	assert.Equal(t, stringID{"", ""}, m.PopIfEqual("a", "c"))
}

// TestSyncMap_Clear checks that clearing works as expected and returns a copy of the original map.
func TestSyncMap_Clear(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	assert.Equal(t, map[string]stringID{}, m.data)
	oldMap := m.Clear()
	assert.Equal(t, map[string]stringID{}, m.data)
	assert.Equal(t, map[string]stringID{}, oldMap)

	m.Put("a", stringID{"b", "b"})
	m.Put("c", stringID{"d", "d"})
	oldMap = m.Clear()

	assert.Equal(t, map[string]stringID{"a": {"b", "b"}, "c": {"d", "d"}}, oldMap)
	assert.Equal(t, map[string]stringID{}, m.data)
}

// TestSyncMap_Values checks that the Values method returns an array of the values.
func TestSyncMap_Values(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	assert.Equal(t, []stringID{}, m.Values())
	m.Put("a", stringID{"b", "b"})
	assert.Equal(t, []stringID{{"b", "b"}}, m.Values())
	m.Put("c", stringID{"d", "d"})
	values := m.Values()
	sort.Slice(values, func(i, j int) bool {
		return values[i].val < values[j].val
	})
	assert.Equal(t, []stringID{{"b", "b"}, {"d", "d"}}, values)
}

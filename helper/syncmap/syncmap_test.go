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

func TestSyncMap_Get(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "b"})
	assert.Equal(t, stringID{"b", "b"}, m.Get("a"))
	assert.Equal(t, stringID{"", ""}, m.Get("c"))
}

func TestSyncMap_Pop(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "b"})
	assert.Equal(t, stringID{"b", "b"}, m.Pop("a"))
	assert.Equal(t, stringID{"", ""}, m.Pop("a"))
	assert.Equal(t, stringID{"", ""}, m.Pop("c"))
}

func TestSyncMap_PopIfEqual(t *testing.T) {
	m := NewSyncMap[string, stringID]()
	m.Put("a", stringID{"b", "c"})
	assert.Equal(t, stringID{"", ""}, m.PopIfEqual("a", "b"))
	assert.Equal(t, stringID{"b", "c"}, m.PopIfEqual("a", "c"))
	assert.Equal(t, stringID{"", ""}, m.PopIfEqual("a", "c"))
}

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
	assert.Equal(t, []stringID{{"b", "b"}, {"d", "d"}}, m.Values())
}

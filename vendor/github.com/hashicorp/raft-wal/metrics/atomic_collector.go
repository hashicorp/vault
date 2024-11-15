// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metrics

import "sync/atomic"

var (
	_ Collector = &AtomicCollector{}
)

// AtomicCollector is a simple Collector that atomically stores
// counters and gauges in memory.
type AtomicCollector struct {
	counters []uint64
	gauges   []uint64

	counterIndex, gaugeIndex map[string]int
}

// NewAtomicCollector creates a collector for the given set of Definitions.
func NewAtomicCollector(defs Definitions) *AtomicCollector {
	c := &AtomicCollector{
		counters:     make([]uint64, len(defs.Counters)),
		gauges:       make([]uint64, len(defs.Gauges)),
		counterIndex: make(map[string]int),
		gaugeIndex:   make(map[string]int),
	}
	for i, d := range defs.Counters {
		if _, ok := c.counterIndex[d.Name]; ok {
			panic("duplicate metrics named " + d.Name)
		}
		c.counterIndex[d.Name] = i
	}
	for i, d := range defs.Gauges {
		if _, ok := c.counterIndex[d.Name]; ok {
			panic("duplicate metrics named " + d.Name)
		}
		if _, ok := c.gaugeIndex[d.Name]; ok {
			panic("duplicate metrics named " + d.Name)
		}
		c.gaugeIndex[d.Name] = i
	}
	return c
}

// IncrementCounter record val occurrences of the named event. Names will
// follow prometheus conventions with lower_case_and_underscores. We don't
// need any additional labels currently.
func (c *AtomicCollector) IncrementCounter(name string, delta uint64) {
	id, ok := c.counterIndex[name]
	if !ok {
		panic("invalid metric name: " + name)
	}
	atomic.AddUint64(&c.counters[id], delta)
}

// SetGauge sets the value of the named gauge overriding any previous value.
func (c *AtomicCollector) SetGauge(name string, val uint64) {
	id, ok := c.gaugeIndex[name]
	if !ok {
		panic("invalid metric name: " + name)
	}
	atomic.StoreUint64(&c.gauges[id], val)
}

// Summary returns a summary of the metrics since startup. Each value is
// atomically loaded but the set is not atomic overall and may represent an
// inconsistent snapshot e.g. with some metrics reflecting the most recent
// operation while others don't.
func (c *AtomicCollector) Summary() Summary {
	s := Summary{
		Counters: make(map[string]uint64, len(c.counters)),
		Gauges:   make(map[string]uint64, len(c.gauges)),
	}
	for name, id := range c.counterIndex {
		s.Counters[name] = atomic.LoadUint64(&c.counters[id])
	}
	for name, id := range c.gaugeIndex {
		s.Gauges[name] = atomic.LoadUint64(&c.gauges[id])
	}
	return s
}

// Summary is a copy of the values recorded so far for each metric.
type Summary struct {
	Counters map[string]uint64
	Gauges   map[string]uint64
}

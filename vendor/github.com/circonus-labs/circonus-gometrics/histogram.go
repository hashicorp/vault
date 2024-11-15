// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"fmt"
	"sync"

	"github.com/circonus-labs/circonusllhist"
)

// Histogram measures the distribution of a stream of values.
type Histogram struct {
	name string
	hist *circonusllhist.Histogram
	rw   sync.RWMutex
}

// Timing adds a value to a histogram
func (m *CirconusMetrics) Timing(metric string, val float64) {
	m.SetHistogramValue(metric, val)
}

// RecordValue adds a value to a histogram
func (m *CirconusMetrics) RecordValue(metric string, val float64) {
	m.SetHistogramValue(metric, val)
}

// RecordCountForValue adds count n for value to a histogram
func (m *CirconusMetrics) RecordCountForValue(metric string, val float64, n int64) {
	hist := m.NewHistogram(metric)

	m.hm.Lock()
	hist.rw.Lock()
	hist.hist.RecordValues(val, n)
	hist.rw.Unlock()
	m.hm.Unlock()
}

// SetHistogramValue adds a value to a histogram
func (m *CirconusMetrics) SetHistogramValue(metric string, val float64) {
	hist := m.NewHistogram(metric)

	m.hm.Lock()
	hist.rw.Lock()
	hist.hist.RecordValue(val)
	hist.rw.Unlock()
	m.hm.Unlock()
}

// GetHistogramTest returns the current value for a gauge. (note: it is a function specifically for "testing", disable automatic submission during testing.)
func (m *CirconusMetrics) GetHistogramTest(metric string) ([]string, error) {
	m.hm.Lock()
	defer m.hm.Unlock()

	if hist, ok := m.histograms[metric]; ok {
		return hist.hist.DecStrings(), nil
	}

	return []string{""}, fmt.Errorf("Histogram metric '%s' not found", metric)
}

// RemoveHistogram removes a histogram
func (m *CirconusMetrics) RemoveHistogram(metric string) {
	m.hm.Lock()
	delete(m.histograms, metric)
	m.hm.Unlock()
}

// NewHistogram returns a histogram instance.
func (m *CirconusMetrics) NewHistogram(metric string) *Histogram {
	m.hm.Lock()
	defer m.hm.Unlock()

	if hist, ok := m.histograms[metric]; ok {
		return hist
	}

	hist := &Histogram{
		name: metric,
		hist: circonusllhist.New(),
	}

	m.histograms[metric] = hist

	return hist
}

// Name returns the name from a histogram instance
func (h *Histogram) Name() string {
	return h.name
}

// RecordValue records the given value to a histogram instance
func (h *Histogram) RecordValue(v float64) {
	h.rw.Lock()
	h.hist.RecordValue(v)
	h.rw.Unlock()
}

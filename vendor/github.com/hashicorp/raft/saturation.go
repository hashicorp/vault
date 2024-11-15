// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import (
	"math"
	"time"

	"github.com/armon/go-metrics"
)

// saturationMetric measures the saturation (percentage of time spent working vs
// waiting for work) of an event processing loop, such as runFSM. It reports the
// saturation as a gauge metric (at most) once every reportInterval.
//
// Callers must instrument their loop with calls to sleeping and working, starting
// with a call to sleeping.
//
// Note: the caller must be single-threaded and saturationMetric is not safe for
// concurrent use by multiple goroutines.
type saturationMetric struct {
	reportInterval time.Duration

	// slept contains time for which the event processing loop was sleeping rather
	// than working in the period since lastReport.
	slept time.Duration

	// lost contains time that is considered lost due to incorrect use of
	// saturationMetricBucket (e.g. calling sleeping() or working() multiple
	// times in succession) in the period since lastReport.
	lost time.Duration

	lastReport, sleepBegan, workBegan time.Time

	// These are overwritten in tests.
	nowFn    func() time.Time
	reportFn func(float32)
}

// newSaturationMetric creates a saturationMetric that will update the gauge
// with the given name at the given reportInterval. keepPrev determines the
// number of previous measurements that will be used to smooth out spikes.
func newSaturationMetric(name []string, reportInterval time.Duration) *saturationMetric {
	m := &saturationMetric{
		reportInterval: reportInterval,
		nowFn:          time.Now,
		lastReport:     time.Now(),
		reportFn:       func(sat float32) { metrics.AddSample(name, sat) },
	}
	return m
}

// sleeping records the time at which the loop began waiting for work. After the
// initial call it must always be proceeded by a call to working.
func (s *saturationMetric) sleeping() {
	now := s.nowFn()

	if !s.sleepBegan.IsZero() {
		// sleeping called twice in succession. Count that time as lost rather than
		// measuring nonsense.
		s.lost += now.Sub(s.sleepBegan)
	}

	s.sleepBegan = now
	s.workBegan = time.Time{}
	s.report()
}

// working records the time at which the loop began working. It must always be
// proceeded by a call to sleeping.
func (s *saturationMetric) working() {
	now := s.nowFn()

	if s.workBegan.IsZero() {
		if s.sleepBegan.IsZero() {
			// working called before the initial call to sleeping. Count that time as
			// lost rather than measuring nonsense.
			s.lost += now.Sub(s.lastReport)
		} else {
			s.slept += now.Sub(s.sleepBegan)
		}
	} else {
		// working called twice in succession. Count that time as lost rather than
		// measuring nonsense.
		s.lost += now.Sub(s.workBegan)
	}

	s.workBegan = now
	s.sleepBegan = time.Time{}
	s.report()
}

// report updates the gauge if reportInterval has passed since our last report.
func (s *saturationMetric) report() {
	now := s.nowFn()
	timeSinceLastReport := now.Sub(s.lastReport)

	if timeSinceLastReport < s.reportInterval {
		return
	}

	var saturation float64
	total := timeSinceLastReport - s.lost
	if total != 0 {
		saturation = float64(total-s.slept) / float64(total)
		saturation = math.Round(saturation*100) / 100
	}
	s.reportFn(float32(saturation))

	s.slept = 0
	s.lost = 0
	s.lastReport = now
}

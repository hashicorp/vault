// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package observations

import (
	"context"
	"sync"
)

// TestObservation represents a recorded observation for testing purposes.
// It contains the type of the observation and associated data.
type TestObservation struct {
	Type string
	Data map[string]interface{}
}

// TestObservationRecorder is an implementation of an observation recorder
// that stores observations in memory for testing purposes.
type TestObservationRecorder struct {
	l                  sync.RWMutex
	observationsByType map[string][]*TestObservation
}

// NewTestObservationRecorder creates a new instance of TestObservationRecorder.
func NewTestObservationRecorder() *TestObservationRecorder {
	return &TestObservationRecorder{
		observationsByType: make(map[string][]*TestObservation),
	}
}

func (t *TestObservationRecorder) RecordObservationFromPlugin(_ context.Context, observationType string, data map[string]interface{}) error {
	t.l.Lock()
	defer t.l.Unlock()
	o := &TestObservation{
		Type: observationType,
		Data: data,
	}

	t.observationsByType[observationType] = append(t.observationsByType[observationType], o)
	return nil
}

// NumObservationsByType returns the number of observations recorded of the given type.
func (t *TestObservationRecorder) NumObservationsByType(observationType string) int {
	t.l.RLock()
	defer t.l.RUnlock()

	return len(t.observationsByType[observationType])
}

// ObservationsByType returns all observations recorded of the given type.
// It returns a copy of the slice. If you only need the count, use NumObservationsByType.
func (t *TestObservationRecorder) ObservationsByType(observationType string) []*TestObservation {
	t.l.RLock()
	defer t.l.RUnlock()

	ofType := t.observationsByType[observationType]
	toReturn := make([]*TestObservation, 0, len(ofType))
	for _, o := range ofType {
		toReturn = append(toReturn, o)
	}
	return toReturn
}

// LastObservationOfType returns the last observation recorded of the given type.
// If no observation of that type was recorded, it returns nil.
func (t *TestObservationRecorder) LastObservationOfType(observationType string) *TestObservation {
	t.l.RLock()
	defer t.l.RUnlock()

	ofType := t.observationsByType[observationType]
	if len(ofType) == 0 {
		return nil
	}
	return ofType[len(ofType)-1]
}

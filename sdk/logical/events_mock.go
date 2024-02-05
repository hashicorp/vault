// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"sync"
)

// MockEventSender is a simple implementation of logical.EventSender that simply stores whatever events it receives,
// meant to be used in testing. It is thread-safe.
type MockEventSender struct {
	sync.Mutex
	Events  []MockEvent
	Stopped bool
}

// MockEvent is a container for an event type + event pair.
type MockEvent struct {
	Type  EventType
	Event *EventData
}

// SendEvent implements the logical.EventSender interface.
func (m *MockEventSender) SendEvent(_ context.Context, eventType EventType, event *EventData) error {
	m.Lock()
	defer m.Unlock()
	if !m.Stopped {
		m.Events = append(m.Events, MockEvent{Type: eventType, Event: event})
	}
	return nil
}

func (m *MockEventSender) Stop() {
	m.Lock()
	defer m.Unlock()
	m.Stopped = true
}

var _ EventSender = (*MockEventSender)(nil)

// NewMockEventSender returns a new MockEventSender ready to be used.
func NewMockEventSender() *MockEventSender {
	return &MockEventSender{}
}

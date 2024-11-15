// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"sync"
	"time"
)

// EventType is a string that uniquely identifies the type of an Event within a
// given Broker.
type EventType string

// An Event is analogous to a log entry.
type Event struct {
	// Type of Event
	Type EventType

	// CreatedAt defines the time the event was Sent
	CreatedAt time.Time

	l sync.RWMutex

	// Formatted used by Formatters to store formatted Event data which Sinks
	// can use when writing.  The keys correspond to different formats (json,
	// text, etc).
	Formatted map[string][]byte

	// Payload is the Event's payload data
	Payload interface{}
}

// FormattedAs sets a formatted value for the event, for the specified format
// type.  Any existing value for the type is overwritten.
func (e *Event) FormattedAs(formatType string, formattedValue []byte) {
	e.l.Lock()
	defer e.l.Unlock()
	if e.Formatted == nil {
		e.Formatted = make(map[string][]byte)
	}
	e.Formatted[formatType] = formattedValue
}

// Format will retrieve the formatted value for the specified format type.  The
// two value return allows the caller to determine the existence of the format
// type.
func (e *Event) Format(formatType string) ([]byte, bool) {
	e.l.RLock()
	defer e.l.RUnlock()

	// while not required, this is a more explicit check for nil
	if e.Formatted == nil {
		return nil, false
	}
	v, ok := e.Formatted[formatType]
	return v, ok
}

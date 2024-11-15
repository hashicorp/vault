// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"context"
)

// Predicate is a func that returns true if we want to keep the Event.
type Predicate func(e *Event) (bool, error)

// Filter is a Node that's used for filtering out events from the Pipeline.
type Filter struct {
	// Predicate is a func that returns true if we want to keep the Event.
	Predicate Predicate
	name      string
}

var _ Node = &Filter{}

// Process will call the Filter's Predicate func to determine whether to return
// the Event or filter it out of the Pipeline (Filtered Events return nil, nil,
// which is a successful response).
func (f *Filter) Process(ctx context.Context, e *Event) (*Event, error) {
	// Use the predicate to see if we want to keep the event.
	keep, err := f.Predicate(e)
	if err != nil {
		return nil, err
	}
	if !keep {
		// Return nil to signal that the event should be discarded.
		return nil, nil
	}

	// return the event
	return e, nil
}

// Reopen is a no op for Filters.
func (f *Filter) Reopen() error {
	return nil
}

// Type describes the type of the node as a Filter.
func (f *Filter) Type() NodeType {
	return NodeTypeFilter
}

// Name returns a representation of the Filter's name
func (f *Filter) Name() string {
	return f.name
}

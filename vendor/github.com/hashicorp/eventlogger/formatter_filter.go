// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// JSONFormatterFilter is a Formatter Node which formats the Event as JSON and
// then may filter the event based on the struct used to format the JSON.  This
// is useful when you want to specify filters based on structure of the
// formatted JSON vs the structure of the event.
type JSONFormatterFilter struct {
	Predicate func(e interface{}) (bool, error)
}

var _ Node = &JSONFormatterFilter{}

// Process formats the Event as JSON and stores that formatted data in
// Event.Formatted with a key of "json" and then may filter the event based on
// the struct used to format the JSON.
func (w *JSONFormatterFilter) Process(ctx context.Context, e *Event) (*Event, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	err := enc.Encode(struct {
		CreatedAt time.Time `json:"created_at"`
		EventType `json:"event_type"`
		Payload   interface{} `json:"payload"`
	}{
		e.CreatedAt,
		e.Type,
		e.Payload,
	})
	if err != nil {
		return nil, err
	}

	e.FormattedAs(JSONFormat, buf.Bytes())

	if w.Predicate != nil {
		// Use the predicate to see if we want to keep the event using it's
		// formatted struct as a parmeter to the predicate.
		keep, err := w.Predicate(e)
		if err != nil {
			return nil, err
		}
		if !keep {
			// Return nil to signal that the event should be discarded.
			return nil, nil
		}
	}
	return e, nil
}

// Reopen is a no op
func (w *JSONFormatterFilter) Reopen() error {
	return nil
}

// Type describes the type of the node as a NodeTypeFormatterFilter.
func (w *JSONFormatterFilter) Type() NodeType {
	return NodeTypeFormatterFilter
}

// Name returns a representation of the FormatterFilter's name
func (w *JSONFormatterFilter) Name() string {
	return "JSONFormatteFilter"
}

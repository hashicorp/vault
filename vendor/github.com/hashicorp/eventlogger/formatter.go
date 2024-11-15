// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

const (
	JSONFormat = "json"
)

// JSONFormatter is a Formatter Node which formats the Event as JSON.
type JSONFormatter struct{}

var _ Node = &JSONFormatter{}

// Process formats the Event as JSON and stores that formatted data in
// Event.Formatted with a key of "json"
func (w *JSONFormatter) Process(ctx context.Context, e *Event) (*Event, error) {
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
	return e, nil
}

// Reopen is a no op
func (w *JSONFormatter) Reopen() error {
	return nil
}

// Type describes the type of the node as a Formatter.
func (w *JSONFormatter) Type() NodeType {
	return NodeTypeFormatter
}

// Name returns a representation of the Formatter's name
func (w *JSONFormatter) Name() string {
	return "JSONFormatter"
}

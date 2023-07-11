// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"fmt"

	"github.com/hashicorp/eventlogger"
	"github.com/jefferai/jsonx"
)

var _ eventlogger.Node = (*AuditFormatterJSONX)(nil)

// AuditFormatterJSONX represents a formatter node which will Process JSON to JSONx format.
type AuditFormatterJSONX struct {
	format auditFormat
}

// NewAuditFormatterJSONX creates a formatter node which can be used to format
// incoming events to JSONX.
// This formatter node requires that a AuditFormatterJSON node exists earlier
// in the pipeline and will attempt to access the JSON encoded data stored by that
// formatter node.
func NewAuditFormatterJSONX() *AuditFormatterJSONX {
	return &AuditFormatterJSONX{format: AuditFormatJSONX}
}

// Reopen is a no-op for this formatter node.
func (_ *AuditFormatterJSONX) Reopen() error {
	return nil
}

// Type describes the type of this node.
func (_ *AuditFormatterJSONX) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to retrieve pre-formatted JSON stored within the event
// and re-encode the data to JSONX.
func (f *AuditFormatterJSONX) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	const op = "event.(AuditFormatterJSONX).Process"
	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	// We expect that JSON has already been parsed for this event.
	jsonBytes, ok := e.Format(AuditFormatJSON.String())
	switch {
	case !ok:
		return nil, fmt.Errorf("%s: pre-formatted JSON required but not found: %w", op, ErrInvalidParameter)
	case jsonBytes == nil:
		return nil, fmt.Errorf("%s: pre-formatted JSON required but was nil: %w", op, ErrInvalidParameter)
	}

	xmlBytes, err := jsonx.EncodeJSONBytes(jsonBytes)
	switch {
	case err != nil:
		return nil, fmt.Errorf("%s: unable to encode JSONX using JSON data: %w", op, err)
	case xmlBytes == nil:
		return nil, fmt.Errorf("%s: encoded JSONX was nil: %w", op, err)
	}

	e.FormattedAs(f.format.String(), xmlBytes)

	return e, nil
}

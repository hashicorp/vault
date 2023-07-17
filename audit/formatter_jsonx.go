// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/eventlogger"
	"github.com/jefferai/jsonx"
)

var _ eventlogger.Node = (*AuditFormatterJSONx)(nil)

// AuditFormatterJSONx represents a formatter node which will Process JSON to JSONx format.
type AuditFormatterJSONx struct{}

// NewAuditFormatterJSONx creates a formatter node which can be used to format
// incoming events to JSONx.
// This formatter node requires that a AuditFormatterJSON node exists earlier
// in the pipeline and will attempt to access the JSON encoded data stored by that
// formatter node.
func NewAuditFormatterJSONx() *AuditFormatterJSONx {
	return &AuditFormatterJSONx{}
}

// Reopen is a no-op for this formatter node.
func (_ *AuditFormatterJSONx) Reopen() error {
	return nil
}

// Type describes the type of this node.
func (_ *AuditFormatterJSONx) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to retrieve pre-formatted JSON stored within the event
// and re-encode the data to JSONx.
func (f *AuditFormatterJSONx) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "audit.(AuditFormatterJSONx).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	// We expect that JSON has already been parsed for this event.
	jsonBytes, ok := e.Format(JSONFormat.String())
	if !ok {
		return nil, fmt.Errorf("%s: pre-formatted JSON required but not found: %w", op, event.ErrInvalidParameter)
	}
	if jsonBytes == nil {
		return nil, fmt.Errorf("%s: pre-formatted JSON required but was nil: %w", op, event.ErrInvalidParameter)
	}

	xmlBytes, err := jsonx.EncodeJSONBytes(jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to encode JSONx using JSON data: %w", op, err)
	}
	if xmlBytes == nil {
		return nil, fmt.Errorf("%s: encoded JSONx was nil: %w", op, err)
	}

	e.FormattedAs(JSONxFormat.String(), xmlBytes)

	return e, nil
}

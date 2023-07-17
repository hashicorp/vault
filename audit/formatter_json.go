// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*AuditFormatterJSON)(nil)

// AuditFormatterJSON represents the formatter node which is used to handle
// formatting audit events as JSON.
type AuditFormatterJSON struct {
	config    FormatterConfig
	formatter Formatter
}

// NewAuditFormatterJSON should be used to create an AuditFormatterJSON.
func NewAuditFormatterJSON(config FormatterConfig, salter Salter) (*AuditFormatterJSON, error) {
	const op = "audit.NewAuditFormatterJSON"

	f, err := NewAuditFormatter(salter)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to create new JSON audit formatter: %w", op, err)
	}

	jsonFormatter := &AuditFormatterJSON{
		config:    config,
		formatter: f,
	}

	return jsonFormatter, nil
}

// Reopen is a no-op for a formatter node.
func (_ *AuditFormatterJSON) Reopen() error {
	return nil
}

// Type describes the type of this node (formatter).
func (_ *AuditFormatterJSON) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// Process will attempt to parse the incoming event data into a corresponding
// audit request/response entry which is serialized to JSON and stored within the event.
func (f *AuditFormatterJSON) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "audit.(AuditFormatterJSON).Process"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	a, ok := e.Payload.(*audit)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, event.ErrInvalidParameter)
	}

	var formatted []byte

	switch a.Subtype {
	case AuditRequestType:
		entry, err := f.formatter.FormatRequest(ctx, f.config, a.Data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse request from audit event: %w", op, err)
		}

		formatted, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format request: %w", op, err)
		}
	case AuditResponseType:
		entry, err := f.formatter.FormatResponse(ctx, f.config, a.Data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse response from audit event: %w", op, err)
		}

		formatted, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format response: %w", op, err)
		}
	default:
		return nil, fmt.Errorf("%s: unknown audit event subtype: %q", op, a.Subtype)
	}

	e.FormattedAs(AuditFormatJSON.String(), formatted)

	return e, nil
}

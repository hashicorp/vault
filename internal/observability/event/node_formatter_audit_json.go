// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"fmt"

	vaultaudit "github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"

	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*AuditFormatterJSON)(nil)

// AuditFormatterJSON represents the formatter node which is used to handle
// formatting audit events as JSON.
type AuditFormatterJSON struct {
	config    vaultaudit.FormatterConfig
	format    auditFormat
	formatter vaultaudit.Formatter
}

// NewAuditFormatterJSON should be used to create an AuditFormatterJSON.
func NewAuditFormatterJSON(config vaultaudit.FormatterConfig, salter vaultaudit.Salter) (*AuditFormatterJSON, error) {
	const op = "event.NewAuditFormatterJSON"

	f, err := vaultaudit.NewAuditFormatter(salter)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to create new JSON audit formatter: %w", op, err)
	}

	jsonFormatter := &AuditFormatterJSON{
		format:    AuditFormatJSON,
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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	const op = "event.(AuditFormatterJSON).Process"
	if e == nil {
		return nil, fmt.Errorf("%s: event is nil: %w", op, ErrInvalidParameter)
	}

	a, ok := e.Payload.(*audit)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, ErrInvalidParameter)
	}

	var formatted []byte

	switch a.Subtype {
	case AuditRequest:
		entry, err := f.formatter.FormatRequest(ctx, f.config, a.Data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse request from audit event: %w", op, err)
		}

		formatted, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format request: %w", op, err)
		}
	case AuditResponse:
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

	e.FormattedAs(f.format.String(), formatted)

	return e, nil
}

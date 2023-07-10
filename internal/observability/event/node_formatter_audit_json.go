// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"fmt"

	vaultaudit "github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/sdk/helper/salt"
)

var _ eventlogger.Node = (*AuditFormatterJSON)(nil)

// AuditFormatterJSON represents the formatter node which is used to handle
// formatting audit events as JSON.
type AuditFormatterJSON struct {
	vaultaudit.FormatterConfig
	format    auditFormat
	formatter *vaultaudit.AuditFormatter
	saltFunc  func(context.Context) (*salt.Salt, error)
}

func (f *AuditFormatterJSON) Salt(ctx context.Context) (*salt.Salt, error) {
	return f.saltFunc(ctx)
}

// AuditFormatterConfig represents configuration that may be required by formatter
// nodes which handle audit events.
type AuditFormatterConfig struct {
	vaultaudit.FormatterConfig
	SaltFunc func(context.Context) (*salt.Salt, error)
}

// NewAuditFormatterJSON should be used to create an AuditFormatterJSON.
func NewAuditFormatterJSON(config *AuditFormatterConfig) *AuditFormatterJSON {
	return &AuditFormatterJSON{
		FormatterConfig: config.FormatterConfig,
		saltFunc:        config.SaltFunc,
		format:          AuditFormatJSON,
		formatter:       &vaultaudit.AuditFormatter{},
	}
}

// Reopen is a no-op for this formatter node.
func (_ *AuditFormatterJSON) Reopen() error {
	return nil
}

// Type describes the type of this node.
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

	a, ok := e.Payload.(audit)
	if !ok {
		return nil, fmt.Errorf("%s: cannot parse event payload: %w", op, ErrInvalidParameter)
	}

	var formatted []byte

	switch a.Subtype {
	case AuditRequest:
		entry, err := f.formatter.FormatRequest(ctx, f.FormatterConfig, a.Data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse request from audit event: %w", op, err)
		}

		formatted, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format request: %w", op, err)
		}
	case AuditResponse:
		entry, err := f.formatter.FormatResponse(ctx, f.FormatterConfig, a.Data)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to parse request from audit event: %w", op, err)
		}

		formatted, err = jsonutil.EncodeJSON(entry)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to format request: %w", op, err)
		}
	default:
		return nil, fmt.Errorf("%s: unknown audit event subtype: %q", op, a.Subtype)
	}

	e.FormattedAs(f.format.String(), formatted)

	return e, nil
}

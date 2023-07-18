// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	_ Formatter = (*EventFormatterWriter)(nil)
	_ Writer    = (*EventFormatterWriter)(nil)
)

// Salt returns a new salt with default configuration and no storage usage, and no error.
func (s *nonPersistentSalt) Salt(_ context.Context) (*salt.Salt, error) {
	return salt.NewNonpersistentSalt(), nil
}

// NewEventFormatterWriter should be used to create a new EventFormatterWriter.
// Deprecated: Please move to using eventlogger.Event via EventFormatter and a sink.
func NewEventFormatterWriter(config FormatterConfig, formatter Formatter, writer Writer) (*EventFormatterWriter, error) {
	switch {
	case formatter == nil:
		return nil, errors.New("cannot create a new audit formatter writer with nil formatter")
	case writer == nil:
		return nil, errors.New("cannot create a new audit formatter writer with nil formatter")
	}

	fw := &EventFormatterWriter{
		Formatter: formatter,
		Writer:    writer,
		config:    config,
	}

	return fw, nil
}

// FormatAndWriteRequest attempts to format the specified logical.LogInput into an RequestEntry,
// and then write the request using the specified io.Writer.
// Deprecated: Please move to using eventlogger.Event via EventFormatter and a sink.
func (f *EventFormatterWriter) FormatAndWriteRequest(ctx context.Context, w io.Writer, in *logical.LogInput) error {
	switch {
	case in == nil || in.Request == nil:
		return fmt.Errorf("request to request-audit a nil request")
	case w == nil:
		return fmt.Errorf("writer for audit request is nil")
	case f.Formatter == nil:
		return fmt.Errorf("no formatter specifed")
	case f.Writer == nil:
		return fmt.Errorf("no writer specified")
	}

	reqEntry, err := f.Formatter.FormatRequest(ctx, in)
	if err != nil {
		return err
	}

	return f.Writer.WriteRequest(w, reqEntry)
}

// FormatAndWriteResponse attempts to format the specified logical.LogInput into an ResponseEntry,
// and then write the response using the specified io.Writer.
// Deprecated: Please move to using eventlogger.Event via EventFormatter and a sink.
func (f *EventFormatterWriter) FormatAndWriteResponse(ctx context.Context, w io.Writer, in *logical.LogInput) error {
	switch {
	case in == nil || in.Request == nil:
		return errors.New("request to response-audit a nil request")
	case w == nil:
		return errors.New("writer for audit request is nil")
	case f.Formatter == nil:
		return errors.New("no formatter specified")
	case f.Writer == nil:
		return errors.New("no writer specified")
	}

	respEntry, err := f.FormatResponse(ctx, in)
	if err != nil {
		return err
	}

	return f.Writer.WriteResponse(w, respEntry)
}

// NewTemporaryFormatter creates a formatter not backed by a persistent salt
func NewTemporaryFormatter(requiredFormat, prefix string) *EventFormatterWriter {
	// We can ignore the error from NewEventFormatter since we are sure the salter isn't nil.
	cfg := FormatterConfig{RequiredFormat: format(requiredFormat)}
	eventFormatter, _ := NewEventFormatter(cfg, &nonPersistentSalt{})

	var w Writer

	switch {
	case strings.EqualFold(requiredFormat, JSONxFormat.String()):
		w = &JSONxWriter{Prefix: prefix}
	default:
		w = &JSONWriter{Prefix: prefix}
	}

	// We can ignore the error from NewEventFormatterWriter since we are sure both
	// the formatter and writer are not nil.
	fw, _ := NewEventFormatterWriter(cfg, eventFormatter, w)

	return fw
}

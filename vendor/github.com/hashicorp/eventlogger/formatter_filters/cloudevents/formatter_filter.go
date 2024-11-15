// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudevents

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-secure-stdlib/strutil"
)

const (
	NodeName    = "cloudevents-formatter-filter" // NodeName defines the name of FormatterFilter nodes
	SpecVersion = "1.0"                          // SpecVersion defines the cloudevents spec version supported
	TextIndent  = "  "                           // TextIndent defines the prefix/indent used when encoding FormatText
)

// ID defines an optional single function interface that Event Payloads may
// implement which returns the cloudevent ID for the event payload. If an Event
// Payload doesn't implement this optional interface then a unique ID is
// generated and used for the cloudevent's ID.
type ID interface {
	// ID returns the cloudevent ID
	ID() string
}

// Data defines an optional single function interface that Event Payloads may
// implement which returns the cloudevent data for the event payload. If an
// Event doesn't implement this optional interface then the entire Event Payload
// is used as the cloudevent data.
type Data interface {
	// Data returns the cloudevent Data.
	Data() interface{}
}

// Event defines type which is used when formatting cloudevents.
//
// For more info on the fields see: https://github.com/cloudevents/spec)
type Event struct {
	// ID identifies the event, cannot be an empty and is required.  The
	// combination of Source + ID must be unique.  Events with the same Source +
	// ID can be assumed to be duplicates by consumers
	ID string `json:"id"`

	// Source identifies the context in which the event happened, it is a
	// URI-reference, cannot be empty and is required.
	Source string `json:"source"`

	// SpecVersion defines the version of CloudEvents that the event is using,
	// it cannot be empty and is required.
	SpecVersion string `json:"specversion"`

	// Type defines the event's type, cannot be empty and is required.
	Type string `json:"type"`

	// Data may include domain-specific information about the event and is
	// optional.
	Data interface{} `json:"data,omitempty"`

	// DataContentType defines the content type of the event's data value and is
	// optional.  If present it must adhere to:
	// https://datatracker.ietf.org/doc/html/rfc2046
	DataContentType string `json:"datacontentype,omitempty"`

	// DataSchema is a URI-reference and is optional.
	DataSchema string `json:"dataschema,omitempty"`

	// Time is in format RFC 3339 (the default for time.Time) and is optional
	Time time.Time `json:"time,omitempty"`

	// Serialized is optional and will contain the serialized data that was
	// "signed" (see: FormatterFilter.Signer)
	Serialized string `json:"serialized,omitempty"`

	// SerializedHmac is optional and will contain the signature of the
	// serialized field (see: FormatterFilter.Signer)
	SerializedHmac string `json:"serialized_hmac,omitempty"`
}

// FormatterFilter is a Node which formats the Event as a CloudEvent in JSON
// format (See: https://github.com/cloudevents/spec)
type FormatterFilter struct {
	// Source identifies the context where the cloudevents happen and is
	// required
	Source *url.URL

	// Schema is the JSON schema for the cloudevent data and is optional
	Schema *url.URL

	// Format defines the format created by the node.  If empty (unspecified),
	// FormatJSON will be used
	Format Format

	// Predicate is a func that returns true if we want to keep the cloudevent.
	// The context parameter is the context of Process(ctx context.Context, e
	// *eventlogger.Event) and the interface{} parameter will be a
	// cloudevents.Event struct.
	Predicate func(ctx context.Context, cloudevent interface{}) (bool, error)

	// Signer provides an optional signer for "signing" formatted events.  If
	// not nil, then formatted events will be "signed" using this signer and
	// the event's Serialized and SerializedHmac fields will be populated.
	//	Serialized: which will contain the serialized data that was "signed"
	//	SerializedHmac: which contains the signature of the serialized field
	Signer Signer

	// SignEventTypes contains a list of event types which should be signed by
	// the Signer
	SignEventTypes []string
}

var _ eventlogger.Node = &FormatterFilter{}

func (f *FormatterFilter) validate() error {
	const op = "cloudevents.(FormatterFilter).validate"
	if f == nil {
		return fmt.Errorf("%s: missing formatter filter: %w", op, eventlogger.ErrInvalidParameter)
	}
	if f.Source == nil || f.Source.String() == "" {
		return fmt.Errorf("%s: missing source: %w", op, eventlogger.ErrInvalidParameter)
	}
	if err := f.Format.validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if f.Schema != nil && f.Schema.String() == "" {
		return fmt.Errorf("%s: an empty schema is not valid: %w", op, eventlogger.ErrInvalidParameter)
	}
	return nil
}

// Process formats the Event as a cloudevent and stores that formatted data in
// Event.Formatted0 with a key of either "cloudevents-json"
// (cloudevents.FormatJSON) or "cloudevents-text" (cloudevents.FormatText) based
// on the FormatterFilter.Format value. If the node has a Predicate, then the
// filter will be applied to the resulting CloudEvent.
func (f *FormatterFilter) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	const op = "cloudevents.(FormatterFilter).Process"
	if err := f.validate(); err != nil {
		return nil, fmt.Errorf("%s: invalid formatter filter %w", op, err)
	}
	if e == nil {
		return nil, fmt.Errorf("%s: missing event: %w", op, eventlogger.ErrInvalidParameter)
	}

	var data interface{}
	if i, ok := e.Payload.(Data); ok {
		data = i.Data()
	} else {
		data = e.Payload
	}
	var id string
	if i, ok := e.Payload.(ID); ok {
		id = i.ID()
		if id == "" {
			return nil, fmt.Errorf("%s: returned ID() is empty: %w", op, eventlogger.ErrInvalidParameter)
		}
	} else {
		var err error
		id, err = newId()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	var schema string
	if f.Schema != nil {
		schema = f.Schema.String()
	}

	ce := Event{
		ID:          id,
		Source:      f.Source.String(),
		SpecVersion: SpecVersion,
		Type:        string(e.Type),
		Data:        data,
		DataSchema:  schema,
		Time:        e.CreatedAt,
	}
	switch f.Format {
	case FormatJSON, FormatUnspecified:
		ce.DataContentType = DataContentTypeCloudEvents
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err := enc.Encode(ce); err != nil {
			return nil, fmt.Errorf("%s: error formatting as JSON: %w", op, err)
		}
		f.sign(ctx, &ce, enc, buf)
		e.FormattedAs(string(FormatJSON), buf.Bytes())
	case FormatText:
		ce.DataContentType = DataContentTypeText
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetIndent("", TextIndent)
		if err := enc.Encode(ce); err != nil {
			return nil, fmt.Errorf("%s: error formatting as text: %w", op, err)
		}
		f.sign(ctx, &ce, enc, buf)
		e.FormattedAs(string(FormatText), buf.Bytes())
	default:
		// this should be unreachable since f.validate() should catch this error
		// condition at the top of the function.
		return nil, fmt.Errorf("%s: %s is not a supported format: %w", op, f.Format, eventlogger.ErrInvalidParameter)
	}

	if f.Predicate != nil {
		// Use the predicate to see if we want to keep the event using it's
		// formatted struct as a parmeter to the predicate.
		keep, err := f.Predicate(ctx, ce)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to filter: %w", op, err)
		}
		if !keep {
			// Return nil to signal that the event should be discarded.
			return nil, nil
		}
	}
	return e, nil
}

// Reopen is a no op
func (f *FormatterFilter) Reopen() error {
	return nil
}

// Type describes the type of the node as a Formatter.
func (f *FormatterFilter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatterFilter
}

// Name returns a representation of the Formatter's name
func (f *FormatterFilter) Name() string {
	return NodeName
}

func (f *FormatterFilter) sign(ctx context.Context, e *Event, enc *json.Encoder, buf *bytes.Buffer) error {
	const op = "cloudevents.(FormatterFilter).sign"
	if e == nil {
		return fmt.Errorf("%s: missing event: %w", op, eventlogger.ErrInvalidParameter)
	}
	if enc == nil {
		return fmt.Errorf("%s: missing encoder: %w", op, eventlogger.ErrInvalidParameter)
	}
	if buf == nil {
		return fmt.Errorf("%s: missing buffer: %w", op, eventlogger.ErrInvalidParameter)
	}
	if f.Signer != nil && strutil.StrListContains(f.SignEventTypes, e.Type) {
		bufHmac, err := f.Signer(ctx, buf.Bytes())
		if err != nil {
			return fmt.Errorf("%s: unable to sign: %w", op, err)
		}
		e.Serialized = base64.RawURLEncoding.EncodeToString(buf.Bytes())
		e.SerializedHmac = bufHmac
		buf.Reset()
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("%s: error formatting as JSON: %w", op, err)
		}
	}
	return nil
}

// Signer defines a function for "signing" an event
type Signer func(context.Context, []byte) (string, error)

// Rotate supports rotating the filter's signer which is used to "sign"
// formatted events
func (f *FormatterFilter) Rotate(s Signer) error {
	const op = "cloudevents.(FormatterFilter).Rotate"
	if s == nil {
		return fmt.Errorf("%s: missing signer: %w", op, eventlogger.ErrInvalidParameter)
	}
	f.Signer = s
	return nil
}

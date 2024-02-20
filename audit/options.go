// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-bexpr"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/mitchellh/pointerstructure"
)

// Option is how options are passed as arguments.
type Option func(*options) error

// options are used to represent configuration for audit related nodes.
type options struct {
	withID              string
	withNow             time.Time
	withSubtype         subtype
	withFormat          format
	withPrefix          string
	withRaw             bool
	withElision         bool
	withOmitTime        bool
	withHMACAccessor    bool
	withHeaderFormatter HeaderFormatter
	withExclusions      []*exclusion
}

// exclusion represents an optional condition and fields which should be excluded
// from audit entries.
type exclusion struct {
	Evaluator *bexpr.Evaluator `json:"condition,omitempty"`
	Fields    []string         `json:"fields,omitempty"`
}

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withNow:          time.Now(),
		withFormat:       JSONFormat,
		withHMACAccessor: true,
	}
}

// getOpts applies each supplied Option and returns the fully configured options.
// Each Option is applied in the order it appears in the argument list, so it is
// possible to supply the same Option numerous times and the 'last write wins'.
func getOpts(opt ...Option) (options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return options{}, err
		}
	}
	return opts, nil
}

// WithID provides an optional ID.
func WithID(id string) Option {
	return func(o *options) error {
		var err error

		id := strings.TrimSpace(id)
		switch {
		case id == "":
			err = errors.New("id cannot be empty")
		default:
			o.withID = id
		}

		return err
	}
}

// WithNow provides an Option to represent 'now'.
func WithNow(now time.Time) Option {
	return func(o *options) error {
		var err error

		switch {
		case now.IsZero():
			err = errors.New("cannot specify 'now' to be the zero time instant")
		default:
			o.withNow = now
		}

		return err
	}
}

// WithSubtype provides an Option to represent the event subtype.
func WithSubtype(s string) Option {
	return func(o *options) error {
		s := strings.TrimSpace(s)
		if s == "" {
			return errors.New("subtype cannot be empty")
		}
		parsed := subtype(s)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withSubtype = parsed
		return nil
	}
}

// WithFormat provides an Option to represent event format.
func WithFormat(f string) Option {
	return func(o *options) error {
		f := strings.TrimSpace(f)
		if f == "" {
			// Return early, we won't attempt to apply this option if its empty.
			return nil
		}

		parsed := format(f)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withFormat = parsed
		return nil
	}
}

// WithPrefix provides an Option to represent a prefix for a file sink.
func WithPrefix(prefix string) Option {
	return func(o *options) error {
		o.withPrefix = prefix

		return nil
	}
}

// WithRaw provides an Option to represent whether 'raw' is required.
func WithRaw(r bool) Option {
	return func(o *options) error {
		o.withRaw = r
		return nil
	}
}

// WithElision provides an Option to represent whether elision (...) is required.
func WithElision(e bool) Option {
	return func(o *options) error {
		o.withElision = e
		return nil
	}
}

// WithOmitTime provides an Option to represent whether to omit time.
func WithOmitTime(t bool) Option {
	return func(o *options) error {
		o.withOmitTime = t
		return nil
	}
}

// WithHMACAccessor provides an Option to represent whether an HMAC accessor is applicable.
func WithHMACAccessor(h bool) Option {
	return func(o *options) error {
		o.withHMACAccessor = h
		return nil
	}
}

// WithHeaderFormatter provides an Option to supply a HeaderFormatter.
// If the HeaderFormatter interface supplied is nil (type or value), the option will not be applied.
func WithHeaderFormatter(f HeaderFormatter) Option {
	return func(o *options) error {
		if f != nil && !reflect.ValueOf(f).IsNil() {
			o.withHeaderFormatter = f
		}

		return nil
	}
}

// WithExclusions provides an Option to supply exclusions in a JSON string format.
// 'condition' is an optional value in the JSON.
// See 'exclusion' type for more information and example below:
// Expected JSON format:
//
//	[
//		{
//			"condition": "\"/request/mount_type\" == transit",
//			"fields": [ "/request/data", "/response/data" ]
//		},
//		{
//			"condition":  "\"/request/mount_type\" == userpass",
//			"fields": [ "/request/data" ]
//		}
//	]
func WithExclusions(e string) Option {
	return func(o *options) error {
		e = strings.TrimSpace(e)
		if e == "" {
			return nil
		}

		var result []*exclusion

		err := json.Unmarshal([]byte(e), &result)
		if err != nil {
			return fmt.Errorf("unable to parse exclusions: %w", err)
		}

		// Validate the exclusions
		for i, exc := range result {
			if err := exc.validate(); err != nil {
				return fmt.Errorf("exclusion (#%d) validation failure: %w", i+1, err)
			}
		}

		o.withExclusions = result

		return nil
	}
}

// UnmarshalJSON handles unmarshalling JSON bytes (string representation) of a collection
// of exclusion types into a Go type.
func (e *exclusion) UnmarshalJSON(b []byte) error {
	// Reference the JSON struct tags for exclusion.
	const keyFields = "fields"
	const keyCondition = "condition"

	var err error

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		return err
	}

	// Parse fields
	f, ok := m[keyFields]
	if !ok {
		return fmt.Errorf("exclusion '%s' missing", keyFields)
	}
	intermediateFields, ok := f.([]any)
	if !ok {
		return fmt.Errorf("unable to parse '%s': expected collection of %s; got: '%v'", keyFields, keyFields, f)
	}
	var fields []string
	for _, v := range intermediateFields {
		v, ok := v.(string)
		if !ok {
			continue
		}
		s := strings.TrimSpace(v)
		if s != "" {
			fields = append(fields, s)
		}
	}
	if len(fields) < 1 {
		return fmt.Errorf("exclusion '%s' cannot be empty", keyFields)
	}

	// Set the fields now, so we can return early if we don't have an optional condition.
	e.Fields = fields

	// Optional condition
	var eval *bexpr.Evaluator
	c, ok := m[keyCondition]
	if !ok {
		// Return early as we've already set the exclusion.Fields
		return nil
	}
	v, ok := c.(string)
	if !ok {
		return nil
	}
	condition := strings.TrimSpace(v)
	if condition == "" {
		// Return early as we've already set the exclusion.Fields
		return nil
	}

	// Create the evaluator and use the current JSON struct tags.
	eval, err = bexpr.CreateEvaluator(condition, bexpr.WithTagName("json"))
	if err != nil {
		return fmt.Errorf("unable to parse expression '%s' from '%s': %w", condition, keyCondition, err)
	}

	// Set the condition to the new evaluator.
	e.Evaluator = eval

	return nil
}

// validate attempts to ensure an exclusion has fields that should be excluded,
// and that they can be represented as JSON pointers.
// NOTE: Due to the dynamic nature of parts of the request and response, we are
// unable to validate that the expression condition works as intended.
// e.g. /request/data/myKey would not appear in a vanilla RequestEntry but is
// valid and should match when the response data contains that key.
func (e *exclusion) validate() error {
	if len(e.Fields) < 1 {
		return fmt.Errorf("exclusion doesn't contain any fields: %w", event.ErrInvalidParameter)
	}

	for _, field := range e.Fields {
		if _, err := pointerstructure.Parse(field); err != nil {
			return fmt.Errorf("unable to parse field '%s': %w", field, err)
		}
	}

	return nil
}

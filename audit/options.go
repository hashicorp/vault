// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withNow:    time.Now(),
		withFormat: JSONFormat,
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
		prefix = strings.TrimSpace(prefix)

		if prefix != "" {
			o.withPrefix = prefix
		}

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

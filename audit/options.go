// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"errors"
	"strings"
	"time"
)

// option is how options are passed as arguments.
type option func(*options) error

// options are used to represent configuration for a audit related nodes.
type options struct {
	withID           string
	withNow          time.Time
	withSubtype      subtype
	withFormat       format
	withPrefix       string
	withRaw          bool
	withElision      bool
	withOmitTime     bool
	withHMACAccessor bool
}

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withNow:          time.Now(),
		withFormat:       jsonFormat,
		withHMACAccessor: true,
	}
}

// getOpts applies each supplied option and returns the fully configured options.
// Each option is applied in the order it appears in the argument list, so it is
// possible to supply the same option numerous times and the 'last write wins'.
func getOpts(opt ...option) (options, error) {
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

// withID provides an optional ID.
func withID(id string) option {
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

// withNow provides an option to represent 'now'.
func withNow(now time.Time) option {
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

// withSubtype provides an option to represent the event subtype.
func withSubtype(s string) option {
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

// withFormat provides an option to represent event format.
func withFormat(f string) option {
	return func(o *options) error {
		f := strings.TrimSpace(strings.ToLower(f))
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

// withPrefix provides an option to represent a prefix for a file sink.
func withPrefix(prefix string) option {
	return func(o *options) error {
		o.withPrefix = prefix

		return nil
	}
}

// withRaw provides an option to represent whether 'raw' is required.
func withRaw(r bool) option {
	return func(o *options) error {
		o.withRaw = r
		return nil
	}
}

// withElision provides an option to represent whether elision (...) is required.
func withElision(e bool) option {
	return func(o *options) error {
		o.withElision = e
		return nil
	}
}

// withOmitTime provides an option to represent whether to omit time.
func withOmitTime(t bool) option {
	return func(o *options) error {
		o.withOmitTime = t
		return nil
	}
}

// withHMACAccessor provides an option to represent whether an HMAC accessor is applicable.
func withHMACAccessor(h bool) option {
	return func(o *options) error {
		o.withHMACAccessor = h
		return nil
	}
}

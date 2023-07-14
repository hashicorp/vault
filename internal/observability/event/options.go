// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
)

// Option is how Options are passed as arguments.
type Option func(*options) error

// options are used to represent configuration for an Event.
type options struct {
	withID       string
	withNow      time.Time
	withSubtype  auditSubtype
	withFormat   auditFormat
	withFileMode *os.FileMode
	withPrefix   string
	withFacility string
	withTag      string
}

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withNow: time.Now(),
	}
}

// getOpts applies all the supplied Option and returns configured options.
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

// NewID is a bit of a modified NewID has been done to stop a circular
// dependency with the errors package that is caused by importing
// boundary/internal/db
func NewID(prefix string) (string, error) {
	const op = "event.NewID"
	if prefix == "" {
		return "", fmt.Errorf("%s: missing prefix: %w", op, ErrInvalidParameter)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("%s: unable to generate ID: %w", op, err)
	}

	return fmt.Sprintf("%s_%s", prefix, id), nil
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

// WithNow provides an option to represent 'now'.
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

// WithSubtype provides an option to represent the subtype.
func WithSubtype(subtype string) Option {
	return func(o *options) error {
		s := strings.TrimSpace(subtype)
		if s == "" {
			return errors.New("subtype cannot be empty")
		}

		parsed := auditSubtype(s)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withSubtype = parsed
		return nil
	}
}

// WithFormat provides an option to represent event format.
func WithFormat(format string) Option {
	return func(o *options) error {
		f := strings.TrimSpace(format)
		if f == "" {
			return errors.New("format cannot be empty")
		}

		parsed := auditFormat(f)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withFormat = parsed
		return nil
	}
}

// WithFileMode provides an option to represent a file mode for a file sink.
// Supplying an empty string or whitespace will prevent this option from being
// applied, but it will not return an error in those circumstances.
func WithFileMode(mode string) Option {
	return func(o *options) error {
		// If supplied file mode is empty, just return early without setting anything.
		// We can assume that this option was called by something that didn't
		// parse the incoming value, perhaps from a config map etc.
		mode = strings.TrimSpace(mode)
		if mode == "" {
			return nil
		}

		// By now we believe we have something that the caller really intended to
		// be parsed into a file mode.
		raw, err := strconv.ParseUint(mode, 8, 32)

		switch {
		case err != nil:
			return fmt.Errorf("unable to parse file mode: %w", err)
		default:
			m := os.FileMode(raw)
			o.withFileMode = &m
		}

		return nil
	}
}

// WithPrefix provides an option to represent a prefix for a file sink.
func WithPrefix(prefix string) Option {
	return func(o *options) error {
		o.withPrefix = prefix
		return nil
	}
}

// WithFacility provides an option to represent a 'facility' for a syslog sink.
func WithFacility(facility string) Option {
	return func(o *options) error {
		o.withFacility = strings.TrimSpace(facility)

		return nil
	}
}

// WithTag provides an option to represent a 'tag' for a syslog sink.
func WithTag(tag string) Option {
	return func(o *options) error {
		o.withTag = strings.TrimSpace(tag)

		return nil
	}
}

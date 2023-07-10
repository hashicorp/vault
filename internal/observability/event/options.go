// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
)

// Option is how Options are passed as arguments.
type Option func(*options) error

// options are used to represent configuration for an Event.
type options struct {
	withID      string
	withNow     time.Time
	withSubtype string
	withFormat  string
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

// WithID allows an optional ID.
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

// WithNow allows an option to represent 'now'.
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

// WithSubtype allows an option to represent the subtype.
func WithSubtype(subtype string) Option {
	return func(o *options) error {
		var err error

		subtype := strings.TrimSpace(subtype)
		switch {
		case subtype == "":
			err = errors.New("subtype cannot be empty")
		default:
			o.withSubtype = subtype
		}

		return err
	}
}

// WithFormat allows an option to represent event format.
func WithFormat(format string) Option {
	return func(o *options) error {
		var err error

		format := strings.TrimSpace(format)
		switch {
		case format == "":
			err = errors.New("format cannot be empty")
		default:
			o.withFormat = format
		}

		return err
	}
}

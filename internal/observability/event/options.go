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
type Option func(*Options) error

// Options are used to represent configuration for an Event.
type Options struct {
	WithID  string
	WithNow time.Time
}

// getDefaultOptions returns Options with their default values.
func getDefaultOptions() Options {
	return Options{
		WithNow: time.Now(),
	}
}

// GetOpts applies all the supplied Option and returns configured Options.
// Each Option is applied in the order it appears in the argument list, so it is
// possible to supply the same Option numerous times and the 'last write wins'.
func GetOpts(opt ...Option) (Options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return Options{}, err
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
	return func(o *Options) error {
		var err error

		id := strings.TrimSpace(id)
		switch {
		case id == "":
			err = errors.New("id cannot be empty")
		default:
			o.WithID = id
		}

		return err
	}
}

// WithNow provides an option to represent 'now'.
func WithNow(now time.Time) Option {
	return func(o *Options) error {
		var err error

		switch {
		case now.IsZero():
			err = errors.New("cannot specify 'now' to be the zero time instant")
		default:
			o.WithNow = now
		}

		return err
	}
}

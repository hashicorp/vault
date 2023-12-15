// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"

	"github.com/hashicorp/go-uuid"
)

// Option is how Options are passed as arguments.
type Option func(*options) error

// Options are used to represent configuration for an Event.
type options struct {
	withID          string
	withNow         time.Time
	withFacility    string
	withTag         string
	withSocketType  string
	withMaxDuration time.Duration
	withFileMode    *os.FileMode
}

// getDefaultOptions returns Options with their default values.
func getDefaultOptions() options {
	fileMode := os.FileMode(0o600)

	return options{
		withNow:         time.Now(),
		withFacility:    "AUTH",
		withTag:         "vault",
		withSocketType:  "tcp",
		withMaxDuration: 2 * time.Second,
		withFileMode:    &fileMode,
	}
}

// getOpts applies all the supplied Option and returns configured Options.
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

// WithFacility provides an Option to represent a 'facility' for a syslog sink.
func WithFacility(facility string) Option {
	return func(o *options) error {
		facility = strings.TrimSpace(facility)

		if facility != "" {
			o.withFacility = facility
		}

		return nil
	}
}

// WithTag provides an Option to represent a 'tag' for a syslog sink.
func WithTag(tag string) Option {
	return func(o *options) error {
		tag = strings.TrimSpace(tag)

		if tag != "" {
			o.withTag = tag
		}

		return nil
	}
}

// WithSocketType provides an Option to represent the socket type for a socket sink.
func WithSocketType(socketType string) Option {
	return func(o *options) error {
		socketType = strings.TrimSpace(socketType)

		if socketType != "" {
			o.withSocketType = socketType
		}

		return nil
	}
}

// WithMaxDuration provides an Option to represent the max duration for writing to a socket.
func WithMaxDuration(duration string) Option {
	return func(o *options) error {
		duration = strings.TrimSpace(duration)

		if duration == "" {
			return nil
		}

		parsed, err := parseutil.ParseDurationSecond(duration)
		if err != nil {
			return err
		}

		o.withMaxDuration = parsed

		return nil
	}
}

// WithFileMode provides an Option to represent a file mode for a file sink.
// Supplying an empty string or whitespace will prevent this Option from being
// applied, but it will not return an error in those circumstances.
func WithFileMode(mode string) Option {
	return func(o *options) error {
		// If supplied file mode is empty, just return early without setting anything.
		// We can assume that this Option was called by something that didn't
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

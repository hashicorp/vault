package audit

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// getDefaultOptions returns options with their default values.
func getDefaultOptions() options {
	return options{
		withNow: time.Now(),
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
			return errors.New("format cannot be empty")
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

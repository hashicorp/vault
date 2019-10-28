package config

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrWaitStringEmpty is the error returned when wait is specified as an empty
	// string.
	ErrWaitStringEmpty = errors.New("wait: cannot be empty")

	// ErrWaitInvalidFormat is the error returned when the wait is specified
	// incorrectly.
	ErrWaitInvalidFormat = errors.New("wait: invalid format")

	// ErrWaitNegative is the error returned with the wait is negative.
	ErrWaitNegative = errors.New("wait: cannot be negative")

	// ErrWaitMinLTMax is the error returned with the minimum wait time is not
	// less than the maximum wait time.
	ErrWaitMinLTMax = errors.New("wait: min must be less than max")
)

// WaitConfig is the Min/Max duration used by the Watcher
type WaitConfig struct {
	// Enabled determines if this wait is enabled.
	Enabled *bool `mapstructure:"bool"`

	// Min and Max are the minimum and maximum time, respectively, to wait for
	// data changes before rendering a new template to disk.
	Min *time.Duration `mapstructure:"min"`
	Max *time.Duration `mapstructure:"max"`
}

// DefaultWaitConfig is the default configuration.
func DefaultWaitConfig() *WaitConfig {
	return &WaitConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *WaitConfig) Copy() *WaitConfig {
	if c == nil {
		return nil
	}

	var o WaitConfig
	o.Enabled = c.Enabled
	o.Min = c.Min
	o.Max = c.Max
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *WaitConfig) Merge(o *WaitConfig) *WaitConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Min != nil {
		r.Min = o.Min
	}

	if o.Max != nil {
		r.Max = o.Max
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *WaitConfig) Finalize() {
	if c.Enabled == nil {
		c.Enabled = Bool(TimeDurationPresent(c.Min))
	}

	if c.Min == nil {
		c.Min = TimeDuration(0 * time.Second)
	}

	if c.Max == nil {
		c.Max = TimeDuration(4 * *c.Min)
	}
}

// GoString defines the printable version of this struct.
func (c *WaitConfig) GoString() string {
	if c == nil {
		return "(*WaitConfig)(nil)"
	}

	return fmt.Sprintf("&WaitConfig{"+
		"Enabled:%s, "+
		"Min:%s, "+
		"Max:%s"+
		"}",
		BoolGoString(c.Enabled),
		TimeDurationGoString(c.Min),
		TimeDurationGoString(c.Max),
	)
}

// ParseWaitConfig parses a string of the format `minimum(:maximum)` into a
// WaitConfig.
func ParseWaitConfig(s string) (*WaitConfig, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return nil, ErrWaitStringEmpty
	}

	parts := strings.Split(s, ":")

	var min, max time.Duration
	var err error

	switch len(parts) {
	case 1:
		min, err = time.ParseDuration(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		max = 4 * min
	case 2:
		min, err = time.ParseDuration(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}

		max, err = time.ParseDuration(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrWaitInvalidFormat
	}

	if min < 0 || max < 0 {
		return nil, ErrWaitNegative
	}

	if max < min {
		return nil, ErrWaitMinLTMax
	}

	var c WaitConfig
	c.Min = TimeDuration(min)
	c.Max = TimeDuration(max)

	return &c, nil
}

// WaitVar implements the Flag.Value interface and allows the user to specify
// a watch interval using Go's flag parsing library.
type WaitVar WaitConfig

// Set sets the value in the format min[:max] for a wait timer.
func (w *WaitVar) Set(value string) error {
	wait, err := ParseWaitConfig(value)
	if err != nil {
		return err
	}

	w.Min = wait.Min
	w.Max = wait.Max

	return nil
}

// String returns the string format for this wait variable
func (w *WaitVar) String() string {
	return fmt.Sprintf("%s:%s", w.Min, w.Max)
}

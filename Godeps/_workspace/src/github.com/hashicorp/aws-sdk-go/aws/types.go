package aws

import (
	"math"
	"strconv"
	"time"
)

// A StringValue is a string which may or may not be present.
type StringValue *string

// String converts a Go string into a StringValue.
func String(v string) StringValue {
	return &v
}

// A BooleanValue is a boolean which may or may not be present.
type BooleanValue *bool

// Boolean converts a Go bool into a BooleanValue.
func Boolean(v bool) BooleanValue {
	return &v
}

// True is the BooleanValue equivalent of the Go literal true.
func True() BooleanValue {
	return Boolean(true)
}

// False is the BooleanValue equivalent of the Go literal false.
func False() BooleanValue {
	return Boolean(false)
}

// An IntegerValue is an integer which may or may not be present.
type IntegerValue *int

// Integer converts a Go int into an IntegerValue.
func Integer(v int) IntegerValue {
	return &v
}

// A LongValue is a 64-bit integer which may or may not be present.
type LongValue *int64

// Long converts a Go int64 into a LongValue.
func Long(v int64) LongValue {
	return &v
}

// A FloatValue is a 32-bit floating point number which may or may not be
// present.
type FloatValue *float32

// Float converts a Go float32 into a FloatValue.
func Float(v float32) FloatValue {
	return &v
}

// A DoubleValue is a 64-bit floating point number which may or may not be
// present.
type DoubleValue *float64

// Double converts a Go float64 into a DoubleValue.
func Double(v float64) DoubleValue {
	return &v
}

// A UnixTimestamp is a Unix timestamp represented as fractional seconds since
// the Unix epoch.
type UnixTimestamp struct {
	Time time.Time
}

// MarshalJSON marshals the timestamp as a float.
func (t UnixTimestamp) MarshalJSON() (text []byte, err error) {
	n := float64(t.Time.UnixNano()) / 1e9
	s := strconv.FormatFloat(n, 'f', -1, 64)
	return []byte(s), nil
}

// UnmarshalJSON unmarshals the timestamp from a float.
func (t *UnixTimestamp) UnmarshalJSON(text []byte) error {
	f, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}

	sec := math.Floor(f)
	nsec := (f - sec) * 1e9

	t.Time = time.Unix(int64(sec), int64(nsec)).UTC()
	return nil
}

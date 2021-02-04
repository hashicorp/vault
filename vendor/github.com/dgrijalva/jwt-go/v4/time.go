package jwt

import (
	"encoding/json"
	"reflect"
	"time"
)

// TimePrecision determines how precisely time is measured
// by this library. When serializing and deserialzing tokens,
// time values are automatically truncated to this precision.
// See the time package's Truncate method for more detail
const TimePrecision = time.Microsecond

// Time is how this library represents time values. It's mostly
// a wrapper for the standard library's time.Time, but adds
// specialized JSON decoding behavior to interop with the way
// time is represented by JWT. Also makes it possible to represent
// nil values.
type Time struct {
	time.Time
}

// NewTime creates a new Time value from a float64, following
// the JWT spec.
func NewTime(t float64) *Time {
	return At(time.Unix(0, int64(t*float64(time.Second))))
}

// Now returns a new Time value using the current time.
// You can override Now by changing the value of TimeFunc
func Now() *Time {
	return At(TimeFunc())
}

// At makes a Time value from a standard library time.Time value
func At(at time.Time) *Time {
	return &Time{at.Truncate(TimePrecision)}
}

// ParseTime is used for creating a Time value from various
// possible representations that can occur in serialization.
func ParseTime(value interface{}) (*Time, error) {
	switch v := value.(type) {
	case int64:
		return NewTime(float64(v)), nil
	case float64:
		return NewTime(v), nil
	case json.Number:
		vv, err := v.Float64()
		if err != nil {
			return nil, err
		}
		return NewTime(vv), nil
	case nil:
		return nil, nil
	default:
		return nil, &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
	}
}

// UnmarshalJSON implements the json package's Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var value json.Number
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	v, err := ParseTime(value)
	*t = *v
	return err
}

// MarshalJSON implements the json package's Marshaler interface
func (t *Time) MarshalJSON() ([]byte, error) {
	f := float64(t.Truncate(TimePrecision).UnixNano()) / float64(time.Second)
	return json.Marshal(f)
}

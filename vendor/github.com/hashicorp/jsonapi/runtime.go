package jsonapi

import (
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"time"
)

// Event represents a lifecycle event in the marshaling or unmarshalling
// process.
type Event int

const (
	// UnmarshalStart is the Event that is sent when deserialization of a payload
	// begins.
	UnmarshalStart Event = iota

	// UnmarshalStop is the Event that is sent when deserialization of a payload
	// ends.
	UnmarshalStop

	// MarshalStart is the Event that is sent sent when serialization of a payload
	// begins.
	MarshalStart

	// MarshalStop is the Event that is sent sent when serialization of a payload
	// ends.
	MarshalStop
)

// Runtime has the same methods as jsonapi package for serialization and
// deserialization but also has a ctx, a map[string]interface{} for storing
// state, designed for instrumenting serialization timings.
type Runtime struct {
	ctx map[string]interface{}
}

// Events is the func type that provides the callback for handling event timings.
type Events func(*Runtime, Event, string, time.Duration)

// Instrumentation is a a global Events variable.  This is the handler for all
// timing events.
var Instrumentation Events

// NewRuntime creates a Runtime for use in an application.
func NewRuntime() *Runtime { return &Runtime{make(map[string]interface{})} }

// WithValue adds custom state variables to the runtime context.
func (r *Runtime) WithValue(key string, value interface{}) *Runtime {
	r.ctx[key] = value

	return r
}

// Value returns a state variable in the runtime context.
func (r *Runtime) Value(key string) interface{} {
	return r.ctx[key]
}

// Instrument is deprecated.
func (r *Runtime) Instrument(key string) *Runtime {
	return r.WithValue("instrument", key)
}

func (r *Runtime) shouldInstrument() bool {
	return Instrumentation != nil
}

// UnmarshalPayload has docs in request.go for UnmarshalPayload.
func (r *Runtime) UnmarshalPayload(reader io.Reader, model interface{}) error {
	return r.instrumentCall(UnmarshalStart, UnmarshalStop, func() error {
		return UnmarshalPayload(reader, model)
	})
}

// UnmarshalManyPayload has docs in request.go for UnmarshalManyPayload.
func (r *Runtime) UnmarshalManyPayload(reader io.Reader, kind reflect.Type) (elems []interface{}, err error) {
	r.instrumentCall(UnmarshalStart, UnmarshalStop, func() error {
		elems, err = UnmarshalManyPayload(reader, kind)
		return err
	})

	return
}

// MarshalPayload has docs in response.go for MarshalPayload.
func (r *Runtime) MarshalPayload(w io.Writer, model interface{}) error {
	return r.instrumentCall(MarshalStart, MarshalStop, func() error {
		return MarshalPayload(w, model)
	})
}

func (r *Runtime) instrumentCall(start Event, stop Event, c func() error) error {
	if !r.shouldInstrument() {
		return c()
	}

	instrumentationGUID, err := newUUID()
	if err != nil {
		return err
	}

	begin := time.Now()
	Instrumentation(r, start, instrumentationGUID, time.Duration(0))

	if err := c(); err != nil {
		return err
	}

	diff := time.Duration(time.Now().UnixNano() - begin.UnixNano())
	Instrumentation(r, stop, instrumentationGUID, diff)

	return nil
}

// citation: http://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, uuid); err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

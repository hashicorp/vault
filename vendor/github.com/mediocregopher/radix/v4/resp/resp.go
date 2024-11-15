// Package resp contains types and utilities useful for interacting with RESP
// protocols, without actually implementing any RESP protocol.
package resp

import (
	"bufio"
	"errors"
	"io"
)

// Opts are used to aid and affect marshaling and unmarshaling of RESP messages.
// Opts are not expected to be thread-safe.
//
// NewOpts should always be used to initialize a new Opts instance, even if some
// or all of the fields are expected to be changed. This way new fields may be
// added in the future without breaking existing usages.
type Opts struct {
	// GetBytes returns a *[]byte from an internal pool, or a newly allocated
	// instance if the pool is empty. The returned instance will have a length
	// of zero.
	//
	// This field may not be nil.
	GetBytes func() *[]byte

	// PutBytes puts a *[]byte back into the pool so it can be re-used later via
	// GetBytes.
	//
	// This field may not be nil.
	PutBytes func(*[]byte)

	// GetReader returns an io.Reader which will read out the given bytes.
	//
	// This field may not be nil.
	GetReader func([]byte) io.Reader

	// GetBufferedReader returns a BufferedReader which will read out the
	// contents of the given io.Reader.
	//
	// This field may not be nil.
	GetBufferedReader func(r io.Reader) BufferedReader

	// GetBufferedWriter returns a BufferedWriter which will buffer writes to
	// the given io.Writer, such that calling Flush on it will ensure all
	// previous writes have been written.
	//
	// This field may not be nil.
	GetBufferedWriter func(w io.Writer) BufferedWriter

	// Deterministic indicates that marshal operations should result in
	// deterministic results. This is largely used for ensuring map key/values
	// are emitted in a deterministic order.
	Deterministic bool

	// DisableErrorBubbling indicates that unmarshaled RESP errors should
	// not be treated as actual errors but like other response type.
	DisableErrorBubbling bool
}

const defaultBytePoolThreshold = 10000000 // ~10MB

// NewOpts returns an Opts instance which is suitable for most use-cases, and
// which may be modified if desired.
func NewOpts() *Opts {
	bp := newBytePool(defaultBytePoolThreshold)
	brp := newByteReaderPool()
	return &Opts{
		GetBytes:          bp.get,
		PutBytes:          bp.put,
		GetReader:         brp.get,
		GetBufferedReader: func(r io.Reader) BufferedReader { return bufio.NewReader(r) },
		GetBufferedWriter: func(w io.Writer) BufferedWriter { return bufio.NewWriter(w) },
	}
}

// Marshaler is the interface implemented by types that can marshal themselves
// into valid RESP messages. Opts may not be nil.
//
// NOTE It's important to keep track of whether a partial RESP message has been
// written to the Writer, and to use ErrConnUsable when returning errors if a
// partial RESP message has not been written.
type Marshaler interface {
	MarshalRESP(io.Writer, *Opts) error
}

// BufferedReader wraps a bufio.Reader.
type BufferedReader interface {
	io.Reader
	ReadSlice(delim byte) (line []byte, err error)
	Peek(n int) ([]byte, error)
	Discard(n int) (discarded int, err error)
	Buffered() int
}

// BufferedWriter wraps a bufio.Writer.
type BufferedWriter interface {
	io.Writer
	Flush() error
}

// Unmarshaler is the interface implemented by types that can unmarshal a RESP
// message of themselves. Opts may not be nil.
//
// NOTE It's important to keep track of whether a partial RESP message has been
// read off the BufferedReader, and to use ErrConnUsable when returning errors
// if a partial RESP message has not been read.
type Unmarshaler interface {
	UnmarshalRESP(BufferedReader, *Opts) error
}

// ErrConnUsable is used to wrap an error encountered while marshaling or
// unmarshaling a message on a connection. It indicates that the network
// connection is still healthy and that there are no partially written/read
// messages on the stream.
type ErrConnUsable struct {
	Err error
}

// ErrConnUnusable takes an existing error and, if it is wrapped in an
// ErrConnUsable, unwraps the ErrConnUsable from around it.
func ErrConnUnusable(err error) error {
	if err == nil {
		return nil
	} else if errConnUsable := (ErrConnUsable{}); errors.As(err, &errConnUsable) {
		return errConnUsable.Err
	}
	return err
}

func (ed ErrConnUsable) Error() string {
	return ed.Err.Error()
}

// Unwrap implements the errors.Wrapper interface.
func (ed ErrConnUsable) Unwrap() error {
	return ed.Err
}

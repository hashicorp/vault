package entropy

import (
	"fmt"
)

type Sourcer interface {
	GetRandom(bytes int) ([]byte, error)
}

type Reader struct {
	source Sourcer
}

func NewReader(source Sourcer) *Reader {
	return &Reader{source}
}

// Read reads exactly len(p) bytes from r into p.
// If r returns an error having read at least len(p) bytes, the error is dropped.
// It returns the number of bytes copied and an error if fewer bytes were read.
// On return, n == len(p) if and only if err == nil.
func (r *Reader) Read(p []byte) (n int, err error) {
	requested := len(p)
	randBytes, err := r.source.GetRandom(requested)
	delivered := copy(p, randBytes)
	if delivered != requested {
		if err != nil {
			return delivered, fmt.Errorf("unable to fill provided buffer with entropy: %w", err)
		}
		return delivered, fmt.Errorf("unable to fill provided buffer with entropy")
	}

	return delivered, nil
}

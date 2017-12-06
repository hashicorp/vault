package lib

import (
	"io"
	"strings"

	"github.com/hashicorp/yamux"
)

var yamuxStreamClosed = yamux.ErrStreamClosed.Error()
var yamuxSessionShutdown = yamux.ErrSessionShutdown.Error()

// IsErrEOF returns true if we get an EOF error from the socket itself, or
// an EOF equivalent error from yamux.
func IsErrEOF(err error) bool {
	if err == io.EOF {
		return true
	}

	errStr := err.Error()
	if strings.Contains(errStr, yamuxStreamClosed) ||
		strings.Contains(errStr, yamuxSessionShutdown) {
		return true
	}

	return false
}

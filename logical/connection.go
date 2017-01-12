package logical

import (
	"crypto/tls"
)

// Connection represents the connection information for a request. This
// is present on the Request structure for credential backends.
type Connection struct {
	// RemoteAddr is the network address that sent the request.
	RemoteAddr string

	// ConnState is the TLS connection state if applicable.
	ConnState *tls.ConnectionState
}

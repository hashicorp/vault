package logical

import (
	"crypto/tls"
)

// Connection represents the connection information for a request. This
// is present on the Request structure for credential backends.
type Connection struct {
	// RemoteAddr is the network address that sent the request.
	RemoteAddr string `json:"remote_addr"`

	// RemotePort is the network port that sent the request.
	RemotePort int `json:"remote_port"`

	// ConnState is the TLS connection state if applicable.
	ConnState *tls.ConnectionState `sentinel:""`
}

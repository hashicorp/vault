package listenerutil

import (
	"context"
	"net"

	"google.golang.org/grpc/test/bufconn"
)

// BufConnListenerDialer implements consul-template's TransportDialer using a
// bufconn listener, to provide a way to Dial the in-memory listener
type BufConnListenerDialer struct {
	listener *bufconn.Listener
}

// NewBufConnListenerDialer returns a new BufConnListenerDialer using an
// existing bufconn.Listener
func NewBufConnListenerDialer(bcl *bufconn.Listener) *BufConnListenerDialer {
	return &BufConnListenerDialer{
		listener: bcl,
	}
}

// Dial connects to the listening end of the bufconn (satisfies
// consul-template's TransportDialer interface). This is essentially the client
// side of the bufconn connection.
func (bcl *BufConnListenerDialer) Dial(network, addr string) (net.Conn, error) {
	return bcl.listener.Dial()
}

// DialContext connects to the listening end of the bufconn (satisfies
// consul-template's TransportDialer interface). This is essentially the client
// side of the bufconn connection.
func (bcl *BufConnListenerDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return bcl.listener.DialContext(ctx)
}

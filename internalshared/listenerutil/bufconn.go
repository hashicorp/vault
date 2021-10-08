package listenerutil

import (
	"net"

	ctconfig "github.com/hashicorp/consul-template/config"
	"google.golang.org/grpc/test/bufconn"
)

var (
	_ net.Listener             = (*BufConnListenerDialer)(nil)
	_ ctconfig.TransportDialer = (*BufConnListenerDialer)(nil)
)

// BufConnListenerDialer implements both a net.Listener and consul
// TransportDialer, to serve both ends of an in-process connection (Dial and
// Accept).
type BufConnListenerDialer struct {
	listener *bufconn.Listener
}

// NewBufConnListenerDialer returns a new BufConnListenerDialer
func NewBufConnListenerDialer() *BufConnListenerDialer {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	return &BufConnListenerDialer{listener: listener}
}

// Accept incoming connections to the bufconn listener
func (bcl *BufConnListenerDialer) Accept() (net.Conn, error) {
	return bcl.listener.Accept()
}

// Close the bufconn
func (bcl *BufConnListenerDialer) Close() error {
	return bcl.listener.Close()
}

// Addr returns the net.Addr of the bufconn listener
func (bcl *BufConnListenerDialer) Addr() net.Addr {
	return bcl.listener.Addr()
}

// Dial and connect to the listening end of the bufconn (satisfies
// consul-template's TransportDialer interface). This is essentially the client
// side of the bufconn connection.
func (bcl *BufConnListenerDialer) Dial(network, addr string) (net.Conn, error) {
	return bcl.listener.Dial()
}

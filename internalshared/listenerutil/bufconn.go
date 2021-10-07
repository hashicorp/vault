package listenerutil

import (
	"net"

	"google.golang.org/grpc/test/bufconn"
)

type BufConnListener struct {
	listener *bufconn.Listener
}

func NewBufConnListener() *BufConnListener {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	return &BufConnListener{listener: listener}
}

func (bcl *BufConnListener) Accept() (net.Conn, error) {
	return bcl.listener.Accept()
}

func (bcl *BufConnListener) Close() error {
	return bcl.listener.Close()
}

func (bcl *BufConnListener) Dial(network, addr string) (net.Conn, error) {
	return bcl.listener.Dial()
}

func (bcl *BufConnListener) Addr() net.Addr {
	return bcl.listener.Addr()
}

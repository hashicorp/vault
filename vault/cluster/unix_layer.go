package cluster

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/atomic"
)

type UnixLayer struct {
	listeners []NetworkListener
	addrs     []*net.UnixAddr
	logger    hclog.Logger

	l       sync.Mutex
	stopped *atomic.Bool
}

var _ NetworkLayer = &UnixLayer{}

// NewUnixLayer returns a UNIXLayer.
func NewUnixLayer(addrs []*net.UnixAddr, logger hclog.Logger) *UnixLayer {
	return &UnixLayer{
		addrs:   addrs,
		logger:  logger,
		stopped: atomic.NewBool(false),
	}
}

func (l *UnixLayer) Addrs() []net.Addr {
	l.l.Lock()
	defer l.l.Unlock()

	if len(l.addrs) == 0 {
		return nil
	}

	ret := make([]net.Addr, len(l.addrs))
	for i, a := range l.addrs {
		ret[i] = a
	}

	return ret
}

func (l *UnixLayer) Listeners() []NetworkListener {
	l.l.Lock()
	defer l.l.Unlock()

	if l.listeners != nil {
		return l.listeners
	}

	listeners := []NetworkListener{}
	for i, laddr := range l.addrs {
		l.logger.Info("starting listener", "listener_address", laddr)

		unixLn, err := net.ListenUnix("unix", laddr)
		if err != nil {
			l.logger.Error("error starting listener", "error", err)
			continue
		}
		if laddr.String() != unixLn.Addr().String() {
			// If we listened on port 0, record the port the OS gave us.
			l.addrs[i] = unixLn.Addr().(*net.UnixAddr)
		}

		listeners = append(listeners, unixLn)
	}

	l.listeners = listeners

	return listeners
}

func (l *UnixLayer) Dial(address string, timeout time.Duration, tlsConfig *tls.Config) (*tls.Conn, error) {
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	return tls.DialWithDialer(dialer, "unix", address, tlsConfig)
}

func (l *UnixLayer) Close() error {
	if l.stopped.Swap(true) {
		return nil
	}
	l.l.Lock()
	defer l.l.Unlock()

	var retErr *multierror.Error
	for _, ln := range l.listeners {
		if err := ln.Close(); err != nil {
			retErr = multierror.Append(retErr, err)
		}
	}

	l.listeners = nil

	return retErr.ErrorOrNil()
}

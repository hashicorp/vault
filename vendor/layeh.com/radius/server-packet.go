package radius

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

type packetResponseWriter struct {
	// listener that received the packet
	conn net.PacketConn
	addr net.Addr
}

func (r *packetResponseWriter) Write(packet *Packet) error {
	raw, err := packet.Encode()
	if err != nil {
		return err
	}
	if _, err := r.conn.WriteTo(raw, r.addr); err != nil {
		return err
	}
	return nil
}

// PacketServer listens for RADIUS requests on a packet-based protocols (e.g.
// UDP).
type PacketServer struct {
	// The address on which the server listens. Defaults to :1812.
	Addr string
	// The network on which the server listens. Defaults to udp.
	Network      string
	SecretSource SecretSource
	Handler      Handler

	// Skip incoming packet authenticity validation.
	// This should only be set to true for debugging purposes.
	InsecureSkipVerify bool

	mu           sync.Mutex
	shuttingDown bool
	ctx          context.Context
	ctxDone      context.CancelFunc
	running      chan struct{}
	listeners    map[net.PacketConn]int
	activeCount  int32
}

// TODO: logger on PacketServer

// Serve accepts incoming connections on conn.
func (s *PacketServer) Serve(conn net.PacketConn) error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return errors.New("radius: nil SecretSource")
	}

	s.mu.Lock()
	if s.shuttingDown {
		s.mu.Unlock()
		return ErrServerShutdown
	}
	var ctx context.Context
	if s.ctx == nil {
		s.ctx, s.ctxDone = context.WithCancel(context.Background())
		ctx = s.ctx
	}
	if s.running == nil {
		s.running = make(chan struct{})
	}
	if s.listeners == nil {
		s.listeners = make(map[net.PacketConn]int)
	}
	s.listeners[conn]++
	s.mu.Unlock()

	type activeKey struct {
		IP         string
		Identifier byte
	}

	var (
		activeLock sync.Mutex
		active     = map[activeKey]struct{}{}
	)

	atomic.AddInt32(&s.activeCount, 1)
	defer func() {
		s.mu.Lock()
		s.listeners[conn]--
		if s.listeners[conn] == 0 {
			delete(s.listeners, conn)
		}
		s.mu.Unlock()

		if atomic.AddInt32(&s.activeCount, -1) == 0 {
			s.mu.Lock()
			s.shuttingDown = false
			close(s.running)
			s.running = nil
			s.ctx = nil
			s.mu.Unlock()
		}
	}()

	for {
		var buff [MaxPacketLength]byte
		n, remoteAddr, err := conn.ReadFrom(buff[:])
		if err != nil {
			s.mu.Lock()
			if s.shuttingDown {
				s.mu.Unlock()
				return nil
			}
			s.mu.Unlock()

			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				return err
			}
			// TODO: log error?
			continue
		}

		buffCopy := make([]byte, n)
		copy(buffCopy, buff[:n])

		atomic.AddInt32(&s.activeCount, 1)
		go func(buff []byte, remoteAddr net.Addr) {
			secret, err := s.SecretSource.RADIUSSecret(ctx, remoteAddr)
			if err != nil {
				// TODO: log only if server is not shutting down?
				return
			}
			if len(secret) == 0 {
				return
			}

			if !s.InsecureSkipVerify && !IsAuthenticRequest(buff, secret) {
				// TODO: log?
				return
			}

			packet, err := Parse(buff, secret)
			if err != nil {
				// TODO: error logger
				return
			}

			key := activeKey{
				IP:         remoteAddr.String(),
				Identifier: packet.Identifier,
			}
			activeLock.Lock()
			if _, ok := active[key]; ok {
				activeLock.Unlock()
				return
			}
			active[key] = struct{}{}
			activeLock.Unlock()

			response := packetResponseWriter{
				conn: conn,
				addr: remoteAddr,
			}

			defer func() {
				activeLock.Lock()
				delete(active, key)
				activeLock.Unlock()

				if atomic.AddInt32(&s.activeCount, -1) == 0 {
					s.mu.Lock()
					s.shuttingDown = false
					close(s.running)
					s.running = nil
					s.ctx = nil
					s.mu.Unlock()
				}
			}()

			request := Request{
				LocalAddr:  conn.LocalAddr(),
				RemoteAddr: remoteAddr,
				Packet:     packet,
				ctx:        ctx,
			}

			s.Handler.ServeRADIUS(&response, &request)
		}(buffCopy, remoteAddr)
	}
}

// ListenAndServe starts a RADIUS server on the address given in s.
func (s *PacketServer) ListenAndServe() error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return errors.New("radius: nil SecretSource")
	}

	addrStr := ":1812"
	if s.Addr != "" {
		addrStr = s.Addr
	}

	network := "udp"
	if s.Network != "" {
		network = s.Network
	}

	pc, err := net.ListenPacket(network, addrStr)
	if err != nil {
		return err
	}
	defer pc.Close()
	return s.Serve(pc)
}

// Shutdown gracefully stops the server. It first closes all listeners (which
// stops accepting new packets) and then waits for running handlers to complete.
//
// Shutdown returns after all handlers have completed, or when ctx is canceled.
// The PacketServer is ready for re-use once the function returns nil.
func (s *PacketServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()

	if len(s.listeners) == 0 {
		s.mu.Unlock()
		return nil
	}

	if !s.shuttingDown {
		s.shuttingDown = true
		s.ctxDone()
		for listener := range s.listeners {
			listener.Close()
		}
	}

	ch := s.running
	s.mu.Unlock()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

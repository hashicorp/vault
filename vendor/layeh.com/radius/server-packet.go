package radius

import (
	"context"
	"errors"
	"log"
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
	encoded, err := packet.Encode()
	if err != nil {
		return err
	}
	if _, err := r.conn.WriteTo(encoded, r.addr); err != nil {
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
	Network string

	// The source from which the secret is obtained for parsing and validating
	// the request.
	SecretSource SecretSource

	// Handler which is called to process the request.
	Handler Handler

	// Skip incoming packet authenticity validation.
	// This should only be set to true for debugging purposes.
	InsecureSkipVerify bool

	// ErrorLog specifies an optional logger for errors
	// around packet accepting, processing, and validation.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	shutdownRequested int32

	mu          sync.Mutex
	ctx         context.Context
	ctxDone     context.CancelFunc
	listeners   map[net.PacketConn]uint
	lastActive  chan struct{} // closed when the last active item finishes
	activeCount int32
}

func (s *PacketServer) initLocked() {
	if s.ctx == nil {
		s.ctx, s.ctxDone = context.WithCancel(context.Background())
		s.listeners = make(map[net.PacketConn]uint)
		s.lastActive = make(chan struct{})
	}
}

func (s *PacketServer) activeAdd() {
	atomic.AddInt32(&s.activeCount, 1)
}

func (s *PacketServer) activeDone() {
	if atomic.AddInt32(&s.activeCount, -1) == -1 {
		close(s.lastActive)
	}
}

func (s *PacketServer) logf(format string, args ...interface{}) {
	if s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Serve accepts incoming connections on conn.
func (s *PacketServer) Serve(conn net.PacketConn) error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return errors.New("radius: nil SecretSource")
	}

	s.mu.Lock()
	s.initLocked()
	if atomic.LoadInt32(&s.shutdownRequested) == 1 {
		s.mu.Unlock()
		return ErrServerShutdown
	}

	s.listeners[conn]++
	s.mu.Unlock()

	type requestKey struct {
		IP         string
		Identifier byte
	}

	var (
		requestsLock sync.Mutex
		requests     = map[requestKey]struct{}{}
	)

	s.activeAdd()
	defer func() {
		s.mu.Lock()
		s.listeners[conn]--
		if s.listeners[conn] == 0 {
			delete(s.listeners, conn)
		}
		s.mu.Unlock()
		s.activeDone()
	}()

	var buff [MaxPacketLength]byte
	for {
		n, remoteAddr, err := conn.ReadFrom(buff[:])
		if err != nil {
			if atomic.LoadInt32(&s.shutdownRequested) == 1 {
				return ErrServerShutdown
			}

			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				return err
			}
			s.logf("radius: could not read packet: %v", err)
			continue
		}

		s.activeAdd()
		go func(buff []byte, remoteAddr net.Addr) {
			defer s.activeDone()

			secret, err := s.SecretSource.RADIUSSecret(s.ctx, remoteAddr)
			if err != nil {
				s.logf("radius: error fetching from secret source: %v", err)
				return
			}
			if len(secret) == 0 {
				s.logf("radius: empty secret returned from secret source")
				return
			}

			if !s.InsecureSkipVerify && !IsAuthenticRequest(buff, secret) {
				s.logf("radius: packet validation failed; bad secret")
				return
			}

			packet, err := Parse(buff, secret)
			if err != nil {
				s.logf("radius: unable to parse packet: %v", err)
				return
			}

			key := requestKey{
				IP:         remoteAddr.String(),
				Identifier: packet.Identifier,
			}

			requestsLock.Lock()
			if _, ok := requests[key]; ok {
				requestsLock.Unlock()
				return
			}
			requests[key] = struct{}{}
			requestsLock.Unlock()

			response := packetResponseWriter{
				conn: conn,
				addr: remoteAddr,
			}

			defer func() {
				requestsLock.Lock()
				delete(requests, key)
				requestsLock.Unlock()
			}()

			request := Request{
				LocalAddr:  conn.LocalAddr(),
				RemoteAddr: remoteAddr,
				Packet:     packet,
				ctx:        s.ctx,
			}

			s.Handler.ServeRADIUS(&response, &request)
		}(append([]byte(nil), buff[:n]...), remoteAddr)
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

// Shutdown gracefully stops the server. It first closes all listeners and then
// waits for any running handlers to complete.
//
// Shutdown returns after nil all handlers have completed. ctx.Err() is
// returned if ctx is canceled.
//
// Any Serve methods return ErrShutdown after Shutdown is called.
func (s *PacketServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.initLocked()
	if atomic.CompareAndSwapInt32(&s.shutdownRequested, 0, 1) {
		for listener := range s.listeners {
			listener.Close()
		}

		s.ctxDone()
		s.activeDone()
	}
	s.mu.Unlock()

	select {
	case <-s.lastActive:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

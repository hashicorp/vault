package radius // import "layeh.com/radius"

import (
	"errors"
	"net"
	"sync"
)

// Handler is a value that can handle a server's RADIUS packet event.
type Handler interface {
	ServeRADIUS(w ResponseWriter, p *Packet)
}

// HandlerFunc is a wrapper that allows ordinary functions to be used as a
// handler.
type HandlerFunc func(w ResponseWriter, p *Packet)

// ServeRADIUS calls h(w, p).
func (h HandlerFunc) ServeRADIUS(w ResponseWriter, p *Packet) {
	h(w, p)
}

var _ Handler = HandlerFunc(nil)

// ResponseWriter is used by Handler when replying to a RADIUS packet.
type ResponseWriter interface {
	// LocalAddr returns the address of the local server that accepted the
	// packet.
	LocalAddr() net.Addr
	// RemoteAddr returns the address of the remote client that sent to packet.
	RemoteAddr() net.Addr

	// Write sends a packet to the sender.
	Write(packet *Packet) error

	// AccessAccept sends an Access-Accept packet to the sender that includes
	// the given attributes.
	AccessAccept(attributes ...*Attribute) error
	// AccessAccept sends an Access-Reject packet to the sender that includes
	// the given attributes.
	AccessReject(attributes ...*Attribute) error
	// AccessAccept sends an Access-Challenge packet to the sender that includes
	// the given attributes.
	AccessChallenge(attributes ...*Attribute) error
}

type responseWriter struct {
	// listener that received the packet
	conn net.PacketConn
	// where the packet came from
	addr net.Addr
	// original packet
	packet *Packet
}

func (r *responseWriter) LocalAddr() net.Addr {
	return r.conn.LocalAddr()
}

func (r *responseWriter) RemoteAddr() net.Addr {
	return r.addr
}

func (r *responseWriter) accessRespond(code Code, attributes ...*Attribute) error {
	packet := Packet{
		Code:          code,
		Identifier:    r.packet.Identifier,
		Authenticator: r.packet.Authenticator,

		Secret: r.packet.Secret,

		Dictionary: r.packet.Dictionary,

		Attributes: attributes,
	}
	return r.Write(&packet)
}

func (r *responseWriter) AccessAccept(attributes ...*Attribute) error {
	// TOOD: do not send if packet was not Access-Request
	return r.accessRespond(CodeAccessAccept, attributes...)
}

func (r *responseWriter) AccessReject(attributes ...*Attribute) error {
	// TOOD: do not send if packet was not Access-Request
	return r.accessRespond(CodeAccessReject, attributes...)
}

func (r *responseWriter) AccessChallenge(attributes ...*Attribute) error {
	// TOOD: do not send if packet was not Access-Request
	return r.accessRespond(CodeAccessChallenge, attributes...)
}

func (r *responseWriter) Write(packet *Packet) error {
	raw, err := packet.Encode()
	if err != nil {
		return err
	}
	if _, err := r.conn.WriteTo(raw, r.addr); err != nil {
		return err
	}
	return nil
}

// Server is a server that listens for and handles RADIUS packets.
type Server struct {
	// Address to bind the server on. If empty, the address defaults to ":1812".
	Addr string
	// Network of the server. Valid values are "udp", "udp4", "udp6". If empty,
	// the network defaults to "udp".
	Network string
	// The shared secret between the client and server.
	Secret []byte

	// TODO: allow a secret function to be defined, which returned the secret
	// that should be used for the given client.

	// Dictionary used when decoding incoming packets.
	Dictionary *Dictionary

	// The packet handler that handles incoming, valid packets.
	Handler Handler
}

// Serve accepts incoming connections on the net.PacketConn pc.
func (s *Server) Serve(pc net.PacketConn) error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}

	type activeKey struct {
		IP         string
		Identifier byte
	}

	var (
		activeLock sync.Mutex
		active     = map[activeKey]struct{}{}
	)

	for {
		buff := make([]byte, 4096)
		n, remoteAddr, err := pc.ReadFrom(buff)
		if err != nil {
			if err.(*net.OpError).Temporary() {
				return err
			}
			continue
		}

		packet, err := Parse(buff[:n], s.Secret, s.Dictionary)
		if err != nil {
			continue
		}

		go func(packet *Packet, remoteAddr net.Addr) {
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

			response := responseWriter{
				conn:   pc,
				addr:   remoteAddr,
				packet: packet,
			}

			s.Handler.ServeRADIUS(&response, packet)

			activeLock.Lock()
			delete(active, key)
			activeLock.Unlock()
		}(packet, remoteAddr)
	}
}

// ListenAndServe starts a RADIUS server on the address given in s.
func (s *Server) ListenAndServe() error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
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

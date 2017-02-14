package radius // import "layeh.com/radius"

import (
	"net"
	"time"
)

// Client is a RADIUS client that can send and receive packets to and from a
// RADIUS server.
type Client struct {
	// Network on which to make the connection. Defaults to "udp".
	Net string

	// Local address to use for outgoing connections (can be nil).
	LocalAddr net.Addr

	// Timeouts for various operations. Default values for each field is 10
	// seconds.
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Interval on which to resend packet.
	Retry time.Duration
}

// Exchange sends the packet to the given server address and waits for a
// response. nil and an error is returned upon failure.
func (c *Client) Exchange(packet *Packet, addr string) (*Packet, error) {
	wire, err := packet.Encode()
	if err != nil {
		return nil, err
	}

	connNet := c.Net
	if connNet == "" {
		connNet = "udp"
	}

	const defaultTimeout = 10 * time.Second
	dialTimeout := c.DialTimeout
	if dialTimeout == 0 {
		dialTimeout = defaultTimeout
	}

	dialer := net.Dialer{
		Timeout:   dialTimeout,
		LocalAddr: c.LocalAddr,
	}
	conn, err := dialer.Dial(connNet, addr)
	if err != nil {
		return nil, err
	}

	writeTimeout := c.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = defaultTimeout
	}
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	conn.Write(wire)

	if c.Retry > 0 {
		retry := time.NewTicker(c.Retry)
		end := make(chan struct{})
		defer close(end)
		go func() {
			for {
				select {
				case <-retry.C:
					conn.Write(wire)
				case <-end:
					retry.Stop()
					return
				}
			}
		}()
	}

	var incoming [maxPacketSize]byte

	readTimeout := c.ReadTimeout
	if readTimeout == 0 {
		readTimeout = defaultTimeout
	}
	conn.SetReadDeadline(time.Now().Add(readTimeout))

	for {
		n, err := conn.Read(incoming[:])
		if err != nil {
			conn.Close()
			return nil, err
		}
		received, err := Parse(incoming[:n], packet.Secret, packet.Dictionary)
		if err == nil && received.IsAuthentic(packet) {
			conn.Close()
			return received, nil
		}
	}
}

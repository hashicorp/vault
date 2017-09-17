package radius

import (
	"context"
	"net"
	"time"
)

// Client is a RADIUS client that can exchange packets with a RADIUS server.
type Client struct {
	// Network on which to make the connection. Defaults to "udp".
	Net string

	// Dialer to use when making the outgoing connections.
	Dialer net.Dialer

	// Interval on which to resend packet (zero or negative value means no
	// retry).
	Retry time.Duration

	InsecureSkipVerify bool
}

// DefaultClient is the RADIUS client used by the Exchange function.
var DefaultClient = &Client{}

// Exchange uses DefaultClient to send the given RADIUS packet to the server at
// address addr and waits for a response.
func Exchange(ctx context.Context, packet *Packet, addr string) (*Packet, error) {
	return DefaultClient.Exchange(ctx, packet, addr)
}

// Exchange sends the packet to the given server and waits for a response. ctx
// must be non-nil.
func (c *Client) Exchange(ctx context.Context, packet *Packet, addr string) (*Packet, error) {
	if ctx == nil {
		panic("nil context")
	}

	wire, err := packet.Encode()
	if err != nil {
		return nil, err
	}

	connNet := c.Net
	if connNet == "" {
		connNet = "udp"
	}

	conn, err := c.Dialer.DialContext(ctx, connNet, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if deadline, deadlineSet := ctx.Deadline(); deadlineSet {
		conn.SetDeadline(deadline)
	}

	conn.Write(wire)

	if c.Retry > 0 {
		retry := time.NewTicker(c.Retry)
		defer retry.Stop()
		end := make(chan struct{})
		defer close(end)
		go func() {
			for {
				select {
				case <-retry.C:
					conn.Write(wire)
				case <-ctx.Done():
					return
				case <-end:
					return
				}
			}
		}()
	}

	var incoming [MaxPacketLength]byte
	for {
		n, err := conn.Read(incoming[:])
		if err != nil {
			return nil, err
		}
		received, err := Parse(incoming[:n], packet.Secret)
		if err == nil && (!c.InsecureSkipVerify && IsAuthenticResponse(incoming[:n], wire, packet.Secret)) {
			return received, nil
		}
	}
}

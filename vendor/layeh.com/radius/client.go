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

	// MaxPacketErrors controls how many packet parsing and validation errors
	// the client will ignore before returning the error from Exchange.
	//
	// If zero, Exchange will drop all packet parsing errors.
	MaxPacketErrors int

	// InsecureSkipVerify controls whether the client should skip verifying
	// response packets received.
	InsecureSkipVerify bool
}

// DefaultClient is the RADIUS client used by the Exchange function.
var DefaultClient = &Client{
	Retry:           time.Second,
	MaxPacketErrors: 10,
}

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
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer conn.Close()

	conn.Write(wire)

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	var retryTimer <-chan time.Time
	if c.Retry > 0 {
		retry := time.NewTicker(c.Retry)
		defer retry.Stop()
		retryTimer = retry.C
	}

	go func() {
		defer conn.Close()
		for {
			select {
			case <-retryTimer:
				conn.Write(wire)
			case <-ctx.Done():
				return
			}
		}
	}()

	var packetErrorCount int

	var incoming [MaxPacketLength]byte
	for {
		n, err := conn.Read(incoming[:])
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			return nil, err
		}

		received, err := Parse(incoming[:n], packet.Secret)
		if err != nil {
			packetErrorCount++
			if c.MaxPacketErrors > 0 && packetErrorCount >= c.MaxPacketErrors {
				return nil, err
			}
			continue
		}

		if !c.InsecureSkipVerify && !IsAuthenticResponse(incoming[:n], wire, packet.Secret) {
			packetErrorCount++
			if c.MaxPacketErrors > 0 && packetErrorCount >= c.MaxPacketErrors {
				return nil, &NonAuthenticResponseError{}
			}
			continue
		}

		return received, nil
	}
}

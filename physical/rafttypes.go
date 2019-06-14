package physical

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/hashicorp/vault/vault/cluster"
)

type ClusterHook interface {
	AddClient(alpn string, client cluster.Client)
	RemoveClient(alpn string)
	AddHandler(alpn string, handler cluster.Handler)
	StopHandler(alpn string)
	TLSConfig(ctx context.Context) (*tls.Config, error)
	Addr() net.Addr
}

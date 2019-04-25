package physical

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/vault/cluster"
)

type NetworkConfig struct {
	Addr      net.Addr
	Cert      []byte
	KeyParams *certutil.ClusterKeyParams
}

type ClusterHook interface {
	AddClient(alpn string, client cluster.Client)
	RemoveClient(alpn string)
	AddHandler(alpn string, handler cluster.Handler)
	StopHandler(alpn string)
	TLSConfig(ctx context.Context) (*tls.Config, error)
	Addr() net.Addr
}

type Clustered interface {
	SetupCluster(context.Context, *NetworkConfig, ClusterHook) error
	TeardownCluster(ClusterHook) error
}

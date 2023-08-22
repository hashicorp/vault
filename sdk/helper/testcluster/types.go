package testcluster

import (
	"crypto/tls"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

type VaultClusterNode interface {
	APIClient() *api.Client
	TLSConfig() *tls.Config
}

type VaultCluster interface {
	Nodes() []VaultClusterNode
	GetBarrierKeys() [][]byte
	GetRecoveryKeys() [][]byte
	GetBarrierOrRecoveryKeys() [][]byte
	SetBarrierKeys([][]byte)
	SetRecoveryKeys([][]byte)
	GetCACertPEMFile() string
	Cleanup()
	ClusterID() string
	NamedLogger(string) hclog.Logger
	SetRootToken(token string)
	GetRootToken() string
}

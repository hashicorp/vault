// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"time"

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

type VaultNodeConfig struct {
	// Not configurable because cluster creator wants to control these:
	//   PluginDirectory string `hcl:"plugin_directory"`
	//   APIAddr              string      `hcl:"api_addr"`
	//   ClusterAddr          string      `hcl:"cluster_addr"`
	//   Storage   *Storage `hcl:"-"`
	//   HAStorage *Storage `hcl:"-"`
	//   DisableMlock bool `hcl:"disable_mlock"`
	//   ClusterName string `hcl:"cluster_name"`

	// Not configurable yet:
	//   Listeners []*Listener `hcl:"-"`
	//   Seals   []*KMS   `hcl:"-"`
	//   Entropy *Entropy `hcl:"-"`
	//   Telemetry *Telemetry `hcl:"telemetry"`
	//   HCPLinkConf *HCPLinkConfig `hcl:"cloud"`
	//   PidFile string `hcl:"pid_file"`
	//   ServiceRegistrationType        string
	//   ServiceRegistrationOptions    map[string]string

	StorageOptions      map[string]string
	AdditionalListeners []VaultNodeListenerConfig

	DefaultMaxRequestDuration      time.Duration `json:"default_max_request_duration"`
	LogFormat                      string        `json:"log_format"`
	LogLevel                       string        `json:"log_level"`
	CacheSize                      int           `json:"cache_size"`
	DisableCache                   bool          `json:"disable_cache"`
	DisablePrintableCheck          bool          `json:"disable_printable_check"`
	EnableUI                       bool          `json:"ui"`
	MaxLeaseTTL                    time.Duration `json:"max_lease_ttl"`
	DefaultLeaseTTL                time.Duration `json:"default_lease_ttl"`
	ClusterCipherSuites            string        `json:"cluster_cipher_suites"`
	PluginFileUid                  int           `json:"plugin_file_uid"`
	PluginFilePermissions          int           `json:"plugin_file_permissions"`
	EnableRawEndpoint              bool          `json:"raw_storage_endpoint"`
	DisableClustering              bool          `json:"disable_clustering"`
	DisablePerformanceStandby      bool          `json:"disable_performance_standby"`
	DisableSealWrap                bool          `json:"disable_sealwrap"`
	DisableIndexing                bool          `json:"disable_indexing"`
	DisableSentinelTrace           bool          `json:"disable_sentinel"`
	EnableResponseHeaderHostname   bool          `json:"enable_response_header_hostname"`
	LogRequestsLevel               string        `json:"log_requests_level"`
	EnableResponseHeaderRaftNodeID bool          `json:"enable_response_header_raft_node_id"`
	LicensePath                    string        `json:"license_path"`
}

type ClusterNode struct {
	APIAddress string `json:"api_address"`
}

type ClusterJson struct {
	Nodes      []ClusterNode `json:"nodes"`
	CACertPath string        `json:"ca_cert_path"`
	RootToken  string        `json:"root_token"`
}

type ClusterOptions struct {
	ClusterName                 string
	KeepStandbysSealed          bool
	SkipInit                    bool
	CACert                      []byte
	NumCores                    int
	TmpDir                      string
	Logger                      hclog.Logger
	VaultNodeConfig             *VaultNodeConfig
	VaultLicense                string
	AdministrativeNamespacePath string
}

type VaultNodeListenerConfig struct {
	Port              int
	ChrootNamespace   string
	RedactAddresses   bool
	RedactClusterName bool
	RedactVersion     bool
}

type CA struct {
	CACert        *x509.Certificate
	CACertBytes   []byte
	CACertPEM     []byte
	CACertPEMFile string
	CAKey         *ecdsa.PrivateKey
	CAKeyPEM      []byte
}

type ClusterStorage interface {
	Start(context.Context, *ClusterOptions) error
	Cleanup() error
	Opts() map[string]interface{}
	Type() string
}

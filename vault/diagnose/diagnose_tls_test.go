package diagnose

import (
	"testing"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func setup(t *testing.T) *vault.Core {
	serverConf := &server.Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                          "tcp",
					Address:                       "127.0.0.1:443",
					ClusterAddress:                "127.0.0.1:8201",
					TLSCertFile:                   "./certs/server.crt",
					TLSKeyFile:                    "./certs/server.key",
					TLSClientCAFile:               "./certs/rootca.crt",
					TLSMinVersion:                 "tls11",
					TLSRequireAndVerifyClientCert: true,
					TLSDisableClientCerts:         true,
				},
				{
					Type:                          "tcp",
					Address:                       "127.0.0.1:443",
					ClusterAddress:                "127.0.0.1:8201",
					TLSCertFile:                   "./certs/server2.crt",
					TLSKeyFile:                    "./certs/server2.key",
					TLSClientCAFile:               "./certs/rootca2.crt",
					TLSMinVersion:                 "tls12",
					TLSRequireAndVerifyClientCert: true,
					TLSDisableClientCerts:         false,
				},
				{
					Type:                          "tcp",
					Address:                       "127.0.0.1:443",
					ClusterAddress:                "127.0.0.1:8201",
					TLSCertFile:                   "./certs/server3.crt",
					TLSKeyFile:                    "./certs/server3.key",
					TLSClientCAFile:               "./certs/rootca3.crt",
					TLSMinVersion:                 "tls13",
					TLSRequireAndVerifyClientCert: false,
					TLSDisableClientCerts:         true,
				},
			},
		},
	}
	conf := &vault.CoreConfig{
		RawConfig: serverConf,
	}
	core := vault.TestCoreWithConfig(t, conf)
	return core
}

func TestTLSConfigChecks(t *testing.T) {}

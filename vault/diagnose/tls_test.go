package diagnose

import (
	"testing"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
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

func TestTLSValidCert(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "0",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTLSFakeCert(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/fakecert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "0",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Errorf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != "tls: failed to find any PEM data in certificate input" {
		t.Errorf("Bad error message: %s", err.Error())
	}
}

// TestTLSTrailingData uses a certificate from:
// https://github.com/golang/go/issues/40545 that contains
// an extra DER sequence, and makes sure a trailing data error
// is returned.
func TestTLSTrailingData(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/trailingdatacert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "0",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Errorf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != "asn1: syntax error: trailing data" {
		t.Errorf("Bad error message: %s", err.Error())
	}
}

func TestTLSExpiredCert(t *testing.T) {
}

func TestTLSMismatchedCryptographicInfo(t *testing.T) {}

func TestTLSContradictoryFlags(t *testing.T) {}

func TestTLSBadCipherSuite(t *testing.T) {}

func TestTLSUnknownAlgorithm(t *testing.T) {}

func TestTLSIncorrectUsageType(t *testing.T) {}

package config

import (
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

func TestLoadConfigFile_AgentCache(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfig("./test-fixtures/config-cache.hcl", logger)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				WrapTTL:   300 * time.Second,
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken: true,
		},
		Listeners: []*Listener{
			&Listener{
				Type: "unix",
				Config: map[string]interface{}{
					"address":      "/path/to/socket",
					"tls_disable":  true,
					"socket_mode":  "configmode",
					"socket_user":  "configuser",
					"socket_group": "configgroup",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":       "127.0.0.1:8400",
					"tls_key_file":  "/path/to/cakey.pem",
					"tls_cert_file": "/path/to/cacert.pem",
				},
			},
		},
		Vault: &Vault{
			Address:       "http://127.0.0.1:1111",
			CACert:        "config_ca_cert",
			CAPath:        "config_ca_path",
			TLSSkipVerify: true,
			ClientCert:    "config_client_cert",
			ClientKey:     "config_client_key",
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfig("./test-fixtures/config-cache-embedded-type.hcl", logger)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	os.Setenv("TEST_AAD_ENV", "aad")
	defer os.Unsetenv("TEST_AAD_ENV")

	config, err := LoadConfig("./test-fixtures/config.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				WrapTTL:   300 * time.Second,
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
				&Sink{
					Type:    "file",
					WrapTTL: 5 * time.Minute,
					DHType:  "curve25519",
					DHPath:  "/tmp/file-foo-dhpath2",
					AAD:     "aad",
					Config: map[string]interface{}{
						"path": "/tmp/file-bar",
					},
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfig("./test-fixtures/config-embedded-type.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_NoAutoAuth(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfig("./test-fixtures/config-cache-no-auto_auth.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Cache: &Cache{},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Bad_AgentCache_InconsisentAutoAuth(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	_, err := LoadConfig("./test-fixtures/bad-config-cache-inconsistent-auto_auth.hcl", logger)
	if err == nil {
		t.Fatal("LoadConfig should return an error when use_auto_auth_token=true and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_NoListeners(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	_, err := LoadConfig("./test-fixtures/bad-config-cache-no-listeners.hcl", logger)
	if err == nil {
		t.Fatal("LoadConfig should return an error when cache section present and no listeners present")
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_NoSink(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-no-sink.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				WrapTTL:   300 * time.Second,
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		Cache: &Cache{
			UseAutoAuthToken: true,
		},
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":     "127.0.0.1:8300",
					"tls_disable": true,
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

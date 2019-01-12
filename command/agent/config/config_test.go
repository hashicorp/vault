package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

func TestLoadConfigFile(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	os.Setenv("TEST_AAD_ENV", "aad")
	defer os.Unsetenv("TEST_AAD_ENV")

	testLoadConfig := func(file string, expectedConfig *Config, expectedError error) {
		config, err := LoadConfig(file, logger)
		if expectedError != nil {
			if err == nil {
				t.Fatal("expected a non-nil error")
			}
			if err.Error() != expectedError.Error() {
				t.Fatalf("Errors differ. Got:\n\t%v\nExpected:\n\t%v", err, expectedError)
			}
			return
		}
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if diff := deep.Equal(config, expectedConfig); diff != nil {
			t.Fatal(diff)
		}
	}

	t.Run("non-embedded-type", func(t *testing.T) {
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

		testLoadConfig("./test-fixtures/config.hcl", expected, nil)
	})

	t.Run("embedded-type", func(t *testing.T) {
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

		testLoadConfig("./test-fixtures/config-embedded-type.hcl", expected, nil)
	})

	t.Run("multiple-wrap-ttls", func(t *testing.T) {
		err := errors.New("error parsing 'auto_auth': error parsing 'sink' stanzas: sink.file 'wrap_ttl' may be in either the 'method' or 'sink' block, but not in both")
		testLoadConfig("./test-fixtures/config-multiple-wrap-ttls.hcl", nil, err)
	})
}
package config

import (
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

	testLoadConfig := func(file string, expected *Config) {
		config, err := LoadConfig(file, logger)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if diff := deep.Equal(config, expected); diff != nil {
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

		testLoadConfig("./test-fixtures/config.hcl", expected)
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

		testLoadConfig("./test-fixtures/config-embedded-type.hcl", expected)
	})

	t.Run("no-sinks", func(t *testing.T) {
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
			PidFile: "./pidfile",
		}

		testLoadConfig("./test-fixtures/config-no-sinks.hcl", expected)
	})
}

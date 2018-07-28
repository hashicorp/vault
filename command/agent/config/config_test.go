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

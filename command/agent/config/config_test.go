package config

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

func TestLoadConfigFile(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfig("./test-fixtures/config.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws-iam",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				&Sink{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
				&Sink{
					Type:    "file",
					WrapTTL: 5 * time.Minute,
					Config: map[string]interface{}{
						"path":     "/tmp/file-bar",
						"wrap_ttl": "5m",
					},
				},
			},
		},
		PidFile: "./pidfile",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

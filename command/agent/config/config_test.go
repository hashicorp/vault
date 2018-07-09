package config

import (
	"reflect"
	"testing"

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
			Vault: &Vault{Address: "http://127.0.0.1:8200", TLSSkipVerify: true, CAPath: "", CACert: ""},
			TokenSinks: []*TokenSink{
				&TokenSink{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
				&TokenSink{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-bar",
					},
				},
			},
		},
		PidFile: "./pidfile",
	}

	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

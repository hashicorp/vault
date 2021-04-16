package cassandra

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/testhelpers/cassandra"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/stretchr/testify/require"
)

var (
	insecureFileMounts = map[string]string{
		"test-fixtures/no_tls/cassandra.yaml": "/etc/cassandra/cassandra.yaml",
	}
	secureFileMounts = map[string]string{
		"test-fixtures/with_tls/cassandra.yaml": "/etc/cassandra/cassandra.yaml",
		"test-fixtures/with_tls/keystore.jks":   "/etc/cassandra/keystore.jks",
		"test-fixtures/with_tls/.cassandra":     "/root/.cassandra/",
	}
)

func TestTLSConnection(t *testing.T) {
	type testCase struct {
		config    map[string]interface{}
		expectErr bool
		// errorMsg is only checked if expectErr is true. This also a partial string match, so if this value shows up
		// anywhere in the error it will pass the assertion
		errorMsg string
	}

	tests := map[string]testCase{
		"tls not specified": {
			config:    map[string]interface{}{},
			expectErr: true,
			errorMsg:  "EOF",
		},
		"unrecognized certificate": {
			config: map[string]interface{}{
				"tls": "true",
			},
			expectErr: true,
			errorMsg:  "certificate signed by unknown authority",
		},
		"insecure TLS": {
			config: map[string]interface{}{
				"tls":          "true",
				"insecure_tls": true,
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			host, cleanup := cassandra.PrepareTestContainer(t,
				cassandra.Version("3.11.9"),
				cassandra.CopyFromTo(secureFileMounts),
				cassandra.SslOpts(&gocql.SslOptions{
					Config:                 &tls.Config{InsecureSkipVerify: true},
					EnableHostVerification: false,
				}),
			)
			defer cleanup()

			// Set values that we don't know until the cassandra container is started
			config := map[string]interface{}{
				"hosts":            host.ConnectionURL(),
				"port":             host.Port,
				"username":         "cassandra",
				"password":         "cassandra",
				"protocol_version": "3",
				"connect_timeout":  "20s",
			}
			// Then add any values specified in the test config. Generally for these tests they shouldn't overlap
			for k, v := range test.config {
				config[k] = v
			}

			db := new()
			initReq := dbplugin.InitializeRequest{
				Config:           config,
				VerifyConnection: true,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := db.Initialize(ctx, initReq)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if err != nil {
				require.Contains(t, err.Error(), test.errorMsg)
			}
		})
	}
}

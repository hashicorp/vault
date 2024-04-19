// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/transit"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTransit_Issue_2958(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": audit.NewFileBackend,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err := client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "/dev/null",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foo", map[string]interface{}{
		"type": "ecdsa-p256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foobar", map[string]interface{}{
		"type": "ecdsa-p384",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/bar", map[string]interface{}{
		"type": "ed25519",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foobar")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/bar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransit_CreateKeyWithAutorotation(t *testing.T) {
	tests := map[string]struct {
		autoRotatePeriod interface{}
		shouldError      bool
		expectedValue    time.Duration
	}{
		"default (no value)": {
			shouldError: false,
		},
		"0 (int)": {
			autoRotatePeriod: 0,
			shouldError:      false,
			expectedValue:    0,
		},
		"0 (string)": {
			autoRotatePeriod: "0",
			shouldError:      false,
			expectedValue:    0,
		},
		"5 seconds": {
			autoRotatePeriod: "5s",
			shouldError:      true,
		},
		"5 hours": {
			autoRotatePeriod: "5h",
			shouldError:      false,
			expectedValue:    5 * time.Hour,
		},
		"negative value": {
			autoRotatePeriod: "-1800s",
			shouldError:      true,
		},
		"invalid string": {
			autoRotatePeriod: "this shouldn't work",
			shouldError:      true,
		},
	}

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client
	err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		t.Fatal(err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			keyNameBytes, err := uuid.GenerateRandomBytes(16)
			if err != nil {
				t.Fatal(err)
			}
			keyName := hex.EncodeToString(keyNameBytes)

			_, err = client.Logical().Write(fmt.Sprintf("transit/keys/%s", keyName), map[string]interface{}{
				"auto_rotate_period": test.autoRotatePeriod,
			})
			switch {
			case test.shouldError && err == nil:
				t.Fatal("expected non-nil error")
			case !test.shouldError && err != nil:
				t.Fatal(err)
			}

			if !test.shouldError {
				resp, err := client.Logical().Read(fmt.Sprintf("transit/keys/%s", keyName))
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				gotRaw, ok := resp.Data["auto_rotate_period"].(json.Number)
				if !ok {
					t.Fatal("returned value is of unexpected type")
				}
				got, err := gotRaw.Int64()
				if err != nil {
					t.Fatal(err)
				}
				want := int64(test.expectedValue.Seconds())
				if got != want {
					t.Fatalf("incorrect auto_rotate_period returned, got: %d, want: %d", got, want)
				}
			}
		})
	}
}

// TestTransit_CreateKey validates transit key creation functionality
func TestTransit_CreateKey(t *testing.T) {
	testCases := map[string]struct {
		creationParams map[string]interface{}
		shouldError    bool
	}{
		"AES-128": {
			creationParams: map[string]interface{}{"type": "aes128-gcm96"},
			shouldError:    false,
		},
		"AES-256": {
			creationParams: map[string]interface{}{"type": "aes256-gcm96"},
			shouldError:    false,
		},
		"CHACHA20": {
			creationParams: map[string]interface{}{"type": "chacha20-poly1305"},
			shouldError:    false,
		},
		"ECDSA-P256": {
			creationParams: map[string]interface{}{"type": "ecdsa-p256"},
			shouldError:    false,
		},
		"ECDSA-P384": {
			creationParams: map[string]interface{}{"type": "ecdsa-p384"},
			shouldError:    false,
		},
		"ECDSA-P521": {
			creationParams: map[string]interface{}{"type": "ecdsa-p521"},
			shouldError:    false,
		},
		"RSA_2048": {
			creationParams: map[string]interface{}{"type": "rsa-2048"},
			shouldError:    false,
		},
		"RSA_3072": {
			creationParams: map[string]interface{}{"type": "rsa-3072"},
			shouldError:    false,
		},
		"RSA_4096": {
			creationParams: map[string]interface{}{"type": "rsa-4096"},
			shouldError:    false,
		},
		"HMAC": {
			creationParams: map[string]interface{}{"type": "hmac", "key_size": 128},
			shouldError:    false,
		},
		"AES-128 CMAC": {
			creationParams: map[string]interface{}{"type": "aes128-cmac"},
			shouldError:    false,
		},
		"AES-256 CMAC": {
			creationParams: map[string]interface{}{"type": "aes256-cmac"},
			shouldError:    false,
		},
		"bad key type": {
			creationParams: map[string]interface{}{"type": "fake-key-type"},
			shouldError:    true,
		},
	}

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client
	err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		t.Fatal(err)
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			keyName, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("error generating key name: %s", err)
			}

			resp, err := client.Logical().Write(fmt.Sprintf("transit/keys/%s", keyName), tt.creationParams)
			if err != nil && !tt.shouldError {
				t.Fatalf("unexpected error creating key: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatal("expected error but got nil")
			}

			if err == nil {
				keyType, ok := resp.Data["type"]
				if !ok {
					t.Fatal("missing key type in response")
				}

				if keyType != tt.creationParams["type"] {
					t.Fatalf("incorrect key type: expected %s, got %s", tt.creationParams["type"], keyType)
				}
			}
		})
	}
}

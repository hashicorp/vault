// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTransit_ConfigSettings(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	doReq := func(req *logical.Request) *logical.Response {
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("got err:\n%#v\nreq:\n%#v\n", err, *req)
		}
		return resp
	}
	doErrReq := func(req *logical.Request) {
		resp, err := b.HandleRequest(context.Background(), req)
		if err == nil {
			if resp == nil || !resp.IsError() {
				t.Fatalf("expected error; req:\n%#v\n", *req)
			}
		}
	}

	// First create a key
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/aes256",
		Data: map[string]interface{}{
			"derived": true,
		},
	}
	doReq(req)

	req.Path = "keys/aes128"
	req.Data["type"] = "aes128-gcm96"
	doReq(req)

	req.Path = "keys/ed"
	req.Data["type"] = "ed25519"
	doReq(req)

	delete(req.Data, "derived")

	req.Path = "keys/p256"
	req.Data["type"] = "ecdsa-p256"
	doReq(req)

	req.Path = "keys/p384"
	req.Data["type"] = "ecdsa-p384"
	doReq(req)

	req.Path = "keys/p521"
	req.Data["type"] = "ecdsa-p521"
	doReq(req)

	delete(req.Data, "type")

	req.Path = "keys/aes128/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/aes256/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/ed/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/p256/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/p384/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/p521/rotate"
	doReq(req)
	doReq(req)
	doReq(req)
	doReq(req)

	req.Path = "keys/aes256/config"
	// Too high
	req.Data["min_decryption_version"] = 7
	doErrReq(req)
	// Too low
	req.Data["min_decryption_version"] = -1
	doErrReq(req)

	delete(req.Data, "min_decryption_version")
	// Too high
	req.Data["min_encryption_version"] = 7
	doErrReq(req)
	// Too low
	req.Data["min_encryption_version"] = 7
	doErrReq(req)

	// Not allowed, cannot decrypt
	req.Data["min_decryption_version"] = 3
	req.Data["min_encryption_version"] = 2
	doErrReq(req)

	// Allowed
	req.Data["min_decryption_version"] = 2
	req.Data["min_encryption_version"] = 3
	doReq(req)
	req.Path = "keys/aes128/config"
	doReq(req)
	req.Path = "keys/ed/config"
	doReq(req)
	req.Path = "keys/p256/config"
	doReq(req)
	req.Path = "keys/p384/config"
	doReq(req)

	req.Path = "keys/p521/config"
	doReq(req)

	req.Data = map[string]interface{}{
		"plaintext": "abcd",
		"input":     "abcd",
		"context":   "abcd",
	}

	maxKeyVersion := 5
	key := "aes256"

	testHMAC := func(ver int, valid bool) {
		req.Path = "hmac/" + key
		delete(req.Data, "hmac")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["hmac"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong hmac version")
		}

		req.Path = "verify/" + key
		delete(req.Data, "key_version")
		req.Data["hmac"] = resp.Data["hmac"]
		doReq(req)
	}

	testEncryptDecrypt := func(ver int, valid bool) {
		req.Path = "encrypt/" + key
		delete(req.Data, "ciphertext")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["ciphertext"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong encryption version")
		}

		req.Path = "decrypt/" + key
		delete(req.Data, "key_version")
		req.Data["ciphertext"] = resp.Data["ciphertext"]
		doReq(req)
	}
	testEncryptDecrypt(5, true)
	testEncryptDecrypt(4, true)
	testEncryptDecrypt(3, true)
	testEncryptDecrypt(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	key = "aes128"
	testEncryptDecrypt(5, true)
	testEncryptDecrypt(4, true)
	testEncryptDecrypt(3, true)
	testEncryptDecrypt(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	delete(req.Data, "plaintext")
	req.Data["input"] = "abcd"
	key = "ed"
	testSignVerify := func(ver int, valid bool) {
		req.Path = "sign/" + key
		delete(req.Data, "signature")
		if ver == maxKeyVersion {
			delete(req.Data, "key_version")
		} else {
			req.Data["key_version"] = ver
		}

		if !valid {
			doErrReq(req)
			return
		}

		resp := doReq(req)
		ct := resp.Data["signature"].(string)
		if strings.Split(ct, ":")[1] != "v"+strconv.Itoa(ver) {
			t.Fatal("wrong signature version")
		}

		req.Path = "verify/" + key
		delete(req.Data, "key_version")
		req.Data["signature"] = resp.Data["signature"]
		doReq(req)
	}
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	delete(req.Data, "context")
	key = "p256"
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	key = "p384"
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)

	key = "p521"
	testSignVerify(5, true)
	testSignVerify(4, true)
	testSignVerify(3, true)
	testSignVerify(2, false)
	testHMAC(5, true)
	testHMAC(4, true)
	testHMAC(3, true)
	testHMAC(2, false)
}

func TestTransit_UpdateKeyConfigWithAutorotation(t *testing.T) {
	tests := map[string]struct {
		initialAutoRotatePeriod interface{}
		newAutoRotatePeriod     interface{}
		shouldError             bool
		expectedValue           time.Duration
	}{
		"default (no value)": {
			initialAutoRotatePeriod: "5h",
			shouldError:             false,
			expectedValue:           5 * time.Hour,
		},
		"0 (int)": {
			initialAutoRotatePeriod: "5h",
			newAutoRotatePeriod:     0,
			shouldError:             false,
			expectedValue:           0,
		},
		"0 (string)": {
			initialAutoRotatePeriod: "5h",
			newAutoRotatePeriod:     0,
			shouldError:             false,
			expectedValue:           0,
		},
		"5 seconds": {
			newAutoRotatePeriod: "5s",
			shouldError:         true,
		},
		"5 hours": {
			newAutoRotatePeriod: "5h",
			shouldError:         false,
			expectedValue:       5 * time.Hour,
		},
		"negative value": {
			newAutoRotatePeriod: "-1800s",
			shouldError:         true,
		},
		"invalid string": {
			newAutoRotatePeriod: "this shouldn't work",
			shouldError:         true,
		},
	}

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
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
				"auto_rotate_period": test.initialAutoRotatePeriod,
			})
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Logical().Write(fmt.Sprintf("transit/keys/%s/config", keyName), map[string]interface{}{
				"auto_rotate_period": test.newAutoRotatePeriod,
			})
			switch {
			case test.shouldError && err == nil:
				t.Fatal("expected non-nil error")
			case !test.shouldError && err != nil:
				t.Fatal(err)
			}

			if !test.shouldError {
				resp, err = client.Logical().Read(fmt.Sprintf("transit/keys/%s", keyName))
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

func TestTransit_KeyUsagesInConfigResponse(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	doReq := func(op logical.Operation, path string, data map[string]interface{}) *logical.Response {
		t.Helper()
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: op,
			Path:      path,
			Data:      data,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("unexpected error for %s %s: err=%v resp=%v", op, path, err, resp)
		}
		return resp
	}

	tests := []struct {
		keyType      string
		creationData map[string]interface{}
		expected     []interface{}
	}{
		{"aes256-gcm96", nil, []interface{}{"aead-encryption"}},
		{"aes128-gcm96", nil, []interface{}{"aead-encryption"}},
		{"chacha20-poly1305", nil, []interface{}{"aead-encryption"}},
		{"ecdsa-p256", nil, []interface{}{"digital-signature"}},
		{"ecdsa-p384", nil, []interface{}{"digital-signature"}},
		{"ecdsa-p521", nil, []interface{}{"digital-signature"}},
		{"ed25519", nil, []interface{}{"digital-signature"}},
		{"rsa-2048", nil, []interface{}{"asymmetric-encryption", "digital-signature"}},
		{"rsa-3072", nil, []interface{}{"asymmetric-encryption", "digital-signature"}},
		{"rsa-4096", nil, []interface{}{"asymmetric-encryption", "digital-signature"}},
		{"hmac", map[string]interface{}{"key_size": 32}, []interface{}{"message-authentication"}},
	}

	for _, tt := range tests {
		t.Run(tt.keyType, func(t *testing.T) {
			keyName := "test-" + tt.keyType

			// Build creation parameters
			createData := map[string]interface{}{"type": tt.keyType}
			for k, v := range tt.creationData {
				createData[k] = v
			}
			doReq(logical.UpdateOperation, "keys/"+keyName, createData)

			// Config endpoint must include key_usages
			resp := doReq(logical.UpdateOperation, "keys/"+keyName+"/config", map[string]interface{}{})
			usages, ok := resp.Data["key_usages"]
			if !ok {
				t.Fatalf("key_usages missing from config response for key type %s", tt.keyType)
			}
			gotSlice, ok := usages.([]string)
			if !ok {
				t.Fatalf("key_usages is not []string for key type %s, got %T", tt.keyType, usages)
			}
			// Convert to []interface{} for comparison
			got := make([]interface{}, len(gotSlice))
			for i, v := range gotSlice {
				got[i] = v
			}
			if len(got) != len(tt.expected) {
				t.Fatalf("key_usages mismatch for %s: got %v, want %v", tt.keyType, got, tt.expected)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("key_usages[%d] mismatch for %s: got %v, want %v", i, tt.keyType, got[i], tt.expected[i])
				}
			}
		})
	}
}


func TestTransit_ConfigWarnsOnIrreversibleFlagDisable(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	doReq := func(op logical.Operation, path string, data map[string]interface{}) *logical.Response {
		t.Helper()
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: op,
			Path:      path,
			Data:      data,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("unexpected error for %s %s: err=%v resp=%v", op, path, err, resp)
		}
		return resp
	}

	hasWarning := func(resp *logical.Response, substr string) bool {
		t.Helper()
		for _, w := range resp.Warnings {
			if strings.Contains(w, substr) {
				return true
			}
		}
		return false
	}

	// Create an exportable key with plaintext backup enabled.
	doReq(logical.UpdateOperation, "keys/irrev", map[string]interface{}{
		"exportable":             true,
		"allow_plaintext_backup": true,
	})

	// Attempt to disable both irreversible flags. Policy must stay enabled,
	// and the response must warn that the supplied values were ignored.
	resp := doReq(logical.UpdateOperation, "keys/irrev/config", map[string]interface{}{
		"exportable":             false,
		"allow_plaintext_backup": false,
	})
	if resp.Data["exportable"] != true {
		t.Fatalf("exportable should remain true after disable attempt, got %#v", resp.Data["exportable"])
	}
	if resp.Data["allow_plaintext_backup"] != true {
		t.Fatalf("allow_plaintext_backup should remain true after disable attempt, got %#v", resp.Data["allow_plaintext_backup"])
	}
	if !hasWarning(resp, "exportable cannot be disabled once set") {
		t.Fatalf("expected exportable disable warning, got %#v", resp.Warnings)
	}
	if !hasWarning(resp, "allow_plaintext_backup cannot be disabled once set") {
		t.Fatalf("expected allow_plaintext_backup disable warning, got %#v", resp.Warnings)
	}
}

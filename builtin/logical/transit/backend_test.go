// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto"
	"crypto/ed25519"
	cryptoRand "crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/observations"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/billing"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

const (
	testPlaintext = "The quick brown fox"
)

func createBackendWithStorage(t testing.TB) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, _ := Backend(context.Background(), config)
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	err := b.Backend.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	b.billingDataCounts = billing.DataProtectionCallCounts{
		Transit: &atomic.Uint64{},
	}
	return b, config.StorageView
}

func createBackendWithSysView(t testing.TB) (*backend, logical.Storage) {
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	conf := &logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	}

	b, _ := Backend(context.Background(), conf)
	if b == nil {
		t.Fatal("failed to create backend")
	}

	err := b.Backend.Setup(context.Background(), conf)
	if err != nil {
		t.Fatal(err)
	}
	b.billingDataCounts = billing.DataProtectionCallCounts{
		Transit: &atomic.Uint64{},
	}

	return b, storage
}

func createBackendWithSysViewWithStorage(t testing.TB, s logical.Storage) *backend {
	sysView := logical.TestSystemView()

	conf := &logical.BackendConfig{
		StorageView: s,
		System:      sysView,
	}

	b, _ := Backend(context.Background(), conf)
	if b == nil {
		t.Fatal("failed to create backend")
	}

	err := b.Backend.Setup(context.Background(), conf)
	if err != nil {
		t.Fatal(err)
	}
	b.billingDataCounts = billing.DataProtectionCallCounts{
		Transit: &atomic.Uint64{},
	}

	return b
}

func createBackendWithForceNoCacheWithSysViewWithStorage(t testing.TB, s logical.Storage) *backend {
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true

	conf := &logical.BackendConfig{
		StorageView: s,
		System:      sysView,
	}

	b, _ := Backend(context.Background(), conf)
	if b == nil {
		t.Fatal("failed to create backend")
	}

	err := b.Backend.Setup(context.Background(), conf)
	if err != nil {
		t.Fatal(err)
	}
	b.billingDataCounts = billing.DataProtectionCallCounts{
		Transit: &atomic.Uint64{},
	}

	return b
}

func createBackendWithObservationRecorder(t testing.TB) (*backend, logical.Storage, *observations.TestObservationRecorder) {
	config := logical.TestBackendConfig()
	obsRecorder := observations.NewTestObservationRecorder()
	config.StorageView = &logical.InmemStorage{}
	config.ObservationRecorder = obsRecorder

	b, _ := Backend(context.Background(), config)
	require.NotNil(t, b)
	err := b.Backend.Setup(context.Background(), config)
	require.NoError(t, err)
	return b, config.StorageView, obsRecorder
}

func factoryWithObservationRecorder(t testing.TB) (logical.Factory, *observations.TestObservationRecorder) {
	obsRecorder := observations.NewTestObservationRecorder()
	return func(ctx context.Context, bc *logical.BackendConfig) (logical.Backend, error) {
		bc.ObservationRecorder = obsRecorder
		return Factory(ctx, bc)
	}, obsRecorder
}

func TestTransit_RSA(t *testing.T) {
	testTransit_RSA(t, "rsa-2048")
	testTransit_RSA(t, "rsa-3072")
	testTransit_RSA(t, "rsa-4096")
}

func testTransit_RSA(t *testing.T, keyType string) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	keyReq := &logical.Request{
		Path:      "keys/rsa",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": keyType,
		},
		Storage: storage,
	}

	resp, err = b.HandleRequest(context.Background(), keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"

	for _, padding := range []keysutil.PaddingScheme{keysutil.PaddingScheme_OAEP, keysutil.PaddingScheme_PKCS1v15, ""} {
		encryptReq := &logical.Request{
			Path:      "encrypt/rsa",
			Operation: logical.UpdateOperation,
			Storage:   storage,
			Data: map[string]interface{}{
				"plaintext": plaintext,
			},
		}

		if padding != "" {
			encryptReq.Data["padding_scheme"] = padding
		}

		resp, err = b.HandleRequest(context.Background(), encryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}

		ciphertext1 := resp.Data["ciphertext"].(string)

		decryptReq := &logical.Request{
			Path:      "decrypt/rsa",
			Operation: logical.UpdateOperation,
			Storage:   storage,
			Data: map[string]interface{}{
				"ciphertext": ciphertext1,
			},
		}
		if padding != "" {
			decryptReq.Data["padding_scheme"] = padding
		}

		resp, err = b.HandleRequest(context.Background(), decryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}

		decryptedPlaintext := resp.Data["plaintext"]

		if plaintext != decryptedPlaintext {
			t.Fatalf("bad: plaintext; expected: %q\nactual: %q", plaintext, decryptedPlaintext)
		}

		// Rotate the key
		rotateReq := &logical.Request{
			Path:      "keys/rsa/rotate",
			Operation: logical.UpdateOperation,
			Storage:   storage,
		}
		resp, err = b.HandleRequest(context.Background(), rotateReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}

		// Encrypt again
		resp, err = b.HandleRequest(context.Background(), encryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}
		ciphertext2 := resp.Data["ciphertext"].(string)

		if ciphertext1 == ciphertext2 {
			t.Fatalf("expected different ciphertexts")
		}

		// See if the older ciphertext can still be decrypted
		resp, err = b.HandleRequest(context.Background(), decryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}
		if resp.Data["plaintext"].(string) != plaintext {
			t.Fatal("failed to decrypt old ciphertext after rotating the key")
		}

		// Decrypt the new ciphertext
		decryptReq.Data = map[string]interface{}{
			"ciphertext": ciphertext2,
		}
		if padding != "" {
			decryptReq.Data["padding_scheme"] = padding
		}

		resp, err = b.HandleRequest(context.Background(), decryptReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}
		if resp.Data["plaintext"].(string) != plaintext {
			t.Fatal("failed to decrypt ciphertext after rotating the key")
		}
	}

	signReq := &logical.Request{
		Path:      "sign/rsa",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"input": plaintext,
		},
	}
	resp, err = b.HandleRequest(context.Background(), signReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	signature := resp.Data["signature"].(string)

	verifyReq := &logical.Request{
		Path:      "verify/rsa",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"input":     plaintext,
			"signature": signature,
		},
	}

	resp, err = b.HandleRequest(context.Background(), verifyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if !resp.Data["valid"].(bool) {
		t.Fatalf("failed to verify the RSA signature")
	}

	signReq.Data = map[string]interface{}{
		"input":          plaintext,
		"hash_algorithm": "invalid",
	}
	resp, err = b.HandleRequest(context.Background(), signReq)
	if err == nil {
		t.Fatal(err)
	}

	signReq.Data = map[string]interface{}{
		"input":          plaintext,
		"hash_algorithm": "sha2-512",
	}
	resp, err = b.HandleRequest(context.Background(), signReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	signature = resp.Data["signature"].(string)

	verifyReq.Data = map[string]interface{}{
		"input":     plaintext,
		"signature": signature,
	}
	resp, err = b.HandleRequest(context.Background(), verifyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if resp.Data["valid"].(bool) {
		t.Fatalf("expected validation to fail")
	}

	verifyReq.Data = map[string]interface{}{
		"input":          plaintext,
		"signature":      signature,
		"hash_algorithm": "sha2-512",
	}
	resp, err = b.HandleRequest(context.Background(), verifyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if !resp.Data["valid"].(bool) {
		t.Fatalf("failed to verify the RSA signature")
	}

	// Take a random hash and sign it using PKCSv1_5_NoOID.
	hash := "P8m2iUWdc4+MiKOkiqnjNUIBa3pAUuABqqU2/KdIE8s="
	signReq.Data = map[string]interface{}{
		"input":               hash,
		"hash_algorithm":      "none",
		"signature_algorithm": "pkcs1v15",
		"prehashed":           true,
	}
	resp, err = b.HandleRequest(context.Background(), signReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	signature = resp.Data["signature"].(string)

	verifyReq.Data = map[string]interface{}{
		"input":               hash,
		"signature":           signature,
		"hash_algorithm":      "none",
		"signature_algorithm": "pkcs1v15",
		"prehashed":           true,
	}
	resp, err = b.HandleRequest(context.Background(), verifyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if !resp.Data["valid"].(bool) {
		t.Fatalf("failed to verify the RSA signature")
	}
}

func TestBackend_basic(t *testing.T) {
	factory, obsRecorder := factoryWithObservationRecorder(t)
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepListPolicy(t, "test", true),
			testAccStepWritePolicy(t, "test", false, obsRecorder),
			testAccStepListPolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false, obsRecorder),
			testAccStepEncrypt(t, "test", testPlaintext, decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEncrypt(t, "test", "", decryptData),
			testAccStepDecrypt(t, "test", "", decryptData),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test", obsRecorder),
			testAccStepWritePolicy(t, "test", false, obsRecorder),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDisableDeletion(t, "test"),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test", obsRecorder),
			testAccStepReadPolicy(t, "test", true, false, obsRecorder),
		},
	})
}

func TestBackend_upsert(t *testing.T) {
	factory, obsRecorder := factoryWithObservationRecorder(t)
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepReadPolicy(t, "test", true, false, obsRecorder),
			testAccStepListPolicy(t, "test", true),
			testAccStepEncryptUpsert(t, "test", testPlaintext, decryptData),
			testAccStepListPolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false, obsRecorder),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
		},
	})
}

func TestBackend_datakey(t *testing.T) {
	factory, obsRecorder := factoryWithObservationRecorder(t)
	dataKeyInfo := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepListPolicy(t, "test", true),
			testAccStepWritePolicy(t, "test", false, obsRecorder),
			testAccStepListPolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false, nil),
			testAccStepWriteDatakey(t, "test", false, 256, dataKeyInfo),
			testAccStepDecryptDatakey(t, "test", dataKeyInfo),
			testAccStepWriteDatakey(t, "test", true, 128, dataKeyInfo),
		},
	})
}

func TestBackend_rotation(t *testing.T) {
	defer os.Setenv("TRANSIT_ACC_KEY_TYPE", "")
	testBackendRotation(t)
	os.Setenv("TRANSIT_ACC_KEY_TYPE", "CHACHA")
	testBackendRotation(t)
}

func testBackendRotation(t *testing.T) {
	decryptData := make(map[string]interface{})
	encryptHistory := make(map[int]map[string]interface{})
	factory, obsRecorder := factoryWithObservationRecorder(t)
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepListPolicy(t, "test", true),
			testAccStepWritePolicy(t, "test", false, obsRecorder),
			testAccStepListPolicy(t, "test", false),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 0, encryptHistory),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 1, encryptHistory),
			testAccStepRotate(t, "test"), // now v2
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 2, encryptHistory),
			testAccStepRotate(t, "test"), // now v3
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 3, encryptHistory),
			testAccStepRotate(t, "test"), // now v4
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 99, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 3, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 99, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepAdjustPolicyMinDecryption(t, "test", 3),
			testAccStepAdjustPolicyMinEncryption(t, "test", 4),
			testAccStepReadPolicyWithVersions(t, "test", false, false, 3, 4, obsRecorder),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 3, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepAdjustPolicyMinDecryption(t, "test", 1),
			testAccStepReadPolicyWithVersions(t, "test", false, false, 1, 4, obsRecorder),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepRewrap(t, "test", decryptData, 4),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test", obsRecorder),
			testAccStepReadPolicy(t, "test", true, false, obsRecorder),
			testAccStepListPolicy(t, "test", true),
		},
	})
}

func TestBackend_basic_derived(t *testing.T) {
	decryptData := make(map[string]interface{})
	factory, obsRecorder := factoryWithObservationRecorder(t)
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: factory,
		Steps: []logicaltest.TestStep{
			testAccStepListPolicy(t, "test", true),
			testAccStepWritePolicy(t, "test", true, obsRecorder),
			testAccStepListPolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, true, obsRecorder),
			testAccStepEncryptContext(t, "test", testPlaintext, "my-cool-context", decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test", obsRecorder),
			testAccStepReadPolicy(t, "test", true, true, obsRecorder),
		},
	})
}

func testAccStepWritePolicy(t *testing.T, name string, derived bool, obsRecorder *observations.TestObservationRecorder) logicaltest.TestStep {
	ts := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
		},
		Check: func(resp *logical.Response) error {
			if obsRecorder == nil {
				return nil
			}
			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyWrite)
			if obs == nil {
				return fmt.Errorf("no observation")
			}
			if name != obs.Data["key_name"] {
				return fmt.Errorf("expected name %s, got %s", name, obs.Data["key_name"])
			}
			return nil
		},
	}
	if os.Getenv("TRANSIT_ACC_KEY_TYPE") == "CHACHA" {
		ts.Data["type"] = "chacha20-poly1305"
	}
	return ts
}

func testAccStepListPolicy(t *testing.T, name string, expectNone bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "keys",
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("missing response")
			}
			if expectNone {
				keysRaw, ok := resp.Data["keys"]
				if ok || keysRaw != nil {
					return fmt.Errorf("response data when expecting none")
				}
				return nil
			}
			if len(resp.Data) == 0 {
				return fmt.Errorf("no data returned")
			}

			var d struct {
				Keys []string `mapstructure:"keys"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if len(d.Keys) > 0 && d.Keys[0] != name {
				return fmt.Errorf("bad name: %#v", d)
			}
			if len(d.Keys) != 1 {
				return fmt.Errorf("only 1 key expected, %d returned", len(d.Keys))
			}
			return nil
		},
	}
}

func testAccStepAdjustPolicyMinDecryption(t *testing.T, name string, minVer int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"min_decryption_version": minVer,
		},
	}
}

func testAccStepAdjustPolicyMinEncryption(t *testing.T, name string, minVer int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"min_encryption_version": minVer,
		},
	}
}

func testAccStepDisableDeletion(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"deletion_allowed": false,
		},
	}
}

func testAccStepEnableDeletion(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"deletion_allowed": true,
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, name string, obsRecorder *observations.TestObservationRecorder) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "keys/" + name,
		Check: func(_ *logical.Response) error {
			if obsRecorder == nil {
				return nil
			}

			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyDelete)
			if obs == nil {
				return fmt.Errorf("expected observation of type %s but got none", ObservationTypeTransitKeyDelete)
			}
			if obs.Data["key_name"] != name {
				return fmt.Errorf("expected name %s, got %s", name, obs.Data["key_name"])
			}
			return nil
		},
	}
}

func testAccStepDeleteNotDisabledPolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "keys/" + name,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("got nil response instead of error")
			}
			if resp.IsError() {
				return nil
			}
			return fmt.Errorf("expected error but did not get one")
		},
	}
}

func testAccStepReadPolicy(t *testing.T, name string, expectNone, derived bool, obsRecorder *observations.TestObservationRecorder) logicaltest.TestStep {
	return testAccStepReadPolicyWithVersions(t, name, expectNone, derived, 1, 0, obsRecorder)
}

func testAccStepReadPolicyWithVersions(t *testing.T, name string, expectNone, derived bool, minDecryptionVersion int, minEncryptionVersion int, obsRecorder *observations.TestObservationRecorder) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil && !expectNone {
				return fmt.Errorf("missing response")
			} else if expectNone {
				if resp != nil {
					return fmt.Errorf("response when expecting none")
				}
				return nil
			}
			var d struct {
				Name                 string           `mapstructure:"name"`
				Key                  []byte           `mapstructure:"key"`
				Keys                 map[string]int64 `mapstructure:"keys"`
				Type                 string           `mapstructure:"type"`
				Derived              bool             `mapstructure:"derived"`
				KDF                  string           `mapstructure:"kdf"`
				DeletionAllowed      bool             `mapstructure:"deletion_allowed"`
				ConvergentEncryption bool             `mapstructure:"convergent_encryption"`
				MinDecryptionVersion int              `mapstructure:"min_decryption_version"`
				MinEncryptionVersion int              `mapstructure:"min_encryption_version"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Name != name {
				return fmt.Errorf("bad name: %#v", d)
			}
			if os.Getenv("TRANSIT_ACC_KEY_TYPE") == "CHACHA" {
				if d.Type != keysutil.KeyType(keysutil.KeyType_ChaCha20_Poly1305).String() {
					return fmt.Errorf("bad key type: %#v", d)
				}
			} else if d.Type != keysutil.KeyType(keysutil.KeyType_AES256_GCM96).String() {
				return fmt.Errorf("bad key type: %#v", d)
			}
			// Should NOT get a key back
			if d.Key != nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Keys == nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.MinDecryptionVersion != minDecryptionVersion {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.MinEncryptionVersion != minEncryptionVersion {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.DeletionAllowed {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Derived != derived {
				return fmt.Errorf("bad: %#v", d)
			}
			if derived && d.KDF != "hkdf_sha256" {
				return fmt.Errorf("bad: %#v", d)
			}

			if obsRecorder == nil {
				return nil
			}
			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyRead)
			if obs == nil {
				return fmt.Errorf("expected key read observation but found none")
			}
			if obs.Data == nil {
				return fmt.Errorf("observation data should not be nil")
			}
			keyName, ok := obs.Data["key_name"]
			if !ok {
				return fmt.Errorf("observation data missing key_name field")
			}
			if keyName != name {
				return fmt.Errorf("observation key_name mismatch: expected %s, got %v", name, keyName)
			}

			return nil
		},
	}
}

func testAccStepEncrypt(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
		},
		Check: func(resp *logical.Response) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			decryptData["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepEncryptUpsert(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
		},
		Check: func(resp *logical.Response) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			decryptData["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepEncryptContext(
	t *testing.T, name, plaintext, context string, decryptData map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
			"context":   base64.StdEncoding.EncodeToString([]byte(context)),
		},
		Check: func(resp *logical.Response) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			decryptData["ciphertext"] = d.Ciphertext
			decryptData["context"] = base64.StdEncoding.EncodeToString([]byte(context))
			return nil
		},
	}
}

func testAccStepDecrypt(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		Check: func(resp *logical.Response) error {
			var d struct {
				Plaintext string `mapstructure:"plaintext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			// Decode the base64
			plainRaw, err := base64.StdEncoding.DecodeString(d.Plaintext)
			if err != nil {
				return err
			}

			if string(plainRaw) != plaintext {
				return fmt.Errorf("plaintext mismatch: %s expect: %s, decryptData was %#v", plainRaw, plaintext, decryptData)
			}
			return nil
		},
	}
}

func testAccStepRewrap(
	t *testing.T, name string, decryptData map[string]interface{}, expectedVer int,
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "rewrap/" + name,
		Data:      decryptData,
		Check: func(resp *logical.Response) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			splitStrings := strings.Split(d.Ciphertext, ":")
			verString := splitStrings[1][1:]
			ver, err := strconv.Atoi(verString)
			if err != nil {
				return fmt.Errorf("error pulling out version from verString %q, ciphertext was %s", verString, d.Ciphertext)
			}
			if ver != expectedVer {
				return fmt.Errorf("did not get expected version")
			}
			decryptData["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepEncryptVX(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
	ver int, encryptHistory map[int]map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
		},
		Check: func(resp *logical.Response) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			splitStrings := strings.Split(d.Ciphertext, ":")
			splitStrings[1] = "v" + strconv.Itoa(ver)
			ciphertext := strings.Join(splitStrings, ":")
			decryptData["ciphertext"] = ciphertext
			encryptHistory[ver] = map[string]interface{}{
				"ciphertext": ciphertext,
			}
			return nil
		},
	}
}

func testAccStepLoadVX(
	t *testing.T, name string, decryptData map[string]interface{},
	ver int, encryptHistory map[int]map[string]interface{},
) logicaltest.TestStep {
	// This is really a no-op to allow us to do data manip in the check function
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *logical.Response) error {
			decryptData["ciphertext"] = encryptHistory[ver]["ciphertext"].(string)
			return nil
		},
	}
}

func testAccStepDecryptExpectFailure(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if !resp.IsError() {
				return fmt.Errorf("expected error")
			}
			return nil
		},
	}
}

func testAccStepRotate(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/rotate",
	}
}

func testAccStepWriteDatakey(t *testing.T, name string,
	noPlaintext bool, bits int,
	dataKeyInfo map[string]interface{},
) logicaltest.TestStep {
	data := map[string]interface{}{}
	subPath := "plaintext"
	if noPlaintext {
		subPath = "wrapped"
	}
	if bits != 256 {
		data["bits"] = bits
	}
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "datakey/" + subPath + "/" + name,
		Data:      data,
		Check: func(resp *logical.Response) error {
			var d struct {
				Plaintext  string `mapstructure:"plaintext"`
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if noPlaintext && len(d.Plaintext) != 0 {
				return fmt.Errorf("received plaintxt when we disabled it")
			}
			if !noPlaintext {
				if len(d.Plaintext) == 0 {
					return fmt.Errorf("did not get plaintext when we expected it")
				}
				dataKeyInfo["plaintext"] = d.Plaintext
				plainBytes, err := base64.StdEncoding.DecodeString(d.Plaintext)
				if err != nil {
					return fmt.Errorf("could not base64 decode plaintext string %q", d.Plaintext)
				}
				if len(plainBytes)*8 != bits {
					return fmt.Errorf("returned key does not have correct bit length")
				}
			}
			dataKeyInfo["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepDecryptDatakey(t *testing.T, name string,
	dataKeyInfo map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      dataKeyInfo,
		Check: func(resp *logical.Response) error {
			var d struct {
				Plaintext string `mapstructure:"plaintext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Plaintext != dataKeyInfo["plaintext"].(string) {
				return fmt.Errorf("plaintext mismatch: got %q, expected %q, decryptData was %#v", d.Plaintext, dataKeyInfo["plaintext"].(string), resp.Data)
			}
			return nil
		},
	}
}

func TestKeyUpgrade(t *testing.T) {
	key, _ := uuid.GenerateRandomBytes(32)
	p := &keysutil.Policy{
		Name: "test",
		Key:  key,
		Type: keysutil.KeyType_AES256_GCM96,
	}

	p.MigrateKeyToKeysMap()

	if p.Key != nil ||
		p.Keys == nil ||
		len(p.Keys) != 1 ||
		!reflect.DeepEqual(p.Keys[strconv.Itoa(1)].Key, key) {
		t.Errorf("bad key migration, result is %#v", p.Keys)
	}
}

func TestDerivedKeyUpgrade(t *testing.T) {
	testDerivedKeyUpgrade(t, keysutil.KeyType_AES256_GCM96)
	testDerivedKeyUpgrade(t, keysutil.KeyType_ChaCha20_Poly1305)
}

func testDerivedKeyUpgrade(t *testing.T, keyType keysutil.KeyType) {
	storage := &logical.InmemStorage{}
	key, _ := uuid.GenerateRandomBytes(32)
	keyContext, _ := uuid.GenerateRandomBytes(32)

	p := &keysutil.Policy{
		Name:    "test",
		Key:     key,
		Type:    keyType,
		Derived: true,
	}

	p.MigrateKeyToKeysMap()
	p.Upgrade(context.Background(), storage, cryptoRand.Reader) // Need to run the upgrade code to make the migration stick

	if p.KDF != keysutil.Kdf_hmac_sha256_counter {
		t.Fatalf("bad KDF value by default; counter val is %d, KDF val is %d, policy is %#v", keysutil.Kdf_hmac_sha256_counter, p.KDF, *p)
	}

	derBytesOld, err := p.GetKey(keyContext, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	derBytesOld2, err := p.GetKey(keyContext, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(derBytesOld, derBytesOld2) {
		t.Fatal("mismatch of same context alg")
	}

	p.KDF = keysutil.Kdf_hkdf_sha256
	if p.NeedsUpgrade() {
		t.Fatal("expected no upgrade needed")
	}

	derBytesNew, err := p.GetKey(keyContext, 1, 64)
	if err != nil {
		t.Fatal(err)
	}

	derBytesNew2, err := p.GetKey(keyContext, 1, 64)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(derBytesNew, derBytesNew2) {
		t.Fatal("mismatch of same context alg")
	}

	if reflect.DeepEqual(derBytesOld, derBytesNew) {
		t.Fatal("match of different context alg")
	}
}

func TestConvergentEncryption(t *testing.T) {
	testConvergentEncryptionCommon(t, 0, keysutil.KeyType_AES256_GCM96)
	testConvergentEncryptionCommon(t, 2, keysutil.KeyType_AES128_GCM96)
	testConvergentEncryptionCommon(t, 2, keysutil.KeyType_AES256_GCM96)
	testConvergentEncryptionCommon(t, 2, keysutil.KeyType_ChaCha20_Poly1305)
	testConvergentEncryptionCommon(t, 3, keysutil.KeyType_AES128_GCM96)
	testConvergentEncryptionCommon(t, 3, keysutil.KeyType_AES256_GCM96)
	testConvergentEncryptionCommon(t, 3, keysutil.KeyType_ChaCha20_Poly1305)
}

func testConvergentEncryptionCommon(t *testing.T, ver int, keyType keysutil.KeyType) {
	b, storage := createBackendWithSysView(t)

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/testkeynonderived",
		Data: map[string]interface{}{
			"derived":               false,
			"convergent_encryption": true,
			"type":                  keyType.String(),
		},
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.IsError() {
		t.Fatalf("bad: expected error response, got %#v", *resp)
	}

	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/testkey",
		Data: map[string]interface{}{
			"derived":               true,
			"convergent_encryption": true,
			"type":                  keyType.String(),
		},
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, resp, "expected populated request")

	p, err := keysutil.LoadPolicy(context.Background(), storage, path.Join("policy", "testkey"))
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("got nil policy")
	}

	if ver > 2 {
		p.ConvergentVersion = -1
	} else {
		p.ConvergentVersion = ver
	}
	err = p.Persist(context.Background(), storage)
	if err != nil {
		t.Fatal(err)
	}
	b.invalidate(context.Background(), "policy/testkey")

	if ver < 3 {
		// There will be an embedded key version of 3, so specifically clear it
		key := p.Keys[strconv.Itoa(p.LatestVersion)]
		key.ConvergentVersion = 0
		p.Keys[strconv.Itoa(p.LatestVersion)] = key
		err = p.Persist(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}
		b.invalidate(context.Background(), "policy/testkey")

		// Verify it
		p, err = keysutil.LoadPolicy(context.Background(), storage, path.Join(p.StoragePrefix, "policy", "testkey"))
		if err != nil {
			t.Fatal(err)
		}
		if p == nil {
			t.Fatal("got nil policy")
		}
		if p.ConvergentVersion != ver {
			t.Fatalf("bad convergent version %d", p.ConvergentVersion)
		}
		key = p.Keys[strconv.Itoa(p.LatestVersion)]
		if key.ConvergentVersion != 0 {
			t.Fatalf("bad convergent key version %d", key.ConvergentVersion)
		}
	}

	// First, test using an invalid length of nonce -- this is only used for v1 convergent
	req.Path = "encrypt/testkey"
	if ver < 2 {
		req.Data = map[string]interface{}{
			"plaintext": "emlwIHphcA==", // "zip zap"
			"nonce":     "Zm9vIGJhcg==", // "foo bar"
			"context":   "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
		}
		resp, err = b.HandleRequest(context.Background(), req)
		if err == nil {
			t.Fatalf("expected error, got nil, version is %d", ver)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if !resp.IsError() {
			t.Fatalf("expected error response, got %#v", *resp)
		}

		// Ensure we fail if we do not provide a nonce
		req.Data = map[string]interface{}{
			"plaintext": "emlwIHphcA==", // "zip zap"
			"context":   "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
		}
		resp, err = b.HandleRequest(context.Background(), req)
		if err == nil && (resp == nil || !resp.IsError()) {
			t.Fatal("expected error response")
		}
	}

	// Now test encrypting the same value twice
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==", // "zip zap"
		"context":   "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
	}
	if ver == 0 {
		req.Data["nonce"] = "b25ldHdvdGhyZWVl" // "onetwothreee"
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext1 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext2 := resp.Data["ciphertext"].(string)

	if ciphertext1 != ciphertext2 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext1, ciphertext2)
	}

	// For sanity, also check a different nonce value...
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==", // "zip zap"
		"context":   "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
	}
	if ver == 0 {
		req.Data["nonce"] = "dHdvdGhyZWVmb3Vy" // "twothreefour"
	} else {
		req.Data["context"] = "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOldandSdd7S"
	}

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext3 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext4 := resp.Data["ciphertext"].(string)

	if ciphertext3 != ciphertext4 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext3, ciphertext4)
	}
	if ciphertext1 == ciphertext3 {
		t.Fatalf("expected different ciphertexts")
	}

	// ...and a different context value
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==", // "zip zap"
		"context":   "qV4h9iQyvn+raODOer4JNAsOhkXBwdT4HZ677Ql4KLqXSU+Jk4C/fXBWbv6xkSYT",
	}
	if ver == 0 {
		req.Data["nonce"] = "dHdvdGhyZWVmb3Vy" // "twothreefour"
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext5 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext6 := resp.Data["ciphertext"].(string)

	if ciphertext5 != ciphertext6 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext5, ciphertext6)
	}
	if ciphertext1 == ciphertext5 {
		t.Fatalf("expected different ciphertexts")
	}
	if ciphertext3 == ciphertext5 {
		t.Fatalf("expected different ciphertexts")
	}

	// If running version 2, check upgrade handling
	if ver == 2 {
		curr, err := keysutil.LoadPolicy(context.Background(), storage, path.Join(p.StoragePrefix, "policy", "testkey"))
		if err != nil {
			t.Fatal(err)
		}
		if curr == nil {
			t.Fatal("got nil policy")
		}
		if curr.ConvergentVersion != 2 {
			t.Fatalf("bad convergent version %d", curr.ConvergentVersion)
		}
		key := curr.Keys[strconv.Itoa(curr.LatestVersion)]
		if key.ConvergentVersion != 0 {
			t.Fatalf("bad convergent key version %d", key.ConvergentVersion)
		}

		curr.ConvergentVersion = 3
		err = curr.Persist(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}
		b.invalidate(context.Background(), "policy/testkey")

		// Different algorithm, should be different value
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			t.Fatalf("got error response: %#v", *resp)
		}
		ciphertext7 := resp.Data["ciphertext"].(string)

		// Now do it via key-specified version
		if len(curr.Keys) != 1 {
			t.Fatalf("unexpected length of keys %d", len(curr.Keys))
		}
		key = curr.Keys[strconv.Itoa(curr.LatestVersion)]
		key.ConvergentVersion = 3
		curr.Keys[strconv.Itoa(curr.LatestVersion)] = key
		curr.ConvergentVersion = 2
		err = curr.Persist(context.Background(), storage)
		if err != nil {
			t.Fatal(err)
		}
		b.invalidate(context.Background(), "policy/testkey")

		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.IsError() {
			t.Fatalf("got error response: %#v", *resp)
		}
		ciphertext8 := resp.Data["ciphertext"].(string)

		if ciphertext7 != ciphertext8 {
			t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext7, ciphertext8)
		}
		if ciphertext6 == ciphertext7 {
			t.Fatalf("expected different ciphertexts")
		}
		if ciphertext3 == ciphertext7 {
			t.Fatalf("expected different ciphertexts")
		}
	}

	// Finally, check operations on empty values
	// First, check without setting a plaintext at all
	req.Data = map[string]interface{}{
		"context": "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
	}
	if ver == 0 {
		req.Data["nonce"] = "dHdvdGhyZWVmb3Vy" // "twothreefour"
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error response, got: %#v", *resp)
	}

	// Now set plaintext to empty
	req.Data = map[string]interface{}{
		"plaintext": "",
		"context":   "pWZ6t/im3AORd0lVYE0zBdKpX6Bl3/SvFtoVTPWbdkzjG788XmMAnOlxandSdd7S",
	}
	if ver == 0 {
		req.Data["nonce"] = "dHdvdGhyZWVmb3Vy" // "twothreefour"
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext7 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext8 := resp.Data["ciphertext"].(string)

	if ciphertext7 != ciphertext8 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext7, ciphertext8)
	}
}

func TestPolicyFuzzing(t *testing.T) {
	var be *backend
	sysView := logical.TestSystemView()
	sysView.CachingDisabledVal = true
	conf := &logical.BackendConfig{
		System: sysView,
	}

	be, _ = Backend(context.Background(), conf)
	be.Setup(context.Background(), conf)
	testPolicyFuzzingCommon(t, be)

	sysView.CachingDisabledVal = true
	be, _ = Backend(context.Background(), conf)
	be.Setup(context.Background(), conf)
	testPolicyFuzzingCommon(t, be)
}

func testPolicyFuzzingCommon(t *testing.T, be *backend) {
	storage := &logical.InmemStorage{}
	wg := sync.WaitGroup{}

	funcs := []string{"encrypt", "decrypt", "rotate", "change_min_version"}
	// keys := []string{"test1", "test2", "test3", "test4", "test5"}
	keys := []string{"test1", "test2", "test3"}

	// This is the goroutine loop
	doFuzzy := func(id int) {
		// Check for panics, otherwise notify we're done
		defer func() {
			wg.Done()
		}()

		// Holds the latest encrypted value for each key
		latestEncryptedText := map[string]string{}

		startTime := time.Now()
		req := &logical.Request{
			Storage: storage,
			Data:    map[string]interface{}{},
		}
		fd := &framework.FieldData{}

		var chosenFunc, chosenKey string

		// t.Errorf("Starting %d", id)
		for {
			// Stop after 10 seconds
			if time.Now().Sub(startTime) > 10*time.Second {
				return
			}

			// Pick a function and a key
			chosenFunc = funcs[rand.Int()%len(funcs)]
			chosenKey = keys[rand.Int()%len(keys)]

			fd.Raw = map[string]interface{}{
				"name": chosenKey,
			}
			fd.Schema = be.pathKeys().Fields

			// Try to write the key to make sure it exists
			_, err := be.pathPolicyWrite(context.Background(), req, fd)
			if err != nil {
				t.Errorf("got an error: %v", err)
			}

			switch chosenFunc {
			// Encrypt our plaintext and store the result
			case "encrypt":
				// t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				fd.Raw["plaintext"] = base64.StdEncoding.EncodeToString([]byte(testPlaintext))
				fd.Schema = be.pathEncrypt().Fields
				resp, err := be.pathEncryptWrite(context.Background(), req, fd)
				if err != nil {
					t.Errorf("got an error: %v, resp is %#v", err, *resp)
				}
				latestEncryptedText[chosenKey] = resp.Data["ciphertext"].(string)

			// Rotate to a new key version
			case "rotate":
				// t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				fd.Schema = be.pathRotate().Fields
				resp, err := be.pathRotateWrite(context.Background(), req, fd)
				if err != nil {
					t.Errorf("got an error: %v, resp is %#v, chosenKey is %s", err, *resp, chosenKey)
				}

			// Decrypt the ciphertext and compare the result
			case "decrypt":
				// t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				ct := latestEncryptedText[chosenKey]
				if ct == "" {
					continue
				}

				fd.Raw["ciphertext"] = ct
				fd.Schema = be.pathDecrypt().Fields
				resp, err := be.pathDecryptWrite(context.Background(), req, fd)
				if err != nil {
					// This could well happen since the min version is jumping around
					if resp.Data["error"].(string) == keysutil.ErrTooOld {
						continue
					}
					t.Errorf("got an error: %v, resp is %#v, ciphertext was %s, chosenKey is %s, id is %d", err, *resp, ct, chosenKey, id)
				}
				ptb64, ok := resp.Data["plaintext"].(string)
				if !ok {
					t.Errorf("no plaintext found, response was %#v", *resp)
					return
				}
				pt, err := base64.StdEncoding.DecodeString(ptb64)
				if err != nil {
					t.Errorf("got an error decoding base64 plaintext: %v", err)
					return
				}
				if string(pt) != testPlaintext {
					t.Errorf("got bad plaintext back: %s", pt)
				}

			// Change the min version, which also tests the archive functionality
			case "change_min_version":
				// t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				resp, err := be.pathPolicyRead(context.Background(), req, fd)
				if err != nil {
					t.Errorf("got an error reading policy %s: %v", chosenKey, err)
				}
				latestVersion := resp.Data["latest_version"].(int)

				// keys start at version 1 so we want [1, latestVersion] not [0, latestVersion)
				setVersion := (rand.Int() % latestVersion) + 1
				fd.Raw["min_decryption_version"] = setVersion
				fd.Schema = be.pathKeysConfig().Fields
				resp, err = be.pathKeysConfigWrite(context.Background(), req, fd)
				if err != nil {
					t.Errorf("got an error setting min decryption version: %v", err)
				}
			}
		}
	}

	// Spawn 1000 of these workers for 10 seconds
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go doFuzzy(i)
	}

	// Wait for them all to finish
	wg.Wait()
}

func TestBadInput(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/test",
	}

	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, resp, "expected populated request")

	req.Path = "decrypt/test"
	req.Data = map[string]interface{}{
		"ciphertext": "vault:v1:abcd",
	}

	_, err = b.HandleRequest(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTransit_AutoRotateKeys(t *testing.T) {
	tests := map[string]struct {
		isDRSecondary   bool
		isPerfSecondary bool
		isStandby       bool
		isLocal         bool
		shouldRotate    bool
	}{
		"primary, no local mount": {
			shouldRotate: true,
		},
		"DR secondary, no local mount": {
			isDRSecondary: true,
			shouldRotate:  false,
		},
		"perf standby, no local mount": {
			isStandby:    true,
			shouldRotate: false,
		},
		"perf secondary, no local mount": {
			isPerfSecondary: true,
			shouldRotate:    false,
		},
		"perf secondary, local mount": {
			isPerfSecondary: true,
			isLocal:         true,
			shouldRotate:    true,
		},
	}

	for name, test := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				var repState consts.ReplicationState
				if test.isDRSecondary {
					repState.AddState(consts.ReplicationDRSecondary)
				}
				if test.isPerfSecondary {
					repState.AddState(consts.ReplicationPerformanceSecondary)
				}
				if test.isStandby {
					repState.AddState(consts.ReplicationPerformanceStandby)
				}

				sysView := logical.TestSystemView()
				sysView.ReplicationStateVal = repState
				sysView.LocalMountVal = test.isLocal

				storage := &logical.InmemStorage{}

				conf := &logical.BackendConfig{
					StorageView: storage,
					System:      sysView,
				}

				b, _ := Backend(context.Background(), conf)
				if b == nil {
					t.Fatal("failed to create backend")
				}

				err := b.Backend.Setup(context.Background(), conf)
				if err != nil {
					t.Fatal(err)
				}

				// Write a key with the default auto rotate value (0/disabled)
				req := &logical.Request{
					Storage:   storage,
					Operation: logical.UpdateOperation,
					Path:      "keys/test1",
				}
				resp, err := b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				require.NotNil(t, resp, "expected populated request")

				// Write a key with an auto rotate value one day in the future
				req = &logical.Request{
					Storage:   storage,
					Operation: logical.UpdateOperation,
					Path:      "keys/test2",
					Data: map[string]interface{}{
						"auto_rotate_period": 24 * time.Hour,
					},
				}
				resp, err = b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				require.NotNil(t, resp, "expected populated request")

				// Run the rotation check and ensure none of the keys have rotated
				b.checkAutoRotateAfter = time.Now()
				if err = b.autoRotateKeys(context.Background(), &logical.Request{Storage: storage}); err != nil {
					t.Fatal(err)
				}
				req = &logical.Request{
					Storage:   storage,
					Operation: logical.ReadOperation,
					Path:      "keys/test1",
				}
				resp, err = b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Data["latest_version"] != 1 {
					t.Fatalf("incorrect latest_version found, got: %d, want: %d", resp.Data["latest_version"], 1)
				}

				req.Path = "keys/test2"
				resp, err = b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Data["latest_version"] != 1 {
					t.Fatalf("incorrect latest_version found, got: %d, want: %d", resp.Data["latest_version"], 1)
				}

				// Update auto rotate period on one key to be one nanosecond
				p, _, err := b.GetPolicy(context.Background(), keysutil.PolicyRequest{
					Storage: storage,
					Name:    "test2",
				}, b.GetRandomReader())
				if err != nil {
					t.Fatal(err)
				}
				if p == nil {
					t.Fatal("expected non-nil policy")
				}
				p.AutoRotatePeriod = time.Nanosecond
				err = p.Persist(context.Background(), storage)
				if err != nil {
					t.Fatal(err)
				}

				// Run the rotation check and validate the state of key rotations
				b.checkAutoRotateAfter = time.Now()
				if err = b.autoRotateKeys(context.Background(), &logical.Request{Storage: storage}); err != nil {
					t.Fatal(err)
				}
				req = &logical.Request{
					Storage:   storage,
					Operation: logical.ReadOperation,
					Path:      "keys/test1",
				}
				resp, err = b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Data["latest_version"] != 1 {
					t.Fatalf("incorrect latest_version found, got: %d, want: %d", resp.Data["latest_version"], 1)
				}
				req.Path = "keys/test2"
				resp, err = b.HandleRequest(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				expectedVersion := 1
				if test.shouldRotate {
					expectedVersion = 2
				}
				if resp.Data["latest_version"] != expectedVersion {
					t.Fatalf("incorrect latest_version found, got: %d, want: %d", resp.Data["latest_version"], expectedVersion)
				}
			},
		)
	}
}

func TestTransit_AEAD(t *testing.T) {
	testTransit_AEAD(t, "aes128-gcm96")
	testTransit_AEAD(t, "aes256-gcm96")
	testTransit_AEAD(t, "chacha20-poly1305")
}

func testTransit_AEAD(t *testing.T, keyType string) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	keyReq := &logical.Request{
		Path:      "keys/aead",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": keyType,
		},
		Storage: storage,
	}

	resp, err = b.HandleRequest(context.Background(), keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA=="                          // "the quick brown fox"
	associated := "U3BoaW54IG9mIGJsYWNrIHF1YXJ0eiwganVkZ2UgbXkgdm93Lgo=" // "Sphinx of black quartz, judge my vow."

	// Basic encrypt/decrypt should work.
	encryptReq := &logical.Request{
		Path:      "encrypt/aead",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"plaintext": plaintext,
		},
	}

	resp, err = b.HandleRequest(context.Background(), encryptReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	ciphertext1 := resp.Data["ciphertext"].(string)

	decryptReq := &logical.Request{
		Path:      "decrypt/aead",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"ciphertext": ciphertext1,
		},
	}

	resp, err = b.HandleRequest(context.Background(), decryptReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	decryptedPlaintext := resp.Data["plaintext"]

	if plaintext != decryptedPlaintext {
		t.Fatalf("bad: plaintext; expected: %q\nactual: %q", plaintext, decryptedPlaintext)
	}

	// Using associated as ciphertext should fail.
	decryptReq.Data["ciphertext"] = associated
	resp, err = b.HandleRequest(context.Background(), decryptReq)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("bad expected error: err: %v\nresp: %#v", err, resp)
	}

	// Redoing the above with additional data should work.
	encryptReq.Data["associated_data"] = associated
	resp, err = b.HandleRequest(context.Background(), encryptReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	ciphertext2 := resp.Data["ciphertext"].(string)
	decryptReq.Data["ciphertext"] = ciphertext2
	decryptReq.Data["associated_data"] = associated

	resp, err = b.HandleRequest(context.Background(), decryptReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	decryptedPlaintext = resp.Data["plaintext"]
	if plaintext != decryptedPlaintext {
		t.Fatalf("bad: plaintext; expected: %q\nactual: %q", plaintext, decryptedPlaintext)
	}

	// Removing the associated_data should break the decryption.
	decryptReq.Data = map[string]interface{}{
		"ciphertext": ciphertext2,
	}
	resp, err = b.HandleRequest(context.Background(), decryptReq)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("bad expected error: err: %v\nresp: %#v", err, resp)
	}

	// Using a valid ciphertext with associated_data should also break the
	// decryption.
	decryptReq.Data["ciphertext"] = ciphertext1
	decryptReq.Data["associated_data"] = associated
	resp, err = b.HandleRequest(context.Background(), decryptReq)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("bad expected error: err: %v\nresp: %#v", err, resp)
	}
}

// Hack: use Transit as a signer.
type transitKey struct {
	public any
	mount  string
	name   string
	t      *testing.T
	client *api.Client
}

func (k *transitKey) Public() crypto.PublicKey {
	return k.public
}

func (k *transitKey) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	hash := opts.(crypto.Hash)
	if hash.String() != "SHA-256" {
		return nil, fmt.Errorf("unknown hash algorithm: %v", opts)
	}

	resp, err := k.client.Logical().Write(k.mount+"/sign/"+k.name, map[string]interface{}{
		"hash_algorithm":      "sha2-256",
		"input":               base64.StdEncoding.EncodeToString(digest),
		"prehashed":           true,
		"signature_algorithm": "pkcs1v15",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}
	require.NotNil(k.t, resp)
	require.NotNil(k.t, resp.Data)
	require.NotNil(k.t, resp.Data["signature"])
	rawSig := resp.Data["signature"].(string)
	sigParts := strings.Split(rawSig, ":")

	decoded, err := base64.StdEncoding.DecodeString(sigParts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature (%v): %w", rawSig, err)
	}

	return decoded, nil
}

func TestTransitPKICSR(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": Factory,
			"pki":     pki.Factory,
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

	// Mount transit, write a key.
	err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("transit/keys/leaf", map[string]interface{}{
		"type": "rsa-2048",
	})
	require.NoError(t, err)

	resp, err := client.Logical().Read("transit/keys/leaf")
	require.NoError(t, err)
	require.NotNil(t, resp)

	keys := resp.Data["keys"].(map[string]interface{})
	require.NotNil(t, keys)
	keyData := keys["1"].(map[string]interface{})
	require.NotNil(t, keyData)
	keyPublic := keyData["public_key"].(string)
	require.NotEmpty(t, keyPublic)

	pemBlock, _ := pem.Decode([]byte(keyPublic))
	require.NotNil(t, pemBlock)
	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	require.NoError(t, err)
	require.NotNil(t, pubKey)

	// Setup a new CSR...
	var reqTemplate x509.CertificateRequest
	reqTemplate.PublicKey = pubKey
	reqTemplate.PublicKeyAlgorithm = x509.RSA
	reqTemplate.Subject.CommonName = "dadgarcorp.com"

	var k transitKey
	k.public = pubKey
	k.mount = "transit"
	k.name = "leaf"
	k.t = t
	k.client = client

	req, err := x509.CreateCertificateRequest(cryptoRand.Reader, &reqTemplate, &k)
	require.NoError(t, err)
	require.NotNil(t, req)

	reqPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: req,
	})
	t.Logf("csr: %v", string(reqPEM))

	// Mount PKI, generate a root, sign this CSR.
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	require.NoError(t, err)

	resp, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "PKI Root X1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	rootCertPEM := resp.Data["certificate"].(string)

	pemBlock, _ = pem.Decode([]byte(rootCertPEM))
	require.NotNil(t, pemBlock)

	rootCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)

	resp, err = client.Logical().Write("pki/issuer/default/sign-verbatim", map[string]interface{}{
		"csr": string(reqPEM),
		"ttl": "10m",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	leafCertPEM := resp.Data["certificate"].(string)
	pemBlock, _ = pem.Decode([]byte(leafCertPEM))
	require.NotNil(t, pemBlock)

	leafCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)
	require.NoError(t, leafCert.CheckSignatureFrom(rootCert))
	t.Logf("root: %v", rootCertPEM)
	t.Logf("leaf: %v", leafCertPEM)
}

func TestTransit_ReadPublicKeyImported(t *testing.T) {
	testTransit_ReadPublicKeyImported(t, "rsa-2048")
	testTransit_ReadPublicKeyImported(t, "ecdsa-p256")
	testTransit_ReadPublicKeyImported(t, "ed25519")
}

func testTransit_ReadPublicKeyImported(t *testing.T, keyType string) {
	generateKeys(t)
	b, s := createBackendWithStorage(t)
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate key ID: %s", err)
	}

	// Get key
	privateKey := getKey(t, keyType)
	publicKeyBytes, err := getPublicKey(privateKey, keyType)
	if err != nil {
		t.Fatalf("failed to extract the public key: %s", err)
	}

	// Import key
	importReq := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/import", keyID),
		Data: map[string]interface{}{
			"public_key": publicKeyBytes,
			"type":       keyType,
		},
	}
	importResp, err := b.HandleRequest(context.Background(), importReq)
	if err != nil || (importResp != nil && importResp.IsError()) {
		t.Fatalf("failed to import public key. err: %s\nresp: %#v", err, importResp)
	}

	// Read key
	readReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "keys/" + keyID,
		Storage:   s,
	}

	readResp, err := b.HandleRequest(context.Background(), readReq)
	if err != nil || (readResp != nil && readResp.IsError()) {
		t.Fatalf("failed to read key. err: %s\nresp: %#v", err, readResp)
	}
}

func TestTransit_SignWithImportedPublicKey(t *testing.T) {
	testTransit_SignWithImportedPublicKey(t, "rsa-2048")
	testTransit_SignWithImportedPublicKey(t, "ecdsa-p256")
	testTransit_SignWithImportedPublicKey(t, "ed25519")
}

func testTransit_SignWithImportedPublicKey(t *testing.T, keyType string) {
	generateKeys(t)
	b, s := createBackendWithStorage(t)
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate key ID: %s", err)
	}

	// Get key
	privateKey := getKey(t, keyType)
	publicKeyBytes, err := getPublicKey(privateKey, keyType)
	if err != nil {
		t.Fatalf("failed to extract the public key: %s", err)
	}

	// Import key
	importReq := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/import", keyID),
		Data: map[string]interface{}{
			"public_key": publicKeyBytes,
			"type":       keyType,
		},
	}
	importResp, err := b.HandleRequest(context.Background(), importReq)
	if err != nil || (importResp != nil && importResp.IsError()) {
		t.Fatalf("failed to import public key. err: %s\nresp: %#v", err, importResp)
	}

	// Sign text
	signReq := &logical.Request{
		Path:      "sign/" + keyID,
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(testPlaintext)),
		},
	}

	_, err = b.HandleRequest(context.Background(), signReq)
	if err == nil {
		t.Fatalf("expected error, should have failed to sign input")
	}
}

func TestTransit_VerifyWithImportedPublicKey(t *testing.T) {
	generateKeys(t)
	keyType := "rsa-2048"
	b, s := createBackendWithStorage(t)
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate key ID: %s", err)
	}

	// Get key
	privateKey := getKey(t, keyType)
	publicKeyBytes, err := getPublicKey(privateKey, keyType)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve public wrapping key
	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}

	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	// generate ciphertext
	importBlob := wrapTargetKeyForImport(t, pubWrappingKey, privateKey, keyType, "SHA256")

	// Import private key
	importReq := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/import", keyID),
		Data: map[string]interface{}{
			"ciphertext": importBlob,
			"type":       keyType,
		},
	}
	importResp, err := b.HandleRequest(context.Background(), importReq)
	if err != nil || (importResp != nil && importResp.IsError()) {
		t.Fatalf("failed to import key. err: %s\nresp: %#v", err, importResp)
	}

	// Sign text
	signReq := &logical.Request{
		Storage:   s,
		Path:      "sign/" + keyID,
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(testPlaintext)),
		},
	}

	signResp, err := b.HandleRequest(context.Background(), signReq)
	if err != nil || (signResp != nil && signResp.IsError()) {
		t.Fatalf("failed to sign plaintext. err: %s\nresp: %#v", err, signResp)
	}

	// Get signature
	signature := signResp.Data["signature"].(string)

	// Import new key as public key
	importPubReq := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/import", "public-key-rsa"),
		Data: map[string]interface{}{
			"public_key": publicKeyBytes,
			"type":       keyType,
		},
	}
	importPubResp, err := b.HandleRequest(context.Background(), importPubReq)
	if err != nil || (importPubResp != nil && importPubResp.IsError()) {
		t.Fatalf("failed to import public key. err: %s\nresp: %#v", err, importPubResp)
	}

	// Verify signed text
	verifyReq := &logical.Request{
		Path:      "verify/public-key-rsa",
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"input":     base64.StdEncoding.EncodeToString([]byte(testPlaintext)),
			"signature": signature,
		},
	}

	verifyResp, err := b.HandleRequest(context.Background(), verifyReq)
	if err != nil || (importResp != nil && verifyResp.IsError()) {
		t.Fatalf("failed to verify signed data. err: %s\nresp: %#v", err, importResp)
	}
}

func TestTransit_ExportPublicKeyImported(t *testing.T) {
	testTransit_ExportPublicKeyImported(t, "rsa-2048")
	testTransit_ExportPublicKeyImported(t, "ecdsa-p256")
	testTransit_ExportPublicKeyImported(t, "ed25519")
}

func testTransit_ExportPublicKeyImported(t *testing.T, keyType string) {
	generateKeys(t)
	b, s := createBackendWithStorage(t)
	keyID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate key ID: %s", err)
	}

	// Get key
	privateKey := getKey(t, keyType)
	publicKeyBytes, err := getPublicKey(privateKey, keyType)
	if err != nil {
		t.Fatalf("failed to extract the public key: %s", err)
	}

	t.Logf("generated key: %v", string(publicKeyBytes))

	// Import key
	importReq := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/import", keyID),
		Data: map[string]interface{}{
			"public_key": publicKeyBytes,
			"type":       keyType,
			"exportable": true,
		},
	}
	importResp, err := b.HandleRequest(context.Background(), importReq)
	if err != nil || (importResp != nil && importResp.IsError()) {
		t.Fatalf("failed to import public key. err: %s\nresp: %#v", err, importResp)
	}

	t.Logf("importing key: %v", importResp)

	// Export key
	exportReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("export/public-key/%s/latest", keyID),
		Storage:   s,
	}

	exportResp, err := b.HandleRequest(context.Background(), exportReq)
	if err != nil || (exportResp != nil && exportResp.IsError()) {
		t.Fatalf("failed to export key. err: %v\nresp: %#v", err, exportResp)
	}

	t.Logf("exporting key: %v", exportResp)

	responseKeys, exist := exportResp.Data["keys"]
	if !exist {
		t.Fatal("expected response data to hold a 'keys' field")
	}

	exportedKeyBytes := responseKeys.(map[string]string)["1"]

	if keyType != "ed25519" {
		exportedKeyBlock, _ := pem.Decode([]byte(exportedKeyBytes))
		publicKeyBlock, _ := pem.Decode(publicKeyBytes)

		if !reflect.DeepEqual(publicKeyBlock.Bytes, exportedKeyBlock.Bytes) {
			t.Fatalf("exported key bytes should have matched with imported key for key type: %v\nexported: %v\nimported: %v", keyType, exportedKeyBlock.Bytes, publicKeyBlock.Bytes)
		}
	} else {
		exportedKey, err := base64.StdEncoding.DecodeString(exportedKeyBytes)
		if err != nil {
			t.Fatalf("error decoding exported key bytes (%v) to base64 for key type %v: %v", exportedKeyBytes, keyType, err)
		}

		publicKeyBlock, _ := pem.Decode(publicKeyBytes)
		publicKeyParsed, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			t.Fatalf("error decoding source key bytes (%v) from PKIX marshaling for key type %v: %v", publicKeyBlock.Bytes, keyType, err)
		}

		if !reflect.DeepEqual([]byte(publicKeyParsed.(ed25519.PublicKey)), exportedKey) {
			t.Fatalf("exported key bytes should have matched with imported key for key type: %v\nexported: %v\nimported: %v", keyType, exportedKey, publicKeyParsed)
		}
	}
}

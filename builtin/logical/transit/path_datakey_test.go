// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

// TestDataKeyWithPaddingScheme validates that we properly leverage padding scheme
// args for the returned keys
func TestDataKeyWithPaddingScheme(t *testing.T) {
	b, s := createBackendWithStorage(t)
	keyName := "test"
	createKeyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + keyName,
		Storage:   s,
		Data: map[string]interface{}{
			"type": "rsa-2048",
		},
	}

	resp, err := b.HandleRequest(context.Background(), createKeyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed key creation: err: %v resp: %#v", err, resp)
	}

	tests := []struct {
		Name                 string
		PaddingScheme        string
		DecryptPaddingScheme string
		ShouldFailToDecrypt  bool
	}{
		{"no-padding-scheme", "", "", false},
		{"oaep", "oaep", "oaep", false},
		{"pkcs1v15", "pkcs1v15", "pkcs1v15", false},
		{"mixed-should-fail", "pkcs1v15", "oaep", true},
		{"mixed-based-on-default-should-fail", "", "pkcs1v15", true},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			dataKeyReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "datakey/wrapped/" + keyName,
				Storage:   s,
				Data:      map[string]interface{}{},
			}
			if len(tc.PaddingScheme) > 0 {
				dataKeyReq.Data["padding_scheme"] = tc.PaddingScheme
			}

			resp, err = b.HandleRequest(context.Background(), dataKeyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("failed data key api: err: %v resp: %#v", err, resp)
			}
			require.NotNil(t, resp, "Got nil nil response")
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			err = mapstructure.Decode(resp.Data, &d)
			require.NoError(t, err, "failed decoding datakey api response")
			require.NotEmpty(t, d.Ciphertext, "ciphertext should not be empty")

			// Attempt to decrypt with data key with the same padding scheme
			decryptReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "decrypt/" + keyName,
				Storage:   s,
				Data: map[string]interface{}{
					"ciphertext": d.Ciphertext,
				},
			}
			if len(tc.DecryptPaddingScheme) > 0 {
				decryptReq.Data["padding_scheme"] = tc.DecryptPaddingScheme
			}

			resp, err = b.HandleRequest(context.Background(), decryptReq)
			if tc.ShouldFailToDecrypt {
				require.Error(t, err, "Should have failed decryption as padding schemes are mixed")
			} else {
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("failed to decrypt data key: err: %v resp: %#v", err, resp)
				}
			}
		})
	}
}

// TestDataKeyWithPaddingSchemeInvalidKeyType validates we fail when we specify a
// padding_scheme value on an invalid key type (non-RSA)
func TestDataKeyWithPaddingSchemeInvalidKeyType(t *testing.T) {
	b, s := createBackendWithStorage(t)
	keyName := "test"
	createKeyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + keyName,
		Storage:   s,
		Data:      map[string]interface{}{},
	}

	resp, err := b.HandleRequest(context.Background(), createKeyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed key creation: err: %v resp: %#v", err, resp)
	}

	dataKeyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "datakey/wrapped/" + keyName,
		Storage:   s,
		Data: map[string]interface{}{
			"padding_scheme": "oaep",
		},
	}

	resp, err = b.HandleRequest(context.Background(), dataKeyReq)
	require.ErrorContains(t, err, "invalid request")
	require.NotNil(t, resp, "response should not be nil")
	require.Contains(t, resp.Error().Error(), "padding_scheme argument invalid: unsupported key")
}

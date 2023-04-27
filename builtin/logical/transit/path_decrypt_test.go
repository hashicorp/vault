// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func TestTransit_BatchDecryption(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := []interface{}{
		map[string]interface{}{"plaintext": "", "reference": "foo"},     // empty string
		map[string]interface{}{"plaintext": "Cg==", "reference": "bar"}, // newline
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA==", "reference": "baz"},
	}
	batchEncryptionData := map[string]interface{}{
		"batch_input": batchEncryptionInput,
	}

	batchEncryptionReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "encrypt/upserted_key",
		Storage:   s,
		Data:      batchEncryptionData,
	}
	resp, err = b.HandleRequest(context.Background(), batchEncryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchResponseItems := resp.Data["batch_results"].([]EncryptBatchResponseItem)
	batchDecryptionInput := make([]interface{}, len(batchResponseItems))
	for i, item := range batchResponseItems {
		batchDecryptionInput[i] = map[string]interface{}{"ciphertext": item.Ciphertext, "reference": item.Reference}
	}
	batchDecryptionData := map[string]interface{}{
		"batch_input": batchDecryptionInput,
	}

	batchDecryptionReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/upserted_key",
		Storage:   s,
		Data:      batchDecryptionData,
	}
	resp, err = b.HandleRequest(context.Background(), batchDecryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchDecryptionResponseItems := resp.Data["batch_results"].([]DecryptBatchResponseItem)
	// This seems fragile
	expectedResult := "[{\"plaintext\":\"\",\"reference\":\"foo\"},{\"plaintext\":\"Cg==\",\"reference\":\"bar\"},{\"plaintext\":\"dGhlIHF1aWNrIGJyb3duIGZveA==\",\"reference\":\"baz\"}]"

	jsonResponse, err := json.Marshal(batchDecryptionResponseItems)
	if err != nil || err == nil && string(jsonResponse) != expectedResult {
		t.Fatalf("bad: expected json response [%s]", jsonResponse)
	}
}

func TestTransit_BatchDecryption_DerivedKey(t *testing.T) {
	var req *logical.Request
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	// Create a derived key.
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/existing_key",
		Storage:   s,
		Data: map[string]interface{}{
			"derived": true,
		},
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Encrypt some values for use in test cases.
	plaintextItems := []struct {
		plaintext, context string
	}{
		{plaintext: "dGhlIHF1aWNrIGJyb3duIGZveA==", context: "dGVzdGNvbnRleHQ="},
		{plaintext: "anVtcGVkIG92ZXIgdGhlIGxhenkgZG9n", context: "dGVzdGNvbnRleHQy"},
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/existing_key",
		Storage:   s,
		Data: map[string]interface{}{
			"batch_input": []interface{}{
				map[string]interface{}{"plaintext": plaintextItems[0].plaintext, "context": plaintextItems[0].context},
				map[string]interface{}{"plaintext": plaintextItems[1].plaintext, "context": plaintextItems[1].context},
			},
		},
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	encryptedItems := resp.Data["batch_results"].([]EncryptBatchResponseItem)

	tests := []struct {
		name           string
		in             []interface{}
		want           []DecryptBatchResponseItem
		shouldErr      bool
		wantHTTPStatus int
		params         map[string]interface{}
	}{
		{
			name:      "nil-input",
			in:        nil,
			shouldErr: true,
		},
		{
			name:      "empty-input",
			in:        []interface{}{},
			shouldErr: true,
		},
		{
			name: "single-item-success",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[0].context},
			},
			want: []DecryptBatchResponseItem{
				{Plaintext: plaintextItems[0].plaintext},
			},
		},
		{
			name: "single-item-invalid-ciphertext",
			in: []interface{}{
				map[string]interface{}{"ciphertext": "xxx", "context": plaintextItems[0].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "invalid ciphertext: no prefix"},
			},
			wantHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "single-item-wrong-context",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
			},
			wantHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "batch-full-success",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[0].context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[1].context},
			},
			want: []DecryptBatchResponseItem{
				{Plaintext: plaintextItems[0].plaintext},
				{Plaintext: plaintextItems[1].plaintext},
			},
		},
		{
			name: "batch-partial-success",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[1].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Plaintext: plaintextItems[1].plaintext},
			},
			wantHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "batch-partial-success-overridden-response",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[1].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Plaintext: plaintextItems[1].plaintext},
			},
			params:         map[string]interface{}{"partial_failure_response_code": http.StatusAccepted},
			wantHTTPStatus: http.StatusAccepted,
		},
		{
			name: "batch-full-failure",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[0].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Error: "cipher: message authentication failed"},
			},
			wantHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "batch-full-failure-overridden-response",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[0].context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Error: "cipher: message authentication failed"},
			},
			params: map[string]interface{}{"partial_failure_response_code": http.StatusAccepted},
			// Full failure, shouldn't affect status code
			wantHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req = &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "decrypt/existing_key",
				Storage:   s,
				Data: map[string]interface{}{
					"batch_input": tt.in,
				},
			}
			for k, v := range tt.params {
				req.Data[k] = v
			}
			resp, err = b.HandleRequest(context.Background(), req)

			didErr := err != nil || (resp != nil && resp.IsError())
			if didErr {
				if !tt.shouldErr {
					t.Fatalf("unexpected error err:%v, resp:%#v", err, resp)
				}
			} else {
				if tt.shouldErr {
					t.Fatal("expected error, but none occurred")
				}

				if rawRespBody, ok := resp.Data[logical.HTTPRawBody]; ok {
					httpResp := &logical.HTTPResponse{}
					err = jsonutil.DecodeJSON([]byte(rawRespBody.(string)), httpResp)
					if err != nil {
						t.Fatalf("failed to unmarshal nested response: err:%v, resp:%#v", err, resp)
					}

					if respStatus, ok := resp.Data[logical.HTTPStatusCode]; !ok || respStatus != tt.wantHTTPStatus {
						t.Fatalf("HTTP response status code mismatch, want:%d, got:%d", tt.wantHTTPStatus, respStatus)
					}

					resp = logical.HTTPResponseToLogicalResponse(httpResp)
				}

				var respItems []DecryptBatchResponseItem
				err = mapstructure.Decode(resp.Data["batch_results"], &respItems)
				if err != nil {
					t.Fatalf("problem decoding response items: err:%v, resp:%#v", err, resp)
				}
				if !reflect.DeepEqual(tt.want, respItems) {
					t.Fatalf("response items mismatch, want:%#v, got:%#v", tt.want, respItems)
				}
			}
		})
	}
}

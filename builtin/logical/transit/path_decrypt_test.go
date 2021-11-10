package transit

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestTransit_BatchDecryption(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := []interface{}{
		map[string]interface{}{"plaintext": ""},     // empty string
		map[string]interface{}{"plaintext": "Cg=="}, // newline
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
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
		batchDecryptionInput[i] = map[string]interface{}{"ciphertext": item.Ciphertext}
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
	expectedResult := "[{\"plaintext\":\"\"},{\"plaintext\":\"Cg==\"},{\"plaintext\":\"dGhlIHF1aWNrIGJyb3duIGZveA==\"}]"

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
	plaintextItems := []BatchRequestItem{
		{Plaintext: "dGhlIHF1aWNrIGJyb3duIGZveA==", Context: "dGVzdGNvbnRleHQ="},
		{Plaintext: "anVtcGVkIG92ZXIgdGhlIGxhenkgZG9n", Context: "dGVzdGNvbnRleHQy"},
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/existing_key",
		Storage:   s,
		Data: map[string]interface{}{
			"batch_input": []interface{}{
				map[string]interface{}{"plaintext": plaintextItems[0].Plaintext, "context": plaintextItems[0].Context},
				map[string]interface{}{"plaintext": plaintextItems[1].Plaintext, "context": plaintextItems[1].Context},
			},
		},
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	encryptedItems := resp.Data["batch_results"].([]EncryptBatchResponseItem)

	tests := []struct {
		name      string
		in        []interface{}
		want      []DecryptBatchResponseItem
		shouldErr bool
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
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[0].Context},
			},
			want: []DecryptBatchResponseItem{
				{Plaintext: plaintextItems[0].Plaintext},
			},
		},
		{
			name: "single-item-invalid-ciphertext",
			in: []interface{}{
				map[string]interface{}{"ciphertext": "xxx", "context": plaintextItems[0].Context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "invalid ciphertext: no prefix"},
			},
		},
		{
			name: "single-item-wrong-context",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].Context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
			},
		},
		{
			name: "batch-full-success",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[0].Context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[1].Context},
			},
			want: []DecryptBatchResponseItem{
				{Plaintext: plaintextItems[0].Plaintext},
				{Plaintext: plaintextItems[1].Plaintext},
			},
		},
		{
			name: "batch-partial-success",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].Context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[1].Context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Plaintext: plaintextItems[1].Plaintext},
			},
		},
		{
			name: "batch-full-failure",
			in: []interface{}{
				map[string]interface{}{"ciphertext": encryptedItems[0].Ciphertext, "context": plaintextItems[1].Context},
				map[string]interface{}{"ciphertext": encryptedItems[1].Ciphertext, "context": plaintextItems[0].Context},
			},
			want: []DecryptBatchResponseItem{
				{Error: "cipher: message authentication failed"},
				{Error: "cipher: message authentication failed"},
			},
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

				respItems := resp.Data["batch_results"].([]DecryptBatchResponseItem)

				if !reflect.DeepEqual(tt.want, respItems) {
					t.Fatalf("response items mismatch, want:%#v, got:%#v", tt.want, respItems)
				}
			}
		})
	}
}

package transit

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

// Case1: If batch decryption input is not base64 encoded, it should fail.
func TestTransit_BatchDecryptionCase1(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := []interface{}{
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
		map[string]interface{}{"plaintext": "Cg=="},
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
	resp, err = b.HandleRequest(batchEncryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchDecryptionData := map[string]interface{}{
		"batch_input": resp.Data["batch_results"],
	}

	batchDecryptionReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/upserted_key",
		Storage:   s,
		Data:      batchDecryptionData,
	}
	resp, err = b.HandleRequest(batchDecryptionReq)
	if err == nil {
		t.Fatalf("expected an error")
	}
}

// Case2: Normal case of batch decryption
func TestTransit_BatchDecryptionCase2(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := []interface{}{
		map[string]interface{}{"plaintext": "Cg=="},
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
	resp, err = b.HandleRequest(batchEncryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchResponseItems := resp.Data["batch_results"].([]BatchResponseItem)
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
	resp, err = b.HandleRequest(batchDecryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchDecryptionResponseItems := resp.Data["batch_results"].([]BatchResponseItem)

	plaintext1 := "dGhlIHF1aWNrIGJyb3duIGZveA=="
	plaintext2 := "Cg=="
	for _, item := range batchDecryptionResponseItems {
		if item.Plaintext != plaintext1 && item.Plaintext != plaintext2 {
			t.Fatalf("bad: plaintext: %q", item.Plaintext)
		}
	}
}

// Case3: Test batch decryption with a derived key
func TestTransit_BatchDecryptionCase3(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	policyData := map[string]interface{}{
		"derived": true,
	}

	policyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/existing_key",
		Storage:   s,
		Data:      policyData,
	}

	resp, err = b.HandleRequest(policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchInput := []interface{}{
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA==", "context": "dGVzdGNvbnRleHQ="},
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA==", "context": "dGVzdGNvbnRleHQ="},
	}

	batchData := map[string]interface{}{
		"batch_input": batchInput,
	}
	batchReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "encrypt/existing_key",
		Storage:   s,
		Data:      batchData,
	}
	resp, err = b.HandleRequest(batchReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchDecryptionInputItems := resp.Data["batch_results"].([]BatchResponseItem)

	batchDecryptionInput := make([]interface{}, len(batchDecryptionInputItems))
	for i, item := range batchDecryptionInputItems {
		batchDecryptionInput[i] = map[string]interface{}{"ciphertext": item.Ciphertext, "context": "dGVzdGNvbnRleHQ="}
	}

	batchDecryptionData := map[string]interface{}{
		"batch_input": batchDecryptionInput,
	}

	batchDecryptionReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/existing_key",
		Storage:   s,
		Data:      batchDecryptionData,
	}
	resp, err = b.HandleRequest(batchDecryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchDecryptionResponseItems := resp.Data["batch_results"].([]BatchResponseItem)

	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA=="
	for _, item := range batchDecryptionResponseItems {
		if item.Plaintext != plaintext {
			t.Fatalf("bad: plaintext. Expected: %q, Actual: %q", plaintext, item.Plaintext)
		}
	}
}

package transit

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

// Case1: If batch decryption input is not base64 encoded, it should fail.
func TestTransit_BatchDecryptionCase1(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := `[{"plaintext":"dGhlIHF1aWNrIGJyb3duIGZveA=="},{"plaintext":"Cg=="}]`
	batchEncryptionInputB64 := base64.StdEncoding.EncodeToString([]byte(batchEncryptionInput))
	batchEncryptionData := map[string]interface{}{
		"batch": batchEncryptionInputB64,
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

	batchDecryptionInput := resp.Data["data"].(string)
	batchDecryptionData := map[string]interface{}{
		"batch": batchDecryptionInput,
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

	batchEncryptionInput := `[{"plaintext":"dGhlIHF1aWNrIGJyb3duIGZveA=="},{"plaintext":"Cg=="}]`
	batchEncryptionInputB64 := base64.StdEncoding.EncodeToString([]byte(batchEncryptionInput))
	batchEncryptionData := map[string]interface{}{
		"batch": batchEncryptionInputB64,
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

	batchDecryptionInput := resp.Data["data"].(string)
	batchDecryptionInputB64 := base64.StdEncoding.EncodeToString([]byte(batchDecryptionInput))
	batchDecryptionData := map[string]interface{}{
		"batch": batchDecryptionInputB64,
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

	var batchDecryptionResponseArray []interface{}
	if err := jsonutil.DecodeJSON([]byte(resp.Data["data"].(string)), &batchDecryptionResponseArray); err != nil {
		t.Fatal(err)
	}

	plaintext1 := "dGhlIHF1aWNrIGJyb3duIGZveA=="
	plaintext2 := "Cg=="
	for _, responseItem := range batchDecryptionResponseArray {
		var item BatchDecryptionItemResponse
		if err := mapstructure.Decode(responseItem, &item); err != nil {
			t.Fatal(err)
		}
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

	batchInput := `[{"plaintext":"dGhlIHF1aWNrIGJyb3duIGZveA==",
"context":"dmlzaGFsCg=="},{"plaintext":"dGhlIHF1aWNrIGJyb3duIGZveA==",
"context":"dmlzaGFsCg=="}]`

	batchInputB64 := base64.StdEncoding.EncodeToString([]byte(batchInput))
	batchData := map[string]interface{}{
		"batch": batchInputB64,
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

	var decryptionRequestItems []BatchDecryptionItemRequest
	var batchResponseArray []interface{}
	if err := jsonutil.DecodeJSON([]byte(resp.Data["data"].(string)), &batchResponseArray); err != nil {
		t.Fatal(err)
	}
	for _, responseItem := range batchResponseArray {
		var item BatchDecryptionItemRequest
		if err := mapstructure.Decode(responseItem, &item); err != nil {
			t.Fatal(err)
		}
		item.Context = "dmlzaGFsCg=="
		decryptionRequestItems = append(decryptionRequestItems, item)
	}

	batchDecryptionInput, err := jsonutil.EncodeJSON(decryptionRequestItems)
	if err != nil {
		t.Fatalf("failed to encode batch decryption input")
	}

	batchDecryptionInputB64 := base64.StdEncoding.EncodeToString(batchDecryptionInput)
	batchDecryptionData := map[string]interface{}{
		"batch": batchDecryptionInputB64,
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

	var batchDecryptionResponseArray []interface{}
	if err := jsonutil.DecodeJSON([]byte(resp.Data["data"].(string)), &batchDecryptionResponseArray); err != nil {
		t.Fatal(err)
	}

	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA=="
	for _, responseItem := range batchDecryptionResponseArray {
		var item BatchDecryptionItemResponse
		if err := mapstructure.Decode(responseItem, &item); err != nil {
			t.Fatal(err)
		}
		if item.Plaintext != plaintext {
			t.Fatalf("bad: plaintext. Expected: %q, Actual: %q", plaintext, item.Plaintext)
		}
	}
}

package transit

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical"
)

// Check the normal flow of rewrap
func TestTransit_BatchRewrapCase1(t *testing.T) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Upsert the key and encrypt the data
	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	encData := map[string]interface{}{
		"plaintext": plaintext,
	}

	// Create a key and encrypt a plaintext
	encReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "encrypt/upserted_key",
		Storage:   s,
		Data:      encData,
	}
	resp, err = b.HandleRequest(context.Background(), encReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Cache the ciphertext
	ciphertext := resp.Data["ciphertext"]
	if !strings.HasPrefix(ciphertext.(string), "vault:v1") {
		t.Fatalf("bad: ciphertext version: expected: 'vault:v1', actual: %s", ciphertext)
	}

	rewrapData := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	// Read the policy and check if the latest version is 1
	policyReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "keys/upserted_key",
		Storage:   s,
	}

	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["latest_version"] != 1 {
		t.Fatalf("bad: latest_version: expected: 1, actual: %d", resp.Data["latest_version"])
	}

	rotateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/upserted_key/rotate",
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), rotateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Read the policy again and the latest version is 2
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["latest_version"] != 2 {
		t.Fatalf("bad: latest_version: expected: 2, actual: %d", resp.Data["latest_version"])
	}

	// Rewrap the ciphertext and check that they are different
	rewrapReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rewrap/upserted_key",
		Storage:   s,
		Data:      rewrapData,
	}

	resp, err = b.HandleRequest(context.Background(), rewrapReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if ciphertext.(string) == resp.Data["ciphertext"].(string) {
		t.Fatalf("bad: ciphertexts are same before and after rewrap")
	}

	if !strings.HasPrefix(resp.Data["ciphertext"].(string), "vault:v2") {
		t.Fatalf("bad: ciphertext version: expected: 'vault:v2', actual: %s", resp.Data["ciphertext"].(string))
	}
}

// Check the normal flow of rewrap with upserted key
func TestTransit_BatchRewrapCase2(t *testing.T) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Upsert the key and encrypt the data
	plaintext := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	encData := map[string]interface{}{
		"plaintext": plaintext,
		"context":   "dmlzaGFsCg==",
	}

	// Create a key and encrypt a plaintext
	encReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "encrypt/upserted_key",
		Storage:   s,
		Data:      encData,
	}
	resp, err = b.HandleRequest(context.Background(), encReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Cache the ciphertext
	ciphertext := resp.Data["ciphertext"]
	if !strings.HasPrefix(ciphertext.(string), "vault:v1") {
		t.Fatalf("bad: ciphertext version: expected: 'vault:v1', actual: %s", ciphertext)
	}

	rewrapData := map[string]interface{}{
		"ciphertext": ciphertext,
		"context":    "dmlzaGFsCg==",
	}

	// Read the policy and check if the latest version is 1
	policyReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "keys/upserted_key",
		Storage:   s,
	}

	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["latest_version"] != 1 {
		t.Fatalf("bad: latest_version: expected: 1, actual: %d", resp.Data["latest_version"])
	}

	rotateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/upserted_key/rotate",
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), rotateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	// Read the policy again and the latest version is 2
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if resp.Data["latest_version"] != 2 {
		t.Fatalf("bad: latest_version: expected: 2, actual: %d", resp.Data["latest_version"])
	}

	// Rewrap the ciphertext and check that they are different
	rewrapReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rewrap/upserted_key",
		Storage:   s,
		Data:      rewrapData,
	}

	resp, err = b.HandleRequest(context.Background(), rewrapReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	if ciphertext.(string) == resp.Data["ciphertext"].(string) {
		t.Fatalf("bad: ciphertexts are same before and after rewrap")
	}

	if !strings.HasPrefix(resp.Data["ciphertext"].(string), "vault:v2") {
		t.Fatalf("bad: ciphertext version: expected: 'vault:v2', actual: %s", resp.Data["ciphertext"].(string))
	}
}

// Batch encrypt plaintexts, rotate the keys and rewrap all the ciphertexts
func TestTransit_BatchRewrapCase3(t *testing.T) {
	var resp *logical.Response
	var err error

	b, s := createBackendWithStorage(t)

	batchEncryptionInput := []interface{}{
		map[string]interface{}{"plaintext": "dmlzaGFsCg=="},
		map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
	}
	batchEncryptionData := map[string]interface{}{
		"batch_input": batchEncryptionInput,
	}
	batchReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "encrypt/upserted_key",
		Storage:   s,
		Data:      batchEncryptionData,
	}
	resp, err = b.HandleRequest(context.Background(), batchReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchEncryptionResponseItems := resp.Data["batch_results"].([]BatchResponseItem)

	batchRewrapInput := make([]interface{}, len(batchEncryptionResponseItems))
	for i, item := range batchEncryptionResponseItems {
		batchRewrapInput[i] = map[string]interface{}{"ciphertext": item.Ciphertext}
	}

	batchRewrapData := map[string]interface{}{
		"batch_input": batchRewrapInput,
	}

	rotateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/upserted_key/rotate",
		Storage:   s,
	}
	resp, err = b.HandleRequest(context.Background(), rotateReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	rewrapReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "rewrap/upserted_key",
		Storage:   s,
		Data:      batchRewrapData,
	}

	resp, err = b.HandleRequest(context.Background(), rewrapReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%v resp:%#v", err, resp)
	}

	batchRewrapResponseItems := resp.Data["batch_results"].([]BatchResponseItem)

	if len(batchRewrapResponseItems) != len(batchEncryptionResponseItems) {
		t.Fatalf("bad: length of input and output or rewrap are not matching; expected: %d, actual: %d", len(batchEncryptionResponseItems), len(batchRewrapResponseItems))
	}

	decReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/upserted_key",
		Storage:   s,
	}

	for i, eItem := range batchEncryptionResponseItems {
		rItem := batchRewrapResponseItems[i]

		if eItem.Ciphertext == rItem.Ciphertext {
			t.Fatalf("bad: rewrap input and output are the same")
		}

		if !strings.HasPrefix(rItem.Ciphertext, "vault:v2") {
			t.Fatalf("bad: invalid version of ciphertext in rewrap response; expected: 'vault:v2', actual: %s", rItem.Ciphertext)
		}

		decReq.Data = map[string]interface{}{
			"ciphertext": rItem.Ciphertext,
		}

		resp, err = b.HandleRequest(context.Background(), decReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%v resp:%#v", err, resp)
		}

		plaintext1 := "dGhlIHF1aWNrIGJyb3duIGZveA=="
		plaintext2 := "dmlzaGFsCg=="
		if resp.Data["plaintext"] != plaintext1 && resp.Data["plaintext"] != plaintext2 {
			t.Fatalf("bad: plaintext. Expected: %q or %q, Actual: %q", plaintext1, plaintext2, resp.Data["plaintext"])
		}
	}
}

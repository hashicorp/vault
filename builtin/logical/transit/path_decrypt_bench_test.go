// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func BenchmarkTransit_BatchDecryption1(b *testing.B) {
	BTransit_BatchDecryption(b, 1)
}

func BenchmarkTransit_BatchDecryption10(b *testing.B) {
	BTransit_BatchDecryption(b, 10)
}

func BenchmarkTransit_BatchDecryption50(b *testing.B) {
	BTransit_BatchDecryption(b, 50)
}

func BenchmarkTransit_BatchDecryption100(b *testing.B) {
	BTransit_BatchDecryption(b, 100)
}

func BenchmarkTransit_BatchDecryption1000(b *testing.B) {
	BTransit_BatchDecryption(b, 1_000)
}

func BenchmarkTransit_BatchDecryption10000(b *testing.B) {
	BTransit_BatchDecryption(b, 10_000)
}

func BTransit_BatchDecryption(b *testing.B, bsize int) {
	b.StopTimer()

	var resp *logical.Response
	var err error

	backend, s := createBackendWithStorage(b)

	batchEncryptionInput := make([]interface{}, 0, bsize)
	for i := 0; i < bsize; i++ {
		batchEncryptionInput = append(
			batchEncryptionInput,
			map[string]interface{}{"plaintext": "dGhlIHF1aWNrIGJyb3duIGZveA=="},
		)
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
	resp, err = backend.HandleRequest(context.Background(), batchEncryptionReq)
	if err != nil || (resp != nil && resp.IsError()) {
		b.Fatalf("err:%v resp:%#v", err, resp)
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		resp, err = backend.HandleRequest(context.Background(), batchDecryptionReq)
		if err != nil || (resp != nil && resp.IsError()) {
			b.Fatalf("err:%v resp:%#v", err, resp)
		}
	}
}

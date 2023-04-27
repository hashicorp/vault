// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func BenchmarkTransit_BatchEncryption1(b *testing.B) {
	BTransit_BatchEncryption(b, 1)
}

func BenchmarkTransit_BatchEncryption10(b *testing.B) {
	BTransit_BatchEncryption(b, 10)
}

func BenchmarkTransit_BatchEncryption50(b *testing.B) {
	BTransit_BatchEncryption(b, 50)
}

func BenchmarkTransit_BatchEncryption100(b *testing.B) {
	BTransit_BatchEncryption(b, 100)
}

func BenchmarkTransit_BatchEncryption1000(b *testing.B) {
	BTransit_BatchEncryption(b, 1_000)
}

func BenchmarkTransit_BatchEncryption10000(b *testing.B) {
	BTransit_BatchEncryption(b, 10_000)
}

func BTransit_BatchEncryption(b *testing.B, bsize int) {
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		resp, err = backend.HandleRequest(context.Background(), batchEncryptionReq)
		if err != nil || (resp != nil && resp.IsError()) {
			b.Fatalf("err:%v resp:%#v", err, resp)
		}
	}
}

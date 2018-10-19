package transit

import (
	"testing"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestTransit_Trim(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	doReq := func(t *testing.T, req *logical.Request) *logical.Response {
		t.Helper()
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("got err:\n%#v\nresp:\n%#v\n", err, resp)
		}
		return resp
	}
	doErrReq := func(t *testing.T, req *logical.Request) {
		t.Helper()
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err == nil && (resp == nil || !resp.IsError()) {
			t.Fatalf("expected error; resp:\n%#v\n", resp)
		}
	}

	// Create a key
	req := &logical.Request{
		Path:      "keys/aes",
		Storage:   storage,
		Operation: logical.UpdateOperation,
	}
	doReq(t, req)

	// Get the policy and check that the archive has correct number of keys
	p, _, err := b.lm.GetPolicy(namespace.RootContext(nil), keysutil.PolicyRequest{
		Storage: storage,
		Name:    "aes",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Archive: 0, 1
	archive, err := p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	// Index "0" in the archive is unused. Hence the length of the archived
	// keys will always be 1 more than the actual number of keys.
	if len(archive.Keys) != 2 {
		t.Fatalf("bad: len of archived keys; expected: 2, actual: %d", len(archive.Keys))
	}

	// Ensure that there are 5 key versions, by rotating the key 4 times
	for i := 0; i < 4; i++ {
		req.Path = "keys/aes/rotate"
		req.Data = nil
		doReq(t, req)
	}

	// Archive: 0, 1, 2, 3, 4, 5
	archive, err = p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.Keys) != 6 {
		t.Fatalf("bad: len of archived keys; expected: 6, actual: %d", len(archive.Keys))
	}

	// Min available version should not be set when min_encryption_version is not
	// set
	req.Path = "keys/aes/trim"
	req.Data = map[string]interface{}{
		"min_available_version": 1,
	}
	doErrReq(t, req)

	// Set min_encryption_version to 4
	req.Path = "keys/aes/config"
	req.Data = map[string]interface{}{
		"min_encryption_version": 4,
	}
	doReq(t, req)

	// Set min_decryption_version to 3
	req.Data = map[string]interface{}{
		"min_decryption_version": 3,
	}
	doReq(t, req)

	// Min available version cannot be greater than min encryption version
	req.Path = "keys/aes/trim"
	req.Data = map[string]interface{}{
		"min_available_version": 5,
	}
	doErrReq(t, req)

	// Min available version cannot be greater than min decryption version
	req.Data["min_available_version"] = 4
	doErrReq(t, req)

	// Min available version cannot be negative
	req.Data["min_available_version"] = -1
	doErrReq(t, req)

	// Min available version should be positive
	req.Data["min_available_version"] = 0
	doErrReq(t, req)

	// Trim all keys before version 3. Index 0 and index 1 will be deleted from
	// archived keys.
	req.Data["min_available_version"] = 3
	doReq(t, req)

	// Archive: 3, 4, 5
	archive, err = p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.Keys) != 3 {
		t.Fatalf("bad: len of archived keys; expected: 3, actual: %d", len(archive.Keys))
	}

	// Min decryption version should not be less than min available version
	req.Path = "keys/aes/config"
	req.Data = map[string]interface{}{
		"min_decryption_version": 1,
	}
	doErrReq(t, req)

	// Min encryption version should not be less than min available version
	req.Data = map[string]interface{}{
		"min_encryption_version": 2,
	}
	doErrReq(t, req)

	// Rotate 5 more times
	for i := 0; i < 5; i++ {
		doReq(t, &logical.Request{
			Path:      "keys/aes/rotate",
			Storage:   storage,
			Operation: logical.UpdateOperation,
		})
	}

	// Archive: 3, 4, 5, 6, 7, 8, 9, 10
	archive, err = p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.Keys) != 8 {
		t.Fatalf("bad: len of archived keys; expected: 8, actual: %d", len(archive.Keys))
	}

	// Set min encryption version to 7
	req.Data = map[string]interface{}{
		"min_encryption_version": 7,
	}
	doReq(t, req)

	// Set min decryption version to 7
	req.Data = map[string]interface{}{
		"min_decryption_version": 7,
	}
	doReq(t, req)

	// Trim all versions before 7
	req.Path = "keys/aes/trim"
	req.Data = map[string]interface{}{
		"min_available_version": 7,
	}
	doReq(t, req)

	// Archive: 7, 8, 9, 10
	archive, err = p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.Keys) != 4 {
		t.Fatalf("bad: len of archived keys; expected: 4, actual: %d", len(archive.Keys))
	}

	// Read the key
	req.Path = "keys/aes"
	req.Operation = logical.ReadOperation
	resp := doReq(t, req)
	keys := resp.Data["keys"].(map[string]int64)
	if len(keys) != 4 {
		t.Fatalf("bad: number of keys; expected: 4, actual: %d", len(keys))
	}

	// Test if moving the min_encryption_version and min_decryption_versions
	// are working fine

	// Set min encryption version to 10
	req.Path = "keys/aes/config"
	req.Operation = logical.UpdateOperation
	req.Data = map[string]interface{}{
		"min_encryption_version": 10,
	}
	doReq(t, req)
	if p.MinEncryptionVersion != 10 {
		t.Fatalf("failed to set min encryption version")
	}

	// Set min decryption version to 9
	req.Data = map[string]interface{}{
		"min_decryption_version": 9,
	}
	doReq(t, req)
	if p.MinDecryptionVersion != 9 {
		t.Fatalf("failed to set min encryption version")
	}

	// Reduce the min decryption version to 8
	req.Data = map[string]interface{}{
		"min_decryption_version": 8,
	}
	doReq(t, req)
	if p.MinDecryptionVersion != 8 {
		t.Fatalf("failed to set min encryption version")
	}

	// Reduce the min encryption version to 8
	req.Data = map[string]interface{}{
		"min_encryption_version": 8,
	}
	doReq(t, req)
	if p.MinDecryptionVersion != 8 {
		t.Fatalf("failed to set min decryption version")
	}

	// Read the key to ensure that the keys are properly copied from the
	// archive into the policy
	req.Path = "keys/aes"
	req.Operation = logical.ReadOperation
	resp = doReq(t, req)
	keys = resp.Data["keys"].(map[string]int64)
	if len(keys) != 3 {
		t.Fatalf("bad: number of keys; expected: 3, actual: %d", len(keys))
	}

	// Ensure that archive has remained unchanged
	// Archive: 7, 8, 9, 10
	archive, err = p.LoadArchive(namespace.RootContext(nil), storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.Keys) != 4 {
		t.Fatalf("bad: len of archived keys; expected: 4, actual: %d", len(archive.Keys))
	}
}

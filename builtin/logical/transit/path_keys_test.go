package transit

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_PathKeys_ExportValidVersionsOnly(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	b = Backend(&logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	})

	// First create a key, v1
	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo",
	}
	req.Data = map[string]interface{}{
		"exportable": true,
	}
	_, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "keys/foo/rotate"
	// v2
	// valid versions: 1, 2
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	// v3
	// valid versions: 1, 2, 3
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	verifyExportCount := func(expectedCount int) {
		req = &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      "key-export/foo",
		}
		rsp, err := b.HandleRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := rsp.Data["keys"]; !ok {
			t.Error("no keys returned from export")
		}

		keys, ok := rsp.Data["keys"].(map[string]string)
		if !ok {
			t.Error("could not cast to keys object")
		}
		if len(keys) != expectedCount {
			t.Errorf("expected %d key, received %d", expectedCount, len(keys))
		}
	}

	verifyExportCount(3)

	// valid versions: 3
	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo/config",
	}
	req.Data = map[string]interface{}{
		"min_decryption_version": 3,
	}
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	verifyExportCount(1)

	// valid versions: 2, 3
	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo/config",
	}
	req.Data = map[string]interface{}{
		"min_decryption_version": 2,
	}
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	verifyExportCount(2)

	req = &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/foo/rotate",
	}
	// v4
	// valid versions: 2, 3, 4
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	verifyExportCount(3)
}

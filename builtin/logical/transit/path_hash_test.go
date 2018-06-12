package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestTransit_Hash(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "hash",
		Data: map[string]interface{}{
			"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
		},
	}

	doRequest := func(req *logical.Request, errExpected bool, expected string) {
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
			}
			return
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		sum, ok := resp.Data["sum"]
		if !ok {
			t.Fatal("no sum key found in returned data")
		}
		if sum.(string) != expected {
			t.Fatal("mismatched hashes")
		}
	}

	// Test defaults -- sha2-256
	doRequest(req, false, "9ecb36561341d18eb65484e833efea61edc74b84cf5e6ae1b81c63533e25fc8f")

	// Test algorithm selection in the path
	req.Path = "hash/sha2-224"
	doRequest(req, false, "ea074a96cabc5a61f8298a2c470f019074642631a49e1c5e2f560865")

	// Reset and test algorithm selection in the data
	req.Path = "hash"
	req.Data["algorithm"] = "sha2-224"
	doRequest(req, false, "ea074a96cabc5a61f8298a2c470f019074642631a49e1c5e2f560865")

	req.Data["algorithm"] = "sha2-384"
	doRequest(req, false, "15af9ec8be783f25c583626e9491dbf129dd6dd620466fdf05b3a1d0bb8381d30f4d3ec29f923ff1e09a0f6b337365a6")

	req.Data["algorithm"] = "sha2-512"
	doRequest(req, false, "d9d380f29b97ad6a1d92e987d83fa5a02653301e1006dd2bcd51afa59a9147e9caedaf89521abc0f0b682adcd47fb512b8343c834a32f326fe9bef00542ce887")

	// Test returning as base64
	req.Data["format"] = "base64"
	doRequest(req, false, "2dOA8puXrWodkumH2D+loCZTMB4QBt0rzVGvpZqRR+nK7a+JUhq8DwtoKtzUf7USuDQ8g0oy8yb+m+8AVCzohw==")

	// Test bad input/format/algorithm
	req.Data["format"] = "base92"
	doRequest(req, true, "")

	req.Data["format"] = "hex"
	req.Data["algorithm"] = "foobar"
	doRequest(req, true, "")

	req.Data["algorithm"] = "sha2-256"
	req.Data["input"] = "foobar"
	doRequest(req, true, "")
}

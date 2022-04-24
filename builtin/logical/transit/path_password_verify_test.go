package transit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestTransit_PasswordVerify(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "password/verify",
		Data: map[string]interface{}{
			"input": "EeYbAzzLsa0K3HXnPYun",
			"hash":  "$2a$10$.MSZmQXuxGmhWX.QS8C/yOFz9.buNE01J.dh1X5EE4c0WNhClNyIW",
		},
	}

	doRequest := func(req *logical.Request, errExpected bool) {
		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: did not get error response: %#v", *resp)
			}
			return
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		pass, ok := resp.Data["password"]
		if !ok {
			t.Fatal("no password key found in returned data")
		}
		password := pass.(string)
		if got, want := password, req.Data["hash"]; got != want {
			t.Fatalf("returned hash %s does not equal expected %s", got, want)
		}
	}

	// Test defaults
	doRequest(req, false)

	// Test wrong input in url
	req.Data["input"] = "1234"
	doRequest(req, true)

	// Test bad input
	delete(req.Data, "input")
	doRequest(req, true)

	// Test bad hash
	req.Data["input"] = "1234"
	delete(req.Data, "hash")
	doRequest(req, true)
}

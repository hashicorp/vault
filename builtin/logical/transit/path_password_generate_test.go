package transit

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/bcrypt"
)

func TestTransit_PasswordGenerate(t *testing.T) {
	b, storage := createBackendWithSysView(t)

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "password/generate",
		Data: map[string]interface{}{
			"input": "idCtBYdWLBkActnuj2k54",
		},
	}

	doRequest := func(req *logical.Request, expectedCost int, errExpected bool) {
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
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Data["input"].(string))); err != nil {
			t.Fatal("bcrypt comparison mismatch")
		}
		if !strings.HasPrefix(password, fmt.Sprintf("$2a$%02d$", expectedCost)) {
			t.Fatalf("%s does not have expected cost %d", password, expectedCost)
		}
	}

	// Test defaults
	doRequest(req, 10, false)

	// Test different cost in body
	req.Data["cost"] = 11
	doRequest(req, 11, false)
	req.Data["cost"] = -1
	doRequest(req, 10, false)

	// Test bad cost value
	req.Data["cost"] = 33
	doRequest(req, 0, true)
	req.Data["cost"] = "9a1d"
	doRequest(req, 0, true)

	// Test bad input
	delete(req.Data, "input")
	doRequest(req, 0, true)
}

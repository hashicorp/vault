package mfa

import (
	"testing"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func setupTestBackend() (*framework.Backend, error) {
	conf := logical.TestBackendConfig()
	conf.StorageView = &logical.InmemStorage{}
	b := makeTestBackend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

// makeTestBackend creates a simple MFA enabled backend.
// Login (before MFA) always succeeds with policy "foo".
// An MFA "test" type is added to mfa.handlers that succeeds
// if MFA method is "accept", otherwise it rejects.
func makeTestBackend() *framework.Backend {
	globalHandlers["test"] = testMFAHandler
	b := &framework.Backend{
		Help: "",

		PathsSpecial: &logical.Paths{
			Root: MFARootPaths(),
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: MFAPaths(nil, testPathLogin()),
	}
	return b
}

func testPathLogin() *framework.Path {
	return &framework.Path{
		Pattern: `login`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type: framework.TypeString,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: testPathLoginHandler,
		},
	}
}

func testPathLoginHandler(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: []string{"foo"},
			Metadata: map[string]string{
				"username": username,
			},
		},
	}, nil
}

func testMFAHandler(req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {
	if d.Get("method").(string) != "accept" {
		return logical.ErrorResponse("Deny access"), nil
	} else {
		return resp, nil
	}
}

func TestMFALogin(t *testing.T) {
	b, err := setupTestBackend()
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepEnableMFA(t),
			testAccStepLogin(t, "user"),
		},
	})
}

func TestMFALoginDenied(t *testing.T) {
	b, err := setupTestBackend()
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepEnableMFA(t),
			testAccStepLoginDenied(t, "user"),
		},
	})
}

func testAccStepEnableMFA(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "mfa_config",
		Data: map[string]interface{}{
			"type": "test",
		},
	}
}

func testAccStepLogin(t *testing.T, username string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"method":   "accept",
			"username": username,
		},
		Unauthenticated: true,
		Check:           logicaltest.TestCheckAuth([]string{"default", "foo"}),
	}
}

func testAccStepLoginDenied(t *testing.T, username string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"method":   "deny",
			"username": username,
		},
		ErrorOk:         true,
		Unauthenticated: true,
		Check:           logicaltest.TestCheckError(),
	}
}

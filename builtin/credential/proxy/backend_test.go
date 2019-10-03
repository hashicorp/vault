package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

type testBackend struct {
	Backend logical.Backend
	Storage logical.Storage
	testing *testing.T
}

func newTestBackend(t *testing.T) *testBackend {
	defaultLeaseTTLVal := time.Hour * 12
	maxLeaseTTLVal := time.Hour * 24

	storage := &logical.InmemStorage{}
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Trace),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
		StorageView: storage,
	})

	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	return &testBackend{
		Backend: b,
		Storage: storage,
		testing: t,
	}
}

func (t *testBackend) HandleRequest(req *logical.Request) (resp *logical.Response, err error) {
	t.testing.Helper()
	req.Storage = t.Storage
	return t.Backend.HandleRequest(context.Background(), req)
}

func (t *testBackend) AssertHandleRequest(req *logical.Request) (resp *logical.Response) {
	t.testing.Helper()
	resp, err := t.HandleRequest(req)
	assertSuccess(t.testing, err, resp)
	return resp
}

func assertSuccess(t *testing.T, err error, resp *logical.Response) {
	t.Helper()
	if err != nil {
		t.Fatalf("request failed. err=%+v\n", err)
	}

	if resp.IsError() {
		t.Fatalf("erroneous response:\n%+v", resp)
	}
}

func createConfigRequest(data map[string]interface{}) *logical.Request {
	return &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config",
		Data:      data,
	}
}

func updateConfigRequest(data map[string]interface{}) *logical.Request {
	req := createConfigRequest(data)
	req.Operation = logical.UpdateOperation
	return req
}

func readConfigRequest() *logical.Request {
	return &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
	}
}

func createRoleRequest(roleName string, data map[string]interface{}) *logical.Request {
	return &logical.Request{
		Operation: logical.CreateOperation,
		Path:      fmt.Sprintf("role/%s", roleName),
		Data:      data,
	}
}

func updateRoleRequest(roleName string, data map[string]interface{}) *logical.Request {
	req := createRoleRequest(roleName, data)
	req.Operation = logical.UpdateOperation
	return req
}

func readRoleRequest(roleName string) *logical.Request {
	return &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("role/%s", roleName),
	}
}

func listRoleRequest() *logical.Request {
	return &logical.Request{
		Operation: logical.ListOperation,
		Path:      "role/",
	}
}

func loginRequest(username, role, userHeader string, headers map[string][]string) *logical.Request {
	if headers == nil {
		headers = make(map[string][]string)
	}
	headers[userHeader] = []string{username}

	return &logical.Request{
		Operation:       logical.UpdateOperation,
		Unauthenticated: true,
		Path:            "login",
		Data:            map[string]interface{}{"role": role},
		Headers:         headers,
	}
}

// assertSerializedEqual checks to see if 'expected' and 'actual' have the
// same representations when serialized to JSON
func assertSerializedEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	expJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal failure: %+v\n", err)
	}

	var expObj interface{}
	if err := json.Unmarshal(expJSON, &expObj); err != nil {
		t.Fatalf("unmarshal failure: %+v\n", err)
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("marshal failure: %+v\n", err)
	}

	var actualObj interface{}
	if err := json.Unmarshal(actualJSON, &actualObj); err != nil {
		t.Fatalf("unmarshal failure: %+v\n", err)
	}
	if !reflect.DeepEqual(actualObj, expObj) {
		t.Fatalf("Unexpected value\ngot: %+v\nexp: %+v", actualObj, expObj)
	}
}

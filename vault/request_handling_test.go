package vault

import (
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func TestRequestHandling_Wrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["kv"] = PassthroughBackendFactory

	meUUID, _ := uuid.GenerateUUID()
	err := core.mount(namespace.RootContext(nil), &MountEntry{
		Table: mountTableType,
		UUID:  meUUID,
		Path:  "wraptest",
		Type:  "kv",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// No duration specified
	req := &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"zip": "zap",
		},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.ReadOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestRequestHandling_LoginWrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["userpass"] = credUserpass.Factory

	// No duration specified
	req := &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req.Path = "auth/userpass/users/test"
	req.Data = map[string]interface{}{
		"password": "foo",
		"policies": "default",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

package proxy

import (
	"testing"
)

func TestRoleReadWrite(t *testing.T) {
	b := newTestBackend(t)

	data := map[string]interface{}{
		"allowed_users": []interface{}{"foo1,foo2"},
		"required_headers": map[string]interface{}{
			"Hdr1": "value1",
			"Hdr2": "value2",
		},
		"policies": []interface{}{"policy1", "policy2"},
		"ttl":      20,
		"max_ttl":  40,
		"period":   7200,
	}
	req := createRoleRequest("foo", data)
	b.AssertHandleRequest(req)

	req = readRoleRequest("foo")
	resp := b.AssertHandleRequest(req)

	assertSerializedEqual(t, data, resp.Data)
}

func TestRoleValidation(t *testing.T) {
	b := newTestBackend(t)

	// missing allowed_users
	req := createRoleRequest("role1", map[string]interface{}{})
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %+v\n", err)
	}

	if !resp.IsError() {
		t.Fatalf("did not get error when required field not set")
	}
}

func TestRoleList(t *testing.T) {
	b := newTestBackend(t)

	data := map[string]interface{}{
		"allowed_users": []interface{}{"foo1,foo2"},
	}
	req := createRoleRequest("role1", data)
	b.AssertHandleRequest(req)

	req = createRoleRequest("role2", data)
	b.AssertHandleRequest(req)

	req = listRoleRequest()
	resp := b.AssertHandleRequest(req)

	exp := map[string][]string{"keys": []string{"role1", "role2"}}
	assertSerializedEqual(t, exp, resp.Data)
}

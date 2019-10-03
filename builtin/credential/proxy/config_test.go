package proxy

import (
	"testing"
)

func TestConfigReadWrite(t *testing.T) {
	b := newTestBackend(t)

	data := map[string]interface{}{
		"user_header": "Foobar",
		"bound_cidrs": []interface{}{"1.2.3.4/8", "5.6.7.8/24"},
	}
	req := createConfigRequest(data)
	b.AssertHandleRequest(req)

	req = readConfigRequest()
	resp := b.AssertHandleRequest(req)

	assertSerializedEqual(t, data, resp.Data)
}

func TestConfigValidation(t *testing.T) {
	b := newTestBackend(t)

	// missing required user header
	req := createConfigRequest(map[string]interface{}{})
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %+v\n", err)
	}

	if !resp.IsError() {
		t.Fatalf("did not get error when required field not set")
	}

	// invalid cidr
	data := map[string]interface{}{
		"user_header": "Foobar",
		"bound_cidrs": "monkey",
	}
	req = createConfigRequest(data)
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %+v\n", err)
	}

	if !resp.IsError() {
		t.Fatalf("did not get error when invalid CIDR provided")
	}
}

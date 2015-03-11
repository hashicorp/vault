package vault

import (
	"reflect"
	"testing"
)

func testSystem(t *testing.T) *SystemBackend {
	c, _ := testUnsealedCore(t)
	return &SystemBackend{c}
}

func TestSystem_mounts(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: ReadOperation,
		Path:      "mounts",
	}
	resp, err := s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"secret/": map[string]string{
			"type":        "generic",
			"description": "generic secret storage",
		},
		"sys/": map[string]string{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
		},
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}

	req = &Request{
		Operation: HelpOperation,
		Path:      "mounts",
	}
	resp, err = s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["help"] != "logical backend mount table" {
		t.Fatalf("got: %#v", resp.Data)
	}
}

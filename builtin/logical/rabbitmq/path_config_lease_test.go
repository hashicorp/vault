package rabbitmq

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_config_lease_RU(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err = b.Setup(config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"ttl":     "10h",
		"max_ttl": "20h",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/lease",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	configReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	if resp.Data["ttl"].(time.Duration) != 36000 {
		t.Fatalf("bad: ttl: expected:36000 actual:%d", resp.Data["ttl"].(time.Duration))
	}
	if resp.Data["max_ttl"].(time.Duration) != 72000 {
		t.Fatalf("bad: ttl: expected:72000 actual:%d", resp.Data["ttl"].(time.Duration))
	}
}

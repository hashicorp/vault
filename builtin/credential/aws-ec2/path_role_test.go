package awsec2

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestAwsEc2_RoleDurationSeconds(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"bound_iam_instance_profile_arn": "testarn",
		"ttl":     "10s",
		"max_ttl": "20s",
		"period":  "30s",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/testrole",
		Data:      roleData,
	}

	resp, err := b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	if int64(resp.Data["ttl"].(time.Duration)) != 10 {
		t.Fatalf("bad: period; expected: 10, actual: %d", resp.Data["ttl"])
	}
	if int64(resp.Data["max_ttl"].(time.Duration)) != 20 {
		t.Fatalf("bad: period; expected: 20, actual: %d", resp.Data["max_ttl"])
	}
	if int64(resp.Data["period"].(time.Duration)) != 30 {
		t.Fatalf("bad: period; expected: 30, actual: %d", resp.Data["period"])
	}
}

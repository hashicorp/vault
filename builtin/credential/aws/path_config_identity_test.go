package awsauth

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathConfigIdentity(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		if resp.IsError() {
			t.Fatalf("failed to read identity config entry")
		} else if resp.Data["iam_alias"] != nil && resp.Data["iam_alias"] != "" {
			t.Fatalf("returned alias is non-empty: %q", resp.Data["alias"])
		}
	}

	data := map[string]interface{}{
		"iam_alias": "invalid",
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("nil response from invalid config/identity request")
	}
	if !resp.IsError() {
		t.Fatalf("received non-error response from invalid config/identity request: %#v", resp)
	}

	data["iam_alias"] = identityAliasIAMFullArn
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("received error response from valid config/identity request: %#v", resp)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("nil response received from config/identity when data expected")
	} else if resp.IsError() {
		t.Fatalf("error response received from reading config/identity: %#v", resp)
	} else if resp.Data["iam_alias"] != identityAliasIAMFullArn {
		t.Fatalf("bad: expected response with iam_alias value of %q; got %#v", identityAliasIAMFullArn, resp)
	}
}

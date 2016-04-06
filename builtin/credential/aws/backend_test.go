package aws

import (
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_ConfigClient(t *testing.T) {
	config := logical.TestBackendConfig()
	storageView := &logical.InmemStorage{}
	config.StorageView = storageView

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		Backend:        b,
		Steps:          []logicaltest.TestStep{},
	})
}

func TestBackend_parseRoleTagValue(t *testing.T) {
	tag := "v1:XwuKhyyBNJc=:a=ami-fce3c696:p=root:t=3h0m0s:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	expected := roleTag{
		Version:  "v1",
		Nonce:    "XwuKhyyBNJc=",
		Policies: []string{"root"},
		MaxTTL:   10800000000000,
		ImageID:  "ami-fce3c696",
		HMAC:     "lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ=",
	}
	actual, err := parseRoleTagValue(tag)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !actual.Equal(&expected) {
		t.Fatalf("err: expected:%#v \ngot: %#v\n", expected, actual)
	}

	tag = "v2:XwuKhyyBNJc=:a=ami-fce3c696:p=root:t=3h0m0s:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	actual, err = parseRoleTagValue(tag)
	if err == nil {
		t.Fatalf("err: expected error due to invalid role tag version", err)
	}

	tag = "v1:XwuKhyyBNJc=:a=ami-fce3c696:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	expected = roleTag{
		Version: "v1",
		Nonce:   "XwuKhyyBNJc=",
		ImageID: "ami-fce3c696",
		HMAC:    "lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ=",
	}
	actual, err = parseRoleTagValue(tag)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	tag = "v1:XwuKhyyBNJc=:p=ami-fce3c696:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	actual, err = parseRoleTagValue(tag)
	if err == nil {
		t.Fatalf("err: expected error due to missing image ID", err)
	}
}

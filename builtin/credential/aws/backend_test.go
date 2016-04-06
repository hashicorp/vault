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

	data := map[string]interface{}{"access_key": "AKIAJBRHKV6EVTTNXDHA",
		"secret_key": "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",
		"region":     "us-east-1",
	}

	stepCreate := logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Data:      data,
	}

	stepUpdate := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data,
	}

	data2 := map[string]interface{}{"access_key": "AKIAJBRHKV6EVTTNXDHA",
		"secret_key": "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",
		"region":     "",
	}
	stepEmptyRegion := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data2,
		ErrorOk:   true,
	}

	data3 := map[string]interface{}{"access_key": "",
		"secret_key": "mCtSM8ZUEQ3mOFVZYPBQkf2sO6F/W7a5TVzrl3Oj",
		"region":     "us-east-1",
	}
	stepInvalidAccessKey := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data3,
		ErrorOk:   true,
	}

	data4 := map[string]interface{}{"access_key": "accesskey",
		"secret_key": "",
		"region":     "us-east-1",
	}
	stepInvalidSecretKey := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data4,
		ErrorOk:   true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			stepCreate,
			stepEmptyRegion,
			stepInvalidAccessKey,
			stepInvalidSecretKey,
			stepUpdate,
		},
	})

	configClientCreateRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storageView,
	}
	_, err = b.HandleRequest(configClientCreateRequest)
	if err != nil {
		t.Fatal(err)
	}

	clientConfig, err := clientConfigEntry(storageView)
	if err != nil {
		t.Fatal(err)
	}
	if clientConfig.AccessKey != data["access_key"] ||
		clientConfig.SecretKey != data["secret_key"] ||
		clientConfig.Region != data["region"] {
		t.Fatalf("bad: expected: %#v\ngot: %#v\n", data, clientConfig)
	}
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

package aws

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestStaticCredsRead(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	// insert a cred to get
	creds := &awsCredentials{
		AccessKeyID:     "foo",
		SecretAccessKey: "bar",
	}
	entry, err := logical.StorageEntryJSON(formatCredsStoragePath("test"), creds)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	b := Backend()

	req := &logical.Request{
		Storage: config.StorageView,
		Data: map[string]interface{}{
			"name": "test",
		},
	}
	resp, err := b.pathStaticCredsRead(context.Background(), req, staticCredsFieldData(req.Data))
	if err != nil {
		t.Fatal(err)
	}
	// resp will have credentials
	if resp == nil {
		t.Fatal("expected a non-nil response, but it was nil")
	}
}

func staticCredsFieldData(data map[string]interface{}) *framework.FieldData {
	schema := map[string]*framework.FieldSchema{
		paramRoleName: {
			Type:        framework.TypeString,
			Description: descRoleName,
		},
	}

	return &framework.FieldData{
		Raw:    data,
		Schema: schema,
	}
}

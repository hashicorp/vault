package aws

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestStaticCredsRead(t *testing.T) {
	// setup
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

	// cases
	cases := []struct {
		name           string
		roleName       string
		expectError    bool
		expectResponse bool
	}{
		{
			name:           "get existing creds",
			roleName:       "test",
			expectResponse: true,
		},
		{
			name:     "get non-existent creds",
			roleName: "this-doesnt-exist",
			// returns nil, nil
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := Backend()

			req := &logical.Request{
				Storage: config.StorageView,
				Data: map[string]interface{}{
					"name": c.roleName,
				},
			}
			resp, err := b.pathStaticCredsRead(context.Background(), req, staticCredsFieldData(req.Data))
			if (c.expectError && (err == nil)) || (!c.expectError && (err != nil)) {
				t.Fatal(err)
			}
			if (c.expectResponse && (resp == nil)) || (!c.expectResponse && (resp != nil)) {
				t.Fatal("expected a non-nil response, but it was nil")
			}
		})
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

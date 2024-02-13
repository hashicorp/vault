// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"testing"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestStaticCredsRead verifies that we can correctly read a cred that exists, and correctly _not read_
// a cred that does not exist.
func TestStaticCredsRead(t *testing.T) {
	// setup
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	bgCTX := context.Background() // for brevity later

	// insert a cred to get
	creds := &awsCredentials{
		AccessKeyID:     "foo",
		SecretAccessKey: "bar",
	}
	entry, err := logical.StorageEntryJSON(formatCredsStoragePath("test"), creds)
	if err != nil {
		t.Fatal(err)
	}
	err = config.StorageView.Put(bgCTX, entry)
	if err != nil {
		t.Fatal(err)
	}

	// cases
	cases := []struct {
		name             string
		roleName         string
		expectedError    error
		expectedResponse *logical.Response
	}{
		{
			name:     "get existing creds",
			roleName: "test",
			expectedResponse: &logical.Response{
				Data: structs.New(creds).Map(),
			},
		},
		{
			name:     "get non-existent creds",
			roleName: "this-doesnt-exist",
			// returns nil, nil
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := Backend(config)

			req := &logical.Request{
				Storage: config.StorageView,
				Data: map[string]interface{}{
					"name": c.roleName,
				},
			}
			resp, err := b.pathStaticCredsRead(bgCTX, req, staticCredsFieldData(req.Data))

			if err != c.expectedError {
				t.Fatalf("got error %q, but expected %q", err, c.expectedError)
			}
			if !reflect.DeepEqual(resp, c.expectedResponse) {
				t.Fatalf("got response %v, but expected %v", resp, c.expectedResponse)
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

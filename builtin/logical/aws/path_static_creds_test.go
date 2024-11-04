// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
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

// Test_awsCredentials_priority verifies that the expiration in the credentials
// is returned as the priority value when it is present, but otherwise the
// priority is now + the rotation period
func Test_awsCredentials_priority(t *testing.T) {
	expiration := time.Date(2023, 10, 24, 15, 21, 0o0, 0o0, time.UTC)
	roleConfig := staticRoleEntry{RotationPeriod: time.Hour}
	t.Run("use credential value", func(t *testing.T) {
		creds := &awsCredentials{
			Expiration: &expiration,
		}
		require.Equal(t, expiration.Unix(), creds.priority(roleConfig))
	})
	t.Run("use role value", func(t *testing.T) {
		hourUnix := time.Now().Add(time.Hour).Unix()
		creds := &awsCredentials{}
		require.InDelta(t, hourUnix, creds.priority(roleConfig), float64(time.Minute/time.Second))
	})
}

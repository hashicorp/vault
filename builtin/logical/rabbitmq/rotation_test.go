// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	tags = "policymaker,monitoring"
	vhostTopicJSON = `{"vhostOne":{"exchangeOneOne":{"write":".*","read":".*"},"exchangeOneTwo":{"write":".*","read":".*" }}}`
)

func TestBackend_Roles_Dynamic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.System = logical.TestSystemView()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		args         map[string]interface{}
		wantErr      bool
		expectedResp map[string]interface{}
	}{
		{
			name: "invalid role with no tags or vhost permissisons",
			args: map[string]interface{}{
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "valid role with tags",
			args: map[string]interface{}{
				"tags": tags,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags": tags,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "roles/test",
				Storage:   config.StorageView,
				Data:      tt.args,
			}
			// Create the role
			resp, err := b.HandleRequest(context.Background(), req)
			if tt.wantErr {
				assert.True(t, resp.IsError(), "expected error")
				return
			}
			assert.False(t, resp.IsError())
			assert.Nil(t, err)

			// Read the role
			req.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
			for k, v := range tt.expectedResp {
				assert.Equal(t, v, resp.Data[k])
			}

			// Delete the role
			req.Operation = logical.DeleteOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
		})
	}
}

func TestBackend_Roles_Static(t *testing.T) {
	config := logical.TestBackendConfig()
	config.System = logical.TestSystemView()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		args         map[string]interface{}
		wantErr      bool
		expectedResp map[string]interface{}
	}{
		{
			name: "invalid static role with no tags or vhost permissisons",
			args: map[string]interface{}{
				"username": "tester",
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "invalid static role with no username",
			args: map[string]interface{}{
				"tags": tags,
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "valid static role with tags",
			args: map[string]interface{}{
				"tags": tags,
				"username": "tester",
				"rotation_period": 3,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags": tags,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &logical.Request{
				Operation: logical.CreateOperation,
				Path:      "static-roles/test",
				Storage:   config.StorageView,
				Data:      tt.args,
			}
			// Create the role
			resp, err := b.HandleRequest(context.Background(), req)
			if tt.wantErr {
				assert.True(t, resp.IsError(), "expected error")
				return
			}
			assert.False(t, resp.IsError())
			assert.Nil(t, err)

			// Read the role
			req.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
			for k, v := range tt.expectedResp {
				assert.Equal(t, v, resp.Data[k])
			}

			// Delete the role
			req.Operation = logical.DeleteOperation
			resp, err = b.HandleRequest(context.Background(), req)
			assert.False(t, resp.IsError())
			assert.Nil(t, err)
		})
	}
}

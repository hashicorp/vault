// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	"github.com/stretchr/testify/assert"
)

const (
	tags           = "policymaker,monitoring"
	vhostTopicJSON = `{"vhostOne":{"exchangeOneOne":{"write":".*","read":".*"},"exchangeOneTwo":{"write":".*","read":".*" }}}`
)

func TestBackend_StaticRoles(t *testing.T) {
	cleanup, connectionURI := prepareRabbitMQTestContainer(t)
	defer cleanup()

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_uri":    connectionURI,
		"username":          "guest",
		"password":          "guest",
		"username_template": "",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
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
				"username":     "tester",
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "invalid static role with missing rotation period",
			args: map[string]interface{}{
				"tags":            tags,
				"username":        "tester",
			},
			wantErr: true,
		},
		{
			name: "invalid static role with no username",
			args: map[string]interface{}{
				"tags":         tags,
				"vhost_topics": vhostTopicJSON,
			},
			wantErr: true,
		},
		{
			name: "valid static role with tags",
			args: map[string]interface{}{
				"tags":            tags,
				"username":        "tester",
				"rotation_period": 7,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags":            tags,
				"username":        "tester",
				"rotation_period": 7.0,
			},
		},
		{
			name: "valid static role with revoke on delete",
			args: map[string]interface{}{
				"tags":                  tags,
				"username":              "tester",
				"revoke_user_on_delete": true,
				"rotation_period": 10,
			},
			wantErr: false,
			expectedResp: map[string]interface{}{
				"tags":                  tags,
				"username":              "tester",
				"revoke_user_on_delete": true,
				"rotation_period": 10.0,
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

func TestBackend_Rotation_StaticRoles(t *testing.T) {
	bgCTX := context.Background()

	cleanup, connectionURI := prepareRabbitMQTestContainer(t)
	defer cleanup()

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend(config)
	if err := b.Setup(bgCTX, config); err != nil {
		t.Fatal(err)
	}
	b.credRotationQueue = queue.New()
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_uri":    connectionURI,
		"username":          "guest",
		"password":          "guest",
		"username_template": "",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	type credToInsert struct {
		name    string
		config  staticRoleEntry
		age     time.Duration
		changed bool
	}
	tests := []struct {
		name  string
		creds []credToInsert
	}{
		{
			name: "refresh one",
			creds: []credToInsert{
				{
					name: "test0",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 2 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     5 * time.Second,
					changed: true,
				},
			},
		},
		{
			name: "refresh some",
			creds: []credToInsert{
				{
					name: "test1",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 2 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     5 * time.Second,
					changed: true,
				},
				{
					name: "test2",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 10 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     5 * time.Second,
					changed: false,
				},
				{
					name: "test3",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 10 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     11 * time.Second,
					changed: true,
				},
				{
					name: "test4",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 15 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     11 * time.Second,
					changed: false,
				},
			},
		},
		{
			name: "refresh none",
			creds: []credToInsert{
				{
					name: "test5",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 10 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     5 * time.Second,
					changed: false,
				},
				{
					name: "test6",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 15 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     5 * time.Second,
					changed: false,
				},
				{
					name: "test7",
					config: staticRoleEntry{
						Username:       "jane-doe",
						RotationPeriod: 30 * time.Second,
						RoleEntry: roleEntry{
							Tags: tags,
							VHosts: map[string]vhostPermission{
								"vhostOne": {
									Write: "*",
									Read:  "*",
								},
							},
						},
					},
					age:     25 * time.Second,
					changed: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for idx, cred := range tt.creds {
				// use refenrence from creds array to allow to update the password outside of range
				err := b.createStaticCredential(bgCTX, config.StorageView, &tt.creds[idx].config, cred.name)
				if err != nil {
					t.Fatalf("couldn't create static credential: %s", err)
				}

				item := &queue.Item{
					Key:      cred.name,
					Value:    tt.creds[idx].config,
					Priority: time.Now().Add(-1 * cred.age).Add(cred.config.RotationPeriod).Unix(),
				}
				err = b.credRotationQueue.Push(item)
				if err != nil {
					t.Fatalf("couldn't push item onto queue: %s", err)
				}
			}

			req := &logical.Request{
				Storage: config.StorageView,
			}
			err = b.rotateExpiredStaticCreds(bgCTX, req)
			if err != nil {
				t.Fatalf("got an error rotating credentials: %s", err)
			}

			for i, cred := range tt.creds {
				entry, err := config.StorageView.Get(bgCTX, rabbitMQStaticRolePath+cred.name)
				if err != nil {
					t.Fatalf("got an error retrieving credentials %d", i)
				}
				var out staticRoleEntry
				err = entry.DecodeJSON(&out)
				if err != nil {
					t.Fatalf("could not unmarshal storage view entry for cred %d to a rabbitmq static role: %s", i, err)
				}

				if cred.changed && out.Password == cred.config.Password {
					t.Fatalf("expected the password for cred %d to have changed, but it hasn't", i)
				} else if !cred.changed && out.Password != cred.config.Password {
					t.Fatalf("expected the password for cred %d to have stayed the same, but it changed", i)
				}
			}
		})
	}
}

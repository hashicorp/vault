// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authmetadata

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type environment struct {
	ctx     context.Context
	storage logical.Storage
	backend logical.Backend
}

func TestAcceptance(t *testing.T) {
	ctx := context.Background()
	storage := &logical.InmemStorage{}
	b, err := backend(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	env := &environment{
		ctx:     ctx,
		storage: storage,
		backend: b,
	}
	t.Run("test initial fields are default", env.TestInitialFieldsAreDefault)
	t.Run("test fields can be unset", env.TestAuthMetadataCanBeUnset)
	t.Run("test defaults can be restored", env.TestDefaultCanBeReused)
	t.Run("test default plus more cannot be selected", env.TestDefaultPlusMoreCannotBeSelected)
	t.Run("test only non-defaults can be selected", env.TestOnlyNonDefaultsCanBeSelected)
	t.Run("test bad field results in useful error", env.TestAddingBadField)
}

func (e *environment) TestInitialFieldsAreDefault(t *testing.T) {
	// On the first read of auth_metadata, when nothing has been touched,
	// we should receive the default field(s) if a read is performed.
	resp, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected non-nil response")
	}
	if !reflect.DeepEqual(resp.Data[authMetadataFields.FieldName], []string{"role_name"}) {
		t.Fatal("expected default field of role_name to be returned")
	}

	// The auth should only have the default metadata.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			"role_name": "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.Alias == nil || resp.Auth.Alias.Metadata == nil {
		t.Fatalf("expected alias metadata")
	}
	if len(resp.Auth.Alias.Metadata) != 1 {
		t.Fatal("expected only 1 field")
	}
	if resp.Auth.Alias.Metadata["role_name"] != "something" {
		t.Fatal("expected role_name to be something")
	}
}

func (e *environment) TestAuthMetadataCanBeUnset(t *testing.T) {
	// We should be able to set the auth_metadata to empty by sending an
	// explicitly empty array.
	resp, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			authMetadataFields.FieldName: []string{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}

	// Now we should receive no fields for auth_metadata.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected non-nil response")
	}
	if !reflect.DeepEqual(resp.Data[authMetadataFields.FieldName], []string{}) {
		t.Fatal("expected no fields to be returned")
	}

	// The auth should have no metadata.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			"role_name": "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.Alias == nil || resp.Auth.Alias.Metadata == nil {
		t.Fatal("expected alias metadata")
	}
	if len(resp.Auth.Alias.Metadata) != 0 {
		t.Fatal("expected 0 fields")
	}
}

func (e *environment) TestDefaultCanBeReused(t *testing.T) {
	// Now if we set it to "default", the default fields should
	// be restored.
	resp, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			authMetadataFields.FieldName: []string{"default"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}

	// Let's make sure we've returned to the default fields.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected non-nil response")
	}
	if !reflect.DeepEqual(resp.Data[authMetadataFields.FieldName], []string{"role_name"}) {
		t.Fatal("expected default field of role_name to be returned")
	}

	// We should again only receive the default field on the login.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			"role_name": "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.Alias == nil || resp.Auth.Alias.Metadata == nil {
		t.Fatal("expected alias metadata")
	}
	if len(resp.Auth.Alias.Metadata) != 1 {
		t.Fatal("expected only 1 field")
	}
	if resp.Auth.Alias.Metadata["role_name"] != "something" {
		t.Fatal("expected role_name to be something")
	}
}

func (e *environment) TestDefaultPlusMoreCannotBeSelected(t *testing.T) {
	// We should not be able to set it to "default" plus 1 optional field.
	_, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			authMetadataFields.FieldName: []string{"default", "remote_addr"},
		},
	})
	if err == nil {
		t.Fatal("expected err")
	}
}

func (e *environment) TestOnlyNonDefaultsCanBeSelected(t *testing.T) {
	// Omit all default fields and just select one.
	resp, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			authMetadataFields.FieldName: []string{"remote_addr"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}

	// Make sure that worked.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected non-nil response")
	}
	if !reflect.DeepEqual(resp.Data[authMetadataFields.FieldName], []string{"remote_addr"}) {
		t.Fatal("expected remote_addr to be returned")
	}

	// Ensure only the selected one is on logins.
	// They both should now appear on the login.
	resp, err = e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			"role_name": "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.Alias == nil || resp.Auth.Alias.Metadata == nil {
		t.Fatal("expected alias metadata")
	}
	if len(resp.Auth.Alias.Metadata) != 1 {
		t.Fatal("expected only 1 field")
	}
	if resp.Auth.Alias.Metadata["remote_addr"] != "http://foo.com" {
		t.Fatal("expected remote_addr to be http://foo.com")
	}
}

func (e *environment) TestAddingBadField(t *testing.T) {
	// Try adding an unsupported field.
	resp, err := e.backend.HandleRequest(e.ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   e.storage,
		Connection: &logical.Connection{
			RemoteAddr: "http://foo.com",
		},
		Data: map[string]interface{}{
			authMetadataFields.FieldName: []string{"asl;dfkj"},
		},
	})
	if err == nil {
		t.Fatal("expected err")
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.IsError() {
		t.Fatal("expected error response")
	}
}

// We expect people to embed the Handler on their
// config so it automatically makes its helper methods
// available and easy to find wherever the config is
// needed. Explicitly naming it in json avoids it
// automatically being named "Handler" by Go's JSON
// marshalling library.
type fakeConfig struct {
	*Handler `json:"auth_metadata_handler"`
}

type fakeBackend struct {
	*framework.Backend
}

// We expect each back-end to explicitly define the fields that
// will be included by default, and optionally available.
var authMetadataFields = &Fields{
	FieldName: "some_field_name",
	Default: []string{
		"role_name", // This would likely never change because the alias is the role name.
	},
	AvailableToAdd: []string{
		"remote_addr", // This would likely change with every new caller.
	},
}

func configPath() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			authMetadataFields.FieldName: FieldSchema(authMetadataFields),
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: func(ctx context.Context, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
					entryRaw, err := req.Storage.Get(ctx, "config")
					if err != nil {
						return nil, err
					}
					conf := &fakeConfig{
						Handler: NewHandler(authMetadataFields),
					}
					if entryRaw != nil {
						if err := entryRaw.DecodeJSON(conf); err != nil {
							return nil, err
						}
					}
					// Note that even if the config entry was nil, we return
					// a populated response to give info on what the default
					// auth metadata is when unconfigured.
					return &logical.Response{
						Data: map[string]interface{}{
							authMetadataFields.FieldName: conf.AuthMetadata(),
						},
					}, nil
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: func(ctx context.Context, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
					entryRaw, err := req.Storage.Get(ctx, "config")
					if err != nil {
						return nil, err
					}
					conf := &fakeConfig{
						Handler: NewHandler(authMetadataFields),
					}
					if entryRaw != nil {
						if err := entryRaw.DecodeJSON(conf); err != nil {
							return nil, err
						}
					}
					// This is where we read in the user's given auth metadata.
					if err := conf.ParseAuthMetadata(fd); err != nil {
						// Since this will only error on bad input, it's best to give
						// a 400 response with the explicit problem included.
						return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
					}
					entry, err := logical.StorageEntryJSON("config", conf)
					if err != nil {
						return nil, err
					}
					if err = req.Storage.Put(ctx, entry); err != nil {
						return nil, err
					}
					return nil, nil
				},
			},
		},
	}
}

func loginPath() *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"role_name": {
				Type:     framework.TypeString,
				Required: true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: func(ctx context.Context, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
					entryRaw, err := req.Storage.Get(ctx, "config")
					if err != nil {
						return nil, err
					}
					conf := &fakeConfig{
						Handler: NewHandler(authMetadataFields),
					}
					if entryRaw != nil {
						if err := entryRaw.DecodeJSON(conf); err != nil {
							return nil, err
						}
					}
					auth := &logical.Auth{
						Alias: &logical.Alias{
							Name: fd.Get("role_name").(string),
						},
					}
					// Here we provide everything and let the method strip out
					// the undesired stuff.
					if err := conf.PopulateDesiredMetadata(auth, map[string]string{
						"role_name":   fd.Get("role_name").(string),
						"remote_addr": req.Connection.RemoteAddr,
					}); err != nil {
						fmt.Println("unable to populate due to " + err.Error())
					}
					return &logical.Response{
						Auth: auth,
					}, nil
				},
			},
		},
	}
}

func backend(ctx context.Context, storage logical.Storage) (logical.Backend, error) {
	b := &fakeBackend{
		Backend: &framework.Backend{
			Paths: []*framework.Path{
				configPath(),
				loginPath(),
			},
		},
	}
	if err := b.Setup(ctx, &logical.BackendConfig{
		StorageView: storage,
		Logger:      hclog.Default(),
	}); err != nil {
		return nil, err
	}
	return b, nil
}

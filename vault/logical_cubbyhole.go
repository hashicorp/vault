// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

const (
	scopeToken  = "per-token"
	scopeEntity = "per-entity"
)

// CubbyholeBackendFactory constructs a new cubbyhole backend
func CubbyholeBackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &CubbyholeBackend{}
	b.Backend = &framework.Backend{
		Help:           strings.TrimSpace(cubbyholeHelp),
		RunningVersion: versions.GetBuiltinVersion(consts.PluginTypeSecrets, "cubbyhole"),
	}

	b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}
	b.Backend.Setup(ctx, conf)

	return b, nil
}

// CubbyholeBackend is used for storing secrets directly into the physical
// backend. The secrets are encrypted in the durable storage.
// This differs from kv in that every token has its own private
// storage view. The view is removed when the token expires.
type CubbyholeBackend struct {
	*framework.Backend

	saltUUID    string
	storageView logical.Storage
}

func (b *CubbyholeBackend) paths() []*framework.Path {
	return []*framework.Path{

		{
			Pattern: "config/scope",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "cubbyhole",
			},

			Fields: map[string]*framework.FieldSchema{
				"scope": {
					Type:        framework.TypeString,
					Description: `Defines the behavior of the physical storage used for cubbyhole\'s secrets. If per token, cubbyhole is destroyed after the token\'s expiration. If per entity, secrets are persisted and are linked to user\'s identity lifetime`,
					Default:     "per-token",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathScopeWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "write",
						OperationSuffix: "scope",
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathScopeRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "scope",
					},
				},
			},
		},
		{
			Pattern: framework.MatchAllRegex("path"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "cubbyhole",
			},

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Specifies the path of the secret.",
				},
			},

			TakesArbitraryInput: true,

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
					Summary: "Retrieve the secret at the specified location.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "write",
					},
					Summary: "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "write",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					Summary: "Deletes the secret at the specified location.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleList,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "list",
					},
					Summary:     "List secret entries at the specified location.",
					Description: "Folders are suffixed with /. The input must be a folder; list on a file will not return a value. The values themselves are not accessible via this command.",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,

			HelpSynopsis:    strings.TrimSpace(cubbyholeHelpSynopsis),
			HelpDescription: strings.TrimSpace(cubbyholeHelpDescription),
		},
	}
}

func (b *CubbyholeBackend) revoke(ctx context.Context, view *BarrierView, saltedToken string) error {
	if saltedToken == "" {
		return fmt.Errorf("client token empty during revocation")
	}

	if err := logical.ClearView(ctx, view.SubView(saltedToken+"/")); err != nil {
		return err
	}

	return nil
}

func (b *CubbyholeBackend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.ClientToken+"/"+req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

func (b *CubbyholeBackend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var pathPrefix string

	scope, err := b.Scope(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	switch scope {
	case scopeToken:
		if req.ClientToken == "" {
			return nil, fmt.Errorf("client token empty")
		}
		pathPrefix = req.ClientToken
	case scopeEntity:
		if req.EntityID == "" {
			return nil, fmt.Errorf("entity ID empty")
		}
		pathPrefix = req.EntityID
	}

	path := data.Get("path").(string)

	if path == "" {
		return nil, fmt.Errorf("missing path")
	}

	// Read the path
	out, err := req.Storage.Get(ctx, pathPrefix+"/"+path)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	// Fast-path the no data case
	if out == nil {
		return nil, nil
	}

	// Decode the data
	var rawData map[string]interface{}
	if err := jsonutil.DecodeJSON(out.Value, &rawData); err != nil {
		return nil, fmt.Errorf("json decoding failed: %w", err)
	}

	// Generate the response
	resp := &logical.Response{
		Data: rawData,
	}

	return resp, nil
}

func (b *CubbyholeBackend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var pathPrefix string

	scope, err := b.Scope(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	switch scope {
	case scopeToken:
		if req.ClientToken == "" {
			return nil, fmt.Errorf("client token empty")
		}
		pathPrefix = req.ClientToken
	case scopeEntity:
		if req.EntityID == "" {
			return nil, fmt.Errorf("entity ID empty")
		}
		pathPrefix = req.EntityID
	}
	// Check that some fields are given
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("missing data fields")
	}

	path := data.Get("path").(string)

	if path == "" {
		return nil, fmt.Errorf("missing path")
	}

	// JSON encode the data
	buf, err := json.Marshal(req.Data)
	if err != nil {
		return nil, fmt.Errorf("json encoding failed: %w", err)
	}

	// Write out a new key
	entry := &logical.StorageEntry{
		Key:   pathPrefix + "/" + path,
		Value: buf,
	}
	if req.WrapInfo != nil && req.WrapInfo.SealWrap {
		entry.SealWrap = true
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to write: %w", err)
	}

	return nil, nil
}

func (b *CubbyholeBackend) handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	path := data.Get("path").(string)

	// Delete the key at the request path
	if err := req.Storage.Delete(ctx, req.ClientToken+"/"+path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *CubbyholeBackend) handleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Right now we only handle directories, so ensure it ends with / We also
	// check if it's empty so we don't end up doing a listing on '<client
	// token>//'
	path := data.Get("path").(string)
	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// List the keys at the prefix given by the request
	keys, err := req.Storage.List(ctx, req.ClientToken+"/"+path)
	if err != nil {
		return nil, err
	}

	// Strip the token
	strippedKeys := make([]string, len(keys))
	for i, key := range keys {
		strippedKeys[i] = strings.TrimPrefix(key, req.ClientToken+"/")
	}

	// Generate the response
	return logical.ListResponse(strippedKeys), nil
}

// Scope returns the lease
func (b *CubbyholeBackend) Scope(ctx context.Context, s logical.Storage) (string, error) {
	entry, err := s.Get(ctx, "config/scope")
	if err != nil {
		return scopeToken, err
	}
	if entry == nil {
		return scopeToken, nil
	}

	var result string
	if err := entry.DecodeJSON(&result); err != nil {
		return scopeToken, err
	}

	return result, nil
}

func (b *CubbyholeBackend) pathScopeWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	scopeRaw := d.Get("scope").(string)

	// Store it
	entry, err := logical.StorageEntryJSON("config/scope", scopeRaw)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *CubbyholeBackend) pathScopeRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	scope, err := b.Scope(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"scope": scope,
		},
	}, nil
}

const cubbyholeHelp = `
The cubbyhole backend reads and writes arbitrary secrets to the backend.
The secrets are encrypted/decrypted by Vault: they are never stored
unencrypted in the backend and the backend never has an opportunity to
see the unencrypted value.

This backend differs from the 'kv' backend in that it is namespaced
per-token. Tokens can only read and write their own values, with no
sharing possible (per-token cubbyholes). This can be useful for implementing
certain authentication workflows, as well as "scratch" areas for individual
clients. When the token is revoked, the entire set of stored values for that
token is also removed.
`

const cubbyholeHelpSynopsis = `
Pass-through secret storage to a token-specific cubbyhole in the storage
backend, allowing you to read/write arbitrary data into secret storage.
`

const cubbyholeHelpDescription = `
The cubbyhole backend reads and writes arbitrary data into secret storage,
encrypting it along the way.

The view into the cubbyhole storage space is different for each token; it is
a per-token cubbyhole. When the token is revoked all values are removed.
`

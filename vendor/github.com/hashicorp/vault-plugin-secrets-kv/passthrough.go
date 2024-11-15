// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
)

// PassthroughBackendFactory returns a PassthroughBackend
// with leases switched off
func PassthroughBackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	return LeaseSwitchedPassthroughBackendFactory(ctx, conf, false)
}

// LeasedPassthroughBackendFactory returns a PassthroughBackend
// with leases switched on
func LeasedPassthroughBackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	return LeaseSwitchedPassthroughBackendFactory(ctx, conf, true)
}

// LeaseSwitchedPassthroughBackendFactory returns a PassthroughBackend
// with leases switched on or off
func LeaseSwitchedPassthroughBackendFactory(ctx context.Context, conf *logical.BackendConfig, leases bool) (logical.Backend, error) {
	b := &PassthroughBackend{
		generateLeases: leases,
	}

	backend := &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        strings.TrimSpace(passthroughHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"*",
			},
		},

		Paths: []*framework.Path{
			{
				Pattern: framework.MatchAllRegex("path"),

				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixKVv1,
				},

				Fields: map[string]*framework.FieldSchema{
					"path": {
						Type:        framework.TypeString,
						Description: "Location of the secret.",
					},
				},

				// The regex and field definition above are purely for the benefit of OpenAPI and generated
				// documentation. The actual request processing functions consult req.Path directly.
				// See documentation of handleReadOrRenew for more detail.

				TakesArbitraryInput: true,

				Operations: map[logical.Operation]framework.OperationHandler{
					logical.ReadOperation: &framework.PathOperation{
						Callback: b.handleReadOrRenew(),
						DisplayAttrs: &framework.DisplayAttributes{
							OperationVerb: "read",
						},
						Responses: map[int][]framework.Response{
							http.StatusOK: {{
								Description: http.StatusText(http.StatusOK),
								Fields:      nil, // dynamic fields
							}},
						},
					},
					logical.CreateOperation: &framework.PathOperation{
						Callback: b.handleWrite(),
						DisplayAttrs: &framework.DisplayAttributes{
							OperationVerb: "write",
						},
						Responses: map[int][]framework.Response{
							http.StatusNoContent: {{
								Description: http.StatusText(http.StatusNoContent),
							}},
						},
					},
					logical.UpdateOperation: &framework.PathOperation{
						Callback: b.handleWrite(),
						DisplayAttrs: &framework.DisplayAttributes{
							OperationVerb: "write",
						},
						Responses: map[int][]framework.Response{
							http.StatusNoContent: {{
								Description: http.StatusText(http.StatusNoContent),
							}},
						},
					},
					logical.DeleteOperation: &framework.PathOperation{
						Callback: b.handleDelete(),
						DisplayAttrs: &framework.DisplayAttributes{
							OperationVerb: "delete",
						},
						Responses: map[int][]framework.Response{
							http.StatusNoContent: {{
								Description: http.StatusText(http.StatusNoContent),
							}},
						},
					},
					logical.ListOperation: &framework.PathOperation{
						Callback: b.handleList(),
						DisplayAttrs: &framework.DisplayAttributes{
							OperationVerb: "list",
						},
					},
				},

				ExistenceCheck: b.handleExistenceCheck(),

				HelpSynopsis:    strings.TrimSpace(passthroughHelpSynopsis),
				HelpDescription: strings.TrimSpace(passthroughHelpDescription),
			},
		},
		Secrets: []*framework.Secret{
			{
				Type: "kv",

				Renew: b.handleReadOrRenew(),
				Revoke: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
					// This is a no-op
					return nil, nil
				},
			},
		},
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}
	backend.Setup(ctx, conf)
	b.Backend = backend

	return b, nil
}

// PassthroughBackend is used storing secrets directly into the physical
// backend. The secrets are encrypted in the durable storage and custom TTL
// information can be specified, but otherwise this backend doesn't do anything
// fancy.
type PassthroughBackend struct {
	*framework.Backend
	generateLeases bool
}

func (b *PassthroughBackend) handleExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		out, err := req.Storage.Get(ctx, req.Path)
		if err != nil {
			return false, fmt.Errorf("existence check failed: %w", err)
		}

		return out != nil, nil
	}
}

// handleReadOrRenew is called both for ReadOperations and RenewOperations. The RenewOperation case only applies when
// using the rather obscure feature of mounting the backend with option leased_passthrough=true, and storing a duration
// within the secret data itself, at the key "ttl" (preferred) or "lease" (deprecated). When that is done, we return a
// renewable lease, and treat invoking the renew as a request to re-read the same path that was originally read. We use
// req.Path, instead of a more conventional data.Get("path").(string), to read the requested path because the data
// attached to a RenewOperation is NOT the original request data - it is the response data! Technically we need only do
// this for this function, but for consistency we also use req.Path throughout the other handlers for this path pattern.
func (b *PassthroughBackend) handleReadOrRenew() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Read the path
		out, err := req.Storage.Get(ctx, req.Path)
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

		var resp *logical.Response
		if b.generateLeases {
			// Generate the response
			resp = b.Secret("kv").Response(rawData, nil)
			resp.Secret.Renewable = false
		} else {
			resp = &logical.Response{
				Secret: &logical.Secret{},
				Data:   rawData,
			}
		}

		// Ensure seal wrapping is carried through if the response is
		// response-wrapped
		if out.SealWrap {
			if resp.WrapInfo == nil {
				resp.WrapInfo = &wrapping.ResponseWrapInfo{}
			}
			resp.WrapInfo.SealWrap = out.SealWrap
		}

		// Check if there is a ttl key
		ttlDuration := b.System().DefaultLeaseTTL()
		ttlRaw, ok := rawData["ttl"]
		if !ok {
			ttlRaw, ok = rawData["lease"]
		}
		if ok {
			dur, err := parseutil.ParseDurationSecond(ttlRaw)
			if err == nil {
				ttlDuration = dur
			}

			if b.generateLeases {
				resp.Secret.Renewable = true
			}
		}

		resp.Secret.TTL = ttlDuration

		return resp, nil
	}
}

func (b *PassthroughBackend) handleWrite() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		if req.Path == "" {
			return logical.ErrorResponse("missing path"), nil
		}

		// Check that some fields are given
		if len(req.Data) == 0 {
			return logical.ErrorResponse("missing data fields"), nil
		}

		// JSON encode the data
		buf, err := json.Marshal(req.Data)
		if err != nil {
			return nil, fmt.Errorf("json encoding failed: %w", err)
		}

		// Write out a new key
		entry := &logical.StorageEntry{
			Key:   req.Path,
			Value: buf,
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, fmt.Errorf("failed to write: %w", err)
		}

		kvEvent(ctx, b.Backend, "write", req.Path, req.Path, true, 1)

		return nil, nil
	}
}

func (b *PassthroughBackend) handleDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Delete the key at the request path
		if err := req.Storage.Delete(ctx, req.Path); err != nil {
			return nil, err
		}

		kvEvent(ctx, b.Backend, "delete", req.Path, "", true, 1)

		return nil, nil
	}
}

func (b *PassthroughBackend) handleList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Right now we only handle directories, so ensure it ends with /; however,
		// some physical backends may not handle the "/" case properly, so only add
		// it if we're not listing the root
		path := req.Path
		if path != "" && !strings.HasSuffix(path, "/") {
			path = path + "/"
		}

		// List the keys at the prefix given by the request
		keys, err := req.Storage.List(ctx, path)
		if err != nil {
			return nil, err
		}

		// Generate the response
		return logical.ListResponse(keys), nil
	}
}

const passthroughHelp = `
The kv backend reads and writes arbitrary secrets to the backend.
The secrets are encrypted/decrypted by Vault: they are never stored
unencrypted in the backend and the backend never has an opportunity to
see the unencrypted value.

TTLs can be set on a per-secret basis. These TTLs will be sent down
when that secret is read, and it is assumed that some outside process will
revoke and/or replace the secret at that path.
`

const passthroughHelpSynopsis = `
Pass-through secret storage to the storage backend, allowing you to
read/write arbitrary data into secret storage.
`

const passthroughHelpDescription = `
The pass-through backend reads and writes arbitrary data into secret storage,
encrypting it along the way.

A TTL can be specified when writing with the "ttl" field. If given, the
duration of leases returned by this backend will be set to this value. This
can be used as a hint from the writer of a secret to the consumer of a secret
that the consumer should re-read the value before the TTL has expired.
However, any revocation must be handled by the user of this backend; the lease
duration does not affect the provided data in any way.
`

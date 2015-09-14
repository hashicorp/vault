package vault

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// logical.Factory
func PassthroughBackendFactory(conf *logical.BackendConfig) (logical.Backend, error) {
	var b PassthroughBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(passthroughHelp),

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: ".*",
				Fields: map[string]*framework.FieldSchema{
					//TODO: Deprecated in 0.3; remove in 0.4
					"lease": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Lease time for this key when read. Ex: 1h",
					},
					"ttl": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "TTL time for this key when read. Ex: 1h",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.handleRead,
					logical.WriteOperation:  b.handleWrite,
					logical.DeleteOperation: b.handleDelete,
					logical.ListOperation:   b.handleList,
				},

				HelpSynopsis:    strings.TrimSpace(passthroughHelpSynopsis),
				HelpDescription: strings.TrimSpace(passthroughHelpDescription),
			},
		},

		Secrets: []*framework.Secret{
			&framework.Secret{
				Type: "generic",

				Renew:  b.handleRead,
				Revoke: b.handleRevoke,
			},
		},
	}

	if conf == nil {
		return nil, fmt.Errorf("Configuation passed into backend is nil")
	}
	b.Backend.Setup(conf)

	return b, nil
}

// PassthroughBackend is used storing secrets directly into the physical
// backend. The secrest are encrypted in the durable storage and custom TTL
// information can be specified, but otherwise this backend doesn't do anything
// fancy.
type PassthroughBackend struct {
	*framework.Backend
}

func (b *PassthroughBackend) handleRevoke(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// This is a no-op
	return nil, nil
}

func (b *PassthroughBackend) handleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Read the path
	out, err := req.Storage.Get(req.Path)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	// Fast-path the no data case
	if out == nil {
		return nil, nil
	}

	// Decode the data
	var rawData map[string]interface{}
	if err := json.Unmarshal(out.Value, &rawData); err != nil {
		return nil, fmt.Errorf("json decoding failed: %v", err)
	}

	// Generate the response
	resp := b.Secret("generic").Response(rawData, nil)
	resp.Secret.Renewable = false

	// Check if there is a ttl key
	var ttl string
	ttl, _ = rawData["lease"].(string)
	if len(ttl) == 0 {
		ttl, _ = rawData["ttl"].(string)
	}

	if len(ttl) != 0 {
		ttlDuration, err := time.ParseDuration(ttl)
		if err == nil {
			resp.Secret.Renewable = true
			resp.Secret.TTL = ttlDuration
		}
	}

	return resp, nil
}

func (b *PassthroughBackend) handleWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Check that some fields are given
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("missing data fields")
	}

	// JSON encode the data
	buf, err := json.Marshal(req.Data)
	if err != nil {
		return nil, fmt.Errorf("json encoding failed: %v", err)
	}

	// Write out a new key
	entry := &logical.StorageEntry{
		Key:   req.Path,
		Value: buf,
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, fmt.Errorf("failed to write: %v", err)
	}

	return nil, nil
}

func (b *PassthroughBackend) handleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Delete the key at the request path
	if err := req.Storage.Delete(req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *PassthroughBackend) handleList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// List the keys at the prefix given by the request
	keys, err := req.Storage.List(req.Path)
	if err != nil {
		return nil, err
	}

	// Generate the response
	return logical.ListResponse(keys), nil
}

const passthroughHelp = `
The generic backend reads and writes arbitrary secrets to the backend.
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

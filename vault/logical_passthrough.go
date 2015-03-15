package vault

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
)

// logical.Factory
func PassthroughBackendFactory(map[string]string) (logical.Backend, error) {
	return new(PassthroughBackend), nil
}

// PassthroughBackend is used storing secrets directly into the physical
// backend. The secrest are encrypted in the durable storage and custom lease
// information can be specified, but otherwise this backend doesn't do anything
// fancy.
type PassthroughBackend struct{}

func (b *PassthroughBackend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	// TODO(mitchellh): help, let's just do it when we migrate to helper/backend

	switch req.Operation {
	case logical.ReadOperation:
		return b.handleRead(req)
	case logical.WriteOperation:
		return b.handleWrite(req)
	case logical.DeleteOperation:
		return b.handleDelete(req)
	case logical.ListOperation:
		return b.handleList(req)
	default:
		return nil, logical.ErrUnsupportedOperation
	}
}

func (b *PassthroughBackend) RootPaths() []string {
	return nil
}

func (b *PassthroughBackend) handleRead(req *logical.Request) (*logical.Response, error) {
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
	var raw map[string]interface{}
	if err := json.Unmarshal(out.Value, &raw); err != nil {
		return nil, fmt.Errorf("json decoding failed: %v", err)
	}

	// Check if there is a lease key
	leaseVal, ok := raw["lease"].(string)
	var lease *logical.Lease
	if ok {
		leaseDuration, err := time.ParseDuration(leaseVal)
		if err == nil {
			lease = &logical.Lease{
				Renewable:    false,
				Revokable:    false,
				Duration:     leaseDuration,
				MaxDuration:  leaseDuration,
				MaxIncrement: 0,
			}
		}
	}

	// Generate the response
	resp := &logical.Response{
		IsSecret: true,
		Lease:    lease,
		Data:     raw,
	}
	return resp, nil
}

func (b *PassthroughBackend) handleWrite(req *logical.Request) (*logical.Response, error) {
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

func (b *PassthroughBackend) handleDelete(req *logical.Request) (*logical.Response, error) {
	// Delete the key at the request path
	if err := req.Storage.Delete(req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *PassthroughBackend) handleList(req *logical.Request) (*logical.Response, error) {
	// List the keys at the prefix given by the request
	keys, err := req.Storage.List(req.Path)
	if err != nil {
		return nil, err
	}

	// Generate the response
	return logical.ListResponse(keys), nil
}

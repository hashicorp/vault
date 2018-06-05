package mock

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// kvPaths is used to test CRUD and List operations. It is a simplified
// version of the passthrough backend that only accepts string values.
func kvPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "kv/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathKVList,
			},
		},
		&framework.Path{
			Pattern: "kv/" + framework.GenericNameRegex("key"),
			Fields: map[string]*framework.FieldSchema{
				"key":     &framework.FieldSchema{Type: framework.TypeString},
				"value":   &framework.FieldSchema{Type: framework.TypeString},
				"version": &framework.FieldSchema{Type: framework.TypeInt},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathKVRead,
				logical.CreateOperation: b.pathKVCreateUpdate,
				logical.UpdateOperation: b.pathKVCreateUpdate,
				logical.DeleteOperation: b.pathKVDelete,
			},
		},
	}
}

func (b *backend) pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

func (b *backend) pathKVRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	version := data.Get("version").(int)

	entry, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	value := string(entry.Value)

	b.Logger().Info("reading value", "key", req.Path, "value", value)
	// Return the secret
	resp := &logical.Response{
		Data: map[string]interface{}{
			"value":   value,
			"version": version,
		},
	}
	if version != 0 {
		resp.Data["version"] = version
	}
	return resp, nil
}

func (b *backend) pathKVCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	value := data.Get("value").(string)

	b.Logger().Info("storing value", "key", req.Path, "value", value)
	entry := &logical.StorageEntry{
		Key:   req.Path,
		Value: []byte(value),
	}

	s := req.Storage
	err := s.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}, nil
}

func (b *backend) pathKVDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathKVList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, "kv/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

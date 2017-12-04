package mock

import (
	"fmt"

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
				"key":   &framework.FieldSchema{Type: framework.TypeString},
				"value": &framework.FieldSchema{Type: framework.TypeString},
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

func (b *backend) pathExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

func (b *backend) pathKVRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(req.Path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	value := string(entry.Value)

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}, nil
}

func (b *backend) pathKVCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	value := data.Get("value").(string)

	entry := &logical.StorageEntry{
		Key:   req.Path,
		Value: []byte(value),
	}

	s := req.Storage
	err := s.Put(entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}, nil
}

func (b *backend) pathKVDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathKVList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List("kv/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

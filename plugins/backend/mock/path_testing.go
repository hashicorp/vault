package mock

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTesting(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "test/ing",
		Fields: map[string]*framework.FieldSchema{
			"value": &framework.FieldSchema{Type: framework.TypeString},
		},
		ExistenceCheck: b.pathTestingExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathTestingRead,
			logical.CreateOperation: b.pathTestingCreate,
		},
	}
}

func (b *backend) pathTestingRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"value": data.Get("value"),
		},
	}, nil
}

func (b *backend) pathTestingCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	val := data.Get("value").(string)

	entry := &logical.StorageEntry{
		Key:   "test/ing",
		Value: []byte(val),
	}

	s := req.Storage
	err := s.Put(entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": data.Get("value"),
		},
	}, nil
}

func (b *backend) pathTestingExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	return false, nil
}

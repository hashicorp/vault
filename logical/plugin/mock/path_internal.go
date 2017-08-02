package mock

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathInternal(b *backend) *framework.Path {
	return &framework.Path{
		Pattern:        "internal",
		Fields:         map[string]*framework.FieldSchema{},
		ExistenceCheck: b.pathTestingExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTestingReadInternal,
		},
	}
}

func (b *backend) pathTestingReadInternal(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"value": b.internal,
		},
	}, nil

}

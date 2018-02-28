package mock

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// pathRaw is used to test raw responses.
func pathRaw(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "raw",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRawRead,
		},
	}
}

func (b *backend) pathRawRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "text/plain",
			logical.HTTPRawBody:     []byte("Response"),
			logical.HTTPStatusCode:  200,
		},
	}, nil

}

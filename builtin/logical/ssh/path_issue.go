package ssh

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameWithAtRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssue,
		},
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	fmt.Println("Hello world")
	return nil, nil
}

package redis

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCreds(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type: framework.TypeString,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.credsRead,
		},
	}
}

func (b *backend) credsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	role, err := getRole(ctx, req.Storage, name)
	if err != nil {
		return logical.ErrorResponse("failed to get role: %s", err), nil
	}
	if role == nil {
		return logical.ErrorResponse("no role named %q found", name), nil
	}

	client, err := b.Client(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("failed to get Redis client: %s", err), nil
	}

	username := fmt.Sprintf("%s-%s", req.DisplayName, req.ID)

	args := []interface{}{"ACL", "SETUSER", username}
	for _, arg := range role.Rules {
		args = append(args, arg)
	}

	password, err := random.DefaultStringGenerator.Generate(ctx, nil)
	if err != nil {
		return logical.ErrorResponse("failed to generate password: %s", err), nil
	}
	args = append(args, "#"+hash(password))

	if _, err := client.Do(ctx, args...).Result(); err != nil {
		return logical.ErrorResponse("failed to create user: %s", err), nil
	}

	resp := b.Secret("creds").Response(
		map[string]interface{}{
			"username": username,
			"password": password,
		},
		map[string]interface{}{
			"username": username,
		},
	)

	return resp, nil
}

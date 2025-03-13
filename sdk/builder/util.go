package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// tokenRevoke removes the token from the Vault storage API and calls the client to revoke the token
func (gb *GenericBackend[CC, C, R]) clientAndRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*C, *R, error) {
	roleName := d.Get("name").(string)

	roleEntry, err := gb.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return nil, nil, errors.New("error retrieving role: role is nil")
	}

	client, err := gb.getClient(ctx, req.Storage)
	if err != nil {
		return nil, nil, err
	}
	return client, roleEntry, nil
}

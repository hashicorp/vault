package ssh

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/strutil"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Structure to hold roles that are allowed to accept any IP address.
type zeroAddressRoles struct {
	Roles []string `json:"roles" mapstructure:"roles"`
}

func pathConfigZeroAddress(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/zeroaddress",
		Fields: map[string]*framework.FieldSchema{
			"roles": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `[Required] Comma separated list of role names which
				allows credentials to be requested for any IP address. CIDR blocks
				previously registered under these roles will be ignored.`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigZeroAddressWrite,
			logical.ReadOperation:   b.pathConfigZeroAddressRead,
			logical.DeleteOperation: b.pathConfigZeroAddressDelete,
		},
		HelpSynopsis:    pathConfigZeroAddressSyn,
		HelpDescription: pathConfigZeroAddressDesc,
	}
}

func (b *backend) pathConfigZeroAddressDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "config/zeroaddress")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConfigZeroAddressRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getZeroAddressRoles(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"roles": entry.Roles,
		},
	}, nil
}

func (b *backend) pathConfigZeroAddressWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roles := d.Get("roles").([]string)
	if len(roles) == 0 {
		return logical.ErrorResponse("Missing roles"), nil
	}

	// Check if the roles listed actually exist in the backend
	for _, item := range roles {
		role, err := b.getRole(ctx, req.Storage, item)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse(fmt.Sprintf("Role %q does not exist", item)), nil
		}
	}

	err := b.putZeroAddressRoles(ctx, req.Storage, roles)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Stores the given list of roles at zeroaddress endpoint
func (b *backend) putZeroAddressRoles(ctx context.Context, s logical.Storage, roles []string) error {
	entry, err := logical.StorageEntryJSON("config/zeroaddress", &zeroAddressRoles{
		Roles: roles,
	})
	if err != nil {
		return err
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

// Retrieves the list of roles from the zeroaddress endpoint.
func (b *backend) getZeroAddressRoles(ctx context.Context, s logical.Storage) (*zeroAddressRoles, error) {
	entry, err := s.Get(ctx, "config/zeroaddress")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result zeroAddressRoles
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Removes a role from the list of roles present in config/zeroaddress path
func (b *backend) removeZeroAddressRole(ctx context.Context, s logical.Storage, roleName string) error {
	zeroAddressEntry, err := b.getZeroAddressRoles(ctx, s)
	if err != nil {
		return err
	}
	if zeroAddressEntry == nil {
		return nil
	}

	zeroAddressEntry.Roles = strutil.StrListDelete(zeroAddressEntry.Roles, roleName)

	return b.putZeroAddressRoles(ctx, s, zeroAddressEntry.Roles)
}

const pathConfigZeroAddressSyn = `
Assign zero address as default CIDR block for select roles.
`

const pathConfigZeroAddressDesc = `
Administrator can choose to make a select few registered roles to accept any IP
address, overriding the CIDR blocks registered during creation of roles. This
doesn't mean that the credentials are created for any IP address. Clients who
have access to these roles are trusted to make valid requests. Access to these
roles should be controlled using Vault policies. It is recommended that all the
roles that are allowed to accept any IP address should have an explicit policy
of deny for unintended clients.

This is a root authenticated endpoint. If backend is mounted at 'ssh' then use
the endpoint 'ssh/config/zeroaddress' to provide the list of allowed roles.
After mounting the backend, use 'path-help' for additional information.
`

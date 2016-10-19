package ssh

import (
	"fmt"
	"strings"

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
				Type: framework.TypeString,
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

func (b *backend) pathConfigZeroAddressDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("config/zeroaddress")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConfigZeroAddressRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getZeroAddressRoles(req.Storage)
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

func (b *backend) pathConfigZeroAddressWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleNames := d.Get("roles").(string)
	if roleNames == "" {
		return logical.ErrorResponse("Missing roles"), nil
	}

	// Check if the roles listed actually exist in the backend
	roles := strings.Split(roleNames, ",")
	for _, item := range roles {
		role, err := b.getRole(req.Storage, item)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse(fmt.Sprintf("Role %q does not exist", item)), nil
		}
	}

	err := b.putZeroAddressRoles(req.Storage, roles)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Stores the given list of roles at zeroaddress endpoint
func (b *backend) putZeroAddressRoles(s logical.Storage, roles []string) error {
	entry, err := logical.StorageEntryJSON("config/zeroaddress", &zeroAddressRoles{
		Roles: roles,
	})
	if err != nil {
		return err
	}
	if err := s.Put(entry); err != nil {
		return err
	}
	return nil
}

// Retrieves the list of roles from the zeroaddress endpoint.
func (b *backend) getZeroAddressRoles(s logical.Storage) (*zeroAddressRoles, error) {
	entry, err := s.Get("config/zeroaddress")
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
func (b *backend) removeZeroAddressRole(s logical.Storage, roleName string) error {
	zeroAddressEntry, err := b.getZeroAddressRoles(s)
	if err != nil {
		return err
	}
	if zeroAddressEntry == nil {
		return nil
	}

	err = zeroAddressEntry.remove(roleName)
	if err != nil {
		return err
	}

	return b.putZeroAddressRoles(s, zeroAddressEntry.Roles)
}

// Removes a given role from the comma separated string
func (r *zeroAddressRoles) remove(roleName string) error {
	var index int
	for i, role := range r.Roles {
		if role == roleName {
			index = i
			break
		}
	}
	length := len(r.Roles)
	if index >= length || index < 0 {
		return fmt.Errorf("invalid index [%d]", index)
	}
	// If slice has zero or one item, remove the item by setting slice to nil.
	if length < 2 {
		r.Roles = nil
		return nil
	}

	// Last item to be deleted
	if length-1 == index {
		r.Roles = r.Roles[:length-1]
		return nil
	}

	// Delete the item by appending all items except the one at index
	r.Roles = append(r.Roles[:index], r.Roles[index+1:]...)
	return nil
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

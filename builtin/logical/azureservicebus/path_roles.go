package azureservicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"sas_policy_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Shared access key policy name",
			},
			"sas_policy_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Shared access primary key",
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Optional: Lease time of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := b.Role(req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"sas_policy_name": role.SASPolicyName,
			"ttl":             role.TTL,
		},
	}, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rolename := data.Get("name").(string)
	policyname := data.Get("sas_policy_name").(string)
	policykey := data.Get("sas_policy_key").(string)

	var ttl time.Duration
	var err error
	ttlraw := data.Get("ttl").(string)
	if ttlraw != "" {
		ttl, err = time.ParseDuration(ttlraw)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid lease time: %s", err)), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+rolename, &roleEntry{
		SASPolicyName: policyname,
		SASPolicyKey:  policykey,
		TTL:           ttl,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleEntry struct {
	SASPolicyName string        `json:"sas_policy_name"`
	SASPolicyKey  string        `json:"sas_policy_key"`
	TTL           time.Duration `json:"ttl"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

Each role corresponds to an existing SAS policy in the Service Bus resource.
Roles can be configured with a role-specific lease time.
`

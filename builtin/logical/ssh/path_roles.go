package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathRoles\n")
	return &framework.Path{
		Pattern: "roles/(?P<name>\\w+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "String representing the policy for the role. See help for more info.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.WriteOperation:  b.pathRoleWrite,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleRead\n")
	entry, err := req.Storage.Get("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"policy": string(entry.Value),
		},
	}, nil
}

func (b *backend) pathRoleWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleWrite\n")
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(d.Get("policy").(string))); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error compacting policy: %s", err)), nil
	}

	err := req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + d.Get("name").(string),
		Value: buf.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleDelete\n")
	err := req.Storage.Delete("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.
`

package aws

import (
	"bytes"
	"encoding/json"
	"fmt"

	"errors"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"arn": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ARN Reference to a managed policy",
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IAM policy document",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: pathRolesDelete,
			logical.ReadOperation:   pathRolesRead,
			logical.UpdateOperation: pathRolesWrite,
		},

		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("policy/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(entries), nil
}

func pathRolesDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func pathRolesRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	val := string(entry.Value)
	if strings.HasPrefix(val, "arn:") {
		return &logical.Response{
			Data: map[string]interface{}{
				"arn": val,
			},
		}, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"policy": val,
		},
	}, nil
}

func useInlinePolicy(d *framework.FieldData) (bool, error) {
	bp := d.Get("policy").(string) != ""
	ba := d.Get("arn").(string) != ""

	if !bp && !ba {
		return false, errors.New("Either policy or arn must be provided")
	}
	if bp && ba {
		return false, errors.New("Only one of policy or arn should be provided")
	}
	return bp, nil
}

func pathRolesWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var buf bytes.Buffer

	uip, err := useInlinePolicy(d)
	if err != nil {
		return nil, err
	}

	if uip {
		if err := json.Compact(&buf, []byte(d.Get("policy").(string))); err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error compacting policy: %s", err)), nil
		}
		// Write the policy into storage
		err := req.Storage.Put(&logical.StorageEntry{
			Key:   "policy/" + d.Get("name").(string),
			Value: buf.Bytes(),
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Write the arn ref into storage
		err := req.Storage.Put(&logical.StorageEntry{
			Key:   "policy/" + d.Get("name").(string),
			Value: []byte(d.Get("arn").(string)),
		})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

const pathListRolesHelpSyn = `List the existing roles in this backend`

const pathListRolesHelpDesc = `Roles will be listed by the role name.`

const pathRolesHelpSyn = `
Read, write and reference IAM policies that access keys can be made for.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create access keys. These roles are associated with IAM policies that
map directly to the route to read the access keys. For example, if the
backend is mounted at "aws" and you create a role at "aws/roles/deploy"
then a user could request access credentials at "aws/creds/deploy".

You can either supply a user inline policy (via the policy argument), or
provide a reference to an existing AWS policy by supplying the full arn
reference (via the arn argument). Inline user policies written are normal
IAM policies. Vault will not attempt to parse these except to validate
that they're basic JSON. No validation is performed on arn references.

To validate the keys, attempt to read an access key after writing the policy.
`

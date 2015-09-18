package transit

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/config",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"min_decryption_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `If set, the minimum version of the key allowed
to be decrypted.`,
			},

			"deletion_allowed": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Whether to allow deletion of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Check if the policy already exists
	policy, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing role named %s could be found", name)),
			logical.ErrInvalidRequest
	}

	persistNeeded := false

	minDecryptionVersion := d.Get("min_decryption_version").(int)
	if minDecryptionVersion != 0 &&
		minDecryptionVersion != policy.MinDecryptionVersion {
		policy.MinDecryptionVersion = minDecryptionVersion
		persistNeeded = true
	}

	allowDeletionInt, ok := d.GetOk("deletion_allowed")
	if ok {
		allowDeletion := allowDeletionInt.(bool)
		if allowDeletion != policy.DeletionAllowed {
			policy.DeletionAllowed = allowDeletion
			persistNeeded = true
		}
	}

	if !persistNeeded {
		return nil, nil
	}

	return nil, policy.Persist(req.Storage, name)
}

const pathConfigHelpSyn = `Configure a named encryption key`

const pathConfigHelpDesc = `
This path is used to configure the named key. Currently, this
supports adjusting the minimum version of the key allowed to
be used for decryption via the min_decryption_version paramter.
`

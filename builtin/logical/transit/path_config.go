package transit

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathConfig() *framework.Path {
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
			logical.UpdateOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Check if the policy already exists
	lp, err := b.policies.getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if lp == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing role named %s could be found", name)),
			logical.ErrInvalidRequest
	}

	lp.Lock()
	defer lp.Unlock()

	// Verify if wasn't deleted before we grabbed the lock
	if lp.policy == nil {
		return nil, fmt.Errorf("no existing role named %s could be found", name)
	}

	resp := &logical.Response{}

	persistNeeded := false

	minDecryptionVersionRaw, ok := d.GetOk("min_decryption_version")
	if ok {
		minDecryptionVersion := minDecryptionVersionRaw.(int)

		if minDecryptionVersion < 0 {
			return logical.ErrorResponse("min decryption version cannot be negative"), nil
		}

		if minDecryptionVersion == 0 {
			minDecryptionVersion = 1
			resp.AddWarning("since Vault 0.3, transit key numbering starts at 1; forcing minimum to 1")
		}

		if minDecryptionVersion > 0 &&
			minDecryptionVersion != lp.policy.MinDecryptionVersion {
			if minDecryptionVersion > lp.policy.LatestVersion {
				return logical.ErrorResponse(
					fmt.Sprintf("cannot set min decryption version of %d, latest key version is %d", minDecryptionVersion, lp.policy.LatestVersion)), nil
			}
			lp.policy.MinDecryptionVersion = minDecryptionVersion
			persistNeeded = true
		}
	}

	allowDeletionInt, ok := d.GetOk("deletion_allowed")
	if ok {
		allowDeletion := allowDeletionInt.(bool)
		if allowDeletion != lp.policy.DeletionAllowed {
			lp.policy.DeletionAllowed = allowDeletion
			persistNeeded = true
		}
	}

	// Add this as a guard here before persisting since we now require the min
	// decryption version to start at 1; even if it's not explicitly set here,
	// force the upgrade
	if lp.policy.MinDecryptionVersion == 0 {
		lp.policy.MinDecryptionVersion = 1
		persistNeeded = true
	}

	if !persistNeeded {
		return nil, nil
	}

	return resp, lp.policy.Persist(req.Storage)
}

const pathConfigHelpSyn = `Configure a named encryption key`

const pathConfigHelpDesc = `
This path is used to configure the named key. Currently, this
supports adjusting the minimum version of the key allowed to
be used for decryption via the min_decryption_version paramter.
`

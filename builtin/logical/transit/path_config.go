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
	lp, err := b.policies.getPolicy(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if lp == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing key named %s could be found", name)),
			logical.ErrInvalidRequest
	}

	// Store some values so we can detect if the policy changed after locking
	lp.RLock()
	currDeletionAllowed := lp.Policy().DeletionAllowed
	currMinDecryptionVersion := lp.Policy().MinDecryptionVersion
	lp.RUnlock()

	// Hold both locks since we want to ensure the policy doesn't change from
	// underneath us
	b.policies.Lock()
	defer b.policies.Unlock()
	lp.Lock()
	defer lp.Unlock()

	// Refresh in case it's changed since before we grabbed the lock
	lp, err = b.policies.refreshPolicy(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if lp == nil {
		return nil, fmt.Errorf("error finding key %s after locking for changes", name)
	}

	// Verify if wasn't deleted before we grabbed the lock
	if lp.Policy() == nil {
		return nil, fmt.Errorf("no existing key named %s could be found", name)
	}

	resp := &logical.Response{}

	// Check for anything to have been updated since we got the write lock
	if currDeletionAllowed != lp.Policy().DeletionAllowed ||
		currMinDecryptionVersion != lp.Policy().MinDecryptionVersion {
		resp.AddWarning("key configuration has changed since this endpoint was called, not updating")
		return resp, nil
	}

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
			minDecryptionVersion != lp.Policy().MinDecryptionVersion {
			if minDecryptionVersion > lp.Policy().LatestVersion {
				return logical.ErrorResponse(
					fmt.Sprintf("cannot set min decryption version of %d, latest key version is %d", minDecryptionVersion, lp.Policy().LatestVersion)), nil
			}
			lp.Policy().MinDecryptionVersion = minDecryptionVersion
			persistNeeded = true
		}
	}

	allowDeletionInt, ok := d.GetOk("deletion_allowed")
	if ok {
		allowDeletion := allowDeletionInt.(bool)
		if allowDeletion != lp.Policy().DeletionAllowed {
			lp.Policy().DeletionAllowed = allowDeletion
			persistNeeded = true
		}
	}

	// Add this as a guard here before persisting since we now require the min
	// decryption version to start at 1; even if it's not explicitly set here,
	// force the upgrade
	if lp.Policy().MinDecryptionVersion == 0 {
		lp.Policy().MinDecryptionVersion = 1
		persistNeeded = true
	}

	if !persistNeeded {
		return nil, nil
	}

	return resp, lp.Policy().Persist(req.Storage)
}

const pathConfigHelpSyn = `Configure a named encryption key`

const pathConfigHelpDesc = `
This path is used to configure the named key. Currently, this
supports adjusting the minimum version of the key allowed to
be used for decryption via the min_decryption_version paramter.
`

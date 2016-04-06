package aws

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathWhitelistIdentityTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "whitelist/identity/tidy$",
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 259200,
				Description: `The amount of extra time that must have passed beyond the identity's
expiration, before it is removed from the backend storage.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathWhitelistIdentityTidyUpdate,
		},

		HelpSynopsis:    pathWhitelistIdentityTidySyn,
		HelpDescription: pathWhitelistIdentityTidyDesc,
	}
}

// pathWhitelistIdentityTidyUpdate is used to delete entries in the whitelist that are expired.
func (b *backend) pathWhitelistIdentityTidyUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	safety_buffer := data.Get("safety_buffer").(int)

	bufferDuration := time.Duration(safety_buffer) * time.Second

	identities, err := req.Storage.List("whitelist/identity/")
	if err != nil {
		return nil, err
	}

	for _, instanceID := range identities {
		identityEntry, err := req.Storage.Get("whitelist/identity/" + instanceID)
		if err != nil {
			return nil, fmt.Errorf("error fetching identity of instanceID %s: %s", instanceID, err)
		}

		if identityEntry == nil {
			return nil, fmt.Errorf("identity entry for instanceID %s is nil", instanceID)
		}

		if identityEntry.Value == nil || len(identityEntry.Value) == 0 {
			return nil, fmt.Errorf("found identity entry for instanceID %s but actual identity is empty", instanceID)
		}

		var result whitelistIdentity
		if err := identityEntry.DecodeJSON(&result); err != nil {
			return nil, err
		}

		if time.Now().After(result.ExpirationTime.Add(bufferDuration)) {
			if err := req.Storage.Delete("whitelist/identity" + instanceID); err != nil {
				return nil, fmt.Errorf("error deleting identity of instanceID %s from storage: %s", instanceID, err)
			}
		}
	}

	return nil, nil
}

const pathWhitelistIdentityTidySyn = `
Clean-up the whitelisted instance identity entries.
`

const pathWhitelistIdentityTidyDesc = `
When an instance identity is whitelisted, the expiration time of the whitelist
entry is set to the least amont 'max_ttl' of the registered AMI, 'max_ttl' of the
role tag and 'max_ttl' of the backend mount.

When this endpoint is invoked all the entries that are expired will be deleted.

A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`

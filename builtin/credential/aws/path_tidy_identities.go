package aws

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTidyIdentities(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy/identities$",
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 259200,
				Description: `The amount of extra time that must have passed beyond the identity's
expiration, before it is removed from the backend storage.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidyIdentitiesUpdate,
		},

		HelpSynopsis:    pathTidyIdentitiesSyn,
		HelpDescription: pathTidyIdentitiesDesc,
	}
}

// tidyWhitelistIdentity is used to delete entries in the whitelist that are expired.
func tidyWhitelistIdentity(s logical.Storage, safety_buffer int) error {
	bufferDuration := time.Duration(safety_buffer) * time.Second

	identities, err := s.List("whitelist/identity/")
	if err != nil {
		return err
	}

	for _, instanceID := range identities {
		identityEntry, err := s.Get("whitelist/identity/" + instanceID)
		if err != nil {
			return fmt.Errorf("error fetching identity of instanceID %s: %s", instanceID, err)
		}

		if identityEntry == nil {
			return fmt.Errorf("identity entry for instanceID %s is nil", instanceID)
		}

		if identityEntry.Value == nil || len(identityEntry.Value) == 0 {
			return fmt.Errorf("found identity entry for instanceID %s but actual identity is empty", instanceID)
		}

		var result whitelistIdentity
		if err := identityEntry.DecodeJSON(&result); err != nil {
			return err
		}

		if time.Now().UTC().After(result.ExpirationTime.Add(bufferDuration)) {
			if err := s.Delete("whitelist/identity" + instanceID); err != nil {
				return fmt.Errorf("error deleting identity of instanceID %s from storage: %s", instanceID, err)
			}
		}
	}

	return nil
}

// pathTidyIdentitiesUpdate is used to delete entries in the whitelist that are expired.
func (b *backend) pathTidyIdentitiesUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, tidyWhitelistIdentity(req.Storage, data.Get("safety_buffer").(int))
}

const pathTidyIdentitiesSyn = `
Clean-up the whitelisted instance identity entries.
`

const pathTidyIdentitiesDesc = `
When an instance identity is whitelisted, the expiration time of the whitelist
entry is set based on the least 'max_ttl' value set on: AMI entry, the role tag
and the backend's mount.

When this endpoint is invoked, all the entries that are expired will be deleted.

A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`

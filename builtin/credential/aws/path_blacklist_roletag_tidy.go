package aws

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathBlacklistRoleTagTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "blacklist/roletag/tidy$",
		Fields: map[string]*framework.FieldSchema{
			"safety_buffer": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 259200, // 72h
				Description: `The amount of extra time that must have passed beyond the roletag's
expiration, before it is removed from the backend storage.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathBlacklistRoleTagTidyUpdate,
		},

		HelpSynopsis:    pathBlacklistRoleTagTidySyn,
		HelpDescription: pathBlacklistRoleTagTidyDesc,
	}
}

// pathBlacklistRoleTagTidyUpdate is used to clean-up the entries in the role tag blacklist.
func (b *backend) pathBlacklistRoleTagTidyUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	// safety_buffer is an optional parameter.
	safety_buffer := data.Get("safety_buffer").(int)
	bufferDuration := time.Duration(safety_buffer) * time.Second

	tags, err := req.Storage.List("blacklist/roletag/")
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		tagEntry, err := req.Storage.Get("blacklist/roletag/" + tag)
		if err != nil {
			return nil, fmt.Errorf("error fetching tag %s: %s", tag, err)
		}

		if tagEntry == nil {
			return nil, fmt.Errorf("tag entry for tag %s is nil", tag)
		}

		if tagEntry.Value == nil || len(tagEntry.Value) == 0 {
			return nil, fmt.Errorf("found entry for tag %s but actual tag is empty", tag)
		}

		var result roleTagBlacklistEntry
		if err := tagEntry.DecodeJSON(&result); err != nil {
			return nil, err
		}

		if time.Now().After(result.ExpirationTime.Add(bufferDuration)) {
			if err := req.Storage.Delete("blacklist/roletag" + tag); err != nil {
				return nil, fmt.Errorf("error deleting tag %s from storage: %s", tag, err)
			}
		}
	}

	return nil, nil
}

const pathBlacklistRoleTagTidySyn = `
Clean-up the blacklisted role tag entries.
`

const pathBlacklistRoleTagTidyDesc = `
When a role tag is blacklisted, the expiration time of the blacklist entry is
determined by the 'max_ttl' present in the role tag. If 'max_ttl' is not provided
in the role tag, the backend mount's 'max_ttl' value will be used to determine
the expiration time of the blacklist entry.

When this endpoint is invoked all the entries that are expired will be deleted.

A 'safety_buffer' (duration in seconds) can be provided, to ensure deletion of
only those entries that are expired before 'safety_buffer' seconds. 
`

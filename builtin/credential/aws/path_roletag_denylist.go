// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathRoletagDenyList() *framework.Path {
	return &framework.Path{
		Pattern: "roletag-denylist/(?P<role_tag>.*)",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "role-tag-deny-list",
		},

		Fields: map[string]*framework.FieldSchema{
			"role_tag": {
				Type: framework.TypeString,
				Description: `Role tag to be deny listed. The tag can be supplied as-is. In order
to avoid any encoding problems, it can be base64 encoded.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRoletagDenyListUpdate,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoletagDenyListRead,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathRoletagDenyListDelete,
			},
		},

		HelpSynopsis:    pathRoletagBlacklistSyn,
		HelpDescription: pathRoletagBlacklistDesc,
	}
}

// Path to list all the deny listed tags.
func (b *backend) pathListRoletagDenyList() *framework.Path {
	return &framework.Path{
		Pattern: "roletag-denylist/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "role-tag-deny-lists",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathRoletagDenyListsList,
			},
		},

		HelpSynopsis:    pathListRoletagDenyListHelpSyn,
		HelpDescription: pathListRoletagDenyListHelpDesc,
	}
}

// Lists all the deny listed role tags.
func (b *backend) pathRoletagDenyListsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.denyListMutex.RLock()
	defer b.denyListMutex.RUnlock()

	tags, err := req.Storage.List(ctx, denyListRoletagStorage)
	if err != nil {
		return nil, err
	}

	// Tags are base64 encoded before indexing to avoid problems
	// with the path separators being present in the tag.
	// Reverse it before returning the list response.
	for i, keyB64 := range tags {
		if key, err := base64.StdEncoding.DecodeString(keyB64); err != nil {
			return nil, err
		} else {
			// Overwrite the result with the decoded string.
			tags[i] = string(key)
		}
	}
	return logical.ListResponse(tags), nil
}

// Fetch an entry from the role tag deny list for a given tag.
// This method takes a role tag in its original form and not a base64 encoded form.
func (b *backend) lockedDenyLististRoleTagEntry(ctx context.Context, s logical.Storage, tag string) (*roleTagBlacklistEntry, error) {
	b.denyListMutex.RLock()
	defer b.denyListMutex.RUnlock()

	return b.nonLockedDenyListRoleTagEntry(ctx, s, tag)
}

func (b *backend) nonLockedDenyListRoleTagEntry(ctx context.Context, s logical.Storage, tag string) (*roleTagBlacklistEntry, error) {
	entry, err := s.Get(ctx, denyListRoletagStorage+base64.StdEncoding.EncodeToString([]byte(tag)))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleTagBlacklistEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Deletes an entry from the role tag deny list for a given tag.
func (b *backend) pathRoletagDenyListDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.denyListMutex.Lock()
	defer b.denyListMutex.Unlock()

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	return nil, req.Storage.Delete(ctx, denyListRoletagStorage+base64.StdEncoding.EncodeToString([]byte(tag)))
}

// If the given role tag is deny listed, returns the details of the deny list entry.
// Returns 'nil' otherwise.
func (b *backend) pathRoletagDenyListRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	entry, err := b.lockedDenyLististRoleTagEntry(ctx, req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"creation_time":   entry.CreationTime.Format(time.RFC3339Nano),
			"expiration_time": entry.ExpirationTime.Format(time.RFC3339Nano),
		},
	}, nil
}

// pathRoletagDenyListUpdate is used to deny list a given role tag.
// Before a role tag is added to the deny list, the correctness of the plaintext part
// in the role tag is verified using the associated HMAC.
func (b *backend) pathRoletagDenyListUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// The role_tag value provided, optionally can be base64 encoded.
	tagInput := data.Get("role_tag").(string)
	if tagInput == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	tag := ""

	// Try to base64 decode the value.
	tagBytes, err := base64.StdEncoding.DecodeString(tagInput)
	if err != nil {
		// If the decoding failed, use the value as-is.
		tag = tagInput
	} else {
		// If the decoding succeeded, use the decoded value.
		tag = string(tagBytes)
	}

	// Parse and verify the role tag from string form to a struct form and verify it.
	rTag, err := b.parseAndVerifyRoleTagValue(ctx, req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if rTag == nil {
		return logical.ErrorResponse("failed to verify the role tag and parse it"), nil
	}

	// Get the entry for the role mentioned in the role tag.
	roleEntry, err := b.role(ctx, req.Storage, rTag.Role)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse("role entry not found"), nil
	}

	b.denyListMutex.Lock()
	defer b.denyListMutex.Unlock()

	// Check if the role tag is already deny listed. If yes, update it.
	blEntry, err := b.nonLockedDenyListRoleTagEntry(ctx, req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if blEntry == nil {
		blEntry = &roleTagBlacklistEntry{}
	}

	currentTime := time.Now()

	// Check if this is a creation of deny list entry.
	if blEntry.CreationTime.IsZero() {
		// Set the creation time for the deny list entry.
		// This should not be updated after setting it once.
		// If deny list operation is invoked more than once, only update the expiration time.
		blEntry.CreationTime = currentTime
	}

	// Decide the expiration time based on the max_ttl values. Since this is
	// restricting access, use the greatest duration, not the least.
	maxDur := rTag.MaxTTL
	if roleEntry.TokenMaxTTL > maxDur {
		maxDur = roleEntry.TokenMaxTTL
	}
	if b.System().MaxLeaseTTL() > maxDur {
		maxDur = b.System().MaxLeaseTTL()
	}

	blEntry.ExpirationTime = currentTime.Add(maxDur)

	entry, err := logical.StorageEntryJSON(denyListRoletagStorage+base64.StdEncoding.EncodeToString([]byte(tag)), blEntry)
	if err != nil {
		return nil, err
	}

	// Store the deny list entry.
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleTagBlacklistEntry struct {
	CreationTime   time.Time `json:"creation_time"`
	ExpirationTime time.Time `json:"expiration_time"`
}

const pathRoletagBlacklistSyn = `
Blacklist a previously created role tag.
`

const pathRoletagBlacklistDesc = `
Add a role tag to the deny list so that it cannot be used by any EC2 instance to perform further
logins. This can be used if the role tag is suspected or believed to be possessed by
an unintended party.

By default, a cron task will periodically look for expired entries in the deny list
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/roletags' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/roletags' endpoint.

Also note that delete operation is supported on this endpoint to remove specific
entries from the deny list.
`

const pathListRoletagDenyListHelpSyn = `
Lists the deny list role tags.
`

const pathListRoletagDenyListHelpDesc = `
Lists all the entries present in the deny list. This will show both the valid
entries and the expired entries in the deny list. Use 'tidy/roletags' endpoint
to clean-up the deny list of role tags based on expiration time.
`

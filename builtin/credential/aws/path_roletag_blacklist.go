package awsauth

import (
	"encoding/base64"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoletagBlacklist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roletag-blacklist/(?P<role_tag>.*)",
		Fields: map[string]*framework.FieldSchema{
			"role_tag": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Role tag to be blacklisted. The tag can be supplied as-is. In order
to avoid any encoding problems, it can be base64 encoded.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRoletagBlacklistUpdate,
			logical.ReadOperation:   b.pathRoletagBlacklistRead,
			logical.DeleteOperation: b.pathRoletagBlacklistDelete,
		},

		HelpSynopsis:    pathRoletagBlacklistSyn,
		HelpDescription: pathRoletagBlacklistDesc,
	}
}

// Path to list all the blacklisted tags.
func pathListRoletagBlacklist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roletag-blacklist/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoletagBlacklistsList,
		},

		HelpSynopsis:    pathListRoletagBlacklistHelpSyn,
		HelpDescription: pathListRoletagBlacklistHelpDesc,
	}
}

// Lists all the blacklisted role tags.
func (b *backend) pathRoletagBlacklistsList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.blacklistMutex.RLock()
	defer b.blacklistMutex.RUnlock()

	tags, err := req.Storage.List("blacklist/roletag/")
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

// Fetch an entry from the role tag blacklist for a given tag.
// This method takes a role tag in its original form and not a base64 encoded form.
func (b *backend) lockedBlacklistRoleTagEntry(s logical.Storage, tag string) (*roleTagBlacklistEntry, error) {
	b.blacklistMutex.RLock()
	defer b.blacklistMutex.RUnlock()

	return b.nonLockedBlacklistRoleTagEntry(s, tag)
}

func (b *backend) nonLockedBlacklistRoleTagEntry(s logical.Storage, tag string) (*roleTagBlacklistEntry, error) {
	entry, err := s.Get("blacklist/roletag/" + base64.StdEncoding.EncodeToString([]byte(tag)))
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

// Deletes an entry from the role tag blacklist for a given tag.
func (b *backend) pathRoletagBlacklistDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.blacklistMutex.Lock()
	defer b.blacklistMutex.Unlock()

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	return nil, req.Storage.Delete("blacklist/roletag/" + base64.StdEncoding.EncodeToString([]byte(tag)))
}

// If the given role tag is blacklisted, returns the details of the blacklist entry.
// Returns 'nil' otherwise.
func (b *backend) pathRoletagBlacklistRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	entry, err := b.lockedBlacklistRoleTagEntry(req.Storage, tag)
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

// pathRoletagBlacklistUpdate is used to blacklist a given role tag.
// Before a role tag is blacklisted, the correctness of the plaintext part
// in the role tag is verified using the associated HMAC.
func (b *backend) pathRoletagBlacklistUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

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
	rTag, err := b.parseAndVerifyRoleTagValue(req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if rTag == nil {
		return logical.ErrorResponse("failed to verify the role tag and parse it"), nil
	}

	// Get the entry for the role mentioned in the role tag.
	roleEntry, err := b.lockedAWSRole(req.Storage, rTag.Role)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse("role entry not found"), nil
	}

	b.blacklistMutex.Lock()
	defer b.blacklistMutex.Unlock()

	// Check if the role tag is already blacklisted. If yes, update it.
	blEntry, err := b.nonLockedBlacklistRoleTagEntry(req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if blEntry == nil {
		blEntry = &roleTagBlacklistEntry{}
	}

	currentTime := time.Now()

	// Check if this is a creation of blacklist entry.
	if blEntry.CreationTime.IsZero() {
		// Set the creation time for the blacklist entry.
		// This should not be updated after setting it once.
		// If blacklist operation is invoked more than once, only update the expiration time.
		blEntry.CreationTime = currentTime
	}

	// Decide the expiration time based on the max_ttl values. Since this is
	// restricting access, use the greatest duration, not the least.
	maxDur := rTag.MaxTTL
	if roleEntry.MaxTTL > maxDur {
		maxDur = roleEntry.MaxTTL
	}
	if b.System().MaxLeaseTTL() > maxDur {
		maxDur = b.System().MaxLeaseTTL()
	}

	blEntry.ExpirationTime = currentTime.Add(maxDur)

	entry, err := logical.StorageEntryJSON("blacklist/roletag/"+base64.StdEncoding.EncodeToString([]byte(tag)), blEntry)
	if err != nil {
		return nil, err
	}

	// Store the blacklist entry.
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleTagBlacklistEntry struct {
	CreationTime   time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`
	ExpirationTime time.Time `json:"expiration_time" structs:"expiration_time" mapstructure:"expiration_time"`
}

const pathRoletagBlacklistSyn = `
Blacklist a previously created role tag.
`

const pathRoletagBlacklistDesc = `
Blacklist a role tag so that it cannot be used by any EC2 instance to perform further
logins. This can be used if the role tag is suspected or believed to be possessed by
an unintended party.

By default, a cron task will periodically look for expired entries in the blacklist
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/roletags' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/roletags' endpoint.

Also note that delete operation is supported on this endpoint to remove specific
entries from the blacklist.
`

const pathListRoletagBlacklistHelpSyn = `
Lists the blacklisted role tags.
`

const pathListRoletagBlacklistHelpDesc = `
Lists all the entries present in the blacklist. This will show both the valid
entries and the expired entries in the blacklist. Use 'tidy/roletags' endpoint
to clean-up the blacklist of role tags based on expiration time.
`

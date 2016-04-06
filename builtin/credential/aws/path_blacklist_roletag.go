package aws

import (
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathBlacklistRoleTag(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "blacklist/roletag$",
		Fields: map[string]*framework.FieldSchema{
			"role_tag": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Role tag that needs be blacklisted",
			},
		},

		ExistenceCheck: b.pathBlacklistRoleTagExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathBlacklistRoleTagUpdate,
			logical.ReadOperation:   b.pathBlacklistRoleTagRead,
			logical.DeleteOperation: b.pathBlacklistRoleTagDelete,
		},

		HelpSynopsis:    pathBlacklistRoleTagSyn,
		HelpDescription: pathBlacklistRoleTagDesc,
	}
}

// Path to list all the blacklisted tags.
func pathListBlacklistRoleTags(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "blacklist/roletags/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathBlacklistRoleTagsList,
		},

		HelpSynopsis:    pathListBlacklistRoleTagsHelpSyn,
		HelpDescription: pathListBlacklistRoleTagsHelpDesc,
	}
}

// Lists all the blacklisted role tags.
func (b *backend) pathBlacklistRoleTagsList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	tags, err := req.Storage.List("blacklist/roletag/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(tags), nil
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
//
// A role should be allowed to be blacklisted even if it was prevously blacklisted.
func (b *backend) pathBlacklistRoleTagExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	return true, nil
}

// Fetch an entry from the role tag blacklist for a given tag.
func blacklistRoleTagEntry(s logical.Storage, tag string) (*roleTagBlacklistEntry, error) {
	entry, err := s.Get("blacklist/roletag/" + tag)
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
func (b *backend) pathBlacklistRoleTagDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	err := req.Storage.Delete("blacklist/roletag/" + tag)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// If the given role tag is blacklisted, returns the details of the blacklist entry.
// Returns 'nil' otherwise.
func (b *backend) pathBlacklistRoleTagRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	entry, err := blacklistRoleTagEntry(req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"creation_time":   entry.CreationTime,
			"expiration_time": entry.ExpirationTime,
		},
	}, nil
}

// pathBlacklistRoleTagUpdate is used to blacklist a given role tag.
// Before a role tag is blacklisted, the correctness of the plaintext part
// in the role tag is verified using the associated HMAC.
func (b *backend) pathBlacklistRoleTagUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	tag := data.Get("role_tag").(string)
	if tag == "" {
		return logical.ErrorResponse("missing role_tag"), nil
	}

	// Parse the role tag from string form to a struct form.
	rTag, err := parseRoleTagValue(tag)
	if err != nil {
		return nil, err
	}

	// Build the plaintext form of the role tag and verify the prepared
	// value using the HMAC.
	verified, err := verifyRoleTagValue(req.Storage, rTag)
	if err != nil {
		return nil, err
	}
	if !verified {
		return logical.ErrorResponse("role tag invalid"), nil
	}

	// Get the entry for the AMI used by the instance.
	imageEntry, err := awsImage(req.Storage, rTag.ImageID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return logical.ErrorResponse("image entry not found"), nil
	}

	blEntry, err := blacklistRoleTagEntry(req.Storage, tag)
	if err != nil {
		return nil, err
	}
	if blEntry == nil {
		blEntry = &roleTagBlacklistEntry{}
	}

	currentTime := time.Now()

	var epoch time.Time
	if blEntry.CreationTime.Equal(epoch) {
		// Set the creation time for the blacklist entry.
		// This should not be updated after setting it once.
		// If blacklist operation is invoked more than once, only update the expiration time.
		blEntry.CreationTime = currentTime
	}

	// If max_ttl is not set for the role tag, fall back on the mount's max_ttl.
	if rTag.MaxTTL == time.Duration(0) {
		rTag.MaxTTL = b.System().MaxLeaseTTL()
	}

	if imageEntry.MaxTTL > time.Duration(0) && rTag.MaxTTL > imageEntry.MaxTTL {
		rTag.MaxTTL = imageEntry.MaxTTL
	}

	// Expiration time is decided by the max_ttl value.
	blEntry.ExpirationTime = currentTime.Add(rTag.MaxTTL)

	entry, err := logical.StorageEntryJSON("blacklist/roletag/"+tag, blEntry)
	if err != nil {
		return nil, err
	}

	// Store it.
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleTagBlacklistEntry struct {
	CreationTime   time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`
	ExpirationTime time.Time `json:"expiration_time" structs:"expiration_time" mapstructure:"expiration_time"`
}

const pathBlacklistRoleTagSyn = `
Blacklist a previously created role tag.
`

const pathBlacklistRoleTagDesc = `
Blacklist a role tag so that it cannot be used by an EC2 instance to perform logins
in the future. This can be used if the role tag is suspected or believed to be possessed
by an unauthorized entity.

The entries in the blacklist are not automatically deleted. Although, they will have an
expiration time set on the entry. There is a separate endpoint 'blacklist/roletag/tidy',
that needs to be invoked to clean-up all the expired entries in the blacklist.
`

const pathListBlacklistRoleTagsHelpSyn = `
List the blacklisted role tags.
`

const pathListBlacklistRoleTagsHelpDesc = `
List all the entries present in the blacklist. This will show both the valid entries and
the expired entries in the blacklist. Use 'blacklist/roletag/tidy' endpoint to clean-up
the blacklist of role tags.
`

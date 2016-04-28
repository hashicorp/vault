package aws

import (
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathWhitelistIdentity(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "whitelist/identity/" + framework.GenericNameRegex("instance_id"),
		Fields: map[string]*framework.FieldSchema{
			"instance_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "EC2 instance ID. A successful login operation from an EC2 instance gets cached in this whitelist, keyed off of instance ID.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathWhitelistIdentityRead,
			logical.DeleteOperation: b.pathWhitelistIdentityDelete,
		},

		HelpSynopsis:    pathWhitelistIdentitySyn,
		HelpDescription: pathWhitelistIdentityDesc,
	}
}

func pathListWhitelistIdentities(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "whitelist/identities/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathWhitelistIdentitiesList,
		},

		HelpSynopsis:    pathListWhitelistIdentitiesHelpSyn,
		HelpDescription: pathListWhitelistIdentitiesHelpDesc,
	}
}

// pathWhitelistIdentitiesList is used to list all the instance IDs that are present
// in the identity whitelist. This will list both valid and expired entries.
func (b *backend) pathWhitelistIdentitiesList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	identities, err := req.Storage.List("whitelist/identity/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(identities), nil
}

// Fetch an item from the whitelist given an instance ID.
func whitelistIdentityEntry(s logical.Storage, instanceID string) (*whitelistIdentity, error) {
	entry, err := s.Get("whitelist/identity/" + instanceID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result whitelistIdentity
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stores an instance ID and the information required to validate further login/renewal attempts from
// the same instance ID.
func setWhitelistIdentityEntry(s logical.Storage, instanceID string, identity *whitelistIdentity) error {
	entry, err := logical.StorageEntryJSON("whitelist/identity/"+instanceID, identity)
	if err != nil {
		return err
	}

	if err := s.Put(entry); err != nil {
		return err
	}
	return nil
}

// pathWhitelistIdentityDelete is used to delete an entry from the identity whitelist given an instance ID.
func (b *backend) pathWhitelistIdentityDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	instanceID := data.Get("instance_id").(string)
	if instanceID == "" {
		return logical.ErrorResponse("missing instance_id"), nil
	}

	return nil, req.Storage.Delete("whitelist/identity/" + instanceID)
}

// pathWhitelistIdentityRead is used to view an entry in the identity whitelist given an instance ID.
func (b *backend) pathWhitelistIdentityRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := data.Get("instance_id").(string)
	if instanceID == "" {
		return logical.ErrorResponse("missing instance_id"), nil
	}

	entry, err := whitelistIdentityEntry(req.Storage, instanceID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(entry).Map(),
	}, nil
}

// Struct to represent each item in the identity whitelist.
type whitelistIdentity struct {
	AmiID                    string    `json:"ami_id" structs:"ami_id" mapstructure:"ami_id"`
	ClientNonce              string    `json:"client_nonce" structs:"client_nonce" mapstructure:"client_nonce"`
	CreationTime             time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`
	DisallowReauthentication bool      `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	PendingTime              string    `json:"pending_time" structs:"pending_time" mapstructure:"pending_time"`
	ExpirationTime           time.Time `json:"expiration_time" structs:"expiration_time" mapstructure:"expiration_time"`
	LastUpdatedTime          time.Time `json:"last_updated_time" structs:"last_updated_time" mapstructure:"last_updated_time"`
}

const pathWhitelistIdentitySyn = `
Read or delete entries in the identity whitelist.
`

const pathWhitelistIdentityDesc = `
Each login from an EC2 instance creates/updates an entry in the identity whitelist.

Entries in this list can be viewed or deleted using this endpoint.

By default, a cron task will periodically looks for expired entries in the whitelist
and delete them. The duration to periodically run this is one hour by default.
However, this can be configured using the 'config/tidy/identities' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/identities' endpoint.
`

const pathListWhitelistIdentitiesHelpSyn = `
List the items present in the identity whitelist.
`

const pathListWhitelistIdentitiesHelpDesc = `
The entries in the identity whitelist is keyed off of the EC2 instance IDs.
This endpoint lists all the entries present in the identity whitelist, both
expired and un-expired entries. Use 'tidy/identities' endpoint to clean-up
the whitelist of identities.
`

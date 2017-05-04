package awsauth

import (
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIdentityWhitelist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "identity-whitelist/" + framework.GenericNameRegex("instance_id"),
		Fields: map[string]*framework.FieldSchema{
			"instance_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `EC2 instance ID. A successful login operation from an EC2 instance
gets cached in this whitelist, keyed off of instance ID.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathIdentityWhitelistRead,
			logical.DeleteOperation: b.pathIdentityWhitelistDelete,
		},

		HelpSynopsis:    pathIdentityWhitelistSyn,
		HelpDescription: pathIdentityWhitelistDesc,
	}
}

func pathListIdentityWhitelist(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "identity-whitelist/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathWhitelistIdentitiesList,
		},

		HelpSynopsis:    pathListIdentityWhitelistHelpSyn,
		HelpDescription: pathListIdentityWhitelistHelpDesc,
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

// pathIdentityWhitelistDelete is used to delete an entry from the identity whitelist given an instance ID.
func (b *backend) pathIdentityWhitelistDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := data.Get("instance_id").(string)
	if instanceID == "" {
		return logical.ErrorResponse("missing instance_id"), nil
	}

	return nil, req.Storage.Delete("whitelist/identity/" + instanceID)
}

// pathIdentityWhitelistRead is used to view an entry in the identity whitelist given an instance ID.
func (b *backend) pathIdentityWhitelistRead(
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

	resp := &logical.Response{
		Data: structs.New(entry).Map(),
	}
	resp.Data["creation_time"] = entry.CreationTime.Format(time.RFC3339Nano)
	resp.Data["expiration_time"] = entry.ExpirationTime.Format(time.RFC3339Nano)
	resp.Data["last_updated_time"] = entry.LastUpdatedTime.Format(time.RFC3339Nano)

	return resp, nil
}

// Struct to represent each item in the identity whitelist.
type whitelistIdentity struct {
	Role                     string    `json:"role" structs:"role" mapstructure:"role"`
	ClientNonce              string    `json:"client_nonce" structs:"client_nonce" mapstructure:"client_nonce"`
	CreationTime             time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`
	DisallowReauthentication bool      `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	PendingTime              string    `json:"pending_time" structs:"pending_time" mapstructure:"pending_time"`
	ExpirationTime           time.Time `json:"expiration_time" structs:"expiration_time" mapstructure:"expiration_time"`
	LastUpdatedTime          time.Time `json:"last_updated_time" structs:"last_updated_time" mapstructure:"last_updated_time"`
}

const pathIdentityWhitelistSyn = `
Read or delete entries in the identity whitelist.
`

const pathIdentityWhitelistDesc = `
Each login from an EC2 instance creates/updates an entry in the identity whitelist.

Entries in this list can be viewed or deleted using this endpoint.

By default, a cron task will periodically look for expired entries in the whitelist
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/identities' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/identities' endpoint.
`

const pathListIdentityWhitelistHelpSyn = `
Lists the items present in the identity whitelist.
`

const pathListIdentityWhitelistHelpDesc = `
The entries in the identity whitelist is keyed off of the EC2 instance IDs.
This endpoint lists all the entries present in the identity whitelist, both
expired and un-expired entries. Use 'tidy/identities' endpoint to clean-up
the whitelist of identities.
`

package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func (b *Backend) programmaticAPIKeys() *framework.Secret {
	return &framework.Secret{
		Type: programmaticAPIKey,
		Fields: map[string]*framework.FieldSchema{
			"public_key": {
				Type:        framework.TypeString,
				Description: "Programmatic API Key Public Key",
			},

			"private_key": {
				Type:        framework.TypeString,
				Description: "Programmatic API Key Private Key",
			},
		},
		Renew:  b.programmaticAPIKeysRenew,
		Revoke: b.programmaticAPIKeyRevoke,
	}
}

func (b *Backend) programmaticAPIKeyCreate(ctx context.Context, s logical.Storage, displayName string, cred *atlasCredentialEntry) (*logical.Response, error) {

	apiKeyDescription, err := genUsername(displayName)
	if err != nil {
		return nil, errwrap.Wrapf("error generating username: {{err}}", err)
	}
	client, err := b.clientMongo(ctx, s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	walID, err := framework.PutWAL(ctx, s, programmaticAPIKey, &walEntry{
		UserName: apiKeyDescription,
	})
	if err != nil {
		return nil, errwrap.Wrapf("error writing WAL entry: {{err}}", err)
	}

	var key *mongodbatlas.APIKey

	switch {
	case isOrgKey(cred.OrganizationID, cred.ProjectID):
		key, err = createOrgKey(ctx, client, apiKeyDescription, cred)
	case isProjectKey(cred.OrganizationID, cred.ProjectID):
		key, err = createProjectAPIKey(ctx, client, apiKeyDescription, cred)
	case isAssignedToProject(cred.OrganizationID, cred.ProjectID):
		key, err = createAndAssignKey(ctx, client, apiKeyDescription, cred)
	}

	if err != nil {
		if walErr := framework.DeleteWAL(ctx, s, walID); walErr != nil {
			dbUserErr := errwrap.Wrapf("error creating programmaticAPIKey: {{err}}", err)
			return nil, errwrap.Wrap(errwrap.Wrapf("failed to delete WAL entry: {{err}}", walErr), dbUserErr)
		}
		return logical.ErrorResponse("Error creating programmatic api key: %s", err), err
	}

	if key == nil {
		return nil, errors.New("error creating credential")
	}

	if err := framework.DeleteWAL(ctx, s, walID); err != nil {
		return nil, errwrap.Wrapf("failed to commit WAL entry: {{err}}", err)
	}

	resp := b.Secret(programmaticAPIKey).Response(map[string]interface{}{
		"public_key":  key.PublicKey,
		"private_key": key.PrivateKey,
		"description": apiKeyDescription,
	}, map[string]interface{}{
		"programmatic_api_key_id": key.ID,
		"project_id":              cred.ProjectID,
		"organization_id":         cred.OrganizationID,
	})

	defaultLease, maxLease := b.getDefaultAndMaxLease()

	// If defined, credential TTL overrides default lease configuration
	if cred.TTL > 0 {
		defaultLease = cred.TTL
	}

	if cred.MaxTTL > 0 {
		maxLease = cred.MaxTTL
	}

	resp.Secret.TTL = defaultLease
	resp.Secret.MaxTTL = maxLease

	return resp, nil
}

func createOrgKey(ctx context.Context, client *mongodbatlas.Client, apiKeyDescription string, credentialEntry *atlasCredentialEntry) (*mongodbatlas.APIKey, error) {
	key, _, err := client.APIKeys.Create(ctx, credentialEntry.OrganizationID,
		&mongodbatlas.APIKeyInput{
			Desc:  apiKeyDescription,
			Roles: credentialEntry.Roles,
		})
	if err != nil {
		return nil, err
	}

	if err := addWhitelistEntry(ctx, client, credentialEntry.OrganizationID, key.ID, credentialEntry); err != nil {
		return nil, err
	}

	return key, nil
}

func createProjectAPIKey(ctx context.Context, client *mongodbatlas.Client, apiKeyDescription string, credentialEntry *atlasCredentialEntry) (*mongodbatlas.APIKey, error) {
	key, _, err := client.ProjectAPIKeys.Create(ctx, credentialEntry.ProjectID,
		&mongodbatlas.APIKeyInput{
			Desc:  apiKeyDescription,
			Roles: credentialEntry.Roles,
		})
	return key, err
}

func createAndAssignKey(ctx context.Context, client *mongodbatlas.Client, apiKeyDescription string, credentialEntry *atlasCredentialEntry) (*mongodbatlas.APIKey, error) {
	key, err := createOrgKey(ctx, client, apiKeyDescription, credentialEntry)
	if err != nil {
		return nil, err
	}

	if _, err := client.ProjectAPIKeys.Assign(ctx, credentialEntry.ProjectID, key.ID, &mongodbatlas.AssignAPIKey{
		Roles: credentialEntry.ProjectRoles,
	}); err != nil {
		return nil, err
	}

	return key, nil
}

func addWhitelistEntry(ctx context.Context, client *mongodbatlas.Client, orgID string, keyID string, cred *atlasCredentialEntry) error {
	var entries []*mongodbatlas.WhitelistAPIKeysReq
	for _, cidrBlock := range cred.CIDRBlocks {
		cidr := &mongodbatlas.WhitelistAPIKeysReq{
			CidrBlock: cidrBlock,
		}
		entries = append(entries, cidr)
	}

	for _, ipAddress := range cred.IPAddresses {
		ip := &mongodbatlas.WhitelistAPIKeysReq{
			IPAddress: ipAddress,
		}
		entries = append(entries, ip)

	}

	if entries != nil {
		_, _, err := client.WhitelistAPIKeys.Create(ctx, orgID, keyID, entries)
		return err

	}

	return nil
}

func (b *Backend) programmaticAPIKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	programmaticAPIKeyIDRaw, ok := req.Secret.InternalData["programmatic_api_key_id"]
	if !ok {
		return nil, fmt.Errorf("secret is missing programmatic api key id internal data")
	}

	programmaticAPIKeyID, ok := programmaticAPIKeyIDRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing programmatic api key id internal data")
	}

	organizationID := ""
	organizationIDRaw, ok := req.Secret.InternalData["organization_id"]
	if ok {
		organizationID, ok = organizationIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("secret is missing organization id internal data")
		}
	}

	projectID := ""
	projectIDRaw, ok := req.Secret.InternalData["project_id"]
	if ok {
		projectID, ok = projectIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("secret is missing project_id internal data")
		}
	}

	var data = map[string]interface{}{
		"organization_id":         organizationID,
		"programmatic_api_key_id": programmaticAPIKeyID,
		"project_id":              projectID,
	}

	// Use the user rollback mechanism to delete this database_user
	if err := b.pathProgrammaticAPIKeyRollback(ctx, req, programmaticAPIKey, data); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *Backend) pathProgrammaticAPIKeyRollback(ctx context.Context, req *logical.Request, _kind string, data interface{}) error {

	var entry walEntry
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	// Get the client
	client, err := b.clientMongo(ctx, req.Storage)
	if err != nil {
		return nil
	}

	switch {
	case isOrgKey(entry.OrganizationID, entry.ProjectID):
		// check if the user exists or not
		_, res, err := client.APIKeys.Get(ctx, entry.OrganizationID, entry.ProgrammaticAPIKeyID)
		// if the user is gone, move along
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}

		// now, delete the api key
		res, err = client.APIKeys.Delete(ctx, entry.OrganizationID, entry.ProgrammaticAPIKeyID)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}
	case isProjectKey(entry.OrganizationID, entry.ProjectID):
		// now, delete the user
		res, err := client.ProjectAPIKeys.Unassign(ctx, entry.ProjectID, entry.ProgrammaticAPIKeyID)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}
	case isAssignedToProject(entry.OrganizationID, entry.ProjectID):
		// check if the user exists or not
		_, res, err := client.APIKeys.Get(ctx, entry.OrganizationID, entry.ProgrammaticAPIKeyID)
		// if the user is gone, move along
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}

		// now, delete the api key
		res, err = client.APIKeys.Delete(ctx, entry.OrganizationID, entry.ProgrammaticAPIKeyID)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}

	}

	return nil
}

func (b *Backend) programmaticAPIKeysRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease (if any)

	defaultLease, maxLease := b.getDefaultAndMaxLease()

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = defaultLease
	resp.Secret.MaxTTL = maxLease
	return resp, nil
}

func (b *Backend) getDefaultAndMaxLease() (time.Duration, time.Duration) {
	maxLease := b.system.MaxLeaseTTL()
	defaultLease := b.system.DefaultLeaseTTL()

	if defaultLease > maxLease {
		maxLease = defaultLease
	}
	return defaultLease, maxLease

}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"go.mongodb.org/atlas/mongodbatlas"
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

func (b *Backend) programmaticAPIKeyCreate(ctx context.Context, s logical.Storage, role string, cred *atlasCredentialEntry) (*logical.Response, error) {

	apiKeyDescription, err := genAPIKeyDescription(role)
	if err != nil {
		return nil, errwrap.Wrapf("error generating username: {{err}}", err)
	}
	client, err := b.clientMongo(ctx, s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	walID, err := framework.PutWAL(ctx, s, programmaticAPIKey, &walEntry{
		Role: apiKeyDescription,
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
		return logical.ErrorResponse("Error creating programmatic api key: %s", err), nil
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
		"role":                    role,
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

	if err := addAccessListEntry(ctx, client, credentialEntry.OrganizationID, key.ID, credentialEntry); err != nil {
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
	if err != nil {
		return nil, err
	}

	orgIDs := map[string]interface{}{}

	// this is the only way to get the orgID needed for this request
	for _, r := range key.Roles {
		if _, ok := orgIDs[r.OrgID]; !ok {
			if len(r.OrgID) > 0 {
				orgIDs[r.OrgID] = 1
			}
		}
	}

	// if we have access list entries and no orgIds then return an error
	if (len(credentialEntry.IPAddresses)+len(credentialEntry.CIDRBlocks)) > 0 && len(orgIDs) == 0 {
		return nil, fmt.Errorf("No organization ID was found on programmatic key roles")
	}

	for orgID := range orgIDs {
		if err := addAccessListEntry(ctx, client, orgID, key.ID, credentialEntry); err != nil {
			return nil, err
		}
	}

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

func addAccessListEntry(ctx context.Context, client *mongodbatlas.Client, orgID string, keyID string, cred *atlasCredentialEntry) error {
	var entries []*mongodbatlas.AccessListAPIKeysReq
	for _, cidrBlock := range cred.CIDRBlocks {
		cidr := &mongodbatlas.AccessListAPIKeysReq{
			CidrBlock: cidrBlock,
		}
		entries = append(entries, cidr)
	}

	for _, ipAddress := range cred.IPAddresses {
		ip := &mongodbatlas.AccessListAPIKeysReq{
			IPAddress: ipAddress,
		}
		entries = append(entries, ip)

	}

	if entries != nil {
		_, _, err := client.AccessListAPIKeys.Create(ctx, orgID, keyID, entries)
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

func (b *Backend) pathProgrammaticAPIKeyRollback(ctx context.Context, req *logical.Request, _ string, data interface{}) error {
	var entry walEntry
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	// Get the client
	client, err := b.clientMongo(ctx, req.Storage)
	if err != nil {
		return nil
	}

	if isOrgKey(entry.OrganizationID, entry.ProjectID) || isAssignedToProject(entry.OrganizationID, entry.ProjectID) {
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
		return nil
	}

	if isProjectKey(entry.OrganizationID, entry.ProjectID) {

		// we need the orgID to delete the Key
		foundKey := mongodbatlas.APIKey{}

		// define the list options to get all the keys
		// currently pull 500 at a time, which is the max according to the docs
		// https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Programmatic-API-Keys/operation/listProjectApiKeys
		opts := &mongodbatlas.ListOptions{
			ItemsPerPage: 500,
		}

		keys, _, err := client.ProjectAPIKeys.List(ctx, entry.ProjectID, opts)
		if err != nil {
			return err
		}
		for _, key := range keys {
			if key.ID == entry.ProgrammaticAPIKeyID {
				foundKey = key
				break
			}
		}

		if foundKey.ID == "" {
			return fmt.Errorf("programmatic key %s not present in fetched page of API keys", entry.ProgrammaticAPIKeyID)
		}

		if len(foundKey.Roles) == 0 {
			return fmt.Errorf("missing roles on programmatic key %s", foundKey.ID)
		}

		// find the first orgID
		orgID := ""
		for _, r := range foundKey.Roles {
			if len(r.OrgID) > 0 {
				orgID = r.OrgID
				break
			}
		}

		// if orgID it's not found, return an error
		if len(orgID) == 0 {
			return fmt.Errorf("missing orgID on programmatic key %s", foundKey.ID)
		}

		// now, delete the user
		res, err := client.ProjectAPIKeys.Unassign(ctx, entry.ProjectID, entry.ProgrammaticAPIKeyID)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}
		// now, delete the api key
		res, err = client.APIKeys.Delete(ctx, orgID, entry.ProgrammaticAPIKeyID)
		if err != nil {
			if res != nil && res.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}
		return nil
	}

	return fmt.Errorf("Programmatic API key %s type not found, not deleting", entry.ProgrammaticAPIKeyID)
}

func (b *Backend) programmaticAPIKeysRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	//check if the role is on the secret
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, errors.New("internal data 'role' not found")
	}

	//get the credential entry
	role := roleRaw.(string)
	cred, err := b.credentialRead(ctx, req.Storage, role)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving credential: {{err}}", err)
	}

	if cred == nil {
		return nil, errors.New("error retrieving credential: credential is nil")
	}

	// Get the lease (if any)
	defaultLease, maxLease := b.getDefaultAndMaxLease()
	if cred.TTL > 0 {
		defaultLease = cred.MaxTTL
	}
	if cred.MaxTTL > 0 {
		maxLease = cred.MaxTTL
	}

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

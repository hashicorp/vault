// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *Backend) pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixMongoDBAtlas,
			OperationSuffix: "role",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the Roles",
				Required:    true,
			},
			"project_id": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf("Project ID the %s API key belongs to.", projectProgrammaticAPIKey),
			},
			"roles": {
				Type:        framework.TypeCommaStringSlice,
				Description: fmt.Sprintf("List of roles that the API Key should be granted. A minimum of one role must be provided. Any roles provided must be valid for the assigned Project, required for %s and %s keys.", orgProgrammaticAPIKey, projectProgrammaticAPIKey),
				Required:    true,
			},
			"organization_id": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf("Organization ID required for an %s API key", orgProgrammaticAPIKey),
			},
			"ip_addresses": {
				Type:        framework.TypeCommaStringSlice,
				Description: fmt.Sprintf("IP address to be added to the access list for the API key. Optional for %s and %s keys.", orgProgrammaticAPIKey, projectProgrammaticAPIKey),
			},
			"cidr_blocks": {
				Type:        framework.TypeCommaStringSlice,
				Description: fmt.Sprintf("Access list entry in CIDR notation to be added for the API key. Optional for %s and %s keys.", orgProgrammaticAPIKey, projectProgrammaticAPIKey),
			},
			"project_roles": {
				Type:        framework.TypeCommaStringSlice,
				Description: fmt.Sprintf("Roles assigned when an %s API Key is assigned to a %s API key", orgProgrammaticAPIKey, projectProgrammaticAPIKey),
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued credential should expire. Defaults to 0, in which case the value will fallback to the system/mount defaults.`,
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "The maximum allowed lifetime of credentials issued using this role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathRolesDelete,
			logical.ReadOperation:   b.pathRolesRead,
			logical.UpdateOperation: b.pathRolesWrite,
		},

		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *Backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "roles/"+d.Get("name").(string))
	return nil, err
}

func (b *Backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.credentialRead(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

func (b *Backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var resp logical.Response

	credentialName := d.Get("name").(string)
	if credentialName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	b.credentialMutex.Lock()
	defer b.credentialMutex.Unlock()
	credentialEntry, err := b.credentialRead(ctx, req.Storage, credentialName)
	if err != nil {
		return nil, err
	}

	if credentialEntry == nil {
		credentialEntry = &atlasCredentialEntry{}
	}

	if organizationIDRaw, ok := d.GetOk("organization_id"); ok {
		credentialEntry.OrganizationID = organizationIDRaw.(string)
	}

	getAPIAccessListArgs(credentialEntry, d)

	if projectIDRaw, ok := d.GetOk("project_id"); ok {
		projectID := projectIDRaw.(string)
		credentialEntry.ProjectID = projectID
	}

	if len(credentialEntry.OrganizationID) == 0 && len(credentialEntry.ProjectID) == 0 {
		return logical.ErrorResponse("organization_id or project_id are required"), nil
	}

	if programmaticKeyRolesRaw, ok := d.GetOk("roles"); ok {
		credentialEntry.Roles = programmaticKeyRolesRaw.([]string)
	} else {
		return logical.ErrorResponse("%s is required for %s and %s keys", "roles", orgProgrammaticAPIKey, projectProgrammaticAPIKey), nil
	}

	if projectRolesRaw, ok := d.GetOk("project_roles"); ok {
		credentialEntry.ProjectRoles = projectRolesRaw.([]string)
	} else {
		if isAssignedToProject(credentialEntry.OrganizationID, credentialEntry.ProjectID) {
			return logical.ErrorResponse("%s is required if both %s and %s are supplied", "roles", "organization_id", "project_id"), nil
		}
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		credentialEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}

	if maxttlRaw, ok := d.GetOk("max_ttl"); ok {
		credentialEntry.MaxTTL = time.Duration(maxttlRaw.(int)) * time.Second
	}

	if credentialEntry.MaxTTL > 0 && credentialEntry.TTL > credentialEntry.MaxTTL {
		return logical.ErrorResponse("ttl exceeds max_ttl"), nil
	}

	if err := setAtlasCredential(ctx, req.Storage, credentialName, credentialEntry); err != nil {
		return nil, err
	}

	return &resp, nil
}

func getAPIAccessListArgs(credentialEntry *atlasCredentialEntry, d *framework.FieldData) {

	if cidrBlocks, ok := d.GetOk("cidr_blocks"); ok {
		credentialEntry.CIDRBlocks = cidrBlocks.([]string)
	}
	if addresses, ok := d.GetOk("ip_addresses"); ok {
		credentialEntry.IPAddresses = addresses.([]string)
	}
}

func setAtlasCredential(ctx context.Context, s logical.Storage, credentialName string, credentialEntry *atlasCredentialEntry) error {
	if credentialName == "" {
		return fmt.Errorf("empty role name")
	}
	if credentialEntry == nil {
		return fmt.Errorf("emtpy credentialEntry")
	}
	entry, err := logical.StorageEntryJSON("roles/"+credentialName, credentialEntry)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("nil result when writing to storage")
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}
	return nil

}

func (b *Backend) credentialRead(ctx context.Context, s logical.Storage, credentialName string) (*atlasCredentialEntry, error) {
	if credentialName == "" {
		return nil, fmt.Errorf("missing credential name")
	}

	entry, err := s.Get(ctx, "roles/"+credentialName)
	if err != nil {
		return nil, err
	}
	var credentialEntry atlasCredentialEntry
	if entry != nil {
		if err := entry.DecodeJSON(&credentialEntry); err != nil {
			return nil, err
		}
		return &credentialEntry, nil
	}
	// Return nil here because all callers expect that if an entry
	// is nil, the method will return nil, nil.
	return nil, nil
}

type atlasCredentialEntry struct {
	ProjectID      string        `json:"project_id"`
	DatabaseName   string        `json:"database_name"`
	Roles          []string      `json:"roles"`
	OrganizationID string        `json:"organization_id"`
	CIDRBlocks     []string      `json:"cidr_blocks"`
	IPAddresses    []string      `json:"ip_addresses"`
	ProjectRoles   []string      `json:"project_roles"`
	TTL            time.Duration `json:"ttl"`
	MaxTTL         time.Duration `json:"max_ttl"`
}

func (r atlasCredentialEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"project_id":      r.ProjectID,
		"database_name":   r.DatabaseName,
		"roles":           r.Roles,
		"organization_id": r.OrganizationID,
		"cidr_blocks":     r.CIDRBlocks,
		"ip_addresses":    r.IPAddresses,
		"project_roles":   r.ProjectRoles,
		"ttl":             r.TTL.Seconds(),
		"max_ttl":         r.MaxTTL.Seconds(),
	}
	return respData
}

const pathRolesHelpSyn = `
Manage the roles used to generate MongoDB Atlas Programmatic API Keys.

`
const pathRolesHelpDesc = `
This path lets you manage the roles used to generate MongoDB Atlas Programmatic API Keys

The "project_id" parameter specifies a project where the Programmatic API Key will be
created.

"organization_id" parameter specifies in which Organization the key will be created.

If both are specified, the key will be created with the "organization_id" and then
assigned to the Project with the provided "project_id".

The "roles" parameter specifies the MongoDB Atlas Programmatic Key roles that should be assigned
to the Programmatic API keys created for a given role. At least one role should be provided
and must be valid for key level (project or org).

"ip_addresses" and "cidr_blocks" are used to add access list entries for the API key.

"project_roles" is used when both "organization_id" and "project_id" are supplied. 
And it's a list of roles that the API Key should be granted. A minimum of one role 
must be provided. Any roles provided must be valid for the assigned Project

To validate the keys, attempt to read an access key after writing the policy.
`
const orgProgrammaticAPIKey = `organization`
const projectProgrammaticAPIKey = `project`
const programmaticAPIKey = `programmatic_api_key`

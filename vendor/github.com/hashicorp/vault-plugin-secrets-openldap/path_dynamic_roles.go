// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldif"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	secretCredsType = "creds"
	dynamicRolePath = "role/"
	dynamicCredPath = "creds/"
)

func (b *backend) pathDynamicRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(dynamicRolePath, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationSuffix: "dynamic-role",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role (lowercase)",
					Required:    true,
				},
				"creation_ldif": {
					Type:        framework.TypeString,
					Description: "LDIF string used to create new entities within the LDAP system. This LDIF can be templated.",
					Required:    true,
				},
				"deletion_ldif": {
					Type:        framework.TypeString,
					Description: "LDIF string used to delete entities created within the LDAP system. This LDIF can be templated.",
					Required:    true,
				},
				"rollback_ldif": {
					Type:        framework.TypeString,
					Description: "LDIF string used to rollback changes in the event of a failure to create credentials. This LDIF can be templated.",
				},
				"username_template": {
					Type:        framework.TypeString,
					Description: "The template used to create a username",
				},
				"default_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default TTL for dynamic credentials",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Max TTL a dynamic credential can be extended to",
				},
			},
			ExistenceCheck: b.pathDynamicRoleExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathDynamicRoleCreateUpdate,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathDynamicRoleCreateUpdate,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathDynamicRoleRead,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathDynamicRoleDelete,
				},
			},
			HelpSynopsis:    staticRoleHelpSynopsis,
			HelpDescription: staticRoleHelpDescription,
		},
		{
			Pattern: strings.TrimSuffix(dynamicRolePath, "/") + optionalGenericNameWithForwardSlashListRegex("path"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "list",
				OperationSuffix: "dynamic-roles",
			},
			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeLowerCaseString,
					Description: "Path of roles to list",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathDynamicRoleList,
				},
			},
			HelpSynopsis:    "List all the dynamic roles Vault is currently managing in LDAP.",
			HelpDescription: "List all the dynamic roles being managed by Vault.",
		},
	}
}

func dynamicSecretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: secretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username of the generated account",
			},
			"password": {
				Type:        framework.TypeString,
				Description: "Password to access the generated account",
			},
			"distinguished_names": {
				Type: framework.TypeStringSlice,
				Description: "List of the distinguished names (DN) created. Each name in this list corresponds to" +
					"each action taken within the creation_ldif statements. This does not de-duplicate entries, " +
					"so this will have one entry for each LDIF statement within creation_ldif.",
			},
		},

		Renew:  b.secretCredsRenew(),
		Revoke: b.secretCredsRevoke(),
	}
}

func (b *backend) pathDynamicRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rawData := data.Raw
	err := convertToDuration(rawData, "default_ttl", "max_ttl")
	if err != nil {
		return nil, fmt.Errorf("failed to convert TTLs to duration: %w", err)
	}

	roleName := data.Get("name").(string)
	dRole, err := retrieveDynamicRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("unable to look for existing role: %w", err)
	}
	if dRole == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, fmt.Errorf("unable to update role: role does not exist")
		}
		dRole = &dynamicRole{}
	}
	err = mapstructure.WeakDecode(rawData, dRole)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	dRole.CreationLDIF = decodeBase64(dRole.CreationLDIF)
	dRole.RollbackLDIF = decodeBase64(dRole.RollbackLDIF)
	dRole.DeletionLDIF = decodeBase64(dRole.DeletionLDIF)

	err = validateDynamicRole(dRole)
	if err != nil {
		return nil, err
	}

	err = storeDynamicRole(ctx, req.Storage, dRole)
	if err != nil {
		return nil, fmt.Errorf("failed to save dynamic role: %w", err)
	}

	return nil, nil
}

func validateDynamicRole(dRole *dynamicRole) error {
	if dRole.CreationLDIF == "" {
		return fmt.Errorf("missing creation_ldif")
	}

	if dRole.DeletionLDIF == "" {
		return fmt.Errorf("missing deletion_ldif")
	}

	err := assertValidLDIFTemplate(dRole.CreationLDIF)
	if err != nil {
		return fmt.Errorf("invalid creation_ldif: %w", err)
	}

	err = assertValidLDIFTemplate(dRole.DeletionLDIF)
	if err != nil {
		return fmt.Errorf("invalid deletion_ldif: %w", err)
	}

	if dRole.RollbackLDIF != "" {
		err = assertValidLDIFTemplate(dRole.RollbackLDIF)
		if err != nil {
			return fmt.Errorf("invalid rollback_ldif: %w", err)
		}
	}

	return nil
}

// convertToDuration all keys in the data map into time.Duration objects. Keys not found in the map will be ignored
func convertToDuration(data map[string]interface{}, keys ...string) error {
	merr := new(multierror.Error)
	for _, key := range keys {
		val, exists := data[key]
		if !exists {
			continue
		}

		dur, err := parseutil.ParseDurationSecond(val)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("invalid duration %s: %w", key, err))
			continue
		}
		data[key] = dur
	}
	return merr.ErrorOrNil()
}

// decodeBase64 attempts to base64 decode the provided string. If the string is not base64 encoded, this
// returns the original string.
// This is equivalent to "if string is base64 encoded, decode it and return, otherwise return the original string"
func decodeBase64(str string) string {
	if str == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}
	return string(decoded)
}

func assertValidLDIFTemplate(rawTemplate string) error {
	// Test the template to ensure there aren't any errors in the template syntax
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	testTemplateData := dynamicTemplateData{
		Username:              "testuser",
		Password:              "testpass",
		DisplayName:           "testdisplayname",
		RoleName:              "testrolename",
		IssueTime:             now.Format(time.RFC3339),
		IssueTimeSeconds:      now.Unix(),
		ExpirationTime:        exp.Format(time.RFC3339),
		ExpirationTimeSeconds: exp.Unix(),
	}

	testLDIF, err := applyTemplate(rawTemplate, testTemplateData)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Test the LDIF to ensure there aren't any errors in the syntax
	entries, err := ldif.Parse(testLDIF)
	if err != nil {
		return fmt.Errorf("LDIF is invalid: %w", err)
	}

	if len(entries.Entries) == 0 {
		return fmt.Errorf("must specify at least one LDIF entry")
	}

	return nil
}

func (b *backend) pathDynamicRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)

	dRole, err := retrieveDynamicRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve dynamic role: %w", err)
	}
	if dRole == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"creation_ldif":     dRole.CreationLDIF,
			"deletion_ldif":     dRole.DeletionLDIF,
			"rollback_ldif":     dRole.RollbackLDIF,
			"username_template": dRole.UsernameTemplate,
			"default_ttl":       dRole.DefaultTTL.Seconds(),
			"max_ttl":           dRole.MaxTTL.Seconds(),
		},
	}
	return resp, nil
}

func (b *backend) pathDynamicRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rolePath := data.Get("path").(string)
	roles, err := req.Storage.List(ctx, dynamicRolePath+rolePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	return logical.ListResponse(roles), nil
}

func (b *backend) pathDynamicRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	roleName := data.Get("name").(string)
	role, err := retrieveDynamicRole(ctx, req.Storage, roleName)
	if err != nil {
		return false, fmt.Errorf("error finding role: %w", err)
	}
	return role != nil, nil
}

func (b *backend) pathDynamicRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)

	err := deleteDynamicRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to delete role: %w", err)
	}
	return nil, nil
}

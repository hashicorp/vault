// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kerberos

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type kerberosConfig struct {
	Keytab             string `json:"keytab"`
	ServiceAccount     string `json:"service_account"`
	AddGroupAliases    bool   `json:"add_group_aliases"`
	RemoveInstanceName bool   `json:"remove_instance_name"`
}

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKerberos,
		},
		Fields: map[string]*framework.FieldSchema{
			"keytab": {
				Type:        framework.TypeString,
				Description: `Base64 encoded keytab`,
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
			"service_account": {
				Type:        framework.TypeString,
				Description: `Service Account`,
			},
			"add_group_aliases": {
				Type: framework.TypeBool,
				Description: `If set to true, returns any groups found in LDAP as 
				a group alias.`,
			},
			"remove_instance_name": {
				Type:        framework.TypeBool,
				Description: `Remove instance/FQDN from keytab principal names.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
		},

		HelpSynopsis:    confHelpSynopsis,
		HelpDescription: confHelpDescription,
	}
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if config, err := b.config(ctx, req.Storage); err != nil {
		return nil, err
	} else if config == nil {
		return nil, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				// keytab is intentionally not returned here because it's sensitive
				"service_account":      config.ServiceAccount,
				"add_group_aliases":    config.AddGroupAliases,
				"remove_instance_name": config.RemoveInstanceName,
			},
		}, nil
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	serviceAccount := data.Get("service_account").(string)
	if serviceAccount == "" {
		return logical.ErrorResponse("data does not contain service_account"), logical.ErrInvalidRequest
	}

	kt := data.Get("keytab").(string)
	if kt == "" {
		return logical.ErrorResponse("data does not contain keytab"), logical.ErrInvalidRequest
	}

	addGroupAliases := data.Get("add_group_aliases").(bool)
	removeInstanceName := data.Get("remove_instance_name").(bool)

	// Check that the keytab is valid by parsing with krb5go
	if _, err := parseKeytab(kt); err != nil {
		return logical.ErrorResponse("invalid keytab: %v", err), logical.ErrInvalidRequest
	}

	config := &kerberosConfig{
		Keytab:             kt,
		ServiceAccount:     serviceAccount,
		AddGroupAliases:    addGroupAliases,
		RemoveInstanceName: removeInstanceName,
	}

	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

const confHelpSynopsis = `Configures the Kerberos keytab and service account.`
const confHelpDescription = `
The keytab must be base64 encoded, use the output of base64 <vault.keytab>.
`

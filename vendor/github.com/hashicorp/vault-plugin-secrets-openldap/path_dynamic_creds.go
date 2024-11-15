// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldif"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault-plugin-secrets-openldap/client"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/template"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/encoding/unicode"
)

func (b *backend) pathDynamicCredsCreate() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(dynamicCredPath, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "request",
				OperationSuffix: "dynamic-role-credentials",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the dynamic role.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathDynamicCredsRead,
				},
			},
			HelpSynopsis: "Request LDAP credentials for a dynamic role. These credentials are " +
				"created within the LDAP system when querying this endpoint.",
			HelpDescription: "This path requests new LDAP credentials for a certain dynamic role. " +
				"The credentials are created within the LDAP system based on the creation_ldif " +
				"specified within the dynamic role configuration. Depending on the LDAP implementation " +
				"the credentials may not be immediately usable due to eventual consistency.",
		},
	}
}

func (b *backend) pathDynamicCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)

	// Get the role and LDAP configs
	dRole, err := retrieveDynamicRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve dynamic role: %w", err)
	}
	if dRole == nil {
		return nil, nil
	}

	config, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("missing LDAP configuration")
	}

	// Generate dynamic data
	username, err := generateUsername(req, dRole)
	if err != nil {
		return nil, fmt.Errorf("failed to generate username: %w", err)
	}
	password, err := b.GeneratePassword(ctx, config)
	if err != nil {
		return nil, err
	}

	// Apply the template & execute
	now := time.Now()
	exp := now.Add(dRole.DefaultTTL)
	templateData := dynamicTemplateData{
		Username:              username,
		Password:              password,
		DisplayName:           req.DisplayName,
		RoleName:              roleName,
		IssueTime:             now.Format(time.RFC3339),
		IssueTimeSeconds:      now.Unix(),
		ExpirationTime:        exp.Format(time.RFC3339),
		ExpirationTimeSeconds: exp.Unix(),
	}
	dns, err := b.executeLDIF(config.LDAP, dRole.CreationLDIF, templateData, false)
	if err != nil {
		// Creation failed, attempt a rollback if one is specified
		if dRole.RollbackLDIF == "" {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		merr := multierror.Append(fmt.Errorf("failed to create user: %w", err))
		_, err = b.executeLDIF(config.LDAP, dRole.RollbackLDIF, templateData, true)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("failed to roll back user creation: %w", err))
		}
		return nil, merr
	}
	respData := map[string]interface{}{
		"username":            username,
		"password":            password,
		"distinguished_names": dns,
	}
	internal := map[string]interface{}{
		"name": roleName,
		// Including the deletion_ldif in the event that the role is deleted while
		// leases are active otherwise leases will fail to revoke
		"deletion_ldif": dRole.DeletionLDIF,
		"template_data": templateData,
	}
	resp := b.Secret(secretCredsType).Response(respData, internal)
	resp.Secret.TTL = dRole.DefaultTTL
	resp.Secret.MaxTTL = dRole.MaxTTL

	return resp, nil
}

// executeLDIF applies the template data against the LDIF template & executes the LDIF statements against the LDAP
// server. If more than one statement is specified within the LDIF string, this will result in multiple operations
// against the LDAP server. This is due to the fact that LDAP does not have transactions, nor any other form of
// atomic operations across multiple LDIF entries. If `continueOnError` is false, this will exit immediately upon
// any error occurring. If true, this will attempt to execute all of the specified LDIF statements and returns an error
// upon completion if any occurred.
func (b *backend) executeLDIF(config *client.Config, ldifTemplate string, templateData dynamicTemplateData, continueOnError bool) (dns []string, err error) {
	rawLDIF, err := applyTemplate(ldifTemplate, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	// Parse the raw LDIF & run it against the LDAP client
	entries, err := ldif.Parse(rawLDIF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated LDIF: %w", err)
	}

	err = b.client.Execute(config, entries.Entries, continueOnError)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statements: %w", err)
	}
	dns = getDNs(entries.Entries)
	return dns, nil
}

func getDNs(entries []*ldif.Entry) []string {
	dns := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry == nil {
			continue
		}

		switch {
		case entry.Entry != nil:
			dns = append(dns, entry.Entry.DN)
		case entry.Add != nil:
			dns = append(dns, entry.Add.DN)
		case entry.Modify != nil:
			dns = append(dns, entry.Modify.DN)
		case entry.Del != nil:
			dns = append(dns, entry.Del.DN)
		}
	}
	return dns
}

func (b *backend) secretCredsRenew() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Retrieve the role to ensure it still exists. If it doesn't, this will reject the renewal request.
		roleNameRaw, ok := req.Secret.InternalData["name"]
		if !ok {
			return nil, fmt.Errorf("missing role name")
		}

		roleName := roleNameRaw.(string)
		dRole, err := retrieveDynamicRole(ctx, req.Storage, roleName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve dynamic role: %w", err)
		}
		if dRole == nil {
			return nil, fmt.Errorf("unable to renew: role does not exist")
		}

		// Update the default TTL & MaxTTL to the latest from the role in the event the role definition has changed
		secret := req.Secret
		secret.TTL = dRole.DefaultTTL
		secret.MaxTTL = dRole.MaxTTL

		resp := &logical.Response{
			Secret: req.Secret,
		}
		return resp, nil
	}
}

func (b *backend) secretCredsRevoke() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		config, err := readConfig(ctx, req.Storage)
		if err != nil {
			return nil, err
		}
		if config == nil {
			return nil, fmt.Errorf("missing LDAP configuration")
		}

		deletionTemplate, err := getString(req.Secret.InternalData, "deletion_ldif")
		if err != nil {
			return nil, fmt.Errorf("broken internal data: unable to retrieve deletion_ldif: %w", err)
		}

		if deletionTemplate == "" {
			return nil, fmt.Errorf("broken internal data: missing deletion_ldif")
		}

		var templateData dynamicTemplateData
		rawTemplateData := req.Secret.InternalData["template_data"]
		switch td := rawTemplateData.(type) {
		case dynamicTemplateData:
			templateData = td
		case map[string]interface{}:
			err := mapstructure.WeakDecode(td, &templateData)
			if err != nil {
				return nil, fmt.Errorf("unable to decode internal data: %w", err)
			}
		default:
			return nil, fmt.Errorf("unable to revoke LDAP dynamic credentials: unrecognized internal data type: %T", td)
		}

		_, err = b.executeLDIF(config.LDAP, deletionTemplate, templateData, true)
		return nil, err
	}
}

type usernameTemplateData struct {
	DisplayName string
	RoleName    string
}

const defaultUsernameTemplate = "v_{{.DisplayName}}_{{.RoleName}}_{{random 10}}_{{unix_time}}"

func generateUsername(req *logical.Request, role *dynamicRole) (string, error) {
	usernameTemplate := role.UsernameTemplate
	if role.UsernameTemplate == "" {
		usernameTemplate = defaultUsernameTemplate
	}
	tmpl, err := template.NewTemplate(
		template.Template(usernameTemplate),
	)
	if err != nil {
		return "", err
	}
	usernameData := usernameTemplateData{
		DisplayName: req.DisplayName,
		RoleName:    role.Name,
	}
	return tmpl.Generate(usernameData)
}

type dynamicTemplateData struct {
	Username              string
	Password              string
	DisplayName           string
	RoleName              string
	IssueTime             string
	IssueTimeSeconds      int64
	ExpirationTime        string
	ExpirationTimeSeconds int64
}

func applyTemplate(rawTemplate string, data dynamicTemplateData) (string, error) {
	tmpl, err := template.NewTemplate(
		template.Template(rawTemplate),
		template.Function("utf16le", encodeUTF16LE),
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	str, err := tmpl.Generate(data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return str, nil
}

func encodeUTF16LE(str string) (string, error) {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	return enc.String(str)
}

func getString(m map[string]interface{}, key string) (string, error) {
	if m == nil {
		return "", fmt.Errorf("nil map")
	}

	val, exists := m[key]
	if !exists {
		return "", nil
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("key %s has %T value, not string", key, val)
	}
	return str, nil
}

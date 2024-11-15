// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	staticAccountStoragePrefix       = "static-account"
	staticAccountPathPrefix          = "static-account"
	gcpServiceAccountInferredProject = "-"
)

func pathStaticAccount(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", staticAccountPathPrefix, framework.GenericNameRegex("name")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationSuffix: "static-account",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Required. Name to refer to this static account in Vault. Cannot be updated.",
			},
			"secret_type": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf("Type of secret generated for this account. Cannot be updated. Defaults to %q", SecretTypeAccessToken),
				Default:     SecretTypeAccessToken,
			},
			"service_account_email": {
				Type:        framework.TypeString,
				Description: "Required. Email of the GCP service account to manage. Cannot be updated.",
			},
			"bindings": {
				Type:        framework.TypeString,
				Description: "Bindings configuration string.",
			},
			"token_scopes": {
				Type:        framework.TypeCommaStringSlice,
				Description: fmt.Sprintf(`List of OAuth scopes to assign to access tokens generated under this account. Ignored if "secret_type" is not "%q"`, SecretTypeAccessToken),
			},
		},
		ExistenceCheck: b.pathStaticAccountExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathStaticAccountDelete,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathStaticAccountRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback:                    b.pathStaticAccountCreate,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathStaticAccountUpdate,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},
		HelpSynopsis:    pathStaticAccountHelpSyn,
		HelpDescription: pathStaticAccountHelpDesc,
	}
}

func pathStaticAccountList(b *backend) *framework.Path {
	// Paths for listing static accounts
	return &framework.Path{
		Pattern: fmt.Sprintf("%ss?/?", staticAccountPathPrefix),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "list",
			OperationSuffix: "static-accounts|static-accounts2",
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathStaticAccountList,
			},
		},
		HelpSynopsis:    pathListStaticAccountHelpSyn,
		HelpDescription: pathListStaticAccountHelpDesc,
	}
}

func (b *backend) pathStaticAccountExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return false, errors.New("static account name is required")
	}

	acct, err := b.getStaticAccount(nameRaw.(string), ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return acct != nil, nil
}

func (b *backend) pathStaticAccountRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}

	acct, err := b.getStaticAccount(nameRaw.(string), ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"service_account_project": acct.Project,
		"service_account_email":   acct.EmailOrId,
		"secret_type":             acct.SecretType,
	}

	if len(acct.Bindings) > 0 {
		data["bindings"] = acct.Bindings.asOutput()
	}
	if acct.TokenGen != nil && acct.SecretType == SecretTypeAccessToken {
		data["token_scopes"] = acct.TokenGen.Scopes
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathStaticAccountDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}
	name := nameRaw.(string)

	b.staticAccountLock.Lock()
	defer b.staticAccountLock.Unlock()

	acct, err := b.getStaticAccount(name, ctx, req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to get static account %q: {{err}}", name), err)
	}
	if acct == nil {
		return nil, nil
	}

	resources := acct.boundResources()

	// Add WALs
	walIds, err := b.addWalsForStaticAccountResources(ctx, req, name, resources)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to create WALs for static account GCP resources %s: {{err}}", name), err)
	}

	// Delete static account
	b.Logger().Debug("deleting static account from storage", "name", name)
	if err := req.Storage.Delete(ctx, fmt.Sprintf("%s/%s", staticAccountStoragePrefix, name)); err != nil {
		return nil, err
	}

	// Try to clean up resources.
	if warnings := b.tryDeleteStaticAccountResources(ctx, req, resources, walIds); len(warnings) > 0 {
		b.Logger().Debug(
			"unable to delete GCP resources for deleted static account but WALs exist to clean up, ignoring errors",
			"static account", name, "errors", warnings)
		return &logical.Response{Warnings: warnings}, nil
	}

	b.Logger().Debug("finished deleting static account from storage", "name", name)
	return nil, nil
}

func (b *backend) pathStaticAccountCreate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	input, warnings, err := b.parseStaticAccountInformation(nil, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if input == nil {
		return nil, fmt.Errorf("plugin error - parse returned unexpected nil input")
	}

	b.staticAccountLock.Lock()
	defer b.staticAccountLock.Unlock()

	// Create and save static account with new resources.
	if err := b.createStaticAccount(ctx, req, input); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if len(warnings) > 0 {
		return &logical.Response{Warnings: warnings}, nil
	}
	return nil, nil
}

func (b *backend) pathStaticAccountUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}
	name := nameRaw.(string)

	b.staticAccountLock.Lock()
	defer b.staticAccountLock.Unlock()

	acct, err := b.getStaticAccount(name, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return nil, fmt.Errorf("unable to find static account %s to update", name)
	}

	initialInput := &inputParams{
		name:                acct.Name,
		secretType:          acct.SecretType,
		rawBindings:         acct.RawBindings,
		bindings:            acct.Bindings,
		project:             acct.Project,
		serviceAccountEmail: acct.EmailOrId,
	}
	if acct.TokenGen != nil {
		initialInput.scopes = acct.TokenGen.Scopes
	}

	updateInput, warnings, err := b.parseStaticAccountInformation(initialInput, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if updateInput == nil {
		return nil, fmt.Errorf("plugin error - parse returned unexpected nil input")
	}

	updateWarns, err := b.updateStaticAccount(ctx, req, acct, updateInput)
	if err != nil {
		return logical.ErrorResponse("unable to update: %s", err), nil
	}
	warnings = append(warnings, updateWarns...)
	if len(warnings) > 0 {
		return &logical.Response{Warnings: warnings}, nil
	}
	return nil, nil
}

func (b *backend) pathStaticAccountList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	accounts, err := req.Storage.List(ctx, fmt.Sprintf("%s/", staticAccountStoragePrefix))
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(accounts), nil
}

func (b *backend) parseStaticAccountInformation(prevValues *inputParams, d *framework.FieldData) (*inputParams, []string, error) {
	var warnings []string

	input := prevValues
	if prevValues == nil {
		input = &inputParams{}
	}

	nameRaw, ok := d.GetOk("name")
	if !ok {
		return nil, nil, fmt.Errorf("name is required")
	}
	input.name = nameRaw.(string)

	ws, err := input.parseOkInputSecretType(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	ws, err = input.parseOkInputServiceAccountEmail(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	ws, err = input.parseOkInputTokenScopes(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	ws, err = input.parseOkInputBindings(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	return input, warnings, nil
}

const pathStaticAccountHelpSyn = `Register and manage a GCP service account to generate credentials under`
const pathStaticAccountHelpDesc = `
This path allows you to register a static GCP service account that you want to generate secrets against.
This creates sets of IAM roles to specific GCP resources. Secrets (either service account keys or
access tokens) are generated under this account. The account must exist at creation of static account creation.

If bindings are specified, Vault will assign IAM permissions to the given service account. Bindings
can be given as a HCL (or JSON) string with the following format:

resource "some/gcp/resource/uri" {
	roles = [
		"roles/role1",
		"roles/role2",
		"roles/role3",
		...
	]
}

The given resource can have the following

* Project-level self link
	Self-link for a resource under a given project
	(i.e. resource name starts with 'projects/...')
	Use if you need to provide a versioned object or
	are directly using resource.self_link.

	Example (Compute instance):
		http://www.googleapis.com/compute/v1/projects/$PROJECT/zones/$ZONE/instances/$INSTANCE_NAME

* Full Resource Name
	A scheme-less URI consisting of a DNS-compatible
	API service name and a resource path (i.e. the
	relative resource name). Useful if you need to
	specify what service this resource is under
	but just want the preferred supported API version.
	Note that if the resource you are using is for
	a non-preferred API with multiple service versions,
	you MUST specify the version.

	Example (IAM service account):
		//$SERVICE.googleapis.com/projects/my-project/serviceAccounts/myserviceaccount@...

* Relative Resource Name:
	A URI path (path-noscheme) without the leading "/".
	It identifies a resource within the API service.
	Use if there is only one service that your
	resource could belong to. If there are multiple
	API versions that support the resource, we will
	attempt to use the preferred version and ask
	for more specific format otherwise.

	Example (Pubsub subscription):
		projects/myproject/subscriptions/mysub
`

const pathListStaticAccountHelpSyn = `List created static accounts.`
const pathListStaticAccountHelpDesc = `List created static accounts.`

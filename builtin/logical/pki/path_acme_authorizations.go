// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeAuthorization(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeAuthorization(b, baseUrl+"/authorization/"+framework.MatchAllRegex("auth_id"), opts)
}

func addFieldsForACMEAuthorization(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["auth_id"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME authorization identifier value",
		Required:    true,
	}

	return fields
}

func patternAcmeAuthorization(b *backend, pattern string, opts acmeWrapperOpts) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEAuthorization(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(opts, b.acmeAuthorizationHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeHelpSync,
		HelpDescription: pathAcmeHelpDesc,
	}
}

func (b *backend) acmeAuthorizationHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, _ *acmeAccount) (*logical.Response, error) {
	authId := fields.Get("auth_id").(string)
	authz, err := b.GetAcmeState().LoadAuthorization(acmeCtx, userCtx, authId)
	if err != nil {
		return nil, fmt.Errorf("failed to load authorization: %w", err)
	}

	var status string
	rawStatus, haveStatus := data["status"]
	if haveStatus {
		var ok bool
		status, ok = rawStatus.(string)
		if !ok {
			return nil, fmt.Errorf("bad type (%T) for value 'status': %w", rawStatus, ErrMalformed)
		}
	}

	if len(data) == 0 {
		return b.acmeAuthorizationFetchHandler(acmeCtx, r, fields, userCtx, data, authz)
	}

	if haveStatus && status == "deactivated" {
		return b.acmeAuthorizationDeactivateHandler(acmeCtx, r, fields, userCtx, data, authz)
	}

	return nil, ErrMalformed
}

func (b *backend) acmeAuthorizationFetchHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, authz *ACMEAuthorization) (*logical.Response, error) {
	return &logical.Response{
		Data: authz.NetworkMarshal(acmeCtx),
	}, nil
}

func (b *backend) acmeAuthorizationDeactivateHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, authz *ACMEAuthorization) (*logical.Response, error) {
	if authz.Status != ACMEAuthorizationPending && authz.Status != ACMEAuthorizationValid {
		return nil, fmt.Errorf("unable to deactivate authorization in '%v' status: %w", authz.Status, ErrMalformed)
	}

	authz.Status = ACMEAuthorizationDeactivated
	for _, challenge := range authz.Challenges {
		challenge.Status = ACMEChallengeInvalid
	}

	if err := b.GetAcmeState().SaveAuthorization(acmeCtx, authz); err != nil {
		return nil, fmt.Errorf("error saving deactivated authorization: %w", err)
	}

	return &logical.Response{
		Data: authz.NetworkMarshal(acmeCtx),
	}, nil
}

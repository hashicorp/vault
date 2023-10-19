// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathAcmeHelpSync = `An endpoint implementing the standard ACME protocol`
	pathAcmeHelpDesc = `This API endpoint implementing a subset of the ACME protocol
 defined in RFC 8555, with its own authentication and argument syntax that
 does not follow conventional Vault operations. An ACME client tool or library
 should be used to interact with these endpoints.`
)

func pathAcmeDirectory(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeDirectory(b, baseUrl+"/directory", opts)
}

func patternAcmeDirectory(b *backend, pattern string, opts acmeWrapperOpts) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(opts, b.acmeDirectoryHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeHelpSync,
		HelpDescription: pathAcmeHelpDesc,
	}
}

func (b *backend) acmeDirectoryHandler(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	rawBody, err := json.Marshal(map[string]interface{}{
		"newNonce":   acmeCtx.baseUrl.JoinPath("new-nonce").String(),
		"newAccount": acmeCtx.baseUrl.JoinPath("new-account").String(),
		"newOrder":   acmeCtx.baseUrl.JoinPath("new-order").String(),
		"revokeCert": acmeCtx.baseUrl.JoinPath("revoke-cert").String(),
		"keyChange":  acmeCtx.baseUrl.JoinPath("key-change").String(),
		// This is purposefully missing newAuthz as we don't support pre-authorization
		"meta": map[string]interface{}{
			"externalAccountRequired": acmeCtx.eabPolicy.IsExternalAccountRequired(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed encoding response: %w", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "application/json",
			logical.HTTPStatusCode:  http.StatusOK,
			logical.HTTPRawBody:     rawBody,
		},
	}, nil
}

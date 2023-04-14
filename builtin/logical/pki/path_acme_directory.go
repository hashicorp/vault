package pki

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathAcmeDirectoryHelpSync = `Read the proper URLs for various ACME operations`
	pathAcmeDirectoryHelpDesc = `Provide an ACME directory response that contains URLS for various ACME operations.`
)

func pathAcmeRootDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "acme/directory")
}

func pathAcmeRoleDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "roles/"+framework.GenericNameRegex("role")+"/acme/directory")
}

func pathAcmeIssuerDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/directory")
}

func pathAcmeIssuerAndRoleDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/directory")
}

func patternAcmeDirectory(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(b.acmeDirectoryHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeDirectoryHelpSync,
		HelpDescription: pathAcmeDirectoryHelpDesc,
	}
}

func (b *backend) acmeDirectoryHandler(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	rawBody, err := json.Marshal(map[string]interface{}{
		"newNonce":   acmeCtx.baseUrl.JoinPath("new-nonce").String(),
		"newAccount": acmeCtx.baseUrl.JoinPath("new-account").String(),
		"newOrder":   acmeCtx.baseUrl.JoinPath("new-order").String(),
		"revokeCert": acmeCtx.baseUrl.JoinPath("revoke-cert").String(),
		"keyChange":  acmeCtx.baseUrl.JoinPath("key-change").String(),
		"meta": map[string]interface{}{
			"externalAccountRequired": false,
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeNonce(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeNonce(b, baseUrl+"/new-nonce", opts)
}

func patternAcmeNonce(b *backend, pattern string, opts acmeWrapperOpts) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.HeaderOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(opts, b.acmeNonceHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(opts, b.acmeNonceHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeHelpSync,
		HelpDescription: pathAcmeHelpDesc,
	}
}

func (b *backend) acmeNonceHandler(ctx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	nonce, _, err := b.GetAcmeState().GetNonce()
	if err != nil {
		return nil, err
	}

	// Header operations return 200, GET return 204.
	httpStatus := http.StatusOK
	if r.Operation == logical.ReadOperation {
		httpStatus = http.StatusNoContent
	}

	return &logical.Response{
		Headers: map[string][]string{
			"Cache-Control": {"no-store"},
			"Replay-Nonce":  {nonce},
			"Link":          genAcmeLinkHeader(ctx),
		},
		Data: map[string]interface{}{
			logical.HTTPStatusCode: httpStatus,
			// Get around Vault limitation of requiring a body set if the status is not http.StatusNoContent
			// for our HEAD request responses.
			logical.HTTPContentType: "",
		},
	}, nil
}

func genAcmeLinkHeader(ctx *acmeContext) []string {
	path := fmt.Sprintf("<%s>;rel=\"index\"", ctx.baseUrl.JoinPath("directory").String())
	return []string{path}
}

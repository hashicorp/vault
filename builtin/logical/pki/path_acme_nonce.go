package pki

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeRootNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "acme/new-nonce")
}

func pathAcmeRoleNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "roles/"+framework.GenericNameRegex("role")+"/acme/new-nonce")
}

func pathAcmeIssuerNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/new-nonce")
}

func pathAcmeIssuerAndRoleNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/new-nonce")
}

func patternAcmeNonce(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.HeaderOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(b.acmeNonceHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(b.acmeNonceHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeDirectoryHelpSync,
		HelpDescription: pathAcmeDirectoryHelpDesc,
	}
}

func (b *backend) acmeNonceHandler(ctx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	nonce, _, err := b.acmeState.GetNonce()
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

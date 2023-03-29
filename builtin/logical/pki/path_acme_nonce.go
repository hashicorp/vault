package pki

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeRootNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "acme/new-nonce", false /* requireRole */, false /* requireIssuer */)
}

func pathAcmeRoleNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "roles/"+framework.GenericNameRegex("role")+"/acme/new-nonce",
		true /* requireRole */, false /* requireIssuer */)
}

func pathAcmeIssuerNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/new-nonce",
		false /* requireRole */, true /* requireIssuer */)
}

func pathAcmeIssuerAndRoleNonce(b *backend) *framework.Path {
	return patternAcmeNonce(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+"/roles/"+framework.GenericNameRegex(
			"role")+"/acme/new-nonce",
		true /* requireRole */, true /* requireIssuer */)
}

func patternAcmeNonce(b *backend, pattern string, requireRole, requireIssuer bool) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	if requireRole {
		fields["role"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The desired role for the acme request`,
			Required:    true,
		}
	}
	if requireIssuer {
		fields[issuerRefParam] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `Reference to an existing issuer name or issuer id`,
			Required:    true,
		}
	}
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

func (b *backend) acmeNonceHandler(ctx acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
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
		},
	}, nil
}

func genAcmeLinkHeader(ctx acmeContext) []string {
	path := fmt.Sprintf("<%s>;rel=\"index\"", ctx.baseUrl.JoinPath("/acme/directory").String())
	return []string{path}
}

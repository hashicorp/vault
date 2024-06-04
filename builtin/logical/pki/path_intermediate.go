// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateIntermediate(b *backend) *framework.Path {
	pattern := "intermediate/generate/" + framework.GenericNameRegex("exported")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "generate",
		OperationSuffix: "intermediate",
	}

	return buildPathGenerateIntermediate(b, pattern, displayAttrs)
}

func pathSetSignedIntermediate(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "intermediate/set-signed",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "set-signed",
			OperationSuffix: "intermediate",
		},

		Fields: map[string]*framework.FieldSchema{
			"certificate": {
				Type: framework.TypeString,
				Description: `PEM-format certificate. This must be a CA
certificate with a public key matching the
previously-generated key from the generation
endpoint. Additional parent CAs may be optionally
appended to the bundle.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImportIssuers,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"mapping": {
								Type:        framework.TypeMap,
								Description: "A mapping of issuer_id to key_id for all issuers included in this request",
								Required:    true,
							},
							"imported_keys": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Net-new keys imported as a part of this request",
								Required:    true,
							},
							"imported_issuers": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Net-new issuers imported as a part of this request",
								Required:    true,
							},
							"existing_keys": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Existing keys specified as part of the import bundle of this request",
								Required:    true,
							},
							"existing_issuers": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Existing issuers specified as part of the import bundle of this request",
								Required:    true,
							},
						},
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathSetSignedIntermediateHelpSyn,
		HelpDescription: pathSetSignedIntermediateHelpDesc,
	}

	return ret
}

func (b *backend) pathGenerateIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	var err error

	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not create intermediate until migration has completed"), nil
	}

	// Nasty hack :-) For cross-signing, we want to use the existing key, but
	// this isn't _actually_ part of the path. Put it into the request
	// parameters as if it was.
	if req.Path == "intermediate/cross-sign" {
		data.Raw["exported"] = "existing"
	}

	// Remove this once https://github.com/golang/go/issues/45990 is fixed
	data.Schema["use_pss"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
	}
	data.Raw["use_pss"] = false

	sc := b.makeStorageContext(ctx, req.Storage)
	exported, format, role, errorResp := getGenerationParams(sc, data)
	if errorResp != nil {
		return errorResp, nil
	}

	keyName, err := getKeyName(sc, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	var resp *logical.Response
	input := &inputBundle{
		role:    role,
		req:     req,
		apiData: data,
	}

	parsedBundle, warnings, err := generateIntermediateCSR(sc, input, b.Backend.GetRandomReader())
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		default:
			return nil, err
		}
	}

	csrb, err := parsedBundle.ToCSRBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw CSR bundle to CSR bundle: %w", err)
	}

	resp = &logical.Response{
		Data: map[string]interface{}{},
	}

	entries, err := getGlobalAIAURLs(ctx, req.Storage)
	if err == nil && len(entries.OCSPServers) == 0 && len(entries.IssuingCertificates) == 0 && len(entries.CRLDistributionPoints) == 0 {
		// If the operator hasn't configured any of the URLs prior to
		// generating this issuer, we should add a warning to the response,
		// informing them they might want to do so and re-generate the issuer.
		resp.AddWarning("This mount hasn't configured any authority information access (AIA) fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls or the newly generated issuer with this information. Since this certificate is an intermediate, it might be useful to regenerate this certificate after fixing this problem for the root mount.")
	}

	switch format {
	case "pem":
		resp.Data["csr"] = csrb.CSR
		if exported {
			resp.Data["private_key"] = csrb.PrivateKey
			resp.Data["private_key_type"] = csrb.PrivateKeyType
		}

	case "pem_bundle":
		resp.Data["csr"] = csrb.CSR
		if exported {
			resp.Data["csr"] = fmt.Sprintf("%s\n%s", csrb.PrivateKey, csrb.CSR)
			resp.Data["private_key"] = csrb.PrivateKey
			resp.Data["private_key_type"] = csrb.PrivateKeyType
		}

	case "der":
		resp.Data["csr"] = base64.StdEncoding.EncodeToString(parsedBundle.CSRBytes)
		if exported {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			resp.Data["private_key_type"] = csrb.PrivateKeyType
		}
	default:
		return nil, fmt.Errorf("unsupported format argument: %s", format)
	}

	if data.Get("private_key_format").(string) == "pkcs8" {
		err = convertRespToPKCS8(resp)
		if err != nil {
			return nil, err
		}
	}

	myKey, _, err := sc.importKey(csrb.PrivateKey, keyName, csrb.PrivateKeyType)
	if err != nil {
		return nil, err
	}
	resp.Data["key_id"] = myKey.ID

	resp = addWarnings(resp, warnings)

	return resp, nil
}

const pathGenerateIntermediateHelpSyn = `
Generate a new CSR and private key used for signing.
`

const pathGenerateIntermediateHelpDesc = `
See the API documentation for more information.
`

const pathSetSignedIntermediateHelpSyn = `
Provide the signed intermediate CA cert.
`

const pathSetSignedIntermediateHelpDesc = `
See the API documentation for more information.
`

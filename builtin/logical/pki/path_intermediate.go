package pki

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateIntermediate(b *backend) *framework.Path {
	return buildPathGenerateIntermediate(b, "intermediate/generate/"+framework.GenericNameRegex("exported"))
}

func pathSetSignedIntermediate(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "intermediate/set-signed",

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

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not create intermediate until migration has completed"), nil
	}

	// Nasty hack :-) For cross-signing, we want to use the existing key, but
	// this isn't _actually_ part of the path. Put it into the request
	// parameters as if it was.
	if req.Path == "intermediate/cross-sign" {
		data.Raw["exported"] = "existing"
	}

	// Nasty hack part two. :-) For generation of CSRs, certutil presently doesn't
	// support configuration of this. However, because we need generation parameters,
	// which create a role and attempt to read this parameter, we need to provide
	// a value (which will be ignored). Hence, we stub in the missing parameter here,
	// including its schema, just enough for it to work..
	data.Schema["signature_bits"] = &framework.FieldSchema{
		Type:    framework.TypeInt,
		Default: 0,
	}
	data.Raw["signature_bits"] = 0

	exported, format, role, errorResp := b.getGenerationParams(ctx, req.Storage, data)
	if errorResp != nil {
		return errorResp, nil
	}

	keyName, err := getKeyName(ctx, req.Storage, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	var resp *logical.Response
	input := &inputBundle{
		role:    role,
		req:     req,
		apiData: data,
	}
	parsedBundle, err := generateIntermediateCSR(ctx, b, input, b.Backend.GetRandomReader())
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

	myKey, _, err := importKey(ctx, b, req.Storage, csrb.PrivateKey, keyName, csrb.PrivateKeyType)
	if err != nil {
		return nil, err
	}
	resp.Data["key_id"] = myKey.ID

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

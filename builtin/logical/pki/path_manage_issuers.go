package pki

import (
	"bytes"
	"context"
	"encoding/pem"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssuerGenerateRoot(b *backend) *framework.Path {
	return buildPathGenerateRoot(b, "issuers/generate/root/"+framework.GenericNameRegex("exported"))
}

func buildPathGenerateRoot(b *backend, pattern string) *framework.Path {
	ret := &framework.Path{
		Pattern: pattern,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCAGenerateRoot,
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateRootHelpSyn,
		HelpDescription: pathGenerateRootHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAKeyGenerationFields(ret.Fields)
	ret.Fields = addCAIssueFields(ret.Fields)
	return ret
}

func pathIssuerGenerateIntermediate(b *backend) *framework.Path {
	return buildPathGenerateIntermediate(b,
		"issuers/generate/intermediate/"+framework.GenericNameRegex("exported"))
}

func buildPathGenerateIntermediate(b *backend, pattern string) *framework.Path {
	ret := &framework.Path{
		Pattern: pattern,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathGenerateIntermediate,
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateIntermediateHelpSyn,
		HelpDescription: pathGenerateIntermediateHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAKeyGenerationFields(ret.Fields)
	ret.Fields["add_basic_constraints"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Whether to add a Basic Constraints
extension with CA: true. Only needed as a
workaround in some compatibility scenarios
with Active Directory Certificate Services.`,
	}

	return ret
}

func pathImportIssuer(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issuers/import/(cert|bundle)",
		Fields: map[string]*framework.FieldSchema{
			"pem_bundle": {
				Type: framework.TypeString,
				Description: `PEM-format, concatenated unencrypted
secret-key (optional) and certificates.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathImportIssuers,
		},

		HelpSynopsis:    pathImportIssuersHelpSyn,
		HelpDescription: pathImportIssuersHelpDesc,
	}
}

func (b *backend) pathImportIssuers(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keysAllowed := strings.HasSuffix(req.Path, "bundle") || req.Path == "config/ca"

	var pemBundle string
	var certificate string
	rawPemBundle, bundleOk := data.GetOk("pem_bundle")
	rawCertificate, certOk := data.GetOk("certificate")
	if bundleOk {
		pemBundle = rawPemBundle.(string)
	}
	if certOk {
		certificate = rawCertificate.(string)
	}

	if len(pemBundle) == 0 && len(certificate) == 0 {
		return logical.ErrorResponse("'pem_bundle' and 'certificate' parameters were empty"), nil
	}
	if len(pemBundle) > 0 && len(certificate) > 0 {
		return logical.ErrorResponse("'pem_bundle' and 'certificate' parameters were both provided"), nil
	}
	if len(certificate) > 0 {
		keysAllowed = false
		pemBundle = certificate
	}

	var createdKeys []string
	var createdIssuers []string
	issuerKeyMap := make(map[string]string)

	// Rather than using certutil.ParsePEMBundle (which restricts the
	// construction of the PEM bundle), we manually parse the bundle instead.
	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block

	var issuers []string
	var keys []string

	for len(bytes.TrimSpace(pemBytes)) > 0 {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "no data found in PEM block"}
		}

		pemBlockString := string(pem.EncodeToMemory(pemBlock))

		switch pemBlock.Type {
		case "CERTIFICATE", "X509 CERTIFICATE":
			// Must be a certificate
			issuers = append(issuers, pemBlockString)
		case "CRL", "X509 CRL":
			// Ignore any CRL entries.
		default:
			// Otherwise, treat them as keys.
			keys = append(keys, pemBlockString)
		}
	}

	if len(keys) > 0 && !keysAllowed {
		return logical.ErrorResponse("private keys found in the PEM bundle but not allowed by the path; use /issuers/import/bundle"), nil
	}

	for _, keyPem := range keys {
		// Handle import of private key.
		key, existing, err := importKey(ctx, req.Storage, keyPem, "")
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if !existing {
			createdKeys = append(createdKeys, key.ID.String())
		}
	}

	for _, certPem := range issuers {
		cert, existing, err := importIssuer(ctx, req.Storage, certPem, "")
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		issuerKeyMap[cert.ID.String()] = cert.KeyID.String()
		if !existing {
			createdIssuers = append(createdIssuers, cert.ID.String())
		}
	}

	if len(createdIssuers) > 0 {
		err := buildCRL(ctx, b, req, true)
		if err != nil {
			return nil, err
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"mapping":          issuerKeyMap,
			"imported_keys":    createdKeys,
			"imported_issuers": createdIssuers,
		},
	}, nil
}

const (
	pathImportIssuersHelpSyn  = `Import the specified issuing certificates.`
	pathImportIssuersHelpDesc = `
This endpoint allows importing the specified issuer certificates.

:type is either the literal value "cert", to only allow importing
certificates, else "bundle" to allow importing keys as well as
certificates.

Depending on the value of :type, the pem_bundle request parameter can
either take PEM-formatted certificates, and, if :type="bundle", unencrypted
secret-keys.
`
)

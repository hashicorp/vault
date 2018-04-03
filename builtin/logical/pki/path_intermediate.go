package pki

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathGenerateIntermediate(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "intermediate/generate/" + framework.GenericNameRegex("exported"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathGenerateIntermediate,
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

func pathSetSignedIntermediate(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "intermediate/set-signed",

		Fields: map[string]*framework.FieldSchema{
			"certificate": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `PEM-format certificate. This must be a CA
certificate with a public key matching the
previously-generated key from the generation
endpoint.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSetSignedIntermediate,
		},

		HelpSynopsis:    pathSetSignedIntermediateHelpSyn,
		HelpDescription: pathSetSignedIntermediateHelpDesc,
	}

	return ret
}

func (b *backend) pathGenerateIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	exported, format, role, errorResp := b.getGenerationParams(data)
	if errorResp != nil {
		return errorResp, nil
	}

	var resp *logical.Response
	input := &dataBundle{
		role:    role,
		req:     req,
		apiData: data,
	}
	parsedBundle, err := generateIntermediateCSR(b, input)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}

	csrb, err := parsedBundle.ToCSRBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw CSR bundle to CSR bundle: {{err}}", err)
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
	}

	if data.Get("private_key_format").(string) == "pkcs8" {
		err = convertRespToPKCS8(resp)
		if err != nil {
			return nil, err
		}
	}

	cb := &certutil.CertBundle{}
	cb.PrivateKey = csrb.PrivateKey
	cb.PrivateKeyType = csrb.PrivateKeyType

	entry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *backend) pathSetSignedIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cert := data.Get("certificate").(string)

	if cert == "" {
		return logical.ErrorResponse("no certificate provided in the \"certificate\" parameter"), nil
	}

	inputBundle, err := certutil.ParsePEMBundle(cert)
	if err != nil {
		switch err.(type) {
		case errutil.InternalError:
			return nil, err
		default:
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if inputBundle.Certificate == nil {
		return logical.ErrorResponse("supplied certificate could not be successfully parsed"), nil
	}

	cb := &certutil.CertBundle{}
	entry, err := req.Storage.Get(ctx, "config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("could not find any existing entry with a private key"), nil
	}

	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, err
	}

	if len(cb.PrivateKey) == 0 || cb.PrivateKeyType == "" {
		return logical.ErrorResponse("could not find an existing private key"), nil
	}

	parsedCB, err := cb.ToParsedCertBundle()
	if err != nil {
		return nil, err
	}
	if parsedCB.PrivateKey == nil {
		return nil, fmt.Errorf("saved key could not be parsed successfully")
	}

	inputBundle.PrivateKey = parsedCB.PrivateKey
	inputBundle.PrivateKeyType = parsedCB.PrivateKeyType
	inputBundle.PrivateKeyBytes = parsedCB.PrivateKeyBytes

	if !inputBundle.Certificate.IsCA {
		return logical.ErrorResponse("the given certificate is not marked for CA use and cannot be used with this backend"), nil
	}

	if err := inputBundle.Verify(); err != nil {
		return nil, errwrap.Wrapf("verification of parsed bundle failed: {{err}}", err)
	}

	cb, err = inputBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw values into cert bundle: {{err}}", err)
	}

	entry, err = logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry.Key = "certs/" + normalizeSerial(cb.SerialNumber)
	entry.Value = inputBundle.CertificateBytes
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	// For ease of later use, also store just the certificate at a known
	// location
	entry.Key = "ca"
	entry.Value = inputBundle.CertificateBytes
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Build a fresh CRL
	err = buildCRL(ctx, b, req)

	return nil, err
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

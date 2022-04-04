package pki

import (
	"context"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Returns the CA in raw format
func pathFetchCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `ca(/pem)?`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchCARawHandler,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCARawHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var serial, pemType, contentType string
	serial = "ca"
	contentType = "application/pkix-cert"
	if req.Path == "ca/pem" {
		pemType = "CERTIFICATE"
		contentType = "application/pem-certificate-chain"
	}

	return b.pathFetchReadRaw(ctx, req, data, serial, pemType, contentType)
}

// Returns the CA chain
func pathFetchCAChain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/ca_chain`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func pathFetchCAChainRaw(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `ca_chain`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchCAChainRawHandler,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCAChainRawHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var serial, pemType, contentType string
	serial = "ca_chain"
	contentType = "application/pkix-cert"

	return b.pathFetchReadRaw(ctx, req, data, serial, pemType, contentType)
}

// Returns the CRL in raw format
func pathFetchCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `crl(/pem)?`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchCRLRawHandler,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCRLRawHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var serial, pemType, contentType string
	serial = "crl"
	contentType = "application/pkix-crl"
	if req.Path == "crl/pem" {
		pemType = "X509 CRL"
		contentType = "application/x-pem-file"
	}

	return b.pathFetchReadRaw(ctx, req, data, serial, pemType, contentType)
}

// Returns any valid (non-revoked) cert in raw format.
func pathFetchValidRaw(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/(?P<serial>[0-9A-Fa-f-:]+)/raw(/pem)?`,
		Fields: map[string]*framework.FieldSchema{
			"serial": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchCertificateRawHandler,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCertificateRawHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var serial, pemType, contentType string

	serial = data.Get("serial").(string)
	contentType = "application/pkix-cert"
	if strings.HasSuffix(req.Path, "/pem") {
		pemType = "CERTIFICATE"
		contentType = "application/pem-certificate-chain"
	}

	return b.pathFetchReadRaw(ctx, req, data, serial, pemType, contentType)
}

// Returns any valid (non-revoked) cert. Since "ca" fits the pattern, this path
// also handles returning the CA cert in a non-raw format.
func pathFetchValid(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/(?P<serial>[0-9A-Fa-f-:]+)`,
		Fields: map[string]*framework.FieldSchema{
			"serial": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// This returns the CRL in a non-raw format
func pathFetchCRLViaCertPath(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/crl`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// This returns the list of serial numbers for certs
func pathFetchListCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathFetchCertList,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCertList(ctx context.Context, req *logical.Request, data *framework.FieldData) (response *logical.Response, retErr error) {
	entries, err := req.Storage.List(ctx, "certs/")
	if err != nil {
		return nil, err
	}
	for i := range entries {
		entries[i] = denormalizeSerial(entries[i])
	}
	return logical.ListResponse(entries), nil
}

func marshalPem(pemType string, certificate []byte) []byte {
	block := pem.Block{
		Type:  pemType,
		Bytes: certificate,
	}

	// This is convoluted on purpose to ensure that we don't have trailing
	// newlines via various paths
	return []byte(strings.TrimSpace(string(pem.EncodeToMemory(&block))))
}

func (b *backend) pathFetchRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (response *logical.Response, retErr error) {
	var serial, pemType string
	var certEntry, revokedEntry *logical.StorageEntry
	var funcErr error
	var certificate []byte
	var fullChain []byte
	var revocationTime int64
	response = &logical.Response{
		Data: map[string]interface{}{},
	}

	// Errors don't cause an immediate exit, because the raw
	// paths still need to return raw output.

	switch {
	case req.Path == "cert/ca_chain":
		serial = "ca_chain"
	case req.Path == "cert/crl":
		serial = "crl"
		pemType = "X509 CRL"
	default:
		serial = data.Get("serial").(string)
		pemType = "CERTIFICATE"
	}
	if len(serial) == 0 {
		response = logical.ErrorResponse("The serial number must be provided")
		goto reply
	}

	if serial == "ca_chain" {
		caInfo, err := fetchCAInfo(ctx, b, req)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				response = logical.ErrorResponse(err.Error())
				goto reply
			default:
				retErr = err
				goto reply
			}
		}

		caChain := caInfo.GetCAChain()
		var certStr string
		for _, ca := range caChain {
			block := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: ca.Bytes,
			}
			certStr = strings.Join([]string{certStr, strings.TrimSpace(string(pem.EncodeToMemory(&block)))}, "\n")
		}
		certificate = []byte(strings.TrimSpace(certStr))

		rawChain := caInfo.GetFullChain()
		var chainStr string
		for _, ca := range rawChain {
			block := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: ca.Bytes,
			}
			chainStr = strings.Join([]string{chainStr, strings.TrimSpace(string(pem.EncodeToMemory(&block)))}, "\n")
		}
		fullChain = []byte(strings.TrimSpace(chainStr))

		goto reply
	}

	certEntry, funcErr = fetchCertBySerial(ctx, req, req.Path, serial)
	if funcErr != nil {
		switch funcErr.(type) {
		case errutil.UserError:
			response = logical.ErrorResponse(funcErr.Error())
			goto reply
		default:
			retErr = funcErr
			goto reply
		}
	}
	if certEntry == nil {
		response = nil
		goto reply
	}

	certificate = certEntry.Value
	if len(pemType) != 0 {
		certificate = marshalPem(pemType, certificate)
	}

	revokedEntry, funcErr = fetchCertBySerial(ctx, req, "revoked/", serial)
	if funcErr != nil {
		switch funcErr.(type) {
		case errutil.UserError:
			response = logical.ErrorResponse(funcErr.Error())
			goto reply
		default:
			retErr = funcErr
			goto reply
		}
	}
	if revokedEntry != nil {
		var revInfo revocationInfo
		err := revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error decoding revocation entry for serial %s: %s", serial, err)), nil
		}
		revocationTime = revInfo.RevocationTime
	}

reply:
	switch {
	case retErr != nil:
		response = nil
		return
	case response == nil:
		return
	case response.IsError():
		return response, nil
	default:
		response.Data["certificate"] = string(certificate)
		response.Data["revocation_time"] = revocationTime

		if len(fullChain) > 0 {
			response.Data["ca_chain"] = string(fullChain)
		}
	}

	return
}

func (b *backend) readCertByAlias(ctx context.Context, req *logical.Request, alias string) ([]byte, error) {
	// ca_chain as an alias needs special handling.
	if alias == "ca_chain" {
		caInfo, err := fetchCAInfo(ctx, b, req)
		if err != nil {
			return nil, err
		}

		caChain := caInfo.GetCAChain()
		var certStr string
		for _, ca := range caChain {
			block := pem.Block{
				Type:  "CERTIFICATE",
				Bytes: ca.Bytes,
			}
			certStr = strings.Join([]string{certStr, strings.TrimSpace(string(pem.EncodeToMemory(&block)))}, "\n")
		}

		return []byte(strings.TrimSpace(certStr)), nil
	}

	certEntry, err := fetchCertBySerial(ctx, req, req.Path, alias)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		return nil, nil
	}

	return certEntry.Value, nil
}

func (b *backend) pathFetchReadRaw(ctx context.Context, req *logical.Request, data *framework.FieldData, serial string, pemType string, contentType string) (response *logical.Response, retErr error) {
	var certificate []byte

	// Errors don't cause an immediate exit, because the raw
	// paths still need to return raw output, even if it is empty.
	// The error will instead be logged (in the event the log surfaces
	// warnings).
	certificate, retErr = b.readCertByAlias(ctx, req, serial)
	if retErr != nil {
		goto reply
	}

	if len(pemType) != 0 {
		certificate = marshalPem(pemType, certificate)
	}

reply:
	response = &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: contentType,
			logical.HTTPRawBody:     certificate,
		},
	}

	if len(certificate) > 0 {
		response.Data[logical.HTTPStatusCode] = 200
	} else {
		response.Data[logical.HTTPStatusCode] = 204
	}

	if retErr != nil {
		if b.Logger().IsWarn() {
			b.Logger().Warn("possible error, but cannot return in raw response. Note that an empty CA probably means none was configured, and an empty CRL is possibly correct; call /crl/rotate to create a non-empty CRL", "error", retErr)

			// Return a 500 response to indicate something went wrong with
			// this request.
			response.Data[logical.HTTPStatusCode] = 500
		}
	}

	retErr = nil

	return
}

const pathFetchHelpSyn = `
Fetch a CA, CRL, CA Chain, or non-revoked certificate.
`

const pathFetchHelpDesc = `
This allows certificates to be fetched. Use /cert/:serial for JSON responses.

Using "ca" or "crl" as the value fetches the appropriate information in DER encoding. Add "/pem" to either to get PEM encoding.

Using "ca_chain" as the value fetches the certificate authority trust chain in PEM encoding.

Otherwise, specify a serial number to fetch the specified certificate. Add "/raw" to get just the certificate in DER form, "/raw/pem" to get the PEM encoded certificate.
`

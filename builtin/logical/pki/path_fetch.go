package pki

import (
	"context"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Returns the CA in raw format
func pathFetchCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `ca(/pem)?`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns the CA chain
func pathFetchCAChain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `(cert/)?ca_chain`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns the CRL in raw format
func pathFetchCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `crl(/pem)?`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns any valid (non-revoked) cert. Since "ca" fits the pattern, this path
// also handles returning the CA cert in a non-raw format.
func pathFetchValid(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/(?P<serial>[0-9A-Fa-f-:]+)`,
		Fields: map[string]*framework.FieldSchema{
			"serial": &framework.FieldSchema{
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

	return logical.ListResponse(entries), nil
}

func (b *backend) pathFetchRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (response *logical.Response, retErr error) {
	var serial, pemType, contentType string
	var certEntry, revokedEntry *logical.StorageEntry
	var funcErr error
	var certificate []byte
	var revocationTime int64
	response = &logical.Response{
		Data: map[string]interface{}{},
	}

	// Some of these need to return raw and some non-raw;
	// this is basically handled by setting contentType or not.
	// Errors don't cause an immediate exit, because the raw
	// paths still need to return raw output.

	switch {
	case req.Path == "ca" || req.Path == "ca/pem":
		serial = "ca"
		contentType = "application/pkix-cert"
		if req.Path == "ca/pem" {
			pemType = "CERTIFICATE"
		}
	case req.Path == "ca_chain" || req.Path == "cert/ca_chain":
		serial = "ca_chain"
		if req.Path == "ca_chain" {
			contentType = "application/pkix-cert"
		}
	case req.Path == "crl" || req.Path == "crl/pem":
		serial = "crl"
		contentType = "application/pkix-crl"
		if req.Path == "crl/pem" {
			pemType = "X509 CRL"
		}
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
		caInfo, err := fetchCAInfo(ctx, req)
		switch err.(type) {
		case errutil.UserError:
			response = logical.ErrorResponse(err.Error())
			goto reply
		case errutil.InternalError:
			retErr = err
			goto reply
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
		goto reply
	}

	certEntry, funcErr = fetchCertBySerial(ctx, req, req.Path, serial)
	if funcErr != nil {
		switch funcErr.(type) {
		case errutil.UserError:
			response = logical.ErrorResponse(funcErr.Error())
			goto reply
		case errutil.InternalError:
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
		block := pem.Block{
			Type:  pemType,
			Bytes: certEntry.Value,
		}
		// This is convoluted on purpose to ensure that we don't have trailing
		// newlines via various paths
		certificate = []byte(strings.TrimSpace(string(pem.EncodeToMemory(&block))))
	}

	revokedEntry, funcErr = fetchCertBySerial(ctx, req, "revoked/", serial)
	if funcErr != nil {
		switch funcErr.(type) {
		case errutil.UserError:
			response = logical.ErrorResponse(funcErr.Error())
			goto reply
		case errutil.InternalError:
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
	case len(contentType) != 0:
		response = &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: contentType,
				logical.HTTPRawBody:     certificate,
			}}
		if retErr != nil {
			if b.Logger().IsWarn() {
				b.Logger().Warn("possible error, but cannot return in raw response. Note that an empty CA probably means none was configured, and an empty CRL is possibly correct", "error", retErr)
			}
		}
		retErr = nil
		if len(certificate) > 0 {
			response.Data[logical.HTTPStatusCode] = 200
		} else {
			response.Data[logical.HTTPStatusCode] = 204
		}
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
	}

	return
}

const pathFetchHelpSyn = `
Fetch a CA, CRL, CA Chain, or non-revoked certificate.
`

const pathFetchHelpDesc = `
This allows certificates to be fetched. If using the fetch/ prefix any non-revoked certificate can be fetched.

Using "ca" or "crl" as the value fetches the appropriate information in DER encoding. Add "/pem" to either to get PEM encoding.

Using "ca_chain" as the value fetches the certificate authority trust chain in PEM encoding.
`

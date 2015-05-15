package pki

import (
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

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

func pathFetchValid(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/(?P<serial>[0-9A-Fa-f-:]+)`,
		Fields: map[string]*framework.FieldSchema{
			"serial": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Certificate serial number, in colon- or hyphen-separated octal",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

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

func pathFetchRevoked(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `revoked/(?P<serial>[0-9A-Fa-f-:]+)`,
		Fields: map[string]*framework.FieldSchema{
			"serial": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Certificate serial number, in colon- or hyphen-separated octal",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchRead,
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchRead(req *logical.Request, data *framework.FieldData) (response *logical.Response, retErr error) {
	var serial string
	var pemType string
	var contentType string
	var certEntry *logical.StorageEntry
	var userErr, intErr, err error
	var certificate []byte
	response = &logical.Response{
		Data: map[string]interface{}{},
	}

	switch {
	case req.Path == "ca" || req.Path == "ca/pem":
		serial = "ca"
		contentType = "application/pkix-cert"
		if req.Path == "ca/pem" {
			pemType = "CERTIFICATE"
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

	_, _, err = fetchCAInfo(req)
	if err != nil {
		response = logical.ErrorResponse("No CA information configured for this backend")
		goto reply
	}

	certEntry, userErr, intErr = fetchCertBySerial(req, req.Path, serial)
	switch {
	case userErr != nil:
		response = logical.ErrorResponse(userErr.Error())
		goto reply
	case intErr != nil:
		retErr = intErr
		goto reply
	}

	switch {
	case strings.HasPrefix(req.Path, "revoked/"):
		var revInfo revocationInfo
		err := certEntry.DecodeJSON(&revInfo)
		if err != nil {
			retErr = fmt.Errorf("Error decoding revocation entry for serial %s: %s", serial, err)
			goto reply
		}
		certificate = revInfo.CertificateBytes
		response.Data["revocation_time"] = revInfo.RevocationTime
	default:
		certificate = certEntry.Value
	}

	if len(pemType) != 0 {
		block := pem.Block{
			Type:  pemType,
			Bytes: certEntry.Value,
		}
		certificate = pem.EncodeToMemory(&block)
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
			b.Logger().Printf("Possible error, but cannot return in raw response: %s. Note that an empty CA probably means none was configured, and an empty CRL is quite possibly correct", retErr)
		}
		retErr = nil
		response.Data[logical.HTTPStatusCode] = 200

	case retErr != nil:
		response = nil
	default:
		response.Data["certificate"] = string(certificate)
	}

	return
}

const pathFetchHelpSyn = `
Fetch a CA, CRL, valid or revoked certificate.
`

const pathFetchHelpDesc = `
This allows certificates to be fetched. If using the fetch/ prefix any valid certificate can be fetched; if using the revoked/ prefix, which requires a root token, revoked certificates can also be fetched.

Using "ca" or "crl" as the value fetches the appropriate information in DER encoding. Add "/pem" to either to get PEM encoding.
`

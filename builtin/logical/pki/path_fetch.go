// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var pathFetchReadSchema = map[int][]framework.Response{
	http.StatusOK: {{
		Description: "OK",
		Fields: map[string]*framework.FieldSchema{
			"certificate": {
				Type:        framework.TypeString,
				Description: `Certificate`,
				Required:    false,
			},
			"revocation_time": {
				Type:        framework.TypeInt64,
				Description: `Revocation time`,
				Required:    false,
			},
			"revocation_time_rfc3339": {
				Type:        framework.TypeString,
				Description: `Revocation time RFC 3339 formatted`,
				Required:    false,
			},
			"issuer_id": {
				Type:        framework.TypeString,
				Description: `ID of the issuer`,
				Required:    false,
			},
			"ca_chain": {
				Type:        framework.TypeString,
				Description: `Issuing CA Chain`,
				Required:    false,
			},
		},
	}},
}

// Returns the CA in raw format
func pathFetchCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `ca(/pem)?`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "ca-der|ca-pem",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns the CA chain
func pathFetchCAChain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `(cert/)?ca_chain`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "ca-chain-pem|cert-ca-chain",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns the CRL in raw format
func pathFetchCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `crl(/pem|/delta(/pem)?)?`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "crl-der|crl-pem|crl-delta|crl-delta-pem",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns the CRL in raw format
func pathFetchUnifiedCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `unified-crl(/pem|/delta(/pem)?)?`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "unified-crl-der|unified-crl-pem|unified-crl-delta|unified-crl-delta-pem",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathFetchRead,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// Returns any valid (non-revoked) cert in raw format.
func pathFetchValidRaw(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `cert/(?P<serial>[0-9A-Fa-f-:]+)/raw(/pem)?`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "cert-raw-der|cert-raw-pem",
		},

		Fields: map[string]*framework.FieldSchema{
			"serial": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
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

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "cert",
		},

		Fields: map[string]*framework.FieldSchema{
			"serial": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// This returns the CRL in a non-raw format
func pathFetchCRLViaCertPath(b *backend) *framework.Path {
	pattern := `cert/(crl|delta-crl)`
	if constants.IsEnterprise {
		pattern = `cert/(crl|delta-crl|unified-crl|unified-delta-crl)`
	}

	return &framework.Path{
		Pattern: pattern,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "cert-crl|cert-delta-crl|cert-unified-crl|cert-unified-delta-crl",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathFetchRead,
				Responses: pathFetchReadSchema,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

// This returns the list of serial numbers for certs
func pathFetchListCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "certs",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathFetchCertList,
			},
		},

		HelpSynopsis:    pathFetchHelpSyn,
		HelpDescription: pathFetchHelpDesc,
	}
}

func (b *backend) pathFetchCertList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (response *logical.Response, retErr error) {
	entries, err := req.Storage.List(ctx, issuing.PathCerts)
	if err != nil {
		return nil, err
	}
	for i := range entries {
		entries[i] = denormalizeSerial(entries[i])
	}
	return logical.ListResponse(entries), nil
}

func (b *backend) pathFetchRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (response *logical.Response, retErr error) {
	var serial, pemType, contentType string
	var certEntry, revokedEntry *logical.StorageEntry
	var funcErr error
	var certificate []byte
	var fullChain []byte
	var revocationTime int64
	var revocationIssuerId string
	var revocationTimeRfc3339 string

	response = &logical.Response{
		Data: map[string]interface{}{},
	}
	sc := b.makeStorageContext(ctx, req.Storage)

	// Some of these need to return raw and some non-raw;
	// this is basically handled by setting contentType or not.
	// Errors don't cause an immediate exit, because the raw
	// paths still need to return raw output.

	modifiedCtx := &IfModifiedSinceHelper{
		req:       req,
		issuerRef: defaultRef,
	}
	switch {
	case req.Path == "ca" || req.Path == "ca/pem" || req.Path == "cert/ca" || req.Path == "cert/ca/raw" || req.Path == "cert/ca/raw/pem":
		modifiedCtx.reqType = ifModifiedCA
		ret, err := sendNotModifiedResponseIfNecessary(modifiedCtx, sc, response)
		if err != nil || ret {
			retErr = err
			goto reply
		}

		serial = "ca"
		contentType = "application/pkix-cert"
		if req.Path == "ca/pem" || req.Path == "cert/ca/raw/pem" {
			pemType = "CERTIFICATE"
			contentType = "application/pem-certificate-chain"
		} else if req.Path == "cert/ca" {
			pemType = "CERTIFICATE"
			contentType = ""
		}
	case req.Path == "ca_chain" || req.Path == "cert/ca_chain":
		serial = "ca_chain"
		if req.Path == "ca_chain" {
			contentType = "application/pkix-cert"
		}
	case req.Path == "crl" || req.Path == "crl/pem" || req.Path == "crl/delta" || req.Path == "crl/delta/pem" || req.Path == "cert/crl" || req.Path == "cert/crl/raw" || req.Path == "cert/crl/raw/pem" || req.Path == "cert/delta-crl" || req.Path == "cert/delta-crl/raw" || req.Path == "cert/delta-crl/raw/pem" || req.Path == "unified-crl" || req.Path == "unified-crl/pem" || req.Path == "unified-crl/delta" || req.Path == "unified-crl/delta/pem" || req.Path == "cert/unified-crl" || req.Path == "cert/unified-crl/raw" || req.Path == "cert/unified-crl/raw/pem" || req.Path == "cert/unified-delta-crl" || req.Path == "cert/unified-delta-crl/raw" || req.Path == "cert/unified-delta-crl/raw/pem":
		config, err := b.CrlBuilder().GetConfigWithUpdate(sc)
		if err != nil {
			retErr = err
			goto reply
		}
		var isDelta bool
		var isUnified bool
		if strings.Contains(req.Path, "delta") {
			isDelta = true
		}
		if strings.Contains(req.Path, "unified") || shouldLocalPathsUseUnified(config) {
			isUnified = true
		}

		modifiedCtx.reqType = ifModifiedCRL
		if !isUnified && isDelta {
			modifiedCtx.reqType = ifModifiedDeltaCRL
		} else if isUnified && !isDelta {
			modifiedCtx.reqType = ifModifiedUnifiedCRL
		} else if isUnified && isDelta {
			modifiedCtx.reqType = ifModifiedUnifiedDeltaCRL
		}

		ret, err := sendNotModifiedResponseIfNecessary(modifiedCtx, sc, response)
		if err != nil || ret {
			retErr = err
			goto reply
		}

		serial = legacyCRLPath
		if !isUnified && isDelta {
			serial = deltaCRLPath
		} else if isUnified && !isDelta {
			serial = unifiedCRLPath
		} else if isUnified && isDelta {
			serial = unifiedDeltaCRLPath
		}

		contentType = "application/pkix-crl"
		if strings.Contains(req.Path, "pem") {
			pemType = "X509 CRL"
			contentType = "application/x-pem-file"
		} else if req.Path == "cert/crl" || req.Path == "cert/delta-crl" || req.Path == "cert/unified-crl" || req.Path == "cert/unified-delta-crl" {
			pemType = "X509 CRL"
			contentType = ""
		}
	case strings.HasSuffix(req.Path, "/pem") || strings.HasSuffix(req.Path, "/raw"):
		serial = data.Get("serial").(string)
		contentType = "application/pkix-cert"
		if strings.HasSuffix(req.Path, "/pem") {
			pemType = "CERTIFICATE"
			contentType = "application/pem-certificate-chain"
		}
	default:
		if ser, ok := data.GetOk("serial"); ok {
			serial = ser.(string)
		}
		pemType = "CERTIFICATE"
	}
	if len(serial) == 0 {
		response = logical.ErrorResponse("The serial number must be provided")
		goto reply
	}

	// Prefer fetchCAInfo to fetchCertBySerial for CA certificates.
	if serial == "ca_chain" || serial == "ca" {
		caInfo, err := sc.fetchCAInfo(defaultRef, issuing.ReadOnlyUsage)
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

		if serial == "ca_chain" {
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
			certificate = fullChain
		} else if serial == "ca" {
			certificate = caInfo.Certificate.Raw

			if len(pemType) != 0 {
				block := pem.Block{
					Type:  pemType,
					Bytes: certificate,
				}

				// This is convoluted on purpose to ensure that we don't have trailing
				// newlines via various paths
				certificate = []byte(strings.TrimSpace(string(pem.EncodeToMemory(&block))))
			}
		}

		goto reply
	}

	certEntry, funcErr = fetchCertBySerial(sc, req.Path, serial)
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
		block := pem.Block{
			Type:  pemType,
			Bytes: certEntry.Value,
		}
		// This is convoluted on purpose to ensure that we don't have trailing
		// newlines via various paths
		certificate = []byte(strings.TrimSpace(string(pem.EncodeToMemory(&block))))
	}

	revokedEntry, funcErr = fetchCertBySerial(sc, "revoked/", serial)
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
		var revInfo revocation.RevocationInfo
		err := revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error decoding revocation entry for serial %s: %s", serial, err)), nil
		}
		revocationTime = revInfo.RevocationTime
		revocationIssuerId = revInfo.CertificateIssuer.String()

		if !revInfo.RevocationTimeUTC.IsZero() {
			revocationTimeRfc3339 = revInfo.RevocationTimeUTC.Format(time.RFC3339Nano)
		}
	}

reply:
	switch {
	case len(contentType) != 0:
		response = &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: contentType,
				logical.HTTPRawBody:     certificate,
			},
		}
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
		response.Data["revocation_time_rfc3339"] = revocationTimeRfc3339
		// Only output this field if we have a value for it as it doesn't make sense for a
		// bunch of code paths that go through here
		if revocationIssuerId != "" {
			response.Data["issuer_id"] = revocationIssuerId
		}

		if len(fullChain) > 0 {
			response.Data["ca_chain"] = string(fullChain)
		}
	}

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

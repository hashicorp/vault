// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/errutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathCreateCsr() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/csr",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "generate-csr-for-key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Required:    true,
				Description: "Name of the key",
			},
			"version": {
				Type:        framework.TypeInt,
				Required:    false,
				Description: "Optional version of key, 'latest' if not set",
			},
			"csr": {
				Type:     framework.TypeString,
				Required: false,
				Description: `PEM encoded CSR template. The information attributes 
will be used as a basis for the CSR with the key in transit. If not set, an empty CSR is returned.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCreateCsrWrite,
			},
		},
		HelpSynopsis:    pathCreateCsrHelpSyn,
		HelpDescription: pathCreateCsrHelpDesc,
	}
}

func (b *backend) pathImportCertChain() *framework.Path {
	return &framework.Path{
		// NOTE: `set-certificate` or `set_certificate`? Paths seem to use different
		// case, such as `transit/wrapping_key` and `transit/cache-config`.
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/set-certificate",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "set-certificate-for-key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Required:    true,
				Description: "Name of the key",
			},
			"version": {
				Type:        framework.TypeInt,
				Required:    false,
				Description: "Optional version of key, 'latest' if not set",
			},
			"certificate_chain": {
				Type:     framework.TypeString,
				Required: true,
				Description: `PEM encoded certificate chain. It should be composed 
by one or more concatenated PEM blocks and ordered starting from the end-entity certificate.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImportCertChainWrite,
			},
		},
		HelpSynopsis:    pathImportCertChainHelpSyn,
		HelpDescription: pathImportCertChainHelpDesc,
	}
}

func (b *backend) pathCreateCsrWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse(fmt.Sprintf("key with provided name '%s' not found", name)), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false) // NOTE: No lock on "read" operations?
	}
	defer p.Unlock()

	// Check if transit key supports signing
	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type '%s' does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Check if key can be derived
	if p.Derived {
		return logical.ErrorResponse("operation not supported on keys with derivation enabled"), logical.ErrInvalidRequest
	}

	// Transit key version
	signingKeyVersion := p.LatestVersion
	// NOTE: BYOK endpoints seem to remove "v" prefix from version,
	// are versions like that also supported?
	if version, ok := d.GetOk("version"); ok {
		signingKeyVersion = version.(int)
	}

	// Read and parse CSR template
	pemCsrTemplate := d.Get("csr").(string)
	csrTemplate, err := parseCsr(pemCsrTemplate)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	pemCsr, err := p.CreateCsr(signingKeyVersion, csrTemplate)
	if err != nil {
		prefixedErr := fmt.Errorf("could not create the csr: %w", err)
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(prefixedErr.Error()), logical.ErrInvalidRequest
		default:
			return nil, prefixedErr
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"name": p.Name,
			"type": p.Type.String(),
			"csr":  string(pemCsr),
		},
	}

	return resp, nil
}

func (b *backend) pathImportCertChainWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse(fmt.Sprintf("key with provided name '%s' not found", name)), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(true) // NOTE: Lock as we are might write to the policy
	}
	defer p.Unlock()

	// Check if transit key supports signing
	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %s does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Check if key can be derived
	if p.Derived {
		return logical.ErrorResponse("operation not supported on keys with derivation enabled"), logical.ErrInvalidRequest
	}

	// Transit key version
	keyVersion := p.LatestVersion
	if version, ok := d.GetOk("version"); ok {
		keyVersion = version.(int)
	}

	// Get certificate chain
	pemCertChain := d.Get("certificate_chain").(string)
	certChain, err := parseCertificateChain(pemCertChain)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	err = p.ValidateAndPersistCertificateChain(ctx, keyVersion, certChain, req.Storage)
	if err != nil {
		prefixedErr := fmt.Errorf("failed to persist certificate chain: %w", err)
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(prefixedErr.Error()), logical.ErrInvalidRequest
		default:
			return nil, prefixedErr
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":              p.Name,
			"type":              p.Type.String(),
			"certificate-chain": pemCertChain,
		},
	}

	return resp, nil
}

func parseCsr(csrStr string) (*x509.CertificateRequest, error) {
	if csrStr == "" {
		return &x509.CertificateRequest{}, nil
	}

	block, _ := pem.Decode([]byte(csrStr))
	if block == nil {
		return nil, errors.New("could not decode PEM certificate request")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

func parseCertificateChain(certChainString string) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	var pemCertBlocks []*pem.Block
	pemBytes := []byte(strings.TrimSpace(certChainString))
	for len(pemBytes) > 0 {
		var pemCertBlock *pem.Block
		pemCertBlock, pemBytes = pem.Decode(pemBytes)
		if pemCertBlock == nil {
			return nil, errors.New("could not decode PEM block in certificate chain")
		}

		switch pemCertBlock.Type {
		case "CERTIFICATE", "X05 CERTIFICATE":
			pemCertBlocks = append(pemCertBlocks, pemCertBlock)
		default:
			// Ignore any other entries
		}
	}

	if len(pemCertBlocks) == 0 {
		return nil, errors.New("provided certificate chain did not contain any valid PEM certificate")
	}

	for _, certBlock := range pemCertBlocks {
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate in certificate chain: %w", err)
		}

		certificates = append(certificates, cert)
	}

	return certificates, nil
}

const pathCreateCsrHelpSyn = `Create a CSR from a key in transit`

const pathCreateCsrHelpDesc = `This path is used to create a CSR from a key in 
transit. If a CSR template is provided, its significant information, expect key 
related data, are included in the CSR otherwise an empty CSR is returned.
`

const pathImportCertChainHelpSyn = `Imports an externally-signed certificate 
chain into an existing key version`

const pathImportCertChainHelpDesc = `This path is used to import an externally-
signed certificate chain into a key in transit. The leaf certificate key has to 
match the selected key in transit.
`

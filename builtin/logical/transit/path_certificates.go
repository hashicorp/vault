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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathCreateCsr() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/csr",
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
are going to be used as a basis for the CSR with the key in transit. If not set, an empty CSR is returned.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			// NOTE: Create and Update?
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathCreateCsrWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "create",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCreateCsrWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "update",
				},
			},
		},
		HelpSynopsis:    pathCreateCsrHelpSyn,
		HelpDescription: pathCreateCsrHelpDesc,
	}
}

func (b *backend) pathImportCertChain() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/set-certificate",
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
				// FIXME: Complete description
				Description: `PEM encoded certificate chain.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			// NOTE: Create and Update?
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathImportCertChainWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "create",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImportCertChainWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "update",
				},
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

	// Read and parse CSR template
	pemCsrTemplate := d.Get("csr").(string)
	csrTemplate, err := parseCsr(pemCsrTemplate)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	signingKeyVersion := p.LatestVersion
	// NOTE: BYOK endpoints seem to remove "v" prefix from version,
	// are versions like that also supported?
	if version, ok := d.GetOk("version"); ok {
		signingKeyVersion = version.(int)
	}

	pemCsr, err := p.CreateCsr(signingKeyVersion, csrTemplate)
	if err != nil {
		return nil, err
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
	// NOTE: A key type that doesn't support signing cannot possible (?) have
	// a certificate, so does it make sense to have this check?
	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %s does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Get certificate chain
	pemCertChain := d.Get("certificate_chain").(string)
	certChain, err := parseCertificateChain(pemCertChain)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	err = validateLeafCertPosition(certChain)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	keyVersion := p.LatestVersion
	if version, ok := d.GetOk("version"); ok {
		keyVersion = version.(int)
	}

	leafCertPublicKeyAlgorithm := certChain[0].PublicKeyAlgorithm
	var keyTypeMatches bool
	switch p.Type {
	case keysutil.KeyType_ECDSA_P256, keysutil.KeyType_ECDSA_P384, keysutil.KeyType_ECDSA_P521:
		if leafCertPublicKeyAlgorithm == x509.ECDSA {
			keyTypeMatches = true
		}
	case keysutil.KeyType_ED25519:
		if leafCertPublicKeyAlgorithm == x509.Ed25519 {
			keyTypeMatches = true
		}
	case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA3072, keysutil.KeyType_RSA4096:
		if leafCertPublicKeyAlgorithm == x509.RSA {
			keyTypeMatches = true
		}
	}
	if !keyTypeMatches {
		// NOTE: Different type "names" might lead to confusion.
		return logical.ErrorResponse(fmt.Sprintf("provided leaf certificate public key type '%s' does not match the transit key type '%s'",
			leafCertPublicKeyAlgorithm.String(), p.Type.String())), logical.ErrInvalidRequest
	}

	// Validate if leaf cert key matches with transit key
	valid, err := p.ValidateLeafCertKeyMatch(keyVersion, leafCertPublicKeyAlgorithm, certChain[0].PublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not validate key match between leaf certificate key and key version in transit: %s", err.Error())
	}
	if !valid {
		return logical.ErrorResponse("leaf certificate public key does match the key version selected"), logical.ErrInvalidRequest
	}

	err = p.PersistCertificateChain(ctx, keyVersion, certChain, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to persist certificate chain: %s", err.Error())
	}

	return nil, nil
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
			return nil, fmt.Errorf("failed to parse certificate in certificate chain: %s", err.Error())
		}

		certificates = append(certificates, cert)
	}

	return certificates, nil
}

func validateLeafCertPosition(certChain []*x509.Certificate) error {
	// NOTE: Validate if Basic Constraints are valid in different condition and return different error message?
	if len(certChain) > 0 && (certChain[0].BasicConstraintsValid && certChain[0].IsCA) {
		return errors.New("leaf certificate not found in the first position of the certificate chain")
	}

	for _, cert := range certChain[1:] {
		if cert.BasicConstraintsValid && !cert.IsCA {
			return errors.New("provided certificate chain contains more than one leaf certificate")
		}
	}

	return nil
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

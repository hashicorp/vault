// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathSignCsr() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/csr",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type: framework.TypeString,
				// NOTE: Required seems to be deprected, still keep it to "improve" readability?
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
				Callback: b.pathSignCsrWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "create",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathSignCsrWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "update",
				},
			},
		},
		// FIXME: Write synposis and description
		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (b *backend) pathSetCertificate() *framework.Path {
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
				Callback: b.pathSetCertificateWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "create",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathSetCertificateWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "update",
				},
			},
		},
		// FIXME: Write synposis and description
		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (b *backend) pathSignCsrWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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
		return logical.ErrorResponse(fmt.Sprintf("key type '%v' does not support signing", p.Type)), logical.ErrInvalidRequest
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

	// FIXME: Remove
	// log.Println("signingKeyVersion: ", signingKeyVersion)

	pemCsr, err := p.SignCsr(signingKeyVersion, csrTemplate)
	if err != nil {
		// FIXME: Error returned
		return nil, err
	}

	// FIXME: Remove
	// log.Printf("CSR:\n%s", pemCsr)

	// `csr` is returned as base64, what's wrong?
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name": p.Name,
			"type": p.Type.String(),
			"csr":  pemCsr,
		},
	}

	return resp, nil
}

func (b *backend) pathSetCertificateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// NOTE: Is this used in multiple places?
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
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Get certificate chain
	pemCertChain := d.Get("certificate_chain").(string)
	certChain, err := parseCertificateChain(pemCertChain)
	if err != nil {
		// FIXME: Error
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Validate there's only a leaf certificate in the chain
	// NOTE: Does it have the be the first? Would make sense
	if hasSingleLeafCert := hasSingleLeafCertificate(certChain); !hasSingleLeafCert {
		return logical.ErrorResponse("expected a single leaf certificate in the certificate chain"), logical.ErrInvalidRequest
	}

	keyVersion := p.LatestVersion
	if version, ok := d.GetOk("version"); ok {
		keyVersion = version.(int)
	}

	// Validate if leaf cert key matches with transit key
	valid, err := p.ValidateLeafCertKeyMatch(keyVersion, certChain[0].PublicKeyAlgorithm, certChain[0].PublicKey)
	if err != nil {
		return nil, err
	}
	if !valid {
		return logical.ErrorResponse("end entity certificate public key does match the transit's public key"), logical.ErrInvalidRequest
	}

	p.PersistCertificateChain(keyVersion, certChain, req.Storage)
	if err != nil {
		// FIXME: should be internal, message
		return nil, errors.New("failed to persist certificate chain")
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

func parseCertificateChain(certChainStr string) ([]*x509.Certificate, error) {
	certificates := []*x509.Certificate{}

	pemCertBlocks := []*pem.Block{}
	rest := []byte(certChainStr)
	for len(rest) != 0 {
		var pemCertBlock *pem.Block
		pemCertBlock, rest = pem.Decode([]byte(rest))
		if pemCertBlock == nil {
			return nil, errors.New("could not decode certificate in certificate chain")
		}

		pemCertBlocks = append(pemCertBlocks, pemCertBlock)
	}

	if len(pemCertBlocks) == 0 {
		return nil, errors.New("no certificates provided in `certificate_chain` parameter")
	}

	// NOTE: This approach or x509.ParseCertificates?
	for _, certBlock := range pemCertBlocks {
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, errors.New("failed to parse certificate in certificate chain")
		}

		certificates = append(certificates, cert)
	}

	return certificates, nil
}

func hasSingleLeafCertificate(certChain []*x509.Certificate) bool {
	var leafCertsCount uint8
	for _, cert := range certChain {
		// FIXME: Remove
		log.Printf("BasicConstraintValid: %t  | IsCA: %t\n", cert.BasicConstraintsValid, cert.IsCA)
		if cert.BasicConstraintsValid && !cert.IsCA {
			leafCertsCount += 1
		}
	}

	// FIXME: Remove
	log.Printf("leafCertsCount: %d\n", leafCertsCount)
	var hasSingleLeafCert bool
	if leafCertsCount == 1 {
		hasSingleLeafCert = true
	}

	return hasSingleLeafCert
}

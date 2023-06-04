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

// NOTE: Or only `pathCsr`?
func (b *backend) pathSignCsr() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/csr",
		// NOTE: Any other field?
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Required:    true,
				Description: "The name of the key",
			},
			"version": {
				Type:     framework.TypeInt,
				Required: false,
				// FIXME: Add description
				Description: `If not set, 'latest' is used.`,
			},
			"csr": {
				Type:     framework.TypeString,
				Required: false,
				// FIXME: Add description
				Description: ``,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			// NOTE: Create and Update?
			logical.CreateOperation: b.pathSignCsrWrite,
			logical.UpdateOperation: b.pathSignCsrWrite,
		},
		// FIXME: Write synposis and description
		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (b *backend) pathSetCertificate() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/set-certificate",
		// NOTE: Any other field?
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Required:    true,
				Description: "The name of the key",
			},
			"version": {
				Type:     framework.TypeInt,
				Required: false,
				// FIXME: Add description
				Description: `If not set, 'latest' is used.`,
			},
			"certificate_chain": {
				Type:     framework.TypeString,
				Required: true,
				// FIXME: Add description
				Description: ``,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			// NOTE: Create and Update?
			logical.CreateOperation: b.pathSetCertificateWrite,
			logical.UpdateOperation: b.pathSetCertificateWrite,
		},
		// FIXME: Write synposis and description
		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

// NOTE: d or data for the framework.FieldData argument?
func (b *backend) pathSignCsrWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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
		p.Lock(false) // NOTE: No lock on "read" operations?
	}
	defer p.Unlock()

	// Check if transit key supports signing
	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Get CSR template
	// NOTE: Use GetOk, or GetErrOk?
	csr := d.Get("csr").(string)
	csrTemplate, err := parseCsrParam(csr)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	signingKeyVersion := p.LatestVersion
	if version, ok := d.GetOk("version"); ok {
		signingKeyVersion = version.(int)
	}

	// FIXME: Remove
	log.Println("signingKeyVersion: ", signingKeyVersion)

	csrBytes, err := p.SignCsr(signingKeyVersion, csrTemplate)
	if err != nil {
		// FIXME: Error returned
		return nil, err
	}

	// FIXME: Remove
	log.Printf("CSR:\n%s", csrBytes)

	// Getting response as base64, is that expected?
	resp := &logical.Response{
		Data: map[string]interface{}{
			"csr": csrBytes,
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
	// FIXME: THIS IS PEM (fix name)
	certChainBytes := d.Get("certificate_chain").(string)
	certChain, err := parseParamCertificateChain(certChainBytes)
	if err != nil {
		// FIXME: Error
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	log.Printf("For some reason we are here\nHere's the number of certificates in the chain: %d", len(certChain))

	// NOTE: Is this check needed?
	if len(certChain) == 0 {
		return logical.ErrorResponse("no certificates provided"), logical.ErrInvalidRequest
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

func parseCsrParam(csr string) (*x509.CertificateRequest, error) {
	if csr == "" {
		return &x509.CertificateRequest{}, nil
	}

	block, _ := pem.Decode([]byte(csr))
	if block == nil {
		return nil, errors.New("failed to decode PEM Certificate Request")
	}

	csrTemplate, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	return csrTemplate, nil
}

// FIXME: Names
func parseParamCertificateChain(certChain string) ([]*x509.Certificate, error) {
	certPemBlocks, _ := pem.Decode([]byte(certChain))
	if certPemBlocks == nil {
		return nil, errors.New("failed to decode PEM Certificate Chain")
	}

	certificates, err := x509.ParseCertificates(certPemBlocks.Bytes)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

func hasSingleLeafCertificate(certChain []*x509.Certificate) bool {
	var leafCertsCount int
	for _, cert := range certChain {
		log.Printf("BasicConstraintValid: %t  | IsCA: %t\n", cert.BasicConstraintsValid, cert.IsCA)
		if cert.BasicConstraintsValid && !cert.IsCA {
			leafCertsCount += 1
		}
	}

	log.Printf("leafCertsCount: %d\n", leafCertsCount)
	var hasSingleLeafCert bool
	if leafCertsCount == 1 {
		hasSingleLeafCert = true
	}

	return hasSingleLeafCert
}

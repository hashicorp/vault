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

// NOTE: d or data for the framework.Fielddata argument?
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

	// Is this check relevant in this scenario?
	// NOTE: It will fail for the empty template CSR
	// if err = csrTemplate.CheckSignature(); err != nil {
	// 	// FIXME: Error returned
	// 	return logical.ErrorResponse("invalid signature in csr provided"), logical.ErrInvalidRequest
	// }

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

	resp := &logical.Response{
		Data: map[string]interface{}{
			"csr": csrBytes,
		},
	}

	return resp, nil
}

func parseCsrParam(csr string) (*x509.CertificateRequest, error) {
	if csr == "" {
		return &x509.CertificateRequest{}, nil
	}

	block, _ := pem.Decode([]byte(csr))
	if block == nil {
		return nil, errors.New("failed to decode CSR PEM")
	}

	csrTemplate, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	return csrTemplate, nil
}

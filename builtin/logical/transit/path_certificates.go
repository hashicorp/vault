// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	// "crypto/rand"
	// "crypto/x509"
	// "encoding/pem"
	"fmt"
	"log"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

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
			// NOTE: Or logical.CreateOperation?
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
		// NOTE: Return custom error or err?
		return nil, fmt.Errorf("no key found with name %s; to import a new key, use the import/ endpoint", name)
	}
	if !b.System().CachingDisabled() {
		p.Lock(false) // NOTE: No lock on "read" operations?
	}
	defer p.Unlock()

	// Check if transit key supports signing
	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	// Create certificate template
	// var csrBytes []bytes
	// csr, isCsrSet := d.GetOk("csr")
	// if isCsrSet {
	// 	csrBlock, _ := pem.Decode(csr.([]byte))
	// 	if csrBlock == nil {
	// 		return logical.ErrorResponse(fmt.Sprintf("invalid csr provided")), logical.ErrInvalidRequest
	// 	}
	// } else {
	// 	csrBytes, err := x509.CreateCertificateRequest(rand.Reader)
	// }
	//
	// csr, err := x509.ParseCertificateRequest(csrBytes)
	// if err != nil {
	// 	return logical.ErrorResponse(fmt.Sprintf("failed to parse csr")), logical.ErrInvalidRequest
	// }

	signingKeyVersion := p.LatestVersion
	if version, ok := d.GetOk("version"); ok {
		signingKeyVersion = version.(int)
	}

	log.Printf("Signing key version: %d", signingKeyVersion)

	// var csrBytes []byte
	// csr, csrSet := d.GetOk("csr")
	// if !csrSet {
	// 	rsaPrivKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	// 	csrBytes = []byte("")
	// 	var err error
	// 	csrTemplate := &x509.CertificateRequest{}
	// 	csrBytes, err = x509.CreateCertificateRequest(rand.Reader, csrTemplate, rsaPrivKey)
	// 	if err != nil {
	// 		log.Printf("ERROR: Failed to create CSR: %v", err)
	// 	}
	//
	// 	csrPem := pem.EncodeToMemory(
	// 		&pem.Block{
	// 			Type:  "CERTIFICATE REQUEST",
	// 			Bytes: csrBytes,
	// 		},
	// 	)
	// 	log.Printf("CSR PEM:\n %s\n", csrPem)
	// } else {
	// 	csrBytes = csr.([]byte)
	// }

	// log.Printf("Empty CSR: %s\n", csrBytes)
	return nil, nil
}

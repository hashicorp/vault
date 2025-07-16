// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathListCertificates creates a path that enables listing of all
// the AWS public certificates registered with Vault.
func (b *backend) pathListCertificates() *framework.Path {
	return &framework.Path{
		Pattern: "config/certificates/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "certificate-configurations",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathCertificatesList,
			},
		},

		HelpSynopsis:    pathListCertificatesHelpSyn,
		HelpDescription: pathListCertificatesHelpDesc,
	}
}

func (b *backend) pathConfigCertificate() *framework.Path {
	return &framework.Path{
		Pattern: "config/certificate/" + framework.GenericNameRegex("cert_name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
		},

		Fields: map[string]*framework.FieldSchema{
			"cert_name": {
				Type:        framework.TypeString,
				Description: "Name of the certificate.",
			},
			"aws_public_cert": {
				Type:        framework.TypeString,
				Description: "Base64 encoded AWS Public cert required to verify PKCS7 signature of the EC2 instance metadata.",
			},
			"type": {
				Type:    framework.TypeString,
				Default: "pkcs7",
				Description: `
Takes the value of either "pkcs7" or "identity", indicating the type of
document which can be verified using the given certificate. The reason is that
the PKCS#7 document will have a DSA digest and the identity signature will have
an RSA signature, and accordingly the public certificates to verify those also
vary. Defaults to "pkcs7".`,
			},
		},

		ExistenceCheck: b.pathConfigCertificateExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigCertificateCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "certificate",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigCertificateCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "certificate",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigCertificateRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "certificate-configuration",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigCertificateDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "certificate-configuration",
				},
			},
		},

		HelpSynopsis:    pathConfigCertificateSyn,
		HelpDescription: pathConfigCertificateDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigCertificateExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return false, fmt.Errorf("missing cert_name")
	}

	entry, err := b.lockedAWSPublicCertificateEntry(ctx, req.Storage, certName)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// pathCertificatesList is used to list all the AWS public certificates registered with Vault
func (b *backend) pathCertificatesList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	certs, err := req.Storage.List(ctx, "config/certificate/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(certs), nil
}

// Decodes the PEM encoded certificate and parses it into a x509 cert
func decodePEMAndParseCertificate(certificate string) (*x509.Certificate, error) {
	// Decode the PEM block and error out if a block is not detected in the first attempt
	decodedPublicCert, rest := pem.Decode([]byte(certificate))
	if len(rest) != 0 {
		return nil, fmt.Errorf("invalid certificate; should be one PEM block only")
	}

	// Check if the certificate can be parsed
	publicCert, err := x509.ParseCertificate(decodedPublicCert.Bytes)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return nil, fmt.Errorf("invalid certificate; failed to parse certificate")
	}
	return publicCert, nil
}

// awsPublicCertificates returns a slice of all the parsed AWS public
// certificates, which are used to verify either the identity, RSA 2048
// or the PKCS7 signatures of the instance identity documents. This method will
// append the certificates registered using `config/certificate/<cert_name>`
// endpoint, along with the default certificates in the backend.
func (b *backend) awsPublicCertificates(ctx context.Context, s logical.Storage, isPkcs bool) ([]*x509.Certificate, error) {
	// Lock at beginning and use internal method so that we are consistent as
	// we iterate through
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	certs := make([]*x509.Certificate, len(defaultCertificates))
	copy(certs, defaultCertificates)

	// Get the list of all the registered certificates
	registeredCerts, err := s.List(ctx, "config/certificate/")
	if err != nil {
		return nil, err
	}

	// Iterate through each certificate, parse and append it to a slice
	for _, cert := range registeredCerts {
		certEntry, err := b.nonLockedAWSPublicCertificateEntry(ctx, s, cert)
		if err != nil {
			return nil, err
		}
		if certEntry == nil {
			return nil, fmt.Errorf("certificate storage has a nil entry under the name: %q", cert)
		}
		// Append relevant certificates only
		if (isPkcs && certEntry.Type == "pkcs7") ||
			(!isPkcs && certEntry.Type == "identity") {
			decodedCert, err := decodePEMAndParseCertificate(certEntry.AWSPublicCert)
			if err != nil {
				return nil, err
			}
			certs = append(certs, decodedCert)
		}
	}

	return certs, nil
}

// lockedSetAWSPublicCertificateEntry is used to store the AWS public key in
// the storage. This method acquires lock before creating or updating a storage
// entry.
func (b *backend) lockedSetAWSPublicCertificateEntry(ctx context.Context, s logical.Storage, certName string, certEntry *awsPublicCert) error {
	if certName == "" {
		return fmt.Errorf("missing certificate name")
	}

	if certEntry == nil {
		return fmt.Errorf("nil AWS public key certificate")
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return b.nonLockedSetAWSPublicCertificateEntry(ctx, s, certName, certEntry)
}

// nonLockedSetAWSPublicCertificateEntry is used to store the AWS public key in
// the storage. This method does not acquire lock before reading the storage.
// If locking is desired, use lockedSetAWSPublicCertificateEntry instead.
func (b *backend) nonLockedSetAWSPublicCertificateEntry(ctx context.Context, s logical.Storage, certName string, certEntry *awsPublicCert) error {
	if certName == "" {
		return fmt.Errorf("missing certificate name")
	}

	if certEntry == nil {
		return fmt.Errorf("nil AWS public key certificate")
	}

	entry, err := logical.StorageEntryJSON("config/certificate/"+certName, certEntry)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("failed to create storage entry for AWS public key certificate")
	}

	return s.Put(ctx, entry)
}

// lockedAWSPublicCertificateEntry is used to get the configured AWS Public Key
// that is used to verify the PKCS#7 signature of the instance identity
// document.
func (b *backend) lockedAWSPublicCertificateEntry(ctx context.Context, s logical.Storage, certName string) (*awsPublicCert, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedAWSPublicCertificateEntry(ctx, s, certName)
}

// nonLockedAWSPublicCertificateEntry reads the certificate information from
// the storage. This method does not acquire lock before reading the storage.
// If locking is desired, use lockedAWSPublicCertificateEntry instead.
func (b *backend) nonLockedAWSPublicCertificateEntry(ctx context.Context, s logical.Storage, certName string) (*awsPublicCert, error) {
	entry, err := s.Get(ctx, "config/certificate/"+certName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var certEntry awsPublicCert
	if err := entry.DecodeJSON(&certEntry); err != nil {
		return nil, err
	}

	// Handle upgrade for certificate type
	persistNeeded := false
	if certEntry.Type == "" {
		certEntry.Type = "pkcs7"
		persistNeeded = true
	}

	if persistNeeded {
		if err := b.nonLockedSetAWSPublicCertificateEntry(ctx, s, certName, &certEntry); err != nil {
			return nil, err
		}
	}

	return &certEntry, nil
}

// pathConfigCertificateDelete is used to delete the previously configured AWS
// Public Key that is used to verify the PKCS#7 signature of the instance
// identity document.
func (b *backend) pathConfigCertificateDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing cert_name"), nil
	}

	return nil, req.Storage.Delete(ctx, "config/certificate/"+certName)
}

// pathConfigCertificateRead is used to view the configured AWS Public Key that
// is used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing cert_name"), nil
	}

	certificateEntry, err := b.lockedAWSPublicCertificateEntry(ctx, req.Storage, certName)
	if err != nil {
		return nil, err
	}
	if certificateEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"aws_public_cert": certificateEntry.AWSPublicCert,
			"type":            certificateEntry.Type,
		},
	}, nil
}

// pathConfigCertificateCreateUpdate is used to register an AWS Public Key that
// is used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing certificate name"), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// Check if there is already a certificate entry registered
	certEntry, err := b.nonLockedAWSPublicCertificateEntry(ctx, req.Storage, certName)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		certEntry = &awsPublicCert{}
	}

	// Check if type information is provided
	certTypeRaw, ok := data.GetOk("type")
	if ok {
		certEntry.Type = strings.ToLower(certTypeRaw.(string))
	} else if req.Operation == logical.CreateOperation {
		certEntry.Type = data.Get("type").(string)
	}

	switch certEntry.Type {
	case "pkcs7":
	case "identity":
	default:
		return logical.ErrorResponse(fmt.Sprintf("invalid certificate type %q", certEntry.Type)), nil
	}

	// Check if the value is provided by the client
	certStrData, ok := data.GetOk("aws_public_cert")
	if ok {
		if certBytes, err := base64.StdEncoding.DecodeString(certStrData.(string)); err == nil {
			certEntry.AWSPublicCert = string(certBytes)
		} else {
			certEntry.AWSPublicCert = certStrData.(string)
		}
	} else {
		// aws_public_cert should be supplied for both create and update operations.
		// If it is not provided, throw an error.
		return logical.ErrorResponse("missing aws_public_cert"), nil
	}

	// If explicitly set to empty string, error out
	if certEntry.AWSPublicCert == "" {
		return logical.ErrorResponse("invalid aws_public_cert"), nil
	}

	// Verify the certificate by decoding it and parsing it
	publicCert, err := decodePEMAndParseCertificate(certEntry.AWSPublicCert)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return logical.ErrorResponse("invalid certificate; failed to decode and parse certificate"), nil
	}

	// If none of the checks fail, save the provided certificate
	if err := b.nonLockedSetAWSPublicCertificateEntry(ctx, req.Storage, certName, certEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Struct awsPublicCert holds the AWS Public Key that is used to verify the PKCS#7 signature
// of the instance identity document.
type awsPublicCert struct {
	AWSPublicCert string `json:"aws_public_cert"`
	Type          string `json:"type"`
}

const pathConfigCertificateSyn = `
Adds the AWS Public Key that is used to verify the PKCS#7 signature of the identity document.
`

const pathConfigCertificateDesc = `
AWS Public Key which is used to verify the PKCS#7 signature of the identity document,
varies by region. The public key(s) can be found in AWS EC2 instance metadata documentation.
The default key that is used to verify the signature is the one that is applicable for
following regions: US East (N. Virginia), US West (Oregon), US West (N. California),
EU (Ireland), EU (Frankfurt), Asia Pacific (Tokyo), Asia Pacific (Seoul), Asia Pacific (Singapore),
Asia Pacific (Sydney), and South America (Sao Paulo).

If the instances belongs to region other than the above, the public key(s) for the
corresponding regions should be registered using this endpoint. PKCS#7 is verified
using a collection of certificates containing the default certificate and all the
certificates that are registered using this endpoint.
`

const pathListCertificatesHelpSyn = `
Lists all the AWS public certificates that are registered with the backend.
`

const pathListCertificatesHelpDesc = `
Certificates will be listed by their respective names that were used during registration.
`

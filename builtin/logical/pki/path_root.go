package pki

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathGenerateRoot(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root/generate/" + framework.GenericNameRegex("exported"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathCAGenerateRoot,
		},

		HelpSynopsis:    pathGenerateRootHelpSyn,
		HelpDescription: pathGenerateRootHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAKeyGenerationFields(ret.Fields)
	ret.Fields = addCAIssueFields(ret.Fields)

	return ret
}

func pathDeleteRoot(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCADeleteRoot,
		},

		HelpSynopsis:    pathDeleteRootHelpSyn,
		HelpDescription: pathDeleteRootHelpDesc,
	}

	return ret
}

func pathSignIntermediate(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root/sign-intermediate",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathCASignIntermediate,
		},

		HelpSynopsis:    pathSignIntermediateHelpSyn,
		HelpDescription: pathSignIntermediateHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAIssueFields(ret.Fields)

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	ret.Fields["use_csr_values"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
		Description: `If true, then:
1) Subject information, including names and alternate
names, will be preserved from the CSR rather than
using values provided in the other parameters to
this path;
2) Any key usages requested in the CSR will be
added to the basic set of key usages used for CA
certs signed by this path; for instance,
the non-repudiation flag.`,
	}

	return ret
}

func pathSignSelfIssued(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root/sign-self-issued",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathCASignSelfIssued,
		},

		Fields: map[string]*framework.FieldSchema{
			"certificate": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `PEM-format self-issued certificate to be signed.`,
			},
		},

		HelpSynopsis:    pathSignSelfIssuedHelpSyn,
		HelpDescription: pathSignSelfIssuedHelpDesc,
	}

	return ret
}

func (b *backend) pathCADeleteRoot(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, req.Storage.Delete(ctx, "config/ca_bundle")
}

func (b *backend) pathCAGenerateRoot(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	entry, err := req.Storage.Get(ctx, "config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		resp := &logical.Response{}
		resp.AddWarning(fmt.Sprintf("Refusing to generate a root certificate over an existing root certificate. If you really want to destroy the original root certificate, please issue a delete against %sroot.", req.MountPoint))
		return resp, nil
	}

	exported, format, role, errorResp := b.getGenerationParams(data)
	if errorResp != nil {
		return errorResp, nil
	}

	maxPathLengthIface, ok := data.GetOk("max_path_length")
	if ok {
		maxPathLength := maxPathLengthIface.(int)
		role.MaxPathLength = &maxPathLength
	}

	input := &dataBundle{
		req:     req,
		apiData: data,
		role:    role,
	}
	parsedBundle, err := generateCert(ctx, b, input, true)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw cert bundle to cert bundle: {{err}}", err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
			"serial_number": cb.SerialNumber,
		},
	}

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = cb.Certificate
		if exported {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}

	case "pem_bundle":
		resp.Data["issuing_ca"] = cb.Certificate

		if exported {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
			resp.Data["certificate"] = fmt.Sprintf("%s\n%s", cb.PrivateKey, cb.Certificate)
		} else {
			resp.Data["certificate"] = cb.Certificate
		}

	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		if exported {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}
	}

	if data.Get("private_key_format").(string) == "pkcs8" {
		err = convertRespToPKCS8(resp)
		if err != nil {
			return nil, err
		}
	}

	// Store it as the CA bundle
	entry, err = logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Also store it as just the certificate identified by serial number, so it
	// can be revoked
	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "certs/" + normalizeSerial(cb.SerialNumber),
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, errwrap.Wrapf("unable to store certificate locally: {{err}}", err)
	}

	// For ease of later use, also store just the certificate at a known
	// location
	entry.Key = "ca"
	entry.Value = parsedBundle.CertificateBytes
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Build a fresh CRL
	err = buildCRL(ctx, b, req)
	if err != nil {
		return nil, err
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the generated certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	return resp, nil
}

func (b *backend) pathCASignIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem" or "der"`,
		), nil
	}

	role := &roleEntry{
		OU:                    data.Get("ou").([]string),
		Organization:          data.Get("organization").([]string),
		Country:               data.Get("country").([]string),
		Locality:              data.Get("locality").([]string),
		Province:              data.Get("province").([]string),
		StreetAddress:         data.Get("street_address").([]string),
		PostalCode:            data.Get("postal_code").([]string),
		TTL:                   time.Duration(data.Get("ttl").(int)) * time.Second,
		AllowLocalhost:        true,
		AllowAnyName:          true,
		AllowIPSANs:           true,
		EnforceHostnames:      false,
		KeyType:               "any",
		AllowedURISANs:        []string{"*"},
		AllowedSerialNumbers:  []string{"*"},
		AllowExpirationPastCA: true,
	}

	if cn := data.Get("common_name").(string); len(cn) == 0 {
		role.UseCSRCommonName = true
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(ctx, req)
	switch caErr.(type) {
	case errutil.UserError:
		return nil, errutil.UserError{Err: fmt.Sprintf(
			"could not fetch the CA certificate (was one set?): %s", caErr)}
	case errutil.InternalError:
		return nil, errutil.InternalError{Err: fmt.Sprintf(
			"error fetching CA certificate: %s", caErr)}
	}

	useCSRValues := data.Get("use_csr_values").(bool)

	maxPathLengthIface, ok := data.GetOk("max_path_length")
	if ok {
		maxPathLength := maxPathLengthIface.(int)
		role.MaxPathLength = &maxPathLength
	}

	input := &dataBundle{
		req:           req,
		apiData:       data,
		signingBundle: signingBundle,
		role:          role,
	}
	parsedBundle, err := signCert(b, input, true, useCSRValues)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}

	if err := parsedBundle.Verify(); err != nil {
		return nil, errwrap.Wrapf("verification of parsed bundle failed: {{err}}", err)
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw signing bundle to cert bundle: {{err}}", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw cert bundle to cert bundle: {{err}}", err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
			"serial_number": cb.SerialNumber,
		},
	}

	if signingBundle.Certificate.NotAfter.Before(parsedBundle.Certificate.NotAfter) {
		resp.AddWarning("The expiration time for the signed certificate is after the CA's expiration time. If the new certificate is not treated as a root, validation paths with the certificate past the issuing CA's expiration time will fail.")
	}

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = signingCB.Certificate
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			resp.Data["ca_chain"] = cb.CAChain
		}

	case "pem_bundle":
		resp.Data["certificate"] = cb.ToPEMBundle()
		resp.Data["issuing_ca"] = signingCB.Certificate
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			resp.Data["ca_chain"] = cb.CAChain
		}

	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(signingBundle.CertificateBytes)

		var caChain []string
		for _, caCert := range parsedBundle.CAChain {
			caChain = append(caChain, base64.StdEncoding.EncodeToString(caCert.Bytes))
		}
		if caChain != nil && len(caChain) > 0 {
			resp.Data["ca_chain"] = cb.CAChain
		}
	}

	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "certs/" + normalizeSerial(cb.SerialNumber),
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, errwrap.Wrapf("unable to store certificate locally: {{err}}", err)
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the signed certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	return resp, nil
}

func (b *backend) pathCASignSelfIssued(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	certPem := data.Get("certificate").(string)
	block, _ := pem.Decode([]byte(certPem))
	if block == nil || len(block.Bytes) == 0 {
		return logical.ErrorResponse("certificate could not be PEM-decoded"), nil
	}
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error parsing certificate: %s", err)), nil
	}
	if len(certs) != 1 {
		return logical.ErrorResponse(fmt.Sprintf("%d certificates found in PEM file, expected 1", len(certs))), nil
	}

	cert := certs[0]
	if !cert.IsCA {
		return logical.ErrorResponse("given certificate is not a CA certificate"), nil
	}
	if !reflect.DeepEqual(cert.Issuer, cert.Subject) {
		return logical.ErrorResponse("given certificate is not self-issued"), nil
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(ctx, req)
	switch caErr.(type) {
	case errutil.UserError:
		return nil, errutil.UserError{Err: fmt.Sprintf(
			"could not fetch the CA certificate (was one set?): %s", caErr)}
	case errutil.InternalError:
		return nil, errutil.InternalError{Err: fmt.Sprintf(
			"error fetching CA certificate: %s", caErr)}
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw signing bundle to cert bundle: {{err}}", err)
	}

	urls := &urlEntries{}
	if signingBundle.URLs != nil {
		urls = signingBundle.URLs
	}
	cert.IssuingCertificateURL = urls.IssuingCertificates
	cert.CRLDistributionPoints = urls.CRLDistributionPoints
	cert.OCSPServer = urls.OCSPServers

	newCert, err := x509.CreateCertificate(rand.Reader, cert, signingBundle.Certificate, cert.PublicKey, signingBundle.PrivateKey)
	if err != nil {
		return nil, errwrap.Wrapf("error signing self-issued certificate: {{err}}", err)
	}
	if len(newCert) == 0 {
		return nil, fmt.Errorf("nil cert was created when signing self-issued certificate")
	}
	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCert,
	})

	return &logical.Response{
		Data: map[string]interface{}{
			"certificate": strings.TrimSpace(string(pemCert)),
			"issuing_ca":  signingCB.Certificate,
		},
	}, nil
}

const pathGenerateRootHelpSyn = `
Generate a new CA certificate and private key used for signing.
`

const pathGenerateRootHelpDesc = `
See the API documentation for more information.
`

const pathDeleteRootHelpSyn = `
Deletes the root CA key to allow a new one to be generated.
`

const pathDeleteRootHelpDesc = `
See the API documentation for more information.
`

const pathSignIntermediateHelpSyn = `
Issue an intermediate CA certificate based on the provided CSR.
`

const pathSignIntermediateHelpDesc = `
see the API documentation for more information.
`

const pathSignSelfIssuedHelpSyn = `
Signs another CA's self-issued certificate.
`

const pathSignSelfIssuedHelpDesc = `
Signs another CA's self-issued certificate. This is most often used for rolling roots; unless you know you need this you probably want to use sign-intermediate instead.

Note that this is a very privileged operation and should be extremely restricted in terms of who is allowed to use it. All values will be taken directly from the incoming certificate and only verification that it is self-issued will be performed.

Configured URLs for CRLs/OCSP/etc. will be copied over and the issuer will be this mount's CA cert. Other than that, all other values will be used verbatim.
`

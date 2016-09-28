package pki

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssue,
		},

		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})
	return ret
}

func pathSign(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSign,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	return ret
}

func pathSignVerbatim(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "sign-verbatim/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSignVerbatim,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "",
		Description: `PEM-format CSR to be signed. Values will be
taken verbatim from the CSR, except for
basic constraints.`,
	}

	return ret
}

// pathIssue issues a certificate and private key from given parameters,
// subject to role restrictions
func (b *backend) pathIssue(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	return b.pathIssueSignCert(req, data, role, false, false)
}

// pathSign issues a certificate from a submitted CSR, subject to role
// restrictions
func (b *backend) pathSign(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	return b.pathIssueSignCert(req, data, role, true, false)
}

// pathSignVerbatim issues a certificate from a submitted CSR, *not* subject to
// role restrictions
func (b *backend) pathSignVerbatim(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	ttl := b.System().DefaultLeaseTTL()
	role := &roleEntry{
		TTL:              ttl.String(),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		AllowIPSANs:      true,
		EnforceHostnames: false,
		KeyType:          "any",
	}

	return b.pathIssueSignCert(req, data, role, true, true)
}

func (b *backend) pathIssueSignCert(
	req *logical.Request, data *framework.FieldData, role *roleEntry, useCSR, useCSRValues bool) (*logical.Response, error) {
	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem", "der", or "pem_bundle"`), nil
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case errutil.UserError:
		return nil, errutil.UserError{Err: fmt.Sprintf(
			"Could not fetch the CA certificate (was one set?): %s", caErr)}
	case errutil.InternalError:
		return nil, errutil.InternalError{Err: fmt.Sprintf(
			"Error fetching CA certificate: %s", caErr)}
	}

	var parsedBundle *certutil.ParsedCertBundle
	var err error
	if useCSR {
		parsedBundle, err = signCert(b, role, signingBundle, false, useCSRValues, req, data)
	} else {
		parsedBundle, err = generateCert(b, role, signingBundle, false, req, data)
	}
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw signing bundle to cert bundle: %s", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw cert bundle to cert bundle: %s", err)
	}

	resp := b.Secret(SecretCertsType).Response(
		map[string]interface{}{
			"serial_number": cb.SerialNumber,
		},
		map[string]interface{}{
			"serial_number": cb.SerialNumber,
		})

	switch format {
	case "pem":
		resp.Data["issuing_ca"] = signingCB.Certificate
		resp.Data["certificate"] = cb.Certificate
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			resp.Data["ca_chain"] = cb.CAChain
		}
		if !useCSR {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}

	case "pem_bundle":
		resp.Data["issuing_ca"] = signingCB.Certificate
		resp.Data["certificate"] = cb.ToPEMBundle()
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			resp.Data["ca_chain"] = cb.CAChain
		}
		if !useCSR {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}

	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(signingBundle.CertificateBytes)

		var caChain []string
		for _, caCert := range parsedBundle.CAChain {
			caChain = append(caChain, base64.StdEncoding.EncodeToString(caCert.Bytes))
		}
		if caChain != nil && len(caChain) > 0 {
			resp.Data["ca_chain"] = caChain
		}

		if !useCSR {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
		}
	}

	resp.Secret.TTL = parsedBundle.Certificate.NotAfter.Sub(time.Now())

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "certs/" + cb.SerialNumber,
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to store certificate locally")
	}

	return resp, nil
}

const pathIssueHelpSyn = `
Request a certificate using a certain role with the provided details.
`

const pathIssueHelpDesc = `
This path allows requesting a certificate to be issued according to the
policy of the given role. The certificate will only be issued if the
requested details are allowed by the role policy.

This path returns a certificate and a private key. If you want a workflow
that does not expose a private key, generate a CSR locally and use the
sign path instead.
`

const pathSignHelpSyn = `
Request certificates using a certain role with the provided details.
`

const pathSignHelpDesc = `
This path allows requesting certificates to be issued according to the
policy of the given role. The certificate will only be issued if the
requested common name is allowed by the role policy.

This path requires a CSR; if you want Vault to generate a private key
for you, use the issue path instead.
`

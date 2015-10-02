package pki

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var issueAndSignSchema = map[string]*framework.FieldSchema{
	"role": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The desired role with configuration for this
request`,
	},

	"common_name": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested common name; if you want more than
one, specify the alternative names in the
alt_names map. If email protection is enabled
in the role, this may be an email address.`,
	},

	"alt_names": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list. If email protection
is enabled for the role, this may contain
email addresses.`,
	},

	"ip_sans": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested IP SANs, if any, in a
comma-delimited list`,
	},

	"ttl": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
	},
}

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),
		Fields:  issueAndSignSchema,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathIssue,
		},

		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}
}

func pathSign(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("role"),
		Fields:  issueAndSignSchema,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathSign,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `PEM-format CSR to be signed.`,
	}

	return ret
}

func (b *backend) pathIssue(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathIssueSignCert(req, data, false)
}

func (b *backend) pathSign(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathIssueSignCert(req, data, true)
}

func (b *backend) pathIssueSignCert(
	req *logical.Request, data *framework.FieldData, useCSR bool) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case certutil.UserError:
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"Could not fetch the CA certificate (was one set?): %s", caErr)}
	case certutil.InternalError:
		return nil, certutil.InternalError{Err: fmt.Sprintf(
			"Error fetching CA certificate: %s", caErr)}
	}

	var parsedBundle *certutil.ParsedCertBundle
	if useCSR {
		parsedBundle, err = signCert(b, role, signingBundle, false, req, data)
	} else {
		parsedBundle, err = generateCert(b, role, signingBundle, false, req, data)
	}
	if err != nil {
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case certutil.InternalError:
			return nil, err
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw cert bundle to cert bundle: %s", err)
	}

	resp := b.Secret(SecretCertsType).Response(
		structs.New(cb).Map(),
		map[string]interface{}{
			"serial_number": cb.SerialNumber,
		})

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

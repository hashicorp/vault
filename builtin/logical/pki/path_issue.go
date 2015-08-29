package pki

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The desired role with configuration for this
request`,
			},

			"common_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested common name; if you want more than
one, specify the alternative names in the
alt_names map`,
			},

			"alt_names": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list`,
			},

			"ip_sans": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested IP SANs, if any, in a
common-delimited list`,
			},

			"ttl": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathIssueCert,
		},

		HelpSynopsis:    pathIssueCertHelpSyn,
		HelpDescription: pathIssueCertHelpDesc,
	}
}

func (b *backend) pathIssueCert(
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

	// Don't allow these on the standard path. Ideally we should determine
	// this internally once we get SudoPrivilege from System() working
	// for non-TokenStore
	delete(req.Data, "ca_type")
	delete(req.Data, "pki_address")

	parsedBundle, err := generateCert(b, role, signingBundle, req, data)
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

const pathIssueCertHelpSyn = `
Request certificates using a certain role with the provided common name.
`

const pathIssueCertHelpDesc = `
This path allows requesting certificates to be issued according to the
policy of the given role. The certificate will only be issued if the
requested common name is allowed by the role policy.
`

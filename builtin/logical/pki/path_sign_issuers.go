package pki

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssuerSignIntermediate(b *backend) *framework.Path {
	pattern := "issuers/" + framework.GenericNameRegex("ref") + "/sign-intermediate"
	return pathIssuerSignIntermediateRaw(b, pattern)
}

func pathSignIntermediate(b *backend) *framework.Path {
	pattern := "root/sign-intermediate"
	return pathIssuerSignIntermediateRaw(b, pattern)
}

func pathIssuerSignIntermediateRaw(b *backend, pattern string) *framework.Path {
	path := &framework.Path{
		Pattern: pattern,
		Fields: map[string]*framework.FieldSchema{
			"ref": {
				Type:        framework.TypeString,
				Description: `Reference to issuer; either "default" for the configured default issuer, an identifier of an issuer, or the name assigned to the issuer.`,
				Default:     "default",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssuerSignIntermediate,
		},

		HelpSynopsis:    pathIssuerSignIntermediateHelpSyn,
		HelpDescription: pathIssuerSignIntermediateHelpDesc,
	}

	path.Fields = addCACommonFields(path.Fields)
	path.Fields = addCAIssueFields(path.Fields)

	path.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	path.Fields["use_csr_values"] = &framework.FieldSchema{
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
the non-repudiation flag;
3) Extensions requested in the CSR will be copied
into the issued certificate.`,
	}

	return path
}

const (
	pathIssuerSignIntermediateHelpSyn  = `Issue an intermediate CA certificate based on the provided CSR.`
	pathIssuerSignIntermediateHelpDesc = `
This API endpoint allows for signing the specified CSR, adding to it a basic
constraint for IsCA=True. This allows the issued certificate to issue its own
leaf certificates.

Note that the resulting certificate is not imported as an issuer in this PKI
mount. This means that you can use the resulting certificate in another Vault
PKI mount point or to issue an external intermediate (e.g., for use with
another X.509 CA).

See the API documentation for more information about required parameters.
`
)

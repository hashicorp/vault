package cert

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the certificate",
			},

			"certificate": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The public certificate that should be trusted.
Must be x509 PEM encoded.`,
			},

			"display_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The display name to use for clients using this
certificate.`,
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-seperated list of policies.",
			},

			"lease": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `Deprecated: use "ttl" instead. TTL time in
seconds. Defaults to system/backend default TTL.`,
			},

			"ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `TTL for tokens issued by this backend.
Defaults to system/backend default TTL time.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCertDelete,
			logical.ReadOperation:   b.pathCertRead,
			logical.UpdateOperation:  b.pathCertWrite,
		},

		HelpSynopsis:    pathCertHelpSyn,
		HelpDescription: pathCertHelpDesc,
	}
}

func (b *backend) Cert(s logical.Storage, n string) (*CertEntry, error) {
	entry, err := s.Get("cert/" + strings.ToLower(n))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result CertEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathCertDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("cert/" + strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathCertRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cert, err := b.Cert(req.Storage, strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, nil
	}

	duration := cert.TTL
	if duration == 0 {
		duration = b.System().DefaultLeaseTTL()
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"certificate":  cert.Certificate,
			"display_name": cert.DisplayName,
			"policies":     strings.Join(cert.Policies, ","),
			"ttl":          duration / time.Second,
		},
	}, nil
}

func (b *backend) pathCertWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	certificate := d.Get("certificate").(string)
	displayName := d.Get("display_name").(string)
	policies := strings.Split(d.Get("policies").(string), ",")
	for i, p := range policies {
		policies[i] = strings.TrimSpace(p)
	}

	// Default the display name to the certificate name if not given
	if displayName == "" {
		displayName = name
	}

	if len(policies) == 0 {
		return logical.ErrorResponse("policies required"), nil
	}
	parsed := parsePEM([]byte(certificate))
	if len(parsed) == 0 {
		return logical.ErrorResponse("failed to parse certificate"), nil
	}

	certEntry := &CertEntry{
		Name:        name,
		Certificate: certificate,
		DisplayName: displayName,
		Policies:    policies,
	}

	// Parse the lease duration or default to backend/system default
	var err error
	maxTTL := b.System().MaxLeaseTTL()
	ttl := time.Duration(d.Get("ttl").(int)) * time.Second
	if ttl == time.Duration(0) {
		ttl = time.Second * time.Duration(d.Get("lease").(int))
	}
	if ttl > maxTTL {
		return logical.ErrorResponse(fmt.Sprintf("Given TTL of %d seconds greater than current mount/system default of %d seconds", ttl/time.Second, maxTTL/time.Second)), nil
	}
	if ttl > time.Duration(0) {
		certEntry.TTL = ttl
	}

	// Store it
	entry, err := logical.StorageEntryJSON("cert/"+name, certEntry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

type CertEntry struct {
	Name        string
	Certificate string
	DisplayName string
	Policies    []string
	TTL         time.Duration
}

const pathCertHelpSyn = `
Manage trusted certificates used for authentication.
`

const pathCertHelpDesc = `
This endpoint allows you to create, read, update, and delete trusted certificates
that are allowed to authenticate.

Deleting a certificate will not revoke auth for prior authenticated connections.
To do this, do a revoke on "login". If you don't need to revoke login immediately,
then the next renew will cause the lease to expire.
`

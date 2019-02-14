package cert

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathCertList,
		},

		HelpSynopsis:    pathCertHelpSyn,
		HelpDescription: pathCertHelpDesc,
	}
}

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

			"allowed_names": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of names.
At least one must exist in either the Common Name or SANs. Supports globbing.  
This parameter is deprecated, please use allowed_common_names, allowed_dns_sans, 
allowed_email_sans, allowed_uri_sans.`,
			},

			"allowed_common_names": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of names.
At least one must exist in the Common Name. Supports globbing.`,
			},

			"allowed_dns_sans": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of DNS names.
At least one must exist in the SANs. Supports globbing.`,
			},

			"allowed_email_sans": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of Email Addresses.
At least one must exist in the SANs. Supports globbing.`,
			},

			"allowed_uri_sans": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of URIs.
At least one must exist in the SANs. Supports globbing.`,
			},

			"allowed_organizational_units": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of Organizational Units names.
At least one must exist in the OU field.`,
			},

			"required_extensions": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated string or array of extensions
formatted as "oid:value". Expects the extension value to be some type of ASN1 encoded string.
All values much match. Supports globbing on "value".`,
			},

			"display_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The display name to use for clients using this
certificate.`,
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies.",
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

			"max_ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `Duration in either an integer number of seconds (3600) or
an integer time unit (60m) after which the
issued token can no longer be renewed.`,
			},

			"period": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `If set, indicates that the token generated using this role
should never expire. The token should be renewed within the
duration specified by this value. At each renewal, the token's
TTL will be set to the value of this parameter.`,
			},
			"bound_cidrs": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can perform the login operation.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathCertDelete,
			logical.ReadOperation:   b.pathCertRead,
			logical.UpdateOperation: b.pathCertWrite,
		},

		HelpSynopsis:    pathCertHelpSyn,
		HelpDescription: pathCertHelpDesc,
	}
}

func (b *backend) Cert(ctx context.Context, s logical.Storage, n string) (*CertEntry, error) {
	entry, err := s.Get(ctx, "cert/"+strings.ToLower(n))
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

func (b *backend) pathCertDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "cert/"+strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathCertList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	certs, err := req.Storage.List(ctx, "cert/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(certs), nil
}

func (b *backend) pathCertRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cert, err := b.Cert(ctx, req.Storage, strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"certificate":                  cert.Certificate,
			"display_name":                 cert.DisplayName,
			"policies":                     cert.Policies,
			"ttl":                          cert.TTL / time.Second,
			"max_ttl":                      cert.MaxTTL / time.Second,
			"period":                       cert.Period / time.Second,
			"allowed_names":                cert.AllowedNames,
			"allowed_common_names":         cert.AllowedCommonNames,
			"allowed_dns_sans":             cert.AllowedDNSSANs,
			"allowed_email_sans":           cert.AllowedEmailSANs,
			"allowed_uri_sans":             cert.AllowedURISANs,
			"allowed_organizational_units": cert.AllowedOrganizationalUnits,
			"required_extensions":          cert.RequiredExtensions,
			"bound_cidrs":                  cert.BoundCIDRs,
		},
	}, nil
}

func (b *backend) pathCertWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(d.Get("name").(string))
	certificate := d.Get("certificate").(string)
	displayName := d.Get("display_name").(string)
	policies := policyutil.ParsePolicies(d.Get("policies"))
	allowedNames := d.Get("allowed_names").([]string)
	allowedCommonNames := d.Get("allowed_common_names").([]string)
	allowedDNSSANs := d.Get("allowed_dns_sans").([]string)
	allowedEmailSANs := d.Get("allowed_email_sans").([]string)
	allowedURISANs := d.Get("allowed_uri_sans").([]string)
	allowedOrganizationalUnits := d.Get("allowed_organizational_units").([]string)
	requiredExtensions := d.Get("required_extensions").([]string)

	var resp logical.Response

	// Parse the ttl (or lease duration)
	systemDefaultTTL := b.System().DefaultLeaseTTL()
	ttl := time.Duration(d.Get("ttl").(int)) * time.Second
	if ttl == 0 {
		ttl = time.Duration(d.Get("lease").(int)) * time.Second
	}
	if ttl > systemDefaultTTL {
		resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds is greater than current mount/system default of %d seconds", ttl/time.Second, systemDefaultTTL/time.Second))
	}

	if ttl < time.Duration(0) {
		return logical.ErrorResponse("ttl cannot be negative"), nil
	}

	// Parse max_ttl
	systemMaxTTL := b.System().MaxLeaseTTL()
	maxTTL := time.Duration(d.Get("max_ttl").(int)) * time.Second
	if maxTTL > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds is greater than current mount/system default of %d seconds", maxTTL/time.Second, systemMaxTTL/time.Second))
	}

	if maxTTL < time.Duration(0) {
		return logical.ErrorResponse("max_ttl cannot be negative"), nil
	}

	if maxTTL != 0 && ttl > maxTTL {
		return logical.ErrorResponse("ttl should be shorter than max_ttl"), nil
	}

	// Parse period
	period := time.Duration(d.Get("period").(int)) * time.Second
	if period > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given period of %d seconds is greater than the backend's maximum TTL of %d seconds", period/time.Second, systemMaxTTL/time.Second))
	}

	if period < time.Duration(0) {
		return logical.ErrorResponse("period cannot be negative"), nil
	}

	// Default the display name to the certificate name if not given
	if displayName == "" {
		displayName = name
	}

	parsed := parsePEM([]byte(certificate))
	if len(parsed) == 0 {
		return logical.ErrorResponse("failed to parse certificate"), nil
	}

	// If the certificate is not a CA cert, then ensure that x509.ExtKeyUsageClientAuth is set
	if !parsed[0].IsCA && parsed[0].ExtKeyUsage != nil {
		var clientAuth bool
		for _, usage := range parsed[0].ExtKeyUsage {
			if usage == x509.ExtKeyUsageClientAuth || usage == x509.ExtKeyUsageAny {
				clientAuth = true
				break
			}
		}
		if !clientAuth {
			return logical.ErrorResponse("non-CA certificates should have TLS client authentication set as an extended key usage"), nil
		}
	}

	parsedCIDRs, err := parseutil.ParseAddrs(d.Get("bound_cidrs"))
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	certEntry := &CertEntry{
		Name:                       name,
		Certificate:                certificate,
		DisplayName:                displayName,
		Policies:                   policies,
		AllowedNames:               allowedNames,
		AllowedCommonNames:         allowedCommonNames,
		AllowedDNSSANs:             allowedDNSSANs,
		AllowedEmailSANs:           allowedEmailSANs,
		AllowedURISANs:             allowedURISANs,
		AllowedOrganizationalUnits: allowedOrganizationalUnits,
		RequiredExtensions:         requiredExtensions,
		TTL:                        ttl,
		MaxTTL:                     maxTTL,
		Period:                     period,
		BoundCIDRs:                 parsedCIDRs,
	}

	// Store it
	entry, err := logical.StorageEntryJSON("cert/"+name, certEntry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}

type CertEntry struct {
	Name                       string
	Certificate                string
	DisplayName                string
	Policies                   []string
	TTL                        time.Duration
	MaxTTL                     time.Duration
	Period                     time.Duration
	AllowedNames               []string
	AllowedCommonNames         []string
	AllowedDNSSANs             []string
	AllowedEmailSANs           []string
	AllowedURISANs             []string
	AllowedOrganizationalUnits []string
	RequiredExtensions         []string
	BoundCIDRs                 []*sockaddr.SockAddrMarshaler
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

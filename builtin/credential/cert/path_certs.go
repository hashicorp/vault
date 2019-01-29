package cert

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/tokenhelper"
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
	path := &framework.Path{
		Pattern: "certs/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeLowerCaseString,
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

			"lease": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `Deprecated: use "ttl" instead. TTL time in
seconds. Defaults to system/backend default TTL.`,
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
	tokenhelper.AddTokenFields(path.Fields)

	return path
}

func (b *backend) Cert(ctx context.Context, s logical.Storage, n string) (*CertEntry, error) {
	entry, err := s.Get(ctx, "cert/"+n)
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

	var needsUpgrade bool
	if result.OldTTL != 0 {
		needsUpgrade = true
		result.TTL = result.OldTTL
		result.OldTTL = 0
	}
	if result.OldMaxTTL != 0 {
		needsUpgrade = true
		result.MaxTTL = result.OldMaxTTL
		result.OldMaxTTL = 0
	}
	if result.OldPeriod != 0 {
		needsUpgrade = true
		result.Period = result.OldPeriod
		result.OldPeriod = 0
	}
	if len(result.OldPolicies) > 0 {
		needsUpgrade = true
		result.Policies = result.OldPolicies
		result.OldPolicies = nil
	}
	if len(result.OldBoundCIDRs) > 0 {
		needsUpgrade = true
		result.BoundCIDRs = result.OldBoundCIDRs
		result.OldBoundCIDRs = nil
	}
	if needsUpgrade && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		entry, err := logical.StorageEntryJSON("cert/"+n, result)
		if err != nil {
			return nil, err
		}
		if err := s.Put(ctx, entry); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (b *backend) pathCertDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "cert/"+d.Get("name").(string))
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
	cert, err := b.Cert(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"certificate":                  cert.Certificate,
		"display_name":                 cert.DisplayName,
		"allowed_names":                cert.AllowedNames,
		"allowed_common_names":         cert.AllowedCommonNames,
		"allowed_dns_sans":             cert.AllowedDNSSANs,
		"allowed_email_sans":           cert.AllowedEmailSANs,
		"allowed_uri_sans":             cert.AllowedURISANs,
		"allowed_organizational_units": cert.AllowedOrganizationalUnits,
		"required_extensions":          cert.RequiredExtensions,
	}
	cert.PopulateTokenData(data)
	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathCertWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	certificate := d.Get("certificate").(string)
	displayName := d.Get("display_name").(string)
	allowedNames := d.Get("allowed_names").([]string)
	allowedCommonNames := d.Get("allowed_common_names").([]string)
	allowedDNSSANs := d.Get("allowed_dns_sans").([]string)
	allowedEmailSANs := d.Get("allowed_email_sans").([]string)
	allowedURISANs := d.Get("allowed_uri_sans").([]string)
	allowedOrganizationalUnits := d.Get("allowed_organizational_units").([]string)
	requiredExtensions := d.Get("required_extensions").([]string)

	certEntry := &CertEntry{}
	var resp logical.Response

	if err := certEntry.ParseTokenFields(req, d); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Parse the ttl (or lease duration)
	systemDefaultTTL := b.System().DefaultLeaseTTL()
	if certEntry.TTL > systemDefaultTTL {
		resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds is greater than current mount/system default of %d seconds", certEntry.TTL/time.Second, systemDefaultTTL/time.Second))
	}

	// Parse max_ttl
	systemMaxTTL := b.System().MaxLeaseTTL()
	if certEntry.MaxTTL > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds is greater than current mount/system default of %d seconds", certEntry.MaxTTL/time.Second, systemMaxTTL/time.Second))
	}

	if certEntry.MaxTTL != 0 && certEntry.TTL > certEntry.MaxTTL {
		return logical.ErrorResponse("ttl should be shorter than max_ttl"), nil
	}

	// Parse period
	if certEntry.Period > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given period of %d seconds is greater than the backend's maximum TTL of %d seconds", certEntry.Period/time.Second, systemMaxTTL/time.Second))
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

	certEntry.Name = name
	certEntry.Certificate = certificate
	certEntry.DisplayName = displayName
	certEntry.AllowedNames = allowedNames
	certEntry.AllowedCommonNames = allowedCommonNames
	certEntry.AllowedDNSSANs = allowedDNSSANs
	certEntry.AllowedEmailSANs = allowedEmailSANs
	certEntry.AllowedURISANs = allowedURISANs
	certEntry.AllowedOrganizationalUnits = allowedOrganizationalUnits
	certEntry.RequiredExtensions = requiredExtensions

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
	tokenhelper.TokenParams

	Name                       string
	Certificate                string
	DisplayName                string
	AllowedNames               []string
	AllowedCommonNames         []string
	AllowedDNSSANs             []string
	AllowedEmailSANs           []string
	AllowedURISANs             []string
	AllowedOrganizationalUnits []string
	RequiredExtensions         []string

	// These token-related fields have been moved to the embedded tokenhelper.TokenParams struct
	OldPolicies   []string                      `json:"Policies"`
	OldTTL        time.Duration                 `json:"TTL"`
	OldMaxTTL     time.Duration                 `json:"MaxTTL"`
	OldPeriod     time.Duration                 `json:"Period"`
	OldBoundCIDRs []*sockaddr.SockAddrMarshaler `json:"BoundCIDRs"`
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListCerts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "certs/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCert,
			OperationSuffix: "certificates",
			Navigation:      true,
			ItemType:        "Certificate",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathCertList,
		},

		HelpSynopsis:    pathCertHelpSyn,
		HelpDescription: pathCertHelpDesc,
	}
}

func pathCerts(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "certs/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCert,
			OperationSuffix: "certificate",
			Action:          "Create",
			ItemType:        "Certificate",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the certificate",
			},

			"certificate": {
				Type: framework.TypeString,
				Description: `The public certificate that should be trusted.
Must be x509 PEM encoded.`,
				DisplayAttrs: &framework.DisplayAttributes{
					EditType: "file",
				},
			},
			"ocsp_enabled": {
				Type:        framework.TypeBool,
				Description: `Whether to attempt OCSP verification of certificates at login`,
			},
			"ocsp_ca_certificates": {
				Type:        framework.TypeString,
				Description: `Any additional CA certificates needed to communicate with OCSP servers`,
				DisplayAttrs: &framework.DisplayAttributes{
					EditType: "file",
				},
			},
			"ocsp_servers_override": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of OCSP server addresses.  If unset, the OCSP server is determined 
from the AuthorityInformationAccess extension on the certificate being inspected.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Description: "A list of OCSP server addresses. If unset, the OCSP server is determined from the AuthorityInformationAccess extension on the certificate being inspected.",
				},
			},
			"ocsp_fail_open": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to true, if an OCSP revocation cannot be made successfully, login will proceed rather than failing.  If false, failing to get an OCSP status fails the request.",
			},
			"ocsp_query_all_servers": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set to true, rather than accepting the first successful OCSP response, query all servers and consider the certificate valid only if all servers agree.",
			},
			"ocsp_this_update_max_age": {
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "If greater than 0, specifies the maximum age of an OCSP thisUpdate field to avoid accepting old responses without a nextUpdate field.",
			},
			"ocsp_max_retries": {
				Type:        framework.TypeInt,
				Default:     4,
				Description: "The number of retries the OCSP client should attempt per query.",
			},
			"allowed_names": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of names.
At least one must exist in either the Common Name or SANs. Supports globbing.  
This parameter is deprecated, please use allowed_common_names, allowed_dns_sans, 
allowed_email_sans, allowed_uri_sans.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Group:       "Constraints",
					Description: "A list of names. At least one must exist in either the Common Name or SANs. Supports globbing. This parameter is deprecated, please use allowed_common_names, allowed_dns_sans, allowed_email_sans, allowed_uri_sans.",
				},
			},

			"allowed_common_names": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of names.
At least one must exist in the Common Name. Supports globbing.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Group:       "Constraints",
					Description: "A list of names. At least one must exist in the Common Name. Supports globbing.",
				},
			},

			"allowed_dns_sans": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of DNS names.
At least one must exist in the SANs. Supports globbing.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:        "Allowed DNS SANs",
					Group:       "Constraints",
					Description: "A list of DNS names. At least one must exist in the SANs. Supports globbing.",
				},
			},

			"allowed_email_sans": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of Email Addresses.
At least one must exist in the SANs. Supports globbing.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:        "Allowed Email SANs",
					Group:       "Constraints",
					Description: "A list of Email Addresses. At least one must exist in the SANs. Supports globbing.",
				},
			},

			"allowed_uri_sans": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of URIs.
At least one must exist in the SANs. Supports globbing.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:        "Allowed URI SANs",
					Group:       "Constraints",
					Description: "A list of URIs. At least one must exist in the SANs. Supports globbing.",
				},
			},

			"allowed_organizational_units": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated list of Organizational Units names.
At least one must exist in the OU field.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Group:       "Constraints",
					Description: "A list of Organizational Units names. At least one must exist in the OU field.",
				},
			},

			"required_extensions": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated string or array of extensions
formatted as "oid:value". Expects the extension value to be some type of ASN1 encoded string.
All values much match. Supports globbing on "value".`,
				DisplayAttrs: &framework.DisplayAttributes{
					Description: "A list of extensions formatted as 'oid:value'. Expects the extension value to be some type of ASN1 encoded string. All values much match. Supports globbing on 'value'.",
				},
			},

			"allowed_metadata_extensions": {
				Type: framework.TypeCommaStringSlice,
				Description: `A comma-separated string or array of oid extensions.
Upon successful authentication, these extensions will be added as metadata if they are present
in the certificate. The metadata key will be the string consisting of the oid numbers
separated by a dash (-) instead of a dot (.) to allow usage in ACL templates.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Description: "A list of OID extensions. Upon successful authentication, these extensions will be added as metadata if they are present in the certificate. The metadata key will be the string consisting of the OID numbers separated by a dash (-) instead of a dot (.) to allow usage in ACL templates.",
				},
			},

			"display_name": {
				Type: framework.TypeString,
				Description: `The display name to use for clients using this
certificate.`,
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},

			"lease": {
				Type:        framework.TypeInt,
				Description: tokenutil.DeprecationText("token_ttl"),
				Deprecated:  true,
			},

			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_ttl"),
				Deprecated:  true,
			},

			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_max_ttl"),
				Deprecated:  true,
			},

			"period": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_period"),
				Deprecated:  true,
			},

			"bound_cidrs": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_bound_cidrs"),
				Deprecated:  true,
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

	tokenutil.AddTokenFields(p.Fields)
	return p
}

func (b *backend) Cert(ctx context.Context, s logical.Storage, n string) (*CertEntry, error) {
	entry, err := s.Get(ctx, trustedCertPath+strings.ToLower(n))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := CertEntry{OcspMaxRetries: defaultOcspMaxRetries} // Specify our defaults if the key is missing
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.TokenTTL == 0 && result.TTL > 0 {
		result.TokenTTL = result.TTL
	}
	if result.TokenMaxTTL == 0 && result.MaxTTL > 0 {
		result.TokenMaxTTL = result.MaxTTL
	}
	if result.TokenPeriod == 0 && result.Period > 0 {
		result.TokenPeriod = result.Period
	}
	if len(result.TokenPolicies) == 0 && len(result.Policies) > 0 {
		result.TokenPolicies = result.Policies
	}
	if len(result.TokenBoundCIDRs) == 0 && len(result.BoundCIDRs) > 0 {
		result.TokenBoundCIDRs = result.BoundCIDRs
	}

	return &result, nil
}

func (b *backend) pathCertDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	defer b.flushTrustedCache()
	err := req.Storage.Delete(ctx, trustedCertPath+strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathCertList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	certs, err := req.Storage.List(ctx, trustedCertPath)
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
		"allowed_metadata_extensions":  cert.AllowedMetadataExtensions,
		"ocsp_ca_certificates":         cert.OcspCaCertificates,
		"ocsp_enabled":                 cert.OcspEnabled,
		"ocsp_servers_override":        cert.OcspServersOverride,
		"ocsp_fail_open":               cert.OcspFailOpen,
		"ocsp_query_all_servers":       cert.OcspQueryAllServers,
		"ocsp_this_update_max_age":     int64(cert.OcspThisUpdateMaxAge.Seconds()),
		"ocsp_max_retries":             cert.OcspMaxRetries,
	}
	cert.PopulateTokenData(data)

	if cert.TTL > 0 {
		data["ttl"] = int64(cert.TTL.Seconds())
	}
	if cert.MaxTTL > 0 {
		data["max_ttl"] = int64(cert.MaxTTL.Seconds())
	}
	if cert.Period > 0 {
		data["period"] = int64(cert.Period.Seconds())
	}
	if len(cert.Policies) > 0 {
		data["policies"] = data["token_policies"]
	}
	if len(cert.BoundCIDRs) > 0 {
		data["bound_cidrs"] = data["token_bound_cidrs"]
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathCertWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	defer b.flushTrustedCache()
	name := strings.ToLower(d.Get("name").(string))

	cert, err := b.Cert(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if cert == nil {
		cert = &CertEntry{
			Name:           name,
			OcspMaxRetries: defaultOcspMaxRetries,
		}
	}

	// Get non tokenutil fields
	if certificateRaw, ok := d.GetOk("certificate"); ok {
		cert.Certificate = certificateRaw.(string)
	}
	if ocspCertificatesRaw, ok := d.GetOk("ocsp_ca_certificates"); ok {
		cert.OcspCaCertificates = ocspCertificatesRaw.(string)
	}
	if ocspEnabledRaw, ok := d.GetOk("ocsp_enabled"); ok {
		cert.OcspEnabled = ocspEnabledRaw.(bool)
	}
	if ocspServerOverrides, ok := d.GetOk("ocsp_servers_override"); ok {
		cert.OcspServersOverride = ocspServerOverrides.([]string)
	}
	if ocspFailOpen, ok := d.GetOk("ocsp_fail_open"); ok {
		cert.OcspFailOpen = ocspFailOpen.(bool)
	}
	if ocspQueryAll, ok := d.GetOk("ocsp_query_all_servers"); ok {
		cert.OcspQueryAllServers = ocspQueryAll.(bool)
	}
	if ocspThisUpdateMaxAge, ok := d.GetOk("ocsp_this_update_max_age"); ok {
		maxAgeDuration, err := parseutil.ParseDurationSecond(ocspThisUpdateMaxAge)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ocsp_this_update_max_age: %w", err)
		}
		cert.OcspThisUpdateMaxAge = maxAgeDuration
	}
	if ocspMaxRetries, ok := d.GetOk("ocsp_max_retries"); ok {
		cert.OcspMaxRetries = ocspMaxRetries.(int)
		if cert.OcspMaxRetries < 0 {
			return nil, fmt.Errorf("ocsp_max_retries can not be a negative number")
		}
	}
	if displayNameRaw, ok := d.GetOk("display_name"); ok {
		cert.DisplayName = displayNameRaw.(string)
	}
	if allowedNamesRaw, ok := d.GetOk("allowed_names"); ok {
		cert.AllowedNames = allowedNamesRaw.([]string)
	}
	if allowedCommonNamesRaw, ok := d.GetOk("allowed_common_names"); ok {
		cert.AllowedCommonNames = allowedCommonNamesRaw.([]string)
	}
	if allowedDNSSANsRaw, ok := d.GetOk("allowed_dns_sans"); ok {
		cert.AllowedDNSSANs = allowedDNSSANsRaw.([]string)
	}
	if allowedEmailSANsRaw, ok := d.GetOk("allowed_email_sans"); ok {
		cert.AllowedEmailSANs = allowedEmailSANsRaw.([]string)
	}
	if allowedURISANsRaw, ok := d.GetOk("allowed_uri_sans"); ok {
		cert.AllowedURISANs = allowedURISANsRaw.([]string)
	}
	if allowedOrganizationalUnitsRaw, ok := d.GetOk("allowed_organizational_units"); ok {
		cert.AllowedOrganizationalUnits = allowedOrganizationalUnitsRaw.([]string)
	}
	if requiredExtensionsRaw, ok := d.GetOk("required_extensions"); ok {
		cert.RequiredExtensions = requiredExtensionsRaw.([]string)
	}
	if allowedMetadataExtensionsRaw, ok := d.GetOk("allowed_metadata_extensions"); ok {
		cert.AllowedMetadataExtensions = allowedMetadataExtensionsRaw.([]string)
	}

	// Get tokenutil fields
	if err := cert.ParseTokenFields(req, d); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(d, "policies", "token_policies", &cert.Policies, &cert.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(d, "ttl", "token_ttl", &cert.TTL, &cert.TokenTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		// Special case here for old lease value
		_, ok := d.GetOk("token_ttl")
		if !ok {
			_, ok = d.GetOk("ttl")
			if !ok {
				ttlRaw, ok := d.GetOk("lease")
				if ok {
					cert.TTL = time.Duration(ttlRaw.(int)) * time.Second
					cert.TokenTTL = cert.TTL
				}
			}
		}

		if err := tokenutil.UpgradeValue(d, "max_ttl", "token_max_ttl", &cert.MaxTTL, &cert.TokenMaxTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(d, "period", "token_period", &cert.Period, &cert.TokenPeriod); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(d, "bound_cidrs", "token_bound_cidrs", &cert.BoundCIDRs, &cert.TokenBoundCIDRs); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	var resp logical.Response

	systemDefaultTTL := b.System().DefaultLeaseTTL()
	if cert.TokenTTL > systemDefaultTTL {
		resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds is greater than current mount/system default of %d seconds", cert.TokenTTL/time.Second, systemDefaultTTL/time.Second))
	}
	systemMaxTTL := b.System().MaxLeaseTTL()
	if cert.TokenMaxTTL > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds is greater than current mount/system default of %d seconds", cert.TokenMaxTTL/time.Second, systemMaxTTL/time.Second))
	}
	if cert.TokenMaxTTL != 0 && cert.TokenTTL > cert.TokenMaxTTL {
		return logical.ErrorResponse("ttl should be shorter than max_ttl"), nil
	}
	if cert.TokenPeriod > systemMaxTTL {
		resp.AddWarning(fmt.Sprintf("Given period of %d seconds is greater than the backend's maximum TTL of %d seconds", cert.TokenPeriod/time.Second, systemMaxTTL/time.Second))
	}

	// Default the display name to the certificate name if not given
	if cert.DisplayName == "" {
		cert.DisplayName = name
	}

	parsed := parsePEM([]byte(cert.Certificate))
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
			return logical.ErrorResponse("nonCA certificates should have TLS client authentication set as an extended key usage"), nil
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON(trustedCertPath+name, cert)
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
	tokenutil.TokenParams

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
	AllowedMetadataExtensions  []string
	BoundCIDRs                 []*sockaddr.SockAddrMarshaler

	OcspCaCertificates   string
	OcspEnabled          bool
	OcspServersOverride  []string
	OcspFailOpen         bool
	OcspQueryAllServers  bool
	OcspThisUpdateMaxAge time.Duration
	OcspMaxRetries       int
}

const pathCertHelpSyn = `
Manage trusted certificates used for authentication.
`

const pathCertHelpDesc = `
This endpoint allows you to create, read, update, and delete trusted certificates
that are allowed to authenticate.

Deleting a certificate will not revoke auth for prior authenticated connections.
To do this, do a revoke on "login". If you don'log need to revoke login immediately,
then the next renew will cause the lease to expire.

`

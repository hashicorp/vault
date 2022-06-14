package pki

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigURLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/urls",
		Fields: map[string]*framework.FieldSchema{
			"issuing_certificates": {
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the issuing certificate attribute. See also RFC 5280 Section 4.2.2.1.`,
			},

			"crl_distribution_points": {
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.`,
			},

			"ocsp_servers": {
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the OCSP servers attribute. See also RFC 5280 Section 4.2.2.1.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathWriteURL,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadURL,
			},
		},

		HelpSynopsis:    pathConfigURLsHelpSyn,
		HelpDescription: pathConfigURLsHelpDesc,
	}
}

func validateURLs(urls []string) string {
	for _, curr := range urls {
		if !govalidator.IsURL(curr) {
			return curr
		}
	}

	return ""
}

func getURLs(ctx context.Context, req *logical.Request) (*certutil.URLEntries, error) {
	entry, err := req.Storage.Get(ctx, "urls")
	if err != nil {
		return nil, err
	}

	entries := &certutil.URLEntries{
		IssuingCertificates:   []string{},
		CRLDistributionPoints: []string{},
		OCSPServers:           []string{},
	}

	if entry == nil {
		return entries, nil
	}

	if err := entry.DecodeJSON(entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func writeURLs(ctx context.Context, req *logical.Request, entries *certutil.URLEntries) error {
	entry, err := logical.StorageEntryJSON("urls", entries)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("unable to marshal entry into JSON")
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}

func (b *backend) pathReadURL(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entries, err := getURLs(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: structs.New(entries).Map(),
	}

	return resp, nil
}

func (b *backend) pathWriteURL(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := getURLs(ctx, req)
	if err != nil {
		return nil, err
	}

	if urlsInt, ok := data.GetOk("issuing_certificates"); ok {
		entries.IssuingCertificates = urlsInt.([]string)
		if badURL := validateURLs(entries.IssuingCertificates); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in issuing certificates: %s", badURL)), nil
		}
	}
	if urlsInt, ok := data.GetOk("crl_distribution_points"); ok {
		entries.CRLDistributionPoints = urlsInt.([]string)
		if badURL := validateURLs(entries.CRLDistributionPoints); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in CRL distribution points: %s", badURL)), nil
		}
	}
	if urlsInt, ok := data.GetOk("ocsp_servers"); ok {
		entries.OCSPServers = urlsInt.([]string)
		if badURL := validateURLs(entries.OCSPServers); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in OCSP servers: %s", badURL)), nil
		}
	}

	return nil, writeURLs(ctx, req, entries)
}

const pathConfigURLsHelpSyn = `
Set the URLs for the issuing CA, CRL distribution points, and OCSP servers.
`

const pathConfigURLsHelpDesc = `
This path allows you to set the issuing CA, CRL distribution points, and
OCSP server URLs that will be encoded into issued certificates. If these
values are not set, no such information will be encoded in the issued
certificates. To delete URLs, simply re-set the appropriate value with an
empty string.

Multiple URLs can be specified for each type; use commas to separate them.
`

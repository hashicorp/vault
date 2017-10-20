package pki

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigURLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/urls",
		Fields: map[string]*framework.FieldSchema{
			"issuing_certificates": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the issuing certificate attribute`,
			},

			"crl_distribution_points": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the CRL distribution points attribute`,
			},

			"ocsp_servers": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of URLs to be used
for the OCSP servers attribute`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathWriteURL,
			logical.ReadOperation:   b.pathReadURL,
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

func getURLs(req *logical.Request) (*urlEntries, error) {
	entry, err := req.Storage.Get("urls")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var entries urlEntries
	if err := entry.DecodeJSON(&entries); err != nil {
		return nil, err
	}

	return &entries, nil
}

func writeURLs(req *logical.Request, entries *urlEntries) error {
	entry, err := logical.StorageEntryJSON("urls", entries)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("Unable to marshal entry into JSON")
	}

	err = req.Storage.Put(entry)
	if err != nil {
		return err
	}

	return nil
}

func (b *backend) pathReadURL(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := getURLs(req)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(entries).Map(),
	}

	return resp, nil
}

func (b *backend) pathWriteURL(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := getURLs(req)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = &urlEntries{
			IssuingCertificates:   []string{},
			CRLDistributionPoints: []string{},
			OCSPServers:           []string{},
		}
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

	return nil, writeURLs(req, entries)
}

type urlEntries struct {
	IssuingCertificates   []string `json:"issuing_certificates" structs:"issuing_certificates" mapstructure:"issuing_certificates"`
	CRLDistributionPoints []string `json:"crl_distribution_points" structs:"crl_distribution_points" mapstructure:"crl_distribution_points"`
	OCSPServers           []string `json:"ocsp_servers" structs:"ocsp_servers" mapstructure:"ocsp_servers"`
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

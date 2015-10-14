package pki

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigURLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/urls/" + framework.GenericNameRegex("urltype"),
		Fields: map[string]*framework.FieldSchema{
			"urltype": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The type of URL to set. Can be "issuing"
to set the issuing certificate URL,
"crl" to set the CRL distribution
points, or "ocsp" to set the OCSP
servers. These values will be
recorded into issues certificates.
Only valid for write and delete
operations; read will return all
URL types.`,
			},
			"urls": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Comma-separated list of URLs to be used
for the type specified.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation:  b.pathWriteURL,
			logical.ReadOperation:   b.pathReadURL,
			logical.DeleteOperation: b.pathDeleteURL,
		},

		HelpSynopsis:    pathConfigURLsHelpSyn,
		HelpDescription: pathConfigURLsHelpDesc,
	}
}

func checkURLType(urlType string) *logical.Response {
	switch urlType {
	case "issuing":
		return nil
	case "crl":
		return nil
	case "ocsp":
		return nil
	default:
		return logical.ErrorResponse(fmt.Sprintf(
			"'%s' is not a valid type; must be 'issuing', 'crl', or 'ocsp'", urlType))
	}
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
	entry, err := logical.StorageEntryJSON("urls", &entries)
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

func (b *backend) pathDeleteURL(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	urlType := data.Get("urlType").(string)
	if err := checkURLType(urlType); err != nil {
		return err, nil
	}

	entries, err := getURLs(req)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return nil, nil
	}

	switch urlType {
	case "issuing":
		entries.IssuingCertificateURLs = []string{}
	case "crl":
		entries.CRLDistributionPoints = []string{}
	case "ocsp":
		entries.OCSPServers = []string{}
	}

	return nil, writeURLs(req, entries)
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
	urlType := data.Get("urlType").(string)
	if err := checkURLType(urlType); err != nil {
		return err, nil
	}

	urls := data.Get("urls").(string)
	if urls == "" {
		return logical.ErrorResponse("'urls' cannot be empty"), nil
	}

	splitURLs := strings.Split(urls, ",")

	entries, err := getURLs(req)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		entries = &urlEntries{
			IssuingCertificateURLs: []string{},
			CRLDistributionPoints:  []string{},
			OCSPServers:            []string{},
		}
	}

	switch urlType {
	case "issuing":
		entries.IssuingCertificateURLs = splitURLs
	case "crl":
		entries.CRLDistributionPoints = splitURLs
	case "ocsp":
		entries.OCSPServers = splitURLs
	}

	return nil, writeURLs(req, entries)
}

type urlEntries struct {
	IssuingCertificateURLs []string `json:"issuing_certificate_urls" structs:"issuing_certificate_urls" mapstructure:"issuing_certificate_urls"`
	CRLDistributionPoints  []string `json:"crl_distribution_points" structs:"crl_distribution_points" mapstructure:"crl_distribution_points"`
	OCSPServers            []string `json:"ocsp_servers" structs:"ocsp_servers" mapstructure:"ocsp_servers"`
}

const pathConfigURLsHelpSyn = `
Set the URLs for the issuing CA, CRL distribution points, and OCSP servers.
`

const pathConfigURLsHelpDesc = `
This path allows you to set the issuing CA, CRL distribution points, and
OCSP server URLs that will be encoded into issued certificates. If these
values are not set (or are set and then deleted), no such information will
be encoded in the issued certificates.

Multiple URLs can be specified for each type; use commas to separate them.
`

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigURLs(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/urls",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

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

			"enable_templating": {
				Type: framework.TypeBool,
				Description: `Whether or not to enabling templating of the
above AIA fields. When templating is enabled the special values '{{issuer_id}}',
'{{cluster_path}}', and '{{cluster_aia_path}}' are available, but the addresses
are not checked for URI validity until issuance time. Using '{{cluster_path}}'
requires /config/cluster's 'path' member to be set on all PR Secondary clusters
and using '{{cluster_aia_path}}' requires /config/cluster's 'aia_path' member
to be set on all PR secondary clusters.`,
				Default: false,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "urls",
				},
				Callback: b.pathWriteURL,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
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
							"enable_templating": {
								Type: framework.TypeBool,
								Description: `Whether or not to enabling templating of the
above AIA fields. When templating is enabled the special values '{{issuer_id}}'
and '{{cluster_path}}' are available, but the addresses are not checked for
URI validity until issuance time. This requires /config/cluster's path to be
set on all PR Secondary clusters.`,
								Default: false,
							},
						},
					}},
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathReadURL,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "urls-configuration",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"issuing_certificates": {
								Type: framework.TypeCommaStringSlice,
								Description: `Comma-separated list of URLs to be used
for the issuing certificate attribute. See also RFC 5280 Section 4.2.2.1.`,
								Required: true,
							},
							"crl_distribution_points": {
								Type: framework.TypeCommaStringSlice,
								Description: `Comma-separated list of URLs to be used
for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.`,
								Required: true,
							},
							"ocsp_servers": {
								Type: framework.TypeCommaStringSlice,
								Description: `Comma-separated list of URLs to be used
for the OCSP servers attribute. See also RFC 5280 Section 4.2.2.1.`,
								Required: true,
							},
							"enable_templating": {
								Type: framework.TypeBool,
								Description: `Whether or not to enable templating of the
above AIA fields. When templating is enabled the special values '{{issuer_id}}'
and '{{cluster_path}}' are available, but the addresses are not checked for
URI validity until issuance time. This requires /config/cluster's path to be
set on all PR Secondary clusters.`,
								Required: true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    pathConfigURLsHelpSyn,
		HelpDescription: pathConfigURLsHelpDesc,
	}
}

func getGlobalAIAURLs(ctx context.Context, storage logical.Storage) (*issuing.AiaConfigEntry, error) {
	entry, err := storage.Get(ctx, "urls")
	if err != nil {
		return nil, err
	}

	entries := &issuing.AiaConfigEntry{
		IssuingCertificates:   []string{},
		CRLDistributionPoints: []string{},
		OCSPServers:           []string{},
		EnableTemplating:      false,
	}

	if entry == nil {
		return entries, nil
	}

	if err := entry.DecodeJSON(entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func writeURLs(ctx context.Context, storage logical.Storage, entries *issuing.AiaConfigEntry) error {
	entry, err := logical.StorageEntryJSON("urls", entries)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("unable to marshal entry into JSON")
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}

func (b *backend) pathReadURL(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entries, err := getGlobalAIAURLs(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"issuing_certificates":    entries.IssuingCertificates,
			"crl_distribution_points": entries.CRLDistributionPoints,
			"ocsp_servers":            entries.OCSPServers,
			"enable_templating":       entries.EnableTemplating,
		},
	}

	return resp, nil
}

func (b *backend) pathWriteURL(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := getGlobalAIAURLs(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if enableTemplating, ok := data.GetOk("enable_templating"); ok {
		entries.EnableTemplating = enableTemplating.(bool)
	}
	if urlsInt, ok := data.GetOk("issuing_certificates"); ok {
		entries.IssuingCertificates = urlsInt.([]string)
	}
	if urlsInt, ok := data.GetOk("crl_distribution_points"); ok {
		entries.CRLDistributionPoints = urlsInt.([]string)
	}
	if urlsInt, ok := data.GetOk("ocsp_servers"); ok {
		entries.OCSPServers = urlsInt.([]string)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"issuing_certificates":    entries.IssuingCertificates,
			"crl_distribution_points": entries.CRLDistributionPoints,
			"ocsp_servers":            entries.OCSPServers,
			"enable_templating":       entries.EnableTemplating,
		},
	}

	if entries.EnableTemplating && !b.UseLegacyBundleCaStorage() {
		sc := b.makeStorageContext(ctx, req.Storage)
		issuers, err := sc.listIssuers()
		if err != nil {
			return nil, fmt.Errorf("unable to read issuers list to validate templated URIs: %w", err)
		}

		if len(issuers) > 0 {
			issuer, err := sc.fetchIssuerById(issuers[0])
			if err != nil {
				return nil, fmt.Errorf("unable to read issuer to validate templated URIs: %w", err)
			}

			_, err = ToURLEntries(sc, issuer.ID, entries)
			if err != nil {
				resp.AddWarning(fmt.Sprintf("issuance may fail: %v\n\nConsider setting the cluster-local address if it is not already set.", err))
			}
		}
	} else if !entries.EnableTemplating {
		if badURL := issuing.ValidateURLs(entries.IssuingCertificates); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in Authority Information Access (AIA) parameter issuing_certificates: %s", badURL)), nil
		}

		if badURL := issuing.ValidateURLs(entries.CRLDistributionPoints); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in Authority Information Access (AIA) parameter crl_distribution_points: %s", badURL)), nil
		}

		if badURL := issuing.ValidateURLs(entries.OCSPServers); badURL != "" {
			return logical.ErrorResponse(fmt.Sprintf(
				"invalid URL found in Authority Information Access (AIA) parameter ocsp_servers: %s", badURL)), nil
		}
	}

	if err := writeURLs(ctx, req.Storage, entries); err != nil {
		return nil, err
	}

	return resp, nil
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

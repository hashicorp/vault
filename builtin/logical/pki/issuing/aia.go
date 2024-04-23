// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const ClusterConfigPath = "config/cluster"

type AiaConfigEntry struct {
	IssuingCertificates   []string `json:"issuing_certificates"`
	CRLDistributionPoints []string `json:"crl_distribution_points"`
	OCSPServers           []string `json:"ocsp_servers"`
	EnableTemplating      bool     `json:"enable_templating"`
}

type ClusterConfigEntry struct {
	Path    string `json:"path"`
	AIAPath string `json:"aia_path"`
}

func GetAIAURLs(ctx context.Context, s logical.Storage, i *IssuerEntry) (*certutil.URLEntries, error) {
	// Default to the per-issuer AIA URLs.
	entries := i.AIAURIs

	// If none are set (either due to a nil entry or because no URLs have
	// been provided), fall back to the global AIA URL config.
	if entries == nil || (len(entries.IssuingCertificates) == 0 && len(entries.CRLDistributionPoints) == 0 && len(entries.OCSPServers) == 0) {
		var err error

		entries, err = GetGlobalAIAURLs(ctx, s)
		if err != nil {
			return nil, err
		}
	}

	if entries == nil {
		return &certutil.URLEntries{}, nil
	}

	return ToURLEntries(ctx, s, i.ID, entries)
}

func GetGlobalAIAURLs(ctx context.Context, storage logical.Storage) (*AiaConfigEntry, error) {
	entry, err := storage.Get(ctx, "urls")
	if err != nil {
		return nil, err
	}

	entries := &AiaConfigEntry{
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

func ToURLEntries(ctx context.Context, s logical.Storage, issuer IssuerID, c *AiaConfigEntry) (*certutil.URLEntries, error) {
	if len(c.IssuingCertificates) == 0 && len(c.CRLDistributionPoints) == 0 && len(c.OCSPServers) == 0 {
		return &certutil.URLEntries{}, nil
	}

	result := certutil.URLEntries{
		IssuingCertificates:   c.IssuingCertificates[:],
		CRLDistributionPoints: c.CRLDistributionPoints[:],
		OCSPServers:           c.OCSPServers[:],
	}

	if c.EnableTemplating {
		cfg, err := GetClusterConfig(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("error fetching cluster-local address config: %w", err)
		}

		for name, source := range map[string]*[]string{
			"issuing_certificates":    &result.IssuingCertificates,
			"crl_distribution_points": &result.CRLDistributionPoints,
			"ocsp_servers":            &result.OCSPServers,
		} {
			templated := make([]string, len(*source))
			for index, uri := range *source {
				if strings.Contains(uri, "{{cluster_path}}") && len(cfg.Path) == 0 {
					return nil, fmt.Errorf("unable to template AIA URLs as we lack local cluster address information (path)")
				}
				if strings.Contains(uri, "{{cluster_aia_path}}") && len(cfg.AIAPath) == 0 {
					return nil, fmt.Errorf("unable to template AIA URLs as we lack local cluster address information (aia_path)")
				}
				if strings.Contains(uri, "{{issuer_id}}") && len(issuer) == 0 {
					// Elide issuer AIA info as we lack an issuer_id.
					return nil, fmt.Errorf("unable to template AIA URLs as we lack an issuer_id for this operation")
				}

				uri = strings.ReplaceAll(uri, "{{cluster_path}}", cfg.Path)
				uri = strings.ReplaceAll(uri, "{{cluster_aia_path}}", cfg.AIAPath)
				uri = strings.ReplaceAll(uri, "{{issuer_id}}", issuer.String())
				templated[index] = uri
			}

			if uri := ValidateURLs(templated); uri != "" {
				return nil, fmt.Errorf("error validating templated %v; invalid URI: %v", name, uri)
			}

			*source = templated
		}
	}

	return &result, nil
}

func GetClusterConfig(ctx context.Context, s logical.Storage) (*ClusterConfigEntry, error) {
	entry, err := s.Get(ctx, ClusterConfigPath)
	if err != nil {
		return nil, err
	}

	var result ClusterConfigEntry
	if entry == nil {
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func ValidateURLs(urls []string) string {
	for _, curr := range urls {
		if !govalidator.IsURL(curr) || strings.Contains(curr, "{{issuer_id}}") || strings.Contains(curr, "{{cluster_path}}") || strings.Contains(curr, "{{cluster_aia_path}}") {
			return curr
		}
	}

	return ""
}

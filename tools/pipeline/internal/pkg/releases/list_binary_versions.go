// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ListBinaryVersionsReq struct {
	VersionsString string
}

type VariantInfo struct {
	Variant string   `json:"variant"`
	OS      []string `json:"os,omitempty"`
}

// - ValidVersions: details for versions that exist
// - InvalidVersions: versions requested by the user but not found upstream
// - AllVersions: original input order
type ListBinaryVersionsRes struct {
	ValidVersions map[string]struct {
		Status   string        `json:"status"`
		Variants []VariantInfo `json:"variants"`
	} `json:"valid_versions"`

	InvalidVersions []string `json:"invalid_versions"`
	AllVersions     []string `json:"all_versions"`
}

func NewListBinaryVersionsReq(s string) *ListBinaryVersionsReq {
	return &ListBinaryVersionsReq{VersionsString: s}
}

// normalizeLabel transforms user-friendly labels into a canonical form
// that matches HashiCorp’s release naming scheme.
func normalizeLabel(label string) (normalized, display string) {
	display = label
	normalized = strings.TrimSuffix(label, "-ce")
	normalized = strings.TrimSuffix(normalized, "-lts")
	if strings.HasSuffix(label, "-ent") {
		normalized = strings.TrimSuffix(normalized, "-ent") + "+ent"
	}
	return normalized, display
}

const baseURL = "https://releases.hashicorp.com/vault/"

// fetchAvailableVersions scrapes the Vault releases index page and returns
func fetchAvailableVersions(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Parse HTML returned by releases.hashicorp.com
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract version folder names from links ("/vault/1.18.1/")
	var versions []string
	doc.Find("body > ul > li > a").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Normalize: trim trailing "/", strip parent path
		href = strings.TrimSuffix(href, "/")
		if idx := strings.LastIndex(href, "/"); idx != -1 {
			href = href[idx+1:]
		}

		if href != "" {
			versions = append(versions, href)
		}
	})

	// Sort + deduplicate for stable deterministic results
	slices.Sort(versions)
	return slices.Compact(versions), nil
}

// fetchAvailableVariants returns all variant names that match a "base version".
func fetchAvailableVariants(available []string, version string) []string {
	var result []string
	prefix := version
	plusPrefix := version + "+" // enterprise or extra flavors

	for _, v := range available {
		if v == prefix || strings.HasPrefix(v, plusPrefix) {
			result = append(result, v)
		}
	}
	return result
}

// fetchAvailableOs inspects each variant directory and extracts the OS
func fetchAvailableOs(ctx context.Context, items []string) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, full := range items {
		result[full] = []string{}
		url := baseURL + full + "/"

		// Fetch the specific variant page
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status for %s: %d", full, resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		// ZIP filename format starts with: vault_<variant>_<os>_<arch>.zip
		prefix := "vault_" + full

		doc.Find("body > ul > li > a").Each(func(_ int, sel *goquery.Selection) {
			name := sel.Text()

			// Ignore unrelated files
			if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ".zip") {
				return
			}

			parts := strings.Split(strings.TrimSuffix(name, ".zip"), "_")
			if len(parts) < 3 {
				return // unexpected format
			}

			// OS is always the second-to-last component ("linux", "darwin", etc.)
			osName := parts[len(parts)-2]

			// Avoid duplicates
			if !slices.Contains(result[full], osName) {
				result[full] = append(result[full], osName)
			}
		})
	}

	return result, nil
}

func (r *ListBinaryVersionsReq) Run(ctx context.Context) (*ListBinaryVersionsRes, error) {
	if r == nil || r.VersionsString == "" {
		return &ListBinaryVersionsRes{}, nil
	}

	// User-provided labels (raw, with suffixes)
	labels := strings.Fields(r.VersionsString)

	// Initialize response container
	res := &ListBinaryVersionsRes{
		ValidVersions: make(map[string]struct {
			Status   string        `json:"status"`
			Variants []VariantInfo `json:"variants"`
		}),
		AllVersions: labels,
	}

	// Fetch the list of *all* versions from HashiCorp
	availableVersions, err := fetchAvailableVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed fetching versions: %w", err)
	}

	// Process each requested version label
	for _, label := range labels {
		normalized, display := normalizeLabel(label)

		// If the base version doesn't exist, the entire entry is invalid
		if !slices.Contains(availableVersions, normalized) {
			res.InvalidVersions = append(res.InvalidVersions, display)
			continue
		}

		// Find all variants associated with this version
		variants := fetchAvailableVariants(availableVersions, normalized)
		if len(variants) == 0 {
			res.InvalidVersions = append(res.InvalidVersions, display)
			continue
		}

		// Fetch OS availability for each variant
		osMap, err := fetchAvailableOs(ctx, variants)
		if err != nil {
			// Non-fatal: mark the version invalid but continue processing others
			slog.WarnContext(ctx, "failed fetching OS", "version", normalized, "error", err)
			res.InvalidVersions = append(res.InvalidVersions, display)
			continue
		}

		// Build variant list
		slices.Sort(variants)
		var vlist []VariantInfo
		for _, v := range variants {
			vlist = append(vlist, VariantInfo{Variant: v, OS: osMap[v]})
		}

		// Add a valid entry for this version label
		res.ValidVersions[display] = struct {
			Status   string        `json:"status"`
			Variants []VariantInfo `json:"variants"`
		}{
			Status:   "valid",
			Variants: vlist,
		}
	}

	slices.Sort(res.InvalidVersions)
	return res, nil
}

// ToJSON returns a pretty-printed JSON representation of the results.
func (r *ListBinaryVersionsRes) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// String provides a concise human-readable summary (counts only).
func (r *ListBinaryVersionsRes) String() string {
	total := 0
	for _, v := range r.ValidVersions {
		total += len(v.Variants)
	}
	return fmt.Sprintf(
		"Listed %d → %d valid (%d variants), %d missing",
		len(r.AllVersions), len(r.ValidVersions), total, len(r.InvalidVersions),
	)
}

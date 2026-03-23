// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	libgitclient "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/zclconf/go-cty/cty"
)

// ListUpdatedVersionsReq represents the request for updating versions.
type ListUpdatedVersionsReq struct {
	VersionsDecodeRes *DecodeRes `json:"versions_decode_res,omitempty"`
}

// ListUpdatedVersionsRes is the response for updating versions.
type ListUpdatedVersionsRes struct {
	Versions string `json:"versions,omitempty"`
}

// Run reads the existing versions.hcl file, updates versions from input,
// applies retention rules, and returns formatted HCL.
func (r *ListUpdatedVersionsReq) Run(ctx context.Context, git *libgitclient.Client, args []string) (*ListUpdatedVersionsRes, error) {
	if err := r.VersionsDecodeRes.Validate(ctx); err != nil {
		return nil, err
	}

	if r.VersionsDecodeRes.Config == nil {
		return nil, errors.New("no versions.hcl information was provided")
	}

	config := r.VersionsDecodeRes.Config
	inputVersions := strings.Fields(strings.Join(args, " "))
	versionRegex := regexp.MustCompile(`^(\d+\.\d+)`)

	// Parse and update versions from input
	updatedKeys := make([]string, 0, len(inputVersions))
	for _, versionStr := range inputVersions {
		matches := versionRegex.FindStringSubmatch(versionStr)
		if len(matches) < 2 {
			return nil, fmt.Errorf("invalid version format: %s", versionStr)
		}
		majorMinor := matches[1] + ".x"
		hasCE := strings.Contains(versionStr, "-ce")
		hasLTS := strings.Contains(versionStr, "-lts")

		if entry, exists := config.ActiveVersion.Versions[majorMinor]; exists {
			entry.CEActive = hasCE
			entry.LTS = hasLTS
		} else {
			config.ActiveVersion.Versions[majorMinor] = &Version{CEActive: hasCE, LTS: hasLTS}
		}
		updatedKeys = append(updatedKeys, majorMinor)
	}

	// Apply new retention and CE activation policy
	filtered := r.applyRetentionPolicy(config.ActiveVersion.Versions, updatedKeys)

	// Generate HCL output
	hf := hclwrite.NewEmptyFile()
	root := hf.Body()
	root.SetAttributeValue("schema", cty.NumberIntVal(int64(config.Schema)))
	root.AppendNewline()

	activeBlock := root.AppendNewBlock("active_versions", nil).Body()

	// Sort versions descending
	versions := make([]string, 0, len(filtered))
	for v := range filtered {
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	for i, version := range versions {
		v := filtered[version]
		block := activeBlock.AppendNewBlock("version", []string{version}).Body()
		block.SetAttributeValue("ce_active", cty.BoolVal(v.CEActive))
		if v.LTS {
			block.SetAttributeValue("lts", cty.BoolVal(true))
		}
		if i < len(versions)-1 {
			activeBlock.AppendNewline()
		}
	}

	content := string(hclwrite.Format(hf.Bytes()))
	return &ListUpdatedVersionsRes{Versions: content}, nil
}

// applyRetentionPolicy implements:
// - n (latest): ce_active = true
// - n-2 and below: ce_active = false
// - n-3 and below: removed unless lts
func (r *ListUpdatedVersionsReq) applyRetentionPolicy(allVersions map[string]*Version, newInputVersions []string) map[string]*Version {
	if len(allVersions) == 0 {
		return make(map[string]*Version)
	}

	// Extract and sort version keys in ascending order
	keys := make([]string, 0, len(allVersions))
	for k := range allVersions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Determine latest version (n): prefer input, fallback to highest existing
	latest := keys[len(keys)-1] // default
	if len(newInputVersions) > 0 {
		sort.Strings(newInputVersions)
		candidate := newInputVersions[len(newInputVersions)-1]
		if idx := sort.SearchStrings(keys, candidate); idx < len(keys) && keys[idx] == candidate {
			latest = candidate
		}
	}

	latestIdx := sort.SearchStrings(keys, latest)
	if latestIdx >= len(keys) || keys[latestIdx] != latest {
		latestIdx = len(keys) - 1
	}

	// Calculate n-1 and n-2 indices (clamp to 0)
	nMinus1Idx := max(latestIdx-1, 0)
	nMinus2Idx := max(latestIdx-2, 0)

	filtered := make(map[string]*Version)

	for i, version := range keys {
		info := allVersions[version]
		clone := &Version{
			CEActive: false,
			LTS:      info.LTS,
		}

		// n: ce_active = true
		if i == latestIdx {
			clone.CEActive = true
			filtered[version] = clone
			continue
		}

		// n-1 and n-2: always keep
		if i == nMinus1Idx || i == nMinus2Idx {
			filtered[version] = clone
			continue
		}

		// n-3 and below: keep only if LTS
		if i < nMinus2Idx && info.LTS {
			filtered[version] = clone
		}
	}

	return filtered
}

// max returns the greater of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

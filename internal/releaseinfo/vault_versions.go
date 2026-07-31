// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releaseinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/helper/useragent"
)

const (
	defaultVaultReleasesURL = "https://releases.hashicorp.com/vault/index.json"
	vaultVersionsCacheTTL   = 1 * time.Hour
)

var (
	vaultReleasesURL = defaultVaultReleasesURL

	vaultVersionsCacheMu     sync.Mutex
	cachedRawVersionKeys     []string
	vaultVersionsCacheExpiry time.Time
)

// FetchVaultVersions retrieves available Vault versions from releases.hashicorp.com,
// filtering by edition and deduplicating to MAJOR.MINOR.PATCH strings.
//
// edition must be either "enterprise" or "community":
//   - "enterprise" returns versions whose build metadata is exactly "+ent"
//     (e.g. "1.19.15+ent") — HSM/FIPS variants such as "+ent.hsm" or
//     "+ent.fips1402" are excluded to avoid duplicates.
//   - "community" returns versions that carry no build metadata at all
//     (e.g. "1.13.9").
//
// Pre-release labels (e.g. "-rc1") are stripped before deduplication so that
// only GA patch versions are presented.
func FetchVaultVersions(ctx context.Context, edition string) ([]string, error) {
	raw, err := fetchRawVersionKeys(ctx)
	if err != nil {
		return nil, err
	}
	return filterVersions(raw, edition), nil
}

// fetchRawVersionKeys returns the full list of version keys from
// releases.hashicorp.com, serving a cached result within the TTL.
func fetchRawVersionKeys(ctx context.Context) ([]string, error) {
	vaultVersionsCacheMu.Lock()
	if cachedRawVersionKeys != nil && time.Now().Before(vaultVersionsCacheExpiry) {
		result := cachedRawVersionKeys
		vaultVersionsCacheMu.Unlock()
		return result, nil
	}
	vaultVersionsCacheMu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, vaultReleasesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating vault versions request: %w", err)
	}
	req.Header.Set("User-Agent", useragent.String())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching vault versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from releases.hashicorp.com", resp.StatusCode)
	}

	var payload struct {
		Versions map[string]json.RawMessage `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decoding vault versions response: %w", err)
	}

	keys := make([]string, 0, len(payload.Versions))
	for v := range payload.Versions {
		keys = append(keys, v)
	}

	vaultVersionsCacheMu.Lock()
	cachedRawVersionKeys = keys
	vaultVersionsCacheExpiry = time.Now().Add(vaultVersionsCacheTTL)
	vaultVersionsCacheMu.Unlock()

	return keys, nil
}

// filterVersions selects version keys matching the requested edition, strips
// pre-release labels, deduplicates, and returns a sorted slice.
func filterVersions(keys []string, edition string) []string {
	seen := make(map[string]struct{})

	for _, key := range keys {
		// Split on '+' to isolate the build-metadata portion.
		corePlusBuild := strings.SplitN(key, "+", 2)
		core := corePlusBuild[0] // e.g. "1.19.15" or "1.14.0-rc1"
		build := ""
		if len(corePlusBuild) == 2 {
			build = corePlusBuild[1] // e.g. "ent", "ent.hsm", "ent.fips1402"
		}

		switch edition {
		case "enterprise":
			// Accept only exactly "+ent" — no HSM, FIPS, or other variants.
			if build != "ent" {
				continue
			}
		default: // "community"
			// Accept only versions with no build metadata at all.
			if build != "" {
				continue
			}
		}

		// Strip pre-release label (everything from '-' onward in the core).
		// "1.14.0-rc1" → "1.14.0"
		if idx := strings.Index(core, "-"); idx != -1 {
			core = core[:idx]
		}

		seen[core] = struct{}{}
	}

	versions := make([]string, 0, len(seen))
	for v := range seen {
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool {
		vi, erri := goversion.NewVersion(versions[i])
		vj, errj := goversion.NewVersion(versions[j])
		if erri != nil || errj != nil {
			return versions[i] < versions[j]
		}
		return vi.LessThan(vj)
	})
	return versions
}

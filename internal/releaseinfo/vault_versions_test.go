// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package releaseinfo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockReleasesServer returns an httptest.Server that serves the given version
// keys as a releases.hashicorp.com-style index.json payload.
func mockReleasesServer(t *testing.T, statusCode int, versionKeys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}
		payload := struct {
			Versions map[string]json.RawMessage `json:"versions"`
		}{
			Versions: make(map[string]json.RawMessage, len(versionKeys)),
		}
		for _, v := range versionKeys {
			payload.Versions[v] = json.RawMessage(`{}`)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

// resetVaultVersionsCache wipes the package-level cache so each test starts clean.
func resetVaultVersionsCache(t *testing.T) {
	t.Helper()
	vaultVersionsCacheMu.Lock()
	cachedRawVersionKeys = nil
	vaultVersionsCacheExpiry = time.Time{}
	vaultVersionsCacheMu.Unlock()
}

// TestFilterVersions_Enterprise verifies that only exact "+ent" versions are
// returned for the enterprise edition, deduplicated and sorted numerically.
func TestFilterVersions_Enterprise(t *testing.T) {
	keys := []string{
		"1.9.7+ent.hsm",               // excluded: not exactly +ent
		"2.0.2+ent",                   // included
		"1.13.4+ent.fips1402",         // excluded: not exactly +ent
		"1.13.9",                      // excluded: community
		"1.14.0-rc1+ent.hsm.fips1402", // excluded
		"1.14.5",                      // excluded: community
		"1.16.20+ent.hsm.fips1402",    // excluded
		"1.16.26+ent",                 // included
		"1.18.6+ent.fips1402",         // excluded
		"1.14.4+ent.hsm.fips1402",     // excluded
		"1.14.9",                      // excluded: community
		"1.19.15+ent",                 // included
		"1.8.10",                      // excluded: community
		"1.2.3+ent",                   // included — tests numeric sort vs 1.11.x
		"1.11.0+ent",                  // included
	}

	got := filterVersions(keys, "enterprise")

	want := []string{"1.2.3", "1.11.0", "1.16.26", "1.19.15", "2.0.2"}
	if len(got) != len(want) {
		t.Fatalf("enterprise: got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("enterprise[%d]: got %q, want %q (full result: %v)", i, got[i], v, got)
		}
	}
}

// TestFilterVersions_Community verifies that only plain versions (no build
// metadata) are returned for the community edition.
func TestFilterVersions_Community(t *testing.T) {
	keys := []string{
		"1.9.7+ent.hsm",
		"2.0.2+ent",
		"1.13.9",
		"1.14.5",
		"1.14.9",
		"1.8.10",
		"1.2.3",
		"1.11.0",
	}

	got := filterVersions(keys, "community")

	want := []string{"1.2.3", "1.8.10", "1.11.0", "1.13.9", "1.14.5", "1.14.9"}
	if len(got) != len(want) {
		t.Fatalf("community: got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("community[%d]: got %q, want %q (full result: %v)", i, got[i], v, got)
		}
	}
}

// TestFilterVersions_UnknownEditionDefaultsToCommunity verifies that any
// unrecognised edition string is treated as community.
func TestFilterVersions_UnknownEditionDefaultsToCommunity(t *testing.T) {
	keys := []string{"1.0.0", "1.0.0+ent", "1.0.0+ent.hsm"}

	got := filterVersions(keys, "potato")

	if len(got) != 1 || got[0] != "1.0.0" {
		t.Errorf("unknown edition: got %v, want [1.0.0]", got)
	}
}

// TestFilterVersions_PreReleaseStripped verifies that pre-release labels are
// stripped and the resulting GA version is deduplicated correctly.
func TestFilterVersions_PreReleaseStripped(t *testing.T) {
	keys := []string{
		"1.14.0-rc1+ent", // pre-release → stripped to 1.14.0+ent → kept once
		"1.14.0+ent",     // GA — same core after stripping, deduplicates
		"1.14.1+ent",
	}

	got := filterVersions(keys, "enterprise")

	want := []string{"1.14.0", "1.14.1"}
	if len(got) != len(want) {
		t.Fatalf("pre-release strip: got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("pre-release strip[%d]: got %q, want %q", i, got[i], v)
		}
	}
}

// TestFilterVersions_NumericSort verifies that version segments are compared
// numerically, not lexicographically (1.2 < 1.11, not 1.11 < 1.2).
func TestFilterVersions_NumericSort(t *testing.T) {
	keys := []string{
		"1.12.0", "1.1.0", "1.11.0", "1.2.0", "1.9.0", "1.10.0",
	}

	got := filterVersions(keys, "community")

	want := []string{"1.1.0", "1.2.0", "1.9.0", "1.10.0", "1.11.0", "1.12.0"}
	if len(got) != len(want) {
		t.Fatalf("numeric sort: got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("numeric sort[%d]: got %q, want %q (full: %v)", i, got[i], v, got)
		}
	}
}

// TestFilterVersions_Empty verifies an empty input returns an empty slice.
func TestFilterVersions_Empty(t *testing.T) {
	got := filterVersions([]string{}, "enterprise")
	if len(got) != 0 {
		t.Errorf("empty input: got %v, want []", got)
	}
}

// TestFetchVaultVersions_HTTP tests the full FetchVaultVersions path against a
// local httptest server, exercising the HTTP fetch and filter together.
func TestFetchVaultVersions_HTTP(t *testing.T) {
	srv := mockReleasesServer(t, http.StatusOK, []string{
		"1.19.15+ent",
		"1.19.15",
		"1.20.0+ent",
		"1.20.0",
		"1.20.0+ent.hsm",
	})
	defer srv.Close()

	// Override the package-level URL and reset cache.
	original := vaultReleasesURL
	vaultReleasesURL = srv.URL
	defer func() { vaultReleasesURL = original }()
	resetVaultVersionsCache(t)
	defer resetVaultVersionsCache(t)

	t.Run("enterprise", func(t *testing.T) {
		resetVaultVersionsCache(t)
		got, err := FetchVaultVersions(context.Background(), "enterprise")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"1.19.15", "1.20.0"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i, v := range want {
			if got[i] != v {
				t.Errorf("[%d] got %q, want %q", i, got[i], v)
			}
		}
	})

	t.Run("community", func(t *testing.T) {
		resetVaultVersionsCache(t)
		got, err := FetchVaultVersions(context.Background(), "community")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"1.19.15", "1.20.0"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i, v := range want {
			if got[i] != v {
				t.Errorf("[%d] got %q, want %q", i, got[i], v)
			}
		}
	})
}

// TestFetchVaultVersions_NonOKStatus verifies that a non-200 response returns an error.
func TestFetchVaultVersions_NonOKStatus(t *testing.T) {
	srv := mockReleasesServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	original := vaultReleasesURL
	vaultReleasesURL = srv.URL
	defer func() { vaultReleasesURL = original }()
	resetVaultVersionsCache(t)
	defer resetVaultVersionsCache(t)

	_, err := FetchVaultVersions(context.Background(), "enterprise")
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

// TestFetchVaultVersions_Cache verifies that a second call within the TTL does
// not make a second HTTP request (server request count stays at 1).
func TestFetchVaultVersions_Cache(t *testing.T) {
	requestCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		payload := `{"versions":{"1.19.0+ent":{},"1.20.0+ent":{}}}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
	}))
	defer srv.Close()

	original := vaultReleasesURL
	vaultReleasesURL = srv.URL
	defer func() { vaultReleasesURL = original }()
	resetVaultVersionsCache(t)
	defer resetVaultVersionsCache(t)

	if _, err := FetchVaultVersions(context.Background(), "enterprise"); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if _, err := FetchVaultVersions(context.Background(), "enterprise"); err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	if requestCount != 1 {
		t.Errorf("expected 1 HTTP request (cache hit on second call), got %d", requestCount)
	}
}
